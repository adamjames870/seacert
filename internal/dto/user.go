package dto

import "time"

type User struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created-at"`
	UpdatedAt   time.Time `json:"updated-at"`
	Forename    string    `json:"forename"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	Nationality string    `json:"nationality"`
}

type ParamsAddUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}
