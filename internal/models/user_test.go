package models

import (
	"strings"
	"testing"
)

func TestValidateNewUser(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	um := &UserModel{DB: db}

	tests := []struct {
		name     string
		info     NewUserInfo
		wantErrs int
	}{
		{
			name: "Valid User",
			info: NewUserInfo{
				UserName:     "testuser",
				Email:        "test@example.com",
				Gender:       "male",
				Password:     "Password123!",
				RepeatedPass: "Password123!",
			},
			wantErrs: 0,
		},
		{
			name: "Short Username",
			info: NewUserInfo{
				UserName:     "te",
				Email:        "test@example.com",
				Gender:       "male",
				Password:     "Password123!",
				RepeatedPass: "Password123!",
			},
			wantErrs: 1,
		},
		{
			name: "Invalid Email",
			info: NewUserInfo{
				UserName:     "testuser",
				Email:        "invalid-email",
				Gender:       "male",
				Password:     "Password123!",
				RepeatedPass: "Password123!",
			},
			wantErrs: 1,
		},
		{
			name: "Password Mismatch",
			info: NewUserInfo{
				UserName:     "testuser2",
				Email:        "test2@example.com",
				Gender:       "male",
				Password:     "Password123!",
				RepeatedPass: "Password456!",
			},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, resp := um.ValidateNewUser(tt.info)
			if len(resp.Messages) != tt.wantErrs {
				t.Errorf("ValidateNewUser() got %v errors, want %v. Errs: %v", len(resp.Messages), tt.wantErrs, resp.Messages)
			}
		})
	}
}

func TestUserModel_InsertAndValidate(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	um := &UserModel{DB: db}

	info := NewUserInfo{
		UserName:     "john_doe",
		Email:        "john@example.com",
		Gender:       "male",
		Password:     "StrongPass123!",
		RepeatedPass: "StrongPass123!",
	}

	newUser, resp := um.ValidateNewUser(info)
	if len(resp.Messages) != 0 {
		t.Fatalf("Validation failed: %v", resp.Messages)
	}

	err := um.InsertUser(newUser)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Test credential validation
	id, errs := um.ValidateUserCredentials("john_doe", "StrongPass123!")
	if id <= 0 || len(errs) != 0 {
		t.Errorf("Credential validation failed: %v", errs)
	}

	// Test duplicate username
	_, resp = um.ValidateNewUser(info)
	foundDuplicate := false
	for _, m := range resp.Messages {
		if strings.Contains(strings.ToLower(m), "already exists") {
			foundDuplicate = true
			break
		}
	}
	if !foundDuplicate {
		t.Errorf("Expected duplicate username error, got: %v", resp.Messages)
	}
}
