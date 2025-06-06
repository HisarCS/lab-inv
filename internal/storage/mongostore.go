package storage

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"lab-inv/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// MongoDB Atlas connection string
	mongoURI = "mongodb+srv://emredayangac:261926Ab@idealab.dxjeojk.mongodb.net/lab_inventory?retryWrites=true&w=majority"
	
	// Database and collection names
	databaseName        = "lab_inventory"
	itemsCollection     = "items"
	locationsCollection = "locations"
	
	// Connection timeout
	connectionTimeout = 30 * time.Second
)

type MongoStore struct {
	client    *mongo.Client
	database  *mongo.Database
	items     *mongo.Collection
	locations *mongo.Collection
}

// NewMongoStore creates a new MongoDB store instance
func NewMongoStore() (*MongoStore, error) {
	log.Println("🔄 Attempting to connect to MongoDB Atlas...")
	log.Printf("📍 Using URI: %s", mongoURI[:50]+"...")
	
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	// Connect to MongoDB Atlas
	log.Println("🔌 Creating MongoDB client...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("❌ Failed to create MongoDB client: %v", err)
		return nil, err
	}

	// Test the connection
	log.Println("🏓 Testing connection with ping...")
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("❌ Failed to ping MongoDB: %v", err)
		return nil, err
	}

	log.Println("✅ Successfully connected to MongoDB Atlas")

	// Get database and collections
	database := client.Database(databaseName)
	items := database.Collection(itemsCollection)
	locations := database.Collection(locationsCollection)

	store := &MongoStore{
		client:    client,
		database:  database,
		items:     items,
		locations: locations,
	}

	// Initialize with sample data if collections are empty
	if err := store.initSampleData(); err != nil {
		log.Printf("Warning: Failed to initialize sample data: %v", err)
	}

	return store, nil
}

// Close closes the MongoDB connection
func (m *MongoStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// initSampleData creates initial data if collections are empty
func (m *MongoStore) initSampleData() error {
	ctx := context.Background()

	// Check if locations collection is empty
	locationCount, err := m.locations.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}

	if locationCount == 0 {
		// Create sample locations
		storageRoom := model.Location{
			ID:       primitive.NewObjectID(),
			Name:     "Storage Room",
			Modified: time.Now(),
		}
		assemblyRoom := model.Location{
			ID:       primitive.NewObjectID(),
			Name:     "Assembly Room",
			Modified: time.Now(),
		}
		electronics := model.Location{
			ID:       primitive.NewObjectID(),
			Name:     "Electronics",
			Modified: time.Now(),
		}

		sampleLocations := []interface{}{storageRoom, assemblyRoom, electronics}

		_, err := m.locations.InsertMany(ctx, sampleLocations)
		if err != nil {
			return err
		}

		log.Println("Created sample locations")

		// Create sample items using the location IDs
		sampleItems := []interface{}{
			model.Item{
				ID:         primitive.NewObjectID(),
				Name:       "Plywood 2mm 900x600mm Sheet",
				LocationID: storageRoom.ID,
				Price:      11.15,
				Number:     25,
				Modified:   time.Now(),
			},
			model.Item{
				ID:         primitive.NewObjectID(),
				Name:       "MDF 4mm 900x600mm Sheet",
				LocationID: storageRoom.ID,
				Price:      5.67,
				Number:     18,
				Modified:   time.Now(),
			},
			model.Item{
				ID:         primitive.NewObjectID(),
				Name:       "Acrylic 5mm 900x600mm Sheet",
				LocationID: storageRoom.ID,
				Price:      19.34,
				Number:     12,
				Modified:   time.Now(),
			},
			model.Item{
				ID:         primitive.NewObjectID(),
				Name:       "Wood Glue",
				LocationID: assemblyRoom.ID,
				Price:      9.0,
				Number:     8,
				Modified:   time.Now(),
			},
			model.Item{
				ID:         primitive.NewObjectID(),
				Name:       "Resistor SMT 200",
				LocationID: electronics.ID,
				Price:      0.2,
				Number:     150,
				Modified:   time.Now(),
			},
		}

		_, err = m.items.InsertMany(ctx, sampleItems)
		if err != nil {
			return err
		}

		log.Println("Created sample items")
	}

	return nil
}

// GetAllItems returns all items from MongoDB
func (m *MongoStore) GetAllItems() ([]model.Item, error) {
	ctx := context.Background()
	
	cursor, err := m.items.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []model.Item{}
	}

	return items, nil
}

// GetAllLocations returns all locations from MongoDB
func (m *MongoStore) GetAllLocations() ([]model.Location, error) {
	ctx := context.Background()
	
	cursor, err := m.locations.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locations []model.Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, err
	}

	if locations == nil {
		locations = []model.Location{}
	}

	return locations, nil
}

