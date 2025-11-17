package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/google/uuid"
)



type RegRepository interface{
	UpdateAdmin(ctx context.Context,adminInfo sqlc.UpdateAdminInfoParams)(bool,sqlc.UniversityAdmin,error)
	CreateUniversity(ctx context.Context, uniInfo sqlc.CreateUniversityParams)(sqlc.University,error)
	RetrievePendingDeans(ctx context.Context, uniId uuid.UUID)([]sqlc.DeanWaitingList,error)
	ApproveDean(ctx context.Context,waitId uuid.UUID,lecturerId uuid.UUID)(bool,sqlc.DeanWaitingList,error)
	RetrievePendingHods(ctx context.Context, deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error)
	ApproveHod(ctx context.Context, waitId uuid.UUID)(bool,sqlc.HodWaitingList,error)
	RetrievePendingLecturers(ctx context.Context, hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error)
	ApproveLecturer(ctx context.Context, waitId uuid.UUID)(bool,sqlc.LecturerWaitingList,error)
	RequestDeanConfirmation(ctx context.Context, dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error)
	RequestHodConfirmation(ctx context.Context, hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error)
	RequestLecturerConfirmation(ctx context.Context, lecturer sqlc.RequestLecturerConfirmationParams)(sqlc.LecturerWaitingList,error)
	CheckDeanConfirmation(ctx context.Context,waitId uuid.UUID)(sqlc.CheckDeanConfirmationRow,error)
	CheckHodConfirmation(ctx context.Context, waitId uuid.UUID)(sqlc.CheckHodConfirmationRow,error)
	CheckLecturerConfirmation(ctx context.Context, waitId uuid.UUID)(sqlc.CheckLecturerConfirmationRow,error)
	CreateFaculty(ctx context.Context, facInfo sqlc.CreateFacultyParams, lecturerId uuid.UUID,startTime time.Time,endTime time.Time)(sqlc.Faculty,error,uuid.UUID)
	CreateDepartment(ctx context.Context,deptInfo sqlc.CreateDepartmentParams,lecturerId uuid.UUID, startTime time.Time, endTime time.Time)(sqlc.Department,uuid.UUID,error)
	CreateDean(ctx context.Context, deanInfo sqlc.CreateDeanParams)(sqlc.CurrentDean,error)
	CreateHod(ctx context.Context, hodInfo sqlc.CreateHodParams)(sqlc.CurrentHod,error)
	RetrieveAdmin(ctx context.Context,adminId uuid.UUID)(bool,sqlc.RetrieveAdminRow,error)
	RetrieveDean(ctx context.Context,deanId uuid.UUID)(bool,sqlc.RetrieveDeanRow,error)
	RetrieveHod(ctx context.Context, hodId uuid.UUID)(bool,sqlc.RetrieveHodRow,error)
	CheckDeanConfirmationWithLecturerId(ctx context.Context,lecturerId uuid.UUID)(sqlc.CheckDeanConfirmationWithLecturerIdRow,error)
	CheckHodConfirmationWithLecturerId(ctx context.Context,lecturerId uuid.UUID)(sqlc.CheckHodConfirmationWithLecturerIdRow,error)
	CreateLecturerUnavailability(ctx context.Context, data []sqlc.CreateLecturerUnavailabilityParams)error

}


type regRepository struct{
	aq *queries.AdminQueries
	sq *queries.StudentQueries
	lq *queries.LecturerQueries
	uq *queries.UniQueries
	dq *queries.DeanQueries
	hq *queries.HodQueries
	fq *queries.FacQueries
	dptq *queries.DeptQueries
	store sqlc.Store

}

func NewRegRepository(aq *queries.AdminQueries,sq *queries.StudentQueries, lq *queries.LecturerQueries, uq *queries.UniQueries, dq *queries.DeanQueries, hq *queries.HodQueries, fq *queries.FacQueries,dptq *queries.DeptQueries, store sqlc.Store)*regRepository{
	return &regRepository{
		aq: aq,
		sq: sq,
		lq: lq,
		uq:uq,
		dq:dq,
		hq:hq,
		fq:fq,
		dptq:dptq,
		store: store,
	}
}

func (rrp *regRepository) UpdateAdmin(ctx context.Context,adminInfo sqlc.UpdateAdminInfoParams)(bool,sqlc.UniversityAdmin,error){
	admin,err := rrp.aq.UpdateAdmin(ctx,adminInfo)
	if err != nil {
		if err == sql.ErrNoRows{
		return false,sqlc.UniversityAdmin{},errors.New("admin not found")
	}
	return true,sqlc.UniversityAdmin{},err
	}
	return true,admin,nil
}


func (rrp *regRepository) CreateUniversity(ctx context.Context, uniInfo sqlc.CreateUniversityParams) (sqlc.University, error) {
    uni, err := rrp.uq.CreateUniversities(ctx, uniInfo)
    if err != nil {
        // Check for unique constraint violation on phone_number
        if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "phone_number") {
            return sqlc.University{}, errors.New("Only one university creation allowed")
        }
        return sqlc.University{}, err
    }
    return uni, nil
}

