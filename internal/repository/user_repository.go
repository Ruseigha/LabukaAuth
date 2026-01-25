package repository

import (
	"context"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id valueobject.UserID) error
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)
}
