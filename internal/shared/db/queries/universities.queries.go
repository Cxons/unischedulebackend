package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
)



type UniQueries struct {
	q *sqlc.Queries
}


func NewUniQueries(q *sqlc.Queries) *UniQueries{
	return &UniQueries{
		q:q,
	}
}

func (uq *UniQueries) CreateUniversities(ctx context.Context, uniInfo sqlc.CreateUniversityParams)(sqlc.University,error){
	return uq.q.CreateUniversity(ctx, uniInfo)
}
func (uq *UniQueries) RetrieveAllUniversities(ctx context.Context)([]sqlc.University,error){
	return uq.q.RetrieveAllUniversities(ctx)
}
func (uq *UniQueries) RetrieveAllUniversitiesWithLimit(ctx context.Context,limit int32)([]sqlc.University,error){
	return uq.q.RetrieveUniversitiesWithLimit(ctx,limit)
}
func (uq *UniQueries) UpdateUniversity(ctx context.Context,uniParams sqlc.UpdateUniversityParams)(sqlc.University,error){
	return uq.q.UpdateUniversity(ctx,uniParams)
}