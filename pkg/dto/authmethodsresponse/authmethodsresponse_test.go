package authmethodsresponse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/corbado/webhook-go/pkg/dto/authmethodsresponse"
)

func TestNew(t *testing.T) {
	dto, err := authmethodsresponse.New("", "invalid")
	assert.ErrorContains(t, err, "status must be either 'exists' or 'not_exists'")
	assert.Nil(t, dto)

	dto, err = authmethodsresponse.New("d5a80602-a771-4532-8cc8-6d4a9003d92a", authmethodsresponse.StatusExists)
	assert.NoError(t, err)
	assert.Equal(t, "d5a80602-a771-4532-8cc8-6d4a9003d92a", dto.ResponseID)
	assert.Equal(t, "exists", dto.Data.Status)
}
