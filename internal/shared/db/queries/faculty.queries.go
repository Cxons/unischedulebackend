package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)


type FacQueries struct {
	q *sqlc.Queries
}


func NewFacQueries(q *sqlc.Queries)*FacQueries{
	return &FacQueries{
		q:q,
	}
}


func (fq *FacQueries) CreateFaculty(ctx context.Context, facInfo sqlc.CreateFacultyParams)(sqlc.Faculty,error){
	return fq.q.CreateFaculty(ctx,facInfo)
}

func (fq *FacQueries) RetrieveAllFaculties(ctx context.Context, uniId uuid.UUID)([]sqlc.Faculty,error){
	return fq.q.RetrieveFacultiesForAUni(ctx, uniId)
}

func (fq *FacQueries) UpdateFaculty(ctx context.Context,facParams sqlc.UpdateFacultyParams)(sqlc.Faculty,error){
	return fq.q.UpdateFaculty(ctx,facParams)
}