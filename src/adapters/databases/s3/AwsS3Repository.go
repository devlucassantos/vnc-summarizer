package s3

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
	"vnc-summarizer/utils/datetime"
)

type AwsS3 struct{}

func NewAwsS3() *AwsS3 {
	return &AwsS3{}
}

func (instance AwsS3) SavePropositionImage(propositionCode int, image []byte) (string, error) {
	log.Info("Starting to register the image of proposition ", propositionCode)

	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")

	awsConfig := aws.Config{
		Region:      awsRegion,
		Credentials: credentials.NewStaticCredentialsProvider(awsAccessKeyId, awsAccessSecretKey, awsSessionToken),
	}

	s3Client := s3.NewFromConfig(awsConfig)

	bucket := os.Getenv("AWS_S3_BUCKET")

	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return "", err
	}

	key := fmt.Sprintf("propositions/%d_%s.png", propositionCode,
		currentDateTime.Format("150405_02012006"))

	uploader := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(image),
		ContentType: aws.String(http.DetectContentType(image)),
	}

	_, err = s3Client.PutObject(context.TODO(), &uploader)
	if err != nil {
		log.Errorf("Error saving image of proposition %d to AWS S3: %s", propositionCode, err.Error())
		return "", err
	}

	imageUrl := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)

	log.Infof("The image of proposition %d was successfully registered: %s", propositionCode, imageUrl)
	return imageUrl, nil
}
