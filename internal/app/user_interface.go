package app

type UserRepository interface {
	Save(user User) error
	FindByLogin(login string) (User, error)
	FindByUUID(uuid string) (bool, []string)
}