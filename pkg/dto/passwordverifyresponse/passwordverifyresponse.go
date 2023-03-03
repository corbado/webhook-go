package passwordverifyresponse

type DTO struct {
	ResponseID string   `json:"responseID"`
	Data       *DTOData `json:"data"`
}

type DTOData struct {
	Success bool `json:"success"`
}

// New returns new response DTO for 'passwordVerify' action with given responseID and status. The responseID
// is used for debugging and can be set to any value you like. It will show up in the webhook log in the
// developer panel.
func New(responseID string, success bool) (*DTO, error) {
	return &DTO{
		ResponseID: responseID,
		Data: &DTOData{
			Success: success,
		},
	}, nil
}
