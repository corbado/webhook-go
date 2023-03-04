package corbado_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corbado "github.com/corbado/webhook-go"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
	"github.com/corbado/webhook-go/pkg/logger"
)

const username = "webhookUsername"
const password = "webhookPassword"

func TestHandler(t *testing.T) {
	webhook, err := corbado.
		NewBuilder().
		SetLogger(logger.NewNull()).
		SetUsername(username).
		SetPassword(password).
		SetAuthMethodsCallback(authMethodsCallback).
		SetPasswordVerifyCallback(passwordVerifyCallback).
		Build()
	require.NoError(t, err)
	require.NotNil(t, webhook)

	standardHandler, err := webhook.GetStandardHandler()
	require.NoError(t, err)
	require.NotNil(t, standardHandler)

	ginHandler, err := webhook.GetGinHandler()
	require.NoError(t, err)
	require.NotNil(t, ginHandler)

	gin.SetMode(gin.ReleaseMode)

	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.POST("/webhook", ginHandler.Handle)

	tests := []struct {
		name                string
		createRequest       func() (*http.Request, error)
		assert              func(resp *http.Response)
		skipStandardHandler bool
		skipGinHandler      bool
	}{
		{
			name: "Missing authentication",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("POST", "/webhook", nil)
				if err != nil {
					return nil, err
				}

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "Invalid authentication",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("POST", "/webhook", nil)
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth("invalidUsername", "invalidPassword")

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "Invalid method",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("GET", "/webhook", nil)
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth(username, password)

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid method 'GET', only POST is allowed", string(body))
			},
			skipGinHandler: true,
		},
		{
			name: "Missing action",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("POST", "/webhook", nil)
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth(username, password)

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, "X-Corbado-Action header missing or empty", string(body))
			},
		},
		{
			name: "Invalid action",
			createRequest: func() (*http.Request, error) {
				body, err := os.ReadFile("testdata/authMethodsRequest.json")
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest("POST", "/webhook", bytes.NewReader(body))
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth(username, password)
				r.Header.Set("X-Corbado-Action", "invalidAction")

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid action given in X-Corbado-Action header ('invalidAction')", string(body))
			},
		},
		{
			name: "Success (authMethods)'",
			createRequest: func() (*http.Request, error) {
				body, err := os.ReadFile("testdata/authMethodsRequest.json")
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest("POST", "/webhook", bytes.NewReader(body))
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth(username, password)
				r.Header.Set("X-Corbado-Action", "authMethods")

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				expectedBody, err := os.ReadFile("testdata/authMethodsResponse.json")
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
			},
		},
		{
			name: "Success (passwordVerify)'",
			createRequest: func() (*http.Request, error) {
				body, err := os.ReadFile("testdata/passwordVerifyRequest.json")
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest("POST", "/webhook", bytes.NewReader(body))
				if err != nil {
					return nil, err
				}

				r.SetBasicAuth(username, password)
				r.Header.Set("X-Corbado-Action", "passwordVerify")

				return r, nil
			},
			assert: func(resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				expectedBody, err := os.ReadFile("testdata/passwordVerifyResponse.json")
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, body)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test standard handler
			if !test.skipStandardHandler {
				r, err := test.createRequest()
				assert.NoError(t, err)
				assert.NotNil(t, r)

				rr := httptest.NewRecorder()
				standardHandler.ServeHTTP(rr, r)
				resp := rr.Result()

				test.assert(resp)
				assert.NoError(t, resp.Body.Close())
			}

			// Test Gin handler
			if !test.skipGinHandler {
				r, err := test.createRequest()
				assert.NoError(t, err)
				assert.NotNil(t, r)

				rr := httptest.NewRecorder()
				ginRouter.ServeHTTP(rr, r)
				resp := rr.Result()

				test.assert(resp)
				assert.NoError(t, resp.Body.Close())
			}
		})
	}
}

func authMethodsCallback(_ string) (authmethodsresponse.Status, error) {
	return authmethodsresponse.StatusExists, nil
}

func passwordVerifyCallback(_ string, _ string) (bool, error) {
	return true, nil
}
