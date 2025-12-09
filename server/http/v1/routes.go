package v1

import (
	"server/http/middleware"
	"server/http/v1/handlers"
	"server/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all API endpoints
func SetupRoutes(app *fiber.App) {
	// Initialize Handlers (Dependency Injection would happen here in main.go, passed in)
	// For this example, we assume they are passed or instantiated here.
	authH := &handlers.AuthHandler{}
	storeH := &handlers.StoreHandler{}
	userH := &handlers.UserHandler{}
	posH := &handlers.POSHandler{}
	opsH := &handlers.OpsHandler{}
	adminH := &handlers.AdminHandler{}

	api := app.Group("/api/v1")

	// =====================================
	// 1. AUTH (Public)
	// =====================================
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/refresh", authH.Refresh)
	auth.Post("/logout", middleware.Protect(), authH.Logout)

	// =====================================
	// 2. STOREFRONT (Public / Semi-Public)
	// =====================================
	store := api.Group("/store")
	store.Get("/catalog/products", storeH.GetProducts)
	store.Get("/catalog/products/:slug", storeH.GetProductDetail)
	store.Post("/cart/sync", storeH.SyncCart) // Guest allowed
	store.Post("/checkout/preview", storeH.CheckoutPreview)
	store.Post("/checkout/place", storeH.CheckoutPlace) // Might require Auth depending on logic

	// =====================================
	// 3. USER PROFILE (Protected: Any Logged In User)
	// =====================================
	me := api.Group("/me", middleware.Protect())
	me.Get("/profile", userH.GetProfile)
	me.Get("/orders", userH.GetOrders)

	// =====================================
	// 4. POS SYSTEM (Protected: STAFF, MANAGER, ADMIN)
	// =====================================
	pos := api.Group("/pos",
		middleware.Protect(),
		middleware.Authorize(domain.RoleStaff, domain.RoleManager, domain.RoleAdmin),
	)

	// Sessions
	pos.Post("/sessions/open", posH.OpenSession)
	pos.Post("/sessions/:id/close", posH.CloseSession)

	// Transactions
	pos.Get("/products/scan", posH.ScanProduct)
	pos.Post("/orders/create", posH.CreateOrder)

	// =====================================
	// 5. OPS & INVENTORY (Protected: MANAGER, ADMIN)
	// =====================================
	ops := api.Group("/ops",
		middleware.Protect(),
		middleware.Authorize(domain.RoleManager, domain.RoleAdmin),
	)

	// Inventory Ledger
	ops.Get("/inventory/movements", opsH.GetMovements)
	ops.Post("/inventory/transfer", opsH.TransferStock)

	// Assembly (Manufacturing)
	ops.Post("/assembly/execute", opsH.ExecuteAssembly)

	// Procurement
	ops.Post("/procurement/po/:id/receive", opsH.ReceivePO)

	// =====================================
	// 6. ADMIN (Protected: ADMIN ONLY)
	// =====================================
	admin := api.Group("/admin",
		middleware.Protect(),
		middleware.Authorize(domain.RoleAdmin),
	)

	// Product Management
	admin.Post("/products", adminH.CreateProduct)

	// User Management (HR)
	admin.Post("/users", adminH.CreateStaff)
}
