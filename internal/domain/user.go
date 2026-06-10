package domain

type User struct {
	ID        string  `json:"id" db:"id"`
	OAuthId   string  `json:"-" db:"oauth_id"`
	Name      string  `json:"name" db:"name"`
	Email     string  `json:"email" db:"email"`
	Picture   string  `json:"picture" db:"picture"`
	CustomUrl *string `json:"custom_url" db:"custom_url"`
}
