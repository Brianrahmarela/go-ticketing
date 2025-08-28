package models

import "time"

type Event struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Title       string     `json:"title" gorm:"size:200;not null"`
	Description string     `json:"description" gorm:"type:text"`
	StartAt     *time.Time `json:"start_at"` // nullable
	EndAt       *time.Time `json:"end_at"`   // nullable
	Capacity    uint       `json:"capacity" gorm:"not null"`
	Price       int64      `json:"price" gorm:"not null"`
	Status      string     `json:"status" gorm:"size:20;not null;default:'active'"` // active, ongoing, finished
	Sold        uint       `json:"sold" gorm:"default:0"`
	Remaining   int        `json:"remaining" gorm:"-"` // gorm:"-" artinya field ini diabaikan oleh GORM

}
