package services

import (
	"errors"

	"go-ticketing/models"
	"go-ticketing/repositories"

	"gorm.io/gorm"
)

type TicketService interface {
	BuyTicket(userID uint, eventID uint, quantity uint) (*models.Ticket, error)
	CancelTicket(ticketID uint, userID uint, role string) error
	GetTicketsByUser(userID uint) ([]models.Ticket, error)
}

type ticketService struct {
	db         *gorm.DB
	ticketRepo repositories.TicketRepository
	eventRepo  repositories.EventRepository
}

func NewTicketService(db *gorm.DB, ticketRepo repositories.TicketRepository, eventRepo repositories.EventRepository) TicketService {
	return &ticketService{db: db, ticketRepo: ticketRepo, eventRepo: eventRepo}
}

func (s *ticketService) BuyTicket(userID uint, eventID uint, quantity uint) (*models.Ticket, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// lock row event
	event, err := s.eventRepo.FindByIDForUpdate(tx, eventID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// cek status event
	if event.Status == "SoldOut" {
		tx.Rollback()
		return nil, errors.New("event is sold out")
	}

	// hitung sisa kursi
	availableSeats := int(event.Capacity - event.Sold)
	if availableSeats < int(quantity) {
		tx.Rollback()
		return nil, errors.New("not enough seats available")
	}

	// update sold
	event.Sold += quantity

	// kalau sudah penuh, tandai sold out
	if event.Sold >= event.Capacity {
		event.Status = "SoldOut"
	}

	if err := s.eventRepo.Update(tx, event); err != nil {
		tx.Rollback()
		return nil, err
	}

	// buat ticket
	ticket := &models.Ticket{
		UserID:    userID,
		EventID:   eventID,
		Quantity:  quantity,
		PricePaid: event.Price * int64(quantity),
		Status:    "purchased",
	}
	if err := s.ticketRepo.Create(tx, ticket); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketService) GetTicketsByUser(userID uint) ([]models.Ticket, error) {
	return s.ticketRepo.FindByUser(userID)
}

func (s *ticketService) CancelTicket(ticketID uint, userID uint, role string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	ticket, err := s.ticketRepo.FindByIDForUpdate(tx, ticketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if role != "admin" && ticket.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized: not your ticket")
	}

	if ticket.Status != "purchased" {
		tx.Rollback()
		return errors.New("cannot cancel this ticket")
	}

	ticket.Status = "cancelled"
	if err := s.ticketRepo.Update(tx, ticket); err != nil {
		tx.Rollback()
		return err
	}

	event, err := s.eventRepo.FindByIDForUpdate(tx, ticket.EventID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if event.Sold >= ticket.Quantity {
		event.Sold -= ticket.Quantity
	}

	if event.Status == "SoldOut" && event.Sold < event.Capacity {
		event.Status = "Active"
	}

	if err := s.eventRepo.Update(tx, event); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
