package computed

import (
	"context"
	"time"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/timetable/repository"
	"github.com/google/uuid"
)


type Computed struct {
	timetableRepository repository.TimetableRepository
}

type modifiedCourseAndVenueData struct{
	CourseId uuid.UUID
	CourseCode string
	CourseTitle string
	CourseCreditiUnit int32
	CourseDuration int32
	DepartmentId uuid.UUID
	UniversityId uuid.UUID
	LecturerId uuid.NullUUID
	SessionsPerWeek int32
	Cohorts []uuid.UUID
	Level int32
	Semester string
	PossibleVenues []uuid.UUID
}

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

func ModifyCohortCourseData(cohortCourseData []sqlc.RetrieveCohortsForAllCoursesRow)map[uuid.UUID][]uuid.UUID{
 cohortCourseMap := make(map[uuid.UUID][]uuid.UUID)
 for _,v := range cohortCourseData{
	courseCohort, ok := cohortCourseMap[v.CourseID]

	if !ok {
		courseCohort = []uuid.UUID{v.CohortID}
		cohortCourseMap[v.CourseID] = courseCohort
	} else{
		courseCohort = append(courseCohort, v.CohortID)
		cohortCourseMap[v.CourseID] = courseCohort
	}
 }

 return cohortCourseMap


}

// format course data into the form modifiedCourseData
func ModifyCourseData(courseData []sqlc.RetrieveAllCoursesAndTheirVenueIdsRow) []modifiedCourseAndVenueData {
    courseDataMap := make(map[uuid.UUID]modifiedCourseAndVenueData)
    for _, v := range courseData {
        course, ok := courseDataMap[v.CourseID]
        if !ok {
            course = modifiedCourseAndVenueData{
                CourseId:         v.CourseID,
                CourseCode:       v.CourseCode,
                CourseTitle:      v.CourseTitle,
				CourseCreditiUnit: v.CourseCreditUnit,
                CourseDuration:   v.CourseDuration,
                SessionsPerWeek:  v.SessionsPerWeek,
                Semester:         v.Semester,
                PossibleVenues:   []uuid.UUID{v.VenueID},
            }
            courseDataMap[v.CourseID] = course
        } else {
            course.PossibleVenues = append(course.PossibleVenues, v.VenueID)
            courseDataMap[v.CourseID] = course
        }
    }

    // Convert map to slice
    modifiedData := make([]modifiedCourseAndVenueData, 0, len(courseDataMap))
    for _, course := range courseDataMap {
        modifiedData = append(modifiedData, course)
    }

    return modifiedData
}




func CreateSessionAtoms(lecturerMap map[uuid.UUID]int,venueMap map[uuid.UUID]int,courseMap map[uuid.UUID]int,cohortMap map[uuid.UUID]int,courseData []modifiedCourseAndVenueData, cohortCourseData map[uuid.UUID][]uuid.UUID)[]SessionAtom{

	sessionAtoms := make([]SessionAtom,0)
	counter := 0
	for _,v := range courseData{
		for range v.SessionsPerWeek{
			counter += 1
			sessionidx := counter
			courseIdx := courseMap[v.CourseId]
			// not handling the error here can hurt me in the future come back to check this if there is a problem
			validLecturerId,_ := utils.NullUUIDToUUid(v.LecturerId)
			lecturerIdx := lecturerMap[validLecturerId]
			sessionDuration := v.CourseDuration
			cohortIds := cohortCourseData[v.CourseId]
			cohortIdxs := make([]int,len(cohortIds))
			venueIdxs := make([]int,len(v.PossibleVenues))

			// creates the cohortsidx slice
			for i,val := range cohortIds{
				cohortIdxs[i] = cohortMap[val]
			}

			// creates the venueidx slice
			for i,val := range v.PossibleVenues{
				venueIdxs[i] = venueMap[val]
			}

			sessionAtoms = append(sessionAtoms, SessionAtom{
				SessionIdx: sessionidx,
				CourseIdx: courseIdx,
				LecturerIdx: lecturerIdx,
				CohortIdxs: cohortIdxs,
				SessionDuration: int(sessionDuration),
				AllowedVenuesIdx: venueIdxs,
			} )



		}
	}
	return sessionAtoms
}



func ComputeLecturerUnavailability(lectUnavailable []sqlc.RetrieveTotalLecturerUnavailabilityRow, lecturerMap map[uuid.UUID]int, days []string, slotDuration time.Duration, startOfDay time.Time, slotsPerDay int)[][]bool{
	totalSlots := len(days) * slotsPerDay

	lecturerUnavailable := make([][]bool,len(lecturerMap))
	for i := range lecturerUnavailable {
		lecturerUnavailable[i] = make([]bool, totalSlots)
	}

	dayIndex := make(map[string]int)
	for i,d := range days{
		dayIndex[d] = i
	}

	for _,row := range lectUnavailable{
		lectidx,exists := lecturerMap[row.LecturerID]
		if !exists{
			continue
		}
		didx,ok := dayIndex[row.Day]
		if !ok{
			continue
		}
		startSlot := int(row.StartTime.Sub(startOfDay)/slotDuration)
		endSlot := int(row.EndTime.Sub(startOfDay)/slotDuration)

		if startSlot < 0 {
			startSlot = 0
		}
		if endSlot > slotsPerDay {
			endSlot = slotsPerDay
		}

		for s:= startSlot; s < endSlot; s++{
			globalslot := didx * slotsPerDay + s
			lecturerUnavailable[lectidx][globalslot] = true
		}
	}

  return lecturerUnavailable
}


