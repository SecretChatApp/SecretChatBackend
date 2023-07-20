package authservices

type UserInput struct {
	Name            string `json:"name" validate:"required,gte=4"`
	Email           string `json:"email" validate:"required,email,isunique=users-email"`
	Password        string `json:"password" validate:"required,gte=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
