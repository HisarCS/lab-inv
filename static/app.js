// Global state
let inventory = {
    items: [],
    locations: []
};

// API configuration
const API_BASE = '/api';

// Application initialization
document.addEventListener('DOMContentLoaded', function() {
    initializeApplication();
});

// Initialize the entire application
function initializeApplication() {
    // Load data from MongoDB
    loadInventoryData();
    
    // Set up event listeners
    setupEventListeners();
    
    // Initialize UI components
    initializeUI();
    
    console.log('Lab Inventory System initialized');
}

// ===== EVENT LISTENERS =====

function setupEventListeners() {
    // Tab switching
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', handleTabSwitch);
    });

    // Modal controls
    document.getElementById('add-item-button').addEventListener('click', openAddItemModal);
    document.getElementById('add-location-button').addEventListener('click', openAddLocationModal);
    
    // Close modal buttons
    document.querySelectorAll('.close-btn, .close-modal').forEach(btn => {
        btn.addEventListener('click', closeAllModals);
    });

    // Form submissions
    document.getElementById('add-item-form').addEventListener('submit', handleAddItem);
    document.getElementById('add-location-form').addEventListener('submit', handleAddLocation);
    document.getElementById('edit-item-form').addEventListener('submit', handleEditItem);
    document.getElementById('edit-location-form').addEventListener('submit', handleEditLocation);
    
    // Search functionality
    document.getElementById('search-input').addEventListener('input', handleSearch);
    
    // Sedat button
    document.getElementById('sedat-button').addEventListener('click', handleSedatClick);

    // Close modal when clicking outside
    window.addEventListener('click', handleModalOutsideClick);
}

// ===== TAB FUNCTIONALITY =====

function handleTabSwitch(event) {
    const clickedTab = event.target;
    const tabName = clickedTab.dataset.tab;
    
    // Remove active class from all tabs and sections
    document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
    document.querySelectorAll('.section').forEach(section => section.classList.remove('active'));
    
    // Add active class to clicked tab and corresponding section
    clickedTab.classList.add('active');
    document.getElementById(tabName + '-section').classList.add('active');
    
    console.log('Switched to tab:', tabName);
}

// ===== MODAL FUNCTIONALITY =====

function openAddItemModal() {
    populateLocationDropdown();
    document.getElementById('add-item-form').reset();
    document.getElementById('add-item-modal').style.display = 'block';
}

function openAddLocationModal() {
    document.getElementById('add-location-form').reset();
    document.getElementById('add-location-modal').style.display = 'block';
}

function closeAllModals() {
    const modals = [
        'add-item-modal',
        'add-location-modal', 
        'edit-item-modal',
        'edit-location-modal'
    ];
    
    modals.forEach(modalId => {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'none';
        }
    });
}

function handleModalOutsideClick(event) {
    if (event.target.classList.contains('modal')) {
        event.target.style.display = 'none';
    }
}

// ===== DATA LOADING =====

async function loadInventoryData() {
    try {
        showLoadingStates();
        
        // Load items and locations in parallel
        const [itemsResponse, locationsResponse] = await Promise.all([
            fetch(`${API_BASE}/items`),
            fetch(`${API_BASE}/locations`)
        ]);

        if (!itemsResponse.ok || !locationsResponse.ok) {
            throw new Error('Failed to load data from server');
        }

        const items = await itemsResponse.json();
        const locations = await locationsResponse.json();

        inventory = {
            items: items || [],
            locations: locations || []
        };

        console.log('Loaded data:', { 
            items: inventory.items.length, 
            locations: inventory.locations.length 
        });

        renderAllData();
        
    } catch (error) {
        console.error('Error loading inventory data:', error);
        showError('Failed to load inventory data. Please check your connection.');
        hideLoadingStates();
    }
}

// ===== RENDERING FUNCTIONS =====

function renderAllData() {
    renderItems(inventory.items);
    renderLocations(inventory.locations);
}

