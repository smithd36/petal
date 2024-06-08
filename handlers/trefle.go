package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)

type Plant struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// FetchPlants retrieves plants from the Trefle API
func FetchPlants(w http.ResponseWriter, r *http.Request) {
    apiKey := os.Getenv("TREFLE_API_KEY")
    if apiKey == "" {
        http.Error(w, "API key not set", http.StatusInternalServerError)
        return
    }

    resp, err := http.Get(fmt.Sprintf("https://trefle.io/api/v1/plants?token=%s", apiKey))
    if err != nil {
        log.Printf("Error fetching plants: %v", err)
        http.Error(w, "Failed to fetch plants", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error response from Trefle API: %v", resp.Status)
        http.Error(w, "Failed to fetch plants", http.StatusInternalServerError)
        return
    }

    var result struct {
        Data []Plant `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("Error decoding response: %v", err)
        http.Error(w, "Failed to decode response", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result.Data)
}
