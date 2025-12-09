package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/dto"
	"server/internal/repository"
	"time"

	"github.com/expr-lang/expr"
	"gorm.io/gorm"
)

type marketingService struct {
	promoRepo repository.PromotionRepository
}

func NewMarketingService(db *gorm.DB) MarketingService {
	return &marketingService{
		promoRepo: repository.NewPromotionRepository(db),
	}
}

func (s *marketingService) GetPromotions(ctx context.Context) ([]domain.Promotion, error) {
	return s.promoRepo.FindAll(ctx)
}

func (s *marketingService) CreatePromotion(ctx context.Context, promo *domain.Promotion) error {
	// Basic validation
	if promo.Name == "" {
		return errors.New("promotion name is required")
	}
	if promo.Conditions == nil {
		promo.Conditions = map[string]interface{}{}
	}
	if promo.Actions == nil {
		promo.Actions = map[string]interface{}{}
	}
	return s.promoRepo.Create(ctx, promo)
}

func (s *marketingService) UpdatePromotion(ctx context.Context, id int, req dto.UpdatePromotionRequest) error {
	promo, err := s.promoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if provided
	if req.Name != nil {
		promo.Name = *req.Name
	}
	if req.Code != nil {
		promo.Code = req.Code
	}
	if req.Description != nil {
		promo.Description = req.Description
	}
	if req.IsActive != nil {
		promo.IsActive = *req.IsActive
	}
	if req.StartsAt != nil {
		start, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			return errors.New("invalid starts_at format")
		}
		promo.StartsAt = &start
	}
	if req.EndsAt != nil {
		end, err := time.Parse(time.RFC3339, *req.EndsAt)
		if err != nil {
			return errors.New("invalid ends_at format")
		}
		promo.EndsAt = &end
	}
	if req.Conditions != nil {
		promo.Conditions = req.Conditions
	}
	if req.Actions != nil {
		promo.Actions = req.Actions
	}
	if req.TotalUsageLimit != nil {
		promo.TotalUsageLimit = req.TotalUsageLimit
	}

	return s.promoRepo.Update(ctx, promo)
}

func (s *marketingService) DeletePromotion(ctx context.Context, id int) error {
	promo, err := s.promoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete: set is_active = false and ends_at = now
	now := time.Now()
	promo.IsActive = false
	promo.EndsAt = &now

	return s.promoRepo.Update(ctx, promo)
}

// EvaluatePromotion checks if a promotion's conditions are met for given data
func (s *marketingService) EvaluatePromotion(ctx context.Context, promo *domain.Promotion, data map[string]interface{}) (bool, error) {
	// Assume Conditions has "expression" key with expr string
	if exprStr, ok := promo.Conditions["expression"].(string); ok {
		program, err := expr.Compile(exprStr, expr.Env(data))
		if err != nil {
			return false, err
		}
		result, err := expr.Run(program, data)
		if err != nil {
			return false, err
		}
		if met, ok := result.(bool); ok {
			return met, nil
		}
	}
	return false, nil
}

// ApplyPromotion applies the actions if conditions met
func (s *marketingService) ApplyPromotion(ctx context.Context, promo *domain.Promotion, data map[string]interface{}) error {
	met, err := s.EvaluatePromotion(ctx, promo, data)
	if err != nil || !met {
		return err
	}
	// Assume Actions has "expression" key with actions to execute
	if exprStr, ok := promo.Actions["expression"].(string); ok {
		program, err := expr.Compile(exprStr, expr.Env(data))
		if err != nil {
			return err
		}
		_, err = expr.Run(program, data)
		return err
	}
	return nil
}

func (s *marketingService) GetSegments(ctx context.Context) ([]string, error) {
	// TODO: Implement customer segmentation
	return []string{"all"}, nil
}

func (s *marketingService) TriggerEmailCampaign(ctx context.Context, segment string, subject, body string) error {
	// TODO: Implement email campaign
	return nil
}