function renderItems(items) {
    const itemsList = document.getElementById('items-list');
    const itemsLoading = document.getElementById('items-loading');
    const itemsEmpty = document.getElementById('items-empty');
    const itemsTable = document.getElementById('items-table');
    
    // Clear existing content
    itemsList.innerHTML = '';
    
    setTimeout(() => {
        hideLoadingStates();
        
        if (!items || items.length === 0) {
            // Show empty state
            itemsEmpty.style.display = 'flex';
            itemsTable.style.display = 'none';
        } else {
            // Show table with data
            itemsEmpty.style.display = 'none';
            itemsTable.style.display = 'table';
            
            items.forEach(item => {
                const locationName = getLocationNameById(item.location_id);
                const row = createItemRow(item, locationName);
                itemsList.appendChild(row);
            });
        }
    }, 200);
}

function renderLocations(locations) {
    const locationsList = document.getElementById('locations-list');
    const locationsLoading = document.getElementById('locations-loading');
    const locationsEmpty = document.getElementById('locations-empty');
    const locationsTable = document.getElementById('locations-table');
    
    // Clear existing content
    locationsList.innerHTML = '';
    
    setTimeout(() => {
        hideLoadingStates();
        
        if (!locations || locations.length === 0) {
            // Show empty state
            locationsEmpty.style.display = 'flex';
            locationsTable.style.display = 'none';
        } else {
            // Show table with data
            locationsEmpty.style.display = 'none';
            locationsTable.style.display = 'table';
            
            locations.forEach(location => {
                const itemCount = getItemCountForLocation(location.id);
                const row = createLocationRow(location, itemCount);
                locationsList.appendChild(row);
            });
        }
    }, 200);
}

// ===== ROW CREATION =====

function createItemRow(item, locationName) {
    const row = document.createElement('tr');
    row.innerHTML = `
        <td>${escapeHtml(item.name)}</td>
        <td>${escapeHtml(locationName)}</td>
        <td>${item.number || 0}</td>
        <td class="price">$${parseFloat(item.price || 0).toFixed(2)}</td>
        <td>
            <button class="action-btn" onclick="editItem('${item.id}')">EDIT</button>
            <button class="action-btn delete" onclick="confirmDeleteItem('${item.id}', '${escapeHtml(item.name)}')">DELETE</button>
        </td>
    `;
    
    return row;
}

function createLocationRow(location, itemCount) {
    const row = document.createElement('tr');
    row.innerHTML = `
        <td>${escapeHtml(location.name)}</td>
        <td>${itemCount}</td>
        <td>
            <button class="action-btn" onclick="editLocation('${location.id}')">EDIT</button>
            <button class="action-btn delete" onclick="confirmDeleteLocation('${location.id}', '${escapeHtml(location.name)}')">DELETE</button>
        </td>
    `;
    
    return row;
}

// ===== FORM HANDLING =====

async function handleAddItem(event) {
    event.preventDefault();
    
    const formData = getItemFormData();
    
    if (!validateItemForm(formData)) {
        return;
    }
    
    try {
        const newItem = await createItemAPI(formData);
        
        await loadInventoryData();
        closeAllModals();
        showSuccess(`Item "${newItem.name}" added successfully!`);
        
    } catch (error) {
        showError(`Failed to add item: ${error.message}`);
    }
}

async function handleAddLocation(event) {
    event.preventDefault();
    
    const formData = getLocationFormData();
    
    if (!validateLocationForm(formData)) {
        return;
    }
    
    try {
        const newLocation = await createLocationAPI(formData);
        
        await loadInventoryData();
        closeAllModals();
        showSuccess(`Location "${newLocation.name}" added successfully!`);
        
    } catch (error) {
        showError(`Failed to add location: ${error.message}`);
    }
}

// ===== FORM DATA EXTRACTION =====

function getItemFormData() {
    return {
        name: document.getElementById('item-name').value.trim(),
        location_id: document.getElementById('item-location').value,
        number: parseInt(document.getElementById('item-number').value) || 0,
        price: parseFloat(document.getElementById('item-price').value) || 0
    };
}

function getLocationFormData() {
    return {
        name: document.getElementById('location-name').value.trim()
    };
}

// ===== FORM VALIDATION =====

function validateItemForm(data) {
    if (!data.name) {
        showError('Item name is required');
        return false;
    }
    
    if (!data.location_id) {
        showError('Location selection is required');
        return false;
    }
    
    if (data.number < 0) {
        showError('Number must be non-negative');
        return false;
    }
    
    if (data.price < 0) {
        showError('Price must be non-negative');
        return false;
    }
    
    return true;
}

