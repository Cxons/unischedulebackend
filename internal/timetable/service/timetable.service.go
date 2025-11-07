package service

import (
	"context"
	"log/slog"
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



func NewTimetableService(repo timetableRepository)*timeTableService{
	return &timeTableService{
		repo: repo,
	}
}



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

		startTime := startOfDay.Add(time.Duration(slotIdxInDay) * slotDuration)

		slotMap[i] = SlotInfo{
			Day:       days[dayIdx],
			StartTime: startTime,
		}
	}

	return slotMap
}

func (tts *timeTableService) CreateATimeTable(ctx context.Context,startOfDay time.Time, endOfDay time.Time,uniId uuid.UUID)(timeTableResponse,string,error){
	slotDuration := time.Hour
	totalDuration := endOfDay.Sub(startOfDay)
	slotsPerDay := int(totalDuration/slotDuration)
	days := []string{"Monday","Tuesday","Wednesday","Thursday","Friday"}
	slotMap := BuildSlotMap(slotsPerDay,days,startOfDay,slotDuration)



  precomputed,_,venueMap,_,coursesMap := tts.computed.ComputePreComputed(ctx,uniId,slotsPerDay,startOfDay,days,slotDuration)
  candidateTimetable := tts.computed.GeneticAlgorithm(precomputed)
  candidateData := sqlc.CreateCandidateParams{
	Fitness: candidateTimetable.Fitness,
	UniversityID: uniId,
	CandidateStatus: "CURRENT",
  }
  sessionPlacements := make([]customSessionPlacement,0)
  for _,val := range candidateTimetable.Placements{
	courseId := uuid.UUID{}
	venueId := uuid.UUID{}
	for key,value := range coursesMap{
		if val.CourseIdx == value{
			courseId = key
		}
	} 
	for key,value := range venueMap{
		if val.VenueIdx == value{
			venueId = key
		}
	}

	sessionPlacements = append(sessionPlacements,customSessionPlacement{
		SessionIdx: int32(val.SessionIdx),
		CourseId: courseId,
		VenueId: venueId,
		Day: slotMap[val.SlotIdx].Day,
		SessionTime: slotMap[val.SlotIdx].StartTime,
		UniversityId: uniId,
	})
  }

  //changes the current status of the latest candidate timetable to deprecated
  deprecateErr := tts.repo.DeprecateLatestCandidate(ctx,uniId)
  if deprecateErr != nil{
	tts.logger.Error("error deprecating latest candidate","err:",deprecateErr)
	return timeTableResponse{},status.InternalServerError.Message,deprecateErr
  }

  err := tts.repo.CreateACandidateTimeTable(ctx,candidateData,sessionPlacements)
  if err != nil{
	restoreErr := tts.repo.RestoreCurrentCandidate(ctx,uniId)
	if restoreErr != nil{
		tts.logger.Error("error restoring current candidate","err:",restoreErr)
	}
	tts.logger.Error("error creating the candidate timetable","err:",err)
	return timeTableResponse{},status.InternalServerError.Message,err
  }

  return timeTableResponse{
	Message: "Timetable created successfully",
	StatusCode: status.Created.Code,
	StatusCodeMessage: status.Created.Message,
  },status.Created.Message,nil


   
}












