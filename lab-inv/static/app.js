// Global variables
let inventory = {
    items: [],
    locations: []
};

// DOM Content Loaded - Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Load data
    loadInventoryData();
    
    // Tab switching
    initializeTabs();
    
    // Modal functionality
    initializeModals();
    
    // Sedat button
    document.getElementById('sedat-button').addEventListener('click', function() {
        alert('Hello from Sedat!');
    });
    
    // Forms submission
    document.getElementById('add-item-form').addEventListener('submit', handleAddItem);
    document.getElementById('add-location-form').addEventListener('submit', handleAddLocation);
    
    // Search functionality
    document.getElementById('search-input').addEventListener('input', function(e) {
        const searchTerm = e.target.value.toLowerCase();
        filterItems(searchTerm);
    });
});

// ===== DATA FUNCTIONS =====

// Load inventory data from localStorage
function loadInventoryData() {
    // Check if we have data in localStorage
    const savedData = localStorage.getItem('labinv_data');
    
    if (savedData) {
        // Use saved data
        inventory = JSON.parse(savedData);
        renderData();
    } else {
        // Create sample data for first run
        createSampleData();
    }
}

// Create sample data for first run
function createSampleData() {
    // Sample locations
    const locations = [
        { id: 1, name: "Storage Room", modified: new Date().toISOString() },
        { id: 2, name: "Assembly Room", modified: new Date().toISOString() },
        { id: 3, name: "Electronics", modified: new Date().toISOString() }
    ];
    
    // Sample items
    const items = [
        { 
            id: 1, 
            name: "Plywood 2mm 900x600mm Sheet", 
            location_id: 1, 
            price: 11.15, 
            modified: new Date().toISOString()
        },
        { 
            id: 2, 
            name: "MDF 4mm 900x600mm Sheet", 
            location_id: 1, 
            price: 5.67, 
            modified: new Date().toISOString()
        },
        { 
            id: 3, 
            name: "Acrylic 5mm 900x600mm Sheet", 
            location_id: 1, 
            price: 19.34, 
            modified: new Date().toISOString()
        },
        { 
            id: 4, 
            name: "Wood Glue", 
            location_id: 2, 
            price: 9.0, 
            modified: new Date().toISOString()
        },
        { 
            id: 5, 
            name: "Resistor SMT 200", 
            location_id: 3, 
            price: 0.2, 
            modified: new Date().toISOString()
        }
    ];
    
    inventory = { items, locations };
    
    // Save to localStorage
    saveInventoryData();
    
    // Render the data
    renderData();
}

// Save inventory data to localStorage
function saveInventoryData() {
    localStorage.setItem('labinv_data', JSON.stringify(inventory));
}

// ===== UI FUNCTIONS =====

// Initialize tab switching
function initializeTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    tabButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons and sections
            tabButtons.forEach(btn => btn.classList.remove('active'));
            document.querySelectorAll('.content-section').forEach(section => {
                section.classList.remove('active');
            });
            
            // Add active class to clicked button and corresponding section
            this.classList.add('active');
            document.getElementById(this.dataset.tab + '-section').classList.add('active');
        });
    });
}

// Initialize modal functionality
function initializeModals() {
    const addItemButton = document.getElementById('add-item-button');
    const addLocationButton = document.getElementById('add-location-button');
    const addItemModal = document.getElementById('add-item-modal');
    const addLocationModal = document.getElementById('add-location-modal');
    const closeButtons = document.querySelectorAll('.close');
    
    addItemButton.addEventListener('click', () => {
        document.getElementById('add-item-form').reset();
        addItemModal.style.display = 'block';
        populateLocationDropdown();
    });
    
    addLocationButton.addEventListener('click', () => {
        document.getElementById('add-location-form').reset();
        addLocationModal.style.display = 'block';
    });
    
    closeButtons.forEach(button => {
        button.addEventListener('click', function() {
            addItemModal.style.display = 'none';
            addLocationModal.style.display = 'none';
        });
    });
    
    window.addEventListener('click', (e) => {
        if (e.target === addItemModal) {
            addItemModal.style.display = 'none';
        }
        if (e.target === addLocationModal) {
            addLocationModal.style.display = 'none';
        }
    });
}

// Render inventory data to the UI
function renderData() {
    renderItems(inventory.items);
    renderLocations(inventory.locations);
}

