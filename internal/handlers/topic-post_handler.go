package handlers

import (
	"fmt"
	"forum/internal/domain"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

/*MARK: CreateTopic
 */
func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "❌ unauthenticated user", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "❌ unauthorized method", http.StatusMethodNotAllowed)
		return
	}
	//.Println("User ID :", user.ID)
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories_string := r.Form["categories"]
	var categories_id []int
	for _, cat := range categories_string {
		id, err := strconv.Atoi(cat)
		if err != nil {
			http.Error(w, "error in id recovery", http.StatusBadRequest)
		}
		categories_id = append(categories_id, id)
	}
	//log.Println("Reçu :", title, content)

	if title == "" || content == "" {
		http.Error(w, "❌ missing required fields", http.StatusBadRequest)
		return
	}

	err := topicPostService.CreateTopic(title, content, int(user.ID), categories_id)
	if err != nil {
		http.Error(w, "❌ error inserting topic"+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// affichage d'un sujet et ses messages
/*MARK: Topic+Posts
 */
func TopicHandler(w http.ResponseWriter, r *http.Request) {
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

	thread, err := topicPostService.GetThreadByID(id)
	if err != nil {
		http.Error(w, "❌ topic not found", http.StatusNotFound)
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

	tmpl := template.Must(template.ParseFiles("internal/templates/topic.html"))
	tmpl.Execute(w, thread)
}

/*MARK: AddPost
 */
func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "❌ unauthenticated user", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "❌ unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	// parser le formulaire
	if err := r.ParseForm(); err != nil {
		http.Error(w, "❌ cannot parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	/*log.Println("Form values:", r.Form)*/

	content := r.FormValue("content")
	topicID, _ := strconv.Atoi(r.FormValue("topic_id"))

	//log.Println("User ID : ", user.ID)
	//userID := r.Context().Value("userID").(int)

	err := topicPostService.AddPost(topicID, content, int(user.ID))
	if err != nil {
		log.Println("❌ AddPost error:", err)
		http.Error(w, "❌ error inserting post:"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("✅ Post inserted, redirecting to thread")
	http.Redirect(w, r, fmt.Sprintf("/thread?id=%d", topicID), http.StatusSeeOther)
}
