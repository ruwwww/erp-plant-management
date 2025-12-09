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

// Catalog
func (h *StoreHandler) GetProducts(c *fiber.Ctx) error {
	var query dto.ProductFilterQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	filter := service.ProductFilterParams{
		Search:       query.Search,
		CategorySlug: query.CategorySlug,
		MinPrice:     query.MinPrice,
		MaxPrice:     query.MaxPrice,
		Page:         query.Page,
		Limit:        query.Limit,
	}

	products, count, err := h.catalogService.GetProducts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  products,
		"total": count,
		"page":  query.Page,
		"limit": query.Limit,
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
	// TODO: Add GetCategories to CatalogService
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Cart & Checkout
func (h *StoreHandler) SyncCart(c *fiber.Ctx) error {
	var req dto.CartSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Map DTO to Domain
	var items []domain.SalesOrderItem
	for _, item := range req.Items {
		items = append(items, domain.SalesOrderItem{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
	}

	result, err := h.cartService.CalculateCart(c.Context(), items, "")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (h *StoreHandler) ApplyCoupon(c *fiber.Ctx) error {
	// TODO: Implement coupon logic
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *StoreHandler) CheckoutPreview(c *fiber.Ctx) error {
	var req dto.CheckoutPreviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var items []domain.SalesOrderItem
	for _, item := range req.Items {
		items = append(items, domain.SalesOrderItem{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		})
	}

	result, err := h.cartService.CalculateCart(c.Context(), items, req.CouponCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// In a real app, we would also calculate shipping based on address here

	return c.JSON(dto.CheckoutPreviewResponse{
		Subtotal:       result.Subtotal,
		DiscountAmount: result.DiscountAmount,
		TaxAmount:      result.TaxAmount,
		ShippingAmount: result.ShippingAmount,
		TotalAmount:    result.TotalAmount,
	})
}

func (h *StoreHandler) CheckoutPlace(c *fiber.Ctx) error {
	var req dto.CheckoutPlaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Construct Order
	order := &domain.SalesOrder{
		// CustomerID: ... (get from context if logged in)
		// ShippingAddressID: req.ShippingAddressID,
		// ...
	}

	// This is a simplified example. Real implementation needs to handle user context, address validation, etc.
	if err := h.orderService.PlaceOrder(c.Context(), order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Order placed successfully", "order_number": order.OrderNumber})
}

func (h *StoreHandler) PaymentWebhook(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
