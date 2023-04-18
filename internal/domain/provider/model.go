package provider

type Repository struct {
	Name    string
	Stars   int
	Private bool
}

type PullRequest struct {
	Id int
}
