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

        // Generate Go struct code
        const structCode = generateGoStruct(parsedJSON, "Data");

        // Update the Go struct output
        document.getElementById("goStructOutput").textContent = structCode;
    } catch (error) {
        alert('Invalid JSON. Please check your input.');
        console.error('Invalid JSON input:', error);
    }
}

function generateGoStruct(jsonData, parentKey) {
    let result = '';

    if (parentKey) {
        result += `// ${parentKey}\n`;
    }

    if (Array.isArray(jsonData)) {
        // Handle array of objects (e.g., the "data" array in your example)
        result += '[]struct {\n';
        jsonData.forEach(item => {
            result += generateGoStruct(item, parentKey); // Recursively generate structs
        });
        result += '}{}\n';
    } else if (typeof jsonData === 'object') {
        // Handle objects (maps)
        result += `type ${parentKey} struct {\n`;
        for (const key in jsonData) {
            if (jsonData.hasOwnProperty(key)) {
                const fieldName = capitalizeFirstLetter(key);
                const fieldValue = jsonData[key];
                let fieldType = getGoType(fieldValue);

                // Recursively handle nested objects or arrays
                if (fieldType === 'struct') {
                    result += `\t${fieldName} ${capitalizeFirstLetter(key)} \`json:"${key}"\`\n`;
                } else {
                    result += `\t${fieldName} ${fieldType} \`json:"${key}"\`\n`;
                }
            }
        }
        result += '}\n';
    } else {
        // Fallback if the value is a primitive type (string, number, etc.)
        result += `type ${parentKey} string\n`; // Default to string type
    }

    return result;
}

function getGoType(value) {
    if (Array.isArray(value)) {
        return '[]struct';
    } else if (typeof value === 'object') {
        return 'struct';
    } else if (typeof value === 'number') {
        return 'int';
    } else if (typeof value === 'boolean') {
        return 'bool';
    }
    return 'string'; // Default type
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}
