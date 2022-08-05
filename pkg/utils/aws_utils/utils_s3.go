package aws_utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"os"
)

func GetS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}

type DownloadObjectFromBucketInput struct {
	BucketID string
	Key      string
	FileName string
}

type DownloadObjectFromBucketOutput struct {
	IsComplete bool
}

func DownloadObjectFromBucket(input *DownloadObjectFromBucketInput) (DownloadObjectFromBucketOutput, error) {

	getObjectinput := s3.GetObjectInput{
		Bucket: aws.String(input.BucketID),
		Key:    aws.String(input.Key),
	}

	file, err := os.Create(input.FileName)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to Create File: %s", input.FileName))
	}
	defer file.Close()

	downloader := GetS3Downloader()

	_, err = downloader.Download(context.TODO(), file, &getObjectinput)
	if err != nil {
		return DownloadObjectFromBucketOutput{}, err

	}

	return DownloadObjectFromBucketOutput{
		IsComplete: true,
	}, nil
}

func GetS3Downloader() *manager.Downloader {
	sess := GetS3Client()
	return manager.NewDownloader(sess)
}

type SaveBytesToS3Input struct {
	S3Bucket string
	Key      string
	Body     *[]byte
	FileType string
}

func SaveBytesToS3(input *SaveBytesToS3Input) error {

	sess := GetS3Client()
	fileReader := bytes.NewReader(*input.Body)

	putObjectInput := s3.PutObjectInput{
		Body:   fileReader,
		Bucket: aws.String(input.S3Bucket),
		Key:    aws.String(input.Key),
	}
	log.Println("Please Wait, Processing Request: ", putObjectInput)
	_, err := sess.PutObject(context.TODO(), &putObjectInput)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Successfully Uploaded Item: %s to Bucket: %s", input.Key, input.S3Bucket))
	return nil
}
