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

type DoctorService interface {
	EditDoctor(id int64, fio, phoneNumber, email, specialization string) (int64, error)
	AddDoctor(fio, phoneNumber, email, specialization string) (int64, error)
	GetAllDoctors() ([]models.Doctor, error)
	GetDoctorByID(id int64) (models.Doctor, error)
}

type DoctorHandler struct {
	service DoctorService
	logger  *zap.SugaredLogger
	port    string
}

func NewDoctorHandler(s DoctorService, logger *zap.SugaredLogger, port string) *DoctorHandler {
	return &DoctorHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// AddDoctor godoc
//
//	@Summary		Добавить врача
//	@Description	Добавить врача в таблицу Doctors
//	@ID				add-doctor
//	@Tags			doctors
//	@Accept			json
//	@Produce		json
//	@Param			Doctor	body	dto.CreateOrEditDoctorRequestDto	true	"Информация о враче"
//
// @Success		201	    {object} 	dto.CreateOrEditDoctorResponseDto	"Created"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/doctors [post]
func (dh *DoctorHandler) AddDoctor(w http.ResponseWriter, r *http.Request) {
	dh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", dh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.ADMIN {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var doctor dto.CreateOrEditDoctorRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil {
		dh.logger.Errorf("doctor handler: json decode %s", err)
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

	id, err := dh.service.AddDoctor(
		doctor.Fio,
		doctor.PhoneNumber,
		doctor.Email,
		doctor.Specialization,
	)
	if err != nil {
		dh.logger.Errorf("doctor handler: add doctor service method: %s", err)
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

	createOrEditDoctorResponseDto := dto.CreateOrEditDoctorResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrEditDoctorResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// EditDoctor godoc
//
//	@Summary		Изменить информацию о враче
//	@Description	Изменить информацию о враче в таблице Doctors
//	@ID				edit-doctor
//	@Tags			doctors
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор врача"
//	@Param			Doctor	body	dto.CreateOrEditDoctorRequestDto	true	"Информация о изменяемом враче"
//
// @Success		200	    {object} 	dto.CreateOrEditDoctorResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/doctors/{id} [patch]
//
//nolint:funlen // it's ok
func (dh *DoctorHandler) EditDoctor(w http.ResponseWriter, r *http.Request) {
	dh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", dh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role == dto.USER {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var doctor dto.CreateOrEditDoctorRequestDto

	param, ok := mux.Vars(r)["id"]
	doctorID, err := strToInt64(param)
	if !ok || err != nil {
		dh.logger.Errorf("doctor handler: %s", err)
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
	err = json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil {
		dh.logger.Errorf("doctor handler: json decode %s", err)

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

	id, err := dh.service.EditDoctor(
		doctorID,
		doctor.Fio,
		doctor.PhoneNumber,
		doctor.Email,
		doctor.Specialization,
	)
	if err != nil {
		dh.logger.Errorf("doctor handler: edit doctor service method: %s", err)
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

	createOrUpdateDoctorResponseDto := dto.CreateOrEditDoctorResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createOrUpdateDoctorResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetAllDoctors godoc
//
//	@Summary		Получить список врачей
//	@Description	Получить всех врачей из таблицы Doctors
//	@ID				get-doctors
//	@Tags			doctors
//	@Accept			json
//	@Produce		json
//
// @Success		200	    {array} 	dto.GetDoctorResponseDto	"OK"
// @Failure		401	    {object}	dto.ErrorDto				"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/doctors [get]
func (dh *DoctorHandler) GetAllDoctors(w http.ResponseWriter, r *http.Request) {
	dh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", dh.port)

	allDoctors, err := dh.service.GetAllDoctors()
	if err != nil {
		dh.logger.Errorf("doctor handler: get all doctors service method: %s", err)
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

	getAllDoctorsResponseDto := make([]*dto.GetDoctorResponseDto, 0)

	for idx := range allDoctors {
		getAllDoctorsResponseDto = append(
			getAllDoctorsResponseDto,
			dto.ConvertToDoctorDto(&allDoctors[idx]),
		)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getAllDoctorsResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetDoctorByID godoc
//
//	@Summary		Получение врача по идентификатору
//	@Description	Получить врача по идентификатору в таблице Doctors
//	@ID				get-doctor-by-id
//	@Tags			doctors
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор врача"
//
// @Success		200	    {object} 	dto.GetDoctorResponseDto		"Врач"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		404	    {object}	dto.ErrorDto					"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/doctors/{id} [get]
func (dh *DoctorHandler) GetDoctorByID(w http.ResponseWriter, r *http.Request) {
	dh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", dh.port)
	param, ok := mux.Vars(r)["id"]
	doctorID, err := strToInt64(param)
	if !ok || err != nil {
		dh.logger.Errorf("doctor handler: %s", err)
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

	doctor, err := dh.service.GetDoctorByID(doctorID)
	if err != nil {
		dh.logger.Errorf("doctor handler: get doctor by id service method: %s", err)
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

	getDoctorResponseDto := dto.ConvertToDoctorDto(&doctor)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getDoctorResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
