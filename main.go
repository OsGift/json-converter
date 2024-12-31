package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"text/template"
)

// Function to generate Go struct from JSON data
func generateGoStruct(jsonData interface{}, parentKey string, structs map[string]string) (string, error) {
	var result string
	val := reflect.ValueOf(jsonData)

	// Log the type of the input data
	log.Printf("Generating struct for: %v", reflect.TypeOf(jsonData))

	// If the value is a map (i.e., a JSON object)
	if val.Kind() == reflect.Map {
		result += "type " + parentKey + " struct {\n"

		// Loop over the keys and types in the JSON object
		for _, key := range val.MapKeys() {
			fieldName := key.String()
			fieldValue := val.MapIndex(key)

			// Capitalize the first letter to make it a valid Go field
			fieldName = strings.Title(fieldName)

			// Get the Go type of the field value
			fieldType := reflect.TypeOf(fieldValue.Interface()).String()

			log.Printf("Processing field: %s with type: %s", fieldName, fieldType)

			// Handle nested objects (maps)
			if fieldType == "map[string]interface {}" { // This is a nested object (map)
				goType := fieldName
				if _, exists := structs[goType]; !exists {
					// Recursively generate the nested struct and add it to the map
					structCode, err := generateGoStruct(fieldValue.Interface(), goType, structs)
					if err != nil {
						return "", err // Return error if recursion fails
					}
					structs[goType] = structCode
				}

				// Reference the nested struct in the parent struct
				result += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, goType, key.String())
			} else if fieldType == "[]interface {}" { // This is an array (slice)
				// If the field is an array, handle the type of its elements
				if fieldValue.Len() > 0 {
					elem := fieldValue.Index(0).Interface() // Take the first element of the array to inspect its type
					elemType := reflect.TypeOf(elem).String()

					// If the element is an object (map), recursively generate a struct for it
					if elemType == "map[string]interface {}" {
						goType := fieldName + "Item" // Unique name for array element type
						if _, exists := structs[goType]; !exists {
							structCode, err := generateGoStruct(elem, goType, structs)
							if err != nil {
								return "", err // Return error if recursion fails
							}
							structs[goType] = structCode
						}

						// Reference the array of objects
						result += fmt.Sprintf("\t%s []%s `json:\"%s\"`\n", fieldName, goType, key.String())
					} else {
						// If it's a primitive array (e.g., []string, []int)
						result += fmt.Sprintf("\t%s []%s `json:\"%s\"`\n", fieldName, elemType, key.String())
					}
				} else {
					// Handle empty arrays gracefully
					result += fmt.Sprintf("\t%s []interface{} `json:\"%s\"`\n", fieldName, key.String())
				}
			} else if fieldType == "float64" { // Handle numbers (float64 is used by JSON library)
				result += fmt.Sprintf("\t%s int `json:\"%s\"`\n", fieldName, key.String())
			} else if fieldType == "bool" { // Handle booleans
				result += fmt.Sprintf("\t%s bool `json:\"%s\"`\n", fieldName, key.String())
			} else { // Default to string for other types (e.g., string, int, etc.)
				result += fmt.Sprintf("\t%s string `json:\"%s\"`\n", fieldName, key.String())
			}
		}

		result += "}\n"
	}

	return result, nil
}

// Handler to render the HTML page
func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Handler to process the JSON and return the Go struct
func convertJSONToStruct(w http.ResponseWriter, r *http.Request) {
	var inputData interface{}

	// Parse JSON from the request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputData); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Log the incoming JSON for debugging purposes
	log.Printf("Received JSON: %v", inputData)

	// Initialize a map to store the generated structs
	structs := make(map[string]string)

	// Generate the Go struct for the root object and any nested structs
	structCode, err := generateGoStruct(inputData, "Data", structs)
	if err != nil {
		log.Printf("Error generating Go struct: %v", err)
		http.Error(w, fmt.Sprintf("Error processing the request: %s", err), http.StatusInternalServerError)
		return
	}

	// Add all nested structs to the output, in the correct order
	for _, structDef := range structs {
		structCode = structDef + "\n" + structCode
	}

	// Send back the response in JSON format
	response := map[string]interface{}{
		"structCode": structCode,
		"data":       inputData,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
	}
}

func main() {
	// Serve static files (HTML, JS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Route for main page
	http.HandleFunc("/", serveIndex)

	// Route to handle JSON conversion
	http.HandleFunc("/convert", convertJSONToStruct)

	// Start the server
	port := "8080"
	fmt.Printf("Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
