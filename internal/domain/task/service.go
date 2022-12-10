package task

type Service interface {
	GetAll() ([]*Task, error)
	Create(dto *CreateTaskDto) (*Task, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetAll() ([]*Task, error) {
	return s.storage.GetAll()
}

func (s *service) Create(dto *CreateTaskDto) (*Task, error) {
	return s.storage.Create(dto)
}
