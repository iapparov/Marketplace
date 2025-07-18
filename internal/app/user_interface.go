package app

type UserRepository interface {
	SaveNewUser(user User) error
	FindByLogin(login string) (User, error) // strings.ToLower(req.Login) допилить
	FindByUUID(uuid string) (User, error)
}