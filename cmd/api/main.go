package main

import (
	"encoding/csv"
	"io"
	"lck-auto-suggest/internal/domain/entity"
	"lck-auto-suggest/internal/infrasturcure/presistence/memory"
	"lck-auto-suggest/internal/interface/http/handler"
	"lck-auto-suggest/internal/interface/repository"
	"lck-auto-suggest/internal/usecase/city"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Initialize repository
	repo := memory.NewMemoryRepository()

	// Load cities from CSV
	if err := loadCitiesFromCSV(repo); err != nil {
		log.Fatal(err)
	}

	// Initialize service and handler
	service := city.NewService(repo)
	cityHandler := handler.NewCityHandler(service)

	// Setup routes
	http.HandleFunc("/suggestions", cityHandler.GetSuggestions)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadCitiesFromCSV(repo repository.CityRepository) error {
	file, err := os.Open("data/cities_canada-usa.tsv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	var cities []entity.City
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			continue
		}

		// Handle potential parsing errors gracefully
		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Error parsing latitude for city %s: %v", record[1], err)
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Error parsing longitude for city %s: %v", record[1], err)
			continue
		}

		pop, err := strconv.ParseInt(record[14], 10, 64)
		if err != nil {
			// If population is invalid, set it to 0 but don't skip the record
			log.Printf("Error parsing population for city %s: %v", record[1], err)
			pop = 0
		}

		city := entity.City{
			ID:         record[0],
			Name:       record[1],
			ASCII:      record[2],
			AltNames:   record[3],
			Latitude:   lat,
			Longitude:  lon,
			Country:    record[8],
			Admin1:     record[10],
			Population: pop,
			Timezone:   record[17],
		}
		cities = append(cities, city)
	}

	log.Printf("Loaded %d cities successfully", len(cities))
	return repo.Load(cities)
}
