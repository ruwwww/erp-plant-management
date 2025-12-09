package handlers

import (
	"math"
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	catalogService     service.CatalogService
	authService        service.AuthService
	userService        service.UserService
	procurementService service.ProcurementService
}

func NewAdminHandler(catalogS service.CatalogService, authS service.AuthService, userS service.UserService, procurementS service.ProcurementService) *AdminHandler {
	return &AdminHandler{
		catalogService:     catalogS,
		authService:        authS,
		userService:        userS,
		procurementService: procurementS,
	}
}

// Products
func (h *AdminHandler) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	filter := dto.ProductFilterParams{
		Search:       c.Query("q"),
		CategorySlug: c.Query("category"),
		Page:         page,
		Limit:        limit,
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		filter.MinPrice, _ = strconv.ParseFloat(minPrice, 64)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		filter.MaxPrice, _ = strconv.ParseFloat(maxPrice, 64)
	}

	if dateStr := c.Query("created_after"); dateStr != "" {
		// Parse ISO8601 / RFC3339 format
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			filter.CreatedAfter = &t
		}
	}

	// Parse 'created_before'
	if dateStr := c.Query("created_before"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			// Adjust to End of Day (23:59:59) if needed, or strict comparison
			filter.CreatedBefore = &t
		}
	}

	statusParam := c.Query("status")
	if statusParam == "active" {
		active := true
		filter.IsActive = &active
	} else if statusParam == "archived" {
		active := false
		filter.IsActive = &active
	}

	products, total, err := h.catalogService.GetProducts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_rows":  total,
			"total_pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func (h *AdminHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	// TODO: Add h.Validator.Struct(&req) here!

	product := req.ToDomain()

	if err := h.catalogService.CreateProduct(c.Context(), product); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"id": product.ID})
}

func (h *AdminHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.catalogService.UpdateProduct(c.Context(), id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (h *AdminHandler) SoftDeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.catalogService.SoftDeleteProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted"})
}

func (h *AdminHandler) RestoreProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.catalogService.RestoreProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product restored"})
}

func (h *AdminHandler) ForceDeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	return h.catalogService.ForceDeleteProduct(c.Context(), id)
}

// Variants
func (h *AdminHandler) GetVariants(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) UpdateVariants(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Media
func (h *AdminHandler) UploadMedia(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) LinkMedia(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) UnlinkMedia(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Users
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) CreateStaff(c *fiber.Ctx) error {
	var req dto.CreateStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user := &domain.User{
		Email:     req.Email,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Role:      domain.UserRole(req.Role),
	}

	if err := h.authService.RegisterStaff(c.Context(), user, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Staff created", "id": user.ID})
}

func (h *AdminHandler) GetUserDetail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	user, err := h.userService.GetUserDetail(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userService.GetUserDetail(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if req.FirstName != "" {
		user.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		user.LastName = &req.LastName
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := h.userService.UpdateProfile(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated"})
}

func (h *AdminHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateUserStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.userService.UpdateUserStatus(c.Context(), id, req.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User status updated"})
}

func (h *AdminHandler) AdminResetPassword(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.AdminResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.authService.ResetPassword(c.Context(), id, req.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password reset"})
}

func (h *AdminHandler) AssignRoles(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// CRM
func (h *AdminHandler) GetSegments(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) TriggerEmailCampaign(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) CreateSupplier(c *fiber.Ctx) error {
	var req dto.CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	supplier := &domain.Supplier{
		Name: req.Name,
	}
	if req.ContactName != "" {
		supplier.ContactName = &req.ContactName
	}
	if req.Email != "" {
		supplier.Email = &req.Email
	}
	if req.Phone != "" {
		supplier.Phone = &req.Phone
	}

	if err := h.procurementService.CreateSupplier(c.Context(), supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Supplier created", "id": supplier.ID})
}

func (h *AdminHandler) GetSuppliers(c *fiber.Ctx) error {
	suppliers, err := h.procurementService.GetSuppliers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(suppliers)
}

func (h *AdminHandler) UpdateSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	supplier, err := h.procurementService.GetSupplier(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.ContactName != "" {
		supplier.ContactName = &req.ContactName
	}
	if req.Email != "" {
		supplier.Email = &req.Email
	}
	if req.Phone != "" {
		supplier.Phone = &req.Phone
	}

	if err := h.procurementService.UpdateSupplier(c.Context(), supplier); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Supplier updated"})
}

func (h *AdminHandler) SoftDeleteSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.SoftDeleteSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier deleted"})
}

func (h *AdminHandler) RestoreSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.RestoreSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier restored"})
}

func (h *AdminHandler) ForceDeleteSupplier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.procurementService.ForceDeleteSupplier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Supplier permanently deleted"})
}

// Promotions
func (h *AdminHandler) GetPromotions(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) CreatePromotion(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Data
func (h *AdminHandler) ExportProducts(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) ImportProducts(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) ImportInventoryAdjustments(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
