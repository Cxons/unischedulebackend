package service

import (
	"context"
	"errors"
	"log/slog"

	regDto "github.com/Cxons/unischedulebackend/internal/registration/dto"
	"github.com/Cxons/unischedulebackend/internal/registration/repository"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	sharedDto "github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)


type RegResponse = sharedDto.ResponseDto
type RegRepo = repository.RegRepository
type UpdateAdminDto = regDto.UpdateAdminDto
type CreateUniversityDto = regDto.CreateUniversityDto
type RequestDeanConfirmationDto = regDto.RequestDeanConfirmationDto
type RequestHodConfirmationDto = regDto.RequestHodConfirmationDto
type RequestLecturerConfirmationDto = regDto.RequestLecturerConfirmationDto

type RegService interface{
	UpdateAdmin(ctx context.Context,adminInfo UpdateAdminDto)(RegResponse,string,error)
	CreateUniversity(ctx context.Context, uniInfo CreateUniversityDto)(RegResponse,string,error)
	RetrievePendingDeans(ctx context.Context, uniId string)(RegResponse,string,error)
}



type regService struct{
	regRepo RegRepo
	logger *slog.Logger
}

func NewRegService(repo RegRepo,logger *slog.Logger)*regService{
	return &regService{
		regRepo: repo,
		logger: logger,
	}
}


func (rs *regService) UpdateAdmin(ctx context.Context,adminInfo UpdateAdminDto)(RegResponse,string,error){
	admin := sqlc.UpdateAdminInfoParams{
		AdminMiddleName: utils.StringToNullString(adminInfo.MiddleName),
		AdminPhoneNumber: utils.StringToNullString(adminInfo.PhoneNumber),
		AdminStaffCard: utils.StringToNullString(adminInfo.StaffCard),
		AdminNumber: utils.StringToNullString(adminInfo.AdminNumber),
		UniversityID: utils.StringToNullUUID(adminInfo.UniversityId),
		AdminID: utils.StringToUUID(adminInfo.AdminId),
	}
	
	adminExists,_,err := rs.regRepo.UpdateAdmin(ctx,admin)

	if !adminExists{
		return RegResponse{},status.NotFound.Message,err
	}

	if err != nil{
		rs.logger.Error("error updating admin details","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Admin updated successfully",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (rs *regService) CreateUniversity(ctx context.Context, uniInfo CreateUniversityDto)(RegResponse,string,error){
	university := sqlc.CreateUniversityParams{
		UniversityName: uniInfo.UniName,
		UniversityLogo: utils.StringToNullString(uniInfo.UniLogo),
		UniversityAbbr: utils.StringToNullString(uniInfo.UniAbbr),
		Email: uniInfo.UniEmail,
		Website: utils.StringToNullString(uniInfo.UniWebsite),
		PhoneNumber: uniInfo.UniPhoneNumber,
		CurrentSession: utils.StringToNullString(uniInfo.CurrentSession),
	}

	uni,err := rs.regRepo.CreateUniversity(ctx,university)
	
	if err != nil {
		rs.logger.Error("Error creating university","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}

	return RegResponse{
		Message: "University created successfully",
		Data:uni,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}



func (rs *regService) RetrievePendingDeans(ctx context.Context, uniId string)(RegResponse,string,error){
	list,err := rs.regRepo.RetrievePendingDeans(ctx,utils.StringToUUID(uniId))

	if err != nil{
		rs.logger.Error("Error retrieving pending deans","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if len(list) == 0 {
		return RegResponse{},status.NotFound.Message,errors.New("no pending deans found")
	}
	return RegResponse{
		Message: "Here are the pending deans",
		Data: list,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *regService) RequestDeanConfirmation(ctx context.Context, dean RequestDeanConfirmationDto )(RegResponse,string,error){

	deanInfo := sqlc.RequestDeanConfirmationParams{
		LecturerID: utils.StringToUUID(dean.LecturerId),
		PotentialFaculty: dean.PotentialFaculty,
		AdditionalMessage: utils.StringToNullString(dean.AdditionalMessage),
		UniversityID: utils.StringToUUID(dean.UniversityId),
	}
	_,err := rs.regRepo.RequestDeanConfirmation(ctx,deanInfo)

	if err != nil{
		rs.logger.Error("Error requesting dean confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Dean request confirmation sent",
		Data: nil,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil

}