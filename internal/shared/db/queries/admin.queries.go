package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)

type AdminQueries struct {
 q *sqlc.Queries
}

func NewAdminQueries(q *sqlc.Queries) *AdminQueries{
	return &AdminQueries{
		q:q,
	}
}

func (aq *AdminQueries) RegisterAdmin(ctx context.Context,admin sqlc.RegisterUniversityAdminParams)(sqlc.UniversityAdmin,error){
	return aq.q.RegisterUniversityAdmin(ctx,admin)
}

func (aq *AdminQueries) RetrieveAdmin(ctx context.Context,adminId uuid.UUID)(sqlc.RetrieveAdminRow,error){
	return aq.q.RetrieveAdmin(ctx,adminId)
}
func (aq *AdminQueries) RetrieveAdminEmail(ctx context.Context,email string)(sqlc.RetrieveAdminEmailRow,error){
	return aq.q.RetrieveAdminEmail(ctx,email)
}

func (aq *AdminQueries) UpdateAdmin(ctx context.Context,adminInfo sqlc.UpdateAdminInfoParams)(sqlc.UniversityAdmin,error){
	return aq.q.UpdateAdminInfo(ctx,adminInfo)
}

func (aq *AdminQueries) RetrievePendingDeans(ctx context.Context,uniId uuid.UUID)([]sqlc.DeanWaitingList,error){
	return aq.q.RetrievePendingDeans(ctx,uniId)
}

func (aq *AdminQueries) ApproveDean(ctx context.Context, waitId uuid.UUID)(sqlc.DeanWaitingList,error){
	return aq.q.ApproveDean(ctx,waitId)
}



