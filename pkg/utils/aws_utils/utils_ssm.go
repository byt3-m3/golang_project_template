package aws_utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log"
)

func GetSSMClient() *ssm.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return ssm.NewFromConfig(cfg)
}

func GetSSMParams(paramName string, isEncrypted bool) (*ssm.GetParameterOutput, error) {
	apiParams := ssm.GetParameterInput{Name: aws.String(paramName), WithDecryption: isEncrypted}
	client := GetSSMClient()

	output, err := client.GetParameter(context.TODO(), &apiParams)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to get params '%s'", paramName))
		log.Println(err)
		return &ssm.GetParameterOutput{}, err
	}
	return output, nil
}
