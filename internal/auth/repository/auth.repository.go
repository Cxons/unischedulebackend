package repository

import (
	"context"
	"database/sql"
	"log/slog"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)

type AuthRepository interface{
 RegisterStudent(ctx context.Context,student sqlc.RegisterStudentParams)(sqlc.Student,error)
 RegisterLecturer(ctx context.Context, lecturer sqlc.RegisterLecturerParams)(sqlc.Lecturer,error)
 RegisterUniversityAdmin(ctx context.Context, admin sqlc.RegisterUniversityAdminParams)(sqlc.UniversityAdmin,error) 
 RetrieveStudentEmail(ctx context.Context,studentEmail string)(bool,sqlc.RetrieveStudentEmailRow,error)
 RetrieveLecturerEmail(ctx context.Context,lecturerEmail string)(bool,sqlc.RetrieveLecturerEmailRow,error)
 RetrieveAdminEmail(ctx context.Context,adminEmail string)(bool,sqlc.RetrieveAdminEmailRow,error)
 AddRefreshToken(ctx context.Context,token sqlc.AddRefreshTokenParams)(sqlc.RefreshToken,error)
 UpdateRefreshToken(ctx context.Context,token sqlc.UpdateRefreshTokenParams)(sqlc.RefreshToken,error)
 CheckAndReturnToken(ctx context.Context,userId uuid.UUID)(bool,sqlc.CheckAndReturnTokenRow,error)
 DeleteRefreshToken(ctx context.Context,userId uuid.UUID)(sqlc.RefreshToken,error)
 InsertOtp(ctx context.Context,otpInfo sqlc.InsertOtpParams)(sqlc.Otp,error)
 RetrieveOtp(ctx context.Context,email string)(bool,sqlc.RetrieveOtpRow,error)
 UpdateOtp(ctx context.Context,otpInfo sqlc.UpdateOtpParams)(sqlc.Otp,error)
}

type authRepository struct{
	lq *queries.LecturerQueries
	sq *queries.StudentQueries
	aq *queries.AdminQueries
	tq *queries.TokenQueries
	oq *queries.OtpQueries
	logger *slog.Logger
}

func NewAuthRepository(sq *queries.StudentQueries, lq *queries.LecturerQueries, aq *queries.AdminQueries, tq *queries.TokenQueries,oq *queries.OtpQueries, logger *slog.Logger) *authRepository {
	return &authRepository{
		lq : lq,
		sq : sq,
		aq : aq,
		tq : tq,
		oq:oq,
		logger: logger,
	}
}

func (r *authRepository) RegisterStudent(ctx context.Context, student sqlc.RegisterStudentParams)(sqlc.Student,error){
	return r.sq.RegisterStudent(ctx,student)
}

func (r *authRepository) RetrieveStudentEmail(ctx context.Context,studentEmail string )(bool,sqlc.RetrieveStudentEmailRow,error){
	student,err := r.sq.RetrieveStudentEmail(ctx,studentEmail)
	if err != nil {
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveStudentEmailRow{},nil
		}
		r.logger.Error("Error retrieving student mail","err:",err)
		return false,sqlc.RetrieveStudentEmailRow{},err
	}
	return true,student,nil
}

func (r *authRepository) RegisterLecturer(ctx context.Context, lecturer sqlc.RegisterLecturerParams)(sqlc.Lecturer,error){
	return r.lq.RegisterLecturer(ctx,lecturer)
}

func (r *authRepository) RetrieveLecturerEmail(ctx context.Context, lecturerEmail string)(bool,sqlc.RetrieveLecturerEmailRow,error){
	lecturer,err :=  r.lq.RetrieveLecturerEmail(ctx,lecturerEmail)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveLecturerEmailRow{},nil
		}
		r.logger.Error("Error retrieving lecturer mail","err:",err)
		return false,sqlc.RetrieveLecturerEmailRow{},err
	}
	return true,lecturer,nil
}

func(r *authRepository) RegisterUniversityAdmin(ctx context.Context, admin sqlc.RegisterUniversityAdminParams)(sqlc.UniversityAdmin,error){
	return r.aq.RegisterAdmin(ctx,admin)
}

func (r *authRepository) RetrieveAdminEmail(ctx context.Context, email string)(bool,sqlc.RetrieveAdminEmailRow,error){
	admin,err := r.aq.RetrieveAdminEmail(ctx,email)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveAdminEmailRow{},nil
		}
		r.logger.Error("Error retrieving admin mail","err:",err)
		return false,sqlc.RetrieveAdminEmailRow{},err
	}
	return true,admin,nil
}
func (r *authRepository) AddRefreshToken(ctx context.Context,token sqlc.AddRefreshTokenParams)(sqlc.RefreshToken,error){
	return r.tq.AddRefreshToken(ctx,token)
}
func (r *authRepository) UpdateRefreshToken(ctx context.Context,token sqlc.UpdateRefreshTokenParams)(sqlc.RefreshToken,error){
	return r.tq.UpdateRefreshToken(ctx,token)
}
func (r * authRepository) CheckAndReturnToken(ctx context.Context,userId uuid.UUID)(bool,sqlc.CheckAndReturnTokenRow,error){
	token,err := r.tq.CheckAndReturnToken(ctx,userId)
	if err != nil {
		if err == sql.ErrNoRows{
			return false,sqlc.CheckAndReturnTokenRow{},nil
		}
		r.logger.Error("Error retrieving refresh token","err:",err)
		return false,sqlc.CheckAndReturnTokenRow{},nil
	}
	return true,token,nil
}
func (r *authRepository) DeleteRefreshToken(ctx context.Context,userId uuid.UUID)(sqlc.RefreshToken,error){
	return r.tq.DeleteRefreshToken(ctx,userId)
}

func (r *authRepository) InsertOtp(ctx context.Context,otpInfo sqlc.InsertOtpParams)(sqlc.Otp,error){
	return r.oq.InsertOtp(ctx,otpInfo)
}

func (r *authRepository) RetrieveOtp(ctx context.Context, email string)(bool,sqlc.RetrieveOtpRow,error){
	otp,err := r.oq.RetrieveOtp(ctx,email)
	if err != nil {
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveOtpRow{},nil
		}
		r.logger.Error("Error retrieving otp with email","err:",err)
		return false,sqlc.RetrieveOtpRow{},err
	}
	return true,otp,nil
}
func (r *authRepository) UpdateOtp(ctx context.Context, otpInfo sqlc.UpdateOtpParams)(sqlc.Otp,error){
	return r.oq.UpdateOtp(ctx,otpInfo)
}