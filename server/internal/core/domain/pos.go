package domain

import "time"

type POSSession struct {
	ID                int              `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int              `gorm:"not null" json:"user_id"`
	User              *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OpeningCash       float64          `gorm:"not null;type:decimal(12,2);default:0" json:"opening_cash"`
	ClosingCashActual *float64         `gorm:"type:decimal(12,2)" json:"closing_cash_actual"`
	ClosingCashSystem float64          `gorm:"not null;type:decimal(12,2);default:0" json:"closing_cash_system"`
	OpenedAt          time.Time        `gorm:"not null;default:current_timestamp" json:"opened_at"`
	ClosedAt          *time.Time       `json:"closed_at"`
	Note              *string          `json:"note"`
	Status            POSSessionStatus `gorm:"not null;default:'OPEN';size:20" json:"status"` // OPEN, CLOSED
}

type CashMove struct {
	ID           int          `gorm:"primaryKey;autoIncrement" json:"id"`
	POSSessionID int          `gorm:"not null" json:"pos_session_id"`
	Amount       float64      `gorm:"not null;type:decimal(12,2)" json:"amount"`
	Type         CashMoveType `gorm:"not null;size:20" json:"type"` // ADD, DROP
	Reason       string       `gorm:"not null;size:255" json:"reason"`
	CreatedAt    time.Time    `gorm:"not null;default:current_timestamp" json:"created_at"`
}
