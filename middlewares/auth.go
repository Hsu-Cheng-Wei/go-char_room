package middlewares

import (
	"chatRoom/models"
	"chatRoom/routes"
	"encoding/json"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
)

func Auth(ctx iris.Context) {
	err := routes.AppRouter.JwtMiddleware.CheckJWT(ctx)

	if err != nil {
		_, _ = ctx.WriteString(err.Error())
		return
	}

	token := ctx.Values().Get(routes.AppRouter.JwtMiddleware.Config.ContextKey).(*jwt.Token)

	data, _ := json.Marshal(token.Claims)

	ctx.Values().Set("user", string(data))
	ctx.Next()
}

func GetUserClaim(ctx iris.Context) (models.UserClaim, error) {
	var user models.UserClaim
	value := ctx.Values().Get("user")
	if value == nil {
		return models.UserClaim{}, errors.New("Can't find user claim")
	}

	_ = json.Unmarshal([]byte(value.(string)), &user)

	return user, nil
}
