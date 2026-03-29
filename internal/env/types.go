package env

import "encoding/json"

type CreateEnvRequest struct {
	Name        string          `json:"name"`
	ProjectID   string          `json:"project_id"`
	Variables   json.RawMessage `json:"variables"`
	Description string          `json:"description"`
}

type UpdateEnvRequest struct {
	Name        *string          `json:"name"`
	Variables   *json.RawMessage `json:"variables"`
	Description *string          `json:"description"`
}
