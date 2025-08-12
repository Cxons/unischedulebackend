package repository

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)



type RegRepository interface{
	UpdateAdmin(ctx context.Context,adminInfo sqlc.UpdateAdminInfoParams)(bool,sqlc.UniversityAdmin,error)
	CreateUniversity(ctx context.Context, uniInfo sqlc.CreateUniversityParams)(sqlc.University,error)
	RetrievePendingDeans(ctx context.Context, uniId uuid.UUID)([]sqlc.DeanWaitingList,error)
	ApproveDean(ctx context.Context,waitId uuid.UUID)(sqlc.DeanWaitingList,error)
	RetrievePendingHods(ctx context.Context, deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error)
	ApproveHod(ctx context.Context, waitId uuid.UUID)(sqlc.HodWaitingList,error)
	RetrievePendingLecturers(ctx context.Context, hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error)
	ApproveLecturer(ctx context.Context, waitId uuid.UUID)(sqlc.LecturerWaitingList,error)
	RequestDeanConfirmation(ctx context.Context, dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error)
	RequestHodConfirmation(ctx context.Context, hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error)
	RequestLecturerConfirmation(ctx context.Context, lecturer sqlc.RequestLecturerConfirmationParams)(sqlc.LecturerWaitingList,error)
}


type regRepository struct{
	aq *queries.AdminQueries
	sq *queries.StudentQueries
	lq *queries.LecturerQueries
	uq *queries.UniQueries
	dq *queries.DeanQueries
	hq *queries.HodQueries

}

func NewRegRepository(aq *queries.AdminQueries,sq *queries.StudentQueries, lq *queries.LecturerQueries, uq *queries.UniQueries, dq *queries.DeanQueries, hq *queries.HodQueries)*regRepository{
	return &regRepository{
		aq: aq,
		sq: sq,
		lq: lq,
		uq:uq,
		dq:dq,
		hq:hq,
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


func (rrp *regRepository) CreateUniversity(ctx context.Context, uniInfo sqlc.CreateUniversityParams)(sqlc.University,error){
	return rrp.uq.CreateUniversities(ctx,uniInfo)
}

func (rrp *regRepository) RetrievePendingDeans(ctx context.Context,uniId uuid.UUID)([]sqlc.DeanWaitingList,error){
	return rrp.aq.RetrievePendingDeans(ctx,uniId)
}

func (rrp *regRepository) ApproveDean(ctx context.Context,waitId uuid.UUID)(sqlc.DeanWaitingList,error){
	return rrp.aq.ApproveDean(ctx,waitId)
}

func (rrp *regRepository)RequestDeanConfirmation(ctx context.Context, dean sqlc.RequestDeanConfirmationParams)(sqlc.DeanWaitingList,error){
	return rrp.dq.RequestDeanConfirmation(ctx,dean)
}

func (rrp *regRepository) RetrievePendingHods(ctx context.Context,deanInfo sqlc.RetrievePendingHodsParams)([]sqlc.HodWaitingList,error){
	return rrp.dq.RetrievePendingHods(ctx,deanInfo)
}

func (rrp *regRepository) ApproveHod(ctx context.Context, waitId uuid.UUID)(sqlc.HodWaitingList,error){
	return rrp.dq.ApproveHod(ctx,waitId)
}

func (rrp *regRepository) RequestHodConfirmation(ctx context.Context, hod sqlc.RequestHodConfirmationParams)(sqlc.HodWaitingList,error){
	return rrp.hq.RequestHodConfirmation(ctx,hod)
}

func (rrp *regRepository) RetrievePendingLecturers(ctx context.Context, hodInfo sqlc.RetrievePendingLecturersParams)([]sqlc.LecturerWaitingList,error){
	return rrp.hq.RetrievePendingLecturers(ctx,hodInfo)
}

func (rrp *regRepository) ApproveLecturer(ctx context.Context, waitId uuid.UUID)(sqlc.LecturerWaitingList,error){
	return rrp.hq.ApproveLecturer(ctx,waitId)
}

func (rrp *regRepository) RequestLecturerConfirmation(ctx context.Context, lecturer sqlc.RequestLecturerConfirmationParams)(sqlc.LecturerWaitingList,error){
	return rrp.lq.RequestLecturerConfirmation(ctx,lecturer)
}