// script.js

const accessToken = localStorage.getItem('accessToken');

// Sample function to simulate fetching data from the API
async function fetchData() {
    // Replace with your actual API URL
    const apiUrl = 'http://192.168.1.9:8009/servicing/service/get-list-v2'; 

    try {
        const response = await fetch(apiUrl, {
            method: 'GET',
            headers: {
                'Authorization': 'Bearer '+accessToken, // Replace with your actual token
                'Content-Type': 'application/json',
            }
        });

        // Check if the response is ok (status code 200)
        if (!response.ok && response.status !== 401) {
            throw new Error('Network response was not ok');
        }else if (response.status === 401) {
            // Handle unauthorized access
            console.log('Unauthorized access');
            // Redirect to login page
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            window.location.href = 'index.html';
        }

        const data = await response.json(); // Parse JSON data
        displayData(data); // Call function to display the data
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}

// Function to display the data in a structured format
function displayData(response) {
    const serviceDataDiv = document.getElementById('serviceData');
    serviceDataDiv.innerHTML = ''; // Clear existing data

    if (response.message === "success") {
        for (const category in response.data) {
            const categoryDiv = document.createElement('div');
            categoryDiv.innerHTML = `<h2>${category}</h2>`;
            const table = document.createElement('table');
            table.className = 'table table-bordered';
            table.innerHTML = `
                <thead class="table-dark">
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Price</th>
                        <th>Description</th>
                        <th>Image</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                </tbody>
            `;
            const tbody = table.querySelector('tbody');

            response.data[category].forEach(item => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${item.id}</td>
                    <td>${item.name}</td>
                    <td>${item.price}</td>
                    <td>${item.description}</td>
                    <td><img src="${item.url}" alt="${item.name}" style="width: 100px;"></td>
                    <td>
                        <button class="btn btn-primary btn-sm" onclick="updateItem(${item.id})">Update</button>
                        <button class="btn btn-danger btn-sm" onclick="deleteItem(${item.id})">Delete</button>
                    </td>
                `;
                tbody.appendChild(row);
            });

            categoryDiv.appendChild(table);
            serviceDataDiv.appendChild(categoryDiv);
        }
    } else {
        serviceDataDiv.innerHTML = `<p>No data available.</p>`;
    }
}

// Function to simulate updating an item
function updateItem(id) {
    console.log(`Update item with ID: ${id}`);
    // Here you would typically open a modal or redirect to an update form
}

// Function to simulate deleting an item
function deleteItem(id) {
    console.log(`Delete item with ID: ${id}`);
    // Here you would typically send a DELETE request to the API
}

// Call the fetchData function when the script loads
fetchData();