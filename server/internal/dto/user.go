package dto

// Profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" validate:"omitempty,min=2"`
	Phone     string `json:"phone" validate:"omitempty,e164"` // +62812...
	Bio       string `json:"bio"`
}

// Addresses
type CreateAddressRequest struct {
	Line1      string  `json:"line1" validate:"required"`
	Line2      string  `json:"line2"`
	City       string  `json:"city" validate:"required"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country" validate:"required,iso3166_1_alpha2"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type UpdateAddressRequest struct {
	Line1      string  `json:"line1"`
	Line2      string  `json:"line2"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type SetDefaultAddressRequest struct {
	IsBilling bool `json:"is_billing"`
}

// Orders
type SubmitReviewRequest struct {
	Rating int    `json:"rating" validate:"required,min=1,max=5"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type RequestReturnRequest struct {
	Reason string   `json:"reason" validate:"required"`
	Images []string `json:"images"` // URLs
}
