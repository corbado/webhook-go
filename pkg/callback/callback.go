package callback

import "github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"

type AuthMethods func(username string) (authmethodsresponse.Status, error)
type PasswordVerify func(username string, password string) (bool, error)
