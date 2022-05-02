package controller

import (
	"chatRoom/handlers"
	"chatRoom/middlewares"
	"chatRoom/routes"
	"chatRoom/services"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
)

func init() {
	api := routes.AppRouter.V1Group.Party("/chat")

	api.RegisterDependency(func(ctx iris.Context) services.IWebsocketService {

		return services.NewWebsocketService(websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		})
	})

	api.ConfigureContainer(chatRoute)
}

func chatRoute(api *iris.APIContainer) {

	api.Get("/start", middlewares.Auth, handlers.StartRoom)
	api.Get("/join/{roomId}", middlewares.Auth, handlers.JoinRoom)
}
