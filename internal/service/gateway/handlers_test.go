package gateway

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	apisixKafkaLog "github.com/naturalselectionlabs/api-gateway/app/apisix/kafkalog"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	apisixHTTPAPI "github.com/naturalselectionlabs/api-gateway/app/apisix/httpapi"
	jwtext "github.com/naturalselectionlabs/api-gateway/app/jwt"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/handlers"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/variables"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema/account"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema/key"
	"github.com/shanbay/gobay"
	"github.com/stretchr/testify/assert"
)

var (
	bapp    *gobay.Application
	handler http.Handler
	once    sync.Once

	expiredAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg3NUQyMjQzMUUxYjQ1OTBDNDg2ODRCQzBBOTY3YTI1NEQ2MjMwMzBiIiwiZXhwIjoxNjkxNTYyOTg1fQ.01fPPdUj6cRthQ-66AdEX3gmPEeKCGNiaiauyWdrP0s"

	validAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHhmMkY2ZTI4NkI2MmRlNEEwQjM2ODczMjcxQzAxNThlMTU4REZhYmU3IiwiY2hhaW5faWQiOjEsImV4cCI6MjAwNjkzODM1Mn0.lUdTv8nHuEu3rGT7BbXV-4GtiKKG98Rz9hCGeUS_apw"
	validAddress   = "0xf2F6e286B62de4A0B36873271C0158e158DFabe7"
)

func init() {
	bapp = setup()

	e := echo.New()
	configureMiddlewares(e, handlers.NewApp(), false)

	handler = e.Server.Handler
}

func setup() *gobay.Application {
	var err error
	once.Do(func() {
		// init app
		curdir, _ := os.Getwd()
		// fix compatibility issues on Windows path resolver
		if runtime.GOOS == "windows" {
			curdir = strings.ReplaceAll(curdir, "\\", "/")
		}
		root := path.Join(curdir, "..", "..")
		extensions := app.Extensions()
		bapp, err = gobay.CreateApp(root, "testing", extensions)
		if err != nil {
			panic(err)
		}
		app.InitExts(bapp)

		config := bapp.Config()

		variables.SIWEDomain = config.GetString("oapi_siwe_domain")

		apisixHTTPAPI.Config.APISixAdminEndpoint = config.GetString("apisix_admin_endpoint")
		apisixHTTPAPI.Config.APISixAdminKey = config.GetString("apisix_admin_key")
	})
	// migrate db
	err = app.EntClient.Schema.Create(context.Background())
	if err != nil {
		panic(err)
	}
	return bapp
}

func tearDown() {
	ctx := context.Background()

	// clear tables
	sqls := []string{
		`DROP SCHEMA public CASCADE;`,
		`CREATE SCHEMA public;`,
		`GRANT ALL ON SCHEMA public TO postgres;`,
		`GRANT ALL ON SCHEMA public TO public;`,
	}
	for _, sql := range sqls {
		if strings.TrimSpace(sql) == "" {
			continue
		}
		_, _ = app.EntExt.DB().Exec(sql)
	}

	// clear redis
	app.RedisExt.Client(ctx).FlushAll()
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
	account, err := model.AccountCreate(ctx, validAddress)
	if err != nil {
		panic(err)
	}
	return account
}