// Render items to the UI
function renderItems(items) {
    const itemsList = document.getElementById('items-list');
    const itemsLoading = document.getElementById('items-loading');
    const itemsEmpty = document.getElementById('items-empty');
    const itemsTable = document.getElementById('items-table');
    
    // Show loading, hide others
    itemsLoading.style.display = 'block';
    itemsEmpty.style.display = 'none';
    itemsTable.style.display = 'none';
    
    // Clear existing items
    itemsList.innerHTML = '';
    
    setTimeout(() => {
        // Hide loading
        itemsLoading.style.display = 'none';
        
        if (items.length === 0) {
            // Show empty state
            itemsEmpty.style.display = 'block';
        } else {
            // Show table and populate it
            itemsTable.style.display = 'table';
            
            items.forEach(item => {
                // Find location name
                const location = inventory.locations.find(loc => loc.id === item.location_id);
                const locationName = location ? location.name : 'Unknown';
                
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${item.name}</td>
                    <td>${locationName}</td>
                    <td>$${item.price.toFixed(2)}</td>
                    <td>
                        <button class="action-button edit-button" data-id="${item.id}">Edit</button>
                        <button class="action-button delete-button" data-id="${item.id}">Delete</button>
                    </td>
                `;
                
                // Add event listeners
                const deleteButton = row.querySelector('.delete-button');
                deleteButton.addEventListener('click', () => {
                    if (confirm(`Are you sure you want to delete "${item.name}"?`)) {
                        deleteItem(item.id);
                    }
                });
                
                itemsList.appendChild(row);
            });
        }
    }, 300); // Simulate loading delay
}

// Render locations to the UI
function renderLocations(locations) {
    const locationsList = document.getElementById('locations-list');
    const locationsLoading = document.getElementById('locations-loading');
    const locationsEmpty = document.getElementById('locations-empty');
    const locationsTable = document.getElementById('locations-table');
    
    // Show loading, hide others
    locationsLoading.style.display = 'block';
    locationsEmpty.style.display = 'none';
    locationsTable.style.display = 'none';
    
    // Clear existing locations
    locationsList.innerHTML = '';
    
    setTimeout(() => {
        // Hide loading
        locationsLoading.style.display = 'none';
        
        if (locations.length === 0) {
            // Show empty state
            locationsEmpty.style.display = 'block';
        } else {
            // Show table and populate it
            locationsTable.style.display = 'table';
            
            locations.forEach(location => {
                // Count items in this location
                const itemCount = inventory.items.filter(item => item.location_id === location.id).length;
                
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${location.name}</td>
                    <td>${itemCount}</td>
                    <td>
                        <button class="action-button edit-button" data-id="${location.id}">Edit</button>
                        <button class="action-button delete-button" data-id="${location.id}">Delete</button>
                    </td>
                `;
                
                // Add event listeners
                const deleteButton = row.querySelector('.delete-button');
                deleteButton.addEventListener('click', () => {
                    if (confirm(`Are you sure you want to delete "${location.name}"?`)) {
                        deleteLocation(location.id);
                    }
                });
                
                locationsList.appendChild(row);
            });
        }
    }, 300); // Simulate loading delay
}

// Populate the location dropdown in the add item form
function populateLocationDropdown() {
    const dropdown = document.getElementById('item-location');
    dropdown.innerHTML = '';
    
    if (inventory.locations.length === 0) {
        const option = document.createElement('option');
        option.value = '';
        option.textContent = 'No locations available';
        option.disabled = true;
        option.selected = true;
        dropdown.appendChild(option);
    } else {
        inventory.locations.forEach(location => {
            const option = document.createElement('option');
            option.value = location.id;
            option.textContent = location.name;
            dropdown.appendChild(option);
        });
    }
}

// Filter items based on search term
function filterItems(searchTerm) {
    if (!searchTerm) {
        renderItems(inventory.items);
    } else {
        const filteredItems = inventory.items.filter(item => 
            item.name.toLowerCase().includes(searchTerm)
        );
        renderItems(filteredItems);
    }
}

// ===== CRUD OPERATIONS =====

// Handle form submission for adding a new item
function handleAddItem(e) {
    e.preventDefault();
    
    const name = document.getElementById('item-name').value;
    const locationId = parseInt(document.getElementById('item-location').value);
    const price = parseFloat(document.getElementById('item-price').value);
    
    // Create new item
    const newItem = {
        id: getNextItemId(),
        name: name,
        location_id: locationId,
        price: price,
        modified: new Date().toISOString()
    };
    
    // Add to inventory
    inventory.items.push(newItem);
    
    // Save and refresh
    saveInventoryData();
    renderItems(inventory.items);
    
    // Close modal and reset form
    document.getElementById('add-item-modal').style.display = 'none';
    document.getElementById('add-item-form').reset();
}

// Handle form submission for adding a new location
function handleAddLocation(e) {
    e.preventDefault();
    
    const name = document.getElementById('location-name').value;
    
    // Create new location
    const newLocation = {
        id: getNextLocationId(),
        name: name,
        modified: new Date().toISOString()
    };
    
    // Add to inventory
    inventory.locations.push(newLocation);
    
    // Save and refresh
    saveInventoryData();
    renderLocations(inventory.locations);
    
    // Close modal and reset form
    document.getElementById('add-location-modal').style.display = 'none';
    document.getElementById('add-location-form').reset();
}

// Delete an item
function deleteItem(id) {
    inventory.items = inventory.items.filter(item => item.id !== id);
    saveInventoryData();
    renderItems(inventory.items);
}

// Delete a location
function deleteLocation(id) {
    // Check if location has items
    const hasItems = inventory.items.some(item => item.location_id === id);
    
    if (hasItems) {
        alert('Cannot delete this location because it contains items.');
        return;
    }
    
    inventory.locations = inventory.locations.filter(location => location.id !== id);
    saveInventoryData();
    renderLocations(inventory.locations);
}

// Get next available item ID
function getNextItemId() {
    return inventory.items.length > 0 
        ? Math.max(...inventory.items.map(item => item.id)) + 1 
        : 1;
}

// Get next available location ID
function getNextLocationId() {
    return inventory.locations.length > 0 
        ? Math.max(...inventory.locations.map(location => location.id)) + 1 
        : 1;
}
