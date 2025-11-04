package service

import "github.com/Cxons/unischedulebackend/internal/timetable/repository"




type timetableRepository repository.TimetableRepository


type TimeTable interface{
	
}







type timeTableService struct{
	repo timetableRepository
}



func NewTimetableService(repo timetableRepository)*timeTableService{
	return &timeTableService{
		repo: repo,
	}
}















