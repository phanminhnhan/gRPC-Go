package cloudinary  

// import (
// 	"mime/multipart"
// )


type mediaUpload interface {
    FileUpload(fileName string) (string, error)
    RemoteUpload(url string) (string, error)
}

type media struct {}

func NewMediaUpload() mediaUpload {
    return &media{}
}

func (*media) FileUpload(fileName string) (string, error) {
    //upload
    uploadUrl, err := ImageUploadHelper(fileName)
    if err != nil {
        return "", err
    }
    return uploadUrl, nil
}

func (*media) RemoteUpload(url string) (string, error) {
    //upload
    uploadUrl, errUrl := ImageUploadHelper(url)
    if errUrl != nil {
        return "", errUrl
    }
    return uploadUrl, nil
}