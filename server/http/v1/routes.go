package v1

import (
	"server/http/middleware"
	"server/http/v1/handlers"
	"server/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers ALL API endpoints based on the Master Spec
func SetupRoutes(
	app *fiber.App,
	module string,
	authH *handlers.AuthHandler,
	storeH *handlers.StoreHandler,
	userH *handlers.UserHandler,
	posH *handlers.POSHandler,
	opsH *handlers.OpsHandler,
	adminH *handlers.AdminHandler,
) {
	api := app.Group("/api/v1")

	// =====================================
	// 1. AUTH (Public) - Always Available
	// =====================================
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/refresh", authH.Refresh)
	auth.Post("/logout", middleware.Protect(), authH.Logout)

	// Forgot Password Flow
	auth.Post("/password-reset", authH.RequestPasswordReset)
	auth.Post("/password-reset/confirm", authH.ConfirmPasswordReset)

	// =====================================
	// 2. STOREFRONT (Public / Semi-Public)
	// =====================================
	if module == "store" || module == "" {
		store := api.Group("/store")

		// Catalog
		store.Get("/catalog/products", storeH.GetProducts)
		store.Get("/catalog/categories", storeH.GetCategories)
		store.Get("/catalog/products/:slug", storeH.GetProductDetail)

		// Cart & Checkout
		store.Post("/cart/sync", storeH.SyncCart) // Sync Guest Cart
		store.Post("/cart/coupons", storeH.ApplyCoupon)
		store.Post("/checkout/preview", storeH.CheckoutPreview)
		store.Post("/checkout/place", storeH.CheckoutPlace) // Might reserve stock

		// Webhooks (Third Party)
		store.Post("/webhooks/payment", storeH.PaymentWebhook)

		// =====================================
		// 3. USER DASHBOARD (Protected: Any Logged In User)
		// =====================================
		me := api.Group("/me", middleware.Protect())

		// Profile
		me.Get("/profile", userH.GetProfile)
		me.Put("/profile", userH.UpdateProfile)

		// Addresses
		me.Get("/addresses", userH.GetAddresses)
		me.Post("/addresses", userH.CreateAddress)
		me.Put("/addresses/:id", userH.UpdateAddress)
		me.Delete("/addresses/:id", userH.DeleteAddress)
		me.Post("/addresses/:id/default", userH.SetDefaultAddress)

		// Orders & History
		me.Get("/orders", userH.GetOrders)
		me.Get("/orders/:number", userH.GetOrderDetail)
		me.Post("/orders/:number/cancel", userH.CancelOrder)
		me.Post("/orders/:number/return", userH.RequestReturn)
		me.Post("/orders/:number/review", userH.SubmitReview)

		// Documents
		me.Get("/invoices", userH.GetInvoices)
		me.Get("/invoices/:number/pdf", userH.DownloadInvoicePDF)

		// Wishlist
		me.Get("/wishlist", userH.GetWishlist)
		me.Post("/wishlist/:variant_id", userH.ToggleWishlist)
	}

	// =====================================
	// 4. POS SYSTEM (Protected: STAFF, MANAGER, ADMIN)
	// =====================================
	if module == "pos" || module == "" {
		pos := api.Group("/pos",
			middleware.Protect(),
			middleware.Authorize(domain.RoleStaff, domain.RoleManager, domain.RoleAdmin),
		)

		// Session Management
		pos.Get("/sessions/active", posH.GetActiveSession)
		pos.Post("/sessions/open", posH.OpenSession)
		pos.Get("/sessions/:id", posH.GetSessionDetails)   // X-Report
		pos.Post("/sessions/:id/close", posH.CloseSession) // Z-Report

		// Cash Management (Safe Drops)
		pos.Post("/cash-moves", posH.RecordCashMove)
		pos.Get("/cash-moves", posH.GetCashMoves)

		// Sales Operations
		pos.Get("/products/scan", posH.ScanProduct)
		pos.Post("/customers/search", posH.SearchCustomer)
		pos.Post("/cart/override-price", posH.OverridePrice) // Manager PIN usually needed
		pos.Post("/orders/create", posH.CreateOrder)

		// Post-Sale Actions
		pos.Get("/orders/recent", posH.GetRecentOrders)
		pos.Post("/orders/:id/print", posH.PrintReceipt)
		pos.Post("/orders/:id/void", posH.VoidOrder)
	}

	// =====================================
	// 5. OPS & INVENTORY (Protected: MANAGER, ADMIN)
	// =====================================
	if module == "internal" || module == "" {
		ops := api.Group("/ops",
			middleware.Protect(),
			middleware.Authorize(domain.RoleManager, domain.RoleAdmin),
		)

		// Inventory Ledger
		ops.Get("/inventory/locations", opsH.GetLocations)
		ops.Post("/inventory/locations", opsH.CreateLocation)
		ops.Get("/inventory/movements", opsH.GetMovements) // The Audit Trail
		ops.Post("/inventory/transfer", opsH.TransferStock)
		ops.Post("/inventory/adjust", opsH.AdjustStock)
		ops.Post("/inventory/bulk-adjust", opsH.BulkAdjustStock)
		ops.Get("/inventory/snapshot", opsH.ExportStockSnapshot)

		// Manufacturing / Assembly
		ops.Get("/assembly/recipes", opsH.GetRecipes)
		ops.Post("/assembly/recipes", opsH.CreateRecipe)
		ops.Delete("/assembly/recipes/:id", opsH.DeleteRecipe)
		ops.Get("/assembly/logs", opsH.GetAssemblyLogs)
		ops.Post("/assembly/execute", opsH.ExecuteAssembly) // The "Make" Button
		ops.Post("/assembly/disassemble", opsH.DisassembleKit)

		// Procurement
		ops.Get("/procurement/po", opsH.GetPurchaseOrders)
		ops.Post("/procurement/po", opsH.CreatePurchaseOrder)
		ops.Post("/procurement/po/:id/receive", opsH.ReceivePurchaseOrder) // Inbound Stock

		// Fulfillment (Shipping)
		ops.Get("/fulfillment/queue", opsH.GetFulfillmentQueue)
		ops.Post("/fulfillment/:id/pack", opsH.PackOrder)
		ops.Post("/fulfillment/:id/ship", opsH.ShipOrder)

		// =====================================
		// 6. ADMIN (Protected: ADMIN ONLY)
		// =====================================
		admin := api.Group("/admin",
			middleware.Protect(),
			middleware.Authorize(domain.RoleAdmin),
		)

		// Product Management
		admin.Get("/products", adminH.GetProducts)
		admin.Post("/products", adminH.CreateProduct)
		admin.Put("/products/:id", adminH.UpdateProduct)
		admin.Delete("/products/:id", adminH.SoftDeleteProduct)
		admin.Post("/products/:id/restore", adminH.RestoreProduct)
		admin.Delete("/products/:id/force", adminH.ForceDeleteProduct)

		// Variants
		admin.Get("/products/:id/variants", adminH.GetVariants)
		admin.Put("/products/:id/variants", adminH.UpdateVariants)

		// Media Pipeline
		admin.Post("/media/upload", adminH.UploadMedia)
		admin.Post("/media/link", adminH.LinkMedia)
		admin.Delete("/media/link", adminH.UnlinkMedia)

		// User Management (HR)
		admin.Get("/users", adminH.GetUsers)
		admin.Post("/users", adminH.CreateStaff)
		admin.Get("/users/:id", adminH.GetUserDetail)
		admin.Put("/users/:id", adminH.UpdateUser)
		admin.Put("/users/:id/status", adminH.UpdateUserStatus) // Ban
		admin.Post("/users/:id/reset-password", adminH.AdminResetPassword)
		admin.Post("/users/:id/roles", adminH.AssignRoles)

		// Supliers
		admin.Get("/suppliers", adminH.GetSuppliers)
		admin.Post("/suppliers", adminH.CreateSupplier)
		admin.Put("/suppliers/:id", adminH.UpdateSupplier)
		admin.Delete("/suppliers/:id", adminH.SoftDeleteSupplier)        // Soft Delete
		admin.Post("/suppliers/:id/restore", adminH.RestoreSupplier)     // Restore
		admin.Delete("/suppliers/:id/force", adminH.ForceDeleteSupplier) // Hard Delete

		// Tags
		admin.Get("/tags", adminH.GetTags)
		admin.Post("/tags", adminH.CreateTag)
		admin.Put("/products/:id/tags", adminH.UpdateProductTags)

		// CRM
		admin.Get("/customers/segments", adminH.GetSegments)
		admin.Post("/customers/email", adminH.TriggerEmailCampaign)

		// Promotions
		admin.Get("/promotions", adminH.GetPromotions)
		admin.Post("/promotions", adminH.CreatePromotion)
		admin.Put("/promotions/:id", adminH.UpdatePromotion)
		admin.Delete("/promotions/:id", adminH.DeletePromotion)

		// Data Import/Export
		admin.Get("/data/products/export", adminH.ExportProducts)
		admin.Post("/data/products/import", adminH.ImportProducts)
		admin.Post("/data/ops/inventory/import", adminH.ImportInventoryAdjustments) // Ledger Import
	}
}
