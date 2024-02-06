package nameresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/character"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/lenshandle"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	goens "github.com/wealdtech/go-ens/v3"
	"go.uber.org/zap"
)

const (
	ErrUnregisterName = "unregistered name"
	ErrUnSupportName  = "unsupport name service resolution"
)

type NameResolver struct {
	ensEthClient       *ethclient.Client
	csbHandleContract  *character.Character
	lensHandleContract *lenshandle.LensHandle
	fcClient           *fcClient
}

type fcClient struct {
	endpointURL *url.URL
	httpClient  *http.Client
}

func (n *NameResolver) Resolve(ctx context.Context, input string) (string, error) {
	splits := strings.Split(input, ".")

	var (
		address string
		err     error
	)

	suffix := splits[len(splits)-1]

	switch {
	case suffix == NameServiceENS.String() && n.ensEthClient != nil:
		address, err = n.resolveENS(ctx, input)
	case suffix == NameServiceCSB.String() && n.csbHandleContract != nil:
		address, err = n.resolveCSB(ctx, input)
	case suffix == NameServiceLens.String() && n.lensHandleContract != nil:
		address, err = n.resolveLens(ctx, input)
	case suffix == NameServiceFarcaster.String() && n.fcClient != nil:
		address, err = n.resolveFarcaster(ctx, input)
	default:
		err = fmt.Errorf("%s:%s", ErrUnSupportName, input)
	}

	return address, err
}

func (n *NameResolver) resolveCSB(_ context.Context, domain string) (string, error) {
	cData, err := n.csbHandleContract.GetCharacterByHandle(&bind.CallOpts{}, strings.TrimSuffix(domain, ".csb"))
	if err != nil {
		return "", fmt.Errorf("failed to get crossbell character by handle: %w", err)
	}

	characterOwner, err := n.csbHandleContract.OwnerOf(&bind.CallOpts{}, cData.CharacterId)
	if err != nil {
		return "", fmt.Errorf("%s", ErrUnregisterName)
	}

	return characterOwner.String(), nil
}

func (n *NameResolver) resolveENS(_ context.Context, domain string) (string, error) {
	address, err := goens.Resolve(n.ensEthClient, domain)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func (n *NameResolver) resolveLens(_ context.Context, domain string) (string, error) {
	label := strings.Split(domain, "."+NameServiceLens.String())[0]
	tokenID, err := n.lensHandleContract.GetTokenId(&bind.CallOpts{}, label)

	if err != nil {
		return "", fmt.Errorf("failed to get lens token id by handle: %w", err)
	}

	owner, err := n.lensHandleContract.OwnerOf(&bind.CallOpts{}, tokenID)
	if err != nil {
		return "", fmt.Errorf("%s", ErrUnregisterName)
	}

	return owner.String(), nil
}

type UserNameProof struct {
	Timestamp uint32 `json:"timestamp"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
	Fid       uint64 `json:"fid"`
	Type      string `json:"type"`
}

func (n *NameResolver) resolveFarcaster(ctx context.Context, domain string) (string, error) {
	var (
		response UserNameProof
		err      error
	)

	fName := strings.Split(domain, "."+NameServiceFarcaster.String())[0]

	params := url.Values{}

	params.Add("name", fName)

	str := fmt.Sprintf("/v1/userNameProofByName?%s", params.Encode())

	onRetry := retry.OnRetry(func(n uint, err error) {
		zap.L().Error("fetch farcaster name", zap.Error(err), zap.Uint("attempts", n))
	})

	retryIf := retry.RetryIf(func(err error) bool {
		return err.Error() != fmt.Errorf("%s", ErrUnregisterName).Error()
	})

	if err = retry.Do(func() error { return n.call(ctx, str, &response) }, retry.Delay(time.Second), retry.Attempts(10), onRetry, retryIf); err != nil {
		return "", err
	}

	return response.Owner, nil
}

func (n *NameResolver) call(ctx context.Context, url string, result any) error {
	url = fmt.Sprintf("%s%s", n.fcClient.endpointURL, url)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	response, err := n.fcClient.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("%s", ErrUnregisterName)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func NewNameResolver(ctx context.Context, config *config.NameService) (*NameResolver, error) {
	var (
		err                error
		ensEthClient       *ethclient.Client
		characterContract  *character.Character
		lensHandleContract *lenshandle.LensHandle
		farcasterClient    *fcClient
	)

	if config.EnsEndpoint != "" {
		ensEthClient, err = ethclient.DialContext(ctx, config.EnsEndpoint)
		if err != nil {
			return nil, fmt.Errorf("dial ens ethereum client: %w", err)
		}
	}

	if config.CsbEndpoint != "" {
		csbEthClient, err := ethclient.DialContext(ctx, config.CsbEndpoint)
		if err != nil {
			return nil, fmt.Errorf("dial csb ethereum client: %w", err)
		}

		characterContract, err = character.NewCharacter(character.AddressCharacter, csbEthClient)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to crossbell character contract: %w", err)
		}
	}

	if config.LensEndpoint != "" {
		lensEthClient, err := ethclient.DialContext(ctx, config.LensEndpoint)
		if err != nil {
			return nil, fmt.Errorf("dial lens ethereum client: %w", err)
		}

		lensHandleContract, err = lenshandle.NewLensHandle(lenshandle.AddressLensHandle, lensEthClient)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to lens handle contract: %w", err)
		}
	}

	if config.FcEndpoint != "" {
		farcasterClient = &fcClient{
			httpClient: http.DefaultClient,
		}

		if farcasterClient.endpointURL, err = url.Parse(config.FcEndpoint); err != nil {
			return nil, fmt.Errorf("parse farcaster endpoint: %w", err)
		}
	}

	return &NameResolver{
		ensEthClient:       ensEthClient,
		csbHandleContract:  characterContract,
		lensHandleContract: lensHandleContract,
		fcClient:           farcasterClient,
	}, nil
}
