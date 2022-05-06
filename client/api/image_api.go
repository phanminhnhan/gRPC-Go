package api

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"time"
	handle "grpc-go/client/handle"
	model "grpc-go/client/model"
	"context"
	"net/http"
)

type ImageClientService struct {}

func NewImageClientService()*ImageClientService{
	return &ImageClientService{}
}

func (*ImageClientService) UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, ":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	defer conn.Close()
	service := handle.NewImageClient(conn)

	res := service.UploadImage(file)

	return c.JSON(http.StatusOK, model.ResponseData{
		Message: "done uploading",
		StattusCode: 200,
		Data: res,
	})
	
}
