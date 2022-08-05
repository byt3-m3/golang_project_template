package aws_utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"log"
)

type SfnExecutionState string

const (
	SfnRunningState   SfnExecutionState = "RUNNING"
	SfnSucceededState SfnExecutionState = "SUCCEEDED"
	SfnAbortedState   SfnExecutionState = "ABORTED"
	SfnFailedState    SfnExecutionState = "FAILED"
	SfnUnknownState   SfnExecutionState = "UNKNOWN"
)

func GetSFNClient() *sfn.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return sfn.NewFromConfig(cfg)
}

type StartStateMachineExecutionInput struct {
	JsonInput     string
	SfnARN        string
	ExecutionName string
}

func StartStateMachineExecution(input *StartStateMachineExecutionInput) (*sfn.StartExecutionOutput, error) {
	sfnInput := sfn.StartExecutionInput{
		Input:           aws.String(input.JsonInput),
		Name:            aws.String(input.ExecutionName),
		StateMachineArn: aws.String(input.SfnARN),
	}

	client := GetSFNClient()
	return client.StartExecution(context.TODO(), &sfnInput)

}

func GetStateMachineExecutionState(executionArn string) (SfnExecutionState, error) {
	client := GetSFNClient()

	input := sfn.DescribeExecutionInput{ExecutionArn: aws.String(executionArn)}

	output, err := client.DescribeExecution(context.TODO(), &input)
	log.Println(fmt.Sprintf("Retrived State for ExecutionID='%s', Status=%s", *output.Name, output.Status))
	if err != nil {
		return "", err
	}

	switch output.Status {
	case "RUNNING":
		return SfnRunningState, nil

	case "SUCCEEDED":
		return SfnSucceededState, nil

	case "ABORTED":
		return SfnAbortedState, nil

	case "FAILED":
		return SfnFailedState, nil

	}

	return SfnUnknownState, nil
}

func AbortStateMachineExecution(executionArn string) (*sfn.StopExecutionOutput, error) {
	client := GetSFNClient()
	sfnInput := &sfn.StopExecutionInput{ExecutionArn: aws.String(executionArn)}

	return client.StopExecution(context.TODO(), sfnInput)
}

type SendSuccessTokenInput struct {
	Token    *string
	JSONData *string
}

func SendSuccessToken(input *SendSuccessTokenInput) error {
	client := GetSFNClient()
	params := &sfn.SendTaskSuccessInput{TaskToken: input.Token, Output: input.JSONData}
	_, err := client.SendTaskSuccess(context.TODO(), params)
	log.Println(fmt.Sprintf("Sent SuccessfulToken: %s", *input.Token))
	if err != nil {
		return err

	}

	return nil
}

type SendFailureTokenInput struct {
	Token  *string
	Error  *string
	Reason *string
}

func SendFailureToken(input *SendFailureTokenInput) error {
	client := GetSFNClient()
	params := &sfn.SendTaskFailureInput{TaskToken: input.Token, Error: input.Error, Cause: input.Reason}
	_, err := client.SendTaskFailure(context.TODO(), params)
	if err != nil {
		return err
	}
	return nil
}
