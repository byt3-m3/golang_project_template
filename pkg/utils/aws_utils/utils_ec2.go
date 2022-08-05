package aws_utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"time"
)

func GetEC2Client() *ec2.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ec2.NewFromConfig(cfg)

}

type CreateEC2InstanceInput struct {
	ImageId      string
	InstanceType string
	MaxCount     int32
	MinCount     int32
	SSHKeyName   string
	Timeout      time.Duration
}

func CreateEC2Instance(input *CreateEC2InstanceInput) (*ec2.RunInstancesOutput, error) {
	client := GetEC2Client()

	ctx := context.Background()
	c, cancelFn := context.WithTimeout(ctx, input.Timeout)
	defer cancelFn()

	return client.RunInstances(c, &ec2.RunInstancesInput{
		ImageId:      aws.String(input.ImageId),
		InstanceType: "t2.micro",
		MinCount:     aws.Int32(input.MinCount),
		MaxCount:     aws.Int32(input.MaxCount),
		KeyName:      aws.String(input.SSHKeyName),
	})

}

func CheckInstanceForNonPendingState(svc *ec2.Client, instanceId string) ec2types.Instance {
	// Get Validates a Single AWS instance status

	iids := []string{instanceId}
	params := &ec2.DescribeInstancesInput{
		InstanceIds: iids,
	}

	dio, err := svc.DescribeInstances(context.TODO(), params)
	if err != nil {
		log.Fatalln(err)
	}

	return dio.Reservations[0].Instances[0]

}

type TerminateEC2InstanceInput struct {
	InstanceId string
}

func TerminateEC2Instance(input *TerminateEC2InstanceInput) error {

	client := GetEC2Client()
	iids := []string{input.InstanceId}
	params := &ec2.TerminateInstancesInput{
		InstanceIds: iids,
		DryRun:      nil,
	}

	_, err := client.TerminateInstances(context.TODO(), params)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Successfully Terminated InstanceID='%s'", input.InstanceId))
	return nil

}
