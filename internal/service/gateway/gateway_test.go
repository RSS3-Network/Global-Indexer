package gateway

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/httpapi"
	apisixKafkaLog "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/kafkalog"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/handlers"
	jwtImpl "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/siwe"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	handler http.Handler

	expiredAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg3NUQyMjQzMUUxYjQ1OTBDNDg2ODRCQzBBOTY3YTI1NEQ2MjMwMzBiIiwiZXhwIjoxNjkxNTYyOTg1fQ.01fPPdUj6cRthQ-66AdEX3gmPEeKCGNiaiauyWdrP0s"

	validAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHhmMkY2ZTI4NkI2MmRlNEEwQjM2ODczMjcxQzAxNThlMTU4REZhYmU3IiwiY2hhaW5faWQiOjEsImV4cCI6MjAwNjkzODM1Mn0.lUdTv8nHuEu3rGT7BbXV-4GtiKKG98Rz9hCGeUS_apw"
	validAddress   = common.HexToAddress("0xf2F6e286B62de4A0B36873271C0158e158DFabe7")

	fakeUserAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

var (
	gatewayApp       *handlers.App
	redisClient      *redis.Client
	databaseClient   *gorm.DB
	jwtClient        *jwtImpl.JWT
	siweClient       *siwe.SIWE
	apisixAPIService *apisixHTTPAPI.HTTPAPIService
)

const (
	JWT_KEY = "abcdefg1234567"
)

func init() {
	// Prepare databaseClient
	dbc, err := dialer.Dial(context.Background(), &config.Database{
		Driver: database.DriverCockroachDB,
		URI:    "postgres://root@localhost:26257/defaultdb",
	})
	if err != nil {
		log.Panic(err)
	}
	err = dbc.Migrate(context.Background())
	if err != nil {
		log.Panic(err)
	}
	databaseClient = dbc.Raw()

	// Prepare redisClient
	rc, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Panic(err)
	}
	redisClient = redis.NewClient(rc)

	// Initialize APISIX configurations
	apisixAPIService, err = apisixHTTPAPI.New(
		"http://localhost:9180",
		"edd1c9f034335f136f87ad84b625c8f1",
	)
	if err != nil {
		log.Panic(err)
	}

	// Prepare JWT
	jwtClient, err = jwtImpl.New(JWT_KEY)
	if err != nil {
		log.Panic(err)
	}

	// Prepare SIWE
	siweClient, err = siwe.New("localhost:3000", redisClient)
	if err != nil {
		log.Panic(err)
	}

	// Prepare echo
	e := echo.New()
	gatewayApp, err = handlers.NewApp(
		apisixAPIService,
		redisClient,
		databaseClient,
		jwtClient,
		siweClient,
	)
	if err != nil {
		log.Panic(err)
	}

	// Configure middlewares
	configureMiddlewares(e, gatewayApp, jwtClient, databaseClient, apisixAPIService)

	handler = e.Server.Handler
}

func setup() {
	// Nothing to do for now
}

func tearDown() {
	ctx := context.Background()

	// clear tables
	sqls := []string{
		`DELETE FROM gateway.consumption_log CASCADE;`,
		`DELETE FROM gateway.key CASCADE;`,
		`DELETE FROM gateway.pending_withdraw_request CASCADE;`,
		`DELETE FROM gateway.account CASCADE;`,
		`DELETE FROM gateway.br_collected CASCADE;`,
		`DELETE FROM gateway.br_deposited CASCADE;`,
		`DELETE FROM gateway.br_withdrawn CASCADE;`,
	}
	for _, sql := range sqls {
		if strings.TrimSpace(sql) == "" {
			continue
		}
		databaseClient.Exec(sql)
	}

	// clear redis
	redisClient.FlushAll(ctx)
}

// opts: authToken, Content-Type
// defaults: address=validAddress; Content-Type=application/json
func getAuth(t *testing.T, opts ...string) *httpexpect.Expect {
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})

	defaults := []string{validAuthToken, "application/json"}
	for index, opt := range opts {
		defaults[index] = opt
	}
	authToken, contentType := defaults[0], defaults[1]

	return e.Builder(func(r *httpexpect.Request) {
		r.WithHeader("Content-Type", contentType)
		if authToken != "" {
			r.WithCookie("auth_token", authToken)
		}
	})
}

