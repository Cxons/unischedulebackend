package computed

import (
	"log/slog"
	"math"
	"math/rand"
	"sort"
	"time"
)

type SessionPlacement struct {
	SessionIdx int
	CourseIdx  int
	VenueIdx   int
	SlotIdx    int
	Conflict   bool
	Score      float64 // the lower the better
}

// just shows the session information to be used in placement during computation
type SessionAtom struct {
	SessionIdx       int
	CourseIdx        int
	LecturerIdx      int
	CohortIdxs       []int
	SessionDuration  int // how long for each session e.g 2 for 2 hours
	AllowedVenuesIdx []int
}

// shows all the necessary things i need to compute before starting the computation
type PreComputed struct {
	TotalSlots          int
	SlotsPerDay         int
	NumVenues           int
	NumLecturers        int
	NumCohorts          int
	NumCourses          int
	SessionAtoms        []SessionAtom
	LecturerUnavailable [][]bool // LecturerUnavailable[lecturerIdx][slot] static forbidden mask (true = unavailable)
	VenueUnavailable    [][]bool // VenueUnavailable[venueIdx][slot] static forbidden mask (true = unavailable)
}

type FeasiblePair struct {
	SlotIdx  int
	VenueIdx int
	Score    float64 // lower = better
	Reasons  string  // for debug
}

type Candidate struct {
	Placements []SessionPlacement
	Fitness    float64 // the higher the better
}

// k here refers to the number of pairs to choose from so if k = 3, it means choose one random from the top 3
func ChooseTopSampleK(topPairs []FeasiblePair, k int, r *rand.Rand) FeasiblePair {
	if len(topPairs) < 1 {
		panic("ChooseTopSampleK was called with a an empty array")
	}
	if k <= 0 {
		k = 1
	}
	if len(topPairs) <= k {
		return topPairs[r.Intn(len(topPairs))]
	}
	return topPairs[r.Intn(k)]
}

