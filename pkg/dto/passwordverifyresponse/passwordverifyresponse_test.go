package passwordverifyresponse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/corbado/webhook-go/pkg/dto/passwordverifyresponse"
)

func TestNew(t *testing.T) {
	dto, err := passwordverifyresponse.New("d5a80602-a771-4532-8cc8-6d4a9003d92a", true)
	assert.NoError(t, err)
	assert.Equal(t, "d5a80602-a771-4532-8cc8-6d4a9003d92a", dto.ResponseID)
	assert.True(t, dto.Data.Success)
}
