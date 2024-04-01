package txmgr

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gicrypto "github.com/naturalselectionlabs/rss3-global-indexer/common/crypto"
	"go.uber.org/zap"
)

const (
	// Geth requires a minimum fee bump of 10% for tx resubmission
	priceBump int64 = 10
)

// new = old * (100 + priceBump) / 100
var priceBumpPercent = big.NewInt(100 + priceBump)
var oneHundred = big.NewInt(100)

type TxManager interface {
	Send(ctx context.Context, candidate TxCandidate) (*types.Receipt, error)
}

type SimpleTxManager struct {
	cfg Config

	chainID        *big.Int
	ethereumClient *ethclient.Client
	from           common.Address
	nonce          *uint64
	nonceLock      sync.RWMutex

	signer gicrypto.SignerFn
}

type TxCandidate struct {
	// TxData is the transaction data to be used in the constructed tx.
	TxData []byte
	// To is the recipient of the constructed tx. Nil means contract creation.
	To *common.Address
	// GasLimit is the gas limit to be used in the constructed tx.
	GasLimit uint64
	// Value is the value to be used in the constructed tx.
	Value *big.Int
}

func (m *SimpleTxManager) Send(ctx context.Context, candidate TxCandidate) (*types.Receipt, error) {
	receipt, err := m.send(ctx, candidate)
	if err != nil {
		m.resetNonce()
	}

	return receipt, err
}

func (m *SimpleTxManager) resetNonce() {
	m.nonceLock.Lock()
	defer m.nonceLock.Unlock()
	m.nonce = nil
}

// send performs the actual transaction creation and sending.
func (m *SimpleTxManager) send(ctx context.Context, candidate TxCandidate) (*types.Receipt, error) {
	if m.cfg.TxSendTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.cfg.TxSendTimeout)

		defer cancel()
	}

	var (
		tx *types.Transaction

		err error
	)

	if tx, err = retry.DoWithData(func() (*types.Transaction, error) {
		tx, err = m.craftTx(ctx, candidate)
		if err != nil {
			zap.L().Warn("Failed to create a transaction, will retry", zap.Error(err))

			return nil, err
		}

		return tx, nil
	}, retry.Delay(2*time.Second), retry.Attempts(30)); err != nil {
		return nil, fmt.Errorf("failed to create the tx: %w", err)
	}

	return m.sendTx(ctx, tx)
}

func (m *SimpleTxManager) craftTx(ctx context.Context, candidate TxCandidate) (*types.Transaction, error) {
	gasTipCap, basefee, err := m.suggestGasPriceCaps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price info: %w", err)
	}

	gasFeeCap := calcGasFeeCap(basefee, gasTipCap)

	rawTx := &types.DynamicFeeTx{
		ChainID:   m.chainID,
		To:        candidate.To,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Data:      candidate.TxData,
		Value:     candidate.Value,
	}

	// If the gas limit is set, we can use that as the gas
	if candidate.GasLimit != 0 {
		rawTx.Gas = candidate.GasLimit
	} else {
		// Calculate the intrinsic gas for the transaction
		gas, err := m.ethereumClient.EstimateGas(ctx, ethereum.CallMsg{
			From:      m.from,
			To:        candidate.To,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Data:      rawTx.Data,
			Value:     rawTx.Value,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}

		rawTx.Gas = gas
	}

	return m.signWithNextNonce(ctx, rawTx)
}

func (m *SimpleTxManager) sendTx(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	sendState := NewSendState(m.cfg.SafeAbortNonceTooLowCount, m.cfg.TxNotInMempoolTimeout)
	receiptChan := make(chan *types.Receipt, 1)
	publishAndWait := func(tx *types.Transaction, bumpFees bool) *types.Transaction {
		wg.Add(1)

		tx, published := m.publishTx(ctx, tx, sendState, bumpFees)

		if published {
			go func() {
				defer wg.Done()
				m.waitForTx(ctx, tx, sendState, receiptChan)
			}()
		} else {
			wg.Done()
		}

		return tx
	}

	// Immediately publish a transaction before starting the resubmission loop
	tx = publishAndWait(tx, false)

	ticker := time.NewTicker(m.cfg.ResubmissionTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Don't resubmit a transaction if it has been mined, but we are waiting for the conf depth.
			if sendState.IsWaitingForConfirmation() {
				continue
			}
			// If we see lots of unrecoverable errors (and no pending transactions) abort sending the transaction.
			if sendState.ShouldAbortImmediately() {
				zap.L().Error("aborting transaction submission")

				return nil, errors.New("aborted transaction sending")
			}

			tx = publishAndWait(tx, true)

		case <-ctx.Done():
			return nil, ctx.Err()

		case receipt := <-receiptChan:
			return receipt, nil
		}
	}
}

