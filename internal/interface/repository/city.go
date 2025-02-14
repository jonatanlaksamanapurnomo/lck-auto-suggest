package repository

import "lck-auto-suggest/internal/domain/entity"

type CityRepository interface {
	Search(query string) ([]entity.City, error)
	GetAll() ([]entity.City, error)
	Load(cities []entity.City) error
}
