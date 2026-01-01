package sender

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/hamillka/team25/reminder/internal/models"
)

type Config struct {
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT"`
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
}

type Sender struct {
	config Config
}

func NewSender(cfg Config) *Sender {
	return &Sender{config: cfg}
}

func (s *Sender) SendReminderEmail(appointment models.Appointment) {
	formattedDateTime := appointment.DateTime.Format("02.01.2006 15:04")

	subject := "Напоминание о записи на приём"
	body := fmt.Sprintf(`
Уважаемый(ая) %s,

Напоминаем вам о записи на приём завтра, %s.

С уважением,
Ваша клиника
`, appointment.PatientFIO, formattedDateTime)

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s", appointment.PatientEmail, subject, body))

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		auth,
		s.config.Username,
		[]string{appointment.PatientEmail},
		message,
	)

	if err != nil {
		log.Printf("Ошибка при отправке email пациенту ID=%d: %v", appointment.PatientID, err)
		return
	}

	log.Printf("Напоминание успешно отправлено пациенту ID=%d на email %s", appointment.PatientID, appointment.PatientEmail)
}
