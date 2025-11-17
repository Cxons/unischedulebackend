package service

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/dto"
	"github.com/Cxons/unischedulebackend/internal/timetable/computed"
	"github.com/Cxons/unischedulebackend/internal/timetable/repository"
	"github.com/Cxons/unischedulebackend/internal/timetable/types"
	status "github.com/Cxons/unischedulebackend/pkg/statuscodes"
	"github.com/google/uuid"
)




type timetableRepository repository.TimetableRepository
type customSessionPlacement = types.CustomSessionPlacement
type timeTableResponse = dto.ResponseDto

type SlotInfo struct {
	Day       string
	StartTime time.Time
}

type timeTableService struct{
	repo timetableRepository
	computed computed.Computed
	logger *slog.Logger
}

type TimetableSession struct {
	Time       string  
	CourseID   uuid.UUID 
	VenueID    uuid.UUID 
	SessionID  uuid.UUID 
	CourseName string 
	VenueName  string 
}

type TimeTableService interface{
	CreateATimeTable(ctx context.Context,startOfDay time.Time, endOfDay time.Time,uniId uuid.UUID)(timeTableResponse,string,error)
	RetrieveTimetableForACohort(ctx context.Context,cohortId uuid.UUID,uniId uuid.UUID)(timeTableResponse,string,error)
	RetrieveTimetableForAStudent(ctx context.Context,studentId uuid.UUID,uniId uuid.UUID,) (timeTableResponse, string, error) 
}


func NewTimetableService(repo timetableRepository,logger *slog.Logger)*timeTableService{
	return &timeTableService{
		repo: repo,
		computed: *computed.NewComputed(repo),
		logger: logger,
	}
}



// BuildSlotMap builds a map from slot index → (day, start time)
// BuildSlotMap builds a map from slot index → (day, start time)
func BuildSlotMap(
    slotsPerDay int,
    days []string,
    startOfDay time.Time,
    slotDuration time.Duration,
) map[int]SlotInfo {
    slotMap := make(map[int]SlotInfo)
    totalSlots := slotsPerDay * len(days)

    for i := 0; i < totalSlots; i++ {
        dayIdx := i / slotsPerDay
        slotIdxInDay := i % slotsPerDay

        // Make sure dayIdx is within bounds
        if dayIdx >= len(days) {
            continue
        }

        startTime := startOfDay.Add(time.Duration(slotIdxInDay) * slotDuration)

        slotMap[i] = SlotInfo{
            Day:       days[dayIdx],
            StartTime: startTime,
        }
    }

    return slotMap
}

