package valueobject

import (
	"regexp"
	"strings"

	domainErrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {

	email = strings.TrimSpace(email)

	if email == "" {
		return Email{}, domainErrors.ErrEmptyEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return Email{}, domainErrors.ErrInvalidEmailFormat
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	return Email{value: normalizedEmail}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

func (e Email) IsEmpty() bool {
	return e.value == ""
}
