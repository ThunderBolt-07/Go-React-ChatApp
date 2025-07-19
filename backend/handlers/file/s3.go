package file

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func MakeAWSClient() (*s3.Client, error) {
	godotenv.Load()
	aws_key := os.Getenv("AWS_ACCESS_KEY_ID")
	if aws_key == "" {
		log.Fatal("uable to load aws key from env")
	}
	aws_secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if aws_secret == "" {
		log.Fatal("uable to load aws secret from env")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			aws_key,
			aws_secret,
			"", // session token, optional
		)),
		//config.WithRegion("us-west-1"), // change as needed
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return nil, err
	}

	//log.Println("config created successfully", cfg)

	client := s3.NewFromConfig(cfg)

	// err = CreateBucket(client)

	// if err != nil {
	// 	log.Println("error while making bucket ", err)
	// 	return nil,err

	// }

	// s, e := GetFileDownlaodUrl(client, "1000043640.jpg")

	// if e != nil {
	// 	log.Println("error e  while getitng url ", e)
	// 	return nil, e
	// }
	// log.Println("success url  ", s)

	return client, nil
}

func CreateBucket(client *s3.Client) error {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*100))
	defer cancel()
	opt, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("gohibucket"),
	})

	if err != nil {

		return err
	}

	log.Println("bucket created ", opt)

	return nil
}

func UploadFile(client *s3.Client, file *os.File, fileName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
	defer cancel()
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("gohibucket"),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return err
	}
	err = s3.NewObjectExistsWaiter(client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String("gohibucket"), Key: aws.String(fileName)}, time.Minute)
	return err
}

func GetFileDownlaodUrl(client *s3.Client, fileName string) (string, error) {
	log.Println("trying to get url....")
	presignClient := s3.NewPresignClient(client)
	bucketName := "gohibucket"

	request, err := GetObject(presignClient, bucketName, fileName)

	if err != nil {
		// log.Println("error occured while getting presign url",err)
		return "", err
	}
	log.Println("signing url success", request)
	return request.URL, nil
}
func GetObject(presignClient *s3.PresignClient, bucket string, fileName string) (*v4.PresignedHTTPRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     &bucket,
		Key:                        &fileName,
		ResponseContentDisposition: aws.String("attachment"),
	}, func(po *s3.PresignOptions) {
		po.Expires = time.Second * 100
	})
	return request, err
}
