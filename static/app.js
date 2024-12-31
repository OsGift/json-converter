//app.js
document.addEventListener('DOMContentLoaded', function () {
    const button = document.getElementById("convertBtn");
    
    // Add click event listener to the button
    button.addEventListener("click", convertJSON);
});

function convertJSON() {
    const jsonInput = document.getElementById("jsonInput").value;

    // Validate if the input is a valid JSON
    try {
        const parsedJSON = JSON.parse(jsonInput);

        // Send the parsed JSON data to the backend for conversion
        fetch('/convert', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(parsedJSON),
        })
        .then(response => response.json())
        .then(data => {
            // Update the Go struct output
            document.getElementById("goStructOutput").textContent = data.structCode;
        })
        .catch(error => {
            console.error('Error during fetch:', error);
            alert('Error processing the request.');
        });

    } catch (error) {
        alert('Invalid JSON. Please check your input.');
        console.error('Invalid JSON input:', error);
    }
}
