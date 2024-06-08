package main

import (
    "log"
    "net/http"
    "os"
    "html/template"

    "github.com/joho/godotenv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/smithd36/petal/handlers"
    "github.com/smithd36/petal/models"
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

    // Retrieve the API key from environment variables
    apiKey := os.Getenv("TREFLE_API_KEY")
    if apiKey == "" {
        log.Fatal("TREFLE_API_KEY is not set in the environment")
    }

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

    r.Get("/register", handlers.RegisterHandler)
    r.Post("/register", handlers.RegisterHandler)
    r.Get("/login", handlers.LoginHandler)
    r.Post("/login", handlers.LoginHandler)
    r.Get("/dashboard", handlers.DashboardHandler)

    r.Get("/roots", handlers.ListRootsHandler)
    r.Get("/roots/new", handlers.CreateRootHandler)
    r.Post("/roots/new", handlers.CreateRootHandler)
    r.Get("/roots/{rootID}", handlers.ViewRootHandler)
    r.Post("/roots/{rootID}/comments", handlers.CreateCommentHandler)

    // Serve static files
    fileServer := http.FileServer(http.Dir("./static"))
    r.Handle("/static/*", http.StripPrefix("/static", fileServer))

    log.Println("Starting server on :8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}