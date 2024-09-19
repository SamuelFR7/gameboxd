package data

import "gameboxd/db"

type User struct {
	Id           string `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password" db:"password_hash"`
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(user CreateUserParams) error {
	_, err := db.Db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (User, error) {
	user := User{}

	err := db.Db.Get(&user, "SELECT * FROM users WHERE users.email = $1", email)
	if err != nil {
		return user, err
	}

	return user, nil
}
