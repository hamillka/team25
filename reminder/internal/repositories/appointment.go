package repositories

import (
	"log"
	"time"

	"github.com/hamillka/team25/reminder/internal/models"
	"github.com/hamillka/team25/reminder/internal/sender"
	"github.com/jmoiron/sqlx"
)

type AppointmentRepository struct {
	db     *sqlx.DB
	sender *sender.Sender
}

func NewAppointmentRepository(db *sqlx.DB, s *sender.Sender) *AppointmentRepository {
	return &AppointmentRepository{db: db, sender: s}
}

func (r *AppointmentRepository) CheckAppointments() {
	now := time.Now()

	targetTime := now.Add(24 * time.Hour)

	log.Printf("Проверяем приёмы запланированные на %v", targetTime.Format("2006-01-02 15:04:05"))

	query := `
		SELECT a.id, a.doctorid, a.patientid, a.datetime, p.fio, p.email
		FROM appointments a
		JOIN patients p ON a.patientid = p.id
		WHERE a.datetime BETWEEN $1 AND $2
	`

	windowStart := targetTime.Add(-30 * time.Minute)
	windowEnd := targetTime.Add(30 * time.Minute)

	rows, err := r.db.Query(query, windowStart, windowEnd)
	if err != nil {
		log.Printf("Ошибка при запросе приёмов: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var appointment models.Appointment
		err := rows.Scan(
			&appointment.ID,
			&appointment.DoctorID,
			&appointment.PatientID,
			&appointment.DateTime,
			&appointment.PatientFIO,
			&appointment.PatientEmail,
		)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			continue
		}

		r.sender.SendReminderEmail(appointment)
	}
}
