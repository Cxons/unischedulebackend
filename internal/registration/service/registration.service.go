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
	"github.com/Cxons/unischedulebackend/pkg/auth/jwt"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
)


type RegResponse = sharedDto.ResponseDto
type RegRepo = repository.RegRepository
type UpdateAdminDto = regDto.UpdateAdminDto
type CreateUniversityDto = regDto.CreateUniversityDto
type CreateFacultyDto = regDto.CreateFacultyDto
type CreateDepartmentDto = regDto.CreateDepartmentDto
type RequestDeanConfirmationDto = regDto.RequestDeanConfirmationDto
type RequestHodConfirmationDto = regDto.RequestHodConfirmationDto
type RequestLecturerConfirmationDto = regDto.RequestLecturerConfirmationDto
type PendingHodDto = regDto.PendingHodDto
type PendingLecturerDto = regDto.PendingLecturerDto
type CreateDeanDto = regDto.CreateDeanDto
type CreateHodDto = regDto.CreateHodDto
var userInfoKey = &sharedDto.UserContext{Value: "UserInfo"}

type RegService interface{
	UpdateAdmin(ctx context.Context,adminInfo UpdateAdminDto)(RegResponse,string,error)
	CreateUniversity(ctx context.Context, uniInfo CreateUniversityDto)(RegResponse,string,error)
	ApproveDean(ctx context.Context,waitId string)(RegResponse,string,error)
	RetrievePendingDeans(ctx context.Context, uniId string)(RegResponse,string,error)
	RequestDeanConfirmation(ctx context.Context, dean RequestDeanConfirmationDto )(RegResponse,string,error)
	CheckDeanConfirmation(ctx context.Context, waitId string)(RegResponse,string,error)
	RetrievePendingHods(ctx context.Context, hodParams PendingHodDto )(RegResponse,string,error)
	ApproveHod(ctx context.Context,waitId string)(RegResponse,string,error)
	RequestHodConfirmation(ctx context.Context, hod RequestHodConfirmationDto)(RegResponse,string,error)
	CheckHodConfirmation(ctx context.Context, waitId string)(RegResponse,string,error)
	RetrievePendingLecturers(ctx context.Context, lecturerParams PendingLecturerDto)(RegResponse,string,error)
	ApproveLecturer(ctx context.Context,waitId string)(RegResponse,string,error)
	RequestLecturerConfirmation(ctx context.Context, lecturer RequestLecturerConfirmationDto)(RegResponse,string,error)
	CheckLecturerConfirmation(ctx context.Context, waitId string)(RegResponse,string,error)
	CreateDepartment(ctx context.Context, deptInfo CreateDepartmentDto)(RegResponse,string,error)
	CreateFaculty(ctx context.Context, facInfo CreateFacultyDto)(RegResponse,string,error)
	CreateDean(ctx context.Context, deanInfo CreateDeanDto )(RegResponse,string,error)
	CreateHod(ctx context.Context, hodInfo CreateHodDto)(RegResponse,string,error)
	CheckCurrentDean(ctx context.Context,deanId string)(bool,error)
	CheckCurrentHod(ctx context.Context,hodId string)(bool,error)
	CheckCurrentAdmin(ctx context.Context,adminId string)(bool,error)
}



type RegServiceStruct struct{
	regRepo RegRepo
	logger *slog.Logger
}

func NewRegService(repo RegRepo,logger *slog.Logger)*RegServiceStruct{
	return &RegServiceStruct{
		regRepo: repo,
		logger: logger,
	}
}