// GetItemByID returns a single item by its ID
func (m *MongoStore) GetItemByID(id string) (model.Item, error) {
	ctx := context.Background()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Item{}, errors.New("invalid item ID format")
	}

	var item model.Item
	err = m.items.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Item{}, errors.New("item not found")
		}
		return model.Item{}, err
	}

	return item, nil
}

// GetLocationByID returns a single location by its ID
func (m *MongoStore) GetLocationByID(id string) (model.Location, error) {
	ctx := context.Background()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Location{}, errors.New("invalid location ID format")
	}

	var location model.Location
	err = m.locations.FindOne(ctx, bson.M{"_id": objectID}).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Location{}, errors.New("location not found")
		}
		return model.Location{}, err
	}

	return location, nil
}

// AddItem adds a new item to MongoDB
func (m *MongoStore) AddItem(createItem model.CreateItem) (model.Item, error) {
	ctx := context.Background()
	
	// Convert location ID string to ObjectID
	locationID, err := primitive.ObjectIDFromHex(createItem.LocationID)
	if err != nil {
		return model.Item{}, errors.New("invalid location ID format")
	}

	// Validate that location exists
	_, err = m.GetLocationByID(createItem.LocationID)
	if err != nil {
		return model.Item{}, errors.New("location does not exist")
	}

	// Create new item
	item := model.Item{
		ID:         primitive.NewObjectID(),
		Name:       createItem.Name,
		LocationID: locationID,
		Price:      createItem.Price,
		Number:     createItem.Number,
		Modified:   time.Now(),
	}

	_, err = m.items.InsertOne(ctx, item)
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

// AddLocation adds a new location to MongoDB
func (m *MongoStore) AddLocation(createLocation model.CreateLocation) (model.Location, error) {
	ctx := context.Background()
	
	// Create new location
	location := model.Location{
		ID:       primitive.NewObjectID(),
		Name:     createLocation.Name,
		Modified: time.Now(),
	}

	_, err := m.locations.InsertOne(ctx, location)
	if err != nil {
		return model.Location{}, err
	}

	return location, nil
}

// UpdateItem updates an existing item in MongoDB
func (m *MongoStore) UpdateItem(id string, updateItem model.CreateItem) (model.Item, error) {
	ctx := context.Background()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Item{}, errors.New("invalid item ID format")
	}

	// Convert location ID string to ObjectID
	locationID, err := primitive.ObjectIDFromHex(updateItem.LocationID)
	if err != nil {
		return model.Item{}, errors.New("invalid location ID format")
	}

	// Validate that location exists
	_, err = m.GetLocationByID(updateItem.LocationID)
	if err != nil {
		return model.Item{}, errors.New("location does not exist")
	}

	// Create updated item
	updatedItem := model.Item{
		ID:         objectID,
		Name:       updateItem.Name,
		LocationID: locationID,
		Price:      updateItem.Price,
		Number:     updateItem.Number,
		Modified:   time.Now(),
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updatedItem}

	result, err := m.items.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Item{}, err
	}

	if result.MatchedCount == 0 {
		return model.Item{}, errors.New("item not found")
	}

	return updatedItem, nil
}

// DeleteItem removes an item from MongoDB
func (m *MongoStore) DeleteItem(id string) error {
	ctx := context.Background()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid item ID format")
	}

	result, err := m.items.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("item not found")
	}

	return nil
}

// DeleteLocation removes a location from MongoDB
func (m *MongoStore) DeleteLocation(id string) error {
	ctx := context.Background()
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid location ID format")
	}

	// Check if any items are using this location
	count, err := m.items.CountDocuments(ctx, bson.M{"location_id": objectID})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete location: it is still being used by items")
	}

	result, err := m.locations.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("location not found")
	}

	return nil
}

// SearchItems searches for items by name
func (m *MongoStore) SearchItems(query string) ([]model.Item, error) {
	ctx := context.Background()
	
	if query == "" {
		return m.GetAllItems()
	}

	// Case-insensitive search using regex
	filter := bson.M{
		"name": bson.M{
			"$regex":   query,
			"$options": "i",
		},
	}

	cursor, err := m.items.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []model.Item{}
	}

	return items, nil
}

// GetItemsWithLocations returns all items with their location names
func (m *MongoStore) GetItemsWithLocations() ([]model.ItemWithLocation, error) {
	ctx := context.Background()
	
	// Use MongoDB aggregation to join items with locations
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "locations",
				"localField":   "location_id",
				"foreignField": "_id",
				"as":           "location_info",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$location_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":      1,
				"name":     1,
				"price":    1,
				"number":   1,
				"modified": 1,
				"location": bson.M{
					"$ifNull": []interface{}{"$location_info.name", "Unknown"},
				},
			},
		},
	}

	cursor, err := m.items.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var itemsWithLocations []model.ItemWithLocation
	if err := cursor.All(ctx, &itemsWithLocations); err != nil {
		return nil, err
	}

	if itemsWithLocations == nil {
		itemsWithLocations = []model.ItemWithLocation{}
	}

	return itemsWithLocations, nil
}
