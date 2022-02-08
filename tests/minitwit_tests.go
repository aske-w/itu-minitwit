package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var baseUrl = "https://localhost:8080/"

// Setup an instance of the entire project using an empty test database.
// Should probably have some way of killing the instance when test ends.
func setupTest() {

}

type RequestType int

const (
	GET RequestType = iota
	POST
)

func (r RequestType) String() string {
	switch r {
	case GET:
		return "GET"
	case POST:
		return "POST"
	}
	return "Unknown request type"
}

type DataEncoding int

const (
	JSON DataEncoding = iota
)

func (e DataEncoding) String() string {
	switch e {
	case JSON:
		return "JSON"
	}
	return "Unknown encoding"
}

func doHttpRequest(url string, requestType RequestType, encoding DataEncoding, body []byte) (*http.Response, error) {
	switch requestType {
	case GET:
		return http.Get(url)
	case POST:
		return http.Post(url, encoding.String(), bytes.NewBuffer(body))
	}
	return nil, fmt.Errorf("unknown request type '%s'", requestType)
}

// Helper functions
func register(username string, password string, password2 string, email string) (*http.Response, error) {
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = fmt.Sprintf("%s@example.com", username)
	}
	data, err := json.Marshal(map[string]string{
		"username":  username,
		"password":  password,
		"password2": password2,
		"email":     email,
	})
	if err != nil {
		return nil, err
	}
	return doHttpRequest(fmt.Sprintf("%s%s", baseUrl, "register"), POST, JSON, data)
	//return result of POST to /register with params
}

func login(username, password string) (*http.Response, error) {
	data, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	return doHttpRequest(fmt.Sprintf("%s%s", baseUrl, "login"), POST, JSON, data)
}

func register_and_login(username, password string) (*http.Response, error) {
	_, err := register(username, password, "", "")
	if err != nil {
		return nil, err
	}
	return login(username, password)
}
