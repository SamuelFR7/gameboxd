package data

import (
	"gameboxd/db"
	"time"
)

type Session struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

const SESSION_EXPIRATION_TIME = 1000 * 60 * 60 * 24 * 30

func generateSessionExpiresAt() time.Time {
	now := time.Now()
	now.Add(time.Duration(SESSION_EXPIRATION_TIME * 1e6))

	return now
}

func CreateSession(userId string) (Session, error) {
	session := Session{}

	lastInsertedId := ""
	err := db.Db.QueryRow("INSERT INTO sessions (user_id, expires_at) VALUES ($1, $2) RETURNING id", userId, generateSessionExpiresAt()).Scan(&lastInsertedId)
	if err != nil {
		return session, err
	}

	err = db.Db.Get(&session, "SELECT * FROM sessions WHERE sessions.id = $1", lastInsertedId)
	if err != nil {
		return session, err
	}

	return session, nil
}
