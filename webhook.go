package corbado

import (
	"crypto/sha256"

	"github.com/pkg/errors"

	"github.com/corbado/webhook-go/pkg/callback"
	"github.com/corbado/webhook-go/pkg/ginhandler"
	"github.com/corbado/webhook-go/pkg/logger"
	"github.com/corbado/webhook-go/pkg/standardhandler"
)

type Webhook interface {
	GetStandardHandler() (*standardhandler.StandardHandler, error)
	GetGinHandler() (*ginhandler.GinHandler, error)
}

type Impl struct {
	logger                 logger.Logger
	usernameHash           [32]byte
	passwordHash           [32]byte
	authMethodsCallback    callback.AuthMethods
	passwordVerifyCallback callback.PasswordVerify
}

var _ Webhook = &Impl{}

// New returns new webhook instance.
func New(
	logger logger.Logger,
	username string,
	password string,
	authMethodsCallback callback.AuthMethods,
	passwordVerifyCallback callback.PasswordVerify,
) (*Impl, error) {
	if logger == nil {
		return nil, errors.New("empty parameter logger")
	}

	if username == "" {
		return nil, errors.New("empty parameter username")
	}

	if password == "" {
		return nil, errors.New("empty parameter password")
	}

	if authMethodsCallback == nil {
		return nil, errors.New("empty parameter authMethodsCallback")
	}

	if passwordVerifyCallback == nil {
		return nil, errors.New("empty parameter passwordVerifyCallback")
	}

	return &Impl{
		logger:                 logger,
		usernameHash:           sha256.Sum256([]byte(username)),
		passwordHash:           sha256.Sum256([]byte(password)),
		authMethodsCallback:    authMethodsCallback,
		passwordVerifyCallback: passwordVerifyCallback,
	}, nil
}

// GetStandardHandler returns standard handler which can be used in standard HTTP library.
func (i *Impl) GetStandardHandler() (*standardhandler.StandardHandler, error) {
	return standardhandler.New(
		i.logger,
		i.usernameHash,
		i.passwordHash,
		i.authMethodsCallback,
		i.passwordVerifyCallback,
	)
}

// GetGinHandler returns Gin handler which can be used in Gin Web Framework.
func (i *Impl) GetGinHandler() (*ginhandler.GinHandler, error) {
	return ginhandler.New(
		i.logger,
		i.usernameHash,
		i.passwordHash,
		i.authMethodsCallback,
		i.passwordVerifyCallback,
	)
}
