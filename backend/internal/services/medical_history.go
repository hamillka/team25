package services

import (
	"github.com/hamillka/team25/backend/internal/models"
)

type MedicalHistoryRepository interface {
	GetHistoryByPatient(id int64) (models.MedicalHistory, error)
	CreateMedicalHistory(chronicDiseases, allergies, bloodType, vaccination string, patientID int64) (int64, error)
	//UpdateMedicalHistory(id int64, chronicDiseases, allergies, bloodType, vaccination string) (int64, error)
	UpdateMedicalHistory(medicalHistory models.MedicalHistory) error
}

type MedicalHistoryService struct {
	repo MedicalHistoryRepository
}

func NewMedicalHistoryService(repository MedicalHistoryRepository) *MedicalHistoryService {
	return &MedicalHistoryService{repo: repository}
}

func (mhs *MedicalHistoryService) GetHistoryByPatient(
	id int64,
) (models.MedicalHistory, error) {
	history, err := mhs.repo.GetHistoryByPatient(id)
	if err != nil {
		return models.MedicalHistory{}, err
	}

	return history, nil
}

func (mhs *MedicalHistoryService) CreateMedicalHistory(
	chronicDiseases,
	allergies,
	bloodType,
	vaccination string,
	patientID int64,
) (int64, error) {
	id, err := mhs.repo.CreateMedicalHistory(chronicDiseases, allergies, bloodType, vaccination, patientID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//
//func (mhs *MedicalHistoryService) UpdateMedicalHistory(
//	id int64,
//	chronicDiseases,
//	allergies,
//	bloodType,
//	vaccination string,
//) (int64, error) {
//	newID, err := mhs.repo.UpdateMedicalHistory(id, chronicDiseases, allergies, bloodType, vaccination)
//	if err != nil {
//		return 0, err
//	}
//
//	return newID, nil
//}
//

func (mhs *MedicalHistoryService) UpdateMedicalHistory(patientID int64, patchData map[string]interface{}) (models.MedicalHistory, error) {
	medicalHistory, err := mhs.repo.GetHistoryByPatient(patientID)
	if err != nil {
		return models.MedicalHistory{}, err
	}

	if chronicDiseases, ok := patchData["chronicDiseases"].(string); ok {
		*medicalHistory.ChronicDiseases = chronicDiseases
	}
	if allergies, ok := patchData["allergies"].(string); ok {
		*medicalHistory.Allergies = allergies
	}
	if bloodType, ok := patchData["bloodType"].(string); ok {
		*medicalHistory.BloodType = bloodType
	}
	if vaccination, ok := patchData["vaccination"].(string); ok {
		*medicalHistory.Vaccination = vaccination
	}

	err = mhs.repo.UpdateMedicalHistory(medicalHistory)
	if err != nil {
		return models.MedicalHistory{}, err
	}

	return medicalHistory, nil
}
