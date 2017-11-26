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

var DefaultUrl = "https://1backend.com:9993"

func init() {
	sa := os.Getenv("SERVICEADDRESS")
	if sa != "" {
		DefaultUrl = sa
	}
}

type GoClient struct {
	Token string
	Url   string
}

func New(token string) GoClient {
	return GoClient{
		Token: token,
		Url:   DefaultUrl,
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