func ComputeVenueUnavaibility(venueUnavailable []sqlc.RetrieveTotalVenueUnavailabilityRow, venueMap map[uuid.UUID]int, days []string, slotDuration time.Duration,startOfDay time.Time, slotsPerDay int)[][]bool{
	totalSlots := len(days) * slotsPerDay

	venUnavailable := make([][]bool,len(venueMap))
	for i := range venUnavailable {
		venUnavailable[i] = make([]bool, totalSlots)
	}

	dayIndex := make(map[string]int)
	for i,d := range days{
		dayIndex[d] = i
	}

	for _,row := range venueUnavailable{
		venueIdx,exists := venueMap[row.VenueID]
		if !exists{
			continue
		}
		// remember all the possible errors that could occur here
		val,_ := utils.NullStringToString(row.Day)
		startTime,_ := utils.NullTimeToTime(row.StartTime)
		endTime,_ := utils.NullTimeToTime(row.EndTime)
		didx,ok := dayIndex[val]
		if !ok{
			continue
		}
		startSlot := int(startTime.Sub(startOfDay)/slotDuration)
		endSlot := int(endTime.Sub(startOfDay)/slotDuration)

		if startSlot < 0 {
			startSlot = 0
		}
		if endSlot > slotsPerDay {
			endSlot = slotsPerDay
		}

		for s:= startSlot; s < endSlot; s++{
			globalslot := didx * slotsPerDay + s
			venUnavailable[venueIdx][globalslot] = true
		}
	}

  return venUnavailable
}


func (c *Computed) ComputePreComputed(ctx context.Context, uniId uuid.UUID,slotsPerDay int, startOfDay time.Time,days[]string,slotDuration time.Duration)(*PreComputed,map[uuid.UUID]int,map[uuid.UUID]int,map[uuid.UUID]int,map[uuid.UUID]int){
	rawCohortData,_ := c.timetableRepository.RetrieveAllCohorts(ctx,uniId)
	rawCoursesData,_ := c.timetableRepository.RetrieveAllCourses(ctx,uniId)
	rawVenuesData,_ := c.timetableRepository.RetrieveAllVenues(ctx,uniId)
	rawLecturersData,_ := c.timetableRepository.RetrieveTotalLecturers(ctx,utils.UuidToNullUUID(uniId))
	rawCohortCourseData,_ := c.timetableRepository.RetrieveCohortsForAllCourses(ctx,uniId)
	rawCourseAndVenueData,_ := c.timetableRepository.RetrieveAllCoursesAndVenues(ctx,uniId)
	cohortMap := MapCohortIdToIdx(rawCohortData)
	venueMap := MapVanueIdToIdx(rawVenuesData)
	lecturerMap := MapLecturerIdToIdx(rawLecturersData)
	coursesMap := MapCoursesIdtoIdx(rawCoursesData)
	courseData := ModifyCourseData(rawCourseAndVenueData)
	cohortCourseData := ModifyCohortCourseData(rawCohortCourseData)
	rawLecturerUnavailability,_ := c.timetableRepository.RetrieveTotalLecturerUnavailability(ctx,utils.UuidToNullUUID(uniId))
	rawVenueUnavailability,_ := c.timetableRepository.RetrieveTotalVenueUnavailability(ctx,uniId)
	venueUnavailability := ComputeVenueUnavaibility(rawVenueUnavailability,venueMap,days,slotDuration,startOfDay,slotsPerDay)
	lecturerUnavailability := ComputeLecturerUnavailability(rawLecturerUnavailability,lecturerMap,days,slotDuration,startOfDay,slotsPerDay)
	sessionAtoms := CreateSessionAtoms(lecturerMap,venueMap,coursesMap,cohortMap,courseData,cohortCourseData)
	totalSlots := slotsPerDay * len(days)
	numVenues,_ := c.timetableRepository.CountNumVenues(ctx,uniId)
	numLecturers,_ := c.timetableRepository.CountNumLecturers(ctx,utils.UuidToNullUUID(uniId))
	numCohorts,_ := c.timetableRepository.CountNumCohorts(ctx,uniId)
	numCourses,_ := c.timetableRepository.CountNumCourses(ctx,uniId)

	return &PreComputed{
		TotalSlots: totalSlots,
		SlotsPerDay: slotsPerDay,
		NumVenues: int(numVenues),
		NumLecturers: int(numLecturers),
		NumCohorts: int(numCohorts),
		NumCourses: int(numCourses),
		SessionAtoms: sessionAtoms,
		LecturerUnavailable: lecturerUnavailability,
		VenueUnavailable: venueUnavailability,
	},cohortMap,venueMap,lecturerMap,coursesMap
}
