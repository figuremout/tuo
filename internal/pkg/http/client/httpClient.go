package client

import (
	"io"
	"net/http"
	"strings"

	jsonUtil "github.com/githubzjm/tuo/internal/pkg/utils/json"
)

var HttpClient *http.Client
var URL string

const (
	UserAgent   = "CLI/version" // TODO
	Connection  = "keep-alive"
	ContentType = "application/json"
)

func InitHttpClient(url string) {
	HttpClient = &http.Client{ // config here
		// Timeout: 3 * time.Second,
	}
	URL = url
}

func Request(method string, api string, body interface{}, header map[string]string, targetResp interface{}) (int, error) {
	payload := strings.NewReader(string(jsonUtil.Dumps(body)))
	req, err := http.NewRequest(method, URL+api, payload)
	if err != nil {
		return -1, err
	}
	SetHeader(req, header)

	// send req
	resp, err := HttpClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	// parse body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// parse the bytes into struct
	if err := jsonUtil.Loads(respBody, targetResp); err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

func SetHeader(req *http.Request, header map[string]string) {
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Connection", Connection)
	req.Header.Set("User-Agent", UserAgent)
	if header == nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
}
