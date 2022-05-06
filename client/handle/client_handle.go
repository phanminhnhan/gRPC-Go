package handle


import(
	"grpc-go/pb"
	"google.golang.org/grpc"
	"context"
	"mime/multipart"
	"log"
	"time"
	"path/filepath"
	"io"
	"bufio"
	"fmt"
)

type imageClient struct {
	service pb.ImageServiceClient
}

func NewImageClient(cc *grpc.ClientConn) *imageClient {
	service := pb.NewImageServiceClient(cc)
	return &imageClient{service}
}



func (c *imageClient) UploadImage(fileIput *multipart.FileHeader)*pb.UploadImageResponse{
	file, err := fileIput.Open()
	// if err != nil {
	// 	log.Fatal("cannot open image file: ", err)
	// }
	// defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := c.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.Image{
				ImageName: fileIput.Filename,
				ImageType: filepath.Ext(fileIput.Filename),
			},
		},
	}
	fmt.Println(req)
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.Info.GetImageId(), res.Info.GetImageSize())
	return res
}