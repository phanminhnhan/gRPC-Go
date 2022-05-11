package handle

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	image "grpc-go/pb"
	cdnService "grpc-go/server/cloudinary"
	"io"
	"log"
	"mime/multipart"
	"os"
	"bytes"
	// "io/ioutil"
)

type Server struct {
	image.UnimplementedImageServiceServer
}

const maxImageSize = 1 << 20

func (s *Server) UploadImage(stream image.ImageService_UploadImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	id := uuid.New()
	imageType := req.GetInfo().GetImageType()
	imageName := req.GetInfo().GetImageName()
	imageData := &bytes.Buffer{}
	imageSize := 0
	fileName := id.String()+imageName

	// open output file
	fo, err := os.Create(fileName)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "failed to create file"))
	}
	defer os.Remove(fileName)
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			 panic(err)
		}
	}()


	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			res := &image.UploadImageResponse{
				Info: &image.Image{
					ImageId:   id.String(),
					ImageName: imageName,
					ImageType: imageType,
					ImageSize: uint32(imageSize),
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				},
			}
			err = stream.SendAndClose(res)
			if err != nil {
				return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
			}
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
		if _, err := imageData.Write(chunk); err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
		// write date to file 
		if _, err := fo.Write(chunk); err != nil {
			return err
		}
	}


	// upFileInfo, _ := fo.Stat()
    // var fileSize int64 = upFileInfo.Size()
    // fileBuffer := make([]byte, fileSize)
    // fo.Read(fileBuffer)

	// fileHeader, err := getFileHeader(fo)
	// if err != nil {
	// 	fmt.Println("err at return file header")
	// }
	// content, err := fileHeader.Open() 
	// if err != nil {
	// 	logError(err)
	// }
	url, err := cdnService.NewMediaUpload().FileUpload(fileName)
	if err != nil {
		fmt.Println("error happened : %w", err)
	}
	fmt.Println("URL:", url)
	
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

func getFileHeader(file *os.File) (*multipart.FileHeader, error) {
	// get file size
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// create *multipart.FileHeader
	return &multipart.FileHeader{
		Filename: fileStat.Name(),
		Size:     fileStat.Size(),
	}, nil
}
