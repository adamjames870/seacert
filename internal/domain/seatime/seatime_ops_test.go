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
		ID:              arg.ID,
		Name:            arg.Name,
		ImoNumber:       arg.ImoNumber,
		Gt:              arg.Gt,
		Flag:            arg.Flag,
		Status:          arg.Status,
		CreatedBy:       arg.CreatedBy,
		ShipTypeID:      arg.ShipTypeID,
		PropulsionPower: arg.PropulsionPower,
	}, nil
}

func (m *mockRepo) CreateSeatime(ctx context.Context, arg sqlc.CreateSeatimeParams) (sqlc.Seatime, error) {
	return sqlc.Seatime{ID: arg.ID}, nil
}

func (m *mockRepo) CreateSeatimePeriod(ctx context.Context, arg sqlc.CreateSeatimePeriodParams) (sqlc.SeatimePeriod, error) {
	return sqlc.SeatimePeriod{ID: arg.ID}, nil
}

func (m *mockRepo) UpdateSeatime(ctx context.Context, arg sqlc.UpdateSeatimeParams) (sqlc.Seatime, error) {
	return sqlc.Seatime{ID: arg.ID}, nil
}

func (m *mockRepo) DeleteSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) error {
	return nil
}

func (m *mockRepo) GetShips(ctx context.Context) ([]sqlc.GetShipsRow, error) {
	return []sqlc.GetShipsRow{
		{
			ID:           uuid.New(),
			Name:         "Admin Ship 1",
			ShipTypeName: "Tanker",
			Status:       "approved",
		},
		{
			ID:           uuid.New(),
			Name:         "Admin Ship 2",
			ShipTypeName: "Bulk",
			Status:       "provisional",
		},
	}, nil
}

func (m *mockRepo) GetShipsForUser(ctx context.Context, createdBy uuid.NullUUID) ([]sqlc.GetShipsForUserRow, error) {
	return []sqlc.GetShipsForUserRow{
		{
			ID:           uuid.New(),
			Name:         "Approved Ship",
			ShipTypeName: "Tanker",
			Status:       "approved",
		},
		{
			ID:           uuid.New(),
			Name:         "My Provisional Ship",
			ShipTypeName: "Bulk",
			Status:       "provisional",
			CreatedBy:    createdBy,
		},
	}, nil
}

func (m *mockRepo) UpdateShipReferences(ctx context.Context, arg sqlc.UpdateShipReferencesParams) error {
	return nil
}

func (m *mockRepo) UpdateShip(ctx context.Context, arg sqlc.UpdateShipParams) (sqlc.Ship, error) {
	return sqlc.Ship{
		ID:              arg.ID,
		Name:            arg.Name,
		ShipTypeID:      arg.ShipTypeID,
		ImoNumber:       arg.ImoNumber,
		Gt:              arg.Gt,
		Flag:            arg.Flag,
		PropulsionPower: arg.PropulsionPower,
	}, nil
}

func (m *mockRepo) UpdateShipStatus(ctx context.Context, arg sqlc.UpdateShipStatusParams) (sqlc.Ship, error) {
	return sqlc.Ship{
		ID:     arg.ID,
		Status: arg.Status,
	}, nil
}

func (m *mockRepo) DeleteShip(ctx context.Context, id uuid.UUID) error {
	return nil
}

var mockGetShipById func(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error)

func (m *mockRepo) GetShipById(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
	if mockGetShipById != nil {
		return mockGetShipById(ctx, id)
	}
	return sqlc.GetShipByIdRow{
		ID:           id,
		Name:         "Mock Ship",
		ShipTypeName: "Mock Type",
		Status:       "approved",
	}, nil
}

func (m *mockRepo) GetVoyageTypes(ctx context.Context) ([]sqlc.VoyageType, error) {
	return nil, nil
}

func (m *mockRepo) GetSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) ([]sqlc.GetSeatimePeriodsRow, error) {
	return nil, nil
}

