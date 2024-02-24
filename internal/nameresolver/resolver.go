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
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/crossbell"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/lens"
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
	csbHandleContract  *crossbell.Character
	lensHandleContract *lens.LensHandle
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

func NewNameResolver(ctx context.Context, config *config.RPCNetwork) (*NameResolver, error) {
	var (
		err                error
		ensEthClient       *ethclient.Client
		characterContract  *crossbell.Character
		lensHandleContract *lens.LensHandle
		farcasterClient    *fcClient
	)

	if config.Ethereum != nil {
		ensEthClient, err = ethclient.DialContext(ctx, config.Ethereum.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("dial ens ethereum client: %w", err)
		}
	}

	if config.Crossbell != nil {
		csbEthClient, err := ethclient.DialContext(ctx, config.Crossbell.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("dial csb ethereum client: %w", err)
		}

		characterContract, err = crossbell.NewCharacter(crossbell.AddressCharacter, csbEthClient)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to crossbell character contract: %w", err)
		}
	}

	if config.Polygon != nil {
		lensEthClient, err := ethclient.DialContext(ctx, config.Polygon.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("dial lens ethereum client: %w", err)
		}

		lensHandleContract, err = lens.NewLensHandle(lens.AddressLensHandle, lensEthClient)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to lens handle contract: %w", err)
		}
	}

	if config.Farcaster != nil {
		farcasterClient = &fcClient{}

		if farcasterClient.endpointURL, err = url.Parse(config.Farcaster.Endpoint); err != nil {
			return nil, fmt.Errorf("parse farcaster endpoint: %w", err)
		}

		var httpClient http.Client

		if config.Farcaster.APIkey != "" {
			httpClient.Transport = NewAuthenticationTransport(config.Farcaster.APIkey)
		} else {
			httpClient = *http.DefaultClient
		}

		farcasterClient.httpClient = &httpClient
	}

	return &NameResolver{
		ensEthClient:       ensEthClient,
		csbHandleContract:  characterContract,
		lensHandleContract: lensHandleContract,
		fcClient:           farcasterClient,
	}, nil
}

type AuthenticationTransport struct {
	APIKey string

	roundTripper http.RoundTripper
}

func (a *AuthenticationTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("api_key", a.APIKey)

	return a.roundTripper.RoundTrip(request)
}

func NewAuthenticationTransport(apiKey string) http.RoundTripper {
	return &AuthenticationTransport{
		APIKey:       apiKey,
		roundTripper: http.DefaultTransport,
	}
}
