package service

import (
	"context"
	"log/slog"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	sharedDto "github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/university/dto"
	"github.com/Cxons/unischedulebackend/internal/university/repository"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/google/uuid"
)


type uniResponse = sharedDto.ResponseDto
type RetrieveAllDepartmentsDto = dto.RetrieveAllDepartmentsDto

type UniService interface{
	 RetrieveAllUniversities(ctx context.Context)(uniResponse,string,error)
	 RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)(uniResponse,string,error)
	  RetrieveAllDepartments(ctx context.Context,param RetrieveAllDepartmentsDto)(uniResponse,string,error)
}

type uniService struct{
	repo repository.UniRepository
	logger *slog.Logger
}


func NewUniService(repo repository.UniRepository,logger *slog.Logger)*uniService{
	return &uniService{
		repo: repo,
		logger: logger,
	}
}


func (uns *uniService) RetrieveAllUniversities(ctx context.Context)(uniResponse,string,error){
	uniData,err := uns.repo.RetrieveAllUniversities(ctx)

	if err != nil{
		uns.logger.Error("Error retrieving all universities","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The universities",
		Data: uniData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)(uniResponse,string,error){
	facData, err := uns.repo.RetrieveAllFaculties(ctx,uniId)
	if err != nil{
		uns.logger.Error("Error retrieving all faculties","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The faculties",
		Data: facData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) RetrieveAllDepartments(ctx context.Context,param RetrieveAllDepartmentsDto)(uniResponse,string,error){
	deptData,err := uns.repo.RetrieveAllDepartments(ctx,sqlc.RetrieveDeptsForAFacultyParams{
		FacultyID: utils.StringToUUID(param.FacultyId),
		UniversityID: utils.StringToUUID(param.UniversityId),
	})

	if err != nil{
		uns.logger.Error("Error retrieving all departments","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The departments",
		Data:deptData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}