package repository

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)



type UniRepository interface{
	RetrieveAllUniversities(ctx context.Context)([]sqlc.University,error)
	RetrieveAllFaculties(ctx context.Context,uniId uuid.UUID)([]sqlc.Faculty,error)
	RetrieveAllDepartments(ctx context.Context,deptParams sqlc.RetrieveDeptsForAFacultyParams)([]sqlc.Department,error)
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