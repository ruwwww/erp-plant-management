package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type POSHandler struct {
	posService service.POSService
}

func NewPOSHandler(posS service.POSService) *POSHandler {
	return &POSHandler{posService: posS}
}

// Session
func (h *POSHandler) GetActiveSession(c *fiber.Ctx) error {
	// userID := c.Locals("userID").(int)
	// session, err := h.posService.GetActiveSession(c.Context(), userID)
	return c.SendStatus(fiber.StatusOK)
}

func (h *POSHandler) OpenSession(c *fiber.Ctx) error {
	var req dto.OpenSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// userID := c.Locals("userID").(int)
	// session, err := h.posService.OpenSession(c.Context(), userID, req.OpeningCash)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Session opened"})
}

func (h *POSHandler) GetSessionDetails(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
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

	return c.JSON(fiber.Map{"message": "Session closed"})
}

// Cash Management
func (h *POSHandler) RecordCashMove(c *fiber.Ctx) error {
	sessionID, _ := strconv.Atoi(c.Params("id"))
	var req dto.CashMoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	moveType := domain.CashMoveType(req.Type) // Ensure type matches
	if err := h.posService.RecordCashMove(c.Context(), sessionID, req.Amount, moveType, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Cash move recorded"})
}

func (h *POSHandler) GetCashMoves(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Sales Operations
func (h *POSHandler) ScanProduct(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *POSHandler) SearchCustomer(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *POSHandler) OverridePrice(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *POSHandler) CreateOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Post-Sale
func (h *POSHandler) GetRecentOrders(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *POSHandler) PrintReceipt(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *POSHandler) VoidOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
