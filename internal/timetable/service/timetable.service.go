package service

import (
	"math"
	"math/rand"
	"sort"

	"github.com/google/uuid"
)



type FixedTimeTableConfig struct {
	DaysPerWeek int
	SlotsPerDay int
	SlotDuration int
	DayStartHour int // e.g 8 for 8am 
	TotalSlots int 
}


type CourseOffering struct {
	OfferingId uuid.UUID
	CourseId uuid.UUID
	UniversityId uuid.UUID
	DepartmentId uuid.UUID
	LecturerId uuid.UUID
	Cohorts []uuid.UUID // array of cohorts(groups) offering the course
	SessionsPerWeek int
	SessionDuration int  // so a 2 hour course is 2 and has 2 slots
	AllowedVenues []uuid.UUID
	CreditUnit int

}

type Session struct {
	SessionId uuid.UUID
	CourseId uuid.UUID
	Cohorts []uuid.UUID
	LecturerId uuid.UUID
	SessionDuration int // how long for each session e.g 2 for 2 hours
	NoPerWeek int// how many times the course should hold
	AllowedVenues []uuid.UUID
}



type Cohort struct {
	CohortId uuid.UUID
	CohortName string
	UniversityId uuid.UUID
	DepartmentId uuid.UUID
	Level int

}


type Venue struct {
	VenueId uuid.UUID
	UniversityId uuid.UUID
	Capacity int // how much the class can hold
}


type Lecturer struct {
	LecturerId uuid.UUID
	UniversityId uuid.UUID
	LecturerUnavailability []bool // the length is the total slots allowed
}

 // there should be a map of uuids to the indexes here..indicates current session placement
type SessionPlacement struct {
	SessionIdx int
	VenueIdx int 
	SlotIdx int
	Conflict bool
}


// just shows the session information to be used in placement during computation
type SessionAtom struct {
	SessionIdx int
	CourseIdx int
	LecturerIdx int
	CohortIdxs []int
	SessionDuration int // how long for each session e.g 2 for 2 hours
	AllowedVenuesIdx []int
}


// shows all the necessary things i need to compute before starting the computation
type PreComputed struct {
	TotalSlots     int
	SlotsPerDay    int
	NumVenues      int
	NumLecturers   int
	NumCohorts     int
	SessionAtoms   []SessionAtom
	LecturerUnavailable   [][]bool // LecturerUnavailable[lecturerIdx][slot] static forbidden mask (true = unavailable)
	VenueUnavailable      [][]bool // VenueUnavailable[venueIdx][slot] static forbidden mask (true = unavailable)
}

type FeasiblePair struct {
	SlotIdx    int
	VenueIdx   int
	Score   float64 // lower = better
	Reasons string  // for debug
}

type Candidate struct {
	Placements []SessionPlacement
}


// k here refers to the number of pairs to choose from so if k = 3, it means choose one random from the top 3
func ChooseTopSampleK( topPairs[]FeasiblePair, k int, r *rand.Rand) FeasiblePair{
	if len(topPairs) < 1 {
		panic("ChooseTopSampleK was called with a an empty array")
	}
	if k <= 0 {
		k = 1;
	}
	if len(topPairs) <= k {
		return topPairs[rand.Intn(len(topPairs))]
	}
	return topPairs[rand.Intn(k)]
}



// next function is to compute top feasible pairs 
func ComputeFeasiblePairs(pre *PreComputed, session *SessionAtom, venueOccupied[][]bool, lecturerOccupied[][]bool, cohortOccupied[][]bool)[]FeasiblePair{
	totalSlots := pre.TotalSlots
	feasible := make([]FeasiblePair,0,totalSlots)

	// iterates through all possible slots
	for start := 0; start < totalSlots; start++{

		// prevents cross day boundary i.e a session crossing a day
		startDaySlot := ( start / pre.SlotsPerDay ) * pre.SlotsPerDay
		EndDaySlot := startDaySlot + pre.SlotsPerDay;

		if start + session.SessionDuration > EndDaySlot{
			continue
		}


		// prevents slotting in periods of unavailable lecturers for the session
		lectOk := true
		s := start + session.SessionDuration
		for si := start; si < s; si++{
			if si >= totalSlots || pre.LecturerUnavailable[session.LecturerIdx][si] || lecturerOccupied[session.LecturerIdx][si]{
				lectOk = false
				break
			}
		}
		if !lectOk{
			continue
		}


		// now prevents cohorts conflict 
		cohortsConflict := false
		for ci := start; ci < s; ci++{
			for _,c := range session.CohortIdxs{
				if cohortOccupied[c][ci]{
					cohortsConflict = true
					break
				}
			}
			if cohortsConflict{
				break
			}
		}

		if cohortsConflict{
			continue
		}


		// prevents venue conflict 
			for _,v := range session.AllowedVenuesIdx{
				venueOk := true

				for vi := start; vi < s; vi++{
				if venueOccupied[v][vi] || pre.VenueUnavailable[v][vi]{
					venueOk = false
					break
				}
			}
			if !venueOk{
				continue
			}
			feasible = append(feasible,FeasiblePair{
				SlotIdx: start,
				VenueIdx: v,
				Score: 0.0,
				Reasons: "",
			})
		}
		
	}
	return feasible
}


