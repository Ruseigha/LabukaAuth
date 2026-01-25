package mongodb

import (
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

type UserDocument struct {
	ID string `bson:"_id,omitempty"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	IsActive bool `bson:"is_active"`
}
func (d *UserDocument) toEntity() (*entity.User, error) {
	// Reconstruct UserID value object
	userID, err := valueobject.NewUserIDFromString(d.ID)
	if err != nil {
		return nil, err
	}

	// Reconstruct Email value object
	email, err := valueobject.NewEmail(d.Email)
	if err != nil {
		return nil, err
	}

	// Reconstruct Password value object (from hash)
	password := valueobject.NewPasswordFromHash(d.Password)
	user := entity.ReconstructUser(
		userID,
		email,
		password,
		d.CreatedAt,
		d.UpdatedAt,
		d.IsActive,
	)

	return user, nil
}

func fromEntity(user *entity.User) *UserDocument {
	return &UserDocument{
		ID:        user.ID().String(),
		Email:     user.Email().String(),
		Password:  user.Password().Hash(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
		IsActive:  user.IsActive(),
	}
}