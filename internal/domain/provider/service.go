package provider

type Service interface {
	GetUserData() (string, error)
	GetRepositories() ([]Repository, error)
	PullRequests()
}
