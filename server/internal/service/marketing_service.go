package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"

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
