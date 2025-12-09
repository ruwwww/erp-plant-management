package domain

import "time"

type Review struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID  int       `gorm:"not null" json:"product_id"`
	CustomerID *int      `json:"customer_id"`
	Rating     int       `gorm:"not null" json:"rating"`
	Title      *string   `gorm:"size:255" json:"title"`
	Body       *string   `json:"body"`
	Approved   bool      `gorm:"not null;default:false" json:"approved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Setting struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"unique;not null;size:255" json:"key"`
	Value       *string   `json:"value"`
	Description *string   `gorm:"size:255" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	ActorID    *int                   `json:"actor_id"`
	Action     string                 `gorm:"not null;size:100" json:"action"`
	EntityType *string                `gorm:"size:100" json:"entity_type"`
	EntityID   *int                   `json:"entity_id"`
	Changes    map[string]interface{} `gorm:"type:jsonb" json:"changes"`
	IPAddress  *string                `gorm:"size:45" json:"ip_address"`
	CreatedAt  time.Time              `json:"created_at"`
}
