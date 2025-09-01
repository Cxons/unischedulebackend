package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	authDto "github.com/Cxons/unischedulebackend/internal/auth/dto"
	repo "github.com/Cxons/unischedulebackend/internal/auth/repository"
	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	sharedDto "github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
	"github.com/Cxons/unischedulebackend/pkg/mail"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type AuthRepository = repo.AuthRepository
type AuthResponse = sharedDto.ResponseDto
type LoginDto = authDto.LoginDto
type RegisterDto = authDto.RegisterDto
type LoginResponseData = authDto.LoginResponseData
type RefreshAccessTokenData = authDto.RefreshAccessTokenData
var STUDENT = constants.STUDENT
var ADMIN = constants.ADMIN
var LECTURER = constants.LECTURER
var ACCESS_TOKEN = constants.ACCESS_TOKEN
var REFRESH_TOKEN = constants.REFRESH_TOKEN

type AuthService interface{
 	Register(ctx context.Context,user RegisterDto)(AuthResponse,string,error)
 	Login(ctx context.Context,user LoginDto)(AuthResponse,string,error)
 	RefreshAccessToken(ctx context.Context,refreshToken *string)(AuthResponse,string,error)
 	SendOtp(ctx context.Context,email string,userType string)(AuthResponse,string,error)
 	VerifyOtp(ctx context.Context,email string,otpStr string)(AuthResponse,string,error)
}


type authService struct{
	repo AuthRepository
	logger *slog.Logger
}

func  NewAuthService(repo AuthRepository,logger *slog.Logger) *authService{
	return &authService{
		repo: repo,
		logger: logger,
	}
}

