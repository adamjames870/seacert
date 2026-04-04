package seatime

import (
	"context"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestCalculateDays(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int32
	}{
		{
			name:     "Same day",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "10 days in January",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			expected: 10,
		},
		{
			name:     "Across month boundary (Jan to Feb)",
			start:    time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: 4, // 30, 31, 1, 2
		},
		{
			name:     "Leap year (Feb 2024)",
			start:    time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 3, // 28, 29, 1
		},
		{
			name:     "Non-leap year (Feb 2023)",
			start:    time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 2, // 28, 1
		},
		{
			name:     "Start after end",
			start:    time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDays(tt.start, tt.end)
			if got != tt.expected {
				t.Errorf("CalculateDays(%v, %v) = %d; want %d", tt.start, tt.end, got, tt.expected)
			}
		})
	}
}

// Mock Repository for testing validation
type mockRepo struct {
	domain.Repository
}

func (m *mockRepo) WithTx(ctx context.Context, fn func(domain.Repository) error) error {
	return fn(m)
}

func (m *mockRepo) CreateShip(ctx context.Context, arg sqlc.CreateShipParams) (sqlc.Ship, error) {
	return sqlc.Ship{
		ID:        arg.ID,
		Name:      arg.Name,
		ImoNumber: arg.ImoNumber,
		Gt:        arg.Gt,
		Flag:      arg.Flag,
	}, nil
}

func (m *mockRepo) CreateSeatime(ctx context.Context, arg sqlc.CreateSeatimeParams) (sqlc.Seatime, error) {
	return sqlc.Seatime{ID: arg.ID}, nil
}

func (m *mockRepo) CreateSeatimePeriod(ctx context.Context, arg sqlc.CreateSeatimePeriodParams) (sqlc.SeatimePeriod, error) {
	return sqlc.SeatimePeriod{ID: arg.ID}, nil
}

func (m *mockRepo) GetShipById(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
	return sqlc.GetShipByIdRow{ID: id}, nil
}

func (m *mockRepo) GetVoyageTypes(ctx context.Context) ([]sqlc.VoyageType, error) {
	return nil, nil
}

func (m *mockRepo) GetSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) ([]sqlc.GetSeatimePeriodsRow, error) {
	return nil, nil
}

func TestCreateSeatimeValidation(t *testing.T) {
	repo := &mockRepo{}
	userId := uuid.New()
	voyageTypeId := uuid.New().String()
	shipTypeId := uuid.New().String()

	tests := []struct {
		name    string
		params  dto.ParamsAddSeatime
		wantErr bool
	}{
		{
			name: "Valid params with ship details",
			params: dto.ParamsAddSeatime{
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				VoyageTypeId:  voyageTypeId,
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Global Ships",
				Capacity:      "Master",
				Ship: &dto.ParamsAddShip{
					Name:       "Ocean Voyager",
					ShipTypeId: shipTypeId,
					ImoNumber:  "IMO1234567",
					Gt:         50000,
					Flag:       "UK",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid start date format",
			params: dto.ParamsAddSeatime{
				StartDate: "01-01-2024",
				EndDate:   "2024-01-10",
			},
			wantErr: true,
		},
		{
			name: "Start date after end date",
			params: dto.ParamsAddSeatime{
				StartDate: "2024-01-10",
				EndDate:   "2024-01-01",
			},
			wantErr: true,
		},
		{
			name: "Missing ship details and ship ID",
			params: dto.ParamsAddSeatime{
				StartDate:    "2024-01-01",
				EndDate:      "2024-01-10",
				VoyageTypeId: voyageTypeId,
			},
			wantErr: true,
		},
		{
			name: "Period outside voyage dates",
			params: dto.ParamsAddSeatime{
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				VoyageTypeId:  voyageTypeId,
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Global Ships",
				Capacity:      "Master",
				Ship: &dto.ParamsAddShip{
					Name:       "Ocean Voyager",
					ShipTypeId: shipTypeId,
					ImoNumber:  "IMO1234567",
					Gt:         50000,
					Flag:       "UK",
				},
				Periods: []dto.ParamsAddSeatimePeriod{
					{
						PeriodTypeId: uuid.New().String(),
						StartDate:    "2023-12-31", // Before voyage start
						EndDate:      "2024-01-05",
						Days:         5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Period after voyage dates",
			params: dto.ParamsAddSeatime{
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				VoyageTypeId:  voyageTypeId,
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Global Ships",
				Capacity:      "Master",
				Ship: &dto.ParamsAddShip{
					Name:       "Ocean Voyager",
					ShipTypeId: shipTypeId,
					ImoNumber:  "IMO1234567",
					Gt:         50000,
					Flag:       "UK",
				},
				Periods: []dto.ParamsAddSeatimePeriod{
					{
						PeriodTypeId: uuid.New().String(),
						StartDate:    "2024-01-05",
						EndDate:      "2024-01-11", // After voyage end
						Days:         7,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid params with existing ship ID",
			params: dto.ParamsAddSeatime{
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				VoyageTypeId:  voyageTypeId,
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Global Ships",
				Capacity:      "Master",
				ShipId:        &[]string{uuid.New().String()}[0],
			},
			wantErr: false,
		},
		{
			name: "Valid params with periods",
			params: dto.ParamsAddSeatime{
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				VoyageTypeId:  voyageTypeId,
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Global Ships",
				Capacity:      "Master",
				Ship: &dto.ParamsAddShip{
					Name:       "Ocean Voyager",
					ShipTypeId: shipTypeId,
					ImoNumber:  "IMO1234567",
					Gt:         50000,
					Flag:       "UK",
				},
				Periods: []dto.ParamsAddSeatimePeriod{
					{
						PeriodTypeId: uuid.New().String(),
						StartDate:    "2024-01-02",
						EndDate:      "2024-01-05",
						Days:         4,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CreateSeatime(context.Background(), repo, tt.params, userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSeatime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !tt.wantErr {
				if tt.params.Ship != nil {
					if res.Ship.Name != tt.params.Ship.Name {
						t.Errorf("Expected ship name %s, got %s", tt.params.Ship.Name, res.Ship.Name)
					}
				}
				if tt.params.ShipId != nil {
					if res.Ship.Id.String() != *tt.params.ShipId {
						t.Errorf("Expected ship id %s, got %s", *tt.params.ShipId, res.Ship.Id.String())
					}
				}
			}
		})
	}
}