func (tts *timeTableService) CreateATimeTable(ctx context.Context, startOfDay time.Time, endOfDay time.Time, uniId uuid.UUID) (timeTableResponse, string, error) {
    slotDuration := time.Hour
    totalDuration := endOfDay.Sub(startOfDay)
    slotsPerDay := int(totalDuration / slotDuration)

    slog.Info("set of values", "startofDay", startOfDay, "end of day", endOfDay, "totalduration", totalDuration, "slotsperday", slotsPerDay)
    
    // Validate slotsPerDay - return empty string for status message
    if slotsPerDay <= 0 {
        return timeTableResponse{}, "", fmt.Errorf("invalid time range: startOfDay=%v, endOfDay=%v", startOfDay, endOfDay)
    }
    
    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
    slotMap := BuildSlotMap(slotsPerDay, days, startOfDay, slotDuration)
    
    // Debug: Check if slotMap is populated
    tts.logger.Info("slotMap contents", "size", len(slotMap), "slotsPerDay", slotsPerDay, "days", len(days))
    if len(slotMap) == 0 {
        return timeTableResponse{}, "", fmt.Errorf("slotMap is empty - check time parameters")
    }

    precomputed, _, venueMap, _, coursesMap := tts.computed.ComputePreComputed(ctx, uniId, slotsPerDay, startOfDay, days, slotDuration)
    candidateTimetable := tts.computed.GeneticAlgorithm(precomputed)
    
    candidateData := sqlc.CreateCandidateParams{
        Fitness:           candidateTimetable.Fitness,
        UniversityID:      uniId,
        CandidateStatus:   "CURRENT",
        StartOfDay:        startOfDay,
        EndOfDay:          endOfDay,
    }
    
    slog.Info("candidate timetable", "val", candidateTimetable.Placements)
    
    sessionPlacements := make([]customSessionPlacement, 0)
    for _, val := range candidateTimetable.Placements {
        courseId := uuid.UUID{}
        venueId := uuid.UUID{}
        
        // Find course ID
        for key, value := range coursesMap {
            if val.CourseIdx == value {
                courseId = key
                break
            }
        }
        
        // Find venue ID
        for key, value := range venueMap {
            if val.VenueIdx == value {
                venueId = key
                break
            }
        }
        
        // Get slot info with bounds checking
        slotInfo, exists := slotMap[val.SlotIdx]
        if !exists {
            tts.logger.Warn("slot index not found in slotMap", "slotIdx", val.SlotIdx, "maxSlot", len(slotMap)-1)
            continue
        }
        
        // Validate slot info
        if slotInfo.Day == "" || slotInfo.StartTime.IsZero() {
            tts.logger.Warn("invalid slot info", "slotIdx", val.SlotIdx, "day", slotInfo.Day, "startTime", slotInfo.StartTime)
            continue
        }

        sessionPlacements = append(sessionPlacements, customSessionPlacement{
            SessionIdx:   int32(val.SessionIdx),
            CourseId:     courseId,
            VenueId:      venueId,
            Day:          slotInfo.Day,
            SessionTime:  slotInfo.StartTime,
            UniversityId: uniId,
        })
    }

    // Validate that we have session placements
    if len(sessionPlacements) == 0 {
        return timeTableResponse{}, "", fmt.Errorf("no valid session placements generated")
    }

    // Log the placements for debugging
    tts.logger.Info("session placements", "count", len(sessionPlacements), "firstPlacement", sessionPlacements[0])

    deprecateErr := tts.repo.DeprecateLatestCandidate(ctx, uniId)
    if deprecateErr != nil {
        tts.logger.Error("error deprecating latest candidate", "err", deprecateErr)
        return timeTableResponse{}, "", deprecateErr
    }

    err := tts.repo.CreateACandidateTimeTable(ctx, candidateData, sessionPlacements)
    if err != nil {
        restoreErr := tts.repo.RestoreCurrentCandidate(ctx, uniId)
        if restoreErr != nil {
            tts.logger.Error("error restoring current candidate", "err", restoreErr)
        }
        tts.logger.Error("error creating the candidate timetable", "err", err)
        return timeTableResponse{}, "", err
    }

    return timeTableResponse{
        Message:           "Timetable created successfully",
        StatusCode:        status.Created.Code,
        StatusCodeMessage: status.Created.Message,
    }, status.Created.Message, nil
}
// var dayOrder = map[string]int{
// 	"Monday":    1,
// 	"Tuesday":   2,
// 	"Wednesday": 3,
// 	"Thursday":  4,
// 	"Friday":    5,
// }


