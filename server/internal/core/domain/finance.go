package domain

import "time"

type Invoice struct {
	ID              int           `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceNumber   *string       `gorm:"unique;size:64" json:"invoice_number"`
	Type            InvoiceType   `gorm:"not null" json:"type"`
	SalesOrderID    *int          `json:"sales_order_id"`
	SalesOrder      *SalesOrder   `gorm:"foreignKey:SalesOrderID" json:"sales_order"`
	PurchaseOrderID *int          `json:"purchase_order_id"`
	IssuedAt        time.Time     `gorm:"type:date;not null;default:current_date" json:"issued_at"`
	DueAt           time.Time     `gorm:"type:date;not null" json:"due_at"`
	TotalAmount     float64       `gorm:"not null;type:decimal(12,2)" json:"total_amount"`
	AmountResidual  float64       `gorm:"not null;type:decimal(12,2)" json:"amount_residual"`
	PDFUrl          *string       `gorm:"size:512" json:"pdf_url"`
	Status          InvoiceStatus `gorm:"not null;default:'DRAFT'" json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// POSSession moved to pos.go

type POSCashMove struct {
	ID           int          `gorm:"primaryKey;autoIncrement" json:"id"`
	POSSessionID int          `gorm:"not null" json:"pos_session_id"`
	UserID       int          `gorm:"not null" json:"user_id"`
	Amount       float64      `gorm:"not null;type:decimal(12,2)" json:"amount"`
	Type         CashMoveType `gorm:"not null" json:"type"`
	Reason       *string      `gorm:"size:255" json:"reason"`
	CreatedAt    time.Time    `json:"created_at"`
}

type FinancialAccount struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"not null;size:100" json:"name"`
	Code        *string `gorm:"unique;size:50" json:"code"`
	Type        string  `gorm:"not null;size:50" json:"type"`
	Currency    string  `gorm:"not null;default:'IDR';size:10" json:"currency"`
	Description *string `json:"description"`
}

type Payment struct {
	ID                   int             `gorm:"primaryKey;autoIncrement" json:"id"`
	SalesOrderID         *int            `json:"sales_order_id"`
	PurchaseOrderID      *int            `json:"purchase_order_id"`
	FinancialAccountID   int             `gorm:"not null" json:"financial_account_id"`
	TransactionReference *string         `gorm:"size:255" json:"transaction_reference"`
	Type                 TransactionType `gorm:"not null;default:'CREDIT'" json:"type"`
	InvoiceID            *int            `json:"invoice_id"`
	Amount               float64         `gorm:"not null;type:decimal(14,2);default:0" json:"amount"`
	Currency             string          `gorm:"not null;default:'IDR';size:10" json:"currency"`
	Method               string          `gorm:"not null;size:50" json:"method"`
	Status               PaymentStatus   `gorm:"not null;default:'UNPAID'" json:"status"`
	PaidAt               *time.Time      `json:"paid_at"`
	CreatedAt            time.Time       `json:"created_at"`
}
