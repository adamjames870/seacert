package seatime

import (
	"github.com/adamjames870/seacert/internal/dto"
)

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
		Ship: dto.Ship{
			Id:              st.Ship.Id.String(),
			CreatedAt:       st.Ship.CreatedAt,
			UpdatedAt:       st.Ship.UpdatedAt,
			Name:            st.Ship.Name,
			ShipTypeId:      st.Ship.ShipTypeId.String(),
			ShipTypeName:    st.Ship.ShipTypeName,
			ImoNumber:       st.Ship.ImoNumber,
			Gt:              st.Ship.Gt,
			Flag:            st.Ship.Flag,
			PropulsionPower: st.Ship.PropulsionPower,
		},
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
