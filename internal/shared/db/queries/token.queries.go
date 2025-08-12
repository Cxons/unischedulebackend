package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)


type TokenQueries struct {
	q *sqlc.Queries
}

func NewTokenQueries(q *sqlc.Queries) *TokenQueries{
	return &TokenQueries{
		q:q,
	}
}

func (tq *TokenQueries) AddRefreshToken(ctx context.Context, tokenInfo sqlc.AddRefreshTokenParams)(sqlc.RefreshToken,error){
	return tq.q.AddRefreshToken(ctx,tokenInfo);
}
func (tq *TokenQueries) UpdateRefreshToken(ctx context.Context,tokenInfo sqlc.UpdateRefreshTokenParams)(sqlc.RefreshToken,error){
	return tq.q.UpdateRefreshToken(ctx,tokenInfo)
}
func (tq *TokenQueries) DeleteRefreshToken(ctx context.Context,userId uuid.UUID)(sqlc.RefreshToken,error){
	return tq.q.DeleteRefreshToken(ctx,userId)
}
func (tq *TokenQueries) CheckAndReturnToken(ctx context.Context,userId uuid.UUID)(sqlc.CheckAndReturnTokenRow,error){
	return tq.q.CheckAndReturnToken(ctx,userId)
}