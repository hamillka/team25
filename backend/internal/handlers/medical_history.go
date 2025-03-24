package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/hamillka/team25/backend/internal/handlers/dto"
	"github.com/hamillka/team25/backend/internal/handlers/middlewares"
	"github.com/hamillka/team25/backend/internal/models"
	"github.com/hamillka/team25/backend/internal/repositories"
	"go.uber.org/zap"
)

type MedicalHistoryService interface {
	GetHistoryByPatient(id int64) (models.MedicalHistory, error)
	CreateMedicalHistory(chronicDiseases, allergies, bloodType, vaccination string, patientID int64) (int64, error)
	UpdateMedicalHistory(patientID int64, patchData map[string]interface{}) (models.MedicalHistory, error)
}

type MedicalHistoryHandler struct {
	service MedicalHistoryService
	logger  *zap.SugaredLogger
	port    string
}

func NewMedicalHistoryHandler(
	s MedicalHistoryService,
	logger *zap.SugaredLogger,
	port string,
) *MedicalHistoryHandler {
	return &MedicalHistoryHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// GetHistoryByPatient godoc
//
//	@Summary		Получение медкарты пациента по его идентификатору
//	@Description	Получить медкарту пациента в таблице medical_history
//	@ID				get-medical-history-by-patient-id
//	@Tags			medicalHistories
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор пациента"
//
// @Success		200	    {object} 	dto.GetMedicalHistoryResponseDto		"Медицинская карта"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto					"Пользователь не имеет доступа"
// @Failure		404	    {object}	dto.ErrorDto					"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients/{id}/medical_history [get]
func (mhh *MedicalHistoryHandler) GetHistoryByPatient(
	w http.ResponseWriter,
	r *http.Request,
) {
	mhh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", mhh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.DOCTOR && role != dto.USER {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	param, ok := mux.Vars(r)["id"]
	patientID, err := strToInt64(param)
	if !ok || err != nil {
		mhh.logger.Errorf("medical history handler: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	medicalHistory, err := mhh.service.GetHistoryByPatient(patientID)
	if err != nil {
		mhh.logger.Errorf("medical history handler: get history by patient id service method: %s", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, repositories.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			errorDto = &dto.ErrorDto{
				Error: "Запись не найдена",
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			errorDto = &dto.ErrorDto{
				Error: "Внутренняя ошибка сервера",
			}
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	getMedicalHistoryResponseDto := dto.ConvertToMedicalHistoryDto(medicalHistory)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getMedicalHistoryResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CreateMedicalHistory godoc
//
// @Summary		Создать медкарту пациента
// @Description	Создать медкарту в таблице medical_histories
// @ID				create-medical-history
// @Tags			medicalHistories
// @Accept			json
// @Produce			json
//
//	@Param 			id				path	integer		true	"Идентификатор пациента"
//
// @Param			MedicalHistory	body	dto.CreateOrEditMedicalHistoryRequestDto	true	"Информация о медкарте"
//
// @Success		201	{object} 	dto.CreateOrEditMedicalHistoryResponseDto	"Created"
// @Failure		400	{object}	dto.ErrorDto								"Некорректные данные"
// @Failure		401	{object}	dto.ErrorDto								"Пользователь не авторизован"
// @Failure		500	{object}	dto.ErrorDto								"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients/{id}/medical_history [post]
func (mhh *MedicalHistoryHandler) CreateMedicalHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	var history dto.CreateOrEditMedicalHistoryRequestDto
	mhh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", mhh.port)

	param := mux.Vars(r)["id"]
	patientID, err := strToInt64(param)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	id, err := mhh.service.CreateMedicalHistory(
		history.ChronicDiseases,
		history.Allergies,
		history.BloodType,
		history.Vaccination,
		patientID,
	)
	if err != nil {
		mhh.logger.Errorf("medical history handler: create history service method: %s", err)
		var errorDto *dto.ErrorDto
		w.WriteHeader(http.StatusInternalServerError)
		errorDto = &dto.ErrorDto{
			Error: "Внутренняя ошибка сервера",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	createOrEditMedicalHistoryResponseDto := dto.CreateOrEditMedicalHistoryResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrEditMedicalHistoryResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// UpdateMedicalHistory godoc
//
//	@Summary		Изменить медкарту
//	@Description	Изменить медкарту в таблице medical_histories
//	@ID				update-medical-history
//	@Tags			medicalHistories
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор пациента"
//	@Param			medicalHistory	body	dto.CreateOrEditMedicalHistoryRequestDto	true	"Информация о медкарте"
//
// @Success		200	    {object} 	dto.CreateOrEditMedicalHistoryResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients/{id}/medical_history [patch]
//
//nolint:funlen // it's ok
func (mhh *MedicalHistoryHandler) UpdateMedicalHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	mhh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", mhh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.DOCTOR {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	param, ok := mux.Vars(r)["id"]
	patientID, err := strToInt64(param)
	if !ok || err != nil {
		mhh.logger.Errorf("medical history handler: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	var patchData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&patchData)
	if err != nil {
		mhh.logger.Errorf("medical history handler: decode json %s", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Некорректные данные",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	medicalHistory, err := mhh.service.UpdateMedicalHistory(
		patientID,
		patchData,
	)
	if err != nil {
		mhh.logger.Errorf("medical history handler: update medical history service method: %s", err)
		var errorDto *dto.ErrorDto
		w.WriteHeader(http.StatusInternalServerError)
		errorDto = &dto.ErrorDto{
			Error: "Внутренняя ошибка сервера",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	createOrUpdateMedicalHistoryResponseDto := dto.CreateOrEditMedicalHistoryResponseDto{
		ID: medicalHistory.ID,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createOrUpdateMedicalHistoryResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
