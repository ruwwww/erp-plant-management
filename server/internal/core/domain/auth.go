package domain

import "time"

type User struct {
	ID           int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string     `gorm:"unique;not null;size:320" json:"email"`
	Phone        *string    `gorm:"size:32" json:"phone"`
	PasswordHash *string    `gorm:"size:255" json:"-"` // Hide from JSON
	FirstName    *string    `gorm:"size:100" json:"first_name"`
	LastName     *string    `gorm:"size:100" json:"last_name"`
	Role         UserRole   `gorm:"not null;default:'CUSTOMER'" json:"role"`
	Locale       *string    `gorm:"size:10;default:'en'" json:"locale"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:current_timestamp" json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type Customer struct {
	ID                int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            *int       `gorm:"unique" json:"user_id"`
	User              *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CompanyName       *string    `gorm:"size:255" json:"company_name"`
	BillingAddressID  *int       `json:"billing_address_id"`
	ShippingAddressID *int       `json:"shipping_address_id"`
	BillingAddress    *Address   `gorm:"foreignKey:BillingAddressID" json:"billing_address,omitempty"`
	ShippingAddress   *Address   `gorm:"foreignKey:ShippingAddressID" json:"shipping_address,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

type Supplier struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"not null;size:255" json:"name"`
	ContactName *string    `gorm:"size:255" json:"contact_name"`
	Email       *string    `gorm:"size:320" json:"email"`
	Phone       *string    `gorm:"size:32" json:"phone"`
	AddressID   *int       `json:"address_id"`
	Address     *Address   `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type Address struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Line1      string    `gorm:"not null;size:255" json:"line1"`
	Line2      *string   `gorm:"size:255" json:"line2"`
	City       string    `gorm:"not null;size:100" json:"city"`
	State      *string   `gorm:"size:100" json:"state"`
	PostalCode *string   `gorm:"size:30" json:"postal_code"`
	Country    string    `gorm:"not null;default:'Indonesia';size:100" json:"country"`
	Latitude   *float64  `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude  *float64  `gorm:"type:decimal(10,7)" json:"longitude"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
