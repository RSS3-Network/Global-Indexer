package txmgr

import (
	"fmt"
	"time"
)

// Config houses parameters for altering the behavior of a SimpleTxManager.
type Config struct {
	// ResubmissionTimeout is the interval at which, if no previously
	// published transaction has been mined, the new tx with a bumped gas
	// price will be published. Only one publication at MaxGasPrice will be
	// attempted.
	ResubmissionTimeout time.Duration

	// The multiplier applied to fee suggestions to put a hard limit on fee increases.
	FeeLimitMultiplier uint64

	// TxSendTimeout is how long to wait for sending a transaction.
	// By default it is unbounded. If set, this is recommended to be at least 20 minutes.
	TxSendTimeout time.Duration

	// TxNotInMempoolTimeout is how long to wait before aborting a transaction send if the transaction does not
	// make it to the mempool. If the tx is in the mempool, TxSendTimeout is used instead.
	TxNotInMempoolTimeout time.Duration

	// NetworkTimeout is the allowed duration for a single network request.
	// This is intended to be used for network requests that can be replayed.
	NetworkTimeout time.Duration

	// RequireQueryInterval is the interval at which the tx manager will
	// query the backend to check for confirmations after a tx at a
	// specific gas price has been published.
	ReceiptQueryInterval time.Duration

	// NumConfirmations specifies how many blocks are need to consider a
	// transaction confirmed.
	NumConfirmations uint64

	// SafeAbortNonceTooLowCount specifies how many ErrNonceTooLow observations
	// are required to give up on a tx at a particular nonce without receiving
	// confirmation.
	SafeAbortNonceTooLowCount uint64
}

func (m Config) Check() error {
	if m.NumConfirmations == 0 {
		return fmt.Errorf("NumConfirmations must not be 0")
	}

	if m.NetworkTimeout == 0 {
		return fmt.Errorf("must provide NetworkTimeout")
	}

	if m.FeeLimitMultiplier == 0 {
		return fmt.Errorf("must provide FeeLimitMultiplier")
	}

	if m.ResubmissionTimeout == 0 {
		return fmt.Errorf("must provide ResubmissionTimeout")
	}

	if m.ReceiptQueryInterval == 0 {
		return fmt.Errorf("must provide ReceiptQueryInterval")
	}

	if m.TxNotInMempoolTimeout == 0 {
		return fmt.Errorf("must provide TxNotInMempoolTimeout")
	}

	if m.SafeAbortNonceTooLowCount == 0 {
		return fmt.Errorf("SafeAbortNonceTooLowCount must not be 0")
	}

	return nil
}
