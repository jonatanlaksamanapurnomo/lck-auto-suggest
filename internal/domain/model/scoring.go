package model

import (
	"lck-auto-suggest/internal/domain/entity"
	"math"
	"strings"
)

type ScoringModel struct{}

func NewScoringModel() *ScoringModel {
	return &ScoringModel{}
}

func (s *ScoringModel) CalculateScore(city entity.City, query string, lat, lon *float64) float64 {
	// Base score from name matching
	nameScore := calculateNameScore(city, query)

	// Location score if coordinates are provided
	locationScore := 1.0
	if lat != nil && lon != nil {
		locationScore = calculateLocationScore(city, *lat, *lon)
	}

	// Population factor (bigger cities get slight boost)
	popScore := calculatePopulationScore(city.Population)

	// Combine scores with weights
	finalScore := (nameScore * 0.6) + (locationScore * 0.3) + (popScore * 0.1)

	return math.Min(math.Max(finalScore, 0.0), 1.0)
}

func calculateNameScore(city entity.City, query string) float64 {
	query = strings.ToLower(strings.TrimSpace(query))
	name := strings.ToLower(city.Name)
	ascii := strings.ToLower(city.ASCII)

	// Exact match gets highest score
	if name == query || ascii == query {
		return 1.0
	}

	// Prefix match gets high score
	if strings.HasPrefix(name, query) || strings.HasPrefix(ascii, query) {
		return 0.9
	}

	// Contains match gets medium score
	if strings.Contains(name, query) || strings.Contains(ascii, query) {
		return 0.7
	}

	// Check alternative names
	if strings.Contains(strings.ToLower(city.AltNames), query) {
		return 0.5
	}

	return 0.0
}

func calculateLocationScore(city entity.City, targetLat, targetLon float64) float64 {
	// Haversine distance calculation
	const earthRadius = 6371.0 // km

	lat1 := targetLat * math.Pi / 180
	lon1 := targetLon * math.Pi / 180
	lat2 := city.Latitude * math.Pi / 180
	lon2 := city.Longitude * math.Pi / 180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	// Convert distance to a score (closer = higher score)
	// Max distance considered is 1000km
	maxDistance := 1000.0
	if distance > maxDistance {
		return 0.0
	}
	return 1.0 - (distance / maxDistance)
}

func calculatePopulationScore(population int64) float64 {
	// Simple logarithmic scaling for population
	if population == 0 {
		return 0.0
	}
	return math.Min(math.Log10(float64(population))/6.0, 1.0)
}
