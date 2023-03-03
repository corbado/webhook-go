package standardhandler

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/corbado/webhook-go/pkg/callback"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsrequest"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
	"github.com/corbado/webhook-go/pkg/dto/passwordverifyrequest"
	"github.com/corbado/webhook-go/pkg/dto/passwordverifyresponse"
	"github.com/corbado/webhook-go/pkg/logger"
)

type StandardHandler struct {
	logger                 logger.Logger
	usernameHash           [32]byte
	passwordHash           [32]byte
	authMethodsCallback    callback.AuthMethods
	passwordVerifyCallback callback.PasswordVerify
}

// New returns standard handler which can be used in standard HTTP library.
func New(
	logger logger.Logger,
	usernameHash [32]byte,
	passwordHash [32]byte,
	authMethodsCallback callback.AuthMethods,
	passwordVerifyCallback callback.PasswordVerify,
) (*StandardHandler, error) {
	if logger == nil {
		return nil, errors.New("empty parameter logger")
	}

	if authMethodsCallback == nil {
		return nil, errors.New("empty parameter authMethodsCallback")
	}

	if passwordVerifyCallback == nil {
		return nil, errors.New("empty parameter passwordVerifyCallback")
	}

	return &StandardHandler{
		logger:                 logger,
		usernameHash:           usernameHash,
		passwordHash:           passwordHash,
		authMethodsCallback:    authMethodsCallback,
		passwordVerifyCallback: passwordVerifyCallback,
	}, nil
}

// ServerHTTP handles the webhook request.
func (s *StandardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("%s %s", r.Method, r.URL.String())

	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))

	usernameMatch := subtle.ConstantTimeCompare(s.usernameHash[:], usernameHash[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(s.passwordHash[:], passwordHash[:]) == 1

	if !usernameMatch || !passwordMatch {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	if r.Method != "POST" {
		s.sendBadRequest(w, "Invalid method '%s', only POST is allowed", r.Method)

		return
	}

	action := r.Header.Get("X-Corbado-Action")
	if action == "" {
		s.sendBadRequest(w, "X-Corbado-Action missing or empty")

		return
	}

	if r.Body == nil {
		s.sendBadRequest(w, "Empty body, provide JSON request")

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	if len(body) == 0 {
		s.sendBadRequest(w, "Empty body, provide JSON request")

		return
	}

	switch action {
	case "authMethods":
		s.handleAuthMethods(w, body)

	case "passwordVerify":
		s.handlePasswordVerify(w, body)

	default:
		s.sendBadRequest(w, "Invalid action given in X-Corbado-Action header ('%s')", action)
	}
}

func (s *StandardHandler) handleAuthMethods(w http.ResponseWriter, body []byte) {
	req, err := authmethodsrequest.NewFromBody(body)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	if req.Data.Username == "" {
		s.sendBadRequest(w, "username must not be empty")

		return
	}

	status, err := s.authMethodsCallback(req.Data.Username)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	resp, err := authmethodsresponse.New("", status)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	s.sendJSON(w, resp)
}

func (s *StandardHandler) handlePasswordVerify(w http.ResponseWriter, body []byte) {
	req, err := passwordverifyrequest.NewFromBody(body)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	if req.Data.Username == "" {
		s.sendBadRequest(w, "username must not be empty")

		return
	}

	if req.Data.Password == "" {
		s.sendBadRequest(w, "password must not be empty")

		return
	}

	success, err := s.passwordVerifyCallback(req.Data.Username, req.Data.Password)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	resp, err := passwordverifyresponse.New("", success)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	s.sendJSON(w, resp)
}

func (s *StandardHandler) sendJSON(w http.ResponseWriter, value any) {
	marshaled, err := json.Marshal(value)
	if err != nil {
		s.sendInternalServerError(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(marshaled); err != nil {
		s.sendInternalServerError(w, errors.WithStack(err))

		return
	}
}

func (s *StandardHandler) sendBadRequest(w http.ResponseWriter, message string, args ...any) {
	w.WriteHeader(http.StatusBadRequest)

	if _, err := w.Write([]byte(fmt.Sprintf(message, args...))); err != nil {
		s.sendInternalServerError(w, errors.WithStack(err))

		return
	}
}

func (s *StandardHandler) sendInternalServerError(w http.ResponseWriter, err error) {
	s.logger.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}
