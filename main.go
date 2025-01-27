package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type MedicineEntry struct {
	ProductCodeNAPPI string
	ProductName      string
	ProductStrength  string
	ProductForm      string
}

var (
	medicineData []MedicineEntry
	dataLock     sync.RWMutex
	
	// Index for efficient product name lookups
	productNameIndex map[string][]int
)

func init() {
	// Path to the NAPPI file
	filePath := "nappi_data.txt"

	// Load the NAPPI data on startup
	log.Println("Loading the NAPPI data...")
	start := time.Now()
	if err := loadNAPPIFile(filePath); err != nil {
		log.Fatalf("Failed to load NAPPI file: %v", err)
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed time:", elapsed)
	
	
	log.Println("Building index for ProductName auto complete data...")
	start = time.Now()
	// Initialize the index (called in init() in your actual code)
	initIndex()
	elapsed = time.Since(start)
	fmt.Println("Elapsed time:", elapsed)
}

// parseFixedWidthLine parses a single line from the fixed-width-delimited file.
func parseFixedWidthLine(line string) MedicineEntry {
	productCodeNAPPI := strings.TrimSpace(line[11:20])
	productName := strings.TrimSpace(line[20:58])
	productStrength := strings.TrimSpace(line[59:75])
	productForm := strings.TrimSpace(line[75:79])
	
	return MedicineEntry{
		ProductCodeNAPPI: productCodeNAPPI,
		ProductName: productName,
		ProductStrength: productStrength,
		ProductForm: productForm,
	}
}

// loadNAPPIFile loads the content of the position-delimited NAPPI file into memory.
func loadNAPPIFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tempData []MedicineEntry

	for scanner.Scan() {
		line := scanner.Text()
		entry := parseFixedWidthLine(line)
		tempData = append(tempData, entry)
		
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Update the global data under a write lock
	dataLock.Lock()
	medicineData = tempData
	dataLock.Unlock()

	fmt.Println("entry count:", len(tempData))

	return nil
}

func initIndex() {
        // Initialize the product name index
        productNameIndex = make(map[string][]int)
        for i, entry := range medicineData {
                productName := strings.ToLower(entry.ProductName) // Case-insensitive search
                productNameIndex[productName] = append(productNameIndex[productName], i) 
        }
		
	fmt.Println("productNameIndex size:", len(productNameIndex))
}

// FindMedicineEntriesByKeywords performs an efficient keyword search 
func FindMedicineEntriesByKeywords(keywords string) []MedicineEntry {

        var results []MedicineEntry
        keywordList := strings.Fields(keywords)

		if len(keywordList) == 0 {
			return results // No keywords provided
		}
		dataLock.RLock()
		defer dataLock.RUnlock()

        // Create a map to track which entries match all keywords
        entryMatches := make(map[int]map[string]bool) 

        /*
        for _, keyword := range keywordList {

			if len(keyword) < 3 {
				continue // Skip short keywords
			}

			keywordLower := strings.ToLower(keyword)
			for productName, indices := range productNameIndex {
				if strings.Contains(productName, keywordLower) {
					for _, i := range indices {
						if _, ok := entryMatches[i]; !ok {
							entryMatches[i] = make(map[string]bool)
						}
						entryMatches[i][keywordLower] = true
					}
				}
			}
		}
        */

		for productName, indices := range productNameIndex {
			
			for _, keyword := range keywordList {

				if len(keyword) < 3 {
					continue // Skip short keywords
				}

				keywordLower := strings.ToLower(keyword)

				if strings.Contains(productName, keywordLower) {
					for _, i := range indices {
						if _, ok := entryMatches[i]; !ok {
							entryMatches[i] = make(map[string]bool)
						}
						entryMatches[i][keywordLower] = true
					}
				}
			}
		}

        // Collect matching entries
        for entryIndex, matchedKeywords := range entryMatches {
			if len(matchedKeywords) == len(keywordList) { // All keywords matched
				results = append(results, medicineData[entryIndex])
			}
		}

        return results
}

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("term")
	if query == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}
	
	if len(query) < 3 {
		http.Error(w, "Search term must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	start := time.Now()

	results := FindMedicineEntriesByKeywords(query)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"results": results}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	
	elapsed := time.Since(start)
	fmt.Println("Elapsed time:", elapsed)
}

func main() {

	httpServer := &http.Server{Addr: ":8080"}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	http.HandleFunc("/autocomplete", autocompleteHandler)

	log.Println("Server is running on port 8080...")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server exited cleanly")
}