function validateLocationForm(data) {
    if (!data.name) {
        showError('Location name is required');
        return false;
    }
    
    return true;
}

// ===== DELETE OPERATIONS =====

function confirmDeleteItem(itemId, itemName) {
    if (confirm(`Are you sure you want to delete "${itemName}"?`)) {
        deleteItem(itemId);
    }
}

function confirmDeleteLocation(locationId, locationName) {
    if (confirm(`Are you sure you want to delete "${locationName}"?`)) {
        deleteLocation(locationId);
    }
}

async function deleteItem(itemId) {
    try {
        await deleteItemAPI(itemId);
        await loadInventoryData();
        showSuccess('Item deleted successfully!');
    } catch (error) {
        showError(`Failed to delete item: ${error.message}`);
    }
}

async function deleteLocation(locationId) {
    try {
        await deleteLocationAPI(locationId);
        await loadInventoryData();
        showSuccess('Location deleted successfully!');
    } catch (error) {
        showError(`Failed to delete location: ${error.message}`);
    }
}

// ===== SEARCH FUNCTIONALITY =====

function handleSearch(event) {
    const searchTerm = event.target.value.toLowerCase().trim();
    
    if (!searchTerm) {
        renderItems(inventory.items);
        return;
    }
    
    const filteredItems = inventory.items.filter(item => 
        item.name.toLowerCase().includes(searchTerm)
    );
    
    renderItems(filteredItems);
}

// ===== API FUNCTIONS =====

async function createItemAPI(itemData) {
    const response = await fetch(`${API_BASE}/items`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(itemData)
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }

    return await response.json();
}

