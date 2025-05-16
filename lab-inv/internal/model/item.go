package model

import (
	"time"
)

// Location represents a physical location where items are stored
type Location struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

// Item represents an inventory item in the Fab Lab
type Item struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	LocationID uint      `json:"location_id"`
	Price      float64   `json:"price"`
	Modified   time.Time `json:"modified"`
}

// CreateItem represents the data needed to create a new item
type CreateItem struct {
	Name       string  `json:"name"`
	LocationID uint    `json:"location_id"`
	Price      float64 `json:"price"`
}

// CreateLocation represents the data needed to create a new location
type CreateLocation struct {
	Name string `json:"name"`
}

// Inventory represents the entire inventory with items and locations
type Inventory struct {
	Items     []Item     `json:"items"`
	Locations []Location `json:"locations"`
}

// ItemWithLocation represents an item with its location information
type ItemWithLocation struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Price    float64   `json:"price"`
	Modified time.Time `json:"modified"`
}

// ToItemWithLocation converts an Item to an ItemWithLocation
func (i *Item) ToItemWithLocation(locationName string) ItemWithLocation {
	return ItemWithLocation{
		ID:       i.ID,
		Name:     i.Name,
		Location: locationName,
		Price:    i.Price,
		Modified: i.Modified,
	}
}

