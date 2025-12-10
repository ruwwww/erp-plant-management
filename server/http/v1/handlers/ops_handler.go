package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type OpsHandler struct {
	inventoryService   service.InventoryService
	assemblyService    service.AssemblyService
	procurementService service.ProcurementService
	fulfillmentService service.FulfillmentService
}

func NewOpsHandler(invS service.InventoryService, asmS service.AssemblyService, procS service.ProcurementService, fulS service.FulfillmentService) *OpsHandler {
	return &OpsHandler{
		inventoryService:   invS,
		assemblyService:    asmS,
		procurementService: procS,
		fulfillmentService: fulS,
	}
}

// Inventory
func (h *OpsHandler) GetLocations(c *fiber.Ctx) error {
	var filters dto.GetLocationsRequest
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	locations, err := h.inventoryService.GetLocations(c.Context(), &filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response DTOs
	responses := make([]dto.LocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = dto.LocationResponse{
			ID:        loc.ID,
			Name:      loc.Name,
			Code:      loc.Code,
			Type:      loc.Type,
			AddressID: loc.AddressID,
			IsActive:  loc.IsActive,
		}
		if loc.Address != nil {
			responses[i].Address = &dto.AddressResponse{
				ID:         loc.Address.ID,
				Line1:      loc.Address.Line1,
				Line2:      loc.Address.Line2,
				City:       loc.Address.City,
				State:      loc.Address.State,
				PostalCode: loc.Address.PostalCode,
				Country:    loc.Address.Country,
			}
		}
	}

	return c.JSON(responses)
}

func (h *OpsHandler) CreateLocation(c *fiber.Ctx) error {
	var req dto.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	loc := &domain.InventoryLocation{
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		AddressID: req.AddressID,
		IsActive:  isActive,
	}

	if err := h.inventoryService.CreateLocation(c.Context(), loc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.LocationResponse{
		ID:        loc.ID,
		Name:      loc.Name,
		Code:      loc.Code,
		Type:      loc.Type,
		AddressID: loc.AddressID,
		IsActive:  loc.IsActive,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *OpsHandler) UpdateLocation(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid location ID"})
	}

	var req dto.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.inventoryService.UpdateLocation(c.Context(), id, &req); err != nil {
		if err.Error() == "location not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Location updated successfully"})
}

func (h *OpsHandler) DeleteLocation(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid location ID"})
	}

	if err := h.inventoryService.DeleteLocation(c.Context(), id); err != nil {
		if strings.Contains(err.Error(), "cannot delete") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Location deleted successfully"})
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

func (h *OpsHandler) GetRecipes(c *fiber.Ctx) error {
	recipes, err := h.assemblyService.GetRecipes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(recipes)
}

func (h *OpsHandler) CreateRecipe(c *fiber.Ctx) error {
	var req dto.CreateRecipeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	recipe := &domain.ProductRecipe{
		ParentVariantID: req.ParentVariantID,
		ChildVariantID:  req.ChildVariantID,
		QuantityNeeded:  req.QuantityNeeded,
	}

	if err := h.assemblyService.CreateRecipe(c.Context(), recipe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(recipe)
}

func (h *OpsHandler) DeleteRecipe(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.assemblyService.DeleteRecipe(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Recipe deleted"})
}

func (h *OpsHandler) GetAssemblyLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	logs, err := h.assemblyService.GetAssemblyLogs(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"logs": logs})
}

func (h *OpsHandler) ExecuteAssembly(c *fiber.Ctx) error {
	var req dto.ExecuteAssemblyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("userID").(int)
	if err := h.assemblyService.AssembleKit(c.Context(), req.VariantID, req.Quantity, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Assembly executed"})
}

func (h *OpsHandler) DisassembleKit(c *fiber.Ctx) error {
	var req dto.ExecuteAssemblyRequest // Reuse DTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("userID").(int)
	if err := h.assemblyService.Disassemble(c.Context(), req.VariantID, req.Quantity, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Disassembly executed"})
}

// Procurement
func (h *OpsHandler) GetPurchaseOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	pos, err := h.procurementService.GetPOs(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pos)
}

func (h *OpsHandler) CreatePurchaseOrder(c *fiber.Ctx) error {
	var req dto.CreatePORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	po := &domain.PurchaseOrder{
		SupplierID: req.SupplierID,
		Notes:      &req.Notes,
		Status:     domain.PODraft,
	}

	if req.ExpectedAt != "" {
		// Parse time, but for simplicity, assume it's handled in service or use time.Parse
	}

	var subtotal float64
	for _, item := range req.Items {
		lineTotal := float64(item.Quantity) * item.UnitCost
		po.Items = append(po.Items, domain.PurchaseOrderItem{
			VariantID:       item.VariantID,
			QuantityOrdered: item.Quantity,
			UnitCost:        item.UnitCost,
			LineTotal:       lineTotal,
		})
		subtotal += lineTotal
	}

	po.SubtotalAmount = subtotal
	po.TotalAmount = subtotal // Add tax/shipping if needed

	// Generate PO number, e.g., PO-2023-001
	po.PONumber = "PO-2023-001" // TODO: generate properly

	if err := h.procurementService.CreatePO(c.Context(), po); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(po)
}

func (h *OpsHandler) ReceivePurchaseOrder(c *fiber.Ctx) error {
	poID, _ := strconv.Atoi(c.Params("id"))

	var req dto.ReceivePORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// First, get the PO to map variant IDs to item IDs
	po, err := h.procurementService.GetPOs(c.Context(), 1, 1000) // TODO: better way to get single PO
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Find the PO
	var targetPO *domain.PurchaseOrder
	for _, p := range po {
		if p.ID == poID {
			targetPO = &p
			break
		}
	}
	if targetPO == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "PO not found"})
	}

	receivedItems := make(map[int]int)
	for _, item := range req.Items {
		// Find item ID by variant ID
		for _, poItem := range targetPO.Items {
			if poItem.VariantID == item.VariantID {
				receivedItems[poItem.ID] = item.QtyReceived
				break
			}
		}
	}

	if err := h.procurementService.ReceivePO(c.Context(), poID, receivedItems); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "PO received"})
}