// next function is to compute top feasible pairs
func ComputeFeasiblePairs(pre *PreComputed, session *SessionAtom, venueOccupied [][]bool, lecturerOccupied [][]bool, cohortOccupied [][]bool) []FeasiblePair {
    totalSlots := pre.TotalSlots
    feasible := make([]FeasiblePair, 0, totalSlots)

    // Add bounds checking for session indices
    if session.LecturerIdx < 0 || session.LecturerIdx >= len(pre.LecturerUnavailable) {
        slog.Error("Invalid lecturer index", "lecturerIdx", session.LecturerIdx, "max", len(pre.LecturerUnavailable))
        return feasible // Return empty if invalid
    }

    // iterates through all possible slots
    for start := 0; start < totalSlots; start++ {
        // prevents cross day boundary i.e a session crossing a day
        startDaySlot := (start / pre.SlotsPerDay) * pre.SlotsPerDay
        EndDaySlot := startDaySlot + pre.SlotsPerDay

        if start+session.SessionDuration > EndDaySlot {
            continue
        }

        // prevents slotting in periods of unavailable lecturers for the session
        lectOk := true
        s := start + session.SessionDuration
        
        // Add bounds checking for lecturer unavailable
        for si := start; si < s; si++ {
            if si >= totalSlots {
                lectOk = false
                break
            }
            // Check if lecturer unavailable array has enough length
            if si >= len(pre.LecturerUnavailable[session.LecturerIdx]) {
                lectOk = false
                break
            }
            if pre.LecturerUnavailable[session.LecturerIdx][si] || 
               (si < len(lecturerOccupied[session.LecturerIdx]) && lecturerOccupied[session.LecturerIdx][si]) {
                lectOk = false
                break
            }
        }
        if !lectOk {
            continue
        }

        // now prevents cohorts conflict
        cohortsConflict := false
        for ci := start; ci < s; ci++ {
            if ci >= totalSlots {
                cohortsConflict = true
                break
            }
            for _, c := range session.CohortIdxs {
                // Add bounds checking for cohort index
                if c < 0 || c >= len(cohortOccupied) {
                    cohortsConflict = true
                    break
                }
                if ci < len(cohortOccupied[c]) && cohortOccupied[c][ci] {
                    cohortsConflict = true
                    break
                }
            }
            if cohortsConflict {
                break
            }
        }

        if cohortsConflict {
            continue
        }

        // prevents venue conflict
        for _, v := range session.AllowedVenuesIdx {
            // Add bounds checking for venue index
            if v < 0 || v >= len(venueOccupied) {
                continue
            }
            
            venueOk := true
            for vi := start; vi < s; vi++ {
                if vi >= totalSlots {
                    venueOk = false
                    break
                }
                // Check venue unavailable bounds
                if v < len(pre.VenueUnavailable) && vi < len(pre.VenueUnavailable[v]) && pre.VenueUnavailable[v][vi] {
                    venueOk = false
                    break
                }
                if vi < len(venueOccupied[v]) && venueOccupied[v][vi] {
                    venueOk = false
                    break
                }
            }
            if !venueOk {
                continue
            }
            feasible = append(feasible, FeasiblePair{
                SlotIdx:  start,
                VenueIdx: v,
                Score:    0.0,
                Reasons:  "",
            })
        }
    }
    return feasible
}
func ComputeLeastBadPair(pre *PreComputed, session *SessionAtom, venueOccupied [][]bool, lecturerOccupied [][]bool, cohortOccupied [][]bool) FeasiblePair {
    totalSlots := pre.TotalSlots

    bestPair := FeasiblePair{
        Score: math.MaxFloat64,
    }

    // Add bounds checking for session indices
    if session.LecturerIdx < 0 || session.LecturerIdx >= len(pre.LecturerUnavailable) {
        slog.Error("Invalid lecturer index in ComputeLeastBadPair", "lecturerIdx", session.LecturerIdx)
        return bestPair
    }

    for start := 0; start < totalSlots; start++ {
        // prevents cross day boundary i.e a session crossing a day
        startDaySlot := (start / pre.SlotsPerDay) * pre.SlotsPerDay
        EndDaySlot := startDaySlot + pre.SlotsPerDay

        if start+session.SessionDuration > EndDaySlot {
            continue
        }

        for _, v := range session.AllowedVenuesIdx {
            // Skip invalid venue indices
            if v < 0 || v >= len(venueOccupied) {
                continue
            }

            conflictScore := 0.0
            s := start + session.SessionDuration

            for si := start; si < s; si++ {
                if si >= totalSlots {
                    conflictScore += 1500 // Penalty for out of bounds
                    continue
                }

                // Lecturer unavailable check with bounds
                if si < len(pre.LecturerUnavailable[session.LecturerIdx]) {
                    if pre.LecturerUnavailable[session.LecturerIdx][si] {
                        conflictScore += 10
                    }
                }

                // Lecturer occupied check with bounds
                if session.LecturerIdx < len(lecturerOccupied) && si < len(lecturerOccupied[session.LecturerIdx]) {
                    if lecturerOccupied[session.LecturerIdx][si] {
                        conflictScore += 1500
                    }
                }

                // Venue unavailable check with bounds
                if v < len(pre.VenueUnavailable) && si < len(pre.VenueUnavailable[v]) {
                    if pre.VenueUnavailable[v][si] {
                        conflictScore += 10
                    }
                }

                // Venue occupied check with bounds
                if v < len(venueOccupied) && si < len(venueOccupied[v]) {
                    if venueOccupied[v][si] {
                        conflictScore += 1500
                    }
                }

                // Cohort occupied check with bounds
                for _, c := range session.CohortIdxs {
                    if c >= 0 && c < len(cohortOccupied) && si < len(cohortOccupied[c]) {
                        if cohortOccupied[c][si] {
                            conflictScore += 500
                        }
                    } else {
                        conflictScore += 500 // Penalty for invalid cohort index
                    }
                }
            }
            
            if conflictScore < bestPair.Score {
                bestPair.Score = conflictScore
                bestPair.VenueIdx = v
                bestPair.SlotIdx = start
            }
        }
    }
    return bestPair
}

