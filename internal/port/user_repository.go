package port

import (
	"github.com/disdreamq/BlogApi/internal/domain"
)

type UserReaderByID interface {
	GetUserByID(userID int64) (*domain.User, error)
}
type UserReaderByEmail interface {
	GetUserByEmail(email string) (*domain.User, error)
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
	UserReaderByID
	UserReaderByEmail
	UserCreater
	UserUpdater
	UserDeleter
}