func setupAccount() model.Account {
	ctx := context.Background()
	acc, err := model.AccountCreate(ctx, validAddress, databaseClient, apisixAPIService)
	if err != nil {
		panic(err)
	}
	return *acc
}

func tearDownAccount() {
	err := databaseClient.Delete(&table.GatewayAccount{
		Address: validAddress,
	}).Error
	if err != nil {
		panic(err)
	}
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sign := signHash([]byte(message))
	signature, err := crypto.Sign(sign.Bytes(), privateKey)

	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return signature, nil
}

func constructMessage(domain string, nonce string, when time.Time, address string) string {
	return fmt.Sprintf("%s wants you to sign in with your Ethereum account:\n%s\n\nSign In With Ethereum to prove you control this wallet.\n\nURI: http://%s\nVersion: 1\nChain ID: 1\nNonce: %s\nIssued At: %s",
		domain,
		address,
		domain,
		nonce,
		when.Format("2006-01-02T15:04:05.999Z"),
	)
}

func Test_SIWEAuth(t *testing.T) {
	setup()
	defer tearDown()

	ctx := context.Background()

	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)
	publicKey := privateKey.PublicKey
	addressRaw := crypto.PubkeyToAddress(publicKey)
	address := addressRaw.Hex()
	_, exist, err := model.AccountGetByAddress(ctx, addressRaw, databaseClient, apisixAPIService)
	assert.False(t, exist)
	assert.NoError(t, err)

	now := time.Now()

	// 400 empty
	getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{}).Expect().Status(http.StatusBadRequest)

	wrongSignature := "0x4033e439566439d0cb5ca8c9412f0932e06f6d125db0415bf342a274f9985aa352083dd361d3b238777356af2b993601601d87f1da3a3da381eaf38394ee957d1c"

	// 401 wrong domain
	nonce := getAuth(t, "").GET("/users/siwe/nonce").Expect().Status(http.StatusOK).Text()
	msg := constructMessage("localhost:3001", nonce.Raw(), now, address)
	getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{
		"message":   msg,
		"signature": wrongSignature,
	}).Expect().Status(http.StatusUnauthorized)

	// 401 invalid signature format
	nonce = getAuth(t, "").GET("/users/siwe/nonce").Expect().Status(http.StatusOK).Text()
	msg = constructMessage("localhost:3000", nonce.Raw(), now, address)
	// using an invalid signature
	getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{
		"message":   msg,
		"signature": "0x.......",
	}).Expect().Status(http.StatusUnauthorized)

	// 401 wrong address
	nonce = getAuth(t, "").GET("/users/siwe/nonce").Expect().Status(http.StatusOK).Text()
	msg = constructMessage("localhost:3000", nonce.Raw(), now, address)
	// using an unmatched signature
	getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{
		"message":   msg,
		"signature": wrongSignature,
	}).Expect().Status(http.StatusUnauthorized)

	// 401 no such nonce
	// using old (consumed) nonce
	// using old message
	signature, err := signMessage(msg, privateKey)
	assert.NotEmpty(t, signature)
	assert.NoError(t, err)
	getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{
		"message":   msg,
		"signature": hexutil.Encode(signature),
	}).Expect().Status(http.StatusUnauthorized)

	// 200: create user (full flow)
	nonce = getAuth(t, "").GET("/users/siwe/nonce").Expect().Status(http.StatusOK).Text()
	msg = constructMessage("localhost:3000", nonce.Raw(), now, address)
	signature, err = signMessage(msg, privateKey)
	assert.NotEmpty(t, signature)
	assert.NoError(t, err)
	cookie := getAuth(t, "").POST("/users/siwe/verify").WithJSON(map[string]any{
		"message":   msg,
		"signature": hexutil.Encode(signature),
	}).Expect().Status(http.StatusOK).Cookie("auth_token").Raw()

	// check cookie
	getAuth(t, cookie.Value).GET("/ru/status").Expect().Status(http.StatusOK)
	// parse jwt
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_KEY), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, address, claims["address"])
	user, exist, err := model.AccountGetByAddress(ctx, addressRaw, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, user.Address.Hex(), claims["address"])

	// check expired cookie
	getAuth(t, expiredAuthToken).GET("/ru/status").Expect().Status(http.StatusUnauthorized)

	// check logout
	getAuth(t, cookie.Value).GET("/users/siwe/logout").Expect().Status(http.StatusOK)
}

