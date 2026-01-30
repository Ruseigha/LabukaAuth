package mocks

import (
	"context"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
	"github.com/Ruseigha/LabukaAuth/internal/repository"
)

type MockUserRepository struct {
	// Function fields - tests can set custom behavior
	CreateFunc        func(ctx context.Context, user *entity.User) error
	FindByIDFunc      func(ctx context.Context, id valueobject.UserID) (*entity.User, error)
	FindByEmailFunc   func(ctx context.Context, email valueobject.Email) (*entity.User, error)
	UpdateFunc        func(ctx context.Context, user *entity.User) error
	DeleteFunc        func(ctx context.Context, id valueobject.UserID) error
	ExistsByEmailFunc func(ctx context.Context, email valueobject.Email) (bool, error)
	ListFunc          func(ctx context.Context, offset, limit int) ([]*entity.User, error)
	CountFunc         func(ctx context.Context) (int64, error)

	// Call tracking - verify what was called
	CreateCalls        int
	FindByIDCalls      int
	FindByEmailCalls   int
	UpdateCalls        int
	DeleteCalls        int
	ExistsByEmailCalls int
	ListCalls          int
	CountCalls         int
}

// Create implements repository.UserRepository
func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	m.CreateCalls++
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

// FindByID implements repository.UserRepository
func (m *MockUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
	m.FindByIDCalls++
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, repository.ErrUserNotFound
}

// FindByEmail implements repository.UserRepository
func (m *MockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	m.FindByEmailCalls++
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, repository.ErrUserNotFound
}

// Update implements repository.UserRepository
func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	m.UpdateCalls++
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

// Delete implements repository.UserRepository
func (m *MockUserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	m.DeleteCalls++
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// ExistsByEmail implements repository.UserRepository
func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	m.ExistsByEmailCalls++
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(ctx, email)
	}
	return false, nil
}

// List implements repository.UserRepository
func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	m.ListCalls++
	if m.ListFunc != nil {
		return m.ListFunc(ctx, offset, limit)
	}
	return []*entity.User{}, nil
}

// Count implements repository.UserRepository
func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	m.CountCalls++
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}
