package env

type CreateEnvRequest struct {
	Name        string            `json:"name" binding:"required"`
	ProjectID   string            `json:"project_id" binding:"required"`
	Variables   map[string]string `json:"variables" binding:"required"`
	Version     int               `json:"version" binding:"required"`
	Description string            `json:"description" binding:"required"`
	CreatedBy   string            `json:"created_by" binding:"required"`
}