// this function builds a candidate timetable
func BuildOneCandidate(r *rand.Rand, pre *PreComputed, k int) *Candidate {
	 if pre == nil {
        slog.Error("PreComputed is nil")
        return &Candidate{Placements: []SessionPlacement{}}
    }
    
    if len(pre.SessionAtoms) == 0 {
        slog.Error("No session atoms in PreComputed")
        return &Candidate{Placements: []SessionPlacement{}}
    }
    
    slog.Info("Building candidate", 
        "totalSessions", len(pre.SessionAtoms),
        "totalVenues", pre.NumVenues,
        "totalLecturers", pre.NumLecturers,
        "totalSlots", pre.TotalSlots,
        "totalCohorts", pre.NumCohorts)
	totalSessions := len(pre.SessionAtoms)
	totalVenues := pre.NumVenues
	totalLecturers := pre.NumLecturers
	totalSlots := pre.TotalSlots
	totalCohorts := pre.NumCohorts

	placements := make([]SessionPlacement, totalSessions)

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
	venueOcc := make([][]bool, totalVenues)
	for i := range venueOcc {
		venueOcc[i] = make([]bool, totalSlots)
	}

	// compute lecturerOccupied
	lecturerOcc := make([][]bool, totalLecturers)
	for i := range lecturerOcc {
		lecturerOcc[i] = make([]bool, totalSlots)
	}

	// compute cohortsOccupied
	cohortOcc := make([][]bool, totalCohorts)
	for i := range cohortOcc {
		cohortOcc[i] = make([]bool, totalSlots)
	}

	// placement of sessions into appropriate slots and venue in order
	for _, sessionIdx := range order {
		if sessionIdx < 0 || sessionIdx >= len(pre.SessionAtoms) {
			continue
		}

		session := &pre.SessionAtoms[sessionIdx]

		feasible := ComputeFeasiblePairs(pre, session, venueOcc, lecturerOcc, cohortOcc)

		// if there are feasible pairs
		if len(feasible) > 0 {
			chosen := ChooseTopSampleK(feasible, k, r)
			placements[sessionIdx] = SessionPlacement{
				SessionIdx: sessionIdx,
				CourseIdx:  session.CourseIdx,
				VenueIdx:   chosen.VenueIdx,
				SlotIdx:    chosen.SlotIdx,
				Conflict:   false,
				Score:      0.0,
			}

			// mark occupancy
			for d := 0; d < session.SessionDuration; d++ {
				// si means slot index
				si := chosen.SlotIdx + d
				if si >= totalSlots {
					break
				}
				if chosen.VenueIdx < len(venueOcc) && si < len(venueOcc[chosen.VenueIdx]) {
					venueOcc[chosen.VenueIdx][si] = true
				}
				if session.LecturerIdx < len(lecturerOcc) && si < len(lecturerOcc[session.LecturerIdx]) {
					lecturerOcc[session.LecturerIdx][si] = true
				}
				for _, c := range session.CohortIdxs {
					if c < len(cohortOcc) && si < len(cohortOcc[c]) {
						cohortOcc[c][si] = true
					}
				}
			}
		} else {
			// fallback: selects least bad pair
			best := ComputeLeastBadPair(pre, session, venueOcc, lecturerOcc, cohortOcc)
			placements[sessionIdx] = SessionPlacement{
				SessionIdx: sessionIdx,
				CourseIdx:  session.CourseIdx,
				VenueIdx:   best.VenueIdx,
				SlotIdx:    best.SlotIdx,
				Conflict:   true,
				Score:      best.Score,
			}

			// mark tempoary occupancy
			for d := 0; d < session.SessionDuration; d++ {
				// si means slot index
				si := best.SlotIdx + d
				if si >= totalSlots {
					break
				}
				if best.VenueIdx < len(venueOcc) && si < len(venueOcc[best.VenueIdx]) {
					venueOcc[best.VenueIdx][si] = true
				}
				if session.LecturerIdx < len(lecturerOcc) && si < len(lecturerOcc[session.LecturerIdx]) {
					lecturerOcc[session.LecturerIdx][si] = true
				}
				for _, c := range session.CohortIdxs {
					if c < len(cohortOcc) && si < len(cohortOcc[c]) {
						cohortOcc[c][si] = true
					}
				}
			}
		}
	}

	return &Candidate{
		Placements: placements,
	}
}

