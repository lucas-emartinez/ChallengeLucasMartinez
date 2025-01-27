package repositories

import (
	"context"
	"sync"
)

// FollowRepository es una implementación en memoria de la interfaz FollowRepository
type FollowRepository struct {
	mu      sync.RWMutex
	follows map[string]map[string]bool
}

// FollowRepository crea una nueva instancia de FollowRepository
func NewFollowRepository() *FollowRepository {
	return &FollowRepository{
		follows: make(map[string]map[string]bool),
	}
}

// Follow establece una relación de seguimiento entre dos usuarios.
func (r *FollowRepository) Follow(ctx context.Context, followerID, followedID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.follows[followerID] == nil {
		r.follows[followerID] = make(map[string]bool)
	}
	r.follows[followerID][followedID] = true
	return nil
}

// IsFollowing verifica si un usuario sigue a otro.
func (r *FollowRepository) IsFollowing(ctx context.Context, followerID, followedID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.follows[followerID] == nil {
		return false, nil
	}
	return r.follows[followerID][followedID], nil
}

// GetFollowers devuelve los seguidores de un usuario.
func (r *FollowRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var followers []string
	for followerID, followed := range r.follows {
		if followed[userID] {
			followers = append(followers, followerID)
		}
	}
	return followers, nil
}
