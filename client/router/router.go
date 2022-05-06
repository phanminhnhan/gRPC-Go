package router

import (
	"github.com/labstack/echo/v4"
	api "grpc-go/client/api"
)

type Router struct {
	Echo *echo.Echo
	Image *api.ImageClientService
}

func (r *Router)SetupRouter(){
	r.Echo.POST("api/image", r.Image.UploadImage)
}