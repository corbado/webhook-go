package authmethodsrequest

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
}

// NewFromBody returns new request DTO for 'authMethod' action from given body.
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

	if dto.Action != "authMethods" {
		validationErrors = append(validationErrors, "field 'action' must be 'authMethods'")
	}

	if dto.Data.Username == "" {
		validationErrors = append(validationErrors, "field 'data.username' is empty")
	}

	if len(validationErrors) > 0 {
		return nil, errors.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	return dto, nil
}
