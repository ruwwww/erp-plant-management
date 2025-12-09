package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/dto"
)

// ==========================================
// 1. CORE & AUTH
// ==========================================

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccess string, newRefresh string, err error)
	Logout(ctx context.Context, tokenString string) error

	// Registration & Password Management
	RegisterStaff(ctx context.Context, user *domain.User, plainPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
	ResetPassword(ctx context.Context, userID int, newPassword string) error
}

type UserService interface {
	// Profile Management
	GetProfile(ctx context.Context, userID int) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error

	// Address Book
	GetAddresses(ctx context.Context, userID int) ([]domain.Address, error)
	AddAddress(ctx context.Context, userID int, addr *domain.Address) error
	UpdateAddress(ctx context.Context, addr *domain.Address) error
	DeleteAddress(ctx context.Context, userID, addressID int) error
	SetDefaultAddress(ctx context.Context, userID, addressID int, isBilling bool) error

	// Customer Features
	GetWishlist(ctx context.Context, userID int) ([]domain.ProductVariant, error)
	ToggleWishlist(ctx context.Context, userID, variantID int) (isAdded bool, err error)

	// HR / Admin User Management
	GetUserList(ctx context.Context, filter UserFilterParams) ([]domain.User, int64, error)
	GetUserDetail(ctx context.Context, targetUserID int) (*domain.User, error)
	UpdateUserStatus(ctx context.Context, userID int, isActive bool) error // Ban/Unban
	AssignRoles(ctx context.Context, userID int, role domain.UserRole) error

	// Soft Delete / Restore
	SoftDeleteUser(ctx context.Context, userID int) error
	RestoreUser(ctx context.Context, userID int) error
	ForceDeleteUser(ctx context.Context, userID int) error
}

type UserFilterParams struct {
	Search   string
	Role     string
	IsActive *bool
	Page     int
	Limit    int
}

// ==========================================
// 2. STORE & CATALOG
// ==========================================

type CatalogService interface {
	// Browsing
	GetProducts(ctx context.Context, filter dto.ProductFilterParams) ([]domain.Product, int64, error)
	GetProductDetail(ctx context.Context, slug string) (*domain.Product, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetVariants(ctx context.Context, productID int) ([]domain.ProductVariant, error)

	// Admin Management
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, id int, req dto.UpdateProductRequest) error
	UpdateVariants(ctx context.Context, productID int, variants []domain.ProductVariant) error
	SoftDeleteProduct(ctx context.Context, id int) error
	RestoreProduct(ctx context.Context, id int) error
	ForceDeleteProduct(ctx context.Context, id int) error // Hard delete

	SoftDeleteVariant(ctx context.Context, id int) error
	RestoreVariant(ctx context.Context, id int) error
	ForceDeleteVariant(ctx context.Context, id int) error

	// Data Operations
	ImportProducts(ctx context.Context, data []byte) error // Process CSV/JSON
	ExportProducts(ctx context.Context) ([]byte, error)    // Return CSV/Excel bytes
}

type MediaService interface {
	Upload(ctx context.Context, file []byte, filename string, mimeType string) (*domain.MediaAsset, error)
	LinkMedia(ctx context.Context, mediaID int, entityType string, entityID int, zone string) error
	UnlinkMedia(ctx context.Context, mediaID int, entityType string, entityID int) error
}

// ==========================================
// 3. SALES: CART, ORDER & POS
// ==========================================

type CartCalculationResult struct {
	Subtotal       float64
	DiscountAmount float64
	TaxAmount      float64
	ShippingAmount float64
	TotalAmount    float64
	Items          []domain.SalesOrderItem
	CouponApplied  *string
}

type CartService interface {
	CalculateCart(ctx context.Context, items []domain.SalesOrderItem, couponCode string) (*CartCalculationResult, error)
}

type OrderService interface {
	// Creation & Lifecycle
	PlaceOrder(ctx context.Context, order *domain.SalesOrder) error
	GetOrder(ctx context.Context, orderNumber string) (*domain.SalesOrder, error)
	GetCustomerHistory(ctx context.Context, userID int, page, limit int) ([]domain.SalesOrder, error)

	// Actions
	CancelOrder(ctx context.Context, orderID int, reason string) error
	ProcessReturn(ctx context.Context, orderID int, items []domain.Return) error
	SubmitReview(ctx context.Context, review *domain.Review) error

	// Admin
	GetOrderList(ctx context.Context, filter dto.OrderFilterParams) ([]domain.SalesOrder, int64, error)
	SoftDeleteOrder(ctx context.Context, orderID int) error
	RestoreOrder(ctx context.Context, orderID int) error
	ForceDeleteOrder(ctx context.Context, orderID int) error

	GetByPOSSession(ctx context.Context, sessionID int) ([]domain.SalesOrder, error)
}

