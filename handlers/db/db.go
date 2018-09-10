package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

const (
	Region        = endpoints.ApNortheast1RegionID
	LinkTableName = "link"
)

type Database interface {
	GetItem(interface{})
	PutItem(interface{})
}

type DB struct {
	instance *dynamodb.DynamoDB
}

type Link struct {
	ShortURL    string `json:"shorten_resource"`
	OriginalURL string `json:"original_url"`
}

func New() DB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Region)}),
	)

	return DB{instance: dynamodb.New(sess)}
}
func (d DB) GetLinkByShortenResource(i interface{}) (string, error) {
	item, err := d.instance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(LinkTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"shorten_resource": {
				S: aws.String(i.(string)),
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get item")
	}
	if item.Item == nil {
		return "", errors.New("no content")
	}

	link := Link{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &link)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal item")
	}

	return link.OriginalURL, nil
}

func (d DB) PutItem(i interface{}) (interface{}, error) {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(LinkTableName),
	}
	item, err := d.instance.PutItem(input)
	if err != nil {
		return nil, err
	}

	return item, nil
}