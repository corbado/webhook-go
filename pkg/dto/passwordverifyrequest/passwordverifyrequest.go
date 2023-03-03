package passwordverifyrequest

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type DTO struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"projectID"`
	Action    string   `json:"action"`
	Data      *DTOData `json:"data"`
}

type DTOData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewFromBody returns new request DTO for 'passwordVerify' action from given body.
func NewFromBody(body []byte) (*DTO, error) {
	if len(body) == 0 {
		return nil, errors.New("passed empty body")
	}

	dto := &DTO{
		Data: &DTOData{},
	}
	if err := json.Unmarshal(body, dto); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal() failed")
	}

	validationErrors := make([]string, 0, 5)
	if dto.ID == "" {
		validationErrors = append(validationErrors, "field 'id' is empty")
	}

	if dto.ProjectID == "" {
		validationErrors = append(validationErrors, "field 'projectID' is empty")
	}

	if dto.Action != "passwordVerify" {
		validationErrors = append(validationErrors, "field 'action' must be 'passwordVerify'")
	}

	if dto.Data.Username == "" {
		validationErrors = append(validationErrors, "field 'data.username' is empty")
	}

	if dto.Data.Password == "" {
		validationErrors = append(validationErrors, "field 'data.password' is empty")
	}

	if len(validationErrors) > 0 {
		return nil, errors.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	return dto, nil
}