// calculate and return the fitness for a candidate timetable.. the lower the better the timetable
func ComputeCandidateFitness(candidate *Candidate) float64 {
	fitnessScore := 0.0

	for sessIdx := 0; sessIdx < len(candidate.Placements); sessIdx++ {
		fitnessScore += candidate.Placements[sessIdx].Score
	}

	// adds the fitness to the candidate object directly
	candidate.Fitness = 1.0 / (1 + fitnessScore)

	return fitnessScore
}

func BuildPopulation(pre *PreComputed, seed int64, populationSize int, K int) []*Candidate {
	r := rand.New(rand.NewSource(seed))
	pop := make([]*Candidate, 0, populationSize)

	for i := 0; i < populationSize; i++ {
		c := BuildOneCandidate(rand.New(rand.NewSource(r.Int63())), pre, K)
		ComputeCandidateFitness(c)
		pop = append(pop, c)
	}

	return pop
}

// basically returns a best fit candidate out of a K options
func Selection(pop []*Candidate, tournamentSize int) *Candidate {
	if len(pop) == 0 {
		return nil
	}
	best := pop[rand.Intn(len(pop))]

	for c := 0; c < tournamentSize; c++ {
		challenger := pop[rand.Intn(len(pop))]

		if challenger.Fitness > best.Fitness {
			best = challenger
		}
	}
	return best
}

// groups sessions according to their courses in ascending index order
func ComputeCourseSessions(pre *PreComputed) [][]SessionAtom {
	sessionLen := len(pre.SessionAtoms)

	sessions := make([]SessionAtom, sessionLen)

	// populates the sessions slice
	copy(sessions, pre.SessionAtoms)

	courseSessions := make([][]SessionAtom, pre.NumCourses)
	for i := range courseSessions {
		courseSessions[i] = make([]SessionAtom, 0)
	}

	// arranges sessions according to their course indexes ascendingly
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CourseIdx < sessions[j].CourseIdx
	})

	for _, session := range sessions {
		if session.CourseIdx >= 0 && session.CourseIdx < len(courseSessions) {
			courseSessions[session.CourseIdx] = append(courseSessions[session.CourseIdx], session)
		}
	}

	return courseSessions
}

// returns 2 good fit parents
func SelectParents(pop []*Candidate, tournamentSize int) (*Candidate, *Candidate) {
	parent1 := Selection(pop, tournamentSize)
	parent2 := Selection(pop, tournamentSize)
	
	// Ensure we don't return nil parents
	if parent1 == nil || parent2 == nil {
		// Fallback: return first two candidates if selection fails
		if len(pop) >= 2 {
			return pop[0], pop[1]
		} else if len(pop) == 1 {
			return pop[0], pop[0]
		} else {
			// This should not happen if population is built correctly
			return &Candidate{}, &Candidate{}
		}
	}
	return parent1, parent2
}

