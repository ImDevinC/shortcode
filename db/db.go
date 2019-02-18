package db

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Database struct {
	db *dynamodb.DynamoDB
}

type shortcodeRow struct {
	Shortcode string `json:"shortCode" dynamodbav:"shortCode"`
	URI       string `json:"URI" dynamodbav:"URI"`
}

const charSet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// New ...
func New() Database {
	rand.Seed(time.Now().UTC().UnixNano())
	db := dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	d := Database{db}
	return d
}

func GenerateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

func (d *Database) FindShortCodeByURI(uri string) (string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("shortcodes"),
		IndexName:              aws.String("URI-index"),
		KeyConditionExpression: aws.String("#URI = :uri"),
		ExpressionAttributeNames: map[string]*string{
			"#URI": aws.String("URI"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":uri": {
				S: aws.String(uri),
			},
		},
	}

	result, err := d.db.Query(input)
	if err != nil {
		return "", err
	}

	if *result.Count == 0 {
		return "", nil
	}

	shortCodeResult := []shortcodeRow{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &shortCodeResult)
	if err != nil {
		return "", err
	}

	return shortCodeResult[0].Shortcode, nil
}

func (d *Database) InsertShortLink(uri string, code string) error {
	fmt.Printf("Creating new shortlink %s\n", code)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("shortcodes"),
		Item: map[string]*dynamodb.AttributeValue{
			"shortcode": {
				S: aws.String(code),
			},
			"URI": {
				S: aws.String(uri),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(shortcode)"),
	}

	_, err := d.db.PutItem(input)
	return err
}

func (d *Database) DeleteShortCode(code string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("shortcodes"),
		Key: map[string]*dynamodb.AttributeValue{
			"URI": {
				S: aws.String(code),
			},
		},
	}

	_, err := d.db.DeleteItem(input)
	return err
}

func (d *Database) GetShortLink(code string) (string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("shortcodes"),
		Key: map[string]*dynamodb.AttributeValue{
			"shortcode": {
				S: aws.String(code),
			},
		},
	}

	result, err := d.db.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", nil
	}

	shortCodeResult := new(shortcodeRow)
	err = dynamodbattribute.UnmarshalMap(result.Item, shortCodeResult)
	if err != nil {
		return "", err
	}

	return shortCodeResult.URI, nil
}
