package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IService interface {
	AuthURL(ctx context.Context) string
	VerifyCode(ctx context.Context, code string) (uint64, error)
}

type Service struct {
	clientId     string
	clientSecret string
	client       *http.Client
}

func NewService(clientId string, clientSecret string) IService {
	return &Service{
		clientId:     clientId,
		clientSecret: clientSecret,
		client:       http.DefaultClient,
	}
}

func (s Service) AuthURL(ctx context.Context) string {
	endpoint := "https://github.com/login/oauth/authorize"
	params := url.Values{
		"client_id":    {"c54992dff1a03482b7de"},
		"redirect_uri": {"http://127.0.0.1:8080/users/github/authorize"},
		"scope":        {"users"},
	}

	return fmt.Sprintf("%s?%s", endpoint, params.Encode())
}

func (s Service) VerifyCode(ctx context.Context, code string) (uint64, error) {
	target := "https://github.com/login/oauth/access_token"
	params := url.Values{
		"client_id":     {s.clientId},
		"client_secret": {s.clientSecret},
		"code":          {code},
	}

	req, err := http.NewRequest(http.MethodPost, target, strings.NewReader(params.Encode()))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return 0, err
	}

	target = "https://api.github.com/user"

	req, err = http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var infoResponse struct {
		Id        uint64 `json:"id"`
		AvatarUrl string `json:"avatar_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Location  string `json:"location"`
	}

	err = json.NewDecoder(resp.Body).Decode(&infoResponse)
	if err != nil {
		return 0, err
	}
	return infoResponse.Id, nil
}
