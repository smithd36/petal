package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/smithd36/petal/models"
)

type DashboardData struct {
	Title   string
	Plants  []models.Plant
	Images  []string
	Message string
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the user ID from context or session
	userID := r.Context().Value("userID")
	if userID == nil {
		log.Printf("User ID is nil")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID to int
	userIDInt, ok := userID.(int)
	if !ok {
		log.Printf("User ID is not of type int")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch images from the user's directory
	userDir := filepath.Join("static", "users", strconv.Itoa(userIDInt))
	var images []string
	var message string

	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		message = "No images uploaded"
	} else {
		entries, err := os.ReadDir(userDir)
		if err != nil {
			log.Printf("Error reading user directory: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				imagePath := filepath.Join("/static/users", strconv.Itoa(userIDInt), entry.Name())
				log.Printf("Image Path: %s", imagePath) // Debug: print the image path
				images = append(images, imagePath)
			}
		}

		if len(images) == 0 {
			message = "No images uploaded"
		}
	}

	data := DashboardData{
		Title:   "Dashboard",
		Images:  images,
		Message: message,
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UploadImageHandler handles the image upload process
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve the file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving an image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the user ID from context
	userID := r.Context().Value("userID")
	if userID == nil {
		log.Printf("User ID is nil")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID to int
	userIDInt, ok := userID.(int)
	if !ok {
		log.Printf("User ID is not of type int")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the user directory if it doesn't exist
	userDir := filepath.Join("static", "users", strconv.Itoa(userIDInt))
	err = os.MkdirAll(userDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating user directory", http.StatusInternalServerError)
		return
	}

	// Create the file path
	filePath := filepath.Join(userDir, header.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the uploaded file to the created file on the filesystem
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	// Redirect to the dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

type TrefleResponse struct {
	Data []models.Plant `json:"data"`
}

func SearchPlantHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("TREFLE_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("https://trefle.io/api/v1/plants/search?token=%s&q=%s", apiKey, query)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch data from Trefle API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var trefleResp TrefleResponse
	if err := json.NewDecoder(resp.Body).Decode(&trefleResp); err != nil {
		http.Error(w, "Failed to decode Trefle API response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	for _, plant := range trefleResp.Data {
		fmt.Fprintf(w, `<button class="list-group-item list-group-item-action" onclick='selectPlant(%s)'>%s</button>`, plantToJSON(plant), plant.CommonName)
	}
}

func plantToJSON(plant models.Plant) string {
	plantJSON, _ := json.Marshal(plant)
	return fmt.Sprintf("'%s'", string(plantJSON))
}
