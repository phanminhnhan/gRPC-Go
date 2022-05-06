package handle


import (
	image "grpc-go/pb"
	"io"
	"log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"bytes"
	"context"
	"github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	image.UnimplementedImageServiceServer
}

const maxImageSize = 1 << 20
func (s *Server)UploadImage(stream image.ImageService_UploadImageServer) error{
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}
	id := uuid.New()
	imageType:= req.GetInfo().GetImageType()
	imageName:= req.GetInfo().GetImageName()
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}
		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}
		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}
	
	res := &image.UploadImageResponse{
		Info: &image.Image{
			ImageId: id.String(),
			ImageName: imageName,
			ImageType: imageType,
			ImageSize: uint32(imageSize),
			CreatedAt:timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	return nil
}


func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
