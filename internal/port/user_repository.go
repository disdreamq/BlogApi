package port

import (
	"github.com/disdreamq/BlogApi/internal/domain"
)

type UserReader interface {
	ReadUser(userID int64) (*domain.User, error)
}

type UserCreater interface {
	CreateUser(username, email, password string) (int64, error)
}

type UserUpdater interface {
	UpdateUser(user *domain.User) error
}
type UserDeleter interface {
	DeleteUser(id int64) error
}

type Hasher interface {
	Hash(password string) (string, error)
	Check(hashed, plain string) error
}

type UserRepository interface {
	UserReader
	UserCreater
	UserUpdater
	UserDeleter
}
