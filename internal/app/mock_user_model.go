package app

import "errors"

type MockUserRepo struct {
    Users map[string]User
}


func (m *MockUserRepo) SaveNewUser(user User) error {
    if _, exists := m.Users[user.Login]; exists {
        return errors.New("user already exists")
    }
    m.Users[user.Login] = user
    return nil
}
func (m *MockUserRepo) FindByLogin(login string) (User, error) {
    user, ok := m.Users[login]
    if !ok {
        return User{}, errors.New("not found")
    }
    return user, nil
}
func (m *MockUserRepo) FindByUUID(uuid string) (User, error) {
    for _, user := range m.Users {
        if user.UUID.String() == uuid {
            return user, nil
        }
    }
    return User{}, errors.New("not found")
}