package seatime

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
)

func MapShipToDto(s Ship) dto.Ship {
	return dto.Ship{
		Id:              s.Id.String(),
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		Name:            s.Name,
		ShipTypeId:      s.ShipTypeId.String(),
		ShipTypeName:    s.ShipTypeName,
		ImoNumber:       s.ImoNumber,
		Gt:              s.Gt,
		Flag:            s.Flag,
		PropulsionPower: s.PropulsionPower,
		Status:          s.Status,
		CreatedBy: func() string {
			if s.CreatedBy != nil {
				return s.CreatedBy.String()
			}
			return ""
		}(),
	}
}

func MapSeatimeDomainToDto(st Seatime) dto.Seatime {
	rv := dto.Seatime{
		Id:             st.Id.String(),
		UserId:         st.UserId.String(),
		VoyageTypeId:   st.VoyageTypeId.String(),
		VoyageTypeName: st.VoyageTypeName,
		CreatedAt:      st.CreatedAt,
		UpdatedAt:      st.UpdatedAt,
		StartDate:      st.StartDate,
		StartLocation:  st.StartLocation,
		EndDate:        st.EndDate,
		EndLocation:    st.EndLocation,
		TotalDays:      st.TotalDays,
		Company:        st.Company,
		Capacity:       st.Capacity,
		IsWatchkeeping: st.IsWatchkeeping,
		Ship:           MapShipToDto(st.Ship),
	}

	for _, p := range st.Periods {
		rv.Periods = append(rv.Periods, dto.SeatimePeriod{
			Id:             p.Id.String(),
			SeatimeId:      p.SeatimeId.String(),
			PeriodTypeId:   p.PeriodTypeId.String(),
			PeriodTypeName: p.PeriodType,
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			Days:           p.Days,
			Remarks:        p.Remarks,
		})
	}

	return rv
}

func MapShipType(st sqlc.ShipType) dto.ShipType {
	return dto.ShipType{
		Id:          st.ID.String(),
		Name:        st.Name,
		Description: domain.FromNullString(st.Description),
	}
}

func MapVoyageType(vt sqlc.VoyageType) dto.VoyageType {
	return dto.VoyageType{
		Id:          vt.ID.String(),
		Name:        vt.Name,
		Description: domain.FromNullString(vt.Description),
	}
}

func MapPeriodType(pt sqlc.SeatimePeriodType) dto.PeriodType {
	return dto.PeriodType{
		Id:          pt.ID.String(),
		Name:        pt.Name,
		Description: domain.FromNullString(pt.Description),
	}
}
