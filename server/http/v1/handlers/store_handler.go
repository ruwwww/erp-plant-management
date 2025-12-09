package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	catalogService service.CatalogService
	cartService    service.CartService
	orderService   service.OrderService
}

func NewStoreHandler(catalogS service.CatalogService, cartS service.CartService, orderS service.OrderService) *StoreHandler {
	return &StoreHandler{
		catalogService: catalogS,
		cartService:    cartS,
		orderService:   orderS,
	}
}

func (h *StoreHandler) GetProducts(c *fiber.Ctx) error {
	var query dto.ProductFilterQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 10
	}

	isActive := true

	filter := dto.ProductFilterParams{
		Search:       query.Search,
		CategorySlug: query.CategorySlug,
		MinPrice:     query.MinPrice,
		MaxPrice:     query.MaxPrice,
		IsActive:     &isActive,
		Page:         page,
		Limit:        limit,
	}

	products, total, err := h.catalogService.GetProducts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,

		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (h *StoreHandler) GetProductDetail(c *fiber.Ctx) error {
	slug := c.Params("slug")
	product, err := h.catalogService.GetProductDetail(c.Context(), slug)
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

func (h *StoreHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.catalogService.GetCategories(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}

// Helper to avoid code duplication
func mapCartItems(reqItems []dto.CartItemRequest) []domain.SalesOrderItem {
	var items []domain.SalesOrderItem
	for _, item := range reqItems {
		items = append(items, domain.SalesOrderItem{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
	}
	return items
}

func (h *StoreHandler) SyncCart(c *fiber.Ctx) error {
	var req dto.CartSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Use Helper
	items := mapCartItems(req.Items)

	result, err := h.cartService.CalculateCart(c.Context(), items, "")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *StoreHandler) ApplyCoupon(c *fiber.Ctx) error {
	var req dto.ApplyCouponRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Usually ApplyCoupon works on a stored cart, but since we are stateless (REST),
	// this endpoint might just validate the coupon and return details.
	// OR, the client just uses CheckoutPreview with the code.
	// For now, let's keep it unimplemented or make it a simple "Check Validity" endpoint.
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *StoreHandler) CheckoutPreview(c *fiber.Ctx) error {
	var req dto.CheckoutPreviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	items := mapCartItems(req.Items)

	result, err := h.cartService.CalculateCart(c.Context(), items, req.CouponCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *StoreHandler) CheckoutPlace(c *fiber.Ctx) error {
	var req dto.CheckoutPlaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	var userID *int
	if uid, ok := c.Locals("user_id").(float64); ok {
		id := int(uid)
		userID = &id
	}

	if userID == nil && req.GuestEmail == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email is required for guest checkout"})
	}

	items := mapCartItems(req.Items)
	cartResult, err := h.cartService.CalculateCart(c.Context(), items, req.CouponCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cart validation failed: " + err.Error()})
	}

	order := &domain.SalesOrder{
		CustomerID: nil,

		GuestEmail: &req.GuestEmail,
		Channel:    domain.ChannelWeb,

		ShippingAddressSnapshot: nil,

		TotalAmount:    cartResult.TotalAmount,
		SubtotalAmount: cartResult.Subtotal,
		DiscountAmount: cartResult.DiscountAmount,
		TaxAmount:      cartResult.TaxAmount,
		ShippingAmount: cartResult.ShippingAmount,

		Items:         cartResult.Items,
		PaymentMethod: req.PaymentMethod,
		Status:        domain.OrderDraft,
	}

	if err := h.orderService.PlaceOrder(c.Context(), order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":      "Order placed",
		"order_number": order.OrderNumber,
		"payment_url":  "http:dummy-payment-url.com",
	})
}

func (h *StoreHandler) PaymentWebhook(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusNotImplemented)

}
