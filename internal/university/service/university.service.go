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
	 FetchApprovedLecturersInDepartment(ctx context.Context,deptId string)(uniResponse,string,error)
	 CreateVenue(ctx context.Context, venueData dto.CreateVenueDto)(uniResponse,string,error)
	 RetrieveAllVenues(ctx context.Context,uniId uuid.UUID)(uniResponse,string,error)
	 FetchCohortsForADepartment(ctx context.Context,deptId uuid.UUID)(uniResponse,string,error)
	FetchAllDepartmentsForAUni(ctx context.Context, uniId uuid.UUID)(uniResponse,string,error)
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
	finalUniData := make([]dto.UniversityResponse,0)
	uniData,err := uns.repo.RetrieveAllUniversities(ctx)
	for _,val := range uniData{
		finalUniData = append(finalUniData, dto.UniversityResponse{
			Id: val.UniversityID,
			Abbr: val.UniversityAbbr.String,
			Email: val.Email,
			Name: val.UniversityName,
			Logo: val.UniversityLogo.String,
			Website: val.Website.String,
			PhoneNumber: val.PhoneNumber,
			Address: val.UniversityAddr.String,
			CurrentSession: val.CurrentSession.String,
		})
	}

	if err != nil{
		uns.logger.Error("Error retrieving all universities","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The universities",
		Data: finalUniData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)(uniResponse,string,error){
	facData, err := uns.repo.RetrieveAllFaculties(ctx,uniId)
	finalFacData := make([]dto.FacultyResponse,0)

	for _,val := range facData{
		finalFacData = append(finalFacData, dto.FacultyResponse{
			Id: val.FacultyID,
			UniversityId: val.UniversityID,
			Name: val.FacultyName,
			Code: val.FacultyCode.String,
		})
	}
	if err != nil{
		uns.logger.Error("Error retrieving all faculties","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The faculties",
		Data: finalFacData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) RetrieveAllDepartments(ctx context.Context,param RetrieveAllDepartmentsDto)(uniResponse,string,error){
	deptData,err := uns.repo.RetrieveAllDepartments(ctx,sqlc.RetrieveDeptsForAFacultyParams{
		FacultyID: utils.StringToUUID(param.FacultyId),
		UniversityID: utils.StringToUUID(param.UniversityId),
	})
	finalDeptData := make([]dto.DepartmentResponse,0)
	for _,val := range deptData{
		finalDeptData = append(finalDeptData, dto.DepartmentResponse{
			Id: val.DepartmentID,
			UniversityId: val.UniversityID,
			FacultyId: val.FacultyID,
			Name: val.DepartmentName,
			Code: val.DepartmentCode.String,
			NumberOfLevels: string(val.NumberOfLevels),
		})
	}

	if err != nil{
		uns.logger.Error("Error retrieving all departments","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	uns.logger.Info("the departments","depts:",finalDeptData)
	return uniResponse{
		Message: "The departments",
		Data:finalDeptData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) FetchApprovedLecturersInDepartment(ctx context.Context,deptId string)(uniResponse,string,error){
	lectData,err := uns.repo.FetchApprovedLecturersInDepartment(ctx,utils.StringToUUID(deptId))
	if err != nil{
		uns.logger.Error("Error retrieving lecturers in a department","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	finalLectData := make([]dto.FetchApprovedLecturersInDepartmentResponse,0)
	for _,val := range lectData{
		finalLectData = append(finalLectData, dto.FetchApprovedLecturersInDepartmentResponse{
			LecturerId: val.LecturerID,
			LecturerFirstName: val.LecturerFirstName,
			LecturerLastName: val.LecturerLastName,
			LecturerMiddleName: val.LecturerMiddleName.String,
			LecturerEmail: val.LecturerEmail,
			LecturerProfilePic: val.LecturerProfilePic.String,
			WaitId: val.WaitID,
			AdditionalMessage: val.AdditionalMessage.String,
			Approved: val.Approved.Bool,
		})
	}

	uns.logger.Info("the lecturers","lecturer",finalLectData)
	return uniResponse{
		Message: "The lecturers",
		Data:finalLectData,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) CreateVenue(ctx context.Context, venueData dto.CreateVenueDto)(uniResponse,string,error){
	actualVenueData := sqlc.CreateVenueParams{
		VenueName: venueData.VenueName,
		VenueLongitude: utils.Float64ToNullFloat64(venueData.VenueLongitude),
		VenueLatitude: utils.Float64ToNullFloat64(venueData.VenueLatitude),
		Location: utils.StringToNullString(venueData.Location),
		VenueImage: utils.StringToNullString(venueData.VenueImage),
		Capacity: venueData.Capacity,
		UniversityID: venueData.UniversityId,
	}
	venueUnavailability := dto.VenueUnavailability{
		Reason: venueData.UnavailabilityReason,
		Day: venueData.UnavailabilityDay,
		StartTime: venueData.UnavailabilityStartTime,
		EndTime: venueData.UnavailabilityEndTime,
	}
	err := uns.repo.CreateVenue(ctx,actualVenueData,venueData.VenueType,venueData.TypeId,venueUnavailability)
	if err != nil{
		uns.logger.Error("error creating venue","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "Venue Created Successfully",
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (uns *uniService) FetchCohortsForADepartment(ctx context.Context,deptId uuid.UUID)(uniResponse,string,error){
	data,err := uns.repo.FetchCohortsForADepartment(ctx,deptId)
	if err != nil{
		uns.logger.Error("error fetching cohorts for a department","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The cohorts",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) RetrieveAllVenues(ctx context.Context,uniId uuid.UUID)(uniResponse,string,error){
	data,err := uns.repo.RetrieveAllVenues(ctx,uniId)
	if err != nil{
		uns.logger.Error("error fetching all venues","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The venues",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (uns *uniService) FetchAllDepartmentsForAUni(ctx context.Context, uniId uuid.UUID)(uniResponse,string,error){
	data,err := uns.repo.FetchAllDepartmentsForAUni(ctx,uniId)
	if err != nil{
		uns.logger.Error("error fetching all departments","err:",err)
		return uniResponse{},status.InternalServerError.Message,err
	}
	return uniResponse{
		Message: "The departments",
		Data: data,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}