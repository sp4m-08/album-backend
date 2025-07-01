package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	BucketName = "album-images-saunok"
	Region     = "eu-north-1"
)

func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(Region)) //config.LoadDefaultConfig()-> automatically loads credentials
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	buffer := bytes.NewBuffer(nil)
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")
	fileExt := filepath.Ext(fileHeader.Filename)
	key = key + fileExt

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(contentType),
		//ACL:         types.ObjectCannedACLPublicRead, // make public
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BucketName, Region, key)
	return url, nil
}
