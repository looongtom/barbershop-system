function displayImage(event, previewId) {
    const input = event.target;
    const file = input.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const img = document.getElementById(previewId);
            img.src = e.target.result;
            img.style.display = 'block';
            img.style.width = '300px'; // Set the desired width
            img.style.height = 'auto'; // Maintain aspect ratio
        }
        reader.readAsDataURL(file);
    }
}

async function getHaircut() {
    const selfImg = document.getElementById('image1').files[0];
    const shapeImg = document.getElementById('image2').files[0];
    const colorImg = document.getElementById('image3').files[0];
    const accountId = "4"; // Replace with the actual account ID if needed

    if (!selfImg) {
        alert('Please upload the source photo.');
        return;
    }

    const formData = new FormData();
    formData.append('self_img', selfImg);
    if (shapeImg) formData.append('shape_img', shapeImg);
    if (colorImg) formData.append('color_img', colorImg);
    formData.append('account_id', accountId);

    try {
        const response = await fetch('http://192.168.1.9:8005/previewimage/upload', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' // Add the actual token if needed
            },
            body: formData
        });

        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }

        const result = await response.json();
        console.log(result);
        alert('Hairstyle generated successfully!');
        // Process the result as needed
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
        alert('Failed to generate the hairstyle.');
    }
}