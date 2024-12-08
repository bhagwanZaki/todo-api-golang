package common

import (
	"context"
	"errors"
	"log"
	"todoGoApi/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageWrappper struct {
	client *cloudinary.Cloudinary
}

func InitiailizeImageWrappper() ( *ImageWrappper, error ){
	log.Println("CLOUDINARY_URL ", utils.CLOUDINARY_URL)
	cld, err := cloudinary.NewFromURL(utils.CLOUDINARY_URL)

	if err != nil {
		Logger(err.Error(),"Image Wrapper Initializer")
		return &ImageWrappper{}, err
	}

	cld.Config.URL.Secure = true
	return &ImageWrappper{client: cld}, nil
}

func (upld *ImageWrappper) Upload(ctx context.Context, file interface{}, fileName string) (*uploader.UploadResult, error) {

	param := uploader.UploadParams{
        PublicID:       fileName,
		Folder: 		"media/feedback",
        UniqueFilename: api.Bool(true),
        Overwrite:      api.Bool(true)};
		
	uploadResult, err := upld.client.Upload.Upload(ctx, file, param)

	if err != nil {
		Logger("Failed to upload image" + err.Error(), "Image Uploader")
		return &uploader.UploadResult{}, err
	}
	
	if uploadResult.Error.Message !=  ""{
		Logger("Failed to upload image" + uploadResult.Error.Message, "Image Uploader")
		return &uploader.UploadResult{}, errors.New(uploadResult.Error.Message)
	}

	return uploadResult, nil
}