func Test_KeyAndRU(t *testing.T) {
	setup()
	defer tearDown()
	setupAccount()
	defer tearDownAccount()

	ctx := context.Background()
	client := getAuth(t)

	// get empty status
	client.GET("/ru/status").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"ru_limit":          0,
		"ru_used_total":     0,
		"ru_used_current":   0,
		"api_calls_total":   0,
		"api_calls_current": 0,
	})

	// generate key
	obj := client.POST("/key").WithJSON(map[string]any{
		"name": "new key",
	}).Expect().Status(http.StatusOK).
		JSON().Object()
	obj.Value("key").String().NotEmpty()
	obj.Value("key").String().Length().IsEqual(36)
	obj.Value("ru_used_total").IsEqual(0)
	obj.Value("ru_used_current").IsEqual(0)
	obj.Value("api_calls_total").IsEqual(0)
	obj.Value("api_calls_current").IsEqual(0)
	obj.Value("name").IsEqual("new key")

	// still empty status
	client.GET("/ru/status").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"ru_limit":          0,
		"ru_used_total":     0,
		"ru_used_current":   0,
		"api_calls_total":   0,
		"api_calls_current": 0,
	})

	keyId := obj.Value("id").Number()
	keyId.IsInt()
	// get key
	client.GET("/key/" + strconv.Itoa(int(keyId.Raw()))).Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"id":                keyId.Raw(),
		"key":               obj.Value("key").String().Raw(),
		"name":              "new key",
		"ru_used_total":     0,
		"ru_used_current":   0,
		"api_calls_total":   0,
		"api_calls_current": 0,
	})

	// incr ru_used and api_calls
	keyUUID, err := uuid.Parse(obj.Value("key").String().Raw())
	assert.NoError(t, err)
	err = databaseClient.Model(&table.GatewayKey{}).Where("key = ?", keyUUID).Updates(map[string]interface{}{
		"api_calls_current": gorm.Expr("api_calls_current + ?", 2),
		"ru_used_current":   gorm.Expr("ru_used_current + ?", 3),
	}).Error
	assert.NoError(t, err)
	client.GET("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"id":                obj.Value("id").Number().Raw(),
		"key":               obj.Value("key").String().Raw(),
		"name":              "new key",
		"ru_used_total":     0,
		"ru_used_current":   3,
		"api_calls_total":   0,
		"api_calls_current": 2,
	})

	// check status
	client.GET("/ru/status").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"ru_limit":          0,
		"ru_used_total":     0,
		"ru_used_current":   3,
		"api_calls_total":   0,
		"api_calls_current": 2,
	})

	// generate new key
	obj = client.POST("/key").WithJSON(map[string]any{
		"name": "new key 2",
	}).Expect().Status(http.StatusOK).
		JSON().Object()
	obj.Value("key").String().NotEmpty()
	obj.Value("key").String().Length().IsEqual(36)
	obj.Value("ru_used_current").IsEqual(0)
	obj.Value("ru_used_total").IsEqual(0)
	obj.Value("api_calls_current").IsEqual(0)
	obj.Value("api_calls_total").IsEqual(0)
	var keyCounts int64
	err = databaseClient.Model(&table.GatewayKey{}).Count(&keyCounts).Error
	assert.NoError(t, err)
	assert.Equal(t, keyCounts, 2)
	user, exist, err := model.AccountGetByAddress(ctx, validAddress, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	assert.True(t, exist)
	client.GET("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"id":                obj.Value("id").Number().Raw(),
		"key":               obj.Value("key").String().Raw(),
		"name":              "new key 2",
		"ru_used_total":     0,
		"ru_used_current":   0,
		"api_calls_total":   0,
		"api_calls_current": 0,
	})

	// status still the same
	client.GET("/ru/status").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"ru_limit":          0,
		"ru_used_total":     0,
		"ru_used_current":   3,
		"api_calls_total":   0,
		"api_calls_current": 2,
	})

	// create new account with key
	fakeUser, err := model.AccountCreate(ctx, fakeUserAddr, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	fakeUserKey, err := model.KeyCreate(ctx, fakeUser.Address, "fake key", databaseClient, apisixAPIService)
	assert.NoError(t, err)
	err = databaseClient.Model(&table.GatewayKey{}).Count(&keyCounts).Error
	assert.NoError(t, err)
	assert.Equal(t, keyCounts, 3)
	err = databaseClient.Model(&table.GatewayKey{}).Where("id = ?", fakeUserKey.ID).Update("ru_used_current", 100).Error
	assert.NoError(t, err)
	_, ruu, _, apicalls, err := fakeUser.GetUsage(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, ruu, 100)
	assert.EqualValues(t, apicalls, 0)

	// get other user key
	client.GET("/key/" + strconv.Itoa(int(fakeUserKey.ID))).Expect().Status(http.StatusNotFound)

	// status still the same
	client.GET("/ru/status").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"ru_limit":          0,
		"ru_used_total":     0,
		"ru_used_current":   3,
		"api_calls_total":   0,
		"api_calls_current": 2,
	})

	// get keys
	objs := client.GET("/keys").Expect().Status(http.StatusOK).
		JSON().Array()
	objs.Length().IsEqual(2)
	dbKeys, err := user.ListKeys(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(dbKeys), 2)
	dbKeyStrs := make([]string, 2)
	for i, k := range dbKeys {
		dbKeyStrs[i] = k.Key.String()
	}
	for _, obj := range objs.Iter() {
		obj.Object().Value("key").String().NotEmpty()
		assert.Contains(t, dbKeyStrs, obj.Object().Value("key").String().Raw())
	}

	// check is RU remain greater than 0
	result, err := user.GetBalance(ctx)
	assert.NoError(t, err)
	assert.False(t, result > 0)
	assert.EqualValues(t, result, -3)

	// rename key
	obj = client.PUT("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).WithJSON(map[string]any{
		"name": "new key name",
	}).Expect().Status(http.StatusOK).
		JSON().Object()
	obj.Value("name").IsEqual("new key name")

	oldKey := obj.Value("key").String().Raw()

	// rotate key
	obj = client.PATCH("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).
		Expect().Status(http.StatusOK).
		JSON().Object()
	obj.Value("key").NotEqual(oldKey)

	// delete key no Auth
	getAuth(t, "").DELETE("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).Expect().Status(http.StatusUnauthorized)
	// delete key
	client.DELETE("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).Expect().Status(http.StatusOK)
	objs = client.GET("/keys").Expect().Status(http.StatusOK).
		JSON().Array()
	objs.Length().IsEqual(1)
	result, err = user.GetBalance(ctx)
	assert.EqualValues(t, result, -3)
	// delete key no key
	client.DELETE("/key/" + strconv.Itoa(int(obj.Value("id").Number().Raw()))).Expect().Status(http.StatusNotFound).
		JSON().Object().Value("msg").String().Contains("Not Found")
}

func Test_RequestWithdraw(t *testing.T) {
	setup()
	defer tearDown()
	setupAccount()
	defer tearDownAccount()

	client := getAuth(t)

	// Get current withdrawal amount
	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": 0,
	})

	var pendingWithdrawRequestCount int64
	err := databaseClient.Model(&table.GatewayPendingWithdrawRequest{}).Count(&pendingWithdrawRequestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, pendingWithdrawRequestCount, int64(0))

	// Create a withdrawal request
	amount1 := float64(10)
	client.POST("/request/withdraw").WithQuery("amount", amount1).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount1,
	})

	err = databaseClient.Model(&table.GatewayPendingWithdrawRequest{}).Count(&pendingWithdrawRequestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, pendingWithdrawRequestCount, int64(1))

	// Update withdrawal request
	amount2 := 20.5
	client.POST("/request/withdraw").WithQuery("amount", amount2).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount2,
	})

	err = databaseClient.Model(&table.GatewayPendingWithdrawRequest{}).Count(&pendingWithdrawRequestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, pendingWithdrawRequestCount, int64(1))

	// Unset withdrawal request
	client.POST("/request/withdraw").WithQuery("amount", 0).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": 0,
	})

	err = databaseClient.Model(&table.GatewayPendingWithdrawRequest{}).Count(&pendingWithdrawRequestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, pendingWithdrawRequestCount, int64(0))

	// Set again

	client.POST("/request/withdraw").WithQuery("amount", amount2).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount2,
	})

	err = databaseClient.Model(&table.GatewayPendingWithdrawRequest{}).Count(&pendingWithdrawRequestCount).Error
	assert.NoError(t, err)
	assert.Equal(t, pendingWithdrawRequestCount, int64(1))

	// Unset again with negative amount
	client.POST("/request/withdraw").WithQuery("amount", -1).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": 0,
	})

}

