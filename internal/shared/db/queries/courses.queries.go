package queries

import (
	"context"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)



type CoursesQueries struct {
	q *sqlc.Queries
}



func NewCoursesQueries(q *sqlc.Queries) *CoursesQueries{
	return &CoursesQueries{
		q:q,
	}
}


func (cq *CoursesQueries) SetStudentCourse(ctx context.Context,studCourseInfo sqlc.SetStudentCourseParams)(sqlc.StudentCoursesOffered,error){
	return cq.q.SetStudentCourse(ctx,studCourseInfo)
}


func (cq *CoursesQueries) CreateCourse(ctx context.Context,courseInfo sqlc.CreateCourseParams)(sqlc.Course,error){
	return cq.q.CreateCourse(ctx,courseInfo)
}

func (cq *CoursesQueries) RetrieveCoursesForDepartment(ctx context.Context, courseParam sqlc.RetrieveCoursesForADepartmentParams)([]sqlc.RetrieveCoursesForADepartmentRow,error){
	return cq.q.RetrieveCoursesForADepartment(ctx,courseParam)
}

func (cq *CoursesQueries) UpdateCourse(ctx context.Context,courseInfo sqlc.UpdateCourseParams)(sqlc.Course,error){
	return cq.q.UpdateCourse(ctx,courseInfo)
}

func (cq *CoursesQueries) DeleteCourse(ctx context.Context,courseId uuid.UUID)error{
	return cq.q.DeleteCourse(ctx,courseId)
}


func (cq *CoursesQueries) SetCourseLecturers(ctx context.Context,param sqlc.SetCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.q.SetCourseLecturers(ctx,param)
}

func (cq *CoursesQueries) UpdateCourseLecturers(ctx context.Context,param sqlc.UpdateCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.q.UpdateCourseLecturers(ctx,param)
}

func (cq *CoursesQueries) CountCoursesForAUni(ctx context.Context,uniId uuid.UUID)(int64,error){
	return cq.q.CountTotalCourses(ctx,uniId)
}

func (cq *CoursesQueries) RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesRow,error){
	return cq.q.RetrieveAllCourses(ctx,uniId)
}

func (cq *CoursesQueries) RetrieveAllCoursesAndVenues(ctx context.Context, uniId uuid.UUID)([]sqlc.RetrieveAllCoursesAndTheirVenueIdsRow,error){
	return cq.q.RetrieveAllCoursesAndTheirVenueIds(ctx,uniId)
}

func (cq *CoursesQueries) CreateCohortCourse(ctx context.Context,params sqlc.CreateCohortCourseParams)(sqlc.CohortCoursesOffered,error){
	return cq.q.CreateCohortCourse(ctx,params)
}


func (cq *CoursesQueries) FetchCoursesForACohort(ctx context.Context, params sqlc.FetchCoursesForACohortParams)([]uuid.UUID,error){
	return cq.q.FetchCoursesForACohort(ctx,params)
}

func (cq *CoursesQueries) FetchSessionsForACohort(ctx context.Context,params sqlc.GetCohortSessionsInCurrentTimetableParams)([]sqlc.GetCohortSessionsInCurrentTimetableRow,error){
	return cq.q.GetCohortSessionsInCurrentTimetable(ctx,params)
}

func (cq *CoursesQueries) FetchSessionsForAStudent(ctx context.Context,studentId uuid.UUID)([]sqlc.GetStudentTimetableSessionsRow,error){
	return cq.q.GetStudentTimetableSessions(ctx,studentId)
}