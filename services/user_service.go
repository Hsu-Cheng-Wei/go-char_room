package services

import "chatRoom/domain"

type UserSignIn struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserService interface {
	GetAll() []domain.User
	SignIn(userInfo UserSignIn) (string, error)
	Insert(user domain.User) (string, error)
	Delete(id string) error
}
