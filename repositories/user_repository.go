package repositories

import (
	"chatRoom/domain"
	"chatRoom/models"
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
	ch := make(chan models.OrmFindResult)

	go func() {
		var user domain.User
		err := u.Db.Where("name=?", userInfo.Name).Find(&user).Error
		ch <- models.OrmFindResult{
			Instance: user,
			Error:    err,
		}
	}()

	result := <-ch

	if result.Error != nil {
		return "", result.Error
	}

	user := result.Instance.(domain.User)
	if user.Password != userInfo.Password {
		return "", errors.New("password is not equal")
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
