package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/stretchr/testify/mock"
)

var (
	nullData            = `{"data":null}`
	errResponse         = `{"error":"it is error","error_code":"internal_error"}`
	validActivityData   = `{"data":{"id":"0x00000000000000000000000000000aadbebe59e6a53f0d552388f00f3cf8e5b1","owner":"0x0000000000000000000000000000000000000091","network":"farcaster","index":0,"from":"0x15426Ef0B2A15B3F7dcA513BD70F88A9481a0320","to":"0x4aF0919907ccdBFF6C80463A47B4Db24599CC0d5","tag":"social","type":"share","platform":"Farcaster","total_actions":2,"actions":[{"tag":"social","type":"share","platform":"Farcaster","from":"0xf0a68CD0e9AC293Ae4E9A5730852C9Cd0eaa56a0","to":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","metadata":{"handle":"tricksloaded","profile_id":"277353","publication_id":"0x00000aAdBebe59E6a53F0d552388f00f3cF8E5b1","target":{"handle":"terrytat","body":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","profile_id":"277460","publication_id":"0x6C14EAe3Ebda408B7d20968fB992b0DF7E8ec2cd"}}},{"tag":"social","type":"share","platform":"Farcaster","from":"0x0000000000000000000000000000000000000091","to":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","metadata":{"handle":"tricksloaded","profile_id":"277353","publication_id":"0x00000aAdBebe59E6a53F0d552388f00f3cF8E5b1","target":{"handle":"terrytat","body":"0xd269Fe371893B81F4c572138F11fCfC91ba09083","profile_id":"277460","publication_id":"0x6C14EAe3Ebda408B7d20968fB992b0DF7E8ec2cd"}}}],"direction":"out","success":true,"timestamp":1711059936},"meta":{"totalPages":1}}`
	validActivitiesData = `{"data":[{"id":"0x000000000000000000000000cea19adbd5060fbd7797a61040c5c378c4e469c3","owner":"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045","network":"farcaster","index":0,"from":"0xb005D4eF6Dd474634f8cF546AAb674E6a9b12A28","to":"0xADD746Be46fF36f10C81d6e3Ba282537f4c68077","tag":"social","type":"share","platform":"Farcaster","total_actions":1,"actions":[{"tag":"social","type":"share","platform":"Farcaster","from":"0x6cAFB2D9cD97e2bb8743C1bD046e716963b115d0","to":"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045","metadata":{"handle":"discovolante","profile_id":"417128","publication_id":"0xCEa19aDbD5060Fbd7797A61040C5C378c4e469c3","target":{"handle":"vitalik.eth","body":"Degen communism: the only correct political ideology\n\nhttps://vitalik.eth.limo/general/2024/04/01/dc.html","profile_id":"5650","publication_id":"0xadE35D02027deef831457C04B2D9C8EfE33AC78A"}}}],"direction":"in","success":true,"timestamp":1711992401},{"id":"0x00000000000000000000000014d46b56c5e64f5a9d88ac2c26118b68b1257402","owner":"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045","network":"farcaster","index":0,"from":"0xe8198c3315335dff97CB69f335E3dFCe5d0ED807","to":"0xADD746Be46fF36f10C81d6e3Ba282537f4c68077","tag":"social","type":"comment","platform":"Farcaster","total_actions":1,"actions":[{"tag":"social","type":"comment","platform":"Farcaster","from":"0x2fcB45C294731Ea7EA1F615ffC280CECe30c16cC","to":"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045","metadata":{"handle":"sebix.eth","body":"Interesting views on taxation, as well as on DAO governance. Things I have problem with: name - \"degen communism\" will not appear to the \"normies\" - people you need to convince; immigration - it's just naive to think large scale, mostly uncontrolled immigration of peoples with completely different worldviews will work.","profile_id":"417182","publication_id":"0x14d46b56c5e64F5a9d88aC2C26118B68b1257402","target":{"handle":"vitalik.eth","body":"Degen communism: the only correct political ideology\n\nhttps://vitalik.eth.limo/general/2024/04/01/dc.html","profile_id":"5650","publication_id":"0xadE35D02027deef831457C04B2D9C8EfE33AC78A"}}}],"direction":"in","success":true,"timestamp":1711983144}],"meta":{"cursor":"0x00000000000000000000000014d46b56c5e64f5a9d88ac2c26118b68b1257402:farcaster"}}`

	nodeMap = map[common.Address]model.RequestMeta{
		common.HexToAddress("0x123"): {
			Method:   "GET",
			Endpoint: "http://localhost:8070",
			Body:     nil,
		},
		common.HexToAddress("0x234"): {
			Method:   "GET",
			Endpoint: "http://localhost:8080",
			Body:     nil,
		},
		common.HexToAddress("0x567"): {
			Method:   "GET",
			Endpoint: "http://localhost:8090",
			Body:     nil,
		},
	}
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) FetchWithMethod(ctx context.Context, _, endpoint string, _ io.Reader) (io.ReadCloser, error) {
	args := m.Called(ctx, endpoint)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestDistributeRequestWithValidData(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(validActivitiesData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response.Data) != validActivitiesData {
		t.Errorf("Expected 'valid', got %s", response.Data)
	}

	if !response.Valid {
		t.Errorf("Expected 'true', got %v", response.Valid)
	}

	if response.Address != common.HexToAddress("0x567") {
		t.Errorf("Expected '0x567', got %v", response.Address.String())
	}
}

func TestDistributeRequestWithNullData(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response.Data) != nullData {
		t.Errorf("Expected 'null', got %s", response.Data)
	}

	if response.Valid {
		t.Errorf("Expected 'false', got %v", response.Valid)
	}
}

func TestDistributeRequestWithNodeError(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(errResponse)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(errResponse)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(errResponse)), nil)

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	if response.Valid {
		t.Errorf("Expected 'false', got %v", response.Valid)
	}
}

func TestDistributeRequestWithError(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(errResponse)), errors.New("error 8070"))
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(errResponse)), errors.New("error 8080"))
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(errResponse)), errors.New("error 8090"))

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	if response.Valid {
		t.Errorf("Expected 'false', got %v", response.Valid)
	}
}

func TestDistributeRequestAndReturnNull(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(nullData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(errResponse)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(errResponse)), errors.New("error 8090"))

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if string(response.Data) != nullData {
		t.Errorf("Expected 'valid', got %s", response.Data)
	}
}

func TestDistributeRequest(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	r := SimpleRouter{httpClient: mockClient}

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8070").Return(io.NopCloser(bytes.NewBufferString(validActivityData)), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(validActivityData)), nil)
	//mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080").Return(io.NopCloser(bytes.NewBufferString(errResponse)), errors.New("error 8090"))
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8090").Return(io.NopCloser(bytes.NewBufferString(validActivityData)), nil)

	response, err := r.DistributeRequest(context.Background(), nodeMap, process)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err = response.Err; err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if !response.Valid {
		t.Errorf("Expected 'false', got %v", response.Valid)
	}

	if string(response.Data) != validActivityData {
		t.Errorf("Expected 'valid', got %s", response.Data)
	}
}

func process(responses []*model.DataResponse) {
	fmt.Println(len(responses))

	for _, res := range responses {
		fmt.Printf("address: %s,valid:%v\n", res.Address.String(), res.Valid)
	}
}
