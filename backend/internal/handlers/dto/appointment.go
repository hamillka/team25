package dto

import (
	"time"

	"github.com/hamillka/team25/backend/internal/models"
)

// CreateOrEditAppointmentRequestDto model info
// @Description Информация о записи при создании или изменении
type CreateOrEditAppointmentRequestDto struct {
	DateTime  time.Time `json:"dateTime"`  // Дата и время записи
	ID        int64     `json:"id"`        // Идентификатор записи
	PatientID int64     `json:"patientId"` // Идентификатор пациента
	DoctorID  int64     `json:"doctorId"`  // Идентификатор врача
}

// CreateOrEditAppointmentResponseDto model info
// @Description Информация о записи при создании или изменении
type CreateOrEditAppointmentResponseDto struct {
	ID int64 `json:"id"` // Идентификатор записи
}

// CancelAppointmentRequestDto model info
// @Description Информация о записи при удалении
type CancelAppointmentRequestDto struct {
	ID int64 `json:"id"` // Идентификатор записи
}

// GetAppointmentResponseDto model info
// @Description Информация о записи при получении
type GetAppointmentResponseDto struct {
	DateTime  time.Time `json:"dateTime"`  // Дата и время записи
	ID        int64     `json:"id"`        // Идентификатор записи
	PatientID int64     `json:"patientId"` // Идентификатор пациента
	DoctorID  int64     `json:"doctorId"`  // Идентификатор врача
}

func ConvertToDto(appointment *models.Appointment) *GetAppointmentResponseDto {
	return &GetAppointmentResponseDto{
		DateTime:  appointment.DateTime,
		ID:        appointment.ID,
		PatientID: appointment.PatientID,
		DoctorID:  appointment.DoctorID,
	}
}
