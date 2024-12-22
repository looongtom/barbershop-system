const accessToken = localStorage.getItem('accessToken');
let currentPage = 1;
const pageSize = 10;

document.getElementById('booking-form').addEventListener('submit', function(event) {
    event.preventDefault();
    currentPage = 1; // Reset to first page on new search
    fetchBookings();
});

document.getElementById('prevPage').addEventListener('click', function(event) {
    event.preventDefault();
    if (currentPage > 1) {
        currentPage--;
        fetchBookings();
    }
});

document.getElementById('nextPage').addEventListener('click', function(event) {
    event.preventDefault();
    currentPage++;
    fetchBookings();
});

function fetchBookings() {
    const bookedDateInput = document.getElementById('booked_date').value;
    const bookedDate = new Date(bookedDateInput);
    const formattedBookedDate = ('0' + bookedDate.getDate()).slice(-2) + '-' + 
                                ('0' + (bookedDate.getMonth() + 1)).slice(-2) + '-' + 
                                bookedDate.getFullYear(); 
    
    const data = {
        page: currentPage,
        pageSize: pageSize,
        booked_date: formattedBookedDate
    };
    
    fetch('http://192.168.1.9:8010/booking/find', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + accessToken,
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        displayData(data); // Call function to display the data
    })
    .catch((error) => {
        console.error('Error:', error);
    });3
}

function displayData(response) {
    console.log('Displaying data:', response);
    const bookingDataDiv = document.getElementById('bookingData');
    bookingDataDiv.innerHTML = ''; // Clear existing data

    // Create table
    const table = document.createElement('table');
    table.classList.add('table', 'table-striped');

    // Create table header
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    const headers = ['Booking ID', 'Date', 'Time', 'Status', 'Barber ID', 'Created At', 'Updated At', 'Actions'];
    headers.forEach(headerText => {
        const th = document.createElement('th');
        th.textContent = headerText;
        headerRow.appendChild(th);
    });
    thead.appendChild(headerRow);
    table.appendChild(thead);

    // Create table body
    const tbody = document.createElement('tbody');
    response.data.forEach(booking => {
        const row = document.createElement('tr');

        const idCell = document.createElement('td');
        idCell.textContent = booking.id;
        row.appendChild(idCell);

        const dateCell = document.createElement('td');
        dateCell.textContent = booking.booked_date;
        row.appendChild(dateCell);

        const timeCell = document.createElement('td');
        timeCell.textContent = booking.time_slot.start_time;
        row.appendChild(timeCell);

        const statusCell = document.createElement('td');
        statusCell.textContent = booking.status;
        row.appendChild(statusCell);

        const barberIdCell = document.createElement('td');
        barberIdCell.textContent = booking.barber_id;
        row.appendChild(barberIdCell);

        const createdAtCell = document.createElement('td');
        createdAtCell.textContent = new Date(booking.created_at * 1000).toLocaleString();
        row.appendChild(createdAtCell);

        const updatedAtCell = document.createElement('td');
        updatedAtCell.textContent = new Date(booking.updated_at * 1000).toLocaleString();
        row.appendChild(updatedAtCell);

        const actionsCell = document.createElement('td');
        const viewButton = document.createElement('button');
        viewButton.textContent = 'View';
        viewButton.classList.add('btn', 'btn-primary', 'btn-sm');
        viewButton.addEventListener('click', () => {
            // Handle view action
            console.log('View booking:', booking.id);
        });

        const updateButton = document.createElement('button');
        updateButton.textContent = 'Update';
        updateButton.classList.add('btn', 'btn-secondary', 'btn-sm');
        updateButton.addEventListener('click', () => {
            // Handle update action
            console.log('Update booking:', booking.id);
        });

        actionsCell.appendChild(viewButton);
        actionsCell.appendChild(updateButton);
        row.appendChild(actionsCell);

        tbody.appendChild(row);
    });
    table.appendChild(tbody);

    // Append table to bookingDataDiv
    bookingDataDiv.appendChild(table);

    bookingDataDiv.appendChild(table);


    // Calculate and display pagination info
    const totalItems = response.totalItems; // Assuming the API response includes the total number of items
    const totalPages = Math.ceil(totalItems / pageSize);
    document.getElementById('pageInfo').textContent = `Page ${currentPage} of ${totalPages}`;

    // Handle pagination visibility
    document.getElementById('prevPage').parentElement.style.display = currentPage > 1 ? 'block' : 'none';
    document.getElementById('nextPage').parentElement.style.display = currentPage < totalPages ? 'block' : 'none';
}