func DetermineBestParent(pre *PreComputed, parent1 *Candidate, parent2 *Candidate, CourseIdx int, lecturerOcc [][]bool, venueOcc [][]bool, cohortOcc [][]bool) []SessionPlacement {
	parent1ConflictScore := 0
	parent1SessionPlacements := make([]SessionPlacement, 0)

	parent2ConflictScore := 0
	parent2SessionPlacements := make([]SessionPlacement, 0)

	for _, session := range parent1.Placements {
		// Add bounds checking for session.SessionIdx
		if session.SessionIdx < 0 || session.SessionIdx >= len(pre.SessionAtoms) {
			continue
		}
		
		hasConflict := false
		lecturerIdx := pre.SessionAtoms[session.SessionIdx].LecturerIdx
		cohortIdxs := pre.SessionAtoms[session.SessionIdx].CohortIdxs
		
		if session.CourseIdx == CourseIdx {
			// Check lecturer bounds
			if lecturerIdx >= 0 && lecturerIdx < len(lecturerOcc) && session.SlotIdx < len(lecturerOcc[lecturerIdx]) {
				if lecturerOcc[lecturerIdx][session.SlotIdx] {
					parent1ConflictScore += 1500
					hasConflict = true
				}
			} else {
				// If bounds are invalid, count as conflict
				parent1ConflictScore += 1500
				hasConflict = true
			}

			// Check venue bounds
			if session.VenueIdx >= 0 && session.VenueIdx < len(venueOcc) && session.SlotIdx < len(venueOcc[session.VenueIdx]) {
				if venueOcc[session.VenueIdx][session.SlotIdx] {
					parent1ConflictScore += 1500
					hasConflict = true
				}
			} else {
				parent1ConflictScore += 1500
				hasConflict = true
			}

			// Check cohort bounds
			for _, cohortIdx := range cohortIdxs {
				if cohortIdx >= 0 && cohortIdx < len(cohortOcc) && session.SlotIdx < len(cohortOcc[cohortIdx]) {
					if cohortOcc[cohortIdx][session.SlotIdx] {
						parent1ConflictScore += 500
						hasConflict = true
					}
				} else {
					parent1ConflictScore += 500
					hasConflict = true
				}
			}
			
			session.Conflict = hasConflict
			parent1SessionPlacements = append(parent1SessionPlacements, session)
		}
	}

	// Apply the same bounds checking for parent2
	for _, session := range parent2.Placements {
		if session.SessionIdx < 0 || session.SessionIdx >= len(pre.SessionAtoms) {
			continue
		}
		
		hasConflict := false
		lecturerIdx := pre.SessionAtoms[session.SessionIdx].LecturerIdx
		cohortIdxs := pre.SessionAtoms[session.SessionIdx].CohortIdxs
		
		if session.CourseIdx == CourseIdx {
			// Check lecturer bounds for parent2
			if lecturerIdx >= 0 && lecturerIdx < len(lecturerOcc) && session.SlotIdx < len(lecturerOcc[lecturerIdx]) {
				if lecturerOcc[lecturerIdx][session.SlotIdx] {
					parent2ConflictScore += 1500
					hasConflict = true
				}
			} else {
				parent2ConflictScore += 1500
				hasConflict = true
			}

			// Check venue bounds for parent2
			if session.VenueIdx >= 0 && session.VenueIdx < len(venueOcc) && session.SlotIdx < len(venueOcc[session.VenueIdx]) {
				if venueOcc[session.VenueIdx][session.SlotIdx] {
					parent2ConflictScore += 1500
					hasConflict = true
				}
			} else {
				parent2ConflictScore += 1500
				hasConflict = true
			}

			// Check cohort bounds for parent2
			for _, cohortIdx := range cohortIdxs {
				if cohortIdx >= 0 && cohortIdx < len(cohortOcc) && session.SlotIdx < len(cohortOcc[cohortIdx]) {
					if cohortOcc[cohortIdx][session.SlotIdx] {
						parent2ConflictScore += 500
						hasConflict = true
					}
				} else {
					parent2ConflictScore += 500
					hasConflict = true
				}
			}
			
			session.Conflict = hasConflict
			parent2SessionPlacements = append(parent2SessionPlacements, session)
		}
	}

	if parent1ConflictScore < parent2ConflictScore {
		return parent1SessionPlacements
	} else {
		return parent2SessionPlacements
	}
}