func (h *OpsHandler) GetSuppliers(c *fiber.Ctx) error {
	suppliers, err := h.procurementService.GetSuppliers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(suppliers)
}

func (h *OpsHandler) GetSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	supplier, err := h.procurementService.GetSupplier(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(supplier)
}

func (h *OpsHandler) CreateSupplier(c *fiber.Ctx) error {
	var req dto.CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	supplier := &domain.Supplier{
		Name:    req.Name,
		Email:   &req.Email,
		Phone:   &req.Phone,
		Address: &req.Address,
		Contact: &req.Contact,
	}

	if err := h.procurementService.CreateSupplier(c.Context(), supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(supplier)
}

func (h *OpsHandler) UpdateSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	supplier, err := h.procurementService.GetSupplier(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	supplier.Name = req.Name
	if req.Contact != "" {
		supplier.Contact = &req.Contact
	}
	if req.Email != "" {
		supplier.Email = &req.Email
	}
	if req.Phone != "" {
		supplier.Phone = &req.Phone
	}
	if req.Address != "" {
		supplier.Address = &req.Address
	}

	if err := h.procurementService.UpdateSupplier(c.Context(), supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(supplier)
}

func (h *OpsHandler) SoftDeleteSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.SoftDeleteSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier deleted"})
}

func (h *OpsHandler) RestoreSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.RestoreSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier restored"})
}

func (h *OpsHandler) ForceDeleteSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.ForceDeleteSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier force deleted"})
}

// Fulfillment
func (h *OpsHandler) GetFulfillmentQueue(c *fiber.Ctx) error {
	orders, err := h.fulfillmentService.GetQueue(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *OpsHandler) PackOrder(c *fiber.Ctx) error {
	orderID, _ := strconv.Atoi(c.Params("id"))
	if err := h.fulfillmentService.PackOrder(c.Context(), orderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order packed"})
}

func (h *OpsHandler) ShipOrder(c *fiber.Ctx) error {
	orderID, _ := strconv.Atoi(c.Params("id"))
	var req dto.ShipOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.fulfillmentService.ShipOrder(c.Context(), orderID, req.Carrier, req.TrackingNumber); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order shipped"})
}
