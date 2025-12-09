package dto

// --- Inventory ---
type CreateLocationRequest struct {
	Name      string `json:"name" validate:"required"`
	Code      string `json:"code" validate:"required,alphanum"`
	AddressID *int   `json:"address_id"`
}

type InventoryTransferRequest struct {
	VariantID      int `json:"variant_id" validate:"required"`
	Quantity       int `json:"quantity" validate:"required,min=1"`
	FromLocationID int `json:"from_location_id" validate:"required"`
	ToLocationID   int `json:"to_location_id" validate:"required"`
}

type InventoryAdjustRequest struct {
	VariantID  int    `json:"variant_id" validate:"required"`
	LocationID int    `json:"location_id" validate:"required"`
	ChangeQty  int    `json:"change_qty" validate:"required,ne=0"` // Can be negative
	Reason     string `json:"reason" validate:"required"`
}

type BulkAdjustItem struct {
	VariantID  int `json:"variant_id" validate:"required"`
	LocationID int `json:"location_id" validate:"required"`
	ActualQty  int `json:"actual_qty" validate:"required,gte=0"` // Sets exact count
}

type BulkAdjustRequest struct {
	Reason string           `json:"reason" validate:"required"`
	Items  []BulkAdjustItem `json:"items" validate:"required,dive"`
}

// --- Assembly ---
type CreateRecipeRequest struct {
	ParentVariantID int     `json:"parent_variant_id" validate:"required"`
	ChildVariantID  int     `json:"child_variant_id" validate:"required"`
	QuantityNeeded  float64 `json:"quantity_needed" validate:"required,gt=0"`
}

type ExecuteAssemblyRequest struct {
	VariantID int `json:"variant_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

// --- Procurement ---
type CreatePORequest struct {
	SupplierID int         `json:"supplier_id" validate:"required"`
	ExpectedAt string      `json:"expected_at"` // ISO8601
	Notes      string      `json:"notes"`
	Items      []POItemDto `json:"items" validate:"required,dive"`
}

type POItemDto struct {
	VariantID int     `json:"variant_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	UnitCost  float64 `json:"unit_cost" validate:"gte=0"`
}

type ReceivePORequest struct {
	Items []POReceiveItem `json:"items" validate:"required,dive"`
}

type POReceiveItem struct {
	VariantID   int `json:"variant_id" validate:"required"`
	QtyReceived int `json:"qty_received" validate:"required,min=1"`
}

// --- Fulfillment ---
type ShipOrderRequest struct {
	Carrier        string `json:"carrier" validate:"required"`
	TrackingNumber string `json:"tracking_number" validate:"required"`
}
