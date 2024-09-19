package data

import "gameboxd/db"

type Rating struct {
	Id     string `json:"id" db:"id"`
	GameId string `json:"gameId" db:"game_id"`
	UserId string `json:"userId" db:"user_id"`
	Rate   int    `json:"rate" db:"rate"`
}

func (r *Rating) CreateRating() error {
	_, err := db.Db.Exec("INSERT INTO ratings (game_id, user_id, rate) VALUES ($1, $2, $3)", r.GameId, r.UserId, r.Rate)
	if err != nil {
		return err
	}

	return nil
}
