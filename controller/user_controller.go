package controller

import (
	"chatRoom/databases"
	"chatRoom/handlers"
	"chatRoom/middlewares"
	"chatRoom/repositories"
	"chatRoom/routes"
	"chatRoom/services"
	"github.com/kataras/iris/v12"
)

func init() {
	api := routes.AppRouter.V1Group.Party("/user")

	api.RegisterDependency(func(ctx iris.Context) services.UserService {

		return repositories.NewUserRepository(&databases.MysqlOrm)
	})

	api.ConfigureContainer(userRoute)
}

func userRoute(api *iris.APIContainer) {
	api.Get("/start", middlewares.Auth, handlers.GetAll)
	api.Post("/signIn", handlers.SignIn)
	api.Post("/signup", handlers.InsertUser)
	api.Delete("/", middlewares.Auth, handlers.DeleteUser)
}
