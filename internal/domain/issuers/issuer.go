package issuers

import (
	"time"

	"github.com/google/uuid"
)

type Issuer struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Website   string    `json:"website"`
}
