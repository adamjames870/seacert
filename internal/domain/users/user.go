package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Forename    string
	Surname     string
	Email       string
	Nationality string
	Role        string
}
