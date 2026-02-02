package security

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	// Hash generates a hash from plain text password
	Hash(password string) (string, error)

	// Compare compares plain text password with hash
	Compare(hashedPassword, password string) error
}

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost // 10
	}

	return &BcryptHasher{
		cost: cost,
	}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *BcryptHasher) Compare(hashedPassword, password string) error {
	// bcrypt.CompareHashAndPassword is constant-time
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
