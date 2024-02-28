package models

import (
	"encoding/json"
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
