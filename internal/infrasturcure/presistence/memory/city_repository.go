package memory

import (
	"lck-auto-suggest/internal/domain/entity"
	"strings"
	"sync"
)

type memoryRepository struct {
	cities []entity.City
	mu     sync.RWMutex
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		cities: make([]entity.City, 0),
	}
}

func (r *memoryRepository) Load(cities []entity.City) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cities = cities
	return nil
}

func (r *memoryRepository) Search(query string) ([]entity.City, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil, nil
	}

	var results []entity.City
	for _, city := range r.cities {
		if strings.Contains(strings.ToLower(city.Name), query) ||
			strings.Contains(strings.ToLower(city.ASCII), query) ||
			strings.Contains(strings.ToLower(city.AltNames), query) {
			results = append(results, city)
		}
	}

	return results, nil
}

func (r *memoryRepository) GetAll() ([]entity.City, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.cities, nil
}
