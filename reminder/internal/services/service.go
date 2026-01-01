package services

type AppointmentRepository interface {
	CheckAppointments()
}

type Service struct {
	repo AppointmentRepository
}

func NewService(repository AppointmentRepository) *Service {
	return &Service{repo: repository}
}

func (s *Service) CheckAppointments() {
	s.repo.CheckAppointments()
}
