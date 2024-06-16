package binder

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UserRegisterRequest struct {
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required"`
	FullName       string `form:"full_name" validate:"required"`
	Gender         string `form:"gender" validate:"required,oneof=Laki-laki Perempuan"`
	DateOfBirth    string `form:"date_of_birth" validate:"required,date"`
	PhoneNumber    string `form:"phone_number" validate:"required"`
	City           string `form:"city" validate:"required"`
	Province       string `form:"province" validate:"required"`
	ProfilePicture string `form:"profile_picture"`
}

type UserCreateRequest struct {
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required"`
	Role           string `form:"role" validate:"required,oneof=User Admin PetugasLapangan"`
	FullName       string `form:"full_name" validate:"required"`
	Gender         string `form:"gender" validate:"required,oneof=Laki-laki Perempuan"`
	DateOfBirth    string `form:"date_of_birth" validate:"required,date"`
	PhoneNumber    string `form:"phone_number" validate:"required"`
	City           string `form:"city" validate:"required"`
	Province       string `form:"province" validate:"required"`
	ProfilePicture string `form:"profile_picture"`
}

type UserUpdateRequest struct {
	ID       string `param:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type UserDeleteRequest struct {
	ID string `param:"id" validate:"required"`
}

type UserFindByIDRequest struct {
	ID string `param:"id" validate:"required"`
}

type UserVerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
