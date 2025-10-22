package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)




type CohortQueries struct {
	q *sqlc.Queries
}


func NewCohortQueries(q *sqlc.Queries) *CohortQueries{
	return &CohortQueries{
		q:q,
	}
}


func (cohq *CohortQueries) CreateCohort(ctx context.Context,cohortParams sqlc.CreateCohortParams)(sqlc.Cohort,error){
	return cohq.q.CreateCohort(ctx,cohortParams)
}


func (cohq *CohortQueries) UpdateCohort(ctx context.Context, updateCohortParams sqlc.UpdateCohortParams)(sqlc.Cohort,error){
	return cohq.q.UpdateCohort(ctx,updateCohortParams)
}

func (cohq *CohortQueries) RetrieveAllCohorts(ctx context.Context,uniId uuid.UUID)([]sqlc.Cohort,error){
	return cohq.q.RetrieveAllCohorts(ctx,uniId)
}

func (cohq *CohortQueries) CountCohortsForAUni(ctx context.Context,uniId uuid.UUID)(int64,error){
	return cohq.q.CountCohortsForOneUni(ctx,uniId)
}