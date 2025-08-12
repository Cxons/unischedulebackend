package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
)



type OtpQueries struct {
	q *sqlc.Queries
}


func NewOtpQueries(q *sqlc.Queries) *OtpQueries{
	return &OtpQueries{
		q:q,
	}
}

func (oq *OtpQueries) InsertOtp(ctx context.Context, otpInfo sqlc.InsertOtpParams)(sqlc.Otp,error){
	return oq.q.InsertOtp(ctx,otpInfo)
}

func (oq *OtpQueries) RetrieveOtp(ctx context.Context,email string)(sqlc.RetrieveOtpRow,error){
	return oq.q.RetrieveOtp(ctx,email)
}

func (oq *OtpQueries) UpdateOtp(ctx context.Context , otpInfo sqlc.UpdateOtpParams)(sqlc.Otp,error){
	return oq.q.UpdateOtp(ctx,otpInfo)
}