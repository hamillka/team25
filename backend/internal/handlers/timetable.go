package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamillka/team25/backend/internal/handlers/dto"
	"github.com/hamillka/team25/backend/internal/models"
	"go.uber.org/zap"
)

type TimetableService interface {
	GetLocationsByDoctor(id int64) ([]models.Office, error)
	GetDoctorsByLocation(id int64) ([]models.Doctor, error)
	GetWorkdaysByDoctor(id int64) ([]*models.Timetable, error)
}

type TimetableHandler struct {
	service TimetableService
	logger  *zap.SugaredLogger
	port    string
}

func NewTimetableHandler(s TimetableService, logger *zap.SugaredLogger, port string) *TimetableHandler {
	return &TimetableHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// GetLocationsByDoctor godoc
//
//	@Summary		Получение рабочих кабинетов по идентификатору врача
//	@Description	Получить кабинеты по идентификатору врача в таблице TimeTable
//	@ID				get-locations-by-doctor-id
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param 			doctor_id	query	integer		true	"Идентификатор врача"
//
// @Success		200	    {array} 	dto.GetOfficeResponseDto		"Информация о кабинетах"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/locations [get]
func (th *TimetableHandler) GetLocationsByDoctor(w http.ResponseWriter, r *http.Request) {
	th.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", th.port)

	doctorID, err := getQueryParam(r, "doctor_id")
	if err != nil {
		th.logger.Errorf("timetable handler: %s", err)
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

	offices, err := th.service.GetLocationsByDoctor(doctorID)
	if err != nil {
		th.logger.Errorf("timetable handler: get location by doctor service method: %s", err)
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

	getLocationsByDoctor := make([]*dto.GetOfficeResponseDto, 0)

	for idx := range offices {
		getLocationsByDoctor = append(getLocationsByDoctor, dto.ConvertToOfficeDto(&offices[idx]))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getLocationsByDoctor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (th *TimetableHandler) GetDoctorsByLocation(w http.ResponseWriter, r *http.Request) {
	th.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", th.port)

	officeID, err := getQueryParam(r, "office_id")
	if err != nil {
		th.logger.Errorf("timetable handler: %s", err)
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

	doctors, err := th.service.GetDoctorsByLocation(officeID)
	if err != nil {
		th.logger.Errorf("timetable handler: get doctors by location service method: %s", err)
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

	getDoctorsByLocation := make([]*dto.GetDoctorResponseDto, 0)

	for idx := range doctors {
		getDoctorsByLocation = append(getDoctorsByLocation, dto.ConvertToDoctorDto(&doctors[idx]))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getDoctorsByLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetTimetableByDoctor godoc
//
//	@Summary		Получение расписания врача по его идентификатору
//	@Description	Получить  расписание врача по его идентификатору в таблице TimeTable
//	@ID				get-timetable-by-doctor-id
//	@Tags			doctors
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор врача"
//
// @Success		200	    {array} 	dto.GetTimetableResponseDto		"Расписание"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/doctors/{id}/workdays [get]
func (th *TimetableHandler) GetTimetableByDoctor(w http.ResponseWriter, r *http.Request) {
	th.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", th.port)

	param, ok := mux.Vars(r)["id"]
	docID, err := strToInt64(param)
	if !ok || err != nil {
		th.logger.Errorf("timetable handler: %s", err)
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

	timetable, err := th.service.GetWorkdaysByDoctor(docID)
	if err != nil {
		th.logger.Errorf("timetable handler: get timetable by doctor service method: %s", err)
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

	getTimetableByDoctor := make([]*dto.GetTimetableResponseDto, 0)

	for idx := range timetable {
		getTimetableByDoctor = append(getTimetableByDoctor, dto.ConvertToTimetableDto(timetable[idx]))
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getTimetableByDoctor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
