package dto

import (
	"github.com/hamillka/team25/backend/internal/models"
)

// GetMedicalHistoryResponseDto model info
// @Description Информация о медкарте при получении
type GetMedicalHistoryResponseDto struct {
	ChronicDiseases string `json:"chronicDiseases"` // Хронические заболевания
	Allergies       string `json:"allergies"`       // Аллергии
	BloodType       string `json:"bloodType"`       // Группа крови
	Vaccination     string `json:"vaccination"`     // Вакцинации
	ID              int64  `json:"id"`              // Идентификатор медкарты
	PatientID       int64  `json:"patientId"`       // Идентификатор пациента
}

// CreateOrEditMedicalHistoryRequestDto model info
// @Description Информация о медкарте при создании или изменении
type CreateOrEditMedicalHistoryRequestDto struct {
	ChronicDiseases string `json:"chronicDiseases"` // Хронические заболевания
	Allergies       string `json:"allergies"`       // Аллергии
	BloodType       string `json:"bloodType"`       // Группа крови
	Vaccination     string `json:"vaccination"`     // Вакцинации
}

// CreateOrEditMedicalHistoryResponseDto model info
// @Description Информация о медкарте при создании или изменении
type CreateOrEditMedicalHistoryResponseDto struct {
	ID int64 `json:"id"` // Идентификатор медкарты
}

func ConvertToMedicalHistoryDto(history models.MedicalHistory) *GetMedicalHistoryResponseDto {
	resp := &GetMedicalHistoryResponseDto{
		ID:        history.ID,
		PatientID: history.PatientID,
	}
	if history.ChronicDiseases == nil {
		resp.ChronicDiseases = ""
	} else {
		resp.ChronicDiseases = *history.ChronicDiseases
	}

	if history.Allergies == nil {
		resp.Allergies = ""
	} else {
		resp.Allergies = *history.Allergies
	}

	if history.BloodType == nil {
		resp.BloodType = ""
	} else {
		resp.BloodType = *history.BloodType
	}

	if history.Vaccination == nil {
		resp.Vaccination = ""
	} else {
		resp.Vaccination = *history.Vaccination
	}

	return resp
}
