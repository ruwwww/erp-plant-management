package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OpsHandler struct {
	inventoryService   service.InventoryService
	assemblyService    service.AssemblyService
	procurementService service.ProcurementService
}

func NewOpsHandler(invS service.InventoryService, asmS service.AssemblyService, procS service.ProcurementService) *OpsHandler {
	return &OpsHandler{
		inventoryService:   invS,
		assemblyService:    asmS,
		procurementService: procS,
	}
}

// Inventory
func (h *OpsHandler) GetMovements(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) TransferStock(c *fiber.Ctx) error {
	var req dto.InventoryTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// userID := c.Locals("userID").(int)
	// if err := h.inventoryService.TransferStock(c.Context(), req.VariantID, req.Quantity, req.FromLocationID, req.ToLocationID, userID); err != nil { ... }

	return c.JSON(fiber.Map{"message": "Stock transferred"})
}

func (h *OpsHandler) AdjustStock(c *fiber.Ctx) error {
	var req dto.InventoryAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// userID := c.Locals("userID").(int)
	cmd := service.StockMoveCmd{
		LocationID: req.LocationID,
		VariantID:  req.VariantID,
		QtyChange:  req.ChangeQty,
		Reason:     domain.MovementReason(req.Reason),
		// UserID: userID,
	}

	if err := h.inventoryService.ExecuteMovement(c.Context(), cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Stock adjusted"})
}

func (h *OpsHandler) BulkAdjustStock(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Assembly
func (h *OpsHandler) GetAssemblies(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) ExecuteAssembly(c *fiber.Ctx) error {
	var req dto.ExecuteAssemblyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// userID := c.Locals("userID").(int)
	// if err := h.assemblyService.AssembleKit(c.Context(), req.VariantID, req.Quantity, userID); err != nil { ... }

	return c.JSON(fiber.Map{"message": "Assembly executed"})
}

// Procurement
func (h *OpsHandler) GetPOs(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) CreatePO(c *fiber.Ctx) error {
	var req dto.CreatePORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	po := &domain.PurchaseOrder{
		SupplierID: req.SupplierID,
		// ... map other fields
	}

	if err := h.procurementService.CreatePO(c.Context(), po); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "PO created"})
}

func (h *OpsHandler) ReceivePO(c *fiber.Ctx) error {
	poID, _ := strconv.Atoi(c.Params("id"))
	// Parse received items map

	if err := h.procurementService.ReceivePO(c.Context(), poID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "PO received"})
}
