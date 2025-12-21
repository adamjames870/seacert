package tests

import (
	"context"
	"testing"

	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestUserContext(t *testing.T) {
	id := uuid.New()
	user := dto.User{
		Id:    id.String(),
		Email: "test@example.com",
	}

	ctx := context.Background()

	// Test UserFromContext when not set
	_, ok := auth.UserFromContext(ctx)
	if ok {
		t.Error("UserFromContext returned ok=true for empty context")
	}

	// Test UserIdFromContext when not set
	_, err := auth.UserIdFromContext(ctx)
	if err == nil {
		t.Error("UserIdFromContext returned nil error for empty context")
	}

	// Test WithUser and UserFromContext
	ctx = auth.WithUser(ctx, user)
	retrievedUser, ok := auth.UserFromContext(ctx)
	if !ok {
		t.Fatal("UserFromContext returned ok=false after WithUser")
	}
	if retrievedUser.Id != user.Id || retrievedUser.Email != user.Email {
		t.Errorf("Retrieved user mismatch: got %+v, want %+v", retrievedUser, user)
	}

	// Test UserIdFromContext
	retrievedId, err := auth.UserIdFromContext(ctx)
	if err != nil {
		t.Fatalf("UserIdFromContext failed: %v", err)
	}
	if retrievedId != id {
		t.Errorf("Retrieved ID mismatch: got %v, want %v", retrievedId, id)
	}
}
