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

type PatientService interface {
	EditPatient(id int64, fio, phoneNumber, email, insurance string) (int64, error)
	AddPatient(fio, phoneNumber, email, insurance string) (int64, error)
	GetAllPatients() ([]models.Patient, error)
	GetPatientByID(id int64) (models.Patient, error)
}

type PatientHandler struct {
	service PatientService
	logger  *zap.SugaredLogger
	port    string
}

func NewPatientHandler(s PatientService, logger *zap.SugaredLogger, port string) *PatientHandler {
	return &PatientHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// AddPatient godoc
//
//	@Summary		Добавить пациента
//	@Description	Добавить пациента в таблицу Patients
//	@ID				add-patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			Patient	body	dto.CreateOrEditPatientRequestDto	true	"Информация о пациенте"
//
// @Success		201	    {object} 	dto.CreateOrEditPatientResponseDto	"Created"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients [post]
func (ph *PatientHandler) AddPatient(w http.ResponseWriter, r *http.Request) {
	ph.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ph.port)

	var patient dto.CreateOrEditPatientRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		ph.logger.Errorf("patient handler: json decode %s", err)
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

	id, err := ph.service.AddPatient(
		patient.Fio,
		patient.PhoneNumber,
		patient.Email,
		patient.Insurance,
	)
	if err != nil {
		ph.logger.Errorf("patient handler: add patient service method: %s", err)
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

	createOrEditPatientResponseDto := dto.CreateOrEditPatientResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrEditPatientResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// EditPatient godoc
//
//	@Summary		Изменить информацию о пациенте
//	@Description	Изменить информацию о пациенте в таблице Patients
//	@ID				edit-patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор пациента"
//	@Param			Patient	body	dto.CreateOrEditPatientRequestDto	true	"Информация о изменяемом пациенте"
//
// @Success		200	    {object} 	dto.CreateOrEditPatientResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients/{id} [patch]
//
//nolint:funlen // it's ok
func (ph *PatientHandler) EditPatient(w http.ResponseWriter, r *http.Request) {
	ph.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ph.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role == dto.DOCTOR {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var patient dto.CreateOrEditPatientRequestDto

	param, ok := mux.Vars(r)["id"]
	patientID, err := strToInt64(param)
	if !ok || err != nil {
		ph.logger.Errorf("patient handler: %s", err)
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
	err = json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		ph.logger.Errorf("patient handler: json decode %s", err)
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

	id, err := ph.service.EditPatient(
		patientID,
		patient.Fio,
		patient.PhoneNumber,
		patient.Email,
		patient.Insurance,
	)
	if err != nil {
		ph.logger.Errorf("patient handler: edit patient service method: %s", err)
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

	createOrUpdatePatientResponseDto := dto.CreateOrEditPatientResponseDto{
		ID: id,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createOrUpdatePatientResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetAllPatients godoc
//
//	@Summary		Получить список пациентов
//	@Description	Получить всех пациентов из таблицы Patients
//	@ID				get-patients
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//
// @Success		200	    {array} 	dto.GetPatientResponseDto	"OK"
// @Failure		401	    {object}	dto.ErrorDto				"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto				"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients [get]
func (ph *PatientHandler) GetAllPatients(w http.ResponseWriter, r *http.Request) {
	ph.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ph.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.ADMIN {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	allPatients, err := ph.service.GetAllPatients()
	if err != nil {
		ph.logger.Errorf("patient handler: get all patients service method: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Error: "Внутренняя ошибка сервера",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	getAllPatientsResponseDto := make([]*dto.GetPatientResponseDto, 0)

	for idx := range allPatients {
		getAllPatientsResponseDto = append(
			getAllPatientsResponseDto,
			dto.ConvertToPatientDto(&allPatients[idx]),
		)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getAllPatientsResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetPatientByID godoc
//
//	@Summary		Получение пациента по идентификатору
//	@Description	Получить пациента по идентификатору в таблице Patients
//	@ID				get-patient-by-id
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор пациента"
//
// @Success		200	    {object} 	dto.GetPatientResponseDto		"Пациент"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto					"Пользователь не имеет доступа"
// @Failure		404	    {object}	dto.ErrorDto					"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/patients/{id} [get]
func (ph *PatientHandler) GetPatientByID(w http.ResponseWriter, r *http.Request) {
	ph.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ph.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role == dto.USER {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	param, ok := mux.Vars(r)["id"]
	patientID, err := strToInt64(param)
	if !ok || err != nil {
		ph.logger.Errorf("patient handler: %s", err)
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

	patient, err := ph.service.GetPatientByID(patientID)
	if err != nil {
		ph.logger.Errorf("patient handler: get patient by id service method: %s", err)
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

	getPatientResponseDto := dto.ConvertToPatientDto(&patient)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getPatientResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
