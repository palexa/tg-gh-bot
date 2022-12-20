package person

type Storage interface {
	Create(dto *CreatePersonDto) (*Person, error)
	Update(dto *UpdatePersonDto) (*Person, error)
	FindOrCreate(dto *CreatePersonDto) (*Person, error)
	FindByTelegramId(telegramId int) (*Person, error)
	FindById(id int) (*Person, error)
	GetAll() ([]*Person, error)
	UpdateGHToken(id int, token string) error
}
