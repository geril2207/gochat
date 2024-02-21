package users

import (
	"context"
	"encoding/json"

	"github.com/geril2207/gochat/packages/db"
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (user User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id    int    `json:"id"`
		Login string `json:"login"`
	}{Id: user.Id, Login: user.Login})
}

func GetUserById(id int) (User, error) {
	var user User

	err := db.Pool.
		QueryRow(context.Background(), "select id,login,password from users where id=$1", id).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}

func GetUserByLogin(login string) (User, error) {
	var user User

	err := db.Pool.
		QueryRow(context.Background(), "select id,login,password from users where login=$1", login).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}

func IsUserExistsInDb(login string) (bool, error) {
	var result int
	err := db.Pool.QueryRow(context.Background(), "select 1 from users where login=$1", login).
		Scan(&result)

	return result == 1, err
}

func InsertUser(login, password string) (User, error) {
	var user User
	err := db.Pool.QueryRow(context.Background(),
		"insert into users(login,password) values($1,$2) RETURNING id,login,password", login, password).
		Scan(&user.Id, &user.Login, &user.Password)

	return user, err
}
