package organization

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	CreatedBy string `json:"created_by" binding:"required"` 
	Email string `json:"email" binding:"required email"`
	Passkey string `json:"passkey" binding:"required"`
}