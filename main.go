package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Game struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

type EnvVars struct {
	TWITCH_CLIENT_ID  string
	TWITCH_SECRET     string
	DATABASE_USER     string
	DATABASE_DB       string
	DATABASE_PASSWORD string
	DATABASE_HOST     string
}

type CountResponse struct {
	Count int `json:"count" db:"count"`
}

func loadEnv() EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	return EnvVars{
		TWITCH_SECRET:     os.Getenv("TWITCH_SECRET"),
		TWITCH_CLIENT_ID:  os.Getenv("TWITCH_CLIENT_ID"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_DB:       os.Getenv("DATABASE_DB"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
	}
}

func getAuthToken(twitchClientId, twitchSecret string) AuthResponse {
	authResponse := AuthResponse{}
	authUrl := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", twitchClientId, twitchSecret)
	resp, err := http.Post(authUrl, "application/json", nil)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {

		log.Fatalln(err)
	}

	return authResponse
}

func getGames(limit, offset int, accessToken, clientId string, httpClient *http.Client) []Game {
	jsonString := fmt.Sprintf("fields id, name, slug; sort id asc; limit %d; offset %d;", limit, offset)
	jsonBytes := []byte(jsonString)
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Client-ID", clientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	gamesResponse := []Game{}

	err = json.Unmarshal(body, &gamesResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return gamesResponse
}

func getTotalGamesDb(db *sqlx.DB) int {
	countResponse := CountResponse{}
	db.Get(&countResponse, "SELECT count(id) FROM games")

	return countResponse.Count
}

func getTotalGamesApi(accessToken, clientId string, httpClient *http.Client) int {
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/count", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Client-ID", clientId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	countResponse := CountResponse{}

	err = json.Unmarshal(body, &countResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return countResponse.Count
}

func main() {
	env := loadEnv()

	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", env.DATABASE_USER, env.DATABASE_DB, env.DATABASE_PASSWORD, env.DATABASE_HOST)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	authResponse := getAuthToken(env.TWITCH_CLIENT_ID, env.TWITCH_SECRET)
	client := &http.Client{}

	dbCount := getTotalGamesDb(db)
	apiCount := getTotalGamesApi(authResponse.AccessToken, env.TWITCH_CLIENT_ID, client)

	if dbCount == apiCount {
		return
	}

	limit := 500
	offset := 0
	lastTime := time.Now().Add(-250 * 1e6)
	total := -1

	for offset < total {
		executionDuration := time.Since(lastTime).Milliseconds()
		if executionDuration < 250 {
			difference := 250 - executionDuration

			time.Sleep(time.Duration(difference * 1e6))
		}

		games := getGames(limit, offset, authResponse.AccessToken, env.TWITCH_CLIENT_ID, client)
		lastTime = time.Now()

		_, err = db.NamedExec("INSERT INTO games (api_id, name, slug) VALUES (:id, :name, :slug)", games)

		if err != nil {
			fmt.Printf("Stopped at offset %d\n", offset)
			log.Fatalln(err)
		}

		offset += 500
	}
}
