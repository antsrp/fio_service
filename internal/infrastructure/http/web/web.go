package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	iweb "github.com/antsrp/fio_service/internal/interfaces/web"
	"go.uber.org/zap"
)

func MakeGETRequest(url string, options map[string]interface{}) *http.Request {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	for k, v := range options {
		q.Add(k, fmt.Sprintf("%v", v))
	}

	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-type", "application/json")

	return req
}

type WebConnection struct {
	client *http.Client
	logger *zap.Logger
}

func CreateNewConnection(log *zap.Logger) *WebConnection {
	c := &WebConnection{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		logger: log,
	}

	return c
}

func (c WebConnection) Do(request *http.Request, i interface{}) (int, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&i)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}

var _ iweb.Connector = WebConnection{}
