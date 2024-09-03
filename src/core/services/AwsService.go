package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"time"
)

func savePropositionImageInAwsS3(propositionCode int, image []byte) (string, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := ""

	awsConfig := aws.Config{
		Region:      awsRegion,
		Credentials: credentials.NewStaticCredentialsProvider(awsAccessKeyId, awsAccessSecretKey, awsSessionToken),
	}

	s3Client := s3.NewFromConfig(awsConfig)

	bucket := os.Getenv("AWS_S3_BUCKET")
	key := fmt.Sprintf("propositions/%d_%s.png", propositionCode, time.Now().Format("150405_02012006"))

	uploader := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(image),
		ContentType: aws.String(http.DetectContentType(image)),
	}

	_, err := s3Client.PutObject(context.TODO(), &uploader)
	if err != nil {
		log.Error("Erro ao enviar a imagem para o AWS S3: ", err)
		return "", err
	}

	imageUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)

	return imageUrl, nil
}
