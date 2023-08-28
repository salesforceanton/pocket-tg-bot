package repository

type Repository interface {
	SaveAccessToken(chatID int64, token string) error
	SaveRequestToken(chatID int64, token string) error
	GetAccessToken(chatID int64) (string, error)
	GetRequestToken(chatID int64) (string, error)
}
