package entity_test

import (
	"testing"
	"time"

	"github.com/Ruseigha/LabukaAuth/internal/domain/entity"
	"github.com/Ruseigha/LabukaAuth/internal/domain/valueobject"
)

// Helper function to create valid test data
// WHY: DRY - reuse across tests
func createValidTestUser(t *testing.T) *entity.User {
	t.Helper() // Marks this as a helper - errors report caller's line

	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create test email: %v", err)
	}

	password, err := valueobject.NewPassword("SecureP@ss123")
	if err != nil {
		t.Fatalf("Failed to create test password: %v", err)
	}

	user, err := entity.NewUser(email, password)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}


func TestNewUser_Success(t *testing.T) {
    email, _ := valueobject.NewEmail("test@example.com")
    password, _ := valueobject.NewPassword("SecureP@ss123")
    
    user, err := entity.NewUser(email, password)
    
    if err != nil {
        t.Errorf("NewUser() unexpected error = %v", err)
        return
    }
    
    // Verify all fields are set correctly
    if user == nil {
        t.Fatal("NewUser() returned nil")
    }
    
    if user.ID().IsEmpty() {
        t.Error("User ID should not be empty")
    }
    
    if !user.Email().Equals(email) {
        t.Errorf("User email = %v, want %v", user.Email(), email)
    }
    
    if user.Password().Hash() != password.Hash() {
        t.Error("User password not set correctly")
    }
    
    if user.CreatedAt().IsZero() {
        t.Error("User CreatedAt should be set")
    }
    
    if user.UpdatedAt().IsZero() {
        t.Error("User UpdatedAt should be set")
    }
    
    if !user.IsActive() {
        t.Error("New user should be active by default")
    }
}


func TestNewUser_InvalidInputs(t *testing.T) {
    validEmail, _ := valueobject.NewEmail("test@example.com")
    validPassword, _ := valueobject.NewPassword("SecureP@ss123")
    emptyEmail := valueobject.Email{}
    emptyPassword := valueobject.Password{}
    
    tests := []struct {
        name      string
        email     valueobject.Email
        password  valueobject.Password
        wantError bool
    }{
        {
            name:      "empty email",
            email:     emptyEmail,
            password:  validPassword,
            wantError: true,
        },
        {
            name:      "empty password",
            email:     validEmail,
            password:  emptyPassword,
            wantError: true,
        },
        {
            name:      "both empty",
            email:     emptyEmail,
            password:  emptyPassword,
            wantError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := entity.NewUser(tt.email, tt.password)
            
            if tt.wantError && err == nil {
                t.Error("NewUser() expected error, got nil")
            }
            
            if !tt.wantError && err != nil {
                t.Errorf("NewUser() unexpected error = %v", err)
            }
            
            if tt.wantError && user != nil {
                t.Error("NewUser() should return nil user on error")
            }
        })
    }
}



func TestUser_Deactivate(t *testing.T) {
    user := createValidTestUser(t)
    
    // Initially active
    if !user.IsActive() {
        t.Error("New user should be active")
    }
    
    // Store original updatedAt
    originalUpdatedAt := user.UpdatedAt()
    
    // Small delay to ensure time difference
    time.Sleep(10 * time.Millisecond)
    
    // Deactivate
    user.Deactivate()
    
    // Should be inactive
    if user.IsActive() {
        t.Error("User should be inactive after Deactivate()")
    }
    
    // UpdatedAt should change
    if !user.UpdatedAt().After(originalUpdatedAt) {
        t.Error("UpdatedAt should be updated after Deactivate()")
    }
}



func TestUser_Activate(t *testing.T) {
    user := createValidTestUser(t)
    
    // Deactivate first
    user.Deactivate()
    if user.IsActive() {
        t.Fatal("User should be inactive")
    }
    
    originalUpdatedAt := user.UpdatedAt()
    time.Sleep(10 * time.Millisecond)
    
    // Activate
    user.Activate()
    
    // Should be active
    if !user.IsActive() {
        t.Error("User should be active after Activate()")
    }
    
    // UpdatedAt should change
    if !user.UpdatedAt().After(originalUpdatedAt) {
        t.Error("UpdatedAt should be updated after Activate()")
    }
}



func TestUser_UpdateEmail(t *testing.T) {
    user := createValidTestUser(t)
    
    newEmail, _ := valueobject.NewEmail("newemail@example.com")
    originalUpdatedAt := user.UpdatedAt()
    time.Sleep(10 * time.Millisecond)
    
    err := user.UpdateEmail(newEmail)
    
    if err != nil {
        t.Errorf("UpdateEmail() unexpected error = %v", err)
    }
    
    if !user.Email().Equals(newEmail) {
        t.Errorf("Email not updated: got %v, want %v", user.Email(), newEmail)
    }
    
    if !user.UpdatedAt().After(originalUpdatedAt) {
        t.Error("UpdatedAt should be updated after UpdateEmail()")
    }
}

func TestUser_UpdateEmail_Invalid(t *testing.T) {
    user := createValidTestUser(t)
    emptyEmail := valueobject.Email{}
    
    err := user.UpdateEmail(emptyEmail)
    
    if err == nil {
        t.Error("UpdateEmail() expected error for empty email")
    }
}

func TestUser_UpdatePassword(t *testing.T) {
    user := createValidTestUser(t)
    
    newPassword, _ := valueobject.NewPassword("NewSecure@123")
    originalUpdatedAt := user.UpdatedAt()
    time.Sleep(10 * time.Millisecond)
    
    err := user.UpdatePassword(newPassword)
    
    if err != nil {
        t.Errorf("UpdatePassword() unexpected error = %v", err)
    }
    
    if user.Password().Hash() != newPassword.Hash() {
        t.Error("Password not updated")
    }
    
    if !user.UpdatedAt().After(originalUpdatedAt) {
        t.Error("UpdatedAt should be updated after UpdatePassword()")
    }
}

func TestUser_CanLogin(t *testing.T) {
    user := createValidTestUser(t)
    
    // Active user can login
    if !user.CanLogin() {
        t.Error("Active user should be able to login")
    }
    
    // Inactive user cannot login
    user.Deactivate()
    if user.CanLogin() {
        t.Error("Inactive user should not be able to login")
    }
}

func TestUser_Validate(t *testing.T) {
    user := createValidTestUser(t)
    
    if err := user.Validate(); err != nil {
        t.Errorf("Validate() unexpected error for valid user: %v", err)
    }
}

func TestReconstructUser(t *testing.T) {
    id := valueobject.NewUserID()
    email, _ := valueobject.NewEmail("test@example.com")
    password, _ := valueobject.NewPassword("SecureP@ss123")
    createdAt := time.Now().UTC()
    updatedAt := time.Now().UTC()
    isActive := true
    
    user := entity.ReconstructUser(id, email, password, createdAt, updatedAt, isActive)
    
    if user == nil {
        t.Fatal("ReconstructUser() returned nil")
    }
    
    if !user.ID().Equals(id) {
        t.Error("Reconstructed user ID mismatch")
    }
    
    if !user.Email().Equals(email) {
        t.Error("Reconstructed user email mismatch")
    }
    
    if user.IsActive() != isActive {
        t.Error("Reconstructed user isActive mismatch")
    }
}