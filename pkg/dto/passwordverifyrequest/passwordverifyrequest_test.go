package passwordverifyrequest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/corbado/webhook-go/pkg/dto/passwordverifyrequest"
)

func TestNewFromBody(t *testing.T) {
	dto, err := passwordverifyrequest.NewFromBody(nil)
	assert.ErrorContains(t, err, "passed empty body")
	assert.Nil(t, dto)

	dto, err = passwordverifyrequest.NewFromBody(readTestDataJSON("broken"))
	assert.ErrorContains(t, err, "invalid character")
	assert.Nil(t, dto)

	dto, err = passwordverifyrequest.NewFromBody(readTestDataJSON("invalid"))
	assert.ErrorContains(t, err,
		"validation failed: field 'id' is empty, field 'projectID' is empty, field 'action' must be 'passwordVerify', field 'data.username' is empty, field 'data.password' is empty")
	assert.Nil(t, dto)

	dto, err = passwordverifyrequest.NewFromBody(readTestDataJSON("valid"))
	assert.NoError(t, err)
	assert.Equal(t, "who-1234567890", dto.ID)
	assert.Equal(t, "pro-1234567890", dto.ProjectID)
	assert.Equal(t, "passwordVerify", dto.Action)
	assert.Equal(t, "testUsername", dto.Data.Username)
	assert.Equal(t, "testPassword", dto.Data.Password)
}

func readTestDataJSON(name string) []byte {
	if name == "" {
		panic("given name is empty")
	}

	data, err := os.ReadFile(fmt.Sprintf("testdata/%s.json", name))
	if err != nil {
		panic(err)
	}

	return data
}
