package aws_utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/byt3-m3/project/pkg/utils/struct_utils"
	"log"
)

type DynamoClient interface {
	GetAWSDynamoDBClient() *dynamodb.Client
	AWSDynamoItemWriter
	AWSDynamoIndexQuerier
	AWSDynamoQuerier
	AWSDynamoScanner
}

type AWSDynamoScanner interface {
	Scan(ctx context.Context, client *dynamodb.Client, input *ScanDynamoInput) (*dynamodb.ScanOutput, error)
}

type AWSDynamoQuerier interface {
	Query(ctx context.Context, input QueryInput) (*dynamodb.QueryOutput, error)
}

type AWSDynamoIndexQuerier interface {
	QueryIndex(ctx context.Context, input *QueryIndexInput) (*dynamodb.QueryOutput, error)
}

type AWSDynamoItemWriter interface {
	SaveItem(ctx context.Context, tableName string, item interface{}) (*dynamodb.PutItemOutput, error)
}

type DynamoClientConfig struct {
	region string
}

type dynamoClient struct {
	client *dynamodb.Client
}

func NewDynamoClient(cfg *DynamoClientConfig) *dynamoClient {
	return &dynamoClient{
		client: GetAWSDynamoDBClient(cfg.region),
	}
}

func (c dynamoClient) GetAWSDynamoDBClient() *dynamodb.Client {
	return c.client
}

func (c dynamoClient) SaveItem(ctx context.Context, tableName string, item interface{}) (*dynamodb.PutItemOutput, error) {
	return saveItemToDynamoTable(ctx, c.client, tableName, item)
}

func (c dynamoClient) QueryIndex(ctx context.Context, input *QueryIndexInput) (*dynamodb.QueryOutput, error) {
	return queryIndex(ctx, c.client, input)
}

func (c dynamoClient) Query(ctx context.Context, input QueryInput) (*dynamodb.QueryOutput, error) {
	return query(ctx, c.client, input)
}

func (c dynamoClient) Scan(ctx context.Context, client *dynamodb.Client, input *ScanDynamoInput) (*dynamodb.ScanOutput, error) {
	return scan(ctx, client, input)
}

func GetAWSDynamoDBClient(region string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg)

}

type QueryIndexInput struct {
	TableName     string
	IndexName     string
	KeyConditions map[string]types.Condition
}

func queryIndex(ctx context.Context, client *dynamodb.Client, input *QueryIndexInput) (*dynamodb.QueryOutput, error) {

	queryInput := dynamodb.QueryInput{
		IndexName:     aws.String(input.IndexName),
		KeyConditions: input.KeyConditions,
		TableName:     aws.String(input.TableName),
	}

	return client.Query(ctx, &queryInput)

}

type QueryInput struct {
	TableName     string
	KeyConditions map[string]types.Condition
}

func query(ctx context.Context, client *dynamodb.Client, input QueryInput) (*dynamodb.QueryOutput, error) {
	queryInput := dynamodb.QueryInput{
		KeyConditionExpression: nil,
		KeyConditions:          input.KeyConditions,
		TableName:              aws.String(input.TableName),
	}

	return client.Query(ctx, &queryInput)

}

type ScanDynamoInput struct {
	TableName string
}

func scan(ctx context.Context, client *dynamodb.Client, input *ScanDynamoInput) (*dynamodb.ScanOutput, error) {
	params := dynamodb.ScanInput{
		TableName: aws.String(input.TableName),
	}

	return client.Scan(ctx, &params)

}

func makePutItemInput(tableName string, item interface{}) *dynamodb.PutItemInput {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling item: %s", err)
	}
	fmt.Println(struct_utils.SerializeStructToJSON(av))
	return &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
}

func saveItemToDynamoTable(ctx context.Context, client *dynamodb.Client, tableName string, item interface{}) (*dynamodb.PutItemOutput, error) {
	input := makePutItemInput(tableName, item)
	return client.PutItem(ctx, input)

}