func (rrp *regRepository) RetrievePendingDeans(ctx context.Context,uniId uuid.UUID)([]sqlc.DeanWaitingList,error){
	return rrp.aq.RetrievePendingDeans(ctx,uniId)
}

func (rrp *regRepository) ApproveDean(ctx context.Context,waitId uuid.UUID,lecturerId uuid.UUID)(bool,sqlc.DeanWaitingList,error){
	deanRow,err := rrp.aq.ApproveDean(ctx,waitId)
	if err != nil {
		if err == sql.ErrNoRows{
		return false,sqlc.DeanWaitingList{},errors.New("wait id not found")
	}
	return true,sqlc.DeanWaitingList{},err
	}
	return true,deanRow,nil
}

func (rrp *regRepository)RequestDeanConfirmation(ctx context.Context, dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error){
	data,err := rrp.dq.RequestDeanConfirmation(ctx,dean)
	 if err != nil {
        // Check for unique constraint violation on phone_number
        if strings.Contains(err.Error(), "unique constraint"){
            return sqlc.DeanWaitingList{}, errors.New("Only one dean request creation allowed")
        }
        return sqlc.DeanWaitingList{}, err
    }
    return data, nil
}

func (rrp *regRepository) CheckDeanConfirmation(ctx context.Context,waitId uuid.UUID)(sqlc.CheckDeanConfirmationRow,error){
	return rrp.dq.CheckDeanConfirmation(ctx,waitId)
}

func (rrp *regRepository) CheckDeanConfirmationWithLecturerId(ctx context.Context,lecturerId uuid.UUID)(sqlc.CheckDeanConfirmationWithLecturerIdRow,error){
	dean,err := rrp.dq.CheckDeanConfirmationWithLecturerId(ctx,lecturerId)
	if err != nil {
		if err == sql.ErrNoRows{
		return sqlc.CheckDeanConfirmationWithLecturerIdRow{},errors.New("dean not confirmed")
	}
	return sqlc.CheckDeanConfirmationWithLecturerIdRow{},err
}
	return dean,nil
}

func (rrp *regRepository) CheckHodConfirmationWithLecturerId(ctx context.Context,lecturerId uuid.UUID)(sqlc.CheckHodConfirmationWithLecturerIdRow,error){
	hod,err := rrp.hq.CheckHodConfirmationWithLecturerId(ctx,lecturerId)
	if err != nil {
		if err == sql.ErrNoRows{
		return sqlc.CheckHodConfirmationWithLecturerIdRow{},errors.New("hod not confirmed")
	}
	return sqlc.CheckHodConfirmationWithLecturerIdRow{},err
}
	return hod,nil
}


func (rrp *regRepository) CreateFaculty(
	ctx context.Context,
	facInfo sqlc.CreateFacultyParams,
	lecturerId uuid.UUID,
	startTime time.Time,
	endTime time.Time,
) (sqlc.Faculty, error,uuid.UUID) {

	var fac sqlc.Faculty
	var deanId uuid.UUID

	err := rrp.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		faculty, err := q.CreateFaculty(ctx, facInfo)
		if err != nil {
			return err
		}

		dean, deanErr := q.InsertCurrentDean(ctx, sqlc.InsertCurrentDeanParams{
			LecturerID:  utils.UuidToNullUUID(lecturerId),
			FacultyID:   utils.UuidToNullUUID(faculty.FacultyID), // ✅ FIXED HERE
			UniversityID: utils.UuidToNullUUID(faculty.UniversityID),
			StartDate:   startTime,
			EndDate:     utils.TimeToNulltime(endTime),
		})
		if deanErr != nil {
			return deanErr
		}
		deanId = dean.DeanID

		fac = faculty
		return nil
	})

	return fac, err,deanId
}

func (rrp *regRepository) RetrievePendingHods(ctx context.Context,deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error){
	return rrp.dq.RetrievePendingHods(ctx,deanInfo)
}

func (rrp *regRepository) ApproveHod(ctx context.Context, waitId uuid.UUID)(bool,sqlc.HodWaitingList,error){
	hodRow,err := rrp.dq.ApproveHod(ctx,waitId)
	if err != nil {
		if err == sql.ErrNoRows{
		return false,sqlc.HodWaitingList{},errors.New("wait id not found")
	}
	return true,sqlc.HodWaitingList{},err
	}
	return true,hodRow,nil

}

func (rrp *regRepository) RequestHodConfirmation(ctx context.Context, hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error){
	return rrp.hq.RequestHodConfirmation(ctx,hod)
}

func (rrp *regRepository) CheckHodConfirmation(ctx context.Context, waitId uuid.UUID)(sqlc.CheckHodConfirmationRow,error){
	return rrp.hq.CheckHodConfirmation(ctx,waitId)
} 

