package db

import (
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Database structure
type Database struct {
	db *dynamodb.DynamoDB
}

type shortcodeRow struct {
	Shortcode string `json:"shortCode" dynamodbav:"shortCode"`
	URI       string `json:"URI" dynamodbav:"URI"`
}

// APIKeyRow structure
type APIKeyRow struct {
	APIKey string `json:"apikey,omitempty" dynamodbav:"apikey"`
	Role   string `json:"role,omitempty" dynamodbav:"role"`
}

const cleanCharSet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// New creates new instance of Database manager
func New() Database {
	rand.Seed(time.Now().UTC().UnixNano())
	db := dynamodb.New(session.New())
	d := Database{db}
	return d
}

// GenerateShortCode creates a new short code with specified
// string length
func GenerateShortCode(length int, clean bool) string {
	chars := charSet
	if clean {
		chars = cleanCharSet
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// FindShortCodeByURI will find an existing shortcode based on
// the provided URI
func (d *Database) FindShortCodeByURI(uri string) (string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("DYNAMO_DB_TABLENAME")),
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

// InsertShortCode will create the new specified shortcode pointing
// to the specified URI
func (d *Database) InsertShortCode(uri string, code string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLENAME")),
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

// DeleteShortCode will remove the specified shortcode
func (d *Database) DeleteShortCode(code string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLENAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"URI": {
				S: aws.String(code),
			},
		},
	}

	_, err := d.db.DeleteItem(input)
	return err
}

// GetShortCode will remove the specified shortcode
func (d *Database) GetShortCode(code string) (string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLENAME")),
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

// CreateNewAPIKey inserts a new API key  with specified role
func (d *Database) CreateNewAPIKey(key string, role string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("shortcode-api-keys"),
		Item: map[string]*dynamodb.AttributeValue{
			"apikey": {
				S: aws.String(key),
			},
			"role": {
				S: aws.String(role),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(apikey)"),
	}

	_, err := d.db.PutItem(input)
	return err
}

// GetAPIKey returns the information for the specified API key
func (d *Database) GetAPIKey(key string) (*APIKeyRow, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("shortcode-api-keys"),
		Key: map[string]*dynamodb.AttributeValue{
			"apikey": {
				S: aws.String(key),
			},
		},
	}

	result, err := d.db.GetItem(input)
	if err != nil {
		return &APIKeyRow{}, err
	}

	if result.Item == nil {
		return &APIKeyRow{}, nil
	}

	apikeyResult := &APIKeyRow{}
	err = dynamodbattribute.UnmarshalMap(result.Item, apikeyResult)
	if err != nil {
		return &APIKeyRow{}, err
	}

	return apikeyResult, nil
}

// DeleteAPIKey will remove the specified API key
func (d *Database) DeleteAPIKey(key string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("shortcode-api-keys"),
		Key: map[string]*dynamodb.AttributeValue{
			"apikey": {
				S: aws.String(key),
			},
		},
	}

	_, err := d.db.DeleteItem(input)
	return err
}