func PrepareTimetableJSON(
	sessions []sqlc.GetCohortSessionsInCurrentTimetableRow,
	courseNameMap map[uuid.UUID]string,
	venueNameMap map[uuid.UUID]string,
	startOfDay time.Time,
	endOfDay time.Time,
	slotDuration time.Duration,
) map[string][]TimetableSession {

	// Define day order
	dayOrder := map[string]int{
		"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4,
	}

	// Sort sessions by day and time
	sort.Slice(sessions, func(i, j int) bool {
		dayI := dayOrder[sessions[i].Day]
		dayJ := dayOrder[sessions[j].Day]

		if dayI == dayJ {
			return sessions[i].SessionTime.Before(sessions[j].SessionTime)
		}
		return dayI < dayJ
	})

	// Precompute all time slots for the day
	var slots []string
	for t := startOfDay; t.Before(endOfDay); t = t.Add(slotDuration) {
		slots = append(slots, t.Format("15:04"))
	}

	// Build a quick lookup for sessions
	sessionLookup := make(map[string]map[string]TimetableSession)
	for _, s := range sessions {
		timeStr := s.SessionTime.Format("15:04")

		if _, exists := sessionLookup[s.Day]; !exists {
			sessionLookup[s.Day] = make(map[string]TimetableSession)
		}

		sessionLookup[s.Day][timeStr] = TimetableSession{
			Time:       timeStr,
			CourseID:   s.CourseID,
			VenueID:    s.VenueID,
			SessionID:  s.SessionID,
			CourseName: courseNameMap[s.CourseID],
			VenueName:  venueNameMap[s.VenueID],
		}
	}

	// Fill grouped timetable with all days and slots
	grouped := make(map[string][]TimetableSession)
	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"} {
		for _, slot := range slots {
			if session, exists := sessionLookup[day][slot]; exists {
				grouped[day] = append(grouped[day], session)
			} else {
				// Empty slot
				grouped[day] = append(grouped[day], TimetableSession{
					Time:       slot,
					CourseName: "",
					VenueName:  "",
				})
			}
		}
	}

	return grouped
}


func PrepareStudentTimetableJSON(
	sessions []sqlc.GetStudentTimetableSessionsRow,
	courseNameMap map[uuid.UUID]string,
	venueNameMap map[uuid.UUID]string,
	startOfDay time.Time,
	endOfDay time.Time,
	slotDuration time.Duration,
) map[string][]TimetableSession {

	// Define day order
	dayOrder := map[string]int{
		"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4,
	}

	// Sort sessions by day and time
	sort.Slice(sessions, func(i, j int) bool {
		dayI := dayOrder[sessions[i].Day]
		dayJ := dayOrder[sessions[j].Day]

		if dayI == dayJ {
			return sessions[i].SessionTime.Before(sessions[j].SessionTime)
		}
		return dayI < dayJ
	})

	// Precompute all time slots for the day
	var slots []string
	for t := startOfDay; t.Before(endOfDay); t = t.Add(slotDuration) {
		slots = append(slots, t.Format("15:04"))
	}

	// Build a quick lookup for sessions
	sessionLookup := make(map[string]map[string]TimetableSession)
	for _, s := range sessions {
		timeStr := s.SessionTime.Format("15:04")

		if _, exists := sessionLookup[s.Day]; !exists {
			sessionLookup[s.Day] = make(map[string]TimetableSession)
		}

		sessionLookup[s.Day][timeStr] = TimetableSession{
			Time:       timeStr,
			CourseID:   s.CourseID,
			VenueID:    s.VenueID,
			SessionID:  s.SessionID,
			CourseName: courseNameMap[s.CourseID],
			VenueName:  venueNameMap[s.VenueID],
		}
	}

	// Fill grouped timetable with all days and slots
	grouped := make(map[string][]TimetableSession)
	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"} {
		for _, slot := range slots {
			if session, exists := sessionLookup[day][slot]; exists {
				grouped[day] = append(grouped[day], session)
			} else {
				// Empty slot
				grouped[day] = append(grouped[day], TimetableSession{
					Time:       slot,
					CourseName: "",
					VenueName:  "",
				})
			}
		}
	}

	return grouped
}