type POSService interface {
	// Session
	OpenSession(ctx context.Context, userID int, openingFloat float64) (*domain.POSSession, error)
	CloseSession(ctx context.Context, sessionID int, closingCashActual float64, note string) error
	GetActiveSession(ctx context.Context, userID int) (*domain.POSSession, error)
	GetSessionDetails(ctx context.Context, sessionID int) (*domain.POSSession, error) // X-Report

	// Cash Management
	RecordCashMove(ctx context.Context, sessionID int, amount float64, moveType domain.CashMoveType, reason string) error
	GetCashMoves(ctx context.Context, sessionID int) ([]domain.POSCashMove, error)

	// Sales Ops
	ScanProduct(ctx context.Context, barcode string) (*domain.ProductVariant, int, error) // Returns stock level
	SearchCustomer(ctx context.Context, query string) ([]domain.User, error)
	OverridePrice(ctx context.Context, variantID int, newPrice float64, managerPIN string) error
	VoidOrder(ctx context.Context, orderID int, managerPIN string) error

	// Utilities
	PrintReceipt(ctx context.Context, orderID int) error // Triggers printer command
}

// ==========================================
// 4. OPS & INVENTORY (THE LEDGER)
// ==========================================

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
	// Core Ledger
	ExecuteMovement(ctx context.Context, cmd StockMoveCmd) error
	GetStockLevel(ctx context.Context, variantID, locationID int) (int, error)
	GetMovements(ctx context.Context, variantID, locationID, page, limit int) ([]domain.StockMovement, error)

	// Operations
	TransferStock(ctx context.Context, variantID, qty, fromLocID, toLocID, userID int) error
	BulkAdjustStock(ctx context.Context, cmds []StockMoveCmd) error

	// Location Management
	GetLocations(ctx context.Context) ([]domain.InventoryLocation, error)
	CreateLocation(ctx context.Context, loc *domain.InventoryLocation) error

	// Data
	ExportStockSnapshot(ctx context.Context) ([]byte, error) // CSV
}

type AssemblyService interface {
	// Recipe Management
	GetRecipes(ctx context.Context) ([]domain.ProductRecipe, error)
	CreateRecipe(ctx context.Context, recipe *domain.ProductRecipe) error
	DeleteRecipe(ctx context.Context, recipeID int) error

	// Production
	AssembleKit(ctx context.Context, variantID, qty int, userID int) error
	Disassemble(ctx context.Context, variantID, qty int, userID int) error
	GetAssemblyLogs(ctx context.Context, page, limit int) ([]domain.StockAssembly, error)
}

type ProcurementService interface {
	GetPOs(ctx context.Context, page, limit int) ([]domain.PurchaseOrder, error)
	CreatePO(ctx context.Context, po *domain.PurchaseOrder) error
	ReceivePO(ctx context.Context, poID int, receivedItems map[int]int) error

	// Supplier Management
	GetSuppliers(ctx context.Context) ([]domain.Supplier, error)
	GetSupplier(ctx context.Context, id int) (*domain.Supplier, error)
	CreateSupplier(ctx context.Context, supplier *domain.Supplier) error
	UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error
	SoftDeleteSupplier(ctx context.Context, id int) error
	RestoreSupplier(ctx context.Context, id int) error
	ForceDeleteSupplier(ctx context.Context, id int) error
}

type FulfillmentService interface {
	GetQueue(ctx context.Context) ([]domain.SalesOrder, error) // Orders ready to ship
	PackOrder(ctx context.Context, orderID int) error
	ShipOrder(ctx context.Context, orderID int, carrier, trackingNumber string) error
}

// ==========================================
// 5. MARKETING & FINANCE
// ==========================================

type MarketingService interface {
	// Promotions
	GetPromotions(ctx context.Context) ([]domain.Promotion, error)
	CreatePromotion(ctx context.Context, promo *domain.Promotion) error
	UpdatePromotion(ctx context.Context, id int, req dto.UpdatePromotionRequest) error
	DeletePromotion(ctx context.Context, id int) error // Soft delete

	// CRM
	GetSegments(ctx context.Context) ([]string, error) // e.g., "Big Spenders", "Inactive"
	TriggerEmailCampaign(ctx context.Context, segment string, subject, body string) error
}

type FinanceService interface {
	GenerateInvoice(ctx context.Context, orderID int) (*domain.Invoice, error)
	RecordPayment(ctx context.Context, invoiceID int, amount float64, method string) error
	GetInvoicePDF(ctx context.Context, invoiceNumber string) (string, error) // Returns URL
}
