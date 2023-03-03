package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	corbado "github.com/corbado/webhook-go"
	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
	"github.com/corbado/webhook-go/pkg/logger"
)

const addr = "localhost:8000"
const webhookUsername = "corbado"
const webhookPassword = "#73KojdPn,f4XksW_]^N"

func main() {
	webhook, err := corbado.
		NewBuilder().
		SetLogger(logger.New()).
		SetUsername(webhookUsername).
		SetPassword(webhookPassword).
		SetAuthMethodsCallback(authMethodsCallback).
		SetPasswordVerifyCallback(passwordVerifyCallback).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	handler, err := webhook.GetGinHandler()
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/corbadoWebhook", handler.Handle)

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// authMethodsCallback is being executed for webhook action "authMethods".
// !!! IMPLEMENT YOUR OWN LOGIC HERE !!!
func authMethodsCallback(_ string) (authmethodsresponse.Status, error) {
	return authmethodsresponse.StatusExists, nil
}

// passwordVerifyCallback is being executed for webhook action "passwordVerify".
// !!! IMPLEMENT YOUR OWN LOGIC HERE !!!
func passwordVerifyCallback(_ string, _ string) (bool, error) {
	return false, nil
}
