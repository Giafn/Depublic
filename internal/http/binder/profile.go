package binder

type CreateProfileRequest struct {
	FullName    string `form:"full_name" validate:"required"`
	Gender      string `form:"gender" validate:"required,oneof=Laki-laki Perempuan"`
	DateOfBirth string `form:"date_of_birth" validate:"required"`
	PhoneNumber string `form:"phone_number" validate:"required"`
	City        string `form:"city" validate:"required"`
	Province    string `form:"province" validate:"required"`
}

type UpdateProfileRequest struct {
	FullName       string `form:"fullname"`
	Gender         string `form:"gender"`
	DateOfBirth    string `form:"date_of_birth"`
	PhoneNumber    string `form:"phone_number"`
	ProfilePicture string `form:"profile_picture"`
	City           string `form:"city"`
	Province       string `form:"province"`
}
