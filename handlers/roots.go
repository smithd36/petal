package handlers

import (
    "html/template"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/smithd36/petal/models"
)

type RootsPageData struct {
    Title    string
    Roots    []models.Root
    Root     models.Root
    Comments []models.Comment
}

func ListRootsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := models.DB.Query("SELECT id, user_id, title, content, created_at FROM roots")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var roots []models.Root
    for rows.Next() {
        var root models.Root
        if err := rows.Scan(&root.ID, &root.UserID, &root.Title, &root.Content, &root.CreatedAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        roots = append(roots, root)
    }

    tmpl, err := template.ParseFiles("templates/base.html", "templates/roots.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data := RootsPageData{
        Title: "Forum Roots",
        Roots: roots,
    }
    if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func CreateRootHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        userID := 1 // TODO: Replace with actual logged-in user ID
        title := r.FormValue("title")
        content := r.FormValue("content")

        _, err := models.DB.Exec("INSERT INTO roots (user_id, title, content, created_at) VALUES (?, ?, ?, ?)", userID, title, content, time.Now())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/roots", http.StatusSeeOther)
        return
    }

    tmpl, err := template.ParseFiles("templates/base.html", "templates/create_root.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data := RootsPageData{Title: "Create Root"}
    if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ViewRootHandler(w http.ResponseWriter, r *http.Request) {
    rootID := chi.URLParam(r, "rootID")

    var root models.Root
    err := models.DB.QueryRow("SELECT id, user_id, title, content, created_at FROM roots WHERE id = ?", rootID).Scan(&root.ID, &root.UserID, &root.Title, &root.Content, &root.CreatedAt)
    if err != nil {
        http.Error(w, "Root not found", http.StatusNotFound)
        return
    }

    rows, err := models.DB.Query("SELECT id, root_id, user_id, content, created_at FROM comments WHERE root_id = ?", rootID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.ID, &comment.RootID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        comments = append(comments, comment)
    }

    tmpl, err := template.ParseFiles("templates/base.html", "templates/root.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data := RootsPageData{
        Title:    root.Title,
        Root:     root,
        Comments: comments,
    }

    if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        rootID := chi.URLParam(r, "rootID")
        userID := 1 // TODO: Replace with actual logged-in user ID
        content := r.FormValue("content")

        _, err := models.DB.Exec("INSERT INTO comments (root_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", rootID, userID, content, time.Now())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/roots/"+rootID, http.StatusSeeOther)
        return
    }
}
