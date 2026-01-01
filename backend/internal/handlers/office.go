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

type OfficeService interface {
	EditOffice(id, number, floor int64) (int64, error)
	AddOffice(number, floor int64) (int64, error)
	GetAllOffices() ([]models.Office, error)
	GetOfficeByID(id int64) (models.Office, error)
}

type OfficeHandler struct {
	service OfficeService
	logger  *zap.SugaredLogger
	port    string
}

func NewOfficeHandler(s OfficeService, logger *zap.SugaredLogger, port string) *OfficeHandler {
	return &OfficeHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

// AddOffice godoc
//
//	@Summary		Добавить кабинет
//	@Description	Добавить кабинет в таблицу offices
//	@ID				add-office
//	@Tags			offices
//	@Accept			json
//	@Produce		json
//	@Param			Office	body	dto.CreateOrEditOfficeRequestDto	true	"Информация о кабинете"
//
// @Success		201	    {object} 	dto.CreateOrEditOfficeResponseDto	"Created"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/offices [post]
func (oh *OfficeHandler) AddOffice(w http.ResponseWriter, r *http.Request) {
	oh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", oh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.ADMIN {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var office dto.CreateOrEditOfficeRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&office)
	if err != nil {
		oh.logger.Errorf("office handler: json decode %s", err)
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

	id, err := oh.service.AddOffice(
		office.Number,
		office.Floor,
	)
	if err != nil {
		oh.logger.Errorf("office handler: add office service method: %s", err)
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

	createOrEditOfficeResponseDto := dto.CreateOrEditOfficeResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrEditOfficeResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// EditOffice godoc
//
//	@Summary		Изменить информацию о кабинете
//	@Description	Изменить информацию о каибинете в таблице Offices
//	@ID				edit-office
//	@Tags			offices
//	@Accept			json
//	@Produce		json
//	@Param 			id		path	integer		true	"Идентификатор кабинета"
//	@Param			Office	body	dto.CreateOrEditOfficeRequestDto	true	"Информация о изменяемом кабинете"
//
// @Success		200	    {object} 	dto.CreateOrEditOfficeResponseDto	"OK"
// @Failure		400	    {object}	dto.ErrorDto						"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto						"Пользователь не авторизован"
// @Failure		403	    {object}	dto.ErrorDto						"Пользователь не имеет доступа"
// @Failure		500	    {object}	dto.ErrorDto						"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/offices/{id} [patch]
//
//nolint:funlen // it's ok
func (oh *OfficeHandler) EditOffice(w http.ResponseWriter, r *http.Request) {
	oh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", oh.port)

	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims) //nolint: errcheck // always props
	role := int64(claims["role"].(float64))
	if role != dto.ADMIN {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var office dto.CreateOrEditOfficeRequestDto

	param, ok := mux.Vars(r)["id"]
	officeID, err := strToInt64(param)
	if !ok || err != nil {
		oh.logger.Errorf("office handler: %s", err)
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
	err = json.NewDecoder(r.Body).Decode(&office)
	if err != nil {
		oh.logger.Errorf("office handler: json decode %s", err)
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

	id, err := oh.service.EditOffice(
		officeID,
		office.Number,
		office.Floor,
	)
	if err != nil {
		oh.logger.Errorf("office handler: edit office service method: %s", err)
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

	createOrUpdateOfficeResponseDto := dto.CreateOrEditOfficeResponseDto{
		ID: id,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(createOrUpdateOfficeResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetAllOffices godoc
//
//	@Summary		Получить список кабинетов
//	@Description	Получить все кабинеты из таблицы Offices
//	@ID				get-offices
//	@Tags			offices
//	@Accept			json
//	@Produce		json
//
// @Success		200	    {array} 	dto.GetOfficeResponseDto	"OK"
// @Failure		401	    {object}	dto.ErrorDto				"Пользователь не авторизован"
// @Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/offices [get]
func (oh *OfficeHandler) GetAllOffices(w http.ResponseWriter, r *http.Request) {
	oh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", oh.port)

	allOffices, err := oh.service.GetAllOffices()
	if err != nil {
		oh.logger.Errorf("office handler: get all offices service method: %s", err)
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

	getAllOfficesResponseDto := make([]*dto.GetOfficeResponseDto, 0)

	for idx := range allOffices {
		getAllOfficesResponseDto = append(
			getAllOfficesResponseDto,
			dto.ConvertToOfficeDto(&allOffices[idx]),
		)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getAllOfficesResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetOfficeByID godoc
//
//	@Summary		Получение информации о кабинете по идентификатору
//	@Description	Получить информацию о кабинете по идентификатору в таблице Offices
//	@ID				get-office-by-id
//	@Tags			offices
//	@Accept			json
//	@Produce		json
//	@Param 			id				path	integer		true	"Идентификатор кабинета"
//
// @Success		200	    {object} 	dto.GetOfficeResponseDto		"Кабинет"
// @Failure		400	    {object}	dto.ErrorDto					"Некорректные данные"
// @Failure		401	    {object}	dto.ErrorDto					"Пользователь не авторизован"
// @Failure		404	    {object}	dto.ErrorDto					"Запись не найдена"
// @Failure		500	    {object}	dto.ErrorDto					"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/api/v1/offices/{id} [get]
func (oh *OfficeHandler) GetOfficeByID(w http.ResponseWriter, r *http.Request) {
	oh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", oh.port)

	param, ok := mux.Vars(r)["id"]
	officeID, err := strToInt64(param)
	if !ok || err != nil {
		oh.logger.Errorf("office handler: %s", err)
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

	office, err := oh.service.GetOfficeByID(officeID)
	if err != nil {
		oh.logger.Errorf("office handler: get office by id service method: %s", err)
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

	getOfficeResponseDto := dto.ConvertToOfficeDto(&office)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getOfficeResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
