package habits

import (
	"time"

	"github.com/google/uuid"
)

type HabitService struct {
	storage *HabitStorage
}

func NewHabitService(storage *HabitStorage) *HabitService {
	return &HabitService{
		storage: storage,
	}
}

func (s *HabitService) CreateHabit(userId string, req HabitReq) (HabitModel, error) {
	habit := HabitModel{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	return habit, s.storage.CreateHabit(userId, habit)
}

func (s *HabitService) GetAllHabits(userId string) ([]HabitModel, error) {
	return s.storage.GetAllHabits(userId)
}

func (s *HabitService) FindHabitById(userId, habitId string) (HabitModel, error) {
	return s.storage.FindHabitById(userId, habitId)
}

func (s *HabitService) DeleteHabit(userId, habitId string) error {
	return s.storage.DeleteHabit(userId, habitId)
}

func (s *HabitService) UpdateHabit(userId, habitId string, req HabitReq) (HabitModel, error) {
	existing, err := s.storage.FindHabitById(userId, habitId)
	if err != nil {
		return HabitModel{}, err
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Color = req.Color
	existing.UpdatedAt = time.Now().Unix()

	err = s.storage.UpdateHabit(userId, habitId, existing)
	if err != nil {
		return HabitModel{}, err
	}

	return existing, nil
}

func (s *HabitService) CreateHabitLog(userId string, req HabitLogReq) (HabitLogModel, error) {
	log := HabitLogModel{
		ID:        uuid.New().String(),
		HabitId:   req.HabitId,
		Date:      req.Date,
		Note:      req.Note,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return log, s.storage.CreateHabitLog(userId, log)
}

func (s *HabitService) GetAllHabitLogs(userId string) ([]HabitLogModel, error) {
	return s.storage.GetAllHabitLogs(userId)
}

func (s *HabitService) FindHabitLogById(userId, logId string) (HabitLogModel, error) {
	return s.storage.FindHabitLogById(userId, logId)
}

func (s *HabitService) DeleteHabitLog(userId, logId string) error {
	return s.storage.DeleteHabitLog(userId, logId)
}

func (s *HabitService) UpdateHabitLog(userId, logId string, req HabitLogReq) (HabitLogModel, error) {
	existing, err := s.storage.FindHabitLogById(userId, logId)
	if err != nil {
		return HabitLogModel{}, err
	}

	existing.HabitId = req.HabitId
	existing.Date = req.Date
	existing.Note = req.Note
	existing.UpdatedAt = time.Now().Unix()

	err = s.storage.UpdateHabitLog(userId, logId, existing)
	if err != nil {
		return HabitLogModel{}, err
	}

	return existing, nil
}