func tearDownAccount() {
	ctx := context.Background()
	_, err := app.EntClient.Account.Delete().Where(account.AddressEQ(validAddress)).Exec(ctx)
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
	bapp := setup()
	defer tearDown()

	ctx := context.Background()

	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).Hex()
	_, err = model.AccountGetByAddress(ctx, address)
	assert.Contains(t, err.Error(), "not found")

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
		return []byte(bapp.Config().GetString("jwt_key")), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, address, claims["address"])
	user, err := model.AccountGetByAddress(ctx, address)
	assert.NoError(t, err)
	assert.Equal(t, user.Address, claims["address"])

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
	err = app.EntClient.Key.Update().AddAPICallsCurrent(2).AddRuUsedCurrent(3).Where(
		key.KeyEQ(keyUUID),
	).Exec(ctx)
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
	assert.Equal(t, len(app.EntClient.Key.Query().AllX(ctx)), 2)
	user, err := model.AccountGetByAddress(ctx, validAddress)
	assert.NoError(t, err)
	keys, err := user.ListKeys(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(keys), 2)
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
	fakeUser, err := model.AccountCreate(ctx, "fake_user")
	assert.NoError(t, err)
	fakeUserKey, err := model.KeyCreate(ctx, fakeUser.ID, fakeUser.Address, "fake key")
	assert.NoError(t, err)
	assert.Equal(t, len(app.EntClient.Key.Query().AllX(ctx)), 3)
	app.EntClient.Key.Update().SetRuUsedCurrent(100).Where(key.IDEQ(fakeUserKey.ID)).ExecX(ctx)
	_, ruu, _, apicalls, err := fakeUser.GetUsage(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, ruu, 100)
	assert.EqualValues(t, apicalls, 0)

	// get other user key
	client.GET("/key/" + strconv.Itoa(fakeUserKey.Key.ID)).Expect().Status(http.StatusNotFound)

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
		dbKeyStrs[i] = k.Key.Key.String()
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

	ctx := context.Background()
	client := getAuth(t)

	// Get current withdrawal amount
	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": 0,
	})

	assert.Equal(t, app.EntClient.PendingWithdrawRequest.Query().CountX(ctx), 0)

	// Create a withdrawal request
	amount1 := float64(10)
	client.POST("/request/withdraw").WithQuery("amount", amount1).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount1,
	})

	assert.Equal(t, app.EntClient.PendingWithdrawRequest.Query().CountX(ctx), 1)

	// Update withdrawal request
	amount2 := 20.5
	client.POST("/request/withdraw").WithQuery("amount", amount2).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount2,
	})

	assert.Equal(t, len(app.EntClient.PendingWithdrawRequest.Query().AllX(ctx)), 1)

	// Unset withdrawal request
	client.POST("/request/withdraw").WithQuery("amount", 0).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": 0,
	})

	assert.Equal(t, app.EntClient.PendingWithdrawRequest.Query().CountX(ctx), 0)

	// Set again

	client.POST("/request/withdraw").WithQuery("amount", amount2).Expect().Status(http.StatusOK)

	client.GET("/request/withdraw").Expect().Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"amount": amount2,
	})

	assert.Equal(t, len(app.EntClient.PendingWithdrawRequest.Query().AllX(ctx)), 1)

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

