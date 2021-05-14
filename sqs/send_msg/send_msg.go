package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
)

func main() {
	queue := "movie-events.fifo"
	region := "us-east-1"
	endpoint := "http://localhost:4566"
	msgToSend := "Best movie ever"

	sqsClient, err := newSQSClient(region, endpoint)
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	q, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{QueueName: &queue})
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	log.Printf("⏰ sending message to queueUrl:'%v'\tregion:'%v'\tendpoint:'%v'\n", *q.QueueUrl, region, endpoint)
	if _, err := sendMessage(sqsClient, msgToSend, *q.QueueUrl); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	log.Printf("✅ success: message sent to queue\n")
}

func newSQSClient(region, endpoint string) (*sqs.Client, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, reg string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			})),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "loading AWS config")
	}
	return sqs.NewFromConfig(awsConfig), nil
}

func sendMessage(sqsClient *sqs.Client, msg, queueURL string) (*sqs.SendMessageOutput, error) {
	sqsMessage := sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(msg),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Predator"),
			},
			"Director": {
				DataType:    aws.String("String"),
				StringValue: aws.String("John McTiernan"),
			},
		},
	}

	output, err := sqsClient.SendMessage(context.Background(), &sqsMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "sending message to queue %v", queueURL)
	}
	return output, nil
}
