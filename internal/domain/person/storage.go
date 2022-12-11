package person

type Storage interface {
	Create(dto CreatePersonDto) (*Person, error)
}
