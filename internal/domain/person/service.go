package person

type Service interface {
	Create(dto *CreatePersonDto) (*Person, error)
	Update(dto *UpdatePersonDto) (*Person, error)
	FindOrCreate(dto *CreatePersonDto) (*Person, error)
	SetGHToken(dto *UpdatePersonDto) error
	GetAll() ([]*Person, error)
	GetGHData() error
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) Create(dto *CreatePersonDto) (*Person, error) {
	return s.storage.Create(dto)
}

func (s *service) Update(dto *UpdatePersonDto) (*Person, error) {
	panic("implement me")
}

func (s *service) SetGHToken(dto *UpdatePersonDto) error {
	err := s.storage.UpdateGHToken(dto.ID, dto.AccessToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetGHData() error {
	panic("implement me")
}

func (s *service) GetAll() ([]*Person, error) {
	return s.storage.GetAll()
}

func (s *service) FindOrCreate(dto *CreatePersonDto) (*Person, error) {
	return s.storage.FindOrCreate(dto)
}