func RepairChildCandidate(pre *PreComputed, childCandidate []SessionPlacement, lecOcc [][]bool, venueOcc [][]bool, cohortOcc [][]bool) {
	for idx := 0; idx < len(childCandidate); idx++ {
		if childCandidate[idx].SessionIdx < 0 || childCandidate[idx].SessionIdx >= len(pre.SessionAtoms) {
			continue
		}
		
		sessionAtom := pre.SessionAtoms[childCandidate[idx].SessionIdx]
		if childCandidate[idx].Conflict {
			leastBadPair := ComputeLeastBadPair(pre, &sessionAtom, venueOcc, lecOcc, cohortOcc)

			for d := 0; d < sessionAtom.SessionDuration; d++ {
				si := childCandidate[idx].SlotIdx + d
				// remove old occupancy with bounds checking
				if sessionAtom.LecturerIdx < len(lecOcc) && si < len(lecOcc[sessionAtom.LecturerIdx]) {
					lecOcc[sessionAtom.LecturerIdx][si] = false
				}
				if childCandidate[idx].VenueIdx < len(venueOcc) && si < len(venueOcc[childCandidate[idx].VenueIdx]) {
					venueOcc[childCandidate[idx].VenueIdx][si] = false
				}
				for _, cohortIdx := range sessionAtom.CohortIdxs {
					if cohortIdx < len(cohortOcc) && si < len(cohortOcc[cohortIdx]) {
						cohortOcc[cohortIdx][si] = false
					}
				}
			}

			childCandidate[idx].SlotIdx = leastBadPair.SlotIdx
			childCandidate[idx].VenueIdx = leastBadPair.VenueIdx
			childCandidate[idx].Score = leastBadPair.Score

			// modify the childCandidate[idx] with the least bad pair
			childCandidate[idx].Conflict = (leastBadPair.Score > 0.0)
			for d := 0; d < sessionAtom.SessionDuration; d++ {
				si := childCandidate[idx].SlotIdx + d
				if si >= pre.TotalSlots {
					break
				}
				// update occupancy with bounds checking
				if sessionAtom.LecturerIdx < len(lecOcc) && si < len(lecOcc[sessionAtom.LecturerIdx]) {
					lecOcc[sessionAtom.LecturerIdx][si] = true
				}
				if childCandidate[idx].VenueIdx < len(venueOcc) && si < len(venueOcc[childCandidate[idx].VenueIdx]) {
					venueOcc[childCandidate[idx].VenueIdx][si] = true
				}
				for _, cohortIdx := range sessionAtom.CohortIdxs {
					if cohortIdx < len(cohortOcc) && si < len(cohortOcc[cohortIdx]) {
						cohortOcc[cohortIdx][si] = true
					}
				}
			}
		}
	}
}

