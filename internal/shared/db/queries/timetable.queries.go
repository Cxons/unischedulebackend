package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)




type TimeTableQueries struct {
	q *sqlc.Queries
}



func NewTimeTableQueries(q *sqlc.Queries)*TimeTableQueries{
	return &TimeTableQueries{
		q:q,
	}
}


func (tmtq *TimeTableQueries) CreateSessionPlacement(ctx context.Context,params sqlc.CreateSessionPlacementsParams)(sqlc.SessionPlacement,error){
	return tmtq.q.CreateSessionPlacements(ctx,params)
}

func (tmtq *TimeTableQueries) CreateCandidate(ctx context.Context, params sqlc.CreateCandidateParams)(sqlc.Candidate,error){
	return tmtq.q.CreateCandidate(ctx,params)
}

func (tmtq *TimeTableQueries) DeprecateLatestCandidate(ctx context.Context,uniId uuid.UUID)error{
	return tmtq.q.DeprecateLatestCandidate(ctx,uniId)
}
func (tmtq *TimeTableQueries) RestoreCurrentCandidate(ctx context.Context,uniId uuid.UUID)error{
	return tmtq.q.RestoreCurrentCandidate(ctx,uniId)
}