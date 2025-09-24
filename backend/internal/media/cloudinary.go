package media

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)


func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    if fileHeader.Size == 0 {
        return "", fmt.Errorf("image file is empty")
    }

    cloudName, apiKey, apiSecret := config.CLOUDINARYCREDINTIALS()
	fmt.Println("cloudName:", cloudName)
	fmt.Println("apiKey:", apiKey)
	fmt.Println("apiSecret:", apiSecret)

    cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
    if err != nil {
        return "", fmt.Errorf("cloudinary init error: %w", err)
    }

    
    file.Seek(0, 0)

    uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
        Folder: "recipes",
    })
    if err != nil {
        return "", fmt.Errorf("cloudinary upload error: %w", err)
    }

    fmt.Printf("Cloudinary Upload Full Response: %+v\n", uploadResult)

    if uploadResult.SecureURL == "" {
        return "", fmt.Errorf("cloudinary returned empty SecureURL")
    }

    return uploadResult.SecureURL, nil
}
