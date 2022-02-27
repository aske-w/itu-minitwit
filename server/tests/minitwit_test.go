package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

var baseUrl = "https://localhost:8080/"

func TestMain(m *testing.M) {
	//setup temporary instance && db here
	exitCode := m.Run()
	//tear down temporary instance && db here
	os.Exit(exitCode)
}

// prepends the baseUrl global variable to endpoint
func formatEndpoint(endpoint string) string {
	return fmt.Sprintf("%s%s", baseUrl, endpoint)
}

func readCloserToString(rc io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	return fmt.Sprintf(buf.String())
}

// Helper functions
func register(username, password, password2, email string) (*http.Response, error) {
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = fmt.Sprintf("%s@example.com", username)
	}

	data := url.Values{
		"username":  {username},
		"password":  {password},
		"password2": {password2},
		"email":     {email},
	}
	return http.PostForm(formatEndpoint("register"), data)
}

func login(username, password string) (*http.Response, error) {
	data := url.Values{
		"username": {username},
		"password": {password},
	}
	return http.PostForm(formatEndpoint("login"), data)
}

func register_and_login(username, password string) (*http.Response, error) {
	_, err := register(username, password, "", "")
	if err != nil {
		return nil, err
	}
	return login(username, password)
}

func logout(session string) (*http.Response, error) {
	request, err := http.NewRequest("GET", formatEndpoint("logout"), bytes.NewBufferString(""))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Cookie", fmt.Sprintf("session=\"%s\"", session))
	client := http.Client{}
	return client.Do(request)
}

// Attemps to get the value of a cookie with the given name in the given cookie array
func getCookie(cookies []*http.Cookie, name string) *string {
	for _, c := range cookies {
		if c.Name == name {
			return &c.Value
		}
	}
	return nil
}

func TestRegister(t *testing.T) {
	//case 1
	resp, err := register("user1", "default", "", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}

	//case 2
	resp, err = register("user1", "default", "", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}

	//case 3
	resp, err = register("", "default", "", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}

	//case 4
	resp, err = register("meh", "", "", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}

	//case 5
	resp, err = register("meh", "x", "y", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}

	//case 6
	resp, err = register("user1", "default", "", "broken")
	if err != nil {
		t.Error(err.Error())
	} else {
		subStr := "You were successfully registered and can login now"
		if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
			t.Errorf("substring '%s' does not exist in response body", subStr)
		}
	}
}

func TestLoginLogout(t *testing.T) {
	//case 1
	username := "user1"
	password := "default"
	resp, err := register_and_login(username, password)
	if err != nil {
		t.Fatal(err.Error())
	}
	subStr := "You were logged in"
	if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
		t.Fatalf("error logging in with username %s and password %s", username, password)
	}
	sess := getCookie(resp.Cookies(), "session")
	resp, err = logout(*sess)
	if err != nil {
		t.Fatal(err.Error())
	}
	subStr = "You were logged out"
	if !(strings.Contains(readCloserToString(resp.Body), subStr)) {
		t.Fatalf("error logging out with username %s", username)
	}

}
