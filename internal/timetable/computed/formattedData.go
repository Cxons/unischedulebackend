package computed

import (
	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/google/uuid"
)

// maps all cohorts ids to idx
func MapCohortIdToIdx(cohorts []sqlc.Cohort)map[uuid.UUID]int{
	cohortMap := make(map[uuid.UUID]int)
	for i,v := range cohorts{
		cohortMap[v.CohortID] = i
	}
	return cohortMap
}

// maps all courses ids to idx
func MapCoursesIdtoIdx(courses []sqlc.RetrieveAllCoursesRow)map[uuid.UUID]int{
	coursesMap := make(map[uuid.UUID]int)
	for i,v := range courses{
		coursesMap[v.CourseID] = i
	}
	return coursesMap
}

// maps all venues ids to idx
func MapVanueIdToIdx(venues []sqlc.RetrieveAllVenuesRow)map[uuid.UUID]int{
	venuesMap := make(map[uuid.UUID]int)
	for i,v := range venues{
		venuesMap[v.VenueID] = i
	}
	return venuesMap
}

// maps all lecturers ids to idx
func MapLecturerIdToIdx(lecturers []sqlc.RetrieveTotalLecturersRow)map[uuid.UUID]int{
	lecturersMap := make(map[uuid.UUID]int)
	for i,v := range lecturers{
		lecturersMap[v.LecturerID] = i
	}
	return lecturersMap
}

// things that i need to do are these
/*
1. Create a table for courses and cohorts
2. fetch all the cohorts attached with a particular course
3. Convert cohort lecturer and course into idxs
4. function to create a session atom which involves looping over courses and their number of times offered per week

*/



// function to get number of venues

//function to get number of lecturers

// function to get number of cohorts

// function to get total number of courses

// function that retrieves all the courses alongside their allowed venues

// function that retrieves all the lecturers and unavailablity

// function that retrieves all the venues and their unavailablities

// function that maps the following ids to idx
/*
cohorts, lecturers , courses,venues
*/