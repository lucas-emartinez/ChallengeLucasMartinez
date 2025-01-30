package repositories

import (
	"context"
	"fmt"
	"log"
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
		mu:      sync.RWMutex{},
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

	followsForFollower, ok := r.follows[followerID] // Obtener el mapa interno y verificar si existe
	if !ok {
		return false, nil // Si no existe la entrada para el seguidor, no sigue a nadie
	}

	isFollowing, ok := followsForFollower[followedID] // Obtener el valor y verificar si existe
	if !ok {
		return false, nil // Si no existe la entrada para el seguido, no lo sigue
	}

	return isFollowing, nil // Si existen ambas entradas, retornar el valor
}

// GetFollowers devuelve los seguidores de un usuario.
func (r *FollowRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	log.Println("Getting followers oon repository")
	r.mu.RLock()
	defer r.mu.RUnlock()
	log.Println("UserID: ", userID)
	if r.follows == nil {
		return nil, fmt.Errorf("follows map is uninitialized")
	}
	log.Println("Follows: ", r.follows)
	var followers []string
	for followerID, followed := range r.follows {
		if followed[userID] {
			followers = append(followers, followerID)
		}
	}
	return followers, nil
}
