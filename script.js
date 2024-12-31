document.getElementById('convert-btn').addEventListener('click', function() {
    const jsonInput = document.getElementById('json-input').value;
    const outputElement = document.getElementById('go-struct-output');
    
    try {
        const jsonObject = JSON.parse(jsonInput);
        const goStruct = convertToGoStruct(jsonObject);
        outputElement.textContent = goStruct;
    } catch (error) {
        outputElement.textContent = 'Invalid JSON format';
    }
});

function convertToGoStruct(jsonObj) {
    const structName = "GeneratedStruct";  // You can customize this if needed
    let goStruct = `type ${structName} struct {\n`;
    
    for (const [key, value] of Object.entries(jsonObj)) {
        const goType = getGoType(value);
        const formattedKey = toGoNamingConvention(key);
        goStruct += `\t${formattedKey} ${goType} \`json:"${key}"\`\n`;
    }
    
    goStruct += '}\n';
    return goStruct;
}

function getGoType(value) {
    if (value === null) return "interface{}";
    switch (typeof value) {
        case 'string':
            return 'string';
        case 'number':
            return Number.isInteger(value) ? 'int' : 'float64';
        case 'boolean':
            return 'bool';
        case 'object':
            return Array.isArray(value) ? '[]interface{}' : 'struct{}'; // Can be further expanded for nested structs
        default:
            return 'interface{}';
    }
}

function toGoNamingConvention(key) {
    return key
        .replace(/_./g, match => match[1].toUpperCase()) // Convert snake_case to CamelCase
        .replace(/^./, match => match.toUpperCase()); // Capitalize the first letter
}
