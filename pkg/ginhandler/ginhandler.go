package ginhandler

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/corbado/webhook-go/pkg/callback"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsrequest"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
	"github.com/corbado/webhook-go/pkg/dto/passwordverifyrequest"
	"github.com/corbado/webhook-go/pkg/dto/passwordverifyresponse"
	"github.com/corbado/webhook-go/pkg/logger"
)

type GinHandler struct {
	logger                 logger.Logger
	usernameHash           [32]byte
	passwordHash           [32]byte
	authMethodsCallback    callback.AuthMethods
	passwordVerifyCallback callback.PasswordVerify
}

// New returns Gin handler which can be used in Gin Web Framework.
func New(
	logger logger.Logger,
	usernameHash [32]byte,
	passwordHash [32]byte,
	authMethodsCallback callback.AuthMethods,
	passwordVerifyCallback callback.PasswordVerify,
) (*GinHandler, error) {
	if logger == nil {
		return nil, errors.New("empty parameter logger")
	}

	if authMethodsCallback == nil {
		return nil, errors.New("empty parameter authMethodsCallback")
	}

	if passwordVerifyCallback == nil {
		return nil, errors.New("empty parameter passwordVerifyCallback")
	}

	return &GinHandler{
		logger:                 logger,
		usernameHash:           usernameHash,
		passwordHash:           passwordHash,
		authMethodsCallback:    authMethodsCallback,
		passwordVerifyCallback: passwordVerifyCallback,
	}, nil
}

// Handle handles the Corbado webhook request.
func (g *GinHandler) Handle(c *gin.Context) {
	g.logger.Debug("%s %s", c.Request.Method, c.Request.URL.String())

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		c.String(http.StatusUnauthorized, "Unauthorized")

		return
	}

	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))

	usernameMatch := subtle.ConstantTimeCompare(g.usernameHash[:], usernameHash[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(g.passwordHash[:], passwordHash[:]) == 1

	if !usernameMatch || !passwordMatch {
		c.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		c.String(http.StatusUnauthorized, "Unauthorized")

		return
	}

	action := c.GetHeader("X-Corbado-Action")
	if action == "" {
		c.String(http.StatusBadRequest, "X-Corbado-Action header missing or empty")

		return
	}

	body, err := c.GetRawData()
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	if len(body) == 0 {
		c.String(http.StatusBadRequest, "Empty body, provide JSON request")

		return
	}

	switch action {
	case "authMethods":
		g.handleAuthMethods(c, body)

	case "passwordVerify":
		g.handlePasswordVerify(c, body)

	default:
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid action given in X-Corbado-Action header ('%s')", action))
	}
}

func (g *GinHandler) handleAuthMethods(c *gin.Context, body []byte) {
	req, err := authmethodsrequest.NewFromBody(body)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	if req.Data.Username == "" {
		c.String(http.StatusBadRequest, "username must not be empty")

		return
	}

	status, err := g.authMethodsCallback(req.Data.Username)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	resp, err := authmethodsresponse.New("", status)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *GinHandler) handlePasswordVerify(c *gin.Context, body []byte) {
	req, err := passwordverifyrequest.NewFromBody(body)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	if req.Data.Username == "" {
		c.String(http.StatusBadRequest, "username must not be empty")

		return
	}

	if req.Data.Password == "" {
		c.String(http.StatusBadRequest, "password must not be empty")

		return
	}

	success, err := g.passwordVerifyCallback(req.Data.Username, req.Data.Password)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	resp, err := passwordverifyresponse.New("", success)
	if err != nil {
		g.sendInternalServerError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (g *GinHandler) sendInternalServerError(c *gin.Context, err error) {
	g.logger.Error(err)
	c.Status(http.StatusInternalServerError)
}
