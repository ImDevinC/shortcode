package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/imdevinc/url_shortener/db"
)

var errorLogger = log.New(os.Stderr, "ERROR", log.Llongfile)
var debugLogger = log.New(os.Stderr, "DEBUG", log.Llongfile)

// NewLinkRequest ...
type NewLinkRequest struct {
	URI string `json:"uri"`
}

func getLinkFromShortCode(database *db.Database, code string) (events.APIGatewayProxyResponse, error) {
	code = strings.TrimPrefix(code, "/")
	link, err := database.GetShortLink(code)
	if err != nil {
		return serverError(err)
	}
	if len(link) == 0 {
		return clientError(http.StatusNotFound)
	}

	headers := make(map[string]string)
	headers["Location"] = link
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers:    headers,
	}, nil
}

func createLinkWithRandomShortCode(database *db.Database, uri string) (string, error) {
	existingCode, err := database.FindShortCodeByURI(uri)
	if err != nil {
		errorLogger.Println(err)
	}

	if len(existingCode) != 0 {
		return existingCode, nil
	}

	count := 0
	var finalCode string
	for count < 5 {
		code := db.GenerateShortCode(5)
		err := database.InsertShortLink(uri, code)
		if err == nil {
			finalCode = code
			break
		} else if strings.HasPrefix(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
			count++
		} else {
			return "", err
		}
	}

	return finalCode, nil
}

func createNewShortCode(database *db.Database, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := b64.StdEncoding.DecodeString(body)
		if err != nil {
			return serverError(err)
		}
		body = string(decoded)
	}

	newLink := new(NewLinkRequest)
	err := json.Unmarshal([]byte(body), newLink)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	var code string
	if len(req.Path) > 0 {
		customCode := strings.TrimPrefix(req.Path, "/")
		err = database.InsertShortLink(newLink.URI, customCode)
		code = customCode
	} else {
		code, err = createLinkWithRandomShortCode(database, newLink.URI)
	}

	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       code,
	}, nil
}

func deleteShortCode(db *db.Database, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Authorization"] != "Password1!" {
		return clientError(http.StatusUnauthorized)
	}

	path := strings.TrimPrefix(req.Path, "/")
	err := db.DeleteShortCode(path)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	database := db.New()
	switch method := req.HTTPMethod; method {
	case "GET":
		return getLinkFromShortCode(&database, req.Path)
	case "POST":
		return createNewShortCode(&database, req)
	case "DELETE":
		return deleteShortCode(&database, req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
