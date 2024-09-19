package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
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

type DbGame struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Slug  string `json:"slug" db:"slug"`
	ApiId int    `json:"api_id" db:"api_id"`
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

var db *sqlx.DB

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

func getAuthToken(twitchClientId, twitchSecret string) (*AuthResponse, error) {
	authResponse := &AuthResponse{}
	authUrl := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", twitchClientId, twitchSecret)
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

func getGames(limit, offset int, accessToken, clientId string, httpClient *http.Client) ([]Game, error) {
	jsonString := fmt.Sprintf("fields id, name, slug; sort created_at asc; limit %d; offset %d;", limit, offset)
	jsonBytes := []byte(jsonString)
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", clientId)
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

	gamesResponse := []Game{}

	err = json.Unmarshal(body, &gamesResponse)
	if err != nil {
		return nil, err
	}

	return gamesResponse, err
}

func getTotalGamesDb(db *sqlx.DB) (int, error) {
	countResponse := CountResponse{}
	err := db.Get(&countResponse, "SELECT count(id) FROM games")
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

func getTotalGamesApi(accessToken, clientId string, httpClient *http.Client) (int, error) {
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/count", nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Client-ID", clientId)
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

	countResponse := CountResponse{}

	err = json.Unmarshal(body, &countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}

type DbCredentials struct {
	User     string
	Database string
	Host     string
	Password string
}

func Connect() error {
	env := loadEnv()

	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", env.DATABASE_USER, env.DATABASE_DB, env.DATABASE_PASSWORD, env.DATABASE_HOST)

	var err error
	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		return err
	}

	return nil
}

func importGames(env EnvVars) error {
        authResponse, err := getAuthToken(env.TWITCH_CLIENT_ID, env.TWITCH_SECRET)
		if err != nil {
            return err
		}
		client := &http.Client{}

		dbCount, err := getTotalGamesDb(db)
		if err != nil {
            return err
		}
		apiCount, err := getTotalGamesApi(authResponse.AccessToken, env.TWITCH_CLIENT_ID, client)
		if err != nil {
            return err
		}

		if dbCount == apiCount {
            return nil
		}

		if dbCount > apiCount {
            return err
		}

		missingGamesQt := apiCount - dbCount
		limit := 500
		callsUntilFinish := int(math.Ceil(float64(missingGamesQt) / float64(limit)))

		offset := dbCount
		lastTime := time.Now().Add(-250 * 1e6)

		for i := 0; i < callsUntilFinish; i++ {
			executionDuration := time.Since(lastTime).Milliseconds()
			if executionDuration < 250 {
				difference := 250 - executionDuration

				time.Sleep(time.Duration(difference * 1e6))
			}

			games, err := getGames(limit, offset, authResponse.AccessToken, env.TWITCH_CLIENT_ID, client)
			if err != nil {
                return err
			}
			lastTime = time.Now()

			_, err = db.NamedExec("INSERT INTO games (api_id, name, slug) VALUES (:id, :name, :slug)", games)

			if err != nil {
				log.Printf("Stopped at offset %d\n", offset)
                return err
			}

			offset += 500
		}

        return nil
}

func main() {
	env := loadEnv()

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

    app.Use(logger.New())
    app.Use(recover.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/games/:slug", func(c *fiber.Ctx) error {
		game := DbGame{}
		slug := c.Params("slug")

		err := db.Get(&game, "SELECT * FROM games WHERE games.slug = $1", slug)
		if err != nil {
            return fiber.ErrNotFound
		}


		return c.JSON(game)
	})
	v1.Post("/ratings/", func(c *fiber.Ctx) error {
		type Body struct {
            GameId string `json:"gameId" db:"game_id"`
            UserId string `json:"userId" db:"user_id"`
            Rate   int    `json:"rate" db:"rate"`
		}
		u := new(Body)

		if err := c.BodyParser(u); err != nil {
            return fiber.ErrUnprocessableEntity
		}

        _, err := db.Exec("INSERT INTO ratings (game_id, user_id, rate) VALUES ($1, $2, $3)", u.GameId, u.UserId, u.Rate)
		if err != nil {
            log.Println(err)
            return fiber.ErrInternalServerError
		}

		return c.SendStatus(fiber.StatusCreated)
	})
	v1.Post("/games/import", func(c *fiber.Ctx) error {
        go importGames(env)

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}
