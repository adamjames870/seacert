package seatime

import (
	"time"

	"github.com/google/uuid"
)

type ShipType struct {
	Id          uuid.UUID
	Name        string
	Description string
}

type VoyageType struct {
	Id          uuid.UUID
	Name        string
	Description string
}

type PeriodType struct {
	Id          uuid.UUID
	Name        string
	Description string
}

type Ship struct {
	Id              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	ShipTypeId      uuid.UUID
	ShipTypeName    string
	ImoNumber       string
	Gt              int32
	Flag            string
	PropulsionPower *int32
}

type Seatime struct {
	Id             uuid.UUID
	UserId         uuid.UUID
	ShipId         uuid.UUID
	Ship           Ship
	VoyageTypeId   uuid.UUID
	VoyageTypeName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	StartDate      time.Time
	StartLocation  string
	EndDate        time.Time
	EndLocation    string
	TotalDays      int32
	Company        string
	Capacity       string
	IsWatchkeeping bool
	Periods        []SeatimePeriod
}

type SeatimePeriod struct {
	Id           uuid.UUID
	SeatimeId    uuid.UUID
	PeriodTypeId uuid.UUID
	PeriodType   string
	StartDate    time.Time
	EndDate      time.Time
	Days         int32
	Remarks      string
}