func (rs *RegServiceStruct) UpdateAdmin(ctx context.Context,adminInfo UpdateAdminDto)(RegResponse,string,error){
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

func (rs *RegServiceStruct) CreateUniversity(ctx context.Context, uniInfo CreateUniversityDto)(RegResponse,string,error){
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



func (rs *RegServiceStruct) RetrievePendingDeans(ctx context.Context, uniId string)(RegResponse,string,error){
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

func (rs *RegServiceStruct) ApproveDean(ctx context.Context,waitId string)(RegResponse,string,error){
	idExists,_,err := rs.regRepo.ApproveDean(ctx,utils.StringToUUID(waitId))

	if err != nil{
		if !idExists{
		return RegResponse{},status.NotFound.Message,err
	}
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Dean approved",
		Data: nil,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}


func (rs *RegServiceStruct) RequestDeanConfirmation(ctx context.Context, dean RequestDeanConfirmationDto )(RegResponse,string,error){
	
	var lecturerId string;

	// retrieves user id from token claims
	claims := ctx.Value(userInfoKey)
	if claims != nil{
		lecturerId = claims.(jwt.CustomClaims).User_id
	}else{
		return RegResponse{},status.InternalServerError.Message,errors.New("problem getiing claims from token")
	}
	
	deanInfo := sqlc.RequestDeanConfirmationParams{
		LecturerID: utils.StringToUUID(lecturerId),
		PotentialFaculty: dean.PotentialFaculty,
		AdditionalMessage: utils.StringToNullString(dean.AdditionalMessage),
		UniversityID: utils.StringToUUID(dean.UniversityId),
	}
	waitInfo,err := rs.regRepo.RequestDeanConfirmation(ctx,deanInfo)

	desiredWaitInfo := sqlc.DeanWaitingList{
		WaitID: waitInfo.WaitID,
	}

	if err != nil{
		rs.logger.Error("Error requesting dean confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Dean request confirmation sent",
		Data: desiredWaitInfo,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}


func (rs *RegServiceStruct) CheckDeanConfirmation(ctx context.Context, waitId string)(RegResponse,string,error){
	confirmationInfo,err := rs.regRepo.CheckDeanConfirmation(ctx,utils.StringToUUID(waitId))
	if err != nil{
		rs.logger.Error("Error checking dean confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if !utils.NullBoolToBool(confirmationInfo.Approved){
		return RegResponse{
			Message: "Dean not confirmed",
			Data: confirmationInfo,
			StatusCode: status.NotAcceptable.Code,
			StatusCodeMessage: status.NotAcceptable.Message,
		},status.OK.Message,nil
	}
	return RegResponse{
		Message: "Dean is confirmed",
		Data: confirmationInfo,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *RegServiceStruct) CreateFaculty(ctx context.Context, facInfo CreateFacultyDto)(RegResponse,string,error){
	faculty := sqlc.CreateFacultyParams{
		FacultyName: facInfo.FacultyName,
		FacultyCode: utils.StringToNullString(facInfo.FacultyCode),
		UniversityID: utils.StringToUUID(facInfo.UniversityId),
		
	}

	uni,err := rs.regRepo.CreateFaculty(ctx,faculty)
	
	if err != nil {
		rs.logger.Error("Error creating faculty","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}

	return RegResponse{
		Message: "Faculty created successfully",
		Data:uni,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (rs *RegServiceStruct) RetrievePendingHods(ctx context.Context, hodParams PendingHodDto )(RegResponse,string,error){
	pendingHods := &sqlc.RetrievePendingHodsParams{
		UniversityID: utils.StringToUUID(hodParams.UniversityId),
		FacultyID: utils.StringToUUID(hodParams.FacultyId),
	}
	list,err := rs.regRepo.RetrievePendingHods(ctx,*pendingHods)

	if err != nil{
		rs.logger.Error("Error retrieving pending hods","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if len(list) == 0 {
		return RegResponse{},status.NotFound.Message,errors.New("no pending hods found")
	}
	return RegResponse{
		Message: "Here are the pending hods",
		Data: list,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *RegServiceStruct) ApproveHod(ctx context.Context,waitId string)(RegResponse,string,error){
	idExists,_,err := rs.regRepo.ApproveHod(ctx,utils.StringToUUID(waitId))

	if err != nil{
		if !idExists{
		return RegResponse{},status.NotFound.Message,err
	}
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Hod approved",
		Data: nil,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *RegServiceStruct) RequestHodConfirmation(ctx context.Context, hod RequestHodConfirmationDto)(RegResponse,string,error){
	var lecturerId string;

	// retrieves user id from token claims
	claims := ctx.Value(userInfoKey)
	if claims != nil{
		lecturerId = claims.(jwt.CustomClaims).User_id
	}else{
		return RegResponse{},status.InternalServerError.Message,errors.New("problem getiing claims from token")
	}
	
	hodInfo := sqlc.RequestHodConfirmationParams{
		LecturerID: utils.StringToUUID(lecturerId),
		PotentialDepartment: hod.PotentialDepartment,
		AdditionalMessage: utils.StringToNullString(hod.AdditionalMessage),
		UniversityID: utils.StringToUUID(hod.UniversityId),
		FacultyID: utils.StringToUUID(hod.FacultyId),
	}
	waitInfo,err := rs.regRepo.RequestHodConfirmation(ctx,hodInfo)

	desiredWaitInfo := sqlc.HodWaitingList{
		WaitID: waitInfo.WaitID,
	}

	if err != nil{
		rs.logger.Error("Error requesting hod confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Hod request confirmation sent",
		Data: desiredWaitInfo,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (rs *RegServiceStruct) CheckHodConfirmation(ctx context.Context, waitId string)(RegResponse,string,error){
	confirmationInfo,err := rs.regRepo.CheckHodConfirmation(ctx,utils.StringToUUID(waitId))
	if err != nil{
		rs.logger.Error("Error checking hod confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if !utils.NullBoolToBool(confirmationInfo.Approved){
		return RegResponse{
			Message: "Hod not confirmed",
			Data: confirmationInfo,
			StatusCode: status.NotAcceptable.Code,
			StatusCodeMessage: status.NotAcceptable.Message,
		},status.OK.Message,nil
	}
	return RegResponse{
		Message: "Hod is confirmed",
		Data: confirmationInfo,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}


func (rs *RegServiceStruct) CreateDepartment(ctx context.Context, deptInfo CreateDepartmentDto)(RegResponse,string,error){
	department := sqlc.CreateDepartmentParams{
		DepartmentName: deptInfo.DepartmentName,
		DepartmentCode: utils.StringToNullString(deptInfo.DepartmentCode),
		UniversityID: utils.StringToUUID(deptInfo.UniversityId),
		FacultyID: utils.StringToUUID(deptInfo.FacultyId),
		NumberOfLevels: int32(deptInfo.NumberOfLevels),
	}

	uni,err := rs.regRepo.CreateDepartment(ctx,department)
	
	if err != nil {
		rs.logger.Error("Error creating department","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}

	return RegResponse{
		Message: "Department created successfully",
		Data:uni,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (rs *RegServiceStruct) RetrievePendingLecturers(ctx context.Context, lecturerParams PendingLecturerDto )(RegResponse,string,error){
	pendingLecturers := &sqlc.RetrievePendingLecturersParams{
		UniversityID: utils.StringToUUID(lecturerParams.UniversityId),
		FacultyID: utils.StringToUUID(lecturerParams.FacultyId),
		DepartmentID: utils.StringToUUID(lecturerParams.DepartmentId),
	}
	list,err := rs.regRepo.RetrievePendingLecturers(ctx,*pendingLecturers)

	if err != nil{
		rs.logger.Error("Error retrieving pending lecturers","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if len(list) == 0 {
		return RegResponse{},status.NotFound.Message,errors.New("no pending lecturers found")
	}
	return RegResponse{
		Message: "Here are the pending lecturers",
		Data: list,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}


func (rs *RegServiceStruct) ApproveLecturer(ctx context.Context,waitId string)(RegResponse,string,error){
	idExists,_,err := rs.regRepo.ApproveLecturer(ctx,utils.StringToUUID(waitId))

	if err != nil{
		if !idExists{
		return RegResponse{},status.NotFound.Message,err
	}
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Lecturer approved",
		Data: nil,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *RegServiceStruct) RequestLecturerConfirmation(ctx context.Context, lecturer RequestLecturerConfirmationDto)(RegResponse,string,error){
	var lecturerId string;

	// retrieves user id from token claims
	claims := ctx.Value(userInfoKey)
	if claims != nil{
		lecturerId = claims.(jwt.CustomClaims).User_id
	}else{
		return RegResponse{},status.InternalServerError.Message,errors.New("problem getiing claims from token")
	}
	
	lecturerInfo := sqlc.RequestLecturerConfirmationParams{
		LecturerID: utils.StringToUUID(lecturerId),
		AdditionalMessage: utils.StringToNullString(lecturer.AdditionalMessage),
		UniversityID: utils.StringToUUID(lecturer.UniversityId),
		FacultyID: utils.StringToUUID(lecturer.FacultyId),
		DepartmentID: utils.StringToUUID(lecturer.DepartmentId),
	}
	waitInfo,err := rs.regRepo.RequestLecturerConfirmation(ctx,lecturerInfo)

	desiredWaitInfo := sqlc.LecturerWaitingList{
		WaitID: waitInfo.WaitID,
	}

	if err != nil{
		rs.logger.Error("Error requesting lecturer confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Lecturer request confirmation sent",
		Data: desiredWaitInfo,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}


func (rs *RegServiceStruct) CheckLecturerConfirmation(ctx context.Context, waitId string)(RegResponse,string,error){
	confirmationInfo,err := rs.regRepo.CheckLecturerConfirmation(ctx,utils.StringToUUID(waitId))
	if err != nil{
		rs.logger.Error("Error checking lecturer confirmation","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	if !utils.NullBoolToBool(confirmationInfo.Approved){
		return RegResponse{
			Message: "Lecturer not confirmed",
			Data: confirmationInfo,
			StatusCode: status.NotAcceptable.Code,
			StatusCodeMessage: status.NotAcceptable.Message,
		},status.OK.Message,nil
	}
	return RegResponse{
		Message: "Lecturer is confirmed",
		Data: confirmationInfo,
		StatusCode: status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	},status.OK.Message,nil
}

func (rs *RegServiceStruct) CreateDean(ctx context.Context, deanInfo CreateDeanDto )(RegResponse,string,error){
	strStartDate,err := utils.StringToTime(deanInfo.StartDate)
	 if err != nil{
		rs.logger.Error("Error converting start date to string","err:",err)
		return RegResponse{},status.BadRequest.Message,err
	 }
	strEndDate,err := utils.StringToNullTime(deanInfo.EndDate)
	 if err != nil{
		rs.logger.Error("Error converting end date to string","err:",err)
		return RegResponse{},status.BadRequest.Message,err
	 }
	deanParams := sqlc.CreateDeanParams{
		LecturerID: utils.StringToNullUUID(deanInfo.LecturerId),
		FacultyID: utils.StringToNullUUID(deanInfo.FacultyId),
		UniversityID: utils.StringToNullUUID(deanInfo.UniversityId),
		StartDate: strStartDate,
		EndDate:strEndDate,
	}
	dean,err := rs.regRepo.CreateDean(ctx,deanParams)

	if err != nil{
		rs.logger.Error("Error creating dean","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Dean Created Successfully",
		Data: dean,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}

func (rs *RegServiceStruct) CreateHod(ctx context.Context, hodInfo CreateHodDto)(RegResponse,string,error){
strStartDate,err := utils.StringToTime(hodInfo.StartDate)
	 if err != nil{
		rs.logger.Error("Error converting start date to string","err:",err)
		return RegResponse{},status.BadRequest.Message,err
	 }
	strEndDate,err := utils.StringToNullTime(hodInfo.EndDate)
	 if err != nil{
		rs.logger.Error("Error converting end date to string","err:",err)
		return RegResponse{},status.BadRequest.Message,err
	 }
	hodParams := sqlc.CreateHodParams{
		LecturerID: utils.StringToNullUUID(hodInfo.LecturerId),
		DepartmentID: utils.StringToNullUUID(hodInfo.DepartmentId),
		UniversityID: utils.StringToNullUUID(hodInfo.UniversityId),
		StartDate: strStartDate,
		EndDate:strEndDate,
	}
	hod,err := rs.regRepo.CreateHod(ctx,hodParams)

	if err != nil{
		rs.logger.Error("Error creating hod","err:",err)
		return RegResponse{},status.InternalServerError.Message,err
	}
	return RegResponse{
		Message: "Hod Created Successfully",
		Data: hod,
		StatusCode: status.Created.Code,
		StatusCodeMessage: status.Created.Message,
	},status.Created.Message,nil
}


// this is used by the middleware
func (rs *RegServiceStruct) CheckCurrentDean(ctx context.Context,deanId string)(bool,error){
	deanExists,_,err := rs.regRepo.RetrieveDean(ctx,utils.StringToUUID(deanId))

	if err != nil {
		rs.logger.Error("Error retrieving dean","err:",err)
		return false,err
	}
	if !deanExists{
		return false,errors.New("dean does not exist")
	}
	return true,nil
}

// this is used by the middleware
func (rs *RegServiceStruct) CheckCurrentHod(ctx context.Context,hodId string)(bool,error){
	hodExists,_,err := rs.regRepo.RetrieveHod(ctx,utils.StringToUUID(hodId))

	if err != nil {
		rs.logger.Error("Error retrieving hod","err:",err)
		return true,err
	}
	if !hodExists{
		return false,nil
	}
	return true,nil
}

// this is used by the middleware
func (rs *RegServiceStruct) CheckCurrentAdmin(ctx context.Context,adminId string)(bool,error){
	adminExists,_,err := rs.regRepo.RetrieveAdmin(ctx,utils.StringToUUID(adminId))

	if err != nil {
		rs.logger.Error("Error retrieving admin","err:",err)
		return false,err
	}
	if !adminExists{
		return false,errors.New("admin does not exist")
	}
	return true,nil
}