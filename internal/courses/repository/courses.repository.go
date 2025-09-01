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


func(cq *courseRepository) SetCourseLecturers(ctx context.Context,param sqlc.SetCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.cq.SetCourseLecturers(ctx,param)
}

func (cq *courseRepository) UpdateCourseLecturers(ctx context.Context,param sqlc.UpdateCourseLecturersParams)(sqlc.CoursesLecturer,error){
	return cq.cq.UpdateCourseLecturers(ctx,param)
}