func (s *authService) Register(ctx context.Context,user RegisterDto)(AuthResponse,string,error){
	exists,_,_ :=s.retrieveUserEmail(ctx,user.Role,user.Email)
	if exists{
		return AuthResponse{},status.Unauthorized.Message,errors.New("email already exists")
	}
	hash,err := bcrypt.GenerateFromPassword([]byte(user.Password),12)
	if err != nil {
		return AuthResponse{},status.InternalServerError.Message,err
	}

	switch user.Role {
	case STUDENT:
		student,err := s.repo.RegisterStudent(ctx, sqlc.RegisterStudentParams{
			StudentFirstName: user.FirstName,
			StudentLastName: user.LastName,
			StudentEmail: user.Email,
			StudentPassword: string(hash),
		})
		if err != nil {
			s.logger.Error("Error registering student","err:",err)
			return AuthResponse{}, status.InternalServerError.Message, err
		}
		_,_,err = s.SendOtp(ctx,user.Email,user.Role)
		if err != nil {
			return AuthResponse{}, status.InternalServerError.Message, err
		}
		return AuthResponse{
			Message: "Thanks for registering " + student.StudentFirstName,
		},"",nil

	case LECTURER:
		lecturer,err := s.repo.RegisterLecturer(ctx,sqlc.RegisterLecturerParams{
			LecturerFirstName: user.FirstName,
			LecturerLastName: user.LastName,
			LecturerEmail: user.Email,
			LecturerPassword: string(hash),
		})
		if err != nil{
			s.logger.Error("error registering lecturer","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err
		}
		_,_,err = s.SendOtp(ctx,user.Email,user.Role)
		if err != nil {
			return AuthResponse{}, status.InternalServerError.Message, err
		}
		return AuthResponse{
			Message: "Thanks for registering " + lecturer.LecturerFirstName,
		},"",nil

	case ADMIN:
		admin,err := s.repo.RegisterUniversityAdmin(ctx,sqlc.RegisterUniversityAdminParams{
			AdminFirstName: user.FirstName,
			AdminLastName: user.LastName,
			AdminEmail: user.Email,
			AdminPassword: string(hash),
		})
		if err != nil{
			s.logger.Error("error registering admin","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err
		}
		_,_,err = s.SendOtp(ctx,user.Email,user.Role)
		if err != nil {
			return AuthResponse{}, status.InternalServerError.Message, err
		}
		return AuthResponse{
				Message: "Thanks for registering." + admin.AdminFirstName + "otp sent",
			},"",err

	default:
		return AuthResponse{},status.BadRequest.Message,errors.New("user type not allowed")

	}
}


func (s *authService) Login(ctx context.Context, user LoginDto)(AuthResponse,string,error){
	var password string;
	var userId uuid.UUID;
	var refreshToken string;
	var accessToken string;
	emailExists,userData,err := s.retrieveUserEmail(ctx,user.Role,user.Email)
	if err != nil {
			s.logger.Error("Error retrieving user email","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err
		}
	switch user.Role{
	case STUDENT: 
		userId = userData.(sqlc.RetrieveStudentEmailRow).StudentID
		password = userData.(sqlc.RetrieveStudentEmailRow).StudentPassword
	case LECTURER:
		userId = userData.(sqlc.RetrieveLecturerEmailRow).LecturerID
		password = userData.(sqlc.RetrieveLecturerEmailRow).LecturerPassword
	case ADMIN:
		userId = userData.(sqlc.RetrieveAdminEmailRow).AdminID
		password = userData.(sqlc.RetrieveAdminEmailRow).AdminPassword
	}
	_,err = s.verifyPassword(emailExists,user.Password,password)
	if err != nil {
			s.logger.Error("Error verifying user email and password","err:",err)
			return AuthResponse{},status.Unauthorized.Message,err
		}
	refreshToken,err = s.manageLoginRefreshToken(ctx,userId,user.Email,user.Role,REFRESH_TOKEN)
	if err != nil {
		s.logger.Error("Error managing login refresh tokens","err:",err)
		return AuthResponse{},status.InternalServerError.Message,err
	}
	accessToken,err = s.createAccessToken(userId.String(),user.Email,user.Role)
	if err != nil {
		s.logger.Error("Error creating access token","err:",err)
		return AuthResponse{},status.InternalServerError.Message,err
	}
	return AuthResponse{
		Message:"Login successful",
		Data: LoginResponseData{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},
	StatusCode: status.Created.Code,
	StatusCodeMessage: status.Created.Message,
	},"",nil
	
}

func (s *authService) retrieveUserEmail(ctx context.Context, role string, email string)(bool,interface{},error){
	switch role {
	case STUDENT:
		return s.repo.RetrieveStudentEmail(ctx,email)
	case LECTURER:
		return s.repo.RetrieveLecturerEmail(ctx,email)
	case ADMIN:
		return s.repo.RetrieveAdminEmail(ctx,email)
	default:
		return false,nil,errors.New("role type not allowed")
	}
}

// makes sure password and email check errors have the same response time
func (s *authService) verifyPassword(emailExistence bool , password string, hash string)(AuthResponse,error){
	hashErr := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
		if !emailExistence || hashErr != nil {
			return AuthResponse{},errors.New("incorrect email or password")
		}
		return AuthResponse{},nil
}



// handles all refresh token operations for login
func (s *authService) manageLoginRefreshToken(ctx context.Context,userId uuid.UUID,email string,role string,typeOfToken string)(string,error){
	// generate token
	token,expires_at,err :=jwt.GenerateToken(userId.String(),email,role,typeOfToken)
			if err != nil {
				return "",err
			}

	// checks if user with refresh token exists
	tokenExists,_,err := s.repo.CheckAndReturnToken(ctx,userId)
	if err != nil {
		return "",err
	}
	
	// if it does not exist add token to database
	if !tokenExists{
			
			_,err = s.repo.AddRefreshToken(ctx,sqlc.AddRefreshTokenParams{
				RefreshToken: token,
				ExpiresAt: expires_at.Time,
				UserID: userId,
			})
			if err != nil {
				return "",err
			}
			return token,nil
		}
		s.logger.Info("expires at","time",expires_at.Time,"currenttime",time.Now())

	// else update user refresh token with new token information
			_,err = s.repo.UpdateRefreshToken(ctx,sqlc.UpdateRefreshTokenParams{
				RefreshToken: token,
				ExpiresAt: expires_at.Time,
				UserID: userId,
			})

			if err != nil {
				return "",err
			}
			return token,nil
}

func (s *authService) createAccessToken(userId string,email string,role string)(string,error){
token,_,err :=jwt.GenerateToken(userId,email,role,ACCESS_TOKEN)
			if err != nil {
				s.logger.Error("Error generating refresh token","err:",err)
				return "",err
			}
		return token,nil
}

func (s *authService) RefreshAccessToken(ctx context.Context,refreshToken *string)(AuthResponse,string,error){
	claims,err := jwt.VerifyToken(*refreshToken)
	if err!= nil{
		if strings.Contains(err.Error(),"token is expired"){
		return AuthResponse{},status.Forbidden.Message,errors.New("otp expired")
		}
		s.logger.Error("Error verifying token","err:",err)
		return AuthResponse{},status.InternalServerError.Message,err
	}
	token,err := s.createAccessToken(claims.User_id,claims.Email,claims.Role)
	if err != nil{
		s.logger.Error("Error creating refresh access token","err:",err)
		return AuthResponse{},"",err
	}
	return AuthResponse{
		Message:"Access token generated",
		Data:RefreshAccessTokenData{
			AccessToken: token,
		},
		StatusCode:status.OK.Code,
		StatusCodeMessage:status.OK.Message,
	},"",nil
}


// generates a random 6 digit code and returns a string value of that code
func (s *authService)generateOTP(length int) (string, error) {
  otp := ""
    for i := 0; i < length; i++ {
        b := make([]byte, 1)
        _, err := rand.Read(b)
        if err != nil {
            return "", err
        }
        otp += fmt.Sprintf("%d", b[0]%10)
    }
    return otp, nil
}


func (s *authService) SendOtp(ctx context.Context,email string, userType string)(AuthResponse,string,error){
		otpCode,err := s.generateOTP(6)
		if err != nil{
			s.logger.Error("Error generating otp","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err
		}

	// formats code to fit string
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
		<h1>One Time Verification Code</h1>
		<p style="color:blue; font-size:25px;">Your otp code is %s</p>
		<p style="font-size:15px;">This code would expire after 5 minutes</p>
	</body>
	</html>
	`,otpCode)

	//sends otp email before an database operations
		msg,err := mail.SendMessage("njokuchukwuma48@gmail.com","OTP CODE",htmlBody,email);
		if err != nil{
			s.logger.Error("Error sending otp mail","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err;
		}
		s.logger.Info("Response from sending otp mail","res:",msg)

		// defines data formats for both updating and inserting otps
		otpInfo := sqlc.InsertOtpParams{
			Otp:otpCode,
			ExpiresAt: time.Now().Add(time.Minute * 5),
			Email:email,
			UserType:userType,
		}
		updateOtpInfo := sqlc.UpdateOtpParams{
			Otp: otpCode,
			ExpiresAt: time.Now().Add(time.Minute * 5),
			Email:email,
		}

		// checks if email with otp already exists
		otpExists,_,err := s.repo.RetrieveOtp(ctx,email)
		if err != nil{
			s.logger.Error("Error retrieving otp","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err;
		}

		// if exists update otp in database
		if otpExists{
			_,err = s.repo.UpdateOtp(ctx,updateOtpInfo)
			if err != nil{
			s.logger.Error("Error updating otp into db","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err;
		}
		}else {
		// if not insert otp into database
		_,err = s.repo.InsertOtp(ctx,otpInfo);
		if err != nil{
			s.logger.Error("Error inserting otp into db","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err;
		}
		}
		
		return AuthResponse{
			Message: "Otp sent successfully",
			Data: nil,
			StatusCode: status.Created.Code,
			StatusCodeMessage: status.Created.Message,
		},"",err
	}

func (s *authService) VerifyOtp(ctx context.Context,email string,otpStr string)(AuthResponse,string,error){
	// checks if email with otp already exists
		_,otp,err := s.repo.RetrieveOtp(ctx,email)
		if err != nil{
			s.logger.Error("Error retrieving otp","err:",err)
			return AuthResponse{},status.InternalServerError.Message,err;
		}
	// rejects expired otps
		if time.Now().After(otp.ExpiresAt){
			return AuthResponse{},status.Forbidden.Message,errors.New("otp expired")
		}
		if otp.Otp != otpStr{
			return AuthResponse{},status.Forbidden.Message,errors.New("incorrect otp")
		}
		 return AuthResponse{
			Message: "Otp is correct",
			Data:nil,
			StatusCode: status.OK.Code,
			StatusCodeMessage: status.OK.Message,
		 },"",err
}




// func (s *authService) sendOtp

