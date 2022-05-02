package main

import (
	_ "chatRoom/controller"
	"chatRoom/routes"
	_ "chatRoom/services"
)

func main() {
	app := routes.AppRouter.App

	_ = app.Listen(":8080")
}
