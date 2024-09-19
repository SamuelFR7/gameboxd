package data

import (
	"gameboxd/api/db"
	"gameboxd/api/services"
	"gameboxd/api/types"
)

type Game struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Slug  string `json:"slug" db:"slug"`
	ApiId int    `json:"apiId" db:"api_id"`
}

func GetTotalGames() (int, error) {
	var count int
	err := db.Db.Get(&count, "SELECT count(id) FROM games")
	if err != nil {
		return 0, err
	}

	return count, nil
}

func ListGames(page int, name string) (types.PaginatedResult, error) {
	limit := 10
	offset := (page - 1) * limit

	games := []Game{}
	var count int
    result := types.PaginatedResult{
        Data: &games,
        TotalCount: count,
    }

	err := db.Db.Select(&games, "SELECT * FROM games WHERE games.name ilike '%' || $1 || '%' OFFSET $2 LIMIT $3", name, offset, limit)
	if err != nil {
		return result, err
	}

	err = db.Db.Get(&count, "SELECT count(id) FROM games WHERE games.name ilike '%' || $1 || '%'", name)
	if err != nil {
		return result, err
	}


	return result, nil
}

func GetGameBySlug(slug string) (Game, error) {
    game := Game{}

    err := db.Db.Get(&game, "SELECT * FROM games WHERE games.slug = $1", slug) 
    if err != nil {
        return game, err 
    }

    return game, nil
}


func CreateMultipleGamesFromIgdb(games []services.IgdbGame) error {
    _, err := db.Db.NamedExec("INSERT INTO games (name, api_id, slug) VALUES (:name, :id, :slug)", &games)
    if err != nil {
        return err
    }

    return nil
}
