package computed

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	sqlc "github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/internal/shared/utils"
	"github.com/Cxons/unischedulebackend/internal/timetable/repository"
	"github.com/google/uuid"
)


type Computed struct {
	timetableRepository repository.TimetableRepository
}
func NewComputed(repo repository.TimetableRepository) *Computed {
    return &Computed{
        timetableRepository: repo,
    }
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

func ModifyCourseData(courseData []sqlc.RetrieveAllCoursesAndTheirVenueIdsRow) []modifiedCourseAndVenueData {
    slog.Info("=== DEBUG ModifyCourseData START ===")
    
    courseDataMap := make(map[uuid.UUID]modifiedCourseAndVenueData)
    
    for i, v := range courseData {
        slog.Info("Processing course-venue row", 
            "index", i,
            "courseId", v.CourseID,
            "lecturerId", v.LecturerID,
            "lecturerValid", v.LecturerID.Valid,
            "lecturerUUID", v.LecturerID.UUID)
            
        course, ok := courseDataMap[v.CourseID]
        if !ok {
            // First time seeing this course - create new entry
            course = modifiedCourseAndVenueData{
                CourseId:          v.CourseID,
                CourseCode:        v.CourseCode,
                CourseTitle:       v.CourseTitle,
                CourseCreditiUnit: v.CourseCreditUnit,
                CourseDuration:    v.CourseDuration,
                SessionsPerWeek:   v.SessionsPerWeek,
                Semester:          v.Semester,
                PossibleVenues:    []uuid.UUID{v.VenueID},
                LecturerId:        v.LecturerID, // PRESERVE THE LECTURER ID
            }
            courseDataMap[v.CourseID] = course
            slog.Info("Created new course entry with lecturer", 
                "courseId", v.CourseID,
                "lecturerValid", v.LecturerID.Valid)
        } else {
            // Course already exists - just add the venue
            course.PossibleVenues = append(course.PossibleVenues, v.VenueID)
            courseDataMap[v.CourseID] = course
            slog.Info("Added venue to existing course", 
                "courseId", v.CourseID,
                "venueId", v.VenueID)
            // NOTE: We don't overwrite the lecturer ID here because it should be the same
        }
    }

    // Convert map to slice
    modifiedData := make([]modifiedCourseAndVenueData, 0, len(courseDataMap))
    for courseId, course := range courseDataMap {
        modifiedData = append(modifiedData, course)
        slog.Info("Final course data", 
            "courseId", courseId,
            "lecturerValid", course.LecturerId.Valid,
            "lecturerUUID", course.LecturerId.UUID,
            "venues", len(course.PossibleVenues))
    }

    slog.Info("=== DEBUG ModifyCourseData END ===", "courses", len(modifiedData))
    return modifiedData
}



func CreateSessionAtoms(lecturerMap map[uuid.UUID]int, venueMap map[uuid.UUID]int, courseMap map[uuid.UUID]int, cohortMap map[uuid.UUID]int, courseData []modifiedCourseAndVenueData, cohortCourseData map[uuid.UUID][]uuid.UUID) []SessionAtom {
    sessionAtoms := make([]SessionAtom, 0)
    counter := 0
	slog.Info("the course data","data",courseData)

    for _, v := range courseData {
        // Validate lecturer exists
        validLecturerId, err := utils.NullUUIDToUUid(v.LecturerId)
        if err != nil {
            slog.Warn("Course has invalid lecturer ID, skipping", "courseId", v.CourseId, "error", err)
            continue
        }

        lecturerIdx, lecturerExists := lecturerMap[validLecturerId]
        if !lecturerExists {
            slog.Warn("Lecturer not found in map, skipping course", "lecturerId", validLecturerId, "courseId", v.CourseId)
            continue
        }

        courseIdx, courseExists := courseMap[v.CourseId]
        if !courseExists {
            slog.Warn("Course not found in map, skipping", "courseId", v.CourseId)
            continue
        }

        // Get cohorts for this course
        cohortIds, hasCohorts := cohortCourseData[v.CourseId]
        if !hasCohorts || len(cohortIds) == 0 {
            slog.Warn("Course has no cohorts, skipping", "courseId", v.CourseId)
            continue
        }

        // Convert cohort IDs to indexes
        cohortIdxs := make([]int, 0, len(cohortIds))
        for _, cohortId := range cohortIds {
            if cohortIdx, exists := cohortMap[cohortId]; exists {
                cohortIdxs = append(cohortIdxs, cohortIdx)
            } else {
                slog.Warn("Cohort not found in map, skipping", "cohortId", cohortId, "courseId", v.CourseId)
            }
        }

        if len(cohortIdxs) == 0 {
            slog.Warn("No valid cohorts found for course, skipping", "courseId", v.CourseId)
            continue
        }

        // Convert venue IDs to indexes
        venueIdxs := make([]int, 0, len(v.PossibleVenues))
        for _, venueId := range v.PossibleVenues {
            if venueIdx, exists := venueMap[venueId]; exists {
                venueIdxs = append(venueIdxs, venueIdx)
            } else {
                slog.Warn("Venue not found in map, skipping", "venueId", venueId, "courseId", v.CourseId)
            }
        }

        if len(venueIdxs) == 0 {
            slog.Warn("No valid venues found for course, skipping", "courseId", v.CourseId)
            continue
        }

        // Create session atoms for each session per week
        for i := 0; i < int(v.SessionsPerWeek); i++ {
            counter++
            sessionAtoms = append(sessionAtoms, SessionAtom{
                SessionIdx:       counter - 1, // Use 0-based indexing
                CourseIdx:        courseIdx,
                LecturerIdx:      lecturerIdx,
                CohortIdxs:       cohortIdxs,
                SessionDuration:  int(v.CourseDuration),
                AllowedVenuesIdx: venueIdxs,
            })
        }
    }

    slog.Info("Created session atoms", "count", len(sessionAtoms))
    return sessionAtoms
}
// Add these helper functions to your computed package
func parseTimeFromString(timeStr string) (time.Time, error) {
    if timeStr == "" {
        return time.Time{}, fmt.Errorf("empty time string")
    }
    
    // Remove any timezone information if present
    timeStr = strings.Split(timeStr, "+")[0]
    timeStr = strings.Split(timeStr, "-")[0]
    timeStr = strings.TrimSpace(timeStr)
    
    // Try different time formats
    formats := []string{
        "15:04:05",
        "15:04:05.999999",
        "15:04",
        "2006-01-02 15:04:05",
        "2006-01-02T15:04:05",
        time.RFC3339,
    }
    
    for _, format := range formats {
        if t, err := time.Parse(format, timeStr); err == nil {
            // Extract only the time part (ignore date)
            return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC), nil
        }
    }
    
    return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
}


