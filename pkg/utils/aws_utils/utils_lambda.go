package aws_utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/byt3-m3/vne_go/pkg/vne_utils"
	"log"
)

func GetLambdaClient() *lambda.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return lambda.NewFromConfig(cfg)

}

type SyncInvokedLambdaFuncInput struct {
	FunctionName string
	Payload      []byte
}
type SyncInvokedLambdaFuncOutput struct {
}

func SyncInvokeLambdaFunc(params SyncInvokedLambdaFuncInput) (SyncInvokedLambdaFuncOutput, error) {
	invokeInput := &lambda.InvokeInput{
		FunctionName:   aws.String(params.FunctionName),
		ClientContext:  nil,
		InvocationType: "",
		LogType:        "",
		Payload:        params.Payload,
		Qualifier:      nil,
	}
	client := GetLambdaClient()
	output, err := client.Invoke(context.TODO(), invokeInput)
	if err != nil {
		log.Println(err)
		return SyncInvokedLambdaFuncOutput{}, err
	}
	fmt.Println(vne_utils.SerializeStructToJSON(output))

	return SyncInvokedLambdaFuncOutput{}, nil
}
