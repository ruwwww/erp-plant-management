package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type POSHandler struct {
	posService   service.POSService
	orderService service.OrderService
}

func NewPOSHandler(posS service.POSService, orderS service.OrderService) *POSHandler {
	return &POSHandler{
		posService:   posS,
		orderService: orderS,
	}
}

func getUserID(c *fiber.Ctx) int {

	if id, ok := c.Locals("user_id").(float64); ok {
		return int(id)
	}

	if id, ok := c.Locals("user_id").(int); ok {
		return id
	}
	return 0
}

func (h *POSHandler) GetActiveSession(c *fiber.Ctx) error {
	userID := getUserID(c)

	session, err := h.posService.GetActiveSession(c.Context(), userID)
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No active session found"})
	}

	return c.JSON(session)
}

func (h *POSHandler) OpenSession(c *fiber.Ctx) error {
	userID := getUserID(c)

	var req dto.OpenSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	session, err := h.posService.OpenSession(c.Context(), userID, req.OpeningCash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}

func (h *POSHandler) GetSessionDetails(c *fiber.Ctx) error {

	sessionID, _ := strconv.Atoi(c.Params("id"))

	session, err := h.posService.GetSessionDetails(c.Context(), sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(session)
}

func (h *POSHandler) CloseSession(c *fiber.Ctx) error {

	sessionID, _ := strconv.Atoi(c.Params("id"))

	var req dto.CloseSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.posService.CloseSession(c.Context(), sessionID, req.ClosingCashActual, req.Note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Session closed successfully"})
}

func (h *POSHandler) RecordCashMove(c *fiber.Ctx) error {

	var req dto.CashMoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	moveType := domain.CashMoveType(req.Type)

	if req.POSSessionID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "pos_session_id is required"})
	}

	if err := h.posService.RecordCashMove(c.Context(), req.POSSessionID, req.Amount, moveType, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Cash move recorded"})
}

func (h *POSHandler) GetCashMoves(c *fiber.Ctx) error {

	sessionID, _ := strconv.Atoi(c.Query("session_id"))

	if sessionID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "session_id query param required"})
	}

	moves, err := h.posService.GetCashMoves(c.Context(), sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(moves)
}

func (h *POSHandler) ScanProduct(c *fiber.Ctx) error {

	barcode := c.Query("barcode")
	if barcode == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Barcode is required"})
	}

	variant, stockLevel, err := h.posService.ScanProduct(c.Context(), barcode)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(fiber.Map{
		"variant": variant,
		"stock":   stockLevel,
	})
}

func (h *POSHandler) SearchCustomer(c *fiber.Ctx) error {
	var req dto.CustomerSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	customers, err := h.posService.SearchCustomer(c.Context(), req.Query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(customers)
}

func (h *POSHandler) OverridePrice(c *fiber.Ctx) error {
	var req dto.OverridePriceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := h.posService.OverridePrice(c.Context(), req.VariantID, req.NewPrice, req.ManagerPIN); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Price override approved"})
}

func (h *POSHandler) CreateOrder(c *fiber.Ctx) error {
	userID := getUserID(c)
	var req dto.CreatePOSOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	items := make([]domain.SalesOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.SalesOrderItem{
			VariantID: item.VariantID,
			Quantity:  item.Quantity,
		}
	}

	order := &domain.SalesOrder{
		POSSessionID:   &req.POSSessionID,
		CustomerID:     req.CustomerID,
		Channel:        domain.ChannelPOS,
		Items:          items,
		PaymentMethod:  req.PaymentMethod,
		Status:         domain.OrderCompleted,
		PaymentStatus:  domain.PaymentPaid,
		ShipmentStatus: domain.ShipmentDelivered,
		CreatedBy:      &userID,
	}

	if err := h.orderService.PlaceOrder(c.Context(), order); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":      "Order created",
		"order_number": order.OrderNumber,
		"order_id":     order.ID,
	})
}

func (h *POSHandler) PrintReceipt(c *fiber.Ctx) error {
	orderID, _ := strconv.Atoi(c.Params("id"))

	if err := h.posService.PrintReceipt(c.Context(), orderID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Print job sent"})
}

func (h *POSHandler) VoidOrder(c *fiber.Ctx) error {
	orderID, _ := strconv.Atoi(c.Params("id"))
	var req struct {
		ManagerPIN string `json:"manager_pin"`
	}
	c.BodyParser(&req)

	if err := h.posService.VoidOrder(c.Context(), orderID, req.ManagerPIN); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Order voided and stock returned"})
}

func (h *POSHandler) GetRecentOrders(c *fiber.Ctx) error {
	userID := getUserID(c)

	session, err := h.posService.GetActiveSession(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active session found. Please open a session first.",
		})
	}

	orders, err := h.orderService.GetByPOSSession(c.Context(), session.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"session_id": session.ID,
		"count":      len(orders),
		"data":       orders,
	})
}
