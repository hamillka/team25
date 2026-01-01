package dto

import "github.com/hamillka/team25/backend/internal/models"

// CreateOrEditPatientRequestDto model info
// @Description Информация о пациенте при создании или изменении
type CreateOrEditPatientRequestDto struct {
	Fio         string `json:"fio"`         // ФИО пациента
	PhoneNumber string `json:"phoneNumber"` // Телефон пациента
	Email       string `json:"email"`       // Почта пациента
	Insurance   string `json:"insurance"`   // Страховка пациента
}

// CreateOrEditPatientResponseDto model info
// @Description Информация о пациенте при создании или изменении
type CreateOrEditPatientResponseDto struct {
	ID int64 `json:"id"` // Идентификатор пациента
}

// GetPatientResponseDto model info
// @Description Информация о пациенте при получении
type GetPatientResponseDto struct {
	Fio         string `json:"fio"`         // ФИО пациента
	PhoneNumber string `json:"phoneNumber"` // Телефон пациента
	Email       string `json:"email"`       // Почта пациента
	Insurance   string `json:"insurance"`   // Страховка пациента
	ID          int64  `json:"id"`          // Идентификатор пациента
}

func ConvertToPatientDto(patient *models.Patient) *GetPatientResponseDto {
	return &GetPatientResponseDto{
		Fio:         patient.Fio,
		PhoneNumber: patient.PhoneNumber,
		Email:       patient.Email,
		Insurance:   patient.Insurance,
		ID:          patient.ID,
	}
}
