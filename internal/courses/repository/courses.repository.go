package repository

import (
	"context"
	"database/sql"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/db/queries"
	"github.com/google/uuid"
)



type CourseRepository interface{
	CreateCourse(ctx context.Context,courseInfo sqlc.CreateCourseParams)(sqlc.Course,error)
	RetrieveCoursesForADepartment(ctx context.Context,courseParam sqlc.RetrieveCoursesForADepartmentParams)(bool,[]sqlc.RetrieveCoursesForADepartmentRow,error)
	UpdateCourse(ctx context.Context,courseInfo sqlc.UpdateCourseParams)(sqlc.Course,error)
	SetStudentCourses(ctx context.Context,studCourseParam []sqlc.SetStudentCourseParams)error
	DeleteCourse(ctx context.Context,courseId uuid.UUID)error
	SetCourseLecturers(ctx context.Context,param sqlc.SetCourseLecturersParams)(sqlc.CoursesLecturer,error)
	UpdateCourseLecturers(ctx context.Context,param sqlc.UpdateCourseLecturersParams)(sqlc.CoursesLecturer,error)
	SetCoursePossibleVenues(ctx context.Context,courseVenueData []sqlc.SetCoursePossibleVenueParams)error
	SetCoursesForACohort(ctx context.Context,uniId uuid.UUID,cohortId uuid.UUID, courses[]uuid.UUID)error
	FetchCoursePossibleVenues(ctx context.Context,courseId uuid.UUID)([]sqlc.FetchCoursePossibleVenuesRow,error)
	DeleteCoursePossibleVenue(ctx context.Context,params sqlc.DeleteCoursePossibleVenueParams)error
	RetrieveCoursesForACohort(ctx context.Context,cohortId uuid.UUID)([]sqlc.RetrieveCoursesForACohortRow,error)
	RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.FetchAllCoursesRow,error)
}
type courseRepository struct {
	store sqlc.Store
	cq *queries.CoursesQueries
}


func NewCourseRepository(cq *queries.CoursesQueries, store sqlc.Store) *courseRepository{
	return &courseRepository{
		cq:cq,
		store:store,
	}
}


func (cq *courseRepository) CreateCourse(ctx context.Context,courseInfo sqlc.CreateCourseParams)(sqlc.Course,error){
	return cq.cq.CreateCourse(ctx,courseInfo)
}

func (cq *courseRepository) DeleteCourse(ctx context.Context,courseId uuid.UUID)error{
	return cq.cq.DeleteCourse(ctx,courseId)
}


// func (cq *courseRepository)

func (cq *courseRepository) RetrieveCoursesForADepartment(ctx context.Context,courseParam sqlc.RetrieveCoursesForADepartmentParams)(bool,[]sqlc.RetrieveCoursesForADepartmentRow,error){
	courses,err := cq.cq.RetrieveCoursesForDepartment(ctx,courseParam)
	if err != nil{
		if err == sql.ErrNoRows{
			return false,nil,err
		}
		return true,nil,err
	}
	return true,courses,nil
}


func(cq *courseRepository) UpdateCourse(ctx context.Context,courseInfo sqlc.UpdateCourseParams)(sqlc.Course,error){
	return cq.cq.UpdateCourse(ctx,courseInfo)
}


func (cq *courseRepository) SetStudentCourses(ctx context.Context,studCourseParam []sqlc.SetStudentCourseParams)error{
	return cq.store.ExecTx(ctx,func (q *sqlc.Queries) error {
		for i := 0; i < len(studCourseParam); i++ {
			_,err := q.SetStudentCourse(ctx,studCourseParam[i])
			if err != nil{
				return err;
			}
		}
		return nil
	})

}

func (cq *courseRepository) DeleteStudentCourse(ctx context.Context, deleteParam []sqlc.RemoveStudentCourseParams)error{
	return cq.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		for _,val := range deleteParam{
			_,err := q.RemoveStudentCourse(ctx,val)
			if err != nil {
				return err
			}
		}
		return nil
	})
}


func(cq *courseRepository) SetCourseLecturers(ctx context.Context,param sqlc.SetCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.cq.SetCourseLecturers(ctx,param)
}

func (cq *courseRepository) UpdateCourseLecturers(ctx context.Context,param sqlc.UpdateCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.cq.UpdateCourseLecturers(ctx,param)
}


func (cq *courseRepository) SetCoursesForACohort(ctx context.Context,uniId uuid.UUID,cohortId uuid.UUID, courses[]uuid.UUID)error{
	return cq.store.ExecTx(ctx,func(q *sqlc.Queries)error{
		for _,v := range courses{
			courseCohortData := sqlc.CreateCohortCourseParams{
				CohortID: cohortId,
				UniversityID: uniId,
				CourseID: v,
			}
			_,createCohortCourseErr := q.CreateCohortCourse(ctx,courseCohortData)
			if createCohortCourseErr != nil{
				return createCohortCourseErr
			}
		}
		return nil
	})
}

func (cq *courseRepository) SetCoursePossibleVenues(ctx context.Context,courseVenueData []sqlc.SetCoursePossibleVenueParams)error{
	return cq.store.ExecTx(ctx,func(q *sqlc.Queries) error {
		for _,val := range courseVenueData{
			err := q.SetCoursePossibleVenue(ctx,val)
			if err != nil{
				return err
			}
		}
		return nil
	})
}

func (cq *courseRepository) FetchCoursePossibleVenues(ctx context.Context,courseId uuid.UUID)([]sqlc.FetchCoursePossibleVenuesRow,error){
	return cq.cq.FetchCoursePossibleVenues(ctx,courseId)
}

func (cq *courseRepository) DeleteCoursePossibleVenue(ctx context.Context,params sqlc.DeleteCoursePossibleVenueParams)error{
	return cq.cq.DeleteCoursePossibleVenue(ctx,params)
}

func (cq *courseRepository) RetrieveCoursesForACohort(ctx context.Context,cohortId uuid.UUID)([]sqlc.RetrieveCoursesForACohortRow,error){
	return cq.cq.RetrieveCoursesForACohort(ctx,cohortId)
}

func (cq *courseRepository) RetrieveAllCourses(ctx context.Context, uniId uuid.UUID)([]sqlc.FetchAllCoursesRow,error){
	return cq.cq.FetchAllCourses(ctx,uniId)
}

// func (cq *courseRepository) Retrieve