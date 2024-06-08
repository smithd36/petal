package handlers

import (
    "log"
    "net/http"
    "html/template"
    "github.com/smithd36/petal/models"
)

type DashboardData struct {
    Title  string
    Plants []models.Plant
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/base.html", "templates/dashboard.html")
    if err != nil {
        log.Printf("Error parsing templates: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Fetch plants for the user (hardcoded user ID for now)
    rows, err := models.DB.Query("SELECT id, plant_name, trefle_id FROM plants WHERE user_id = ?", 1)
    if err != nil {
        log.Printf("Error fetching plants: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    plants := []models.Plant{}
    for rows.Next() {
        var plant models.Plant
        err := rows.Scan(&plant.ID, &plant.PlantName, &plant.TrefleID)
        if err != nil {
            log.Printf("Error scanning plant: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        plants = append(plants, plant)
    }

    data := DashboardData{
        Title:  "Dashboard",
        Plants: plants,
    }

    err = tmpl.ExecuteTemplate(w, "base.html", data)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
