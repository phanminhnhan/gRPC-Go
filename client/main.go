package main

import (
	router "grpc-go/client/router"
	"github.com/labstack/echo/v4"
	api "grpc-go/client/api"
)



func main() {
	e := echo.New()
	image := api.NewImageClientService()
	r := router.Router{
		Echo: e,
		Image: image,
	}
	r.SetupRouter()
	e.Logger.Fatal(e.Start(":8888"))
}