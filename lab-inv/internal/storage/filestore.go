// lab-inv/internal/storage/filestore.go

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"lab-inv/internal/model"
)

const (
	// Default data directory
	defaultDataDir = "./data"
	
	// File names
	inventoryFileName = "inventory.json"
)

// FileStore handles storing inventory data in JSON files
type FileStore struct {
	dataDir        string
	inventoryFile  string
	inventory      model.Inventory
	nextItemID     uint
	nextLocationID uint
	mu             sync.RWMutex
}

// NewFileStore creates a new FileStore instance
func NewFileStore(dataDir string) (*FileStore, error) {
	if dataDir == "" {
		dataDir = defaultDataDir
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	fs := &FileStore{
		dataDir:        dataDir,
		inventoryFile:  filepath.Join(dataDir, inventoryFileName),
		inventory:      model.Inventory{
			Items:     []model.Item{},
			Locations: []model.Location{},
		},
		nextItemID:     1,
		nextLocationID: 1,
	}

	// Load existing data if any
	if err := fs.loadInventory(); err != nil {
		// If the file doesn't exist, that's fine - we'll create it
		if !os.IsNotExist(err) {
			return nil, err
		}
		
		// Initialize with sample data if starting fresh
		fs.initSampleData()
		if err := fs.saveInventory(); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

// loadInventory reads the inventory from the JSON file
func (fs *FileStore) loadInventory() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.inventoryFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &fs.inventory); err != nil {
		return fmt.Errorf("failed to unmarshal inventory data: %w", err)
	}

	// Determine the next IDs
	fs.updateNextIDs()

	return nil
}

// saveInventory writes the inventory to the JSON file
func (fs *FileStore) saveInventory() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := json.MarshalIndent(fs.inventory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal inventory data: %w", err)
	}

	if err := os.WriteFile(fs.inventoryFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write inventory file: %w", err)
	}

	return nil
}

// updateNextIDs determines the next available IDs based on existing items and locations
func (fs *FileStore) updateNextIDs() {
	fs.nextItemID = 1
	fs.nextLocationID = 1

	for _, item := range fs.inventory.Items {
		if item.ID >= fs.nextItemID {
			fs.nextItemID = item.ID + 1
		}
	}

	for _, location := range fs.inventory.Locations {
		if location.ID >= fs.nextLocationID {
			fs.nextLocationID = location.ID + 1
		}
	}
}

// initSampleData creates some initial data for a new inventory
func (fs *FileStore) initSampleData() {
	// Create sample locations
	storageRoom := model.Location{
		ID:       fs.nextLocationID,
		Name:     "Storage Room",
		Modified: time.Now(),
	}
	fs.nextLocationID++

	assemblyRoom := model.Location{
		ID:       fs.nextLocationID,
		Name:     "Assembly Room",
		Modified: time.Now(),
	}
	fs.nextLocationID++

	electronics := model.Location{
		ID:       fs.nextLocationID,
		Name:     "Electronics",
		Modified: time.Now(),
	}
	fs.nextLocationID++

	fs.inventory.Locations = []model.Location{storageRoom, assemblyRoom, electronics}

	// Create sample items
	fs.inventory.Items = []model.Item{
		{
			ID:         fs.nextItemID,
			Name:       "Plywood 2mm 900x600mm Sheet",
			LocationID: storageRoom.ID,
			Price:      11.15,
			Modified:   time.Now(),
		},
		{
			ID:         fs.nextItemID + 1,
			Name:       "MDF 4mm 900x600mm Sheet",
			LocationID: storageRoom.ID,
			Price:      5.67,
			Modified:   time.Now(),
		},
		{
			ID:         fs.nextItemID + 2,
			Name:       "Acrylic 5mm 900x600mm Sheet",
			LocationID: storageRoom.ID,
			Price:      19.34,
			Modified:   time.Now(),
		},
		{
			ID:         fs.nextItemID + 3,
			Name:       "Wood Glue",
			LocationID: assemblyRoom.ID,
			Price:      9.0,
			Modified:   time.Now(),
		},
		{
			ID:         fs.nextItemID + 4,
			Name:       "Resistor SMT 200",
			LocationID: electronics.ID,
			Price:      0.2,
			Modified:   time.Now(),
		},
	}
	fs.nextItemID += 5
}

// GetAllItems returns all items in the inventory
func (fs *FileStore) GetAllItems() []model.Item {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	// Return a copy to prevent potential race conditions
	items := make([]model.Item, len(fs.inventory.Items))
	copy(items, fs.inventory.Items)
	
	return items
}

// GetAllLocations returns all locations in the inventory
func (fs *FileStore) GetAllLocations() []model.Location {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	// Return a copy to prevent potential race conditions
	locations := make([]model.Location, len(fs.inventory.Locations))
	copy(locations, fs.inventory.Locations)
	
	return locations
}

// GetItemByID returns an item by its ID
func (fs *FileStore) GetItemByID(id uint) (model.Item, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	for _, item := range fs.inventory.Items {
		if item.ID == id {
			return item, nil
		}
	}
	
	return model.Item{}, errors.New("item not found")
}

