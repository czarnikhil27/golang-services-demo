package main

import (
	"fmt"
	"os"
	"os/signal"
	"queue/config"
	constants "queue/utils"
	"sync"
	"syscall"

	"queue/internal/s3Client"
	"queue/internal/snsclient"
	sqsclient "queue/internal/sqsClient"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func initEnv() (*config.Config, error) {
	config, err := config.LoadConfig("env.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config file")

		return nil, err
	}
	return config, nil
}
func initAWSSession(config *config.Config) *session.Session {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("eu-north-1"),
			Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyID, config.AwsSecretAccessKey, ""),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return session
}
func main() {
	config, err := initEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load env variables")
	}
	awsSession := initAWSSession(config)

	snsSVC := sns.New(awsSession, nil)
	if snsSVC == nil {
		log.Fatal().Err(err).Msgf("Unable to initialize sns instance")
	}
	sqsSVC := sqs.New(awsSession, nil)
	if sqsSVC == nil {
		log.Fatal().Err(err).Msgf("Unable to initialize sqs instance")
	}
	publisher := snsclient.NewSNSClient(snsSVC)

	for i := 0; i <= 5; i++ {
		publisher.SendMessage(constants.SNS_TOPIC, fmt.Sprintf("message %v", i))
	}

	// // consumer code
	// // this is a infinitely running go routine
	consumer := sqsclient.NewSQSClient(sqsSVC)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Subscribe(constants.SQS_QUEUE_1)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	wg.Wait()

	// s3
	s3SVC := s3.New(awsSession)
	s3Client := s3Client.NewS3Client(s3SVC)
	file, err := s3Client.ListFiles("nikhilbucket1234", "")
	if err != nil {
		log.Err(err).Msgf("Unable to read file")
	}
	fmt.Println(file)
}
