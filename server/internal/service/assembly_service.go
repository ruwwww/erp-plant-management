package service

import (
	"context"
)

type AssemblyServiceImpl struct {
	// Dependencies
}

func NewAssemblyService() AssemblyService {
	return &AssemblyServiceImpl{}
}

func (s *AssemblyServiceImpl) AssembleKit(ctx context.Context, variantID, qty int, userID int) error {
	// TODO: Implement logic
	return nil
}

func (s *AssemblyServiceImpl) Disassemble(ctx context.Context, variantID, qty int, userID int) error {
	// TODO: Implement logic
	return nil
}
