package handlers

import (
	"forum/internal/domain"
	"log"
	"net/http"
)

type Datas struct {
	Topics     []domain.Topic
	Categories []domain.Category
	IsLoggedIn bool
	Error      string
	Email      string
}

/*var tpl = template.Must(template.ParseFiles("internal/templates/home.html"))*/

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := false
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		isLoggedIn = false
	} else {
		isLoggedIn = true
	}

	if r.URL.Path != "/" {
		http.Error(w, "❌ not found", http.StatusNotFound)
		return
	}

	topics, err := topicPostService.GetAllTopics()
	if err != nil {
		log.Println("❌ error fetching topics:", err)
		http.Error(w, "❌ error fetching topics", http.StatusInternalServerError)
		return
	}
	for i := range topics {
		topics[i].Categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
		if err != nil {
			log.Println("❌ error fetching categories:", err)
			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
			return
		}
		likes, dislikes, _ := reactionService.GetReactionCounts("topics", int64(topics[i].ID))
		topics[i].Likes = likes
		topics[i].Dislikes = dislikes
	}
	categories, err := categoryService.GetAllCategories()
	if err != nil {
		log.Println("❌ error fetching categories:", err)
		http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
		return
	}

	datas := Datas{
		Topics:     topics,
		IsLoggedIn: isLoggedIn,
		Categories: categories,
	}
	log.Printf("Nombre de topics: %d\n", len(topics))
	RenderTemplate(w, "home.html", datas)
}