func (m *SimpleTxManager) suggestGasPriceCaps(ctx context.Context) (*big.Int, *big.Int, error) {
	cCtx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)
	defer cancel()

	tip, err := m.ethereumClient.SuggestGasTipCap(cCtx)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch the suggested gas tip cap: %w", err)
	} else if tip == nil {
		return nil, nil, fmt.Errorf("the suggested tip was nil")
	}

	cCtx, cancel = context.WithTimeout(ctx, m.cfg.NetworkTimeout)
	defer cancel()

	head, err := m.ethereumClient.HeaderByNumber(cCtx, nil)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch the suggested basefee: %w", err)
	} else if head.BaseFee == nil {
		return nil, nil, fmt.Errorf("txmgr does not support pre-london blocks that do not have a basefee")
	}

	return tip, head.BaseFee, nil
}

func calcGasFeeCap(baseFee, gasTipCap *big.Int) *big.Int {
	return new(big.Int).Add(
		gasTipCap,
		new(big.Int).Mul(baseFee, big.NewInt(2)),
	)
}

func (m *SimpleTxManager) signWithNextNonce(ctx context.Context, rawTx *types.DynamicFeeTx) (*types.Transaction, error) {
	m.nonceLock.Lock()
	defer m.nonceLock.Unlock()

	if m.nonce == nil {
		// Fetch the sender's nonce from the latest known block (nil `blockNumber`)
		childCtx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)
		defer cancel()

		nonce, err := m.ethereumClient.NonceAt(childCtx, m.from, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}

		m.nonce = &nonce
	} else {
		*m.nonce++
	}

	rawTx.Nonce = *m.nonce
	ctx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)

	defer cancel()

	tx, err := m.signer(ctx, m.from, types.NewTx(rawTx))

	if err != nil {
		// decrement the nonce, so we can retry signing with the same nonce next time
		// signWithNextNonce is called
		*m.nonce--
	}

	return tx, err
}

func (m *SimpleTxManager) publishTx(ctx context.Context, tx *types.Transaction, sendState *SendState, bumpFeesImmediately bool) (*types.Transaction, bool) {
	for {
		if bumpFeesImmediately {
			newTx, err := m.increaseGasPrice(ctx, tx)
			if err != nil {
				return tx, false
			}

			tx = newTx
			sendState.bumpCount++
		}

		bumpFeesImmediately = true // bump fees next loop

		if sendState.IsWaitingForConfirmation() {
			// there is a chance the previous tx goes into "waiting for confirmation" state
			// during the increaseGasPrice call; continue waiting rather than resubmit the tx
			return tx, false
		}

		cCtx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)
		err := m.ethereumClient.SendTransaction(cCtx, tx)

		cancel()
		sendState.ProcessSendError(err)

		if err == nil {
			return tx, true
		}

		zap.L().Error("sending transaction error", zap.Error(err), zap.String("hash", tx.Hash().String()))

		switch {
		case errStringMatch(err, core.ErrNonceTooLow):
			zap.L().Warn("nonce too low", zap.Error(err))
		case errStringMatch(err, context.Canceled):
			zap.L().Warn("transaction send cancelled", zap.Error(err))
		case errStringMatch(err, txpool.ErrAlreadyKnown):
			zap.L().Warn("resubmitted already known transaction", zap.Error(err))
		case errStringMatch(err, txpool.ErrReplaceUnderpriced):
			zap.L().Warn("transaction replacement is underpriced", zap.Error(err))
			continue // retry with fee bump
		case errStringMatch(err, txpool.ErrUnderpriced):
			zap.L().Warn("transaction is underpriced", zap.Error(err))
			continue // retry with fee bump
		default:
			zap.L().Error("unable to publish transaction", zap.Error(err))
		}

		// on non-underpriced error return immediately; will retry on next resubmission timeout
		return tx, false
	}
}

