package aws_utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

func GetSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return sqs.NewFromConfig(cfg)
}

func GetSQSQueueURL(svc *sqs.Client, queueName *string) (string, error) {
	urlResult, err := svc.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: queueName,
	})

	if err != nil {
		return "", err
	}
	queueURL := urlResult.QueueUrl
	return *queueURL, nil
}

type SQSEvent struct {
	EventType    string  `json:"eventType,omitempty"`
	Data         string  `json:"data,omitempty"`
	EventReceipt *string `json:"eventReceipt,omitempty"`
}

func MakeSQSDeleteMessageInput(receipt *string, queueName *string) (*sqs.DeleteMessageInput, error) {

	sqsClient := GetSQSClient()

	queueUrl, err := GetSQSQueueURL(sqsClient, queueName)
	if err != nil {
		return &sqs.DeleteMessageInput{}, err
	}
	return &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: receipt,
	}, nil
}

func DeleteSQSMessage(input *sqs.DeleteMessageInput) error {
	sqsClient := GetSQSClient()

	_, err := sqsClient.DeleteMessage(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Println("Removed Message From Queue", input.QueueUrl)
	return nil
}

func RecvMessageFromQueue(queueName *string, maxNumberOfMessage int32, visibilityTimeout int32) ([]SQSEvent, error) {

	svc := GetSQSClient()

	queueURL, err := GetSQSQueueURL(svc, queueName)
	if err != nil {
		return []SQSEvent{}, err
	}

	sqsParams := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: *aws.Int32(maxNumberOfMessage),
		VisibilityTimeout:   *aws.Int32(visibilityTimeout),
	}

	msgResult, _ := svc.ReceiveMessage(context.TODO(), sqsParams)

	var messages []SQSEvent
	for _, message := range msgResult.Messages {

		var model SQSEvent
		fmt.Println(*message.Body)
		err := json.Unmarshal([]byte(*message.Body), &model)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Received Event: ", model)

		model.EventReceipt = message.ReceiptHandle
		messages = append(messages, model)

	}

	log.Println("Processed", len(messages), "Messages")
	return messages, nil
}

//type SendSQSIntegrationEventInput struct {
//	QueueName        string
//	IntegrationEvent common.IntegrationEvent
//}
//
//func SendSQSIntegrationEvent(input *SendSQSIntegrationEventInput) (*sqs.SendMessageOutput, error) {
//	client := GetSQSClient()
//
//	queueUrl, err := GetSQSQueueURL(client, &input.QueueName)
//	if err != nil {
//		return &sqs.SendMessageOutput{}, err
//	}
//
//	integrationEventStr, err := vne_utils.SerializeStructToJSON(input.IntegrationEvent)
//	if err != nil {
//		return &sqs.SendMessageOutput{}, err
//	}
//
//	sqsInput := sqs.SendMessageInput{
//		MessageBody: aws.String(integrationEventStr),
//		QueueUrl:    &queueUrl,
//	}
//	sendMessageOutput, err := client.SendMessage(context.TODO(), &sqsInput)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(sendMessageOutput)
//	return sendMessageOutput, err
//}
