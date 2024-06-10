package handlers

import (
	"html/template"
	"net/http"
	"strconv"
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
	rows, err := models.DB.Query(`
		SELECT roots.id, roots.user_id, users.username, roots.title, roots.content, roots.created_at
		FROM roots
		JOIN users ON roots.user_id = users.id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roots []models.Root
	for rows.Next() {
		var root models.Root
		if err := rows.Scan(&root.ID, &root.UserID, &root.Username, &root.Title, &root.Content, &root.CreatedAt); err != nil {
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
		userID := r.Context().Value("userID").(int)
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
	err := models.DB.QueryRow("SELECT roots.id, roots.user_id, users.username, roots.title, roots.content, roots.created_at FROM roots JOIN users ON roots.user_id = users.id WHERE roots.id = ?", rootID).Scan(&root.ID, &root.UserID, &root.Username, &root.Title, &root.Content, &root.CreatedAt)
	if err != nil {
		http.Error(w, "Root not found", http.StatusNotFound)
		return
	}

	rows, err := models.DB.Query("SELECT comments.id, comments.root_id, comments.user_id, users.username, comments.content, comments.created_at FROM comments JOIN users ON comments.user_id = users.id WHERE comments.root_id = ?", rootID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.RootID, &comment.UserID, &comment.Username, &comment.Content, &comment.CreatedAt); err != nil {
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

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rootIDStr := chi.URLParam(r, "rootID")
		rootID, err := strconv.Atoi(rootIDStr)
		if err != nil {
			http.Error(w, "Invalid root ID", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Content is required", http.StatusBadRequest)
			return
		}

		// Insert the comment into the database
		_, err = models.DB.Exec("INSERT INTO comments (root_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", rootID, userID, content, time.Now())
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}

		// Redirect to the root view page
		http.Redirect(w, r, "/roots/"+rootIDStr, http.StatusSeeOther)
	}
}
