package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
)





type VenueQueries struct {
	q *sqlc.Queries
}


func NewVenueQueries(q *sqlc.Queries)  *VenueQueries{
	return &VenueQueries{
		q:q,
	}
}



func (vq *VenueQueries) CreateVenue(ctx context.Context, venueInfo sqlc.CreateVenueParams)(sqlc.Venue,error){
	return vq.q.CreateVenue(ctx,venueInfo)
}


func (vq *VenueQueries) SetFacultyVenue(ctx context.Context, param sqlc.SetFacultyVenueParams)error{
	return vq.q.SetFacultyVenue(ctx,param)
}


func (vq *VenueQueries) SetDepartmentVenue(ctx context.Context, param sqlc.SetDepartmentVenueParams)error{
	return vq.q.SetDepartmentVenue(ctx,param)
}