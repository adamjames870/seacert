package dto

import (
	"time"
)

type ShipType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VoyageType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PeriodType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Ship struct {
	Id              string    `json:"id"`
	CreatedAt       time.Time `json:"created-at"`
	UpdatedAt       time.Time `json:"updated-at"`
	Name            string    `json:"name"`
	ShipTypeId      string    `json:"ship-type-id"`
	ShipTypeName    string    `json:"ship-type-name"`
	ImoNumber       string    `json:"imo-number"`
	Gt              int32     `json:"gt"`
	Flag            string    `json:"flag"`
	PropulsionPower *int32    `json:"propulsion-power,omitzero"`
}

type Seatime struct {
	Id             string          `json:"id"`
	UserId         string          `json:"user-id"`
	Ship           Ship            `json:"ship"`
	VoyageTypeId   string          `json:"voyage-type-id"`
	VoyageTypeName string          `json:"voyage-type-name"`
	CreatedAt      time.Time       `json:"created-at"`
	UpdatedAt      time.Time       `json:"updated-at"`
	StartDate      time.Time       `json:"start-date"`
	StartLocation  string          `json:"start-location"`
	EndDate        time.Time       `json:"end-date"`
	EndLocation    string          `json:"end-location"`
	TotalDays      int32           `json:"total-days"`
	Company        string          `json:"company"`
	Capacity       string          `json:"capacity"`
	IsWatchkeeping bool            `json:"is-watchkeeping"`
	Periods        []SeatimePeriod `json:"periods"`
}

type SeatimePeriod struct {
	Id             string    `json:"id"`
	SeatimeId      string    `json:"seatime-id"`
	PeriodTypeId   string    `json:"period-type-id"`
	PeriodTypeName string    `json:"period-type-name"`
	StartDate      time.Time `json:"start-date"`
	EndDate        time.Time `json:"end-date"`
	Days           int32     `json:"days"`
	Remarks        string    `json:"remarks"`
}

type ParamsAddSeatime struct {
	UserId         string                   `json:"-"`
	ShipId         *string                  `json:"ship-id,omitempty"`
	Ship           *ParamsAddShip           `json:"ship,omitempty"`
	VoyageTypeId   string                   `json:"voyage-type-id" validate:"required"`
	StartDate      string                   `json:"start-date" validate:"required"`
	StartLocation  string                   `json:"start-location" validate:"required"`
	EndDate        string                   `json:"end-date" validate:"required"`
	EndLocation    string                   `json:"end-location" validate:"required"`
	TotalDays      int32                    `json:"total-days" validate:"required"`
	Company        string                   `json:"company" validate:"required"`
	Capacity       string                   `json:"capacity" validate:"required"`
	IsWatchkeeping bool                     `json:"is-watchkeeping"`
	Periods        []ParamsAddSeatimePeriod `json:"periods,omitempty"`
}

type ParamsAddShip struct {
	Name            string `json:"name" validate:"required"`
	ShipTypeId      string `json:"ship-type-id" validate:"required"`
	ImoNumber       string `json:"imo-number" validate:"required"`
	Gt              int32  `json:"gt" validate:"required"`
	Flag            string `json:"flag" validate:"required"`
	PropulsionPower *int32 `json:"propulsion-power,omitempty"`
}

type ParamsAddSeatimePeriod struct {
	PeriodTypeId string `json:"period-type-id" validate:"required"`
	StartDate    string `json:"start-date" validate:"required"`
	EndDate      string `json:"end-date" validate:"required"`
	Days         int32  `json:"days" validate:"required"`
	Remarks      string `json:"remarks,omitempty"`
}

type SeatimeLookups struct {
	ShipTypes   []ShipType   `json:"ship-types"`
	VoyageTypes []VoyageType `json:"voyage-types"`
	PeriodTypes []PeriodType `json:"period-types"`
}
