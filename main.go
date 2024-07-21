package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := "meu-bucket"
	key := "test.txt"
	filepath := "test.txt"

	config, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("teste", "teste", "teste")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localstack:4566"}, nil
			})))
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	fileOpen, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fileOpen.Close()

	fileInfo, err := fileOpen.Stat()
	if err != nil {
		panic(err)
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          fileOpen,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String("text/plain"),
	}

	_, err = s3Client.PutObject(context.Background(), input)
	if err != nil {
		panic(err)
	}

	println("Objeto criado com sucesso")
}