func errStringMatch(err, target error) bool {
	if err == nil && target == nil {
		return true
	} else if err == nil || target == nil {
		return false
	}

	return strings.Contains(err.Error(), target.Error())
}

func (m *SimpleTxManager) waitForTx(ctx context.Context, tx *types.Transaction, sendState *SendState, receiptChan chan *types.Receipt) {
	t := time.Now()
	// Poll for the transaction to be ready & then send the result to receiptChan
	receipt, err := m.waitMined(ctx, tx, sendState)
	if err != nil {
		// this will happen if the tx was successfully replaced by a tx with bumped fees
		zap.L().Info("Transaction receipt not found", zap.Error(err), zap.String("hash", tx.Hash().String()))
		return
	}
	select {
	case receiptChan <- receipt:
		zap.L().Info("Transaction receipt return", zap.Any("cost", time.Since(t).Milliseconds()))
	default:
	}
}

func (m *SimpleTxManager) waitMined(ctx context.Context, tx *types.Transaction, sendState *SendState) (*types.Receipt, error) {
	txHash := tx.Hash()
	queryTicker := time.NewTicker(m.cfg.ReceiptQueryInterval)

	defer queryTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
			if receipt := m.queryReceipt(ctx, txHash, sendState); receipt != nil {
				return receipt, nil
			}
		}
	}
}

func (m *SimpleTxManager) queryReceipt(ctx context.Context, txHash common.Hash, sendState *SendState) *types.Receipt {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)
	defer cancel()

	receipt, err := m.ethereumClient.TransactionReceipt(ctx, txHash)

	if errors.Is(err, ethereum.NotFound) {
		sendState.TxNotMined(txHash)
		zap.L().Info("Transaction not yet mined", zap.Any("hash", txHash.String()))

		return nil
	} else if err != nil {
		zap.L().Error("Receipt retrieval failed", zap.Any("hash", txHash.String()), zap.Error(err))

		return nil
	} else if receipt == nil {
		zap.L().Warn("Receipt and error are both nil", zap.Any("hash", txHash.String()))

		return nil
	}

	// Receipt is confirmed to be valid from this point on
	sendState.TxMined(txHash)

	txHeight := receipt.BlockNumber.Uint64()
	tipHeight, err := m.ethereumClient.BlockNumber(ctx)

	if err != nil {
		return nil
	}

	// The transaction is considered confirmed when
	// txHeight+numConfirmations-1 <= tipHeight. Note that the -1 is
	// needed to account for the fact that confirmations have an
	// inherent off-by-one, i.e. when using 1 confirmation the
	// transaction should be confirmed when txHeight is equal to
	// tipHeight. The equation is rewritten in this form to avoid
	// underflows.
	if txHeight+m.cfg.NumConfirmations <= tipHeight+1 {
		zap.L().Info("Transaction confirmed", zap.Any("hash", txHash.String()))
		return receipt
	}

	// Safe to subtract since we know the LHS above is greater.
	confsRemaining := (txHeight + m.cfg.NumConfirmations) - (tipHeight + 1)

	zap.L().Debug("Transaction not yet confirmed", zap.Any("hash", txHash.String()), zap.Uint64("confsRemaining", confsRemaining))

	return nil
}

