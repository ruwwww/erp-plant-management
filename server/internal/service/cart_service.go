package service

import (
	"context"
	"server/internal/core/domain"
)

type CartServiceImpl struct {
	marketingService MarketingService
}

func NewCartService(marketingS MarketingService) CartService {
	return &CartServiceImpl{
		marketingService: marketingS,
	}
}

func (s *CartServiceImpl) CalculateCart(ctx context.Context, items []domain.SalesOrderItem, couponCode string) (*CartCalculationResult, error) {
	// Calculate subtotal
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.UnitPrice * float64(item.Quantity)
	}

	// Prepare data for promotion evaluation
	data := map[string]interface{}{
		"cart": map[string]interface{}{
			"total": subtotal,
			"items": items,
		},
		"discount": 0.0,
	}

	// Get active promotions
	promos, err := s.marketingService.GetPromotions(ctx)
	if err != nil {
		return nil, err
	}

	// Apply promotions (simplified: apply all that match, no exclusivity yet)
	totalDiscount := 0.0
	for _, promo := range promos {
		err := s.marketingService.ApplyPromotion(ctx, &promo, data)
		if err != nil {
			// Log error but continue
			continue
		}
		if disc, ok := data["discount"].(float64); ok {
			totalDiscount += disc
		}
	}

	// For now, assume no tax or shipping
	result := &CartCalculationResult{
		Subtotal:       subtotal,
		DiscountAmount: totalDiscount,
		TaxAmount:      0,
		ShippingAmount: 0,
		TotalAmount:    subtotal - totalDiscount,
		Items:          items,
	}

	if couponCode != "" {
		result.CouponApplied = &couponCode
		// TODO: Apply coupon discount
	}

	return result, nil
}
