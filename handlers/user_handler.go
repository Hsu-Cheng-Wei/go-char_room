package handlers

import (
	"bytes"
	"chatRoom/domain"
	"chatRoom/models"
	"chatRoom/routes"
	"chatRoom/services"
	"encoding/gob"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

type UserDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetAll(ctx *context.Context, userService services.UserService) []UserDto {
	users := userService.GetAll()
	var result []UserDto
	for _, user := range users {
		result = append(result, UserDto{
			ID:    "id",
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return result
}

func SignIn(ctx iris.Context, userAuth services.UserSignIn, userService services.UserService) {
	id, err := userService.SignIn(userAuth)
	if err != nil {
		ctx.JSON(iris.Map{
			"status":  201,
			"message": "Sign in fail",
		})
		return
	}
	var data bytes.Buffer
	encode := gob.NewEncoder(&data)
	err = encode.Encode(models.UserClaim{
		ID:   id,
		Name: userAuth.Name,
	})
	if err != nil {
		panic(err)
	}

	token := routes.AppRouter.GetJwtSignedToken(models.UserClaim{
		ID:   id,
		Name: userAuth.Name,
	})

	token = "bearer " + token

	ctx.SetCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
		Path:  "/",
	})
	ctx.JSON(iris.Map{
		"status":  200,
		"message": "success",
		"token":   token,
	})
}

func InsertUser(userService services.UserService, user domain.User) string {
	id, err := userService.Insert(user)

	if err != nil {
		panic(err)
	}

	return id
}

func DeleteUser(userService services.UserService, id string) error {
	return userService.Delete(id)
}
