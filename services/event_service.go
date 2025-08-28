package services

import (
	"go-ticketing/models"
	"go-ticketing/repositories"
)

type EventService interface {
	Create(event *models.Event) error
	Update(event *models.Event) error
	Delete(id uint) error
	GetByID(id uint) (*models.Event, error)
	GetAll() ([]models.Event, error)
}

type eventService struct {
	repo repositories.EventRepository
}

func NewEventService(repo repositories.EventRepository) EventService {
	return &eventService{repo}
}

func (s *eventService) Create(event *models.Event) error {
	return s.repo.Create(event)
}

func (s *eventService) Update(event *models.Event) error {
	return s.repo.Update(nil, event) // tx nil → update langsung
}

func (s *eventService) Delete(id uint) error {
	ev, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Update(nil, &models.Event{ID: ev.ID, Status: "Deleted"})
}

func (s *eventService) GetByID(id uint) (*models.Event, error) {
	return s.repo.FindByID(id)
}

func (s *eventService) GetAll() ([]models.Event, error) {
	return s.repo.ListAll()
}
