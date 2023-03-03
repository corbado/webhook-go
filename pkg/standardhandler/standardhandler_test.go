package standardhandler_test

import (
	"bytes"
	"crypto/sha256"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
	"github.com/corbado/webhook-go/pkg/logger"
	"github.com/corbado/webhook-go/pkg/standardhandler"
)

const username = "webhookUsername"
const password = "webhookPassword"

//nolint:funlen
func TestImpl_ServeHTTP(t *testing.T) {
	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))

	handler, err := standardhandler.New(
		logger.NewNull(),
		usernameHash,
		passwordHash,
		authMethodsCallback,
		passwordVerifyCallback,
	)
	require.NoError(t, err)
	require.NotNil(t, handler)

	tests := []struct {
		name          string
		createRequest func() (*http.Request, error)
		assert        func(resp *http.Response)
	}{
		{
			name: "Missing authentication",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("GET", "/", nil)
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
				r, err := http.NewRequest("GET", "/", nil)
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
				r, err := http.NewRequest("GET", "/", nil)
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
		},
		{
			name: "Missing action",
			createRequest: func() (*http.Request, error) {
				r, err := http.NewRequest("POST", "/", nil)
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
				assert.Equal(t, "X-Corbado-Action missing or empty", string(body))
			},
		},
		{
			name: "Invalid action",
			createRequest: func() (*http.Request, error) {
				body, err := os.ReadFile("testdata/authMethodsRequest.json")
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest("POST", "/", bytes.NewReader(body))
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

				r, err := http.NewRequest("POST", "/", bytes.NewReader(body))
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

				r, err := http.NewRequest("POST", "/", bytes.NewReader(body))
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
			r, err := test.createRequest()
			assert.NoError(t, err)
			assert.NotNil(t, r)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, r)
			resp := rr.Result()
			defer resp.Body.Close()

			test.assert(resp)
		})
	}
}

func authMethodsCallback(_ string) (authmethodsresponse.Status, error) {
	return authmethodsresponse.StatusExists, nil
}

func passwordVerifyCallback(_ string, _ string) (bool, error) {
	return true, nil
}
