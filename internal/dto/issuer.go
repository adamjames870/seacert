package dto

import "time"

type Issuer struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Website   string    `json:"website"`
}

type ParamsAddIssuer struct {
	Name    string  `json:"name"`
	Country *string `json:"country,omitempty"`
	Website *string `json:"website,omitempty"`
}
