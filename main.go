package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/smithd36/petal/handlers"
	authMiddleware "github.com/smithd36/petal/middleware"
	"github.com/smithd36/petal/models"
	"github.com/smithd36/petal/utils"
)

type PageData struct {
	Title string
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	// Initialize the JWT key
	utils.InitializeJWTKey()

	models.InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
		if err != nil {
			log.Printf("Error parsing templates: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		data := PageData{Title: "Welcome"}
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})

	// About route GET
	r.Get("/about", handlers.AboutHandler)

	// Register route GET
	r.Get("/register", handlers.RegisterHandler)

	// Register route POST
	r.Post("/register", handlers.RegisterHandler)

	// Login route GET
	r.Get("/login", handlers.LoginHandler)

	// Login route POST
	r.Post("/login", handlers.LoginHandler)

	// View a single root publicly viewable
	r.Get("/roots/{rootID}", handlers.ViewRootHandler)

	// Roots should be publicly viewable
	r.Get("/roots", handlers.ListRootsHandler)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.JWTAuth)

		// For the dashboard
		r.Get("/dashboard", handlers.DashboardHandler)

		// For uploading images
		r.Post("/upload_image", handlers.UploadImageHandler)

		// For creating roots
		r.Get("/roots/new", handlers.CreateRootHandler)

		// For creating roots
		r.Post("/roots/new", handlers.CreateRootHandler)

		// For adding comments
		r.Post("/roots/{rootID}/comments", handlers.AddCommentHandler)
	})

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
