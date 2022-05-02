package repositories

import (
	"chatRoom/domain"
	"chatRoom/services"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (u *UserRepository) Insert(user domain.User) (string, error) {

	id := uuid.New().String()

	user.ID = id

	err := u.Db.Create(&user).Error

	return id, err
}

func (u *UserRepository) GetAll() []domain.User {
	var users []domain.User
	u.Db.Find(&users)

	return users
}

func (u *UserRepository) SignIn(userInfo services.UserSignIn) (string, error) {
	user := domain.User{
		Name: userInfo.Name,
	}
	u.Db.Find(&user)

	if user.Password != userInfo.Password {
		return "", errors.New("Password does not correct")
	}

	return user.ID, nil
}

func (u *UserRepository) Delete(id string) error {
	_, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	u.Db.Delete(&domain.User{ID: id})

	u.Db.Delete(&domain.User{ID: id})

	return nil
}
