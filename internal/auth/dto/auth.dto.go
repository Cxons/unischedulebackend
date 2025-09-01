package dto

type RegisterDto struct {
 	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName string `json:"last_name" validate:"required,alpha"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Role string `json:"role" validate:"required,oneof=student lecturer admin"`
}

type LoginDto struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Role string `json:"role" validate:"required,oneof=student lecturer admin"`
}

type SendOtpDto struct {
	Email string `json:"email" validate:"required,email"`
	UserType string `json:"usertype" validate:"required,oneof=student lecturer admin"`
}

type VerifyOtpDto struct{
	Email string `json:"email" validate:"required,email"`
	Otp string `json:"otp" validate:"required"`
}


type LoginResponseData struct {
		AccessToken string
		RefreshToken string
}
type RefreshAccessTokenData struct {
	AccessToken string
}

// type RegisterLecturerDto struct {
// 	First_Name string `json:"first_name" validate:"required,alpha"`
// 	Last_Name string `json:"last_name" validate:"required,alpha"`
// 	Email string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8,max=64"`
// }
// type RegisterUniversityAdminDto struct {
// 	First_Name string `json:"first_name" validate:"required,alpha"`
// 	Last_Name string `json:"last_name" validate:"required,alpha"`
// 	Email string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8,max=64"`
// }