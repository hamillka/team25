package main

import (
	"time"

	cfg "github.com/hamillka/team25/reminder/internal/config"
	"github.com/hamillka/team25/reminder/internal/db"
	"github.com/hamillka/team25/reminder/internal/logger"
	"github.com/hamillka/team25/reminder/internal/repositories"
	"github.com/hamillka/team25/reminder/internal/sender"
	"github.com/hamillka/team25/reminder/internal/services"
)

func main() {
	config, err := cfg.New()
	logger := logger.CreateLogger(config.Log)
	sender := sender.NewSender(config.Sender)

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

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	repo := repositories.NewAppointmentRepository(dbConn, sender)
	service := services.NewService(repo)

	service.CheckAppointments()
	// Продолжаем проверять по расписанию
	for range ticker.C {
		service.CheckAppointments()
	}
}
