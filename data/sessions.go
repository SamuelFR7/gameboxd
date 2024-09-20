package data

import (
	"gameboxd/db"
	"log"
	"time"
)

type Session struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

const SESSION_EXPIRATION_TIME = time.Hour * 24 * 30

func generateSessionExpiresAt() time.Time {
	now := time.Now()
	expiresAt := now.Add(SESSION_EXPIRATION_TIME)

	return expiresAt
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

func ValidateSession(sessionId string) (Session, error) {
	session := Session{}
	err := db.Db.Get(&session, "SELECT * FROM sessions WHERE sessions.id = $1 AND sessions.expires_at > NOW()", sessionId)
	if err != nil {
		return session, err
	}

	return session, nil
}

func DestroySession(sessionId string) error {
	_, err := db.Db.Exec("DELETE FROM sessions WHERE sessions.id = $1", sessionId)
	if err != nil {
		return err
	}

	return nil
}
