package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号、邮箱或密码不正确")
)

type IUserService interface {
	Login(ctx context.Context, email string, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	EditProfile(ctx context.Context, u domain.Profile) error
	QueryProfile(ctx context.Context, userId uint64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type UserService struct {
	repo                   repository.UserRepository
	compareHashAndPassword func(hashedPassword []byte, password []byte) error
}

func NewUserService(repo repository.UserRepository) IUserService {
	return &UserService{
		repo:                   repo,
		compareHashAndPassword: bcrypt.CompareHashAndPassword,
	}
}

func NewUserServiceForTest(repo repository.UserRepository, compareHashAndPassword func(hashedPassword []byte, password []byte) error) IUserService {
	return &UserService{
		repo:                   repo,
		compareHashAndPassword: compareHashAndPassword,
	}
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = svc.compareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) EditProfile(ctx context.Context, u domain.Profile) error {
	return svc.repo.UpdateProfile(ctx, u)
}

func (svc *UserService) QueryProfile(ctx context.Context, userId uint64) (domain.User, error) {
	return svc.repo.QueryProfile(ctx, userId)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}

	u = domain.User{Phone: phone}
	err = svc.repo.Create(ctx, u)
	if err != nil && !errors.Is(err, repository.ErrUserDuplicate) {
		return u, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return svc.compareHashAndPassword(hashedPassword, password)
}
