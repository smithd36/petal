package handlers

import (
    "log"
    "net/http"
    "html/template"
    "github.com/smithd36/petal/models"
    "golang.org/x/crypto/bcrypt"
)

type AuthPageData struct {
    Title string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
        if err != nil {
            log.Printf("Error parsing templates: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        data := AuthPageData{Title: "Register"}
        err = tmpl.ExecuteTemplate(w, "base.html", data)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    } else if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            log.Printf("Error hashing password: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        _, err = models.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
        if err != nil {
            log.Printf("Error inserting user: %v", err)
            http.Error(w, "Unable to register user", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
        if err != nil {
            log.Printf("Error parsing templates: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        data := AuthPageData{Title: "Login"}
        err = tmpl.ExecuteTemplate(w, "base.html", data)
        if err != nil {
            log.Printf("Error executing template: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    } else if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        var hashedPassword string
        err := models.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
        if err != nil {
            log.Printf("Error fetching user: %v", err)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
        if err != nil {
            log.Printf("Invalid credentials: %v", err)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    }
}
