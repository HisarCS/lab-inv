package main

import (
	"encoding/json"
	
	"log"
	"net/http"
	"strconv"
	"strings"

	"lab-inv/internal/model"
	"lab-inv/internal/storage"
)

var store *storage.MongoStore

func main() {
	// Initialize MongoDB store
	var err error
	store, err = storage.NewMongoStore()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer store.Close()

	log.Println("Lab Inventory System starting...")

	// Set up HTTP routes
	setupRoutes()

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// setupRoutes configures all HTTP routes
func setupRoutes() {
	// Serve static files from /static directory
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/api/items", handleItems)
	http.HandleFunc("/api/items/", handleItemByID)
	http.HandleFunc("/api/locations", handleLocations)
	http.HandleFunc("/api/locations/", handleLocationByID)
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/items-with-locations", handleItemsWithLocations)
}

// handleItems handles GET (list all) and POST (create) for items
func handleItems(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	switch r.Method {
	case http.MethodGet:
		items, err := store.GetAllItems()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(w, items)

	case http.MethodPost:
		var createItem model.CreateItem
		if err := json.NewDecoder(r.Body).Decode(&createItem); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input
		if createItem.Name == "" {
			http.Error(w, "Item name is required", http.StatusBadRequest)
			return
		}
		if createItem.LocationID == "" {
			http.Error(w, "Location ID is required", http.StatusBadRequest)
			return
		}
		if createItem.Price < 0 {
			http.Error(w, "Price must be non-negative", http.StatusBadRequest)
			return
		}

		item, err := store.AddItem(createItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendJSON(w, item)

	case http.MethodOptions:
		// Handle preflight CORS requests
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleItemByID handles GET, PUT, DELETE for individual items
func handleItemByID(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if path == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := store.GetItemByID(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		sendJSON(w, item)

	case http.MethodPut:
		var updateItem model.CreateItem
		if err := json.NewDecoder(r.Body).Decode(&updateItem); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input
		if updateItem.Name == "" {
			http.Error(w, "Item name is required", http.StatusBadRequest)
			return
		}
		if updateItem.LocationID == "" {
			http.Error(w, "Location ID is required", http.StatusBadRequest)
			return
		}
		if updateItem.Price < 0 {
			http.Error(w, "Price must be non-negative", http.StatusBadRequest)
			return
		}

		item, err := store.UpdateItem(path, updateItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendJSON(w, item)

	case http.MethodDelete:
		err := store.DeleteItem(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodOptions:
		// Handle preflight CORS requests
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLocations handles GET (list all) and POST (create) for locations
func handleLocations(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	switch r.Method {
	case http.MethodGet:
		locations, err := store.GetAllLocations()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSON(w, locations)

	case http.MethodPost:
		var createLocation model.CreateLocation
		if err := json.NewDecoder(r.Body).Decode(&createLocation); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input
		if createLocation.Name == "" {
			http.Error(w, "Location name is required", http.StatusBadRequest)
			return
		}

		location, err := store.AddLocation(createLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendJSON(w, location)

	case http.MethodOptions:
		// Handle preflight CORS requests
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLocationByID handles GET, PUT, DELETE for individual locations
func handleLocationByID(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/locations/")
	if path == "" {
		http.Error(w, "Location ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		location, err := store.GetLocationByID(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		sendJSON(w, location)

	case http.MethodDelete:
		err := store.DeleteLocation(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodOptions:
		// Handle preflight CORS requests
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSearch handles item search requests
func handleSearch(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	items, err := store.SearchItems(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, items)
}

// handleItemsWithLocations returns items with location names joined
func handleItemsWithLocations(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemsWithLocations, err := store.GetItemsWithLocations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, itemsWithLocations)
}

// enableCORS sets CORS headers for frontend requests
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// parseInt safely converts string to int
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
