package models

import "time"

type Ticket struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	EventID   uint       `json:"event_id" gorm:"not null"`
	Quantity  uint       `json:"quantity" gorm:"not null"`
	PricePaid int64      `json:"price_paid" gorm:"not null"`
	Status    string     `json:"status" gorm:"size:20;not null;default:'purchased'"` // purchased, cancelled
}