func Crossover(pre *PreComputed, parent1 *Candidate, parent2 *Candidate, CourseSessions [][]SessionAtom) (*Candidate, [][]bool, [][]bool, [][]bool) {
	totalSlots := pre.TotalSlots
	childCandidate := make([]SessionPlacement, 0, totalSlots)

	venueOccupied := make([][]bool, pre.NumVenues)
	for i := range venueOccupied {
		venueOccupied[i] = make([]bool, totalSlots)
	}

	lecturerOccupied := make([][]bool, pre.NumLecturers)
	for i := range lecturerOccupied {
		lecturerOccupied[i] = make([]bool, totalSlots)
	}

	cohortOccupied := make([][]bool, pre.NumCohorts)
	for i := range cohortOccupied {
		cohortOccupied[i] = make([]bool, totalSlots)
	}

	for courseIdx := range CourseSessions {
		placements := DetermineBestParent(pre, parent1, parent2, courseIdx, lecturerOccupied, venueOccupied, cohortOccupied)
		for _, placement := range placements {
			childCandidate = append(childCandidate, placement)
			
			// Add bounds checking for session index
			if placement.SessionIdx < 0 || placement.SessionIdx >= len(pre.SessionAtoms) {
				continue
			}
			
			lecturerIdx := pre.SessionAtoms[placement.SessionIdx].LecturerIdx
			cohortIdxs := pre.SessionAtoms[placement.SessionIdx].CohortIdxs

			sessionAtom := pre.SessionAtoms[placement.SessionIdx]
			for d := 0; d < sessionAtom.SessionDuration; d++ {
				si := d + placement.SlotIdx
				if si >= pre.TotalSlots {
					break
				}
				// mark occupancy for lecturer venue and cohorts with bounds checking
				if lecturerIdx < len(lecturerOccupied) && si < len(lecturerOccupied[lecturerIdx]) {
					lecturerOccupied[lecturerIdx][si] = true
				}
				if placement.VenueIdx < len(venueOccupied) && si < len(venueOccupied[placement.VenueIdx]) {
					venueOccupied[placement.VenueIdx][si] = true
				}
				for _, cohortidx := range cohortIdxs {
					if cohortidx < len(cohortOccupied) && si < len(cohortOccupied[cohortidx]) {
						cohortOccupied[cohortidx][si] = true
					}
				}
			}
		}
	}

	// then i would run a repair function on the child and return the child timetable
	RepairChildCandidate(pre, childCandidate, lecturerOccupied, venueOccupied, cohortOccupied)

	// calculate fitness of child candidate
	actualChildCandidate := &Candidate{
		Placements: childCandidate,
	}
	return actualChildCandidate, lecturerOccupied, venueOccupied, cohortOccupied
}

func Mutation(pre *PreComputed, childCandidate *Candidate, r *rand.Rand, lecOcc [][]bool, venueOcc [][]bool, cohortOcc [][]bool) {
	mutationRate := 0.05

	for placementIdx, placement := range childCandidate.Placements {
		if r.Float64() < mutationRate {
			// Add bounds checking
			if placementIdx < 0 || placementIdx >= len(childCandidate.Placements) {
				continue
			}
			if childCandidate.Placements[placementIdx].SessionIdx < 0 || childCandidate.Placements[placementIdx].SessionIdx >= len(pre.SessionAtoms) {
				continue
			}
			
			sessionAtom := pre.SessionAtoms[childCandidate.Placements[placementIdx].SessionIdx]

			// remove occupancy
			for d := 0; d < sessionAtom.SessionDuration; d++ {
				si := placement.SlotIdx + d
				if si >= pre.TotalSlots {
					break
				}
				if sessionAtom.LecturerIdx < len(lecOcc) && si < len(lecOcc[sessionAtom.LecturerIdx]) {
					lecOcc[sessionAtom.LecturerIdx][si] = false
				}
				if placement.VenueIdx < len(venueOcc) && si < len(venueOcc[placement.VenueIdx]) {
					venueOcc[placement.VenueIdx][si] = false
				}
				for _, cohortidx := range sessionAtom.CohortIdxs {
					if cohortidx < len(cohortOcc) && si < len(cohortOcc[cohortidx]) {
						cohortOcc[cohortidx][si] = false
					}
				}
			}

			//  should run and add some variety for 30% of the chosen placements
			if r.Float64() > 0.3 {
				// find a random new slot - FIXED: Add bounds checking for random slot
				maxSlot := pre.TotalSlots - sessionAtom.SessionDuration
				if maxSlot <= 0 {
					maxSlot = 1
				}
				randomSlot := r.Intn(maxSlot)
				
				if len(sessionAtom.AllowedVenuesIdx) > 0 {
					randomVenue := sessionAtom.AllowedVenuesIdx[r.Intn(len(sessionAtom.AllowedVenuesIdx))]
					childCandidate.Placements[placementIdx].SlotIdx = randomSlot
					childCandidate.Placements[placementIdx].VenueIdx = randomVenue
					childCandidate.Placements[placementIdx].Conflict = true
				}
			} else {
				// for the other 70%
				leastBadPair := ComputeLeastBadPair(pre, &sessionAtom, venueOcc, lecOcc, cohortOcc)
				childCandidate.Placements[placementIdx].SlotIdx = leastBadPair.SlotIdx
				childCandidate.Placements[placementIdx].VenueIdx = leastBadPair.VenueIdx
				childCandidate.Placements[placementIdx].Score = leastBadPair.Score
				childCandidate.Placements[placementIdx].Conflict = (leastBadPair.Score > 0.0)
			}

			// reoccupy occupancy
			for d := 0; d < sessionAtom.SessionDuration; d++ {
				si := childCandidate.Placements[placementIdx].SlotIdx + d
				if si >= pre.TotalSlots {
					break
				}
				if sessionAtom.LecturerIdx < len(lecOcc) && si < len(lecOcc[sessionAtom.LecturerIdx]) {
					lecOcc[sessionAtom.LecturerIdx][si] = true
				}
				if childCandidate.Placements[placementIdx].VenueIdx < len(venueOcc) && si < len(venueOcc[childCandidate.Placements[placementIdx].VenueIdx]) {
					venueOcc[childCandidate.Placements[placementIdx].VenueIdx][si] = true
				}
				for _, cohortidx := range sessionAtom.CohortIdxs {
					if cohortidx < len(cohortOcc) && si < len(cohortOcc[cohortidx]) {
						cohortOcc[cohortidx][si] = true
					}
				}
			}
		}
	}
}

