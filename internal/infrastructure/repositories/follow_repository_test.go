package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFollow verifica el método Follow.
func TestFollow(t *testing.T) {
	repo := NewFollowRepository()
	ctx := context.Background()

	// Caso 1: Un usuario sigue a otro
	err := repo.Follow(ctx, "user1", "user2")
	assert.NoError(t, err, "Follow should not return an error")

	// Verificar que user1 sigue a user2
	isFollowing, err := repo.IsFollowing(ctx, "user1", "user2")
	assert.NoError(t, err, "IsFollowing should not return an error")
	assert.True(t, isFollowing, "user1 should be following user2")

	// Caso 2: Un usuario intenta seguirse a sí mismo
	err = repo.Follow(ctx, "user1", "user1")
	assert.NoError(t, err, "Follow should not return an error when following oneself")

	// Verificar que user1 se sigue a sí mismo
	isFollowing, err = repo.IsFollowing(ctx, "user1", "user1")
	assert.NoError(t, err, "IsFollowing should not return an error")
	assert.True(t, isFollowing, "user1 should be following themselves")
}

// TestIsFollowing verifica el método IsFollowing.
func TestIsFollowing(t *testing.T) {
	repo := NewFollowRepository()
	ctx := context.Background()

	// Configurar datos de prueba
	err := repo.Follow(ctx, "user1", "user2")
	assert.NoError(t, err, "Follow should not return an error")

	// Caso 1: Verificar que un usuario no sigue a otro
	isFollowing, err := repo.IsFollowing(ctx, "user2", "user1")
	assert.NoError(t, err, "IsFollowing should not return an error")
	assert.False(t, isFollowing, "user2 should not be following user1")

	// Caso 2: Verificar que un usuario sigue a otro
	isFollowing, err = repo.IsFollowing(ctx, "user1", "user2")
	assert.NoError(t, err, "IsFollowing should not return an error")
	assert.True(t, isFollowing, "user1 should be following user2")
}

// TestGetFollowers verifica el método GetFollowers.
func TestGetFollowers(t *testing.T) {
	repo := NewFollowRepository()
	ctx := context.Background()

	// Configurar datos de prueba
	err := repo.Follow(ctx, "user2", "user1")
	assert.NoError(t, err, "Follow should not return an error")
	err = repo.Follow(ctx, "user3", "user1")
	assert.NoError(t, err, "Follow should not return an error")

	// Caso 1: Obtener los seguidores de un usuario que tiene seguidores
	followers, err := repo.GetFollowers(ctx, "user1")
	assert.NoError(t, err, "GetFollowers should not return an error")
	assert.Len(t, followers, 2, "user1 should have 2 followers")
	assert.Contains(t, followers, "user2", "user2 should be a follower of user1")
	assert.Contains(t, followers, "user3", "user3 should be a follower of user1")

	// Caso 2: Obtener los seguidores de un usuario que no tiene seguidores
	followers, err = repo.GetFollowers(ctx, "user4")
	assert.NoError(t, err, "GetFollowers should not return an error")
	assert.Empty(t, followers, "user4 should have no followers")
}
