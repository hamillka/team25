package models

import "time"

type Appointment struct {
	ID           int
	DoctorID     int
	PatientID    int
	DateTime     time.Time
	PatientFIO   string
	PatientEmail string
}
