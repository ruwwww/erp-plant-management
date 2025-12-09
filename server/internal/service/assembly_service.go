package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"
)

type AssemblyServiceImpl struct {
	recipeRepo       repository.RecipeRepository
	assemblyRepo     repository.AssemblyRepository
	inventoryService InventoryService
}

func NewAssemblyService(recipeRepo repository.RecipeRepository, assemblyRepo repository.AssemblyRepository, inventoryService InventoryService) AssemblyService {
	return &AssemblyServiceImpl{
		recipeRepo:       recipeRepo,
		assemblyRepo:     assemblyRepo,
		inventoryService: inventoryService,
	}
}

func (s *AssemblyServiceImpl) GetRecipes(ctx context.Context) ([]domain.ProductRecipe, error) {
	return s.recipeRepo.FindAll(ctx)
}

func (s *AssemblyServiceImpl) CreateRecipe(ctx context.Context, recipe *domain.ProductRecipe) error {
	return s.recipeRepo.Create(ctx, recipe)
}

func (s *AssemblyServiceImpl) DeleteRecipe(ctx context.Context, recipeID int) error {
	return s.recipeRepo.Delete(ctx, recipeID)
}

func (s *AssemblyServiceImpl) AssembleKit(ctx context.Context, variantID, qty int, userID int) error {
	// Get recipe for this variant
	recipes, err := s.recipeRepo.GetRecipe(ctx, variantID)
	if err != nil {
		return err
	}
	if len(recipes) == 0 {
		return errors.New("no recipe found for variant")
	}

	// Check stock for components
	for _, r := range recipes {
		currentQty, err := s.inventoryService.GetStockLevel(ctx, r.ChildVariantID, 1) // Assume location 1 for now
		if err != nil {
			return err
		}
		required := r.QuantityNeeded * float64(qty)
		if float64(currentQty) < required {
			return errors.New("insufficient stock for component")
		}
	}

	// Consume components
	for _, r := range recipes {
		change := -int(r.QuantityNeeded * float64(qty))
		cmd := StockMoveCmd{
			LocationID: 1, // TODO: Get from config or param
			VariantID:  r.ChildVariantID,
			QtyChange:  change,
			Reason:     domain.ReasonAssemblyConsumption,
			UserID:     userID,
		}
		if err := s.inventoryService.ExecuteMovement(ctx, cmd); err != nil {
			return err
		}
	}

	// Produce output
	cmd := StockMoveCmd{
		LocationID: 1,
		VariantID:  variantID,
		QtyChange:  qty,
		Reason:     domain.ReasonAssemblyOutput,
		UserID:     userID,
	}
	if err := s.inventoryService.ExecuteMovement(ctx, cmd); err != nil {
		return err
	}

	// Log assembly
	assembly := &domain.StockAssembly{
		AssemblyNumber:   "ASM-" + time.Now().Format("20060102-150405"), // TODO: Generate properly
		VariantID:        variantID,
		QuantityProduced: qty,
		TotalCost:        0, // TODO: Calculate
		CreatedBy:        &userID,
		CreatedAt:        time.Now(),
	}
	return s.assemblyRepo.Create(ctx, assembly)
}

func (s *AssemblyServiceImpl) Disassemble(ctx context.Context, variantID, qty int, userID int) error {
	// Reverse of assemble: consume kit, produce components
	recipes, err := s.recipeRepo.GetRecipe(ctx, variantID)
	if err != nil {
		return err
	}
	if len(recipes) == 0 {
		return errors.New("no recipe found for variant")
	}

	// Check stock for kit
	currentQty, err := s.inventoryService.GetStockLevel(ctx, variantID, 1)
	if err != nil {
		return err
	}
	if currentQty < qty {
		return errors.New("insufficient stock for kit")
	}

	// Consume kit
	cmd := StockMoveCmd{
		LocationID: 1,
		VariantID:  variantID,
		QtyChange:  -qty,
		Reason:     domain.ReasonAssemblyConsumption, // Or new reason
		UserID:     userID,
	}
	if err := s.inventoryService.ExecuteMovement(ctx, cmd); err != nil {
		return err
	}

	// Produce components
	for _, r := range recipes {
		change := int(r.QuantityNeeded * float64(qty))
		cmd := StockMoveCmd{
			LocationID: 1,
			VariantID:  r.ChildVariantID,
			QtyChange:  change,
			Reason:     domain.ReasonAssemblyOutput,
			UserID:     userID,
		}
		if err := s.inventoryService.ExecuteMovement(ctx, cmd); err != nil {
			return err
		}
	}

	// Log disassembly (similar to assembly)
	assembly := &domain.StockAssembly{
		AssemblyNumber:   "DIS-" + time.Now().Format("20060102-150405"),
		VariantID:        variantID,
		QuantityProduced: -qty, // Negative for disassembly
		TotalCost:        0,
		CreatedBy:        &userID,
		CreatedAt:        time.Now(),
	}
	return s.assemblyRepo.Create(ctx, assembly)
}

func (s *AssemblyServiceImpl) GetAssemblyLogs(ctx context.Context, page, limit int) ([]domain.StockAssembly, error) {
	return s.assemblyRepo.GetLogs(ctx, page, limit)
}
