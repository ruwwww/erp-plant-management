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
func (h *OpsHandler) GetLocations(c *fiber.Ctx) error {
	locations, err := h.inventoryService.GetLocations(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(locations)
}

func (h *OpsHandler) CreateLocation(c *fiber.Ctx) error {
	var req dto.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	loc := &domain.InventoryLocation{
		Name:      req.Name,
		Type:      "WAREHOUSE", // Default, or add to DTO
		AddressID: req.AddressID,
		IsActive:  true,
	}

	if err := h.inventoryService.CreateLocation(c.Context(), loc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(loc)
}

func (h *OpsHandler) GetMovements(c *fiber.Ctx) error {
	variantID := c.QueryInt("variant_id", 0)
	locationID := c.QueryInt("location_id", 0)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	movements, err := h.inventoryService.GetMovements(c.Context(), variantID, locationID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"movements": movements})
}

func (h *OpsHandler) TransferStock(c *fiber.Ctx) error {
	var req dto.InventoryTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("userID").(int)
	if err := h.inventoryService.TransferStock(c.Context(), req.VariantID, req.Quantity, req.FromLocationID, req.ToLocationID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Stock transferred"})
}

func (h *OpsHandler) AdjustStock(c *fiber.Ctx) error {
	var req dto.InventoryAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("userID").(int)
	cmd := service.StockMoveCmd{
		LocationID: req.LocationID,
		VariantID:  req.VariantID,
		QtyChange:  req.ChangeQty,
		Reason:     domain.MovementReason(req.Reason),
		UserID:     userID,
	}

	if err := h.inventoryService.ExecuteMovement(c.Context(), cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Stock adjusted"})
}

func (h *OpsHandler) BulkAdjustStock(c *fiber.Ctx) error {
	var req dto.BulkAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("userID").(int)
	var cmds []service.StockMoveCmd

	for _, item := range req.Items {
		currentQty, err := h.inventoryService.GetStockLevel(c.Context(), item.VariantID, item.LocationID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get current stock"})
		}

		change := item.ActualQty - currentQty
		if change != 0 {
			cmds = append(cmds, service.StockMoveCmd{
				LocationID: item.LocationID,
				VariantID:  item.VariantID,
				QtyChange:  change,
				Reason:     domain.MovementReason(req.Reason),
				UserID:     userID,
			})
		}
	}

	if err := h.inventoryService.BulkAdjustStock(c.Context(), cmds); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Bulk stock adjusted"})
}

func (h *OpsHandler) ExportStockSnapshot(c *fiber.Ctx) error {
	data, err := h.inventoryService.ExportStockSnapshot(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=stock_snapshot.csv")
	return c.Send(data)
}

// Assembly
func (h *OpsHandler) GetRecipes(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) CreateRecipe(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) DeleteRecipe(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) GetAssemblyLogs(c *fiber.Ctx) error {
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

func (h *OpsHandler) DisassembleKit(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Procurement
func (h *OpsHandler) GetPurchaseOrders(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) CreatePurchaseOrder(c *fiber.Ctx) error {
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

func (h *OpsHandler) ReceivePurchaseOrder(c *fiber.Ctx) error {
	poID, _ := strconv.Atoi(c.Params("id"))
	// Parse received items map

	if err := h.procurementService.ReceivePO(c.Context(), poID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "PO received"})
}

// Fulfillment
func (h *OpsHandler) GetFulfillmentQueue(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) PackOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *OpsHandler) ShipOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
