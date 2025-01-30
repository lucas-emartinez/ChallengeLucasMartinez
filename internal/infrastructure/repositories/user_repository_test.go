package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.users)
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := NewUserRepository()
	userID := "123"

	user, err := repo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)

	// Verificar que el usuario se guarda en el repositorio
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	savedUser, exists := repo.users[userID]
	assert.True(t, exists)
	assert.Equal(t, user, savedUser)
}
