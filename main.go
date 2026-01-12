package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Item represents a product in the inventory
type Item struct {
	ID          int     `json:"id"`
	SKU         string  `json:"sku"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	LastUpdated string  `json:"last_updated"`
}

var (
	inventory []Item
	mu        sync.RWMutex
)

func main() {
	// Load inventory on startup
	start := time.Now()
	if err := loadInventory("inventory.csv"); err != nil {
		log.Fatalf("Failed to load inventory: %v", err)
	}
	fmt.Printf("Loaded %d items in %v\n", len(inventory), time.Since(start))

	// Setup routes
	http.HandleFunc("/api/inventory", handleInventory)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// handleInventory handles the GET /api/inventory endpoint
func handleInventory(w http.ResponseWriter, r *http.Request) {
	// enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// loadInventory reads the CSV file and populates the inventory slice using goroutines
func loadInventory(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Read header
	if _, err := reader.Read(); err != nil {
		return err
	}

	// Channel to send lines to workers
	lines := make(chan []string, 100)
	// Channel to collect parsed items
	results := make(chan Item, 100)
	// Error channel
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	
	// Start workers (simulating heavy parsing)
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range lines {
				item, err := parseRecord(record)
				if err != nil {
					// In a real app we might handle errors differently, 
					// here we just log and continue or stop.
					// For now let's just log print it to not stop everything
					fmt.Printf("Error parsing record: %v\n", err)
					continue
				}
				results <- item
			}
		}()
	}

	// Feeder goroutine
	go func() {
		defer close(lines)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			lines <- record
		}
	}()

	// Closer goroutine
	go func() {
		wg.Wait()
		close(results)
		close(errCh)
	}()

	// Collect results
	for item := range results {
		mu.Lock()
		inventory = append(inventory, item)
		mu.Unlock()
	}

	// Check if any critical error occurred during reading
	if err, ok := <-errCh; ok {
		return err
	}

	return nil
}

func parseRecord(record []string) (Item, error) {
	if len(record) < 7 {
		return Item{}, fmt.Errorf("insufficient fields")
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return Item{}, fmt.Errorf("invalid ID: %v", err)
	}

	stock, err := strconv.Atoi(record[4])
	if err != nil {
		return Item{}, fmt.Errorf("invalid stock: %v", err)
	}

	price, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return Item{}, fmt.Errorf("invalid price: %v", err)
	}

	return Item{
		ID:          id,
		SKU:         record[1],
		ProductName: record[2],
		Category:    record[3],
		Stock:       stock,
		Price:       price,
		LastUpdated: record[6],
	}, nil
}
