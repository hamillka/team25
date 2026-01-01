package dto

import "github.com/hamillka/team25/backend/internal/models"

// GetTimetableResponseDto model info
// @Description Информация о расписании при получении
type GetTimetableResponseDto struct {
	WorkDay  int64 `json:"workDay"`  // Номер рабочего дня
	ID       int64 `json:"id"`       // Идентификатор записи в таблице
	DoctorID int64 `json:"doctorId"` // Идентификатор врача
	OfficeID int64 `json:"officeId"` // Идентификатор кабинета
}

func ConvertToTimetableDto(tt *models.Timetable) *GetTimetableResponseDto {
	return &GetTimetableResponseDto{
		WorkDay:  tt.WorkDay,
		ID:       tt.ID,
		DoctorID: tt.DoctorID,
		OfficeID: tt.OfficeID,
	}
}
