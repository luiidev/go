package requests

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
}
