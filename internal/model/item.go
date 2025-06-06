package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location represents a physical location where items are stored
type Location struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Modified time.Time          `json:"modified" bson:"modified"`
}

// Item represents an inventory item in the Lab
type Item struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	LocationID primitive.ObjectID `json:"location_id" bson:"location_id"`
	Price      float64            `json:"price" bson:"price"`
	Number     int                `json:"number" bson:"number"` // Quantity/Count
	Modified   time.Time          `json:"modified" bson:"modified"`
}

// CreateItem represents the data needed to create a new item
type CreateItem struct {
	Name       string  `json:"name"`
	LocationID string  `json:"location_id"` // String to receive from frontend
	Price      float64 `json:"price"`
	Number     int     `json:"number"` // Quantity/Count
}

// CreateLocation represents the data needed to create a new location
type CreateLocation struct {
	Name string `json:"name"`
}

// ItemWithLocation represents an item with its location name (for display)
type ItemWithLocation struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Location string             `json:"location" bson:"location"`
	Price    float64            `json:"price" bson:"price"`
	Number   int                `json:"number" bson:"number"` // Quantity/Count
	Modified time.Time          `json:"modified" bson:"modified"`
}

// Inventory represents the entire inventory with items and locations (for compatibility)
type Inventory struct {
	Items     []Item     `json:"items"`
	Locations []Location `json:"locations"`
}

// Config represents application configuration
type Config struct {
	MongoURI     string `json:"mongo_uri"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
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
