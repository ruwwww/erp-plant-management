package handlers

import (
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService    service.UserService
	orderService   service.OrderService
	financeService service.FinanceService
}

func NewUserHandler(userS service.UserService, orderS service.OrderService, financeS service.FinanceService) *UserHandler {
	return &UserHandler{
		userService:    userS,
		orderService:   orderS,
		financeService: financeS,
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	user, err := h.userService.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user := &domain.User{
		ID:        userID,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Phone:     &req.Phone,
		// Bio: req.Bio, // If User struct has Bio
	}

	if err := h.userService.UpdateProfile(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Profile updated"})
}

// Addresses
func (h *UserHandler) GetAddresses(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	addresses, err := h.userService.GetAddresses(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(addresses)
}

func (h *UserHandler) CreateAddress(c *fiber.Ctx) error {
	var req dto.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	addr := &domain.Address{
		Line1:      req.Line1,
		Line2:      &req.Line2,
		City:       req.City,
		State:      &req.State,
		PostalCode: &req.PostalCode,
		Country:    req.Country,
		Latitude:   &req.Latitude,
		Longitude:  &req.Longitude,
	}

	if err := h.userService.AddAddress(c.Context(), userID, addr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Address created", "id": addr.ID})
}

func (h *UserHandler) UpdateAddress(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	addressID, _ := strconv.Atoi(c.Params("id"))

	var req dto.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get existing address to update
	addresses, err := h.userService.GetAddresses(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var addr *domain.Address
	for i := range addresses {
		if addresses[i].ID == addressID {
			addr = &addresses[i]
			break
		}
	}
	if addr == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found"})
	}

	// Update fields if provided
	if req.Line1 != "" {
		addr.Line1 = req.Line1
	}
	if req.Line2 != "" {
		addr.Line2 = &req.Line2
	}
	if req.City != "" {
		addr.City = req.City
	}
	if req.State != "" {
		addr.State = &req.State
	}
	if req.PostalCode != "" {
		addr.PostalCode = &req.PostalCode
	}
	if req.Country != "" {
		addr.Country = req.Country
	}
	if req.Latitude != 0 {
		addr.Latitude = &req.Latitude
	}
	if req.Longitude != 0 {
		addr.Longitude = &req.Longitude
	}

	if err := h.userService.UpdateAddress(c.Context(), addr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Address updated"})
}

func (h *UserHandler) DeleteAddress(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	addressID, _ := strconv.Atoi(c.Params("id"))

	if err := h.userService.DeleteAddress(c.Context(), userID, addressID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Address deleted"})
}

func (h *UserHandler) SetDefaultAddress(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	addressID, _ := strconv.Atoi(c.Params("id"))

	var req dto.SetDefaultAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.userService.SetDefaultAddress(c.Context(), userID, addressID, req.IsBilling); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Default address set"})
}

// Orders
func (h *UserHandler) GetOrders(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) GetOrderDetail(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) CancelOrder(c *fiber.Ctx) error {
	orderID, _ := strconv.Atoi(c.Params("id"))
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.orderService.CancelOrder(c.Context(), orderID, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Order cancelled"})
}

func (h *UserHandler) RequestReturn(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) SubmitReview(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Documents
func (h *UserHandler) GetInvoices(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) DownloadInvoicePDF(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

// Wishlist
func (h *UserHandler) GetWishlist(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) ToggleWishlist(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
