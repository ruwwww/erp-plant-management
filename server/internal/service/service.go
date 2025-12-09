package service

import (
	"context"
	"server/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)

	RegisterStaff(ctx context.Context, user *domain.User, plainPassword string) error

	ResetPassword(ctx context.Context, userID int, newPassword string) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
	AddAddress(ctx context.Context, userID int, addr *domain.Address) error
	UpdateAddress(ctx context.Context, addr *domain.Address) error
	DeleteAddress(ctx context.Context, userID, addressID int) error
	SetDefaultAddress(ctx context.Context, userID, addressID int, isBilling bool) error
}

type ProductFilterParams struct {
	Search       string
	CategorySlug string
	MinPrice     float64
	MaxPrice     float64
	IsActive     *bool
	Page         int
	Limit        int
}

type CatalogService interface {
	GetProducts(ctx context.Context, filter ProductFilterParams) ([]domain.Product, int64, error)

	GetProductDetail(ctx context.Context, slug string) (*domain.Product, error)

	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	SoftDeleteProduct(ctx context.Context, id int) error
	RestoreProduct(ctx context.Context, id int) error
}

type CartCalculationResult struct {
	Subtotal       float64
	DiscountAmount float64
	TaxAmount      float64
	ShippingAmount float64
	TotalAmount    float64
	Items          []domain.SalesOrderItem
}

type CartService interface {
	CalculateCart(ctx context.Context, items []domain.SalesOrderItem, couponCode string) (*CartCalculationResult, error)
}

type OrderService interface {
	PlaceOrder(ctx context.Context, order *domain.SalesOrder) error

	CancelOrder(ctx context.Context, orderID int, reason string) error

	ProcessReturn(ctx context.Context, orderID int, items []domain.Return) error
}

type POSService interface {
	OpenSession(ctx context.Context, userID int, openingFloat float64) (*domain.POSSession, error)

	CloseSession(ctx context.Context, sessionID int, closingCashActual float64, note string) error

	RecordCashMove(ctx context.Context, sessionID int, amount float64, moveType domain.CashMoveType, reason string) error

	GetActiveSession(ctx context.Context, userID int) (*domain.POSSession, error)
}

type StockMoveCmd struct {
	LocationID    int
	VariantID     int
	QtyChange     int
	Reason        domain.MovementReason
	ReferenceID   int
	ReferenceType string
	UserID        int
}

type InventoryService interface {
	ExecuteMovement(ctx context.Context, cmd StockMoveCmd) error

	TransferStock(ctx context.Context, variantID, qty, fromLocID, toLocID, userID int) error

	GetStockLevel(ctx context.Context, variantID, locationID int) (int, error)

	BulkAdjustStock(ctx context.Context, cmds []StockMoveCmd) error
}

type AssemblyService interface {
	AssembleKit(ctx context.Context, variantID, qty int, userID int) error

	Disassemble(ctx context.Context, variantID, qty int, userID int) error
}

type ProcurementService interface {
	CreatePO(ctx context.Context, po *domain.PurchaseOrder) error

	ReceivePO(ctx context.Context, poID int, receivedItems map[int]int) error
}

type FinanceService interface {
	GenerateInvoice(ctx context.Context, orderID int) (*domain.Invoice, error)

	RecordPayment(ctx context.Context, invoiceID int, amount float64, method string) error
}
