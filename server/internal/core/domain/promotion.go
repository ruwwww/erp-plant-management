package domain

import "time"

type Promotion struct {
	ID              int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string                 `gorm:"not null;size:255" json:"name"`
	Code            *string                `gorm:"unique;size:64" json:"code"`
	Description     *string                `json:"description"`
	Priority        int                    `gorm:"not null;default:0" json:"priority"`
	IsExclusive     bool                   `gorm:"not null;default:false" json:"is_exclusive"`
	IsActive        bool                   `gorm:"not null;default:true" json:"is_active"`
	Conditions      map[string]interface{} `gorm:"type:jsonb;not null;default:'{}'" json:"conditions"`
	Actions         map[string]interface{} `gorm:"type:jsonb;not null;default:'{}'" json:"actions"`
	StartsAt        *time.Time             `json:"starts_at"`
	EndsAt          *time.Time             `json:"ends_at"`
	TotalUsageLimit *int                   `json:"total_usage_limit"`
	PerUserLimit    *int                   `gorm:"default:1" json:"per_user_limit"`
	GLAccountCode   *string                `gorm:"default:'SALES_DISC';size:50" json:"gl_account_code"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type PromotionUsage struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PromotionID  int       `gorm:"not null" json:"promotion_id"`
	CustomerID   int       `gorm:"not null" json:"customer_id"`
	SalesOrderID int       `gorm:"not null" json:"sales_order_id"`
	UsedAt       time.Time `gorm:"not null;default:current_timestamp" json:"used_at"`
}

type SalesOrderPromotion struct {
	ID             int                    `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID   int                    `gorm:"not null" json:"sales_order_id"`
	PromotionID    int                    `gorm:"not null" json:"promotion_id"`
	DiscountAmount float64                `gorm:"not null;type:decimal(12,2)" json:"discount_amount"`
	MetaData       map[string]interface{} `gorm:"type:jsonb" json:"meta_data"`
}