func SortPopulation(pop []*Candidate) {
	sort.Slice(pop, func(i int, j int) bool {
		return pop[i].Fitness > pop[j].Fitness
	})
}

func BuildNextGeneration(pre *PreComputed,
	previousPopulation []*Candidate,
	courseSessions [][]SessionAtom,
	r *rand.Rand,
) []*Candidate {
	if len(previousPopulation) == 0 {
		return previousPopulation
	}
	
	topK := int(0.1 * float64(len(previousPopulation)))
	if topK < 1 {
		topK = 1
	}

	newPopulation := make([]*Candidate, 0)

	SortPopulation(previousPopulation)

	for i := 0; i < len(previousPopulation); i++ {
		// choose the top k in previous population
		if i < topK {
			newPopulation = append(newPopulation, previousPopulation[i])
		} else {
			parent1, parent2 := SelectParents(previousPopulation, topK)
			candidate, lecOcc, venueOcc, cohortOcc := Crossover(pre, parent1, parent2, courseSessions)
			Mutation(pre, candidate, r, lecOcc, venueOcc, cohortOcc)
			ComputeCandidateFitness(candidate)
			newPopulation = append(newPopulation, candidate)
		}
	}
	return newPopulation
}

func SelectBestCandidateFromPopulation(pop []*Candidate) *Candidate {
	if len(pop) == 0 {
		return nil
	}
	
	bestIdx := 0
	currentFitness := math.Inf(-1)
	for idx, cand := range pop {
		if cand.Fitness > currentFitness {
			currentFitness = cand.Fitness
			bestIdx = idx
		}
	}
	slog.Info("the best population","val",pop[bestIdx])
	return pop[bestIdx]
}

func (c *Computed) GeneticAlgorithm(pre *PreComputed) *Candidate {
	seed := int64(50)
	populationSize := 100
	k := 5
	numberOfGeneration := 100
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	population := BuildPopulation(pre, seed, populationSize, k)
	courseSessions := ComputeCourseSessions(pre)

	for i := 0; i < numberOfGeneration; i++ {
		population = BuildNextGeneration(pre, population, courseSessions, r)
	}
	return SelectBestCandidateFromPopulation(population)
}