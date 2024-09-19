package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	twitchApiUrl = "https://id.twitch.tv/oauth2/token"
	igdbApiUrl   = "https://api.igdb.com/v4"
)

var (
	twitchClientId = os.Getenv("TWITCH_CLIENT_ID")
	twitchSecret   = os.Getenv("TWITCH_SECRET")
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type countResponse struct {
	Count int `json:"count"`
}

type IgdbGame struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

func GetAuthToken() (*authResponse, error) {
	authResponse := &authResponse{}
	authUrl := fmt.Sprintf("%s?client_id=%s&client_secret=%s&grant_type=client_credentials", twitchApiUrl, twitchClientId, twitchSecret)
	resp, err := http.Post(authUrl, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

func GetTotalGames(accessToken string, httpClient *http.Client) (int, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", igdbApiUrl, "/games/count"), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	countResponse := countResponse{}

	err = json.Unmarshal(body, &countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func GetGames(offset, limit int, accessToken string, httpClient *http.Client) ([]IgdbGame, error) {
	jsonString := fmt.Sprintf("fields id, name, slug; sort created_at asc; limit %d; offset %d;", limit, offset)
	jsonBytes := []byte(jsonString)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", igdbApiUrl, "/games"), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", twitchClientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gamesResponse := []IgdbGame{}

	err = json.Unmarshal(body, &gamesResponse)
	if err != nil {
		return nil, err
	}

	return gamesResponse, err

}
