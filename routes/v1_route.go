package routes

import (
	"chatRoom/enviroment"
	"fmt"
	jwtMiddleWare "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	jwt "github.com/kataras/iris/v12/middleware/jwt"
	"strings"
	"time"
)

var AppRouter Router

type Router struct {
	App *iris.Application

	V1Group router.Party

	JwtMiddleware *jwtMiddleWare.Middleware

	JwtSigned string

	JwtSigner *jwt.Signer
}

func (r Router) GetJwtSignedToken(key interface{}) string {

	token, error := r.JwtSigner.Sign(key)

	if error != nil {
		panic(error)
	}

	return string(token[:])
}

func init() {
	v1Router()
}

func v1Router() {
	AppRouter.App = iris.New()

	AppRouter.JwtSigner = jwt.NewSigner(jwt.HS256, enviroment.JwtSigned, 24*time.Hour)
	AppRouter.V1Group = AppRouter.App.Party(enviroment.ApiV1)
	AppRouter.JwtSigned = enviroment.JwtSigned
	AppRouter.JwtMiddleware = jwtConfig()
}

func jwtConfig() *jwtMiddleWare.Middleware {
	return jwtMiddleWare.New(jwtMiddleWare.Config{

		Extractor: FromAuthHeaderAndCookies,

		ValidationKeyGetter: func(token *jwtMiddleWare.Token) (interface{}, error) {
			return []byte(enviroment.JwtSigned), nil
		},
		ContextKey:    "iris.jwt.claims",
		SigningMethod: jwtMiddleWare.SigningMethodHS256,
	})
}

func FromAuthHeaderAndCookies(ctx iris.Context) (string, error) {
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		auth = ctx.GetCookie("Authorization")
		if auth == "" {
			return "", nil // No error, just no token
		}
	}

	authHeaderParts := strings.Split(auth, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
