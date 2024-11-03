// account.js
const username = localStorage.getItem('username');
if (username) {
    // Display username in profile-bt and hide login-bt
    const profileButton = document.getElementById('profile-bt');
    const loginButton = document.getElementById('login-bt');
    profileButton.textContent = username;
    profileButton.style.display = 'block';
    loginButton.style.display = 'none';
    }

document.addEventListener("DOMContentLoaded", function() {
    document.addEventListener('click', function(event) {
        if (event.target && event.target.id === 'sign-in-button') {
            event.preventDefault(); // Prevent the default form submission

            const username = document.getElementById('user').value;
            const password = document.getElementById('pass').value;

            console.log('Username:', username);
            console.log('Password:', password);

            // Replace with your API endpoint
            const apiEndpoint = 'http://localhost:8008/auth/login';

            // Create the request payload
            const payload = {
                username: username,
                password: password
            };

            // Send the request to the API
            fetch(apiEndpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            })
            .then(response => response.json())
            .then(data => {
                // Handle the response data
                console.log('Success:', data);
                // You can add more logic here to handle successful login
                // save accesstoken in local storage
                localStorage.setItem('accessToken', data.accessToken);
                localStorage.setItem('refreshToken', data.refreshToken);
                localStorage.setItem('username', username);
                window.location.href = 'index.html';

            })
            .catch(error => {
                console.error('Error:', error);
                // Handle the error
                alert('Invalid username or password');
            });
        }

        if (event.target && event.target.id === 'logout-bt') {
            const token = localStorage.getItem('accessToken'); // Assuming the token is stored in localStorage
            fetch('http://localhost:8008/auth/logout', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(response => {
                if (response.ok) {
                    // Clear localStorage and update UI
                    localStorage.removeItem('username');
                    localStorage.removeItem('token');
                    document.getElementById('profile-bt').style.display = 'none';
                    document.getElementById('login-bt').style.display = 'block';
                    console.log('Logout successful');
                } else {
                    console.error('Logout failed');
                }
            })
            .catch(error => console.error('Error:', error));
        }
    });
});
