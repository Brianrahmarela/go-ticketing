package repositories

import (
	"database/sql"
	"go-ticketing/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TicketRepository defines data access methods for Ticket model
type TicketRepository interface {
	Create(tx *gorm.DB, ticket *models.Ticket) error
	FindByID(id uint) (*models.Ticket, error)
	FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Ticket, error)
	FindByUser(userID uint) ([]models.Ticket, error)
	FindByEvent(eventID uint) ([]models.Ticket, error)
	SumPurchasedByEvent(eventID uint) (uint, error)
	Update(tx *gorm.DB, ticket *models.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) Create(tx *gorm.DB, ticket *models.Ticket) error {
	if tx != nil {
		return tx.Create(ticket).Error
	}
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindByID(id uint) (*models.Ticket, error) {
	var t models.Ticket
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *ticketRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Ticket, error) {
	var t models.Ticket
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ticketRepository) FindByUser(userID uint) ([]models.Ticket, error) {
	var list []models.Ticket
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ticketRepository) FindByEvent(eventID uint) ([]models.Ticket, error) {
	var list []models.Ticket
	if err := r.db.Where("event_id = ?", eventID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ticketRepository) SumPurchasedByEvent(eventID uint) (uint, error) {
	var total sql.NullInt64
	row := r.db.Model(&models.Ticket{}).
		Select("SUM(quantity)").
		Where("event_id = ? AND status = ?", eventID, "purchased").
		Row()
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return uint(total.Int64), nil
}

func (r *ticketRepository) Update(tx *gorm.DB, ticket *models.Ticket) error {
	if tx != nil {
		return tx.Save(ticket).Error
	}
	return r.db.Save(ticket).Error
}
func (r *ticketRepository) FindByUserID(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}
