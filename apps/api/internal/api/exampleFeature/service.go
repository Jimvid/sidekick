package climbs

type ClimbService struct {
	storage *ClimbStorage
}

func NewClimbService(storage *ClimbStorage) *ClimbService {
	return &ClimbService{
		storage: storage,
	}
}

func (s *ClimbService) Create(climb ClimbModel) (ClimbModel, error) {
	climb, err := s.storage.Create(climb)
	if err != nil {
		return ClimbModel{}, err
	}
	return climb, nil
}

func (s *ClimbService) GetAll(userId string) ([]ClimbModel, error) {
	climbs, err := s.storage.GetAll(userId)
	if err != nil {
		return climbs, err
	}
	return climbs, nil
}

func (s *ClimbService) FindById(userId, climbId string) (ClimbModel, error) {
	climb, err := s.storage.FindById(userId, climbId)
	if err != nil {
		return ClimbModel{}, err
	}
	return climb, nil
}

func (s *ClimbService) Delete(userId, climbId string) error {
	err := s.storage.Delete(userId, climbId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ClimbService) Update(userId, climbId string, data ClimbUpdateReq) (ClimbModel, error) {
	climb, err := s.storage.Update(userId, climbId, data)
	if err != nil {
		return climb, err
	}
	return climb, nil
}
