package main

import (
	"net/http"
	"time"

	cfg "github.com/hamillka/team25/backend/internal/config"
	"github.com/hamillka/team25/backend/internal/db"
	"github.com/hamillka/team25/backend/internal/handlers"
	"github.com/hamillka/team25/backend/internal/logger"
	"github.com/hamillka/team25/backend/internal/repositories"
	"github.com/hamillka/team25/backend/internal/services"
)

// @title DicDoc Service
// @version 1.0
// @description DicDoc Service
// @po
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						auth-x
//	@description				Authorization check
func main() {
	config, err := cfg.New()
	logger := logger.CreateLogger(config.Log)

	defer func() {
		err = logger.Sync()
		if err != nil {
			logger.Errorf("Error while syncing logger: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	dbInstance := db.NewConn(&config.DB, 25, logger)
	dbConn := dbInstance.GetConn()

	defer func() {
		err = dbConn.Close()
		if err != nil {
			logger.Errorf("Error while closing connection to db: %v", err)
		}
	}()

	ar := repositories.NewAppointmentRepository(dbConn)
	dr := repositories.NewDoctorRepository(dbConn)
	or := repositories.NewOfficeRepository(dbConn)
	pr := repositories.NewPatientRepository(dbConn)
	ur := repositories.NewUserRepository(dbConn)
	tr := repositories.NewTimetableRepository(dbConn)
	mhr := repositories.NewMedicalHistoryRepository(dbConn)

	as := services.NewAppointmentService(ar)
	ds := services.NewDoctorService(dr, tr)
	os := services.NewOfficeService(or)
	ps := services.NewPatientService(pr)
	us := services.NewUserService(ur, dr, pr)
	ts := services.NewTimetableService(tr)
	mhs := services.NewMedicalHistoryService(mhr)

	r := handlers.Router(as, ds, os, ps, us, ts, mhs, logger, config.Port)
	ch := handlers.NewCors()

	port := config.Port

	server := &http.Server{
		Addr:              "localhost:" + port,
		Handler:           ch(r),
		ReadHeaderTimeout: time.Duration(config.Timeout) * time.Second,
	}

	logger.Info("Server is started on port ", port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatalf("Error while starting server: %v", err)
	}
}
