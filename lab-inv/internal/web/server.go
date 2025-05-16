// lab-inv/internal/web/server.go

package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"lab-inv/internal/storage"
)

// Server represents the web server for the inventory application
type Server struct {
	store  *storage.FileStore
	router *http.ServeMux
	server *http.Server
}

// NewServer creates a new web server instance
func NewServer(store *storage.FileStore, staticDir string) *Server {
	s := &Server{
		store:  store,
		router: http.NewServeMux(),
	}

	// Set up routes
	s.routes(staticDir)

	return s
}

// routes sets up all the HTTP routes for the server
func (s *Server) routes(staticDir string) {
	// API endpoints
	s.router.HandleFunc("/api/items", s.handleItems)
	s.router.HandleFunc("/api/items/", s.handleItemByID)
	s.router.HandleFunc("/api/items/search", s.handleSearchItems)
	s.router.HandleFunc("/api/items/with-locations", s.handleItemsWithLocations)
	s.router.HandleFunc("/api/locations", s.handleLocations)
	s.router.HandleFunc("/api/locations/", s.handleLocationByID)

	// Serve static files
	fs := http.FileServer(http.Dir(staticDir))
	s.router.Handle("/", fs)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting server at %s", addr)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	return s.server.Close()
}

// handleItems handles GET and POST requests for /api/items
func (s *Server) handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetItems(w, r)
	case http.MethodPost:
		s.handleCreateItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleItemByID handles GET, PUT and DELETE requests for /api/items/{id}
func (s *Server) handleItemByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id, err := parseIDFromURL(r.URL.Path, "/api/items/")
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGetItem(w, r, id)
	case http.MethodPut:
		s.handleUpdateItem(w, r, id)
	case http.MethodDelete:
		s.handleDeleteItem(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLocations handles GET and POST requests for /api/locations
func (s *Server) handleLocations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetLocations(w, r)
	case http.MethodPost:
		s.handleCreateLocation(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLocationByID handles GET, PUT and DELETE requests for /api/locations/{id}
func (s *Server) handleLocationByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id, err := parseIDFromURL(r.URL.Path, "/api/locations/")
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGetLocation(w, r, id)
	case http.MethodPut:
		s.handleUpdateLocation(w, r, id)
	case http.MethodDelete:
		s.handleDeleteLocation(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSearchItems handles GET requests for /api/items/search
func (s *Server) handleSearchItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.handleSearch(w, r)
}

// handleItemsWithLocations handles GET requests for /api/items/with-locations
func (s *Server) handleItemsWithLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.handleGetItemsWithLocations(w, r)
}

// Other handler method signatures - these will be implemented in handlers.go
func (s *Server) handleGetItems(w http.ResponseWriter, r *http.Request)                {}
func (s *Server) handleCreateItem(w http.ResponseWriter, r *http.Request)              {}
func (s *Server) handleGetItem(w http.ResponseWriter, r *http.Request, id uint)        {}
func (s *Server) handleUpdateItem(w http.ResponseWriter, r *http.Request, id uint)     {}
func (s *Server) handleDeleteItem(w http.ResponseWriter, r *http.Request, id uint)     {}
func (s *Server) handleGetLocations(w http.ResponseWriter, r *http.Request)            {}
func (s *Server) handleCreateLocation(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) handleGetLocation(w http.ResponseWriter, r *http.Request, id uint)    {}
func (s *Server) handleUpdateLocation(w http.ResponseWriter, r *http.Request, id uint) {}
func (s *Server) handleDeleteLocation(w http.ResponseWriter, r *http.Request, id uint) {}
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request)                  {}
func (s *Server) handleGetItemsWithLocations(w http.ResponseWriter, r *http.Request)   {}

// Helper function to parse IDs from URL paths
func parseIDFromURL(path string, prefix string) (uint, error) {
	// Implementation will be in handlers.go
	return 0, fmt.Errorf("not implemented")
}
