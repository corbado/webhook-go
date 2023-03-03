package authmethodsresponse

import "github.com/pkg/errors"

type DTO struct {
	ResponseID string   `json:"responseID"`
	Data       *DTOData `json:"data"`
}

type DTOData struct {
	Status string `json:"status"`
}

type Status string

const (
	StatusExists    Status = "exists"
	StatusNotExists Status = "not_exists"
)

// New returns new response DTO for 'authMethods' action with given responseID and status. The responseID
// is used for debugging and can be set to any value you like. It will show up in the webhook log in the
// developer panel.
func New(responseID string, status Status) (*DTO, error) {
	if status != StatusExists && status != StatusNotExists {
		return nil, errors.Errorf("status must be either '%s' or '%s'", StatusExists, StatusNotExists)
	}

	return &DTO{
		ResponseID: responseID,
		Data: &DTOData{
			Status: string(status),
		},
	}, nil
}
