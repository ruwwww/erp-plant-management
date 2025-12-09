package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type UserServiceImpl struct {
	userRepo repository.Repository[domain.User]
	addrRepo repository.Repository[domain.Address]
}

func NewUserService(userRepo repository.Repository[domain.User], addrRepo repository.Repository[domain.Address]) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		addrRepo: addrRepo,
	}
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	// 1. Call Repository
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = nil

	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *UserServiceImpl) GetAddresses(ctx context.Context, userID int) ([]domain.Address, error) {
	return nil, nil
}

func (s *UserServiceImpl) AddAddress(ctx context.Context, userID int, addr *domain.Address) error {
	// TODO: Link address to user
	return s.addrRepo.Create(ctx, addr)
}

func (s *UserServiceImpl) UpdateAddress(ctx context.Context, addr *domain.Address) error {
	return s.addrRepo.Update(ctx, addr)
}

func (s *UserServiceImpl) DeleteAddress(ctx context.Context, userID, addressID int) error {
	return s.addrRepo.Delete(ctx, addressID)
}

func (s *UserServiceImpl) SetDefaultAddress(ctx context.Context, userID, addressID int, isBilling bool) error {
	// TODO: Implement logic
	return nil
}

func (s *UserServiceImpl) GetWishlist(ctx context.Context, userID int) ([]domain.ProductVariant, error) {
	return nil, nil
}

func (s *UserServiceImpl) ToggleWishlist(ctx context.Context, userID, variantID int) (bool, error) {
	return false, nil
}

func (s *UserServiceImpl) GetUserList(ctx context.Context, filter UserFilterParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}

func (s *UserServiceImpl) GetUserDetail(ctx context.Context, targetUserID int) (*domain.User, error) {
	return nil, nil
}

func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, userID int, isActive bool) error {
	return nil
}

func (s *UserServiceImpl) AssignRoles(ctx context.Context, userID int, role domain.UserRole) error {
	return nil
}
