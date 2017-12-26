package goclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// @todo think about formats other than JSON

// Url is the address of a 1backend installation.
// The default value is https://1backend.com:9993, change this to your own installation be either changing this variable
// or by having an environment variable with the name 1BACKEND_URL on your host.
var Url = "https://1backend.com:9993"

// CallerId is a secret key, specific to each project, which gets translated to a project name for namespacing
// by the proxy. This enables 1backend apps to namespace their database based on the caller's identity.
// It comes from the environment variable name CALLER_ID.
//
// To read more see: https://github.com/1backend/1backend/blob/master/docs/namespacing.md
var CallerId = ""

func init() {
	sa := os.Getenv("1BACKEND_URL")
	if sa != "" {
		Url = sa
	}
	CallerId = os.Getenv("CALLER_ID")
}

type GoClient struct {
	Token string
	Url   string
	// CallerId - by default this comes from the CALLER_ID environment variable.
	// You can modify it if you want to pass on a caller id coming from the current request header.
	CallerId string
}

func New(token string) GoClient {
	return GoClient{
		Token:    token,
		Url:      Url,
		CallerId: CallerId,
	}
}

func (g GoClient) Call(author, projectName, method, path string, input map[string]interface{}, result interface{}) error {
	url := fmt.Sprintf("%v/app/%v/%v%v", g.Url, author, projectName, path)
	body := []byte{}
	if method == "POST" || method == "PUT" {
		bs, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bs
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header = make(http.Header)
	req.Header.Set("token", g.Token)
	req.Header.Set("caller-id", g.CallerId)
	if method == "GET" || method == "PUT" {
		for key, value := range input {
			req.Form.Add(key, fmt.Sprintf("%v", value))
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(responseBody, result)
}
