package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestExistingCode(t *testing.T) {
	const existingCode = "IMDEVINC"

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "{\"URI\": \"http://imdevinc.com\"}",
	}

	resp, err := handleRequest(request)
	if err != nil {
		t.Error(err.Error())
	} else if resp.StatusCode != 200 {
		t.Error("Response code is not 200")
	} else if resp.Body != existingCode {
		t.Errorf("Invalid code returned %s, should be %s", resp.Body, existingCode)
	}
}

func TestGetShortCode(t *testing.T) {
	const code = "IMDEVINC"
	const expectedURL = "http://imdevinc.com"

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       fmt.Sprintf("/%s", code),
	}

	resp, err := handleRequest(request)
	if err != nil {
		t.Error(err.Error())
	} else if resp.StatusCode != 301 {
		t.Errorf("Response code is not 301, %d", resp.StatusCode)
	} else if resp.Headers["Location"] != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, resp.Headers["Location"])
	}
}

func TestCustomShortCode(t *testing.T) {
	const code = "IMDEVINC"
	const URL = "http://imdevinc.com"

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Path:       fmt.Sprintf("/%s", code),
		Body:       fmt.Sprintf("{\"URI\": \"%s\"}", URL),
	}

	resp, err := handleRequest(request)
	if err != nil {
		t.Error(err.Error())
	} else if resp.StatusCode != 200 {
		t.Errorf("Response code is not 200, %d", resp.StatusCode)
	} else if resp.Body != code {
		t.Errorf("Invalid code returned %s, should be %s", resp.Body, code)
	}
}

func TestDeleteShortCut(t *testing.T) {
	const code = "IMDEVINC"

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "DELETE",
		Path:       fmt.Sprintf("/%s", code),
	}

	resp, err := handleRequest(request)
	if err != nil {
		t.Error(err.Error())
	} else if resp.StatusCode != 200 {
		t.Errorf("Response code is not 200, %d", resp.StatusCode)
	}
}

func TestCreateNewShortCode(t *testing.T) {
	const URL = "https://blog.imdevinc.com"

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"URI\": \"%s\"}", URL),
	}

	resp, err := handleRequest(request)
	if err != nil {
		t.Error(err.Error())
	} else if resp.StatusCode != 200 {
		t.Errorf("Response code is not 200, %d", resp.StatusCode)
	}

	t.Logf("New code created: %s", resp.Body)
}