func (m *mockRepo) GetSeatimeByUserId(ctx context.Context, userID uuid.UUID) ([]sqlc.GetSeatimeByUserIdRow, error) {
	if mockGetSeatimeByUserId != nil {
		return mockGetSeatimeByUserId(ctx, userID)
	}
	return nil, nil
}

var mockGetSeatimeByUserId func(ctx context.Context, userID uuid.UUID) ([]sqlc.GetSeatimeByUserIdRow, error)

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
			res, err := CreateSeatime(context.Background(), repo, tt.params, userId, false)
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

func TestGetShips(t *testing.T) {
	repo := &mockRepo{}
	userId := uuid.New()

	t.Run("Admin sees all ships", func(t *testing.T) {
		ships, err := GetShips(context.Background(), repo, &userId, true)
		if err != nil {
			t.Fatalf("GetShips failed: %v", err)
		}
		if len(ships) != 2 {
			t.Errorf("Expected 2 ships for admin, got %d", len(ships))
		}
	})

	t.Run("User sees filtered ships", func(t *testing.T) {
		ships, err := GetShips(context.Background(), repo, &userId, false)
		if err != nil {
			t.Fatalf("GetShips failed: %v", err)
		}
		if len(ships) != 2 {
			t.Errorf("Expected 2 ships for user, got %d", len(ships))
		}
		for _, s := range ships {
			if s.Status == "provisional" && (s.CreatedBy == nil || *s.CreatedBy != userId) {
				t.Errorf("User should not see others' provisional ships")
			}
		}
	})
}

func TestResolveShip(t *testing.T) {
	repo := &mockRepo{}

	t.Run("Valid resolve", func(t *testing.T) {
		params := dto.ParamsResolveShip{
			ProvisionalId: uuid.New().String(),
			ReplacementId: uuid.New().String(),
		}
		err := ResolveShip(context.Background(), repo, params)
		if err != nil {
			t.Errorf("ResolveShip failed: %v", err)
		}
	})

	t.Run("Invalid UUIDs", func(t *testing.T) {
		params := dto.ParamsResolveShip{
			ProvisionalId: "invalid",
			ReplacementId: uuid.New().String(),
		}
		err := ResolveShip(context.Background(), repo, params)
		if err == nil {
			t.Errorf("ResolveShip should fail for invalid UUID")
		}
	})
}

