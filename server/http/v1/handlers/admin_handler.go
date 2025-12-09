package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	catalogService service.CatalogService
	authService    service.AuthService
	userService    service.UserService
}

func NewAdminHandler(catalogS service.CatalogService, authS service.AuthService, userS service.UserService) *AdminHandler {
	return &AdminHandler{
		catalogService: catalogS,
		authService:    authS,
		userService:    userS,
	}
}

// Products
func (h *AdminHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product := &domain.Product{
		Name: req.Name,
		SKU:  req.SKU,
		Slug: req.Slug,
		// Description: &req.Description,
		// CategoryID:  &req.CategoryID,
		BasePrice: req.BasePrice,
		// ...
	}

	if err := h.catalogService.CreateProduct(c.Context(), product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product created", "id": product.ID})
}

func (h *AdminHandler) UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product := &domain.Product{
		ID:   id,
		Name: req.Name,
		// ...
	}

	if err := h.catalogService.UpdateProduct(c.Context(), product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product updated"})
}

func (h *AdminHandler) DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.catalogService.SoftDeleteProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted"})
}

// Users
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

func (h *AdminHandler) UpdateUserStatus(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *AdminHandler) AssignRole(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
