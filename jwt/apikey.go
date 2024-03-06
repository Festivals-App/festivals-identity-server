package token

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID      int    `json:"api_key_id" sql:"api_key_id"`
	Key     string `json:"api_key" sql:"api_key"`
	Comment string `json:"api_key_comment" sql:"api_key_comment"`
}

func APIKeyClient(clientCert string, clientKey string) (*http.Client, error) {

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
	return client, nil
}

func LoadAPIKeys(endpoint string, serviceKey string, client *http.Client) ([]string, error) {

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("X-Request-ID", uuid.New().String())
	request.Header.Set("Service-Key", serviceKey)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve api keys from identity service with error response: " + string(resBody))
	}

	var data map[string]interface{}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, err
	}

	key_map := data["data"].([]APIKey)
	var rawKeys = []string{}
	for _, key := range key_map {
		rawKeys = append(rawKeys, key.Key)
	}
	return rawKeys, nil
}
