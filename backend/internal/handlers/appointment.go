package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/hamillka/team25/backend/internal/handlers/dto"
	"github.com/hamillka/team25/backend/internal/handlers/middlewares"
	"github.com/hamillka/team25/backend/internal/models"
	"github.com/hamillka/team25/backend/internal/repositories"
	"go.uber.org/zap"
)

var ErrValidate = errors.New("validation error")

type AppointmentService interface {
	CreateAppointment(patientID, doctorID int64, dateTime time.Time) (int64, error)
	CancelAppointment(id int64) error
	GetAppointmentsByPatient(id int64) ([]*models.Appointment, error)
	GetAppointmentsByDoctor(id int64) ([]*models.Appointment, error)
	EditAppointment(id, doctorID, patientID int64, dateTime time.Time) (int64, error)
	GetAppointmentByID(id int64) (*models.Appointment, error)
	GetAppointmentsByPatientAndDoctor(patientID, doctorID int64) ([]*models.Appointment, error)
	GetAllAppointments() ([]*models.Appointment, error)
}

type AppointmentHandler struct {
	service AppointmentService
	logger  *zap.SugaredLogger
	port    string
}

func NewAppointmentHandler(s AppointmentService, logger *zap.SugaredLogger, port string) *AppointmentHandler {
	return &AppointmentHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// CreateAppointment godoc
//
//	@Summary		Создать запись
//	@Description	Создать запись пациента к врачу
//	@ID				create-appointment
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param			Appointment	body	dto.CreateOrEditAppointmentRequestDto	true	"Информация о создаваемой записи"
//
// @Success		201	    {object} 	dto.CreateOrEditAppointmentResponseDto	"Created"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/appointments [post]
func (ah *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	ah.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ah.port)

	var appointment dto.CreateOrEditAppointmentRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&appointment)
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

	id, err := ah.service.CreateAppointment(
		appointment.PatientID,
		appointment.DoctorID,
		appointment.DateTime,
	)
	if err != nil {
		ah.logger.Errorf("appointment handler: create appointment service method: %s", err)
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

	createOrEditAppointmentResponseDto := dto.CreateOrEditAppointmentResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrEditAppointmentResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// EditAppointment godoc
//
//	@Summary		Изменить запись
//	@Description	Изменить запись в таблице Appointments
//	@ID				edit-appointment
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор записи"
//	@Param			Appointment	body	dto.CreateOrEditAppointmentRequestDto	true	"Информация о изменяемой записи"
//
// @Success		200	    {object} 	dto.CreateOrEditAppointmentResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		404	    {object}	dto.ErrorDto						"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/appointments/{id} [put]
//
//nolint:funlen // it's ok
func (ah *AppointmentHandler) EditAppointment(w http.ResponseWriter, r *http.Request) {
	ah.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ah.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role == dto.DOCTOR {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var appointment dto.CreateOrEditAppointmentRequestDto
	param, ok := mux.Vars(r)["id"]
	appointmentID, err := strToInt64(param)
	if !ok || err != nil {
		ah.logger.Errorf("appointment handler: %s", err)
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
	err = json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		ah.logger.Errorf("appointment handler: decode json %s", err)
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

	id, err := ah.service.EditAppointment(
		appointmentID,
		appointment.DoctorID,
		appointment.PatientID,
		appointment.DateTime,
	)
	if err != nil {
		ah.logger.Errorf("appointment handler: edit appointment service method: %s", err)
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

	createOrUpdateAppointmentResponseDto := dto.CreateOrEditAppointmentResponseDto{
		ID: id,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createOrUpdateAppointmentResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CancelAppointment godoc
//
//	@Summary		Удалить запись
//	@Description	Удалить запись в таблице Appointments
//	@ID				cancel-appointment
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор записи"
//
// @Success		204	    {object} 	dto.ErrorDto			"Запись успешно удалена"
// @Failure		400	    {object}	dto.ErrorDto			"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto			"Пользователь не авторизован"
// @Failure		404	    {object}	dto.ErrorDto			"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto			"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/appointments/{id} [delete]
func (ah *AppointmentHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	ah.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ah.port)

	param, ok := mux.Vars(r)["id"]
	appointmentID, err := strToInt64(param)
	if !ok || err != nil {
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

	err = ah.service.CancelAppointment(appointmentID)
	if err != nil {
		ah.logger.Errorf("appointment handler: cancel appointment service method: %s", err)
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

	w.WriteHeader(http.StatusNoContent)
}

// GetAppointmentByID godoc
//
//	@Summary		Получение записи по идентификатору
//	@Description	Получить запись по идентификатору в таблице Appointments
//	@ID				get-appointments-by-id
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор записи"
//
// @Success		200	    {object} 	dto.GetAppointmentResponseDto	"Запись"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		404	    {object}	dto.ErrorDto					"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/appointments/{id} [get]
func (ah *AppointmentHandler) GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	ah.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ah.port)

	param, ok := mux.Vars(r)["id"]
	appointmentID, err := strToInt64(param)
	if !ok || err != nil {
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

	appointment, err := ah.service.GetAppointmentByID(appointmentID)
	if err != nil {
		ah.logger.Errorf("appointment handler: get appointment by id service method: %s", err)
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

	getAppointmentResponseDto := dto.ConvertToDto(appointment)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getAppointmentResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetAppointments godoc
//
//	@Summary		Получить записи
//	@Description	Получить все записи из таблицы Appointments
//	@ID				get-appointments
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param 			doctor_id	query	integer		false	"Идентификатор врача"
//	@Param 			patient_id	query	integer		false	"Идентификатор пациента"
//
// @Success		200	    {array} 	dto.GetAppointmentResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto				"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto				"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/appointments [get]
func (ah *AppointmentHandler) GetAppointments( //nolint: cyclop // no other ways
	w http.ResponseWriter,
	r *http.Request,
) {
	ah.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", ah.port)

	doctorID, errDoc := getQueryParam(r, "doctor_id")
	patientID, errPat := getQueryParam(r, "patient_id")
	getAppointmentsResponseDto := make([]*dto.GetAppointmentResponseDto, 0)
	var (
		appointments []*models.Appointment
		err          error
	)

	switch {
	case errPat == nil && errDoc == nil && patientID != 0 && doctorID != 0:
		appointments, err = ah.service.GetAppointmentsByPatientAndDoctor(patientID, doctorID)
		if err != nil {
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
	case errDoc == nil && doctorID != 0:
		appointments, err = ah.service.GetAppointmentsByDoctor(doctorID)
		if err != nil {
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
	case errPat == nil && patientID != 0:
		appointments, err = ah.service.GetAppointmentsByPatient(patientID)
		if err != nil {
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
	default:
		appointments, err = ah.service.GetAllAppointments()
		if err != nil {
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
	}

	for _, appointment := range appointments {
		getAppointmentsResponseDto = append(getAppointmentsResponseDto, dto.ConvertToDto(appointment))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getAppointmentsResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func strToInt64(val string) (int64, error) {
	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, ErrValidate
	}

	return result, nil
}

func getQueryParam(r *http.Request, key string) (int64, error) {
	var val string
	if val = r.URL.Query().Get(key); val == "" {
		return 0, nil
	}

	result, err := strToInt64(val)
	if err != nil {
		return 0, err
	}

	return result, nil
}
