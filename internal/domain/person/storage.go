package person

type Storage interface {
	Create(dto *CreatePersonDto) (*Person, error)
	Update(dto *UpdatePersonDto) (*Person, error)
	FindOrCreate(dto *CreatePersonDto) (*Person, error)
	FindByTelegramId(telegramId int) (*Person, error)
	GetAll() ([]*Person, error)
}
