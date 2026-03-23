package project

type CreateProjectRequest struct {
	Name           string `json:"name" binding:"required"`
	Passkey        string `json:"passkey" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
	CreatedBy      string `json:"created_by" binding:"required"`
}

type UpdateProjectRequest struct {
	Name    *string `json:"name"`
	Passkey *string `json:"passkey"`
}