func ComputeLeastBadPair(pre *PreComputed, session *SessionAtom,venueOccupied[][]bool,lecturerOccupied[][]bool,cohortOccupied[][]bool) FeasiblePair{
	totalSlots := pre.TotalSlots

	bestPair := FeasiblePair{
		Score: math.MaxInt64,
	}

	for start := 0; start < totalSlots; start++{
		
		// prevents cross day boundary i.e a session crossing a day
		startDaySlot := ( start / pre.SlotsPerDay ) * pre.SlotsPerDay
		EndDaySlot := startDaySlot + pre.SlotsPerDay;

		if start + session.SessionDuration > EndDaySlot{
			continue
		}
		
		for _,v := range session.AllowedVenuesIdx{
			conflictScore := 0.0

			s := start + session.SessionDuration

			for si := start; si < s; si++{

				// this is a hard constraint i.e lecturer not available during this time
				if pre.LecturerUnavailable[session.LecturerIdx][si]{
					conflictScore += 1500
				}


				// this is a medium constraint this is during the making of the  partially formed timetable
				if lecturerOccupied[session.LecturerIdx][si]{
					conflictScore += 1000
				}

				// for times when the venue is not available so medium constraint
				if pre.VenueUnavailable[v][si] {
					conflictScore += 1500
				}


				if venueOccupied[v][si]{
					conflictScore += 1000
				}

				// decently hard constraint also manages cohorts conflicts
				for _,c := range session.CohortIdxs{
					if cohortOccupied[c][si]{
					conflictScore += 1500
					}
				}
			}
			if conflictScore < bestPair.Score{
				bestPair.Score = conflictScore
				bestPair.VenueIdx = v
				bestPair.SlotIdx = start
			}

		}
	}
	return bestPair
}


// this function builds a candidate timetable 
func BuildOneCandidate(r *rand.Rand, pre *PreComputed, k int) *Candidate{
	totalSessions := len(pre.SessionAtoms)
	totalVenues := pre.NumVenues
	totalLecturers := pre.NumLecturers
	totalSlots := pre.TotalSlots
	totalCohorts := pre.NumCohorts

	placements := make([]SessionPlacement,totalSessions)
   


	// create an order slice to be used to order session placement 
	order := make([]int, totalSessions)
	
	// input appropriate matching indexes into order array
	for i := 0; i < totalSessions; i++ {
		order[i] = i
	}

	// sort the order array according giving preference to sessions with less allowed venues and longer durations
	sort.Slice(order, func(i, j int) bool {
		a := pre.SessionAtoms[order[i]]
		b := pre.SessionAtoms[order[j]]
		scoreA := float64(len(a.AllowedVenuesIdx))*0.5 + float64(a.SessionDuration)*1.0
		scoreB := float64(len(b.AllowedVenuesIdx))*0.5 + float64(b.SessionDuration)*1.0
		return scoreA < scoreB // less allowed venues = harder -> come earlier
	})


	// compute venueOccupied 
	venueOcc := make([][]bool,totalVenues)
	for i := range venueOcc{
		venueOcc[i] = make([]bool,totalSlots)
	}


	// compute lecturerOccupied
	lecturerOcc := make([][]bool,totalLecturers)
	for i := range lecturerOcc{
		lecturerOcc[i] = make([]bool,totalSlots)
	}


	// compute cohortsOccupied
	cohortOcc := make([][]bool,totalCohorts)
	for i := range cohortOcc{
		cohortOcc[i] = make([]bool,totalSlots)
	}


	// placement of sessions into appropriate slots and venue in order
	for _,sessionIdx := range order{

		session := &pre.SessionAtoms[sessionIdx]

		feasible := ComputeFeasiblePairs(pre,session,venueOcc,lecturerOcc,cohortOcc)

		

		// if there are feasible pairs 
		if len(feasible) > 0{
			chosen := ChooseTopSampleK(feasible,k,r)
			placements[sessionIdx] = SessionPlacement{
				SessionIdx: sessionIdx,
				VenueIdx: chosen.VenueIdx,
				SlotIdx: chosen.SlotIdx,
				Conflict: false,
		}

		// mark occupancy
		for d:= 0; d < session.SessionDuration; d++{
			// si means slot index
			si := chosen.SlotIdx + d
			venueOcc[chosen.VenueIdx][si] = true
			lecturerOcc[session.LecturerIdx][si] = true
			for _,c := range session.CohortIdxs{
				cohortOcc[c][si] = true
			}
		}
		}else {
			// fallback: selects least bad pair
			best := ComputeLeastBadPair(pre,session,venueOcc,lecturerOcc,cohortOcc)
			placements[sessionIdx] = SessionPlacement{
				SessionIdx: sessionIdx,
				VenueIdx: best.VenueIdx,
				SlotIdx: best.SlotIdx,
				Conflict: true,
		}

		// mark tempoary occupancy
		for d:= 0; d < session.SessionDuration; d++{
			// si means slot index
			si := best.SlotIdx + d
			if si >= totalSlots{
				break
			}
			venueOcc[best.VenueIdx][si] = true
			lecturerOcc[session.LecturerIdx][si] = true
			for _,c := range session.CohortIdxs{
				cohortOcc[c][si] = true
			}
		}
		}
	}

	return &Candidate{
		Placements: placements,
	}

}


// func BuildPopulation()








