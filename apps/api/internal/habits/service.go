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

func (s *HabitService) Create(userId string, req HabitReq) (HabitModel, error) {
	habit := HabitModel{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	return habit, s.storage.Create(userId, habit)
}

func (s *HabitService) GetAll(userId string) ([]HabitModel, error) {
	return s.storage.GetAll(userId)
}

func (s *HabitService) FindById(userId, habitId string) (HabitModel, error) {
	return s.storage.FindById(userId, habitId)
}

func (s *HabitService) Delete(userId, habitId string) error {
	return s.storage.Delete(userId, habitId)
}

func (s *HabitService) Update(userId, habitId string, req HabitReq) (HabitModel, error) {
	existing, err := s.storage.FindById(userId, habitId)
	if err != nil {
		return HabitModel{}, err
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Color = req.Color
	existing.UpdatedAt = time.Now().Unix()

	err = s.storage.Update(userId, habitId, existing)
	if err != nil {
		return HabitModel{}, err
	}

	return existing, nil
}
