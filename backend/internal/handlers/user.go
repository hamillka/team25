package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/team25/backend/internal/handlers/dto"
	"github.com/hamillka/team25/backend/internal/handlers/middlewares"
	"github.com/hamillka/team25/backend/internal/models"
	"go.uber.org/zap"
)

type UserService interface {
	Login() error
	Register(fio, phoneNumber, email, insurance, specialization, login, password string, role int64) (int64, error)
	CheckUserRole(id int64) (int64, error)
	GetUserByLoginAndPassword(login, password string) (models.User, error)
}

type UserHandler struct {
	service UserService
	logger  *zap.SugaredLogger
	port    string
}

func NewUserHandler(s UserService, logger *zap.SugaredLogger, port string) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  logger,
		port:    port,
	}
}

func createToken(role int64) (string, error) {
	payload := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 720).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(middlewares.Secret)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Login godoc
//
// @Summary			Вход в систему
// @Description		Войти в систему по логину и паролю
// @ID				login
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			loginInfo	body	dto.UserLoginRequestDto	true	"Логин и пароль пользователя"
//
// @Success		200	{object} 	dto.UserLoginResponseDto					"Информация о пользователе и JWT-токен"
// @Failure		400	{object}	dto.ErrorDto								"Некорректные данные"
// @Failure		401	{object}	dto.ErrorDto								"Пользователь не авторизован"
// @Failure		404	{object}	dto.ErrorDto								"Пользователь не найден"
// @Failure		500	{object}	dto.ErrorDto								"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/auth/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	uh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", uh.port)

	var userLoginDto dto.UserLoginRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userLoginDto)
	if err != nil {
		uh.logger.Errorf("user handler: json decode %s", err)
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

	user, err := uh.service.GetUserByLoginAndPassword(userLoginDto.Login, userLoginDto.Password)
	if err != nil {
		uh.logger.Errorf("user handler: get user by login and password service method: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		errorDto := &dto.ErrorDto{
			Error: "Неверное имя пользователя или пароль",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	t, err := createToken(user.Role)
	if err != nil {
		uh.logger.Errorf("user handler: create token method: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Error: "Возникла внутренняя ошибка при авторизации",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userReqDto := dto.ConvertToUserDto(&user)

	userResponseDto := dto.UserLoginResponseDto{
		JWTToken: t,
		User:     *userReqDto,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Register godoc
//
// @Summary			Регистрация пользователя
// @Description		Добавить пользователя в таблицу Users
// @ID				register
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			registerInfo	body	dto.UserRegisterRequestDto	true	"Информация о пользователе"
//
// @Success		201	{object} 	dto.UserRegisterResponseDto					"Идентификатор нового пользователя"
// @Failure		400	{object}	dto.ErrorDto								"Некорректные данные"
// @Failure		500	{object}	dto.ErrorDto								"Внутренняя ошибка сервера"
// @Security		ApiKeyAuth
// @Router			/auth/register [post]
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	uh.logger.Infof("CURRENT APP RUNNING ON PORT : %s\n", uh.port)

	var userDto dto.UserRegisterRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		uh.logger.Errorf("user handler: json decode %s", err)

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

	_, err = uh.service.Register(
		userDto.FIO,
		userDto.PhoneNumber,
		userDto.Email,
		userDto.Insurance,
		userDto.Specialization,
		userDto.Login,
		userDto.Password,
		userDto.Role,
	)
	if err != nil {
		uh.logger.Errorf("user handler: register service method: %s", err)
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
	w.WriteHeader(http.StatusCreated)
}
