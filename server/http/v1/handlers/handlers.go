package handlers

import "github.com/gofiber/fiber/v2"

// --- Auth Handler ---
type AuthHandler struct { /* AuthService */
}

func (h *AuthHandler) Login(c *fiber.Ctx) error   { return c.SendStatus(200) }
func (h *AuthHandler) Refresh(c *fiber.Ctx) error { return c.SendStatus(200) }
func (h *AuthHandler) Logout(c *fiber.Ctx) error  { return c.SendStatus(200) }

// --- Store Handler (Public) ---
type StoreHandler struct { /* CatalogService, CartService */
}

func (h *StoreHandler) GetProducts(c *fiber.Ctx) error      { return c.SendStatus(200) }
func (h *StoreHandler) GetProductDetail(c *fiber.Ctx) error { return c.SendStatus(200) }
func (h *StoreHandler) SyncCart(c *fiber.Ctx) error         { return c.SendStatus(200) }
func (h *StoreHandler) CheckoutPreview(c *fiber.Ctx) error  { return c.SendStatus(200) }
func (h *StoreHandler) CheckoutPlace(c *fiber.Ctx) error    { return c.SendStatus(200) }

// --- User Handler (My Account) ---
type UserHandler struct { /* UserService */
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error { return c.SendStatus(200) }
func (h *UserHandler) GetOrders(c *fiber.Ctx) error  { return c.SendStatus(200) }

// --- POS Handler ---
type POSHandler struct { /* POSService, OrderService */
}

func (h *POSHandler) OpenSession(c *fiber.Ctx) error  { return c.SendStatus(200) }
func (h *POSHandler) CloseSession(c *fiber.Ctx) error { return c.SendStatus(200) }
func (h *POSHandler) CreateOrder(c *fiber.Ctx) error  { return c.SendStatus(200) }
func (h *POSHandler) ScanProduct(c *fiber.Ctx) error  { return c.SendStatus(200) }

// --- Ops Handler (Inventory) ---
type OpsHandler struct { /* InventoryService, AssemblyService */
}

func (h *OpsHandler) GetMovements(c *fiber.Ctx) error    { return c.SendStatus(200) }
func (h *OpsHandler) TransferStock(c *fiber.Ctx) error   { return c.SendStatus(200) }
func (h *OpsHandler) ExecuteAssembly(c *fiber.Ctx) error { return c.SendStatus(200) }
func (h *OpsHandler) ReceivePO(c *fiber.Ctx) error       { return c.SendStatus(200) }

// --- Admin Handler ---
type AdminHandler struct { /* ProductService, UserService */
}

func (h *AdminHandler) CreateProduct(c *fiber.Ctx) error { return c.SendStatus(201) }
func (h *AdminHandler) CreateStaff(c *fiber.Ctx) error   { return c.SendStatus(201) }