func Test_HealthCheck(t *testing.T) {
	setup()
	defer tearDown()

	noUser := getAuth(t, "")

	noUser.GET("/health").Expect().Status(http.StatusBadRequest)

	noUser.GET("/health").WithQuery("type", "readiness").Expect().Status(http.StatusOK)
	noUser.GET("/health").WithQuery("type", "liveness").Expect().Status(http.StatusOK)
}

func Test_JWTSign(t *testing.T) {
	token, err := jwtClient.SignToken(&jwtImpl.User{
		Address: validAddress,
		ChainId: 1,
		Expires: 2006938352,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	t.Log(token)
}

func Test_ProcessAccessLog(t *testing.T) {
	setup()
	defer tearDown()

	ctx := context.Background()

	// Create test account and add some RU
	fakeUser, err := model.AccountCreate(ctx, fakeUserAddr, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	err = databaseClient.Model(&table.GatewayAccount{}).Where("address = ?", fakeUser.Address).Update("ru_limit", 100).Error
	assert.NoError(t, err)
	fakeUser, exist, err := model.AccountGetByAddress(ctx, fakeUser.Address, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, fakeUser.RuLimit, int64(100))

	// Create test key
	fakeUserKey, err := model.KeyCreate(ctx, fakeUser.Address, "fake key", databaseClient, apisixAPIService)
	assert.NoError(t, err)
	fakeUserKey, exist, err = fakeUser.GetKey(ctx, fakeUserKey.ID)
	assert.True(t, exist)
	assert.Equal(t, fakeUserKey.RuUsedCurrent, int64(0))

	// Mock some request logs
	consumer := fmt.Sprintf("key_%d", fakeUserKey.ID)
	time1, err := time.Parse(time.RFC3339, "2023-11-01T08:13:18Z")
	assert.NoError(t, err)
	time2, err := time.Parse(time.RFC3339, "2023-11-01T08:13:27Z")
	assert.NoError(t, err)
	time3, err := time.Parse(time.RFC3339, "2023-11-01T08:13:43Z")
	assert.NoError(t, err)

	requestLogs := []apisixKafkaLog.AccessLog{
		{ // Should be billed
			Consumer:  &consumer,
			Timestamp: time1,
			ClientIP:  "172.26.0.1",
			RouteID:   "484047074917089994",
			URI:       "/data/accounts/activities", // 10 RU
			Host:      "127.0.0.1",
			Status:    200,
		},
		{ // Should not be billed
			Consumer:  &consumer,
			Timestamp: time2,
			ClientIP:  "172.26.0.1",
			RouteID:   "484047074917089994",
			URI:       "/data/accounts/activities",
			Host:      "127.0.0.1",
			Status:    500,
		},
		{ // Should not panic
			URI:       "/data/accounts/activities",
			Timestamp: time3,
			ClientIP:  "172.26.0.1",
			RouteID:   "484047074917089994",
			Host:      "127.0.0.1",
			Status:    401,
		},
	}

	for _, reqLog := range requestLogs {
		gatewayApp.ProcessAccessLog(reqLog)
	}

	// Check RU consumption
	fakeUser, exist, err = model.AccountGetByAddress(ctx, fakeUser.Address, databaseClient, apisixAPIService)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, fakeUser.RuLimit, int64(100))
	fakeUserKey, exist, err = fakeUser.GetKey(ctx, fakeUserKey.ID)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, fakeUserKey.RuUsedCurrent, int64(10))

	_, ruUsedCurrent, _, apiCallsCurrent, err := fakeUser.GetUsage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, ruUsedCurrent, int64(10))
	assert.Equal(t, apiCallsCurrent, int64(2))

}
