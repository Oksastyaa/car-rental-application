package seeder

import (
	"car-rental-application/internal/models"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Make struct {
	Name string `json:"name"`
}

type CarAPIResponse struct {
	Make  Make   `json:"make"`
	Model string `json:"name"`
}

type APIResponse struct {
	Data []CarAPIResponse `json:"data"`
}

// FetchCarsFromApi fetches cars data from an API
func FetchCarsFromApi() ([]CarAPIResponse, error) {
	apiUrl := os.Getenv("RAPID_API_URL")
	apiKey := os.Getenv("RAPID_API_KEY")
	apiHost := os.Getenv("RAPID_API_HOST")

	if apiUrl == "" || apiKey == "" || apiHost == "" {
		return nil, fmt.Errorf("missing API configuration: check RAPID_API_URL, RAPID_API_KEY, RAPID_API_HOST")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logrus.Infof("API Response: %s", string(body))

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return apiResponse.Data, nil
}

// SeedCars seeds the cars table with sample data
func SeedCars(db *gorm.DB) {
	carData, err := FetchCarsFromApi()
	if err != nil {
		logrus.Fatalf("Error fetching cars from API: %v", err)
	}

	if len(carData) == 0 {
		logrus.Fatalf("No car data found from API.")
	}

	var cars []models.Car
	for _, car := range carData {
		carModel := models.Car{
			Name:              car.Model,
			Brands:            car.Make.Name,
			RentalCost:        float64(rand.Intn(50000) + 100000),
			StockAvailability: rand.Intn(5) + 1,
		}
		cars = append(cars, carModel)
	}

	if len(cars) == 0 {
		logrus.Fatalf("Error seeding cars: empty slice found")
	}

	if err := db.Create(&cars).Error; err != nil {
		logrus.Fatalf("Error seeding cars: %v", err)
	}

	logrus.Info("Cars seeded successfully")
}
