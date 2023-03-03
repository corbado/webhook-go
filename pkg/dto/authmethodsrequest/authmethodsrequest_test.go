package authmethodsrequest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/corbado/webhook-go/pkg/dto/authmethodsrequest"
)

func TestNewFromBody(t *testing.T) {
	dto, err := authmethodsrequest.NewFromBody(nil)
	assert.ErrorContains(t, err, "passed empty body")
	assert.Nil(t, dto)

	dto, err = authmethodsrequest.NewFromBody(readTestDataJSON("broken"))
	assert.ErrorContains(t, err, "invalid character")
	assert.Nil(t, dto)

	dto, err = authmethodsrequest.NewFromBody(readTestDataJSON("invalid"))
	assert.ErrorContains(t, err, "validation failed: field 'id' is empty, field 'projectID' is empty, field 'action' must be 'authMethods', field 'data.username' is empty")
	assert.Nil(t, dto)

	dto, err = authmethodsrequest.NewFromBody(readTestDataJSON("valid"))
	assert.NoError(t, err)
	assert.Equal(t, "who-1234567890", dto.ID)
	assert.Equal(t, "pro-1234567890", dto.ProjectID)
	assert.Equal(t, "authMethods", dto.Action)
	assert.Equal(t, "testUsername", dto.Data.Username)
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
