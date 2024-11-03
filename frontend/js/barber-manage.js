// script.js

const accessToken = localStorage.getItem('accessToken');

// Sample function to simulate fetching data from the API
async function fetchData() {
    // Replace with your actual API URL
    const apiUrl = 'http://localhost:8008/barber/get-list'; 

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
    const barberDataDiv = document.getElementById('barberData');
    barberDataDiv.innerHTML = ''; // Clear existing data

    if (response.message === "success" && response.status === 200) {
        const table = document.createElement('table');
        table.className = 'table table-bordered';
        table.innerHTML = `
            <thead class="table-dark">
                <tr>
                    <th>ID</th>
                    <th>Full Name</th>
                    <th>Username</th>
                    <th>Email</th>
                    <th>Phone Number</th>
                    <th>About</th>
                    <th>Date of Birth</th>
                    <th>Avatar</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
            </tbody>
        `;
        const tbody = table.querySelector('tbody');

        response.data.forEach(barber => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${barber.id}</td>
                <td>${barber.fullName}</td>
                <td>${barber.username}</td>
                <td>${barber.email}</td>
                <td>${barber.phoneNumber}</td>
                <td>${barber.about}</td>
                <td>${barber.dob || 'N/A'}</td>
                <td><img src="${barber.avatar}" alt="${barber.fullName}" style="width: 100px;"></td>
                <td>
                    <button class="btn btn-primary btn-sm" onclick="updateBarber(${barber.id})">Update</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteBarber(${barber.id})">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });

        barberDataDiv.appendChild(table);
    } else {
        barberDataDiv.innerHTML = `<p>No data available.</p>`;
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

// Function to create a new barber
async function createBarber(event) {
    event.preventDefault();
    const form = document.getElementById('createBarberForm');
    const formData = new FormData(form);


    // Append additional fields if necessary
    formData['username']=formData.get('username');
    formData['email']=formData.get('email');
    formData['phoneNumber']=formData.get('phoneNumber');
    formData['fullName']=formData.get('fullName'); //encode to UTF-8
    formData['dob']=formData.get('dob');
    formData['password']=formData.get('password');

    formData.set('role', '2');

    const dataObject = Object.fromEntries(formData.entries());
    dataObject.role = parseInt(dataObject.role, 10);
    const requestBody = JSON.stringify(dataObject);

    const apiUrl = 'http://localhost:8008/auth/register'; // Replace with your actual API URL

    try {
        const response = await fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: requestBody,
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        if (response.ok) {
            const data = await response.json();
            console.log('Barber created successfully:', data);
            // close the modal
            const createBarberModal = bootstrap.Modal.getInstance(document.getElementById('createBarberModal'));
            if (createBarberModal) {
                createBarberModal.hide();
            }            // display success message alert
            alert("Barber created successfully");

            
            fetchData(); // Refresh the barber list
        }
    } catch (error) {
        console.error('Error creating barber:', error);
    }
}

// Add event listener to the create barber form
document.getElementById('createBarberForm').addEventListener('submit', createBarber);

// Call the fetchBarberData function when the script loads

// Call the fetchData function when the script loads
fetchData();