func TestUpdateSeatimeValidation(t *testing.T) {
	repo := &mockRepo{}
	userId := uuid.New()
	seatimeId := uuid.New().String()
	shipId := uuid.New().String()
	voyageTypeId := uuid.New().String()

	tests := []struct {
		name    string
		params  dto.ParamsUpdateSeatime
		wantErr bool
	}{
		{
			name: "Valid update",
			params: dto.ParamsUpdateSeatime{
				Id:             seatimeId,
				ShipId:         &shipId,
				VoyageTypeId:   voyageTypeId,
				StartDate:      "2024-01-01",
				EndDate:        "2024-01-10",
				StartLocation:  "London",
				EndLocation:    "New York",
				TotalDays:      10,
				Company:        "Updated Co",
				Capacity:       "Chief Mate",
				IsWatchkeeping: true,
			},
			wantErr: false,
		},
		{
			name: "Invalid ID",
			params: dto.ParamsUpdateSeatime{
				Id: "invalid-uuid",
			},
			wantErr: true,
		},
		{
			name: "Period outside voyage dates",
			params: dto.ParamsUpdateSeatime{
				Id:            seatimeId,
				ShipId:        &shipId,
				VoyageTypeId:  voyageTypeId,
				StartDate:     "2024-01-01",
				EndDate:       "2024-01-10",
				StartLocation: "London",
				EndLocation:   "New York",
				TotalDays:     10,
				Company:       "Updated Co",
				Capacity:      "Chief Mate",
				Periods: []dto.ParamsAddSeatimePeriod{
					{
						PeriodTypeId: uuid.New().String(),
						StartDate:    "2023-12-31",
						EndDate:      "2024-01-05",
						Days:         5,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockGetSeatimeByUserId = func(ctx context.Context, userID uuid.UUID) ([]sqlc.GetSeatimeByUserIdRow, error) {
					id, _ := uuid.Parse(tt.params.Id)
					var sId uuid.UUID
					if tt.params.ShipId != nil {
						sId, _ = uuid.Parse(*tt.params.ShipId)
					}
					voyageTypeId, _ := uuid.Parse(tt.params.VoyageTypeId)
					return []sqlc.GetSeatimeByUserIdRow{
						{
							ID:           id,
							UserID:       userID,
							ShipID:       sId,
							VoyageTypeID: voyageTypeId,
							StartDate:    time.Now(),
							EndDate:      time.Now(),
						},
					}, nil
				}
			}
			_, err := UpdateSeatime(context.Background(), repo, tt.params, userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSeatime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStandaloneShipOps(t *testing.T) {
	repo := &mockRepo{}
	userId := uuid.New()
	shipTypeId := uuid.New().String()

	t.Run("CreateShipStandalone as User", func(t *testing.T) {
		mockGetShipById = func(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
			return sqlc.GetShipByIdRow{
				ID:     id,
				Status: "provisional",
				Name:   "Test Ship",
			}, nil
		}
		params := dto.ParamsAddShip{
			Name:       "Test Ship",
			ShipTypeId: shipTypeId,
			ImoNumber:  "IMO123",
			Gt:         100,
			Flag:       "UK",
		}
		s, err := CreateShipStandalone(context.Background(), repo, params, userId, false)
		if err != nil {
			t.Fatalf("CreateShipStandalone failed: %v", err)
		}
		if s.Status != "provisional" {
			t.Errorf("Expected provisional status, got %s", s.Status)
		}
	})

	t.Run("CreateShipStandalone as Admin", func(t *testing.T) {
		mockGetShipById = func(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
			return sqlc.GetShipByIdRow{
				ID:     id,
				Status: "approved",
				Name:   "Admin Ship",
			}, nil
		}
		params := dto.ParamsAddShip{
			Name:       "Admin Ship",
			ShipTypeId: shipTypeId,
			ImoNumber:  "IMO456",
			Gt:         100,
			Flag:       "UK",
		}
		s, err := CreateShipStandalone(context.Background(), repo, params, userId, true)
		if err != nil {
			t.Fatalf("CreateShipStandalone failed: %v", err)
		}
		if s.Status != "approved" {
			t.Errorf("Expected approved status, got %s", s.Status)
		}
	})

	t.Run("UpdateShip as Owner", func(t *testing.T) {
		shipId := uuid.New()
		mockGetShipById = func(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
			return sqlc.GetShipByIdRow{
				ID:         id,
				Status:     "provisional",
				Name:       "Updated Name",
				CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
				ShipTypeID: uuid.MustParse(shipTypeId),
			}, nil
		}
		params := dto.ParamsUpdateShip{
			Id:         shipId.String(),
			Name:       "Updated Name",
			ShipTypeId: shipTypeId,
		}
		s, err := UpdateShip(context.Background(), repo, params, userId, false)
		if err != nil {
			t.Fatalf("UpdateShip failed: %v", err)
		}
		if s.Name != "Updated Name" {
			t.Errorf("Expected name Updated Name, got %s", s.Name)
		}
	})

	t.Run("UpdateShip Forbidden for Approved", func(t *testing.T) {
		shipId := uuid.New()
		mockGetShipById = func(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
			return sqlc.GetShipByIdRow{
				ID:         id,
				Status:     "approved",
				CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
				ShipTypeID: uuid.MustParse(shipTypeId),
			}, nil
		}
		params := dto.ParamsUpdateShip{
			Id:         shipId.String(),
			Name:       "Updated Name",
			ShipTypeId: shipTypeId,
		}
		_, err := UpdateShip(context.Background(), repo, params, userId, false)
		if err != domain.ErrForbidden {
			t.Errorf("Expected ErrForbidden, got %v", err)
		}
	})

	t.Run("ApproveShip", func(t *testing.T) {
		shipId := uuid.New()
		err := ApproveShip(context.Background(), repo, shipId)
		if err != nil {
			t.Errorf("ApproveShip failed: %v", err)
		}
	})
}
