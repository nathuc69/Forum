package handlers

import (
	"forum/internal/domain"
	"html/template"
	"net/http"
	"strconv"
)

type ThreadData struct {
	Topic      domain.Topic
	Posts      []domain.Post
	Categories []domain.Category
	IsLoggedIn bool
}

func ThreadHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := false
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		isLoggedIn = false
	} else {
		isLoggedIn = true
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "❌ missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "❌ invalid ID", http.StatusBadRequest)
		return
	}

	// Récupération du thread (topic + posts)
	thread, err := topicPostService.GetThreadByID(id)
	if err != nil {
		http.Error(w, "❌ topic not found", http.StatusNotFound)
		return
	}

	thread.Categories, err = categoryService.GetCategoriesByTopicID(id)
	if err != nil {
		http.Error(w, "❌ categories not found", http.StatusNotFound)
		return
	}

	// Injection des likes/dislikes pour le topic
	likes, dislikes, _ := reactionService.GetReactionCounts("topics", int64(thread.Topic.ID))
	thread.Topic.Likes = likes
	thread.Topic.Dislikes = dislikes

	// Injection des likes/dislikes pour chaque post
	for i := range thread.Posts {
		plikes, pdislikes, _ := reactionService.GetReactionCounts("posts", int64(thread.Posts[i].ID))
		thread.Posts[i].Likes = plikes
		thread.Posts[i].Dislikes = pdislikes
	}

	thread.IsLoggedIn = isLoggedIn
	// Rendu du template thread.html
	tmpl := template.Must(template.ParseFiles("internal/templates/thread.html"))
	err = tmpl.Execute(w, thread)
	if err != nil {
		http.Error(w, "❌ template error: "+err.Error(), http.StatusInternalServerError)
	}
}
