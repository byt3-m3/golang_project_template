package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetAWSSessionWithSharedCreds() *session.Session {

	config := aws.Config{
		Region: aws.String("us-east-1"),
	}
	//creds := credentials.NewEnvCredentials()
	//profile := "default"
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	//SharedConfigState: session.SharedConfigEnable,
	//	Profile: profile,
	//}))

	sess, _ := session.NewSession(&config)
	return sess
}
