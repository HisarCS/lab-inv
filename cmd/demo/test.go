// lab-inv/cmd/demo/test.go

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"lab-inv/internal/model"
	"lab-inv/internal/storage"
)

func main() {
	// Create a new FileStore instance
	dataDir := "./data"
	fileStore, err := storage.NewFileStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize FileStore: %v", err)
	}
	
	fmt.Println("Lab Inventory System")
	fmt.Println("-------------------")
	fmt.Println("Data directory:", dataDir)

	// Show inventory summary
	items := fileStore.GetAllItems()
	locations := fileStore.GetAllLocations()

	fmt.Printf("Total items: %d\n", len(items))
	fmt.Printf("Total locations: %d\n", len(locations))

	// Start the interactive loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nCommands:")
		fmt.Println("1. List all items")
		fmt.Println("2. List all locations")
		fmt.Println("3. Add new item")
		fmt.Println("4. Add new location")
		fmt.Println("5. Search items")
		fmt.Println("6. Delete item")
		fmt.Println("7. Delete location")
		fmt.Println("0. Exit")
		fmt.Print("\nEnter command: ")

		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "0", "exit", "quit":
			fmt.Println("Goodbye!")
			return

		case "1":
			listAllItems(fileStore)

		case "2":
			listAllLocations(fileStore)

		case "3":
			addNewItem(fileStore, scanner)

		case "4":
			addNewLocation(fileStore, scanner)

		case "5":
			searchItems(fileStore, scanner)

		case "6":
			deleteItem(fileStore, scanner)

		case "7":
			deleteLocation(fileStore, scanner)

		default:
			fmt.Println("Unknown command")
		}
	}
}

// listAllItems displays all items with their locations
func listAllItems(fileStore *storage.FileStore) {
	items := fileStore.GetItemsWithLocations()

	if len(items) == 0 {
		fmt.Println("No items found")
		return
	}

	fmt.Println("\nAll Items:")
	fmt.Printf("%-5s | %-30s | %-20s | %-10s\n", "ID", "Name", "Location", "Price")
	fmt.Println(strings.Repeat("-", 70))
	
	for _, item := range items {
		fmt.Printf("%-5d | %-30s | %-20s | $%.2f\n", 
			item.ID, item.Name, item.Location, item.Price)
	}
}

// listAllLocations displays all locations
func listAllLocations(fileStore *storage.FileStore) {
	locations := fileStore.GetAllLocations()

	if len(locations) == 0 {
		fmt.Println("No locations found")
		return
	}

	fmt.Println("\nAll Locations:")
	fmt.Printf("%-5s | %-30s\n", "ID", "Name")
	fmt.Println(strings.Repeat("-", 40))
	
	for _, location := range locations {
		fmt.Printf("%-5d | %-30s\n", location.ID, location.Name)
	}
}

// addNewItem adds a new item to the inventory
func addNewItem(fileStore *storage.FileStore, scanner *bufio.Scanner) {
	// List available locations first
	locations := fileStore.GetAllLocations()
	if len(locations) == 0 {
		fmt.Println("No locations available. Please add a location first.")
		return
	}

	fmt.Println("\nAvailable Locations:")
	for _, loc := range locations {
		fmt.Printf("%d. %s\n", loc.ID, loc.Name)
	}

	// Get item details
	var name string
	var locationID uint
	var price float64

	fmt.Print("\nEnter item name: ")
	scanner.Scan()
	name = scanner.Text()
	if name == "" {
		fmt.Println("Item name cannot be empty")
		return
	}

	fmt.Print("Enter location ID: ")
	scanner.Scan()
	locIDStr := scanner.Text()
	locID, err := strconv.ParseUint(locIDStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid location ID")
		return
	}
	locationID = uint(locID)

	// Verify location exists
	locationExists := false
	for _, loc := range locations {
		if loc.ID == locationID {
			locationExists = true
			break
		}
	}
	if !locationExists {
		fmt.Println("Location with that ID does not exist")
		return
	}

	fmt.Print("Enter price: ")
	scanner.Scan()
	priceStr := scanner.Text()
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("Invalid price")
		return
	}

	// Create the item
	newItem := model.CreateItem{
		Name:       name,
		LocationID: locationID,
		Price:      price,
	}

	item, err := fileStore.AddItem(newItem)
	if err != nil {
		fmt.Printf("Failed to add item: %v\n", err)
		return
	}

	fmt.Printf("Item added successfully with ID %d\n", item.ID)
}

// addNewLocation adds a new location
func addNewLocation(fileStore *storage.FileStore, scanner *bufio.Scanner) {
	fmt.Print("\nEnter location name: ")
	scanner.Scan()
	name := scanner.Text()
	
	if name == "" {
		fmt.Println("Location name cannot be empty")
		return
	}

	// Create the location
	newLocation := model.CreateLocation{
		Name: name,
	}

	location, err := fileStore.AddLocation(newLocation)
	if err != nil {
		fmt.Printf("Failed to add location: %v\n", err)
		return
	}

	fmt.Printf("Location added successfully with ID %d\n", location.ID)
}

// searchItems searches for items by name
func searchItems(fileStore *storage.FileStore, scanner *bufio.Scanner) {
	fmt.Print("\nEnter search term: ")
	scanner.Scan()
	query := scanner.Text()

	items := fileStore.SearchItems(query)
	if len(items) == 0 {
		fmt.Println("No items found matching your search")
		return
	}

	fmt.Printf("\nFound %d items:\n", len(items))
	fmt.Printf("%-5s | %-30s | %-20s | %-10s\n", "ID", "Name", "Location", "Price")
	fmt.Println(strings.Repeat("-", 70))
	
	for _, item := range items {
		// Get location name
		location, err := fileStore.GetLocationByID(item.LocationID)
		locationName := "Unknown"
		if err == nil {
			locationName = location.Name
		}

		fmt.Printf("%-5d | %-30s | %-20s | $%.2f\n", 
			item.ID, item.Name, locationName, item.Price)
	}
}

// deleteItem deletes an item by ID
func deleteItem(fileStore *storage.FileStore, scanner *bufio.Scanner) {
	fmt.Print("\nEnter item ID to delete: ")
	scanner.Scan()
	idStr := scanner.Text()
	
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid item ID")
		return
	}

	// Confirm deletion
	fmt.Print("Are you sure you want to delete this item? (y/n): ")
	scanner.Scan()
	confirm := scanner.Text()
	
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	err = fileStore.DeleteItem(uint(id))
	if err != nil {
		fmt.Printf("Failed to delete item: %v\n", err)
		return
	}

	fmt.Println("Item deleted successfully")
}

// deleteLocation deletes a location by ID
func deleteLocation(fileStore *storage.FileStore, scanner *bufio.Scanner) {
	fmt.Print("\nEnter location ID to delete: ")
	scanner.Scan()
	idStr := scanner.Text()
	
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid location ID")
		return
	}

	// Confirm deletion
	fmt.Print("Are you sure you want to delete this location? (y/n): ")
	scanner.Scan()
	confirm := scanner.Text()
	
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	err = fileStore.DeleteLocation(uint(id))
	if err != nil {
		fmt.Printf("Failed to delete location: %v\n", err)
		return
	}

	fmt.Println("Location deleted successfully")
}
