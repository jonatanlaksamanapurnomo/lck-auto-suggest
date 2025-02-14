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
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	var cities []entity.City
	lineNum := 1 // Start after header

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++

		// Skip malformed records but log them
		if err != nil {
			log.Printf("Warning: Skipping line %d due to error: %v", lineNum, err)
			continue
		}

		// Ensure we have minimum required fields
		if len(record) < 18 {
			log.Printf("Warning: Skipping line %d due to insufficient fields (got %d fields)", lineNum, len(record))
			continue
		}

		// Safely parse fields with error handling
		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Warning: Invalid latitude on line %d: %v", lineNum, err)
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Warning: Invalid longitude on line %d: %v", lineNum, err)
			continue
		}

		// Parse population with default value if invalid
		pop, err := strconv.ParseInt(record[14], 10, 64)
		if err != nil {
			pop = 0 // Default to 0 if population is invalid
			log.Printf("Warning: Invalid population on line %d: %v, using 0", lineNum, err)
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

	log.Printf("Successfully loaded %d cities (skipped %d malformed records)", len(cities), lineNum-len(cities)-1)
	return repo.Load(cities)
}
