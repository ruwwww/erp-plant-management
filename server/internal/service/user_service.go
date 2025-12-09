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

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, user *domain.User) error {
	return s.userRepo.Update(ctx, user)
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
