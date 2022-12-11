package person

type CreatePersonDto struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	UserName   string `json:"user_name,omitempty"`
	TelegramId int64  `json:"telegram_id,omitempty"`
}

type UpdatePersonDto struct {
	ID          string `json:"_id,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	TelegramId  int64  `json:"telegram_id,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
