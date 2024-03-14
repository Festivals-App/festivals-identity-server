package token

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Validation interface {
	ValidateAccessToken(token string) (string, error)
}

type ValidationService struct {
	Key      *rsa.PublicKey
	APIKeys  *[]string
	Client   *http.Client
	Endpoint string
}

func NewValidationService(endpoint string, clientCert string, clientKey string, serviceKey string) *ValidationService {

	client, err := validationClient(clientCert, clientKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create validation client.")
	}

	keys, err := loadAPIKeys(endpoint, serviceKey, client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load API keys from identity service.")
	}

	vaidationKey, err := loadValidationKey(endpoint, serviceKey, client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load API keys from identity service.")
	}

	return &ValidationService{Key: vaidationKey, APIKeys: &keys, Client: client, Endpoint: endpoint}
}

// ValidateAccessToken parses and validates the given access token
// returns the custom claim present in the token payload
func (validator *ValidationService) ValidateAccessToken(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error().Msg("Unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth token")
		}
		return validator.Key, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Unable to parse claims")
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid || claims.UserID == "" {
		return nil, errors.New("invalid token: authentication failed")
	}
	return claims, nil
}

func validationClient(clientCert string, clientKey string) (*http.Client, error) {

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

func loadValidationKey(endpoint string, serviceKey string, client *http.Client) (*rsa.PublicKey, error) {

	request, err := http.NewRequest(http.MethodGet, endpoint+"/validation-key", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-x509-user-cert; charset=utf-8")
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
		return nil, errors.New("failed to retrieve validation key from identity service with error response: " + string(resBody))
	}

	return jwt.ParseRSAPublicKeyFromPEM(resBody)
}

func loadAPIKeys(endpoint string, serviceKey string, client *http.Client) ([]string, error) {

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