func ComputeLecturerUnavailability(lectUnavailable []sqlc.RetrieveTotalLecturerUnavailabilityRow, lecturerMap map[uuid.UUID]int, days []string, slotDuration time.Duration, startOfDay time.Time, slotsPerDay int) [][]bool {
    totalSlots := len(days) * slotsPerDay

    lecturerUnavailable := make([][]bool, len(lecturerMap))
    for i := range lecturerUnavailable {
        lecturerUnavailable[i] = make([]bool, totalSlots)
    }

    dayIndex := make(map[string]int)
    for i, d := range days {
        dayIndex[d] = i
    }

    for _, row := range lectUnavailable {
        lectidx, exists := lecturerMap[row.LecturerID]
        if !exists {
            slog.Debug("Lecturer not found in map", "lecturerId", row.LecturerID)
            continue
        }
        
        didx, ok := dayIndex[row.Day]
        if !ok {
            slog.Debug("Day not in configured days", "day", row.Day, "lecturerId", row.LecturerID)
            continue
        }
		startTime,_ := parseTimeFromString(row.StartTime)
		endTime,_ := parseTimeFromString(row.EndTime)
        
        startSlot := int(startTime.Sub(startOfDay) / slotDuration)
        endSlot := int(endTime.Sub(startOfDay) / slotDuration)

        // Validate and clamp slot ranges
        if startSlot < 0 {
            startSlot = 0
        }
        if endSlot > slotsPerDay {
            endSlot = slotsPerDay
        }
        if startSlot >= endSlot {
            slog.Debug("Invalid time range for lecturer unavailability", 
                "lecturerId", row.LecturerID, "startSlot", startSlot, "endSlot", endSlot)
            continue
        }

        // Mark unavailable slots
        for s := startSlot; s < endSlot; s++ {
            globalslot := didx*slotsPerDay + s
            if globalslot < totalSlots {
                lecturerUnavailable[lectidx][globalslot] = true
            }
        }
    }

    return lecturerUnavailable
}