// GetLocationByID returns a location by its ID
func (fs *FileStore) GetLocationByID(id uint) (model.Location, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	for _, location := range fs.inventory.Locations {
		if location.ID == id {
			return location, nil
		}
	}
	
	return model.Location{}, errors.New("location not found")
}

// AddItem adds a new item to the inventory
func (fs *FileStore) AddItem(item model.CreateItem) (model.Item, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Validate location exists
	locationExists := false
	for _, loc := range fs.inventory.Locations {
		if loc.ID == item.LocationID {
			locationExists = true
			break
		}
	}
	
	if !locationExists {
		return model.Item{}, errors.New("location does not exist")
	}
	
	// Create new item
	newItem := model.Item{
		ID:         fs.nextItemID,
		Name:       item.Name,
		LocationID: item.LocationID,
		Price:      item.Price,
		Modified:   time.Now(),
	}
	
	fs.nextItemID++
	fs.inventory.Items = append(fs.inventory.Items, newItem)
	
	// Save changes
	if err := fs.saveInventory(); err != nil {
		return model.Item{}, err
	}
	
	return newItem, nil
}

// AddLocation adds a new location to the inventory
func (fs *FileStore) AddLocation(location model.CreateLocation) (model.Location, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	newLocation := model.Location{
		ID:       fs.nextLocationID,
		Name:     location.Name,
		Modified: time.Now(),
	}
	
	fs.nextLocationID++
	fs.inventory.Locations = append(fs.inventory.Locations, newLocation)
	
	// Save changes
	if err := fs.saveInventory(); err != nil {
		return model.Location{}, err
	}
	
	return newLocation, nil
}

// UpdateItem updates an existing item
func (fs *FileStore) UpdateItem(id uint, update model.CreateItem) (model.Item, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Find the item
	var foundIndex = -1
	for i, item := range fs.inventory.Items {
		if item.ID == id {
			foundIndex = i
			break
		}
	}
	
	if foundIndex == -1 {
		return model.Item{}, errors.New("item not found")
	}
	
	// Validate location exists
	locationExists := false
	for _, loc := range fs.inventory.Locations {
		if loc.ID == update.LocationID {
			locationExists = true
			break
		}
	}
	
	if !locationExists {
		return model.Item{}, errors.New("location does not exist")
	}
	
	// Update the item
	fs.inventory.Items[foundIndex].Name = update.Name
	fs.inventory.Items[foundIndex].LocationID = update.LocationID
	fs.inventory.Items[foundIndex].Price = update.Price
	fs.inventory.Items[foundIndex].Modified = time.Now()
	
	updatedItem := fs.inventory.Items[foundIndex]
	
	// Save changes
	if err := fs.saveInventory(); err != nil {
		return model.Item{}, err
	}
	
	return updatedItem, nil
}

// DeleteItem removes an item from the inventory
func (fs *FileStore) DeleteItem(id uint) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Find the item
	var foundIndex = -1
	for i, item := range fs.inventory.Items {
		if item.ID == id {
			foundIndex = i
			break
		}
	}
	
	if foundIndex == -1 {
		return errors.New("item not found")
	}
	
	// Remove the item
	fs.inventory.Items = append(fs.inventory.Items[:foundIndex], fs.inventory.Items[foundIndex+1:]...)
	
	// Save changes
	if err := fs.saveInventory(); err != nil {
		return err
	}
	
	return nil
}

// DeleteLocation removes a location from the inventory
func (fs *FileStore) DeleteLocation(id uint) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Find the location
	var foundIndex = -1
	for i, location := range fs.inventory.Locations {
		if location.ID == id {
			foundIndex = i
			break
		}
	}
	
	if foundIndex == -1 {
		return errors.New("location not found")
	}
	
	// Check if any items are using this location
	for _, item := range fs.inventory.Items {
		if item.LocationID == id {
			return errors.New("cannot delete location: it is still being used by items")
		}
	}
	
	// Remove the location
	fs.inventory.Locations = append(fs.inventory.Locations[:foundIndex], fs.inventory.Locations[foundIndex+1:]...)
	
	// Save changes
	if err := fs.saveInventory(); err != nil {
		return err
	}
	
	return nil
}

// SearchItems searches for items by name
func (fs *FileStore) SearchItems(query string) []model.Item {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	if query == "" {
		return fs.GetAllItems()
	}
	
	query = strings.ToLower(query)
	var results []model.Item
	
	for _, item := range fs.inventory.Items {
		if strings.Contains(strings.ToLower(item.Name), query) {
			results = append(results, item)
		}
	}
	
	return results
}

// GetItemsWithLocations returns all items with their location names
func (fs *FileStore) GetItemsWithLocations() []model.ItemWithLocation {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	var itemsWithLocations []model.ItemWithLocation
	
	for _, item := range fs.inventory.Items {
		locationName := "Unknown"
		
		// Find the location name
		for _, loc := range fs.inventory.Locations {
			if loc.ID == item.LocationID {
				locationName = loc.Name
				break
			}
		}
		
		itemWithLocation := item.ToItemWithLocation(locationName)
		itemsWithLocations = append(itemsWithLocations, itemWithLocation)
	}
	
	return itemsWithLocations
}
