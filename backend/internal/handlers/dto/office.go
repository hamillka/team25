package dto

import "github.com/hamillka/team25/backend/internal/models"

// CreateOrEditOfficeRequestDto model info
// @Description Информация о кабинете при создании или изменении
type CreateOrEditOfficeRequestDto struct {
	Number int64 `json:"number"` // Номер кабинета
	Floor  int64 `json:"floor"`  // Этаж
}

// CreateOrEditOfficeResponseDto model info
// @Description Информация о кабинете при создании или изменении
type CreateOrEditOfficeResponseDto struct {
	ID int64 `json:"id"` // Идентификатор кабинета
}

// GetOfficeResponseDto model info
// @Description Информация о кабинете при получении
type GetOfficeResponseDto struct {
	ID     int64 `json:"id"`     // Идентификатор кабинета
	Number int64 `json:"number"` // Номер кабинета
	Floor  int64 `json:"floor"`  // Этаж
}

func ConvertToOfficeDto(office *models.Office) *GetOfficeResponseDto {
	return &GetOfficeResponseDto{
		ID:     office.ID,
		Number: office.Number,
		Floor:  office.Floor,
	}
}
