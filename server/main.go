package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"url_shortener/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

// NewLinkRequest ...
type NewLinkRequest struct {
	URI string `json:"uri"`
}

func getLinkFromShortCode(database *db.Database, code string) (events.APIGatewayProxyResponse, error) {
	code = strings.TrimPrefix(code, "/")
	link, err := database.GetShortCode(code)
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
		logger.Warn().Err(err)
	}

	if len(existingCode) != 0 {
		return existingCode, nil
	}

	count := 0
	var finalCode string
	for count < 5 {
		code := db.GenerateShortCode(5, true)
		err := database.InsertShortCode(uri, code)
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
	// contentType := strings.ToLower(req.Headers["Content-Type"])
	// if !strings.HasPrefix(contentType, "application/json") {
	// 	fmt.Printf("%+v", req)
	// 	return clientError(http.StatusNotAcceptable)
	// }

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
	path := strings.TrimPrefix(req.Path, "/")
	if len(path) > 0 {
		err = database.InsertShortCode(newLink.URI, path)
		code = path
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

func deleteShortCode(database *db.Database, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !isAuthorized(database, req) {
		return clientError(http.StatusUnauthorized)
	}

	path := strings.TrimPrefix(req.Path, "/")
	err := database.DeleteShortCode(path)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func isAuthorized(database *db.Database, req events.APIGatewayProxyRequest) bool {
	auth := req.Headers["Authorization"]
	if len(auth) == 0 {
		return false
	}

	result, err := database.GetAPIKey(auth)
	if err != nil {
		logger.Info().Err(err).Msg("Failed to get API key from database")
		return false
	}

	return result.Role == "admin"
}

func generateAPIKey(database *db.Database, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := b64.StdEncoding.DecodeString(body)
		if err != nil {
			return serverError(err)
		}
		body = string(decoded)
	}

	apikeyInput := &db.APIKeyRow{}
	err := json.Unmarshal([]byte(body), apikeyInput)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	if len(apikeyInput.Role) == 0 {
		apikeyInput.Role = "admin"
	}

	count := 0
	var finalCode string
	for count < 5 {
		code := db.GenerateShortCode(26, false)
		err := database.CreateNewAPIKey(code, apikeyInput.Role)
		if err == nil {
			finalCode = code
			break
		} else if strings.HasPrefix(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
			count++
		} else {
			return serverError(err)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       finalCode,
	}, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
	logger.Error().Err(err)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
