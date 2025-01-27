package repositories

import (
	"context"
	"sync"

	"ChallengeUALA/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Como en el challenge nos dan como validos cualquier usuario, vamos a devolver siempre un usuario
	// en caso de que no exista
	_, ok := r.users[id]
	if !ok {
		r.users[id] = &domain.User{
			ID: id,
		}
	}

	return r.users[id], nil
}
