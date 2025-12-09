package service

import (
	"context"
	"server/internal/core/domain"
	"server/internal/repository"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
	addrRepo repository.Repository[domain.Address]
}

func NewUserService(userRepo repository.UserRepository, addrRepo repository.Repository[domain.Address]) UserService {
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
	return s.addrRepo.Find(ctx, "user_id = ?", userID)
}

func (s *UserServiceImpl) AddAddress(ctx context.Context, userID int, addr *domain.Address) error {
	// addr.UserID = userID // If field exists
	return s.addrRepo.Create(ctx, addr)
}

func (s *UserServiceImpl) UpdateAddress(ctx context.Context, addr *domain.Address) error {
	return s.addrRepo.Update(ctx, addr)
}

func (s *UserServiceImpl) DeleteAddress(ctx context.Context, userID, addressID int) error {
	return s.addrRepo.Delete(ctx, addressID)
}

func (s *UserServiceImpl) SetDefaultAddress(ctx context.Context, userID, addressID int, isBilling bool) error {
	// Logic to unset other defaults and set this one
	return nil
}

func (s *UserServiceImpl) GetWishlist(ctx context.Context, userID int) ([]domain.ProductVariant, error) {
	return nil, nil
}

func (s *UserServiceImpl) ToggleWishlist(ctx context.Context, userID, variantID int) (bool, error) {
	return true, nil
}

func (s *UserServiceImpl) GetUserList(ctx context.Context, filter UserFilterParams) ([]domain.User, int64, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, int64(len(users)), nil
}

func (s *UserServiceImpl) GetUserDetail(ctx context.Context, targetUserID int) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, targetUserID)
}

func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, userID int, isActive bool) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.IsActive = isActive
	return s.userRepo.Update(ctx, user)
}

func (s *UserServiceImpl) AssignRoles(ctx context.Context, userID int, role domain.UserRole) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Role = role
	return s.userRepo.Update(ctx, user)
}

func (s *UserServiceImpl) SoftDeleteUser(ctx context.Context, userID int) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *UserServiceImpl) RestoreUser(ctx context.Context, userID int) error {
	return s.userRepo.Restore(ctx, userID)
}

func (s *UserServiceImpl) ForceDeleteUser(ctx context.Context, userID int) error {
	return s.userRepo.ForceDelete(ctx, userID)
}
