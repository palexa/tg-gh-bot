package task

type Storage interface {
	GetAll() ([]*Task, error)
	Create(dto *CreateTaskDto) (*Task, error)
}
