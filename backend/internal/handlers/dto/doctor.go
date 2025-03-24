package dto

import "github.com/hamillka/team25/backend/internal/models"

// CreateOrEditDoctorRequestDto model info
// @Description Информация о враче при создании или изменении
type CreateOrEditDoctorRequestDto struct {
	Fio            string `json:"fio"`            // ФИО врача
	PhoneNumber    string `json:"phoneNumber"`    // Телефон врача
	Email          string `json:"email"`          // Почта врача
	Specialization string `json:"specialization"` // Специализация врача
}

// CreateOrEditDoctorResponseDto model info
// @Description Информация о враче при создании или изменении
type CreateOrEditDoctorResponseDto struct {
	ID int64 `json:"id"` // Идентификатор врача
}

// GetDoctorResponseDto model info
// @Description Информация о враче при получении
type GetDoctorResponseDto struct {
	Fio            string `json:"fio"`            // ФИО врача
	PhoneNumber    string `json:"phoneNumber"`    // Телефон врача
	Email          string `json:"email"`          // Почта врача
	ID             int64  `json:"id"`             // Идентификатор врача
	Specialization string `json:"specialization"` // Специализация врача
}

func ConvertToDoctorDto(doctor *models.Doctor) *GetDoctorResponseDto {
	return &GetDoctorResponseDto{
		Fio:            doctor.Fio,
		PhoneNumber:    doctor.PhoneNumber,
		Email:          doctor.Email,
		ID:             doctor.ID,
		Specialization: doctor.Specialization,
	}
}