func updateFees(oldTip, oldFeeCap, newTip, newBaseFee *big.Int) (*big.Int, *big.Int) {
	newFeeCap := calcGasFeeCap(newBaseFee, newTip)

	thresholdTip := calcThresholdValue(oldTip)
	thresholdFeeCap := calcThresholdValue(oldFeeCap)

	if newTip.Cmp(thresholdTip) >= 0 && newFeeCap.Cmp(thresholdFeeCap) >= 0 {
		zap.L().Debug("Using new tip and feecap")
		return newTip, newFeeCap
	} else if newTip.Cmp(thresholdTip) >= 0 && newFeeCap.Cmp(thresholdFeeCap) < 0 {
		// Tip has gone up, but basefee is flat or down.
		// TODO(CLI-3714): Do we need to recalculate the FC here?
		zap.L().Debug("Using new tip and threshold feecap")
		return newTip, thresholdFeeCap
	} else if newTip.Cmp(thresholdTip) < 0 && newFeeCap.Cmp(thresholdFeeCap) >= 0 {
		// Basefee has gone up, but the tip hasn't. Recalculate the feecap because if the tip went up a lot
		// not enough of the feecap may be dedicated to paying the basefee.
		zap.L().Debug("Using threshold tip and recalculated feecap")

		return thresholdTip, calcGasFeeCap(newBaseFee, thresholdTip)
	}

	// TODO(CLI-3713): Should we skip the bump in this case?
	zap.L().Debug("Using threshold tip and threshold feecap")

	return thresholdTip, thresholdFeeCap
}

func calcThresholdValue(x *big.Int) *big.Int {
	threshold := new(big.Int).Mul(priceBumpPercent, x)
	threshold = threshold.Div(threshold, oneHundred)

	return threshold
}

func (m *SimpleTxManager) increaseGasPrice(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	zap.L().Info("bumping gas price for tx", zap.String("hash", tx.Hash().String()), zap.Uint64("tip", tx.GasTipCap().Uint64()), zap.Uint64("fee", tx.GasFeeCap().Uint64()), zap.Uint64("gaslimit", tx.Gas()))

	tip, basefee, err := m.suggestGasPriceCaps(ctx)

	if err != nil {
		zap.L().Warn("failed to get suggested gas tip and basefee", zap.Error(err))
		return nil, err
	}

	bumpedTip, bumpedFee := updateFees(tx.GasTipCap(), tx.GasFeeCap(), tip, basefee)

	// Make sure increase is at most [FeeLimitMultiplier] the suggested values
	maxTip := new(big.Int).Mul(tip, big.NewInt(int64(m.cfg.FeeLimitMultiplier)))
	if bumpedTip.Cmp(maxTip) > 0 {
		return nil, fmt.Errorf("bumped tip 0x%s is over %dx multiple of the suggested value", bumpedTip.Text(16), m.cfg.FeeLimitMultiplier)
	}

	maxFee := calcGasFeeCap(new(big.Int).Mul(basefee, big.NewInt(int64(m.cfg.FeeLimitMultiplier))), maxTip)

	if bumpedFee.Cmp(maxFee) > 0 {
		return nil, fmt.Errorf("bumped fee 0x%s is over %dx multiple of the suggested value", bumpedFee.Text(16), m.cfg.FeeLimitMultiplier)
	}

	rawTx := &types.DynamicFeeTx{
		ChainID:    tx.ChainId(),
		Nonce:      tx.Nonce(),
		GasTipCap:  bumpedTip,
		GasFeeCap:  bumpedFee,
		To:         tx.To(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}

	zap.L().Info("re-estimate gas", zap.String("hash", tx.Hash().String()), zap.Uint64("gasTipCap", bumpedTip.Uint64()), zap.Uint64("gasFeeCap", bumpedFee.Uint64()))

	rawTx.Gas = tx.Gas()

	ctx, cancel := context.WithTimeout(ctx, m.cfg.NetworkTimeout)
	defer cancel()

	newTx, err := m.signer(ctx, m.from, types.NewTx(rawTx))

	if err != nil {
		zap.L().Warn("failed to sign new transaction", zap.Error(err))
		return tx, nil
	}

	return newTx, nil
}

func NewSimpleTxManager(conf Config, chainID *big.Int, nonce *uint64, ethereumClient *ethclient.Client, from common.Address, singer gicrypto.SignerFn) (*SimpleTxManager, error) {
	if err := conf.Check(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &SimpleTxManager{
		cfg: conf,

		chainID:        chainID,
		ethereumClient: ethereumClient,
		from:           from,
		nonce:          nonce,

		signer: singer,
	}, nil
}