func (rrp *regRepository) CreateDepartment(ctx context.Context,deptInfo sqlc.CreateDepartmentParams,lecturerId uuid.UUID,startTime time.Time,endTime time.Time)(sqlc.Department,uuid.UUID,error){
		var dept sqlc.Department
		var hodId uuid.UUID
	err := rrp.store.ExecTx(ctx,func(q *sqlc.Queries)error{
	
		department,createDeptErr := q.CreateDepartment(ctx,deptInfo)
		if createDeptErr != nil{
			return createDeptErr
		}
		for i := range department.NumberOfLevels{
			cohort := sqlc.CreateCohortParams{
				CohortName: department.DepartmentName + " " + strconv.Itoa(int(i + 1)) + "00 level",
				CohortLevel: i+1,
				CohortDepartmentID: department.DepartmentID,
				CohortFacultyID: department.FacultyID,
				CohortUniversityID: department.UniversityID,
			}
			_,createCohortErr := q.CreateCohort(ctx,cohort)
			if createCohortErr != nil{
				return createCohortErr
			}
		}
		hod, hodErr := q.InsertCurrentHod(ctx, sqlc.InsertCurrentHodParams{
			LecturerID:  utils.UuidToNullUUID(lecturerId),
			DepartmentID:   utils.UuidToNullUUID(department.DepartmentID),
			UniversityID: utils.UuidToNullUUID(department.UniversityID),
			StartDate:   startTime,
			EndDate:     utils.TimeToNulltime(endTime),
		})
		if hodErr != nil {
			return hodErr
		}
		hodId = hod.HodID
		dept = department

		return nil
	})
	return dept,hodId,err
}

func (rrp *regRepository) RetrievePendingLecturers(ctx context.Context, hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error){
	return rrp.hq.RetrievePendingLecturers(ctx,hodInfo)
}

func (rrp *regRepository) ApproveLecturer(ctx context.Context, waitId uuid.UUID)(bool,sqlc.LecturerWaitingList,error){
	lecturerRow,err := rrp.hq.ApproveLecturer(ctx,waitId)
	if err != nil {
		if err == sql.ErrNoRows{
		return false,sqlc.LecturerWaitingList{},errors.New("wait id not found")
	}
	return true,sqlc.LecturerWaitingList{},err
	}
	return true,lecturerRow,nil
}

func (rrp *regRepository) RequestLecturerConfirmation(ctx context.Context, lecturer sqlc.RequestLecturerConfirmationParams)(sqlc.LecturerWaitingList,error){
	return rrp.lq.RequestLecturerConfirmation(ctx,lecturer)
}

func (rrp *regRepository) CheckLecturerConfirmation(ctx context.Context, waitId uuid.UUID)(sqlc.CheckLecturerConfirmationRow,error){
	return rrp.lq.CheckLecturerConfirmation(ctx,waitId)
}

func (rrp *regRepository) CreateDean(ctx context.Context, deanInfo sqlc.CreateDeanParams)(sqlc.CurrentDean,error){
	return rrp.dq.CreateDean(ctx,deanInfo)
}

func (rrp *regRepository) CreateHod(ctx context.Context, hodInfo sqlc.CreateHodParams)(sqlc.CurrentHod,error){
	return rrp.hq.CreateHod(ctx,hodInfo)
}

func (rrp *regRepository) RetrieveAdmin(ctx context.Context,adminId uuid.UUID)(bool,sqlc.RetrieveAdminRow,error){
	admin,err := rrp.aq.RetrieveAdmin(ctx,adminId)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveAdminRow{},nil
		}
		return true,sqlc.RetrieveAdminRow{},err
	}
	return true,admin,nil
}

func (rrp *regRepository) RetrieveDean(ctx context.Context,deanId uuid.UUID)(bool,sqlc.RetrieveDeanRow,error){
	dean,err := rrp.dq.RetrieveDean(ctx,deanId)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveDeanRow{},nil
		}
		return true,sqlc.RetrieveDeanRow{},err
	}
	return true,dean,nil
}

func (rrp *regRepository) RetrieveHod(ctx context.Context, hodId uuid.UUID)(bool,sqlc.RetrieveHodRow,error){
	hod,err := rrp.hq.RetrieveHod(ctx,hodId)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,sqlc.RetrieveHodRow{},nil
		}
		return true,sqlc.RetrieveHodRow{},err
	}
	return true,hod,nil
	
}


func (rrp *regRepository) CreateLecturerUnavailability(ctx context.Context, data []sqlc.CreateLecturerUnavailabilityParams)error{
	return rrp.store.ExecTx(ctx,func(q *sqlc.Queries) error {
		for _,val := range data{
			err := q.CreateLecturerUnavailability(ctx,val)
			if err != nil{
				return err
			}
		}
		return  nil
	})
}