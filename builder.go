package corbado

import (
	"github.com/pkg/errors"

	"github.com/corbado/webhook-go/pkg/callback"
	"github.com/corbado/webhook-go/pkg/logger"
)

type Builder struct {
	logger                 logger.Logger
	username               string
	password               string
	authMethodsCallback    callback.AuthMethods
	passwordVerifyCallback callback.PasswordVerify
}

// NewBuilder returns new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// SetLogger sets given logger
func (b *Builder) SetLogger(logger logger.Logger) *Builder {
	b.logger = logger

	return b
}

// SetUsername sets given username on builder.
func (b *Builder) SetUsername(username string) *Builder {
	b.username = username

	return b
}

// SetPassword sets given password on builder.
func (b *Builder) SetPassword(password string) *Builder {
	b.password = password

	return b
}

// SetAuthMethodsCallback sets given callback on builder.
func (b *Builder) SetAuthMethodsCallback(authMethodsCallback callback.AuthMethods) *Builder {
	b.authMethodsCallback = authMethodsCallback

	return b
}

// SetPasswordVerifyCallback sets given callback on builder.
func (b *Builder) SetPasswordVerifyCallback(passwordVerifyCallback callback.PasswordVerify) *Builder {
	b.passwordVerifyCallback = passwordVerifyCallback

	return b
}

// Build builds a webhook instance, first validating all given parameters.
func (b *Builder) Build() (Webhook, error) {
	if b.logger == nil {
		return nil, errors.New("logger cannot be empty, call SetLogger() with logger")
	}

	if b.username == "" {
		return nil, errors.New("username cannot be empty, call SetUsername() with username")
	}

	if b.password == "" {
		return nil, errors.New("password cannot be empty, call SetPassword() with password")
	}

	if b.authMethodsCallback == nil {
		return nil, errors.New("authMethodsCallback cannot be empty, call SetAuthMethodsCallback() with callback")
	}

	if b.passwordVerifyCallback == nil {
		return nil, errors.New("passwordVerifyCallback cannot be empty, call SetPasswordVerifyCallback() with callback")
	}

	return New(b.logger, b.username, b.password, b.authMethodsCallback, b.passwordVerifyCallback)
}
