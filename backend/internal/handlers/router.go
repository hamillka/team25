package handlers

import (
	_ "github.com/hamillka/team25/backend/api"
	"github.com/hamillka/team25/backend/internal/handlers/middlewares"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

//nolint:funlen // it's ok
func Router(
	as AppointmentService,
	ds DoctorService,
	os OfficeService,
	ps PatientService,
	us UserService,
	ts TimetableService,
	mhs MedicalHistoryService,
	logger *zap.SugaredLogger,
	port string,
) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	secure := router.PathPrefix("/auth").Subrouter()
	fun := router.PathPrefix("/api/v1").Subrouter()

	fun.Use(middlewares.AuthMiddleware)

	ah := NewAppointmentHandler(as, logger, port)
	dh := NewDoctorHandler(ds, logger, port)
	oh := NewOfficeHandler(os, logger, port)
	ph := NewPatientHandler(ps, logger, port)
	uh := NewUserHandler(us, logger, port)
	th := NewTimetableHandler(ts, logger, port)
	mhh := NewMedicalHistoryHandler(mhs, logger, port)

	secure.HandleFunc("/login", uh.Login).Methods("POST")
	secure.HandleFunc("/register", uh.Register).Methods("POST")

	fun.HandleFunc("/appointments", ah.GetAppointments).Methods("GET")
	fun.HandleFunc("/appointments/{id}", ah.GetAppointmentByID).Methods("GET")
	fun.HandleFunc("/appointments/{id}", ah.CancelAppointment).Methods("DELETE")
	fun.HandleFunc("/appointments", ah.CreateAppointment).Methods("POST")
	fun.HandleFunc("/appointments/{id}", ah.EditAppointment).Methods("PUT")

	fun.HandleFunc("/doctors/{id}", dh.EditDoctor).Methods("PATCH")
	fun.HandleFunc("/doctors", dh.AddDoctor).Methods("POST")
	fun.HandleFunc("/doctors", dh.GetAllDoctors).Methods("GET")
	fun.HandleFunc("/doctors/{id}", dh.GetDoctorByID).Methods("GET")

	fun.HandleFunc("/offices/{id}", oh.EditOffice).Methods("PATCH")
	fun.HandleFunc("/offices", oh.AddOffice).Methods("POST")
	fun.HandleFunc("/offices", oh.GetAllOffices).Methods("GET")
	fun.HandleFunc("/offices/{id}", oh.GetOfficeByID).Methods("GET")

	fun.HandleFunc("/patients/{id}", ph.EditPatient).Methods("PATCH")
	fun.HandleFunc("/patients/{id}/medical_history", mhh.GetHistoryByPatient).Methods("GET")
	fun.HandleFunc("/patients", ph.AddPatient).Methods("POST")
	fun.HandleFunc("/patients", ph.GetAllPatients).Methods("GET")
	fun.HandleFunc("/patients/{id}", ph.GetPatientByID).Methods("GET")

	fun.HandleFunc("/locations", th.GetLocationsByDoctor).Methods("GET")
	fun.HandleFunc("/doctors/{id}/workdays", th.GetTimetableByDoctor).Methods("GET")

	fun.HandleFunc("/patients/{id}/medical_history", mhh.CreateMedicalHistory).Methods("POST")
	fun.HandleFunc("/patients/{id}/medical_history", mhh.UpdateMedicalHistory).Methods("PATCH")

	return router
}