async function createLocationAPI(locationData) {
    const response = await fetch(`${API_BASE}/locations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(locationData)
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }

    return await response.json();
}

async function deleteItemAPI(itemId) {
    const response = await fetch(`${API_BASE}/items/${itemId}`, {
        method: 'DELETE'
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }
}

async function deleteLocationAPI(locationId) {
    const response = await fetch(`${API_BASE}/locations/${locationId}`, {
        method: 'DELETE'
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }
}

// ===== UTILITY FUNCTIONS =====

function getLocationNameById(locationId) {
    const location = inventory.locations.find(loc => loc.id === locationId);
    return location ? location.name : 'Unknown';
}

function getItemCountForLocation(locationId) {
    return inventory.items.filter(item => item.location_id === locationId).length;
}

function populateLocationDropdown() {
    const dropdown = document.getElementById('item-location');
    dropdown.innerHTML = '<option value="">Select location</option>';
    
    inventory.locations.forEach(location => {
        const option = document.createElement('option');
        option.value = location.id;
        option.textContent = location.name;
        dropdown.appendChild(option);
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ===== UI STATE MANAGEMENT =====

function showLoadingStates() {
    document.getElementById('items-loading').style.display = 'flex';
    document.getElementById('locations-loading').style.display = 'flex';
}

function hideLoadingStates() {
    document.getElementById('items-loading').style.display = 'none';
    document.getElementById('locations-loading').style.display = 'none';
}

function initializeUI() {
    // Check database connection status
    checkDatabaseConnection();
}

async function checkDatabaseConnection() {
    try {
        const response = await fetch(`${API_BASE}/items`);
        const statusElement = document.getElementById('db-status');
        
        if (response.ok) {
            statusElement.textContent = 'Database Connected';
            statusElement.classList.remove('error');
        } else {
            throw new Error('Connection failed');
        }
    } catch (error) {
        const statusElement = document.getElementById('db-status');
        statusElement.textContent = 'Database Error';
        statusElement.classList.add('error');
        console.error('Database connection error:', error);
    }
}

// ===== NOTIFICATION SYSTEM =====

function showSuccess(message) {
    console.log('SUCCESS:', message);
    // Could be enhanced with toast notifications
}

function showError(message) {
    console.error('ERROR:', message);
    alert(`Error: ${message}`);
}

// ===== EDIT FUNCTIONALITY =====

function editItem(itemId) {
    const item = inventory.items.find(item => item.id === itemId);
    if (!item) {
        showError('Item not found');
        return;
    }
    
    // Populate edit form with current data
    document.getElementById('edit-item-name').value = item.name;
    document.getElementById('edit-item-location').value = item.location_id;
    document.getElementById('edit-item-number').value = item.number || 0;
    document.getElementById('edit-item-price').value = item.price || 0;
    
    // Store item ID for update
    document.getElementById('edit-item-form').dataset.itemId = itemId;
    
    // Populate location dropdown and show modal
    populateEditLocationDropdown();
    document.getElementById('edit-item-modal').style.display = 'block';
}

function editLocation(locationId) {
    const location = inventory.locations.find(loc => loc.id === locationId);
    if (!location) {
        showError('Location not found');
        return;
    }
    
    // Populate edit form with current data
    document.getElementById('edit-location-name').value = location.name;
    
    // Store location ID for update
    document.getElementById('edit-location-form').dataset.locationId = locationId;
    
    // Show modal
    document.getElementById('edit-location-modal').style.display = 'block';
}

async function handleEditItem(event) {
    event.preventDefault();
    
    const form = event.target;
    const itemId = form.dataset.itemId;
    
    const formData = getEditItemFormData();
    
    if (!validateItemForm(formData)) {
        return;
    }
    
    try {
        const updatedItem = await updateItemAPI(itemId, formData);
        
        await loadInventoryData();
        closeAllModals();
        showSuccess(`Item "${updatedItem.name}" updated successfully!`);
        
    } catch (error) {
        showError(`Failed to update item: ${error.message}`);
    }
}

async function handleEditLocation(event) {
    event.preventDefault();
    
    const form = event.target;
    const locationId = form.dataset.locationId;
    
    const formData = getEditLocationFormData();
    
    if (!validateLocationForm(formData)) {
        return;
    }
    
    try {
        const updatedLocation = await updateLocationAPI(locationId, formData);
        
        await loadInventoryData();
        closeAllModals();
        showSuccess(`Location "${updatedLocation.name}" updated successfully!`);
        
    } catch (error) {
        showError(`Failed to update location: ${error.message}`);
    }
}

function getEditItemFormData() {
    return {
        name: document.getElementById('edit-item-name').value.trim(),
        location_id: document.getElementById('edit-item-location').value,
        number: parseInt(document.getElementById('edit-item-number').value) || 0,
        price: parseFloat(document.getElementById('edit-item-price').value) || 0
    };
}

function getEditLocationFormData() {
    return {
        name: document.getElementById('edit-location-name').value.trim()
    };
}

function populateEditLocationDropdown() {
    const dropdown = document.getElementById('edit-item-location');
    dropdown.innerHTML = '<option value="">Select location</option>';
    
    inventory.locations.forEach(location => {
        const option = document.createElement('option');
        option.value = location.id;
        option.textContent = location.name;
        dropdown.appendChild(option);
    });
}

async function updateItemAPI(itemId, itemData) {
    const response = await fetch(`${API_BASE}/items/${itemId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(itemData)
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }

    return await response.json();
}

async function updateLocationAPI(locationId, locationData) {
    const response = await fetch(`${API_BASE}/locations/${locationId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(locationData)
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
    }

    return await response.json();
}

// ===== SEDAT BUTTON =====

function handleSedatClick() {
    const messages = [
        'Hello from Sedat!',
        'MongoDB is working perfectly!',
        'Professional Lab Inventory System',
        'Made with Go + MongoDB + JavaScript',
        'Sedat approves this application!'
    ];
    
    const randomMessage = messages[Math.floor(Math.random() * messages.length)];
    alert(randomMessage);
}

// ===== DEBUGGING AND DEVELOPMENT =====

// Expose some functions to global scope for debugging
window.labInventory = {
    loadData: loadInventoryData,
    getInventory: () => inventory,
    checkConnection: checkDatabaseConnection,
    editItem: editItem,
    editLocation: editLocation,
    closeModals: closeAllModals,
    version: '2.1.0-mongodb-with-edit'
};

console.log('Lab Inventory System loaded - MongoDB Edition');
console.log('Use window.labInventory for debugging');