func (tts *timeTableService) RetrieveTimetableForACohort(ctx context.Context, cohortId uuid.UUID, uniId uuid.UUID) (timeTableResponse, string, error) {
    timetable, err := tts.repo.FetchSessionsForACohort(ctx, sqlc.GetCohortSessionsInCurrentTimetableParams{
        CohortID:      cohortId,
        UniversityID: uniId,
    })
    if err != nil {
        tts.logger.Error("error retrieving timetable for cohort", "err", err)
        return timeTableResponse{}, status.InternalServerError.Message, err
    }

    // Check if timetable is empty
    if len(timetable) == 0 {
        return timeTableResponse{
            Message:           "No timetable found for this cohort",
            Data:              make(map[string][]TimetableSession),
            StatusCode:        status.NotFound.Code,
            StatusCodeMessage: status.NotFound.Message,
        }, status.NotFound.Message, nil
    }

    courses, retrieveAllCourseErr := tts.repo.RetrieveAllCourses(ctx, uniId)
    if retrieveAllCourseErr != nil {
        tts.logger.Error("error retrieving all courses", "err", retrieveAllCourseErr)
        return timeTableResponse{}, status.InternalServerError.Message, retrieveAllCourseErr
    }

    venues, retrieveAllVenuesErr := tts.repo.RetrieveAllVenues(ctx, uniId)
    if retrieveAllVenuesErr != nil {
        tts.logger.Error("error retrieving all venues", "err", retrieveAllVenuesErr)
        return timeTableResponse{}, status.InternalServerError.Message, retrieveAllVenuesErr
    }

    courseNameMap := make(map[uuid.UUID]string)
    venuesNameMap := make(map[uuid.UUID]string)
    
    for _, val := range courses {
        courseNameMap[val.CourseID] = val.CourseTitle
    }
    for _, val := range venues {
        venuesNameMap[val.VenueID] = val.VenueName
    }

    formattedTimetable := PrepareTimetableJSON(timetable, courseNameMap, venuesNameMap, timetable[0].StartOfDay, timetable[0].EndOfDay, time.Hour)
	// slog.Info("formattedtimetable","val",formattedTimetable)
    
    return timeTableResponse{
        Message:           "Cohort Timetable",
        Data:              formattedTimetable,
        StatusCode:        status.OK.Code,
        StatusCodeMessage: status.OK.Message,
    }, status.OK.Message, nil
}

func (tts *timeTableService) RetrieveTimetableForAStudent(
	ctx context.Context,
	studentId uuid.UUID,
	uniId uuid.UUID,
) (timeTableResponse, string, error) {
	timetable, err := tts.repo.FetchSessionsForAStudent(ctx,studentId)
	if err != nil {
		tts.logger.Error("error retrieving timetable for student", "err", err)
		return timeTableResponse{}, status.InternalServerError.Message, err
	}

	if len(timetable) == 0 {
		return timeTableResponse{
			Message:           "No timetable found for this student",
			StatusCode:        status.NotFound.Code,
			StatusCodeMessage: status.NotFound.Message,
		}, status.NotFound.Message, nil
	}

	courses, retrieveAllCourseErr := tts.repo.RetrieveAllCourses(ctx, uniId)
	if retrieveAllCourseErr != nil {
		tts.logger.Error("error retrieving all courses", "err", retrieveAllCourseErr)
		return timeTableResponse{}, status.InternalServerError.Message, retrieveAllCourseErr
	}

	venues, retrieveAllVenuesErr := tts.repo.RetrieveAllVenues(ctx, uniId)
	if retrieveAllVenuesErr != nil {
		tts.logger.Error("error retrieving all venues", "err", retrieveAllVenuesErr)
		return timeTableResponse{}, status.InternalServerError.Message, retrieveAllVenuesErr
	}

	courseNameMap := make(map[uuid.UUID]string)
	venuesNameMap := make(map[uuid.UUID]string)
	for _, val := range courses {
		courseNameMap[val.CourseID] = val.CourseTitle
	}
	for _, val := range venues {
		venuesNameMap[val.VenueID] = val.VenueName
	}

	formattedTimetable := PrepareStudentTimetableJSON(
		timetable,
		courseNameMap,
		venuesNameMap,
		timetable[0].StartOfDay,
		timetable[0].EndOfDay,
		time.Hour,
	)

	return timeTableResponse{
		Message:           "Student Timetable",
		Data:              formattedTimetable,
		StatusCode:        status.OK.Code,
		StatusCodeMessage: status.OK.Message,
	}, status.OK.Message, nil
}













