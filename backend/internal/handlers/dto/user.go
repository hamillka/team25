package dto

import "github.com/hamillka/team25/backend/internal/models"

const (
	ADMIN = iota
	USER
	DOCTOR
)

// UserDto model info
// @Description Информация о пользователе
type UserDto struct {
	PatientID *int64 `json:"patientId"` // Идентификатор пациента (если пользователь - пациент)
	DoctorID  *int64 `json:"doctorId"`  // Идентификатор врач (если пользователь - врач)
	Login     string `json:"login"`     // Логин пользователя
	Password  string `json:"password"`  // Пароль пользователя
	Role      int64  `json:"role"`      // Роль пользователя (Админ, Пациент, Врач)
	ID        int64  `json:"id"`        // Идентификатор пользователя
}

// UserLoginRequestDto model info
// @Description Информация о пользователе при попытке входа
type UserLoginRequestDto struct {
	Login    string `json:"login"`    // Логин пользователя
	Password string `json:"password"` // Пароль пользователя
}

// UserLoginResponseDto model info
// @Description Информация о пользователе при попытке входа
type UserLoginResponseDto struct {
	JWTToken string  `json:"jwtToken"` // JWT-токен
	User     UserDto `json:"user"`     // Информация о пользователе
}

// UserRegisterRequestDto model info
// @Description Информация о пользователе при попытке регистрации
type UserRegisterRequestDto struct {
	FIO            string `json:"fio"`            // ФИО пользователя
	PhoneNumber    string `json:"phoneNumber"`    // Телефон пользователя
	Email          string `json:"email"`          // Почта пользователя
	Insurance      string `json:"insurance"`      // Страховка пользователя (если пользователь - пациент)
	Specialization string `json:"specialization"` // Специализация пользователя (если пользователь - врач)
	Login          string `json:"login"`          // Логин пользователя
	Password       string `json:"password"`       // Пароль пользователя
	Role           int64  `json:"role"`           // Роль пользователя
}

// UserRegisterResponseDto model info
// @Description Информация о пользователе при попытке регистрации
type UserRegisterResponseDto struct {
	ID int64 `json:"id"` // Идентификатор пользователя
}

func ConvertToUserDto(user *models.User) *UserDto {
	return &UserDto{
		PatientID: user.PatientID,
		DoctorID:  user.DoctorID,
		Login:     user.Login,
		Password:  user.Password,
		Role:      user.Role,
		ID:        user.ID,
	}
}
