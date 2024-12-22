// script.js

const accessToken = localStorage.getItem('accessToken');

document.addEventListener('DOMContentLoaded', function() {
    fetch('http://192.168.1.9:8008/barber/get-list')
        .then(response => response.json())
        .then(data => {
            if (data.status === 200 && data.message === "success") {
                const barberSelect = document.getElementById('barber_id');
                data.data.forEach(barber => {
                    const option = document.createElement('option');
                    option.value = barber.id;
                    option.textContent = barber.fullName; // Use 'fullName' property
                    barberSelect.appendChild(option);
                });
            } else {
                console.error('Error fetching barbers:', data.message);
            }
        })
        .catch(error => {
            console.error('Error fetching barbers:', error);
        });
});

document.getElementById('timeslot-form').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const barberId = parseInt(document.getElementById('barber_id').value, 10);
    const bookedDateInput = document.getElementById('booked_date').value;
    const bookedDate = new Date(bookedDateInput);
    const formattedBookedDate = ('0' + bookedDate.getDate()).slice(-2) + '-' + 
                                ('0' + (bookedDate.getMonth() + 1)).slice(-2) + '-' + 
                                bookedDate.getFullYear();    
    const startTime = document.getElementById('timeInput').value;
    
    const data = {
        barber_id: barberId,
        booked_date: formattedBookedDate,
        start_time: startTime
    };
    
    fetch('http://192.168.1.9:8003/timeslot/find', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer' +accessToken,
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        // Handle the response data here
        displayData(data); // Call function to display the data
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});

var barberName = '';
document.getElementById('barber_id').addEventListener('change', function() {
    const selectedOption = this.options[this.selectedIndex];
    barberName = selectedOption.textContent;
    console.log('Selected Barber Name:', barberName);
    // You can use barberName as needed
});


// Sample function to simulate fetching data from the API
// async function fetchData() {
//     // Replace with your actual API URL
//     const apiUrl = 'http://192.168.1.9:8003/timeslot/find'; 

//     try {
//         const response = await fetch(apiUrl, {
//             method: 'POST',
//             headers: {
//                 'Authorization': 'Bearer '+accessToken, // Replace with your actual token
//                 'Content-Type': 'application/json',
//             }
//         });

//         // Check if the response is ok (status code 200)
//         if (!response.ok && response.status !== 401) {
//             throw new Error('Network response was not ok');
//         }else if (response.status === 401) {
//             // Handle unauthorized access
//             console.log('Unauthorized access');
//             // Redirect to login page
//             localStorage.removeItem('username');
//             localStorage.removeItem('token');
//             window.location.href = 'index.html';
//         }

//         const data = await response.json(); // Parse JSON data
//         displayData(data); // Call function to display the data
//     } catch (error) {
//         console.error('Error fetching data:', error);
//     }
// }

// Function to display the data in a structured format
function displayData(response) {
    console.log('Displaying data:', response);
    const timeslotDataDiv = document.getElementById('timeslotData');
    timeslotDataDiv.innerHTML = ''; // Clear existing data

    if (response.message === "success" || response.status === 200) {
        const table = document.createElement('table');
        table.className = 'table table-bordered';
        table.innerHTML = `
            <thead class="table-dark">
                <tr>
                    <th>ID</th>
                    <th>Start Time</th>
                    <th>Date</th>
                    <th>Status</th>
                    <th>Barber</th>
                    <th>Created at</th>
                    <th>Updated at</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
            </tbody>
        `;
        const tbody = table.querySelector('tbody');

        response.data.forEach(timeslot => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${timeslot.id}</td>
                <td>${timeslot.start_time}</td>
                <td>${timeslot.booked_date}</td>
                <td>${timeslot.status}</td>
                <td>${barberName}</td>
                <td>${formatTimestamp(timeslot.created_at)}</td>
                <td>${formatTimestamp(timeslot.updated_at)}</td>
                <td>
                    <button class="btn btn-primary btn-sm" onclick="updateBarber(${timeslot.id})">Update</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteBarber(${timeslot.id})">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });

        timeslotDataDiv.appendChild(table);
    } else {
        timeslotDataDiv.innerHTML = `<p>No data available.</p>`;
    }
}

function formatTimestamp(timestamp) {
    const date = new Date(timestamp * 1000); // Convert from seconds to milliseconds
    const day = ('0' + date.getDate()).slice(-2);
    const month = ('0' + (date.getMonth() + 1)).slice(-2);
    const year = date.getFullYear();
    const hours = ('0' + date.getHours()).slice(-2);
    const minutes = ('0' + date.getMinutes()).slice(-2);
    const seconds = ('0' + date.getSeconds()).slice(-2);
    return `${day}-${month}-${year} ${hours}:${minutes}:${seconds}`;
}

const startHour = 8;
const endHour = 18;
const datalist = document.getElementById('timeOptions');

for (let hour = startHour; hour <= endHour; hour++) {
  for (let minute of [0, 30]) {
    const timeValue = `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`;
    const option = document.createElement('option');
    option.value = timeValue;
    datalist.appendChild(option);
  }
}

document.getElementById('timeInput').addEventListener('input', function() {
    const timeInput = document.getElementById('timeInput');
    const timeOptions = document.getElementById('timeOptions').options;
    for (let i = 0; i < timeOptions.length; i++) {
        if (timeOptions[i].value === timeInput.value) {
            timeInput.value = timeOptions[i].value;
            break;
        }
    }
});

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

    const apiUrl = 'http://192.168.1.9:8008/auth/register'; // Replace with your actual API URL

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
