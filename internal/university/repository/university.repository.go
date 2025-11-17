package repository

import (
	"context"
	"log/slog"

	"github.com/Cxons/unischedulebackend/internal/shared/constants"
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/university/dto"
	"github.com/google/uuid"
)



type UniRepository interface{
	RetrieveAllUniversities(ctx context.Context)([]sqlc.University,error)
	RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)([]sqlc.Faculty,error)
	RetrieveAllDepartments(ctx context.Context,deptParams sqlc.RetrieveDeptsForAFacultyParams)([]sqlc.Department,error)
	FetchApprovedLecturersInDepartment(ctx context.Context, deptId uuid.UUID)([]sqlc.FetchApprovedLecturersInDepartmentRow,error)
	CreateVenue(ctx context.Context,venueInfo sqlc.CreateVenueParams,venueType string, id uuid.UUID,unavailabilityData dto.VenueUnavailability)error
	RetrieveAllVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllVenuesRow,error)
	FetchCohortsForADepartment(ctx context.Context,deptId uuid.UUID)([]sqlc.FetchCohortsForADepartmentRow,error)
	FetchAllDepartmentsForAUni(ctx context.Context,uniId uuid.UUID)([]sqlc.FetchAllDepartmentsForAUniRow,error)
}



type uniRepository struct{
	aq *queries.AdminQueries
	sq *queries.StudentQueries
	lq *queries.LecturerQueries
	uq *queries.UniQueries
	dq *queries.DeanQueries
	hq *queries.HodQueries
	fq *queries.FacQueries
	dptq *queries.DeptQueries
	cohq *queries.CohortQueries
	vq *queries.VenueQueries
	store sqlc.Store
}


func NewUniRepository(
	aq *queries.AdminQueries,
	sq *queries.StudentQueries,
	lq *queries.LecturerQueries,
	uq *queries.UniQueries,
	dq *queries.DeanQueries,
	hq *queries.HodQueries,
	fq *queries.FacQueries,
	dptq *queries.DeptQueries,
	cohq *queries.CohortQueries,
	vq *queries.VenueQueries,
	store sqlc.Store,
)*uniRepository{
	return &uniRepository{
		aq: aq,
		sq: sq,
		lq: lq,
		uq: uq,
		dq: dq,
		hq: hq,
		fq: fq,
		dptq: dptq,
		cohq: cohq,
		vq: vq,
		store: store,
	}
}


func (unp *uniRepository) RetrieveAllUniversities(ctx context.Context)([]sqlc.University,error){
	return unp.uq.RetrieveAllUniversities(ctx)
}


func (unp *uniRepository) RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)([]sqlc.Faculty,error){
	return unp.fq.RetrieveAllFaculties(ctx,uniId)
}

func (unp *uniRepository) RetrieveAllDepartments(ctx context.Context,deptParams sqlc.RetrieveDeptsForAFacultyParams)([]sqlc.Department,error){
	return unp.dptq.RetrieveAllDepartments(ctx,deptParams)
}

func (unp *uniRepository) FetchAllDepartmentsForAUni(ctx context.Context,uniId uuid.UUID)([]sqlc.FetchAllDepartmentsForAUniRow,error){
	return unp.dptq.FetchAllDepartmentsForAUni(ctx,uniId)
}

func (unp *uniRepository) FetchApprovedLecturersInDepartment(ctx context.Context, deptId uuid.UUID)([]sqlc.FetchApprovedLecturersInDepartmentRow,error){
	return unp.lq.FetchApprovedLecturersInDepartment(ctx,deptId)
}

func (unp *uniRepository) RetrieveAllVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllVenuesRow,error){
	return unp.vq.RetrieveAllVenues(ctx,uniId)
}

func (unp *uniRepository) CreateVenue(ctx context.Context,venueInfo sqlc.CreateVenueParams,venueType string, id uuid.UUID,unavailabilityData dto.VenueUnavailability)error{
	return unp.store.ExecTx(ctx,func(q *sqlc.Queries) error {
		venue,err := q.CreateVenue(ctx,venueInfo)
		if err != nil{
			slog.Error("error creating actual venue","err:",err)
			return err
		}
		if venueType == constants.FACULTY{
			err := q.SetFacultyVenue(ctx,sqlc.SetFacultyVenueParams{
				VenueID: venue.VenueID,
				FacultyID: id,
				UniversityID: venue.UniversityID,
			})
			if err != nil{
				slog.Error("error setting faculty venue","err:",err)
				return err
			}
		}
		if venueType == constants.DEPARTMENT{
			err := q.SetDepartmentVenue(ctx,sqlc.SetDepartmentVenueParams{
				VenueID: venue.VenueID,
				DepartmentID: id,
				UniversityID: venue.UniversityID,
			})
			if err != nil{
				slog.Error("error setting department venue","err:",err)
				return err
			}
		}
		slog.Info("unavailability data","datavenueId",venue.VenueID)
		unavailabilityErr := q.CreateVenueUnavailablity(ctx,sqlc.CreateVenueUnavailablityParams{
			VenueID: venue.VenueID,
			Reason: utils.StringToNullString(unavailabilityData.Reason),
			StartTime: utils.TimeToNulltime(unavailabilityData.StartTime),
			EndTime: utils.TimeToNulltime(unavailabilityData.EndTime),
			Day: utils.StringToNullString(unavailabilityData.Day),
			UniversityID: venue.UniversityID,
		})
		if unavailabilityErr != nil{
			return unavailabilityErr
		}
		return nil
	})
}


func (unp *uniRepository) FetchCohortsForADepartment(ctx context.Context,deptId uuid.UUID)([]sqlc.FetchCohortsForADepartmentRow,error){
	return unp.cohq.FetchCohortsForADepartment(ctx,deptId)
}