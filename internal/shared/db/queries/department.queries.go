package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)




type DeptQueries struct {
	q *sqlc.Queries
}



func NewDeptQueries(q *sqlc.Queries) *DeptQueries{
	return &DeptQueries{
		q:q,
	}
}

func (dq *DeptQueries) CreateDeparment(ctx context.Context, deptInfo sqlc.CreateDepartmentParams)(sqlc.Department,error){
	return dq.q.CreateDepartment(ctx,deptInfo)
}

func (dq *DeptQueries) RetrieveAllDepartments(ctx context.Context,deptParams sqlc.RetrieveDeptsForAFacultyParams)([]sqlc.Department,error){
	return dq.q.RetrieveDeptsForAFaculty(ctx,deptParams)
}

func (dq *DeptQueries) UpdateDepartment(ctx context.Context, deptInfo sqlc.UpdateDepartmentParams)(sqlc.Department,error){
	return dq.q.UpdateDepartment(ctx, deptInfo)
}

func (dq *DeptQueries) FetchAllDepartmentsForAUni(ctx context.Context,uniId uuid.UUID)([]sqlc.FetchAllDepartmentsForAUniRow,error){
	return dq.q.FetchAllDepartmentsForAUni(ctx,uniId)
}