func ComputeVenueUnavaibility(venueUnavailable []sqlc.RetrieveTotalVenueUnavailabilityRow, venueMap map[uuid.UUID]int, days []string, slotDuration time.Duration, startOfDay time.Time, slotsPerDay int) [][]bool {
    totalSlots := len(days) * slotsPerDay

    venUnavailable := make([][]bool, len(venueMap))
    for i := range venUnavailable {
        venUnavailable[i] = make([]bool, totalSlots)
    }

    dayIndex := make(map[string]int)
    for i, d := range days {
        dayIndex[d] = i
    }

    for _, row := range venueUnavailable {
		startTimeF,_ := parseTimeFromString(row.StartTime)
		endTimeF,_ := parseTimeFromString(row.EndTime)

        venueIdx, exists := venueMap[row.VenueID]
        if !exists {
            continue
        }
        
        // Check if we have valid day and time data
        // if !row.Day.Valid || startTime ==  || !row.EndTime.Valid {
        //     slog.Debug("Skipping venue unavailability row with null values", 
        //         "venueId", row.VenueID, 
        //         "dayValid", row.Day.Valid,
        //         "startTimeValid", row.StartTime.Valid,
        //         "endTimeValid", row.EndTime.Valid)
        //     continue
        // }
        
        dayStr := row.Day.String
        startTime := startTimeF
        endTime := endTimeF
        
        didx, ok := dayIndex[dayStr]
        if !ok {
            slog.Debug("Day not in configured days", "day", dayStr, "venueId", row.VenueID)
            continue
        }
        
        // Calculate slots
        startSlot := int(startTime.Sub(startOfDay) / slotDuration)
        endSlot := int(endTime.Sub(startOfDay) / slotDuration)

        // Validate and clamp slot ranges
        if startSlot < 0 {
            startSlot = 0
        }
        if endSlot > slotsPerDay {
            endSlot = slotsPerDay
        }
        if startSlot >= endSlot {
            slog.Debug("Invalid time range for venue unavailability", 
                "venueId", row.VenueID, "startSlot", startSlot, "endSlot", endSlot)
            continue
        }

        // Mark unavailable slots
        for s := startSlot; s < endSlot; s++ {
            globalslot := didx*slotsPerDay + s
            if globalslot < totalSlots {
                venUnavailable[venueIdx][globalslot] = true
            }
        }
    }

    return venUnavailable
}
func (c *Computed) ComputePreComputed(ctx context.Context, uniId uuid.UUID, slotsPerDay int, startOfDay time.Time, days []string, slotDuration time.Duration) (*PreComputed, map[uuid.UUID]int, map[uuid.UUID]int, map[uuid.UUID]int, map[uuid.UUID]int) {
    slog.Info("=== START ComputePreComputed DEBUG ===")
    slog.Info("Parameters", "universityId", uniId, "slotsPerDay", slotsPerDay, "days", days, "startOfDay", startOfDay, "slotDuration", slotDuration)

    // 1. DEBUG: Check raw courses data
    slog.Info("=== STEP 1: Checking raw courses ===")
    rawCoursesData, err := c.timetableRepository.RetrieveAllCourses(ctx, uniId)
    if err != nil {
        slog.Error("❌ Failed to retrieve courses", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Raw courses retrieved", "count", len(rawCoursesData))
    for i, course := range rawCoursesData {
        slog.Info("Raw course", 
            "index", i,
            "courseId", course.CourseID,
            "courseCode", course.CourseCode,
            "courseTitle", course.CourseTitle,
            "sessionsPerWeek", course.SessionsPerWeek,
            "duration", course.CourseDuration)
    }

    // 2. DEBUG: Check raw course-venue relationships (THIS IS CRITICAL)
    slog.Info("=== STEP 2: Checking course-venue relationships ===")
    rawCourseAndVenueData, err := c.timetableRepository.RetrieveAllCoursesAndVenues(ctx, uniId)
    if err != nil {
        slog.Error("❌ Failed to retrieve course-venue relationships", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Course-venue relationships retrieved", "count", len(rawCourseAndVenueData))
    
    targetCourseId := uuid.MustParse("1e71ea2b-d69e-4855-92e5-13a7095fb508")
    foundCourseVenueData := false
    for i, rel := range rawCourseAndVenueData {
        slog.Info("Course-venue relationship", 
            "index", i,
            "courseId", rel.CourseID,
            "venueId", rel.VenueID,
            "lecturerId", rel.LecturerID,
            "lecturerValid", rel.LecturerID.Valid,
            "lecturerUUID", rel.LecturerID.UUID)
        
        if rel.CourseID == targetCourseId {
            foundCourseVenueData = true
            slog.Info("✅ TARGET COURSE FOUND in course-venue data", 
                "courseId", rel.CourseID,
                "lecturerValid", rel.LecturerID.Valid,
                "lecturerUUID", rel.LecturerID.UUID)
        }
    }
    if !foundCourseVenueData {
        slog.Error("❌ TARGET COURSE NOT FOUND in course-venue relationships", "courseId", targetCourseId)
    }

    // 3. DEBUG: Check raw lecturers
    slog.Info("=== STEP 3: Checking lecturers ===")
    rawLecturersData, err := c.timetableRepository.RetrieveTotalLecturers(ctx, utils.UuidToNullUUID(uniId))
    if err != nil {
        slog.Error("❌ Failed to retrieve lecturers", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Lecturers retrieved", "count", len(rawLecturersData))
    targetLecturerId := uuid.MustParse("ef5ddd4c-3784-488c-90bb-392dc21b41c5")
    foundLecturer := false
    for i, lecturer := range rawLecturersData {
        slog.Info("Raw lecturer", 
            "index", i,
            "lecturerId", lecturer.LecturerID,
            "name", lecturer.LecturerFirstName)
        if lecturer.LecturerID == targetLecturerId {
            foundLecturer = true
            slog.Info("✅ TARGET LECTURER FOUND", "lecturerId", targetLecturerId)
        }
    }
    if !foundLecturer {
        slog.Error("❌ TARGET LECTURER NOT FOUND in raw data", "lecturerId", targetLecturerId)
    }

    // 4. DEBUG: Check cohorts for courses
    slog.Info("=== STEP 4: Checking cohort-course relationships ===")
    rawCohortCourseData, err := c.timetableRepository.RetrieveCohortsForAllCourses(ctx, uniId)
    if err != nil {
        slog.Error("❌ Failed to retrieve cohort-course relationships", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Cohort-course relationships retrieved", "count", len(rawCohortCourseData))
    foundCohorts := false
    for _, rel := range rawCohortCourseData {
        if rel.CourseID == targetCourseId {
            slog.Info("✅ TARGET COURSE HAS COHORT", "courseId", rel.CourseID, "cohortId", rel.CohortID)
            foundCohorts = true
        }
    }
    if !foundCohorts {
        slog.Error("❌ TARGET COURSE HAS NO COHORTS", "courseId", targetCourseId)
    }

    // 5. DEBUG: Check venues
    slog.Info("=== STEP 5: Checking venues ===")
    rawVenuesData, err := c.timetableRepository.RetrieveAllVenues(ctx, uniId)
    if err != nil {
        slog.Error("❌ Failed to retrieve venues", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Venues retrieved", "count", len(rawVenuesData))
    for i, venue := range rawVenuesData {
        slog.Info("Raw venue", 
            "index", i,
            "venueId", venue.VenueID,
            "venueName", venue.VenueName)
    }

    // 6. DEBUG: Check cohorts
    slog.Info("=== STEP 6: Checking cohorts ===")
    rawCohortData, err := c.timetableRepository.RetrieveAllCohorts(ctx, uniId)
    if err != nil {
        slog.Error("❌ Failed to retrieve cohorts", "error", err, "universityId", uniId)
        return nil, nil, nil, nil, nil
    }
    slog.Info("Cohorts retrieved", "count", len(rawCohortData))

    // Create mapping indexes
    slog.Info("=== STEP 7: Creating mapping indexes ===")
    cohortMap := MapCohortIdToIdx(rawCohortData)
    venueMap := MapVanueIdToIdx(rawVenuesData)
    lecturerMap := MapLecturerIdToIdx(rawLecturersData)
    coursesMap := MapCoursesIdtoIdx(rawCoursesData)

    // Log mapping statistics for debugging
    slog.Info("Mapping statistics",
        "cohorts", len(cohortMap),
        "venues", len(venueMap),
        "lecturers", len(lecturerMap),
        "courses", len(coursesMap))

    // Check if target lecturer exists in lecturerMap
    if lecturerIdx, exists := lecturerMap[targetLecturerId]; exists {
        slog.Info("✅ TARGET LECTURER FOUND in lecturerMap", "lecturerId", targetLecturerId, "index", lecturerIdx)
    } else {
        slog.Error("❌ TARGET LECTURER NOT FOUND in lecturerMap", "lecturerId", targetLecturerId)
        slog.Info("Available lecturers in map:")
        for lid, idx := range lecturerMap {
            slog.Info("Available lecturer", "lecturerId", lid, "index", idx)
        }
    }

    // Check if target course exists in coursesMap
    if courseIdx, exists := coursesMap[targetCourseId]; exists {
        slog.Info("✅ TARGET COURSE FOUND in coursesMap", "courseId", targetCourseId, "index", courseIdx)
    } else {
        slog.Error("❌ TARGET COURSE NOT FOUND in coursesMap", "courseId", targetCourseId)
    }

    // 8. DEBUG: Process course data
    slog.Info("=== STEP 8: Processing course data ===")
    courseData := ModifyCourseData(rawCourseAndVenueData)
    slog.Info("Processed course data", "count", len(courseData))
    
    foundProcessedCourse := false
    for i, course := range courseData {
        slog.Info("Processed course", 
            "index", i,
            "courseId", course.CourseId,
            "courseTitle", course.CourseTitle,
            "hasLecturer", course.LecturerId.Valid,
            "lecturerId", course.LecturerId.UUID,
            "sessionsPerWeek", course.SessionsPerWeek,
            "duration", course.CourseDuration,
            "possibleVenues", len(course.PossibleVenues))
        
        if course.CourseId == targetCourseId {
            foundProcessedCourse = true
            slog.Info("✅ TARGET COURSE FOUND in processed data", 
                "courseId", course.CourseId,
                "lecturerValid", course.LecturerId.Valid,
                "lecturerUUID", course.LecturerId.UUID)
            
            if !course.LecturerId.Valid {
                slog.Error("❌ PROBLEM: LecturerId.Valid is FALSE in processed data!")
            }
            if course.LecturerId.UUID == targetLecturerId {
                slog.Info("✅ Lecturer UUID matches target")
            } else {
                slog.Error("❌ Lecturer UUID mismatch", "expected", targetLecturerId, "actual", course.LecturerId.UUID)
            }
        }
    }
    if !foundProcessedCourse {
        slog.Error("❌ TARGET COURSE NOT FOUND in processed course data", "courseId", targetCourseId)
    }

    // 9. DEBUG: Process cohort data
    slog.Info("=== STEP 9: Processing cohort data ===")
    cohortCourseData := ModifyCohortCourseData(rawCohortCourseData)
    slog.Info("Cohort course data", "coursesWithCohorts", len(cohortCourseData))
    if cohorts, exists := cohortCourseData[targetCourseId]; exists {
        slog.Info("✅ TARGET COURSE HAS COHORTS in processed data", "courseId", targetCourseId, "cohortCount", len(cohorts))
        for _ , cohortId := range cohorts {
            if cohortIdx, exists := cohortMap[cohortId]; exists {
                slog.Info("Cohort found in cohortMap", "cohortId", cohortId, "index", cohortIdx)
            } else {
                slog.Warn("Cohort not found in cohortMap", "cohortId", cohortId)
            }
        }
    } else {
        slog.Error("❌ TARGET COURSE HAS NO COHORTS in processed data", "courseId", targetCourseId)
    }

    // 10. DEBUG: Create session atoms
    slog.Info("=== STEP 10: Creating session atoms ===")
    sessionAtoms := CreateSessionAtoms(lecturerMap, venueMap, coursesMap, cohortMap, courseData, cohortCourseData)
    slog.Info("Session atoms created", "count", len(sessionAtoms))

    // If no session atoms, debug why
    if len(sessionAtoms) == 0 {
        slog.Error("❌ CRITICAL: NO SESSION ATOMS CREATED")
        slog.Info("Debugging why session atoms are not created...")
        
        // Manually check each condition that would skip the course
        for _, course := range courseData {
            if course.CourseId == targetCourseId {
                slog.Info("=== MANUAL DEBUG FOR TARGET COURSE ===")
                
                // Check lecturer
                if !course.LecturerId.Valid {
                    slog.Error("❌ Skip reason: LecturerId.Valid = false")
                } else {
                    validLecturerId := course.LecturerId.UUID
                    if _, exists := lecturerMap[validLecturerId]; !exists {
                        slog.Error("❌ Skip reason: Lecturer not found in lecturerMap", "lecturerId", validLecturerId)
                    } else {
                        slog.Info("✅ Lecturer check PASSED")
                    }
                }
                
                // Check course in courseMap
                if _, exists := coursesMap[course.CourseId]; !exists {
                    slog.Error("❌ Skip reason: Course not found in courseMap")
                } else {
                    slog.Info("✅ Course map check PASSED")
                }
                
                // Check cohorts
                cohorts, hasCohorts := cohortCourseData[course.CourseId]
                if !hasCohorts || len(cohorts) == 0 {
                    slog.Error("❌ Skip reason: No cohorts")
                } else {
                    slog.Info("✅ Cohorts check PASSED", "cohortCount", len(cohorts))
                }
                
                // Check venues
                if len(course.PossibleVenues) == 0 {
                    slog.Error("❌ Skip reason: No possible venues")
                } else {
                    slog.Info("✅ Venues check PASSED", "venueCount", len(course.PossibleVenues))
                }
                
                // Check sessions per week
                if course.SessionsPerWeek <= 0 {
                    slog.Error("❌ Skip reason: Invalid sessions per week", "sessions", course.SessionsPerWeek)
                } else {
                    slog.Info("✅ Sessions per week check PASSED", "sessions", course.SessionsPerWeek)
                }
                
                // Check duration
                if course.CourseDuration <= 0 {
                    slog.Error("❌ Skip reason: Invalid duration", "duration", course.CourseDuration)
                } else {
                    slog.Info("✅ Duration check PASSED", "duration", course.CourseDuration)
                }
            }
        }
        
        return nil, nil, nil, nil, nil
    }

    // [Rest of your function remains the same...]
    totalSlots := slotsPerDay * len(days)
    numVenues := len(venueMap)
    numLecturers := len(lecturerMap)
    numCohorts := len(cohortMap)
    numCourses := len(coursesMap)

    // Retrieve availability data (with error handling)
    rawLecturerUnavailability, err := c.timetableRepository.RetrieveTotalLecturerUnavailability(ctx, utils.UuidToNullUUID(uniId))
    if err != nil {
        slog.Warn("Failed to retrieve lecturer unavailability, using empty data", "error", err)
        rawLecturerUnavailability = []sqlc.RetrieveTotalLecturerUnavailabilityRow{}
    }

    rawVenueUnavailability, err := c.timetableRepository.RetrieveTotalVenueUnavailability(ctx, uniId)
    if err != nil {
        slog.Warn("Failed to retrieve venue unavailability, using empty data", "error", err)
        rawVenueUnavailability = []sqlc.RetrieveTotalVenueUnavailabilityRow{}
    }

    // Compute availability matrices
    venueUnavailability := ComputeVenueUnavaibility(rawVenueUnavailability, venueMap, days, slotDuration, startOfDay, slotsPerDay)
    lecturerUnavailability := ComputeLecturerUnavailability(rawLecturerUnavailability, lecturerMap, days, slotDuration, startOfDay, slotsPerDay)

    // Create and validate PreComputed structure
    pre := &PreComputed{
        TotalSlots:          totalSlots,
        SlotsPerDay:         slotsPerDay,
        NumVenues:           numVenues,
        NumLecturers:        numLecturers,
        NumCohorts:          numCohorts,
        NumCourses:          numCourses,
        SessionAtoms:        sessionAtoms,
        LecturerUnavailable: lecturerUnavailability,
        VenueUnavailable:    venueUnavailability,
    }

    slog.Info("✅ PreComputed data successfully created", 
        "totalSlots", pre.TotalSlots,
        "sessionAtoms", len(pre.SessionAtoms),
        "lecturers", pre.NumLecturers,
        "venues", pre.NumVenues,
        "cohorts", pre.NumCohorts,
        "courses", pre.NumCourses)

    slog.Info("=== END ComputePreComputed DEBUG ===")
    return pre, cohortMap, venueMap, lecturerMap, coursesMap
}