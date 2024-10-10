package requests

type StoreUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
	Email     string `json:"email" gorm:"unique" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"required,min=8"`
}
