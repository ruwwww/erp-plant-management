package repository

import (
	"context"
	"server/internal/core/domain"
	"server/internal/dto"
)

// 1. Base Generic Interface
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id any) (*T, error)
	FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
	Find(ctx context.Context, condition interface{}, args ...interface{}) ([]T, error)
	FindAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error
}

// 2. Auth & User
type UserRepository interface {
	Repository[domain.User]
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	ForceDelete(ctx context.Context, id int) error
}

// 3. Catalog (Product & Category)
type ProductRepository interface {
	Repository[domain.Product]
	// GetFullProduct loads Variants, Category, and Supplier
	GetFullProduct(ctx context.Context, slug string) (*domain.Product, error)
	// Search supports complex filtering
	Search(ctx context.Context, filter dto.ProductFilterParams) ([]domain.Product, int64, error)
	// SoftDelete soft deletes a product
	SoftDelete(ctx context.Context, id int) error
	// Restore restores a soft-deleted product
	Restore(ctx context.Context, id int) error
	// ForceDelete permanently deletes a product
	ForceDelete(ctx context.Context, id int) error
}

type VariantRepository interface {
	Repository[domain.ProductVariant]
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	ForceDelete(ctx context.Context, id int) error
}

type TagRepository interface {
	Repository[domain.Tag]
	FindBySlug(ctx context.Context, slug string) (*domain.Tag, error)
}

type CategoryRepository interface {
	Repository[domain.Category]
	GetTree(ctx context.Context) ([]domain.Category, error)
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	ForceDelete(ctx context.Context, id int) error
}

// 4. Sales & POS
type OrderRepository interface {
	Repository[domain.SalesOrder]
	// GetFullOrder loads Items, Customer, and Payment info
	GetFullOrder(ctx context.Context, orderNumber string) (*domain.SalesOrder, error)
	GetFullOrderByID(ctx context.Context, id int) (*domain.SalesOrder, error)
	Search(ctx context.Context, filter dto.OrderFilterParams) ([]domain.SalesOrder, int64, error)
	GetByPOSSession(ctx context.Context, sessionID int) ([]domain.SalesOrder, error)
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	ForceDelete(ctx context.Context, id int) error
}

type POSSessionRepository interface {
	Repository[domain.POSSession]
	FindActiveSession(ctx context.Context, userID int) (*domain.POSSession, error)
}

type CashMoveRepository interface {
	Repository[domain.POSCashMove]
	GetBySession(ctx context.Context, sessionID int) ([]domain.POSCashMove, error)
}

// 5. Inventory & Ops (The Ledger)
type InventoryRepository interface {
	Repository[domain.Stock]

	// GetStock finds specific stock for a variant at a location
	GetStock(ctx context.Context, variantID, locationID int) (*domain.Stock, error)

	// UpdateStockAtomic safely increments/decrements quantity to prevent race conditions
	UpdateStockAtomic(ctx context.Context, variantID, locationID, qtyChange int) error
}

type MovementRepository interface {
	Repository[domain.StockMovement]
	GetHistory(ctx context.Context, variantID, locationID int, limit int) ([]domain.StockMovement, error)
}

type AssemblyRepository interface {
	Repository[domain.StockAssembly]
	GetRecipe(ctx context.Context, variantID int) ([]domain.ProductRecipe, error)
}

type SupplierRepository interface {
	Repository[domain.Supplier]
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	ForceDelete(ctx context.Context, id int) error
}

type PurchaseOrderRepository interface {
	Repository[domain.PurchaseOrder]
	GetFullPO(ctx context.Context, id int) (*domain.PurchaseOrder, error)
}

// 6. Finance
type InvoiceRepository interface {
	Repository[domain.Invoice]
	FindByOrder(ctx context.Context, orderID int) (*domain.Invoice, error)
}

// 7. Marketing & Promotions
type PromotionRepository interface {
	Repository[domain.Promotion]
	GetActivePromotions(ctx context.Context) ([]domain.Promotion, error)
}
