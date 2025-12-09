package dto

type PaymentWebhookRequest struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
	Signature     string `json:"signature"`
	ProviderRef   string `json:"provider_ref"`
}

type LogisticsWebhookRequest struct {
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"` // DELIVERED, FAILED
	Timestamp      string `json:"timestamp"`
}
