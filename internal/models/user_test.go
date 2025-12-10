package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserToResponse(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "user",
		IsActive:  true,
	}

	response := user.ToResponse()

	if response.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, response.Email)
	}
	if response.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, response.FirstName)
	}
	if response.Role != user.Role {
		t.Errorf("Expected role %s, got %s", user.Role, response.Role)
	}
}

func TestListParams(t *testing.T) {
	params := ListParams{
		Page:   1,
		Limit:  10,
		Search: "test",
	}

	if params.Page != 1 {
		t.Errorf("Expected page 1, got %d", params.Page)
	}
	if params.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", params.Limit)
	}
}
