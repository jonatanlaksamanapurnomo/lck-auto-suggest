package city

import (
	"lck-auto-suggest/internal/domain/entity"
	"lck-auto-suggest/internal/domain/model"
	"lck-auto-suggest/internal/interface/repository"
	"sort"
)

type Service struct {
	repo    repository.CityRepository
	scoring *model.ScoringModel
}

func NewService(repo repository.CityRepository) *Service {
	return &Service{
		repo:    repo,
		scoring: model.NewScoringModel(),
	}
}

type SuggestionResult struct {
	Suggestions []entity.ScoredCity `json:"suggestions"`
}

func (s *Service) GetSuggestions(query string, lat, lon *float64) (*SuggestionResult, error) {
	// Search for matching cities
	cities, err := s.repo.Search(query)
	if err != nil {
		return nil, err
	}

	// Score each city
	scoredCities := make([]entity.ScoredCity, 0, len(cities))
	for _, city := range cities {
		score := s.scoring.CalculateScore(city, query, lat, lon)
		if score > 0 {
			scoredCity := entity.ScoredCity{
				Name:      formatCityName(city),
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
				Score:     score,
			}
			scoredCities = append(scoredCities, scoredCity)
		}
	}

	// Sort by score descending
	sort.Slice(scoredCities, func(i, j int) bool {
		return scoredCities[i].Score > scoredCities[j].Score
	})

	// Limit to top results
	const maxResults = 5
	if len(scoredCities) > maxResults {
		scoredCities = scoredCities[:maxResults]
	}

	return &SuggestionResult{
		Suggestions: scoredCities,
	}, nil
}

func formatCityName(city entity.City) string {
	// Format: "City, State, Country"
	if city.Admin1 != "" {
		return city.Name + ", " + city.Admin1 + ", " + city.Country
	}
	return city.Name + ", " + city.Country
}