//func Test_Proxy(t *testing.T) {
//	// mock upstream
//	httpmock.Activate()
//	defer httpmock.DeactivateAndReset()
//	responder := httpmock.NewStringResponder(200, `{"id": "aaa"}`)
//	httpmock.RegisterResponder(
//		"GET",
//		"http://rss3-api.pregod:8080/v1/activities/aaa",
//		responder,
//	)
//	httpmock.RegisterResponder(
//		"POST",
//		"http://rss3-api.pregod:8080/v1/accounts/activities",
//		responder,
//	)
//	httpmock.RegisterResponder(
//		"GET",
//		"http://rss3-api.pregod:8080/v1/mastodon/aaa/activities",
//		httpmock.NewStringResponder(404, ""),
//	)
//
//	setup()
//	defer tearDown()
//	acc := setupAccount()
//	defer tearDownAccount()
//
//	ctx := context.Background()
//	noUser := getAuth(t, "")
//
//	// use an invalid key
//	key := uuid.New().String()
//	noUser.GET("/data/v1/activities/aaa").WithHeader(app.KeyHeader, key).Expect().Status(http.StatusUnauthorized)
//
//	// use a valid key
//	k, err := model.KeyCreate(ctx, acc.ID, acc.Address, "valid key")
//	assert.NoError(t, err)
//	assert.Equal(t, len(app.EntClient.Key.Query().AllX(ctx)), 1)
//	assert.Equal(t, k.RuUsedCurrent, int64(0))
//	assert.Equal(t, k.APICallsCurrent, int64(0))
//	noUser.GET("/data/v1/activities/aaa").WithHeader(app.KeyHeader, k.Key.Key.String()).Expect().Status(http.StatusBadRequest).JSON().Object().HasValue("msg", "No enough RU")
//
//	// create a new key
//	newK, err := model.KeyCreate(ctx, acc.ID, acc.Address, "new key")
//	assert.NoError(t, err)
//	assert.Equal(t, len(app.EntClient.Key.Query().AllX(ctx)), 2)
//	assert.Equal(t, newK.RuUsedCurrent, int64(0))
//	assert.Equal(t, newK.APICallsCurrent, int64(0))
//
//	// use a valid key with sufficient RU
//	// + 2 RU
//	app.EntClient.Account.Update().AddRuLimit(2).Where(account.IDEQ(acc.ID)).ExecX(ctx)
//	assert.EqualValues(t, 2, app.EntClient.Account.GetX(ctx, acc.ID).RuLimit)
//	// - 1 RU
//	noUser.GET("/data/v1/activities/aaa").WithHeader(app.KeyHeader, k.Key.Key.String()).Expect().Status(http.StatusOK)
//	assert.EqualValues(t, 2, app.EntClient.Account.GetX(ctx, acc.ID).RuLimit)
//	k, err = model.KeyGetByKey(ctx, k.Key.Key.String())
//	assert.NoError(t, err)
//	assert.Equal(t, k.RuUsedCurrent, int64(1))
//	assert.Equal(t, k.APICallsCurrent, int64(1))
//	acc, err = k.GetAccount(ctx)
//	assert.NoError(t, err)
//	balance, err := acc.GetBalance(ctx)
//	assert.NoError(t, err)
//	assert.EqualValues(t, 1, balance)
//	_, ruUsed, _, apiCallMade, err := acc.GetUsage(ctx)
//	assert.NoError(t, err)
//	assert.EqualValues(t, 1, ruUsed)
//	assert.EqualValues(t, 1, apiCallMade)
//
//	// make a POST request and consume 10 RU
//	// carry over balance = 1 RU
//	// - 10 RU
//	noUser.POST("/data/v1/accounts/activities").WithHeader(app.KeyHeader, k.Key.Key.String()).Expect().Status(http.StatusOK)
//	k, err = model.KeyGetByKey(ctx, k.Key.Key.String())
//	assert.NoError(t, err)
//	acc, err = k.GetAccount(ctx)
//	assert.NoError(t, err)
//	balance, err = acc.GetBalance(ctx)
//	assert.NoError(t, err)
//	assert.EqualValues(t, -9, balance)
//	_, ruUsed, _, apiCallMade, err = acc.GetUsage(ctx)
//	assert.NoError(t, err)
//	// total RU used = 11
//	assert.EqualValues(t, 11, ruUsed)
//	assert.EqualValues(t, 2, apiCallMade)
//	newK, err = model.KeyGetByKey(ctx, newK.Key.Key.String())
//	assert.NoError(t, err)
//	assert.Equal(t, newK.RuUsedCurrent, int64(0))
//	assert.Equal(t, newK.APICallsCurrent, int64(0))
//
//	// use a valid key with insufficient RU
//	// carry over balance = -9 RU
//	noUser.GET("/data/v1/mastodon/aaa/activities").WithHeader(app.KeyHeader, newK.Key.Key.String()).Expect().Status(http.StatusBadRequest).JSON().Object().HasValue("msg", "No enough RU")
//	// + 10 RU
//	app.EntClient.Account.Update().AddRuLimit(10).Where(account.IDEQ(acc.ID)).ExecX(ctx)
//	assert.EqualValues(t, 12, app.EntClient.Account.GetX(ctx, acc.ID).RuLimit)
//	acc, err = newK.GetAccount(ctx)
//	assert.NoError(t, err)
//	balance, err = acc.GetBalance(ctx)
//	assert.NoError(t, err)
//	assert.EqualValues(t, 1, balance)
//	noUser.GET("/data/v1/mastodon/aaa/activities").WithHeader(app.KeyHeader, newK.Key.Key.String()).Expect().Status(http.StatusNotFound)
//	acc, err = newK.GetAccount(ctx)
//	assert.NoError(t, err)
//	balance, err = acc.GetBalance(ctx)
//	assert.NoError(t, err)
//	// Balance should not change as the API call was not successful
//	assert.EqualValues(t, 1, balance)
//}

func Test_JWTSign(t *testing.T) {
	token, err := app.JwtExt.SignToken(&jwtext.User{
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
	fakeUser, err := model.AccountCreate(ctx, "fake_user")
	assert.NoError(t, err)
	app.EntClient.Account.UpdateOneID(fakeUser.ID).SetRuLimit(100).ExecX(ctx)
	assert.Equal(t, app.EntClient.Account.GetX(ctx, fakeUser.ID).RuLimit, int64(100))

	// Create test key
	fakeUserKey, err := model.KeyCreate(ctx, fakeUser.ID, fakeUser.Address, "fake key")
	assert.NoError(t, err)
	assert.Equal(t, app.EntClient.Key.GetX(ctx, fakeUserKey.ID).RuUsedCurrent, int64(0))

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

	for _, log := range requestLogs {
		handlers.ProcessAccessLog(log)
	}

	// Check RU consumption
	assert.Equal(t, app.EntClient.Account.GetX(ctx, fakeUser.ID).RuLimit, int64(100))
	assert.Equal(t, app.EntClient.Key.GetX(ctx, fakeUserKey.ID).RuUsedCurrent, int64(10))

	_, ruUsedCurrent, _, apiCallsCurrent, err := fakeUser.GetUsage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, ruUsedCurrent, int64(10))
	assert.Equal(t, apiCallsCurrent, int64(2))

}
