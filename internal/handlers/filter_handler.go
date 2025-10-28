package handlers

import (
	"fmt"
	"forum/internal/domain"
	"log"
	"net/http"
)

/*type Datas struct {
	Topics     []domain.Topic
	Categories []domain.Category
	IsLoggedIn bool
}*/

func FilterTopicByUser(w http.ResponseWriter, r *http.Request) {

	isLoggedIn := false
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		isLoggedIn = false
	} else {
		isLoggedIn = true
	}

	if r.Method != http.MethodPost {
		http.Error(w, "❌ unauthorized method", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur de parsing du formulaire", http.StatusBadRequest)
		fmt.Println("err de parse")
		return
	}

	var topics []domain.Topic
	var categories []domain.Category
	var err error

	filter := r.FormValue("filters")
	FormCategory := r.FormValue("categories")
	fmt.Println(FormCategory)

	fmt.Println(filter)

	if filter == "messages" && FormCategory == "" {
		topics, err = filterService.FilterTopic(int(user.ID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if topics == nil {
			fmt.Println("utilisateur sans topics")
		}
		for i := range topics {
			categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
			if err != nil {
				log.Println("❌ error fetching categories:", err)
				http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
				return
			}
		}
	} else if filter == "" && FormCategory != "" {
		topics, err = filterService.FilterByCategorie(FormCategory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if topics == nil {
			fmt.Println("aucun topic pour cette categories")
		}
		categories, err = categoryService.GetAllCategories()
		if err != nil {
			log.Println("❌ error fetching categories:", err)
			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
			return
		}
	} else if filter == "messages" && FormCategory != "" {
		topics, err = filterService.FilterByCategorieAndUserId(FormCategory, int(user.ID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if topics == nil {
			fmt.Println("aucun topic pour cette categories")
		}
		categories, err = categoryService.GetAllCategories()
		if err != nil {
			log.Println("❌ error fetching categories:", err)
			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
			return
		}
	} else if filter == "likes" && FormCategory == "" {
		topics, err = filterService.GetLikedTopics(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if topics == nil {
			fmt.Println("utilisateur sans topics")
		}
		for i := range topics {
			categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
			if err != nil {
				log.Println("❌ error fetching categories:", err)
				http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	for i := range topics {
		topics[i].Categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
		if err != nil {
			log.Println("❌ error fetching categories:", err)
			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
			return
		}
	}

	datas := Datas{
		Topics:     topics,
		IsLoggedIn: isLoggedIn,
		Categories: categories,
	}
	log.Printf("Nombre de topics: %d\n", len(topics))
	RenderTemplate(w, "home.html", datas)
}

// package handlers

// import (
// 	"fmt"
// 	"forum/internal/domain"
// 	"log"
// 	"net/http"
// )

// /*type Datas struct {
// 	Topics     []domain.Topic
// 	Categories []domain.Category
// 	IsLoggedIn bool
// }*/

// func FilterTopicByUser(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "❌ unauthorized method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Erreur de parsing du formulaire", http.StatusBadRequest)
// 		fmt.Println("err de parse")
// 		return
// 	}
// 	var user *domain.User
// 	isLoggedIn := false
// 	var topics []domain.Topic
// 	var categories []domain.Category
// 	cookie, err := r.Cookie("session_token")
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		user, _ = userService.Home(cookie.Value)
// 		if user != nil {
// 			isLoggedIn = true
// 		}
// 	}

// 	filter := r.FormValue("filters")
// 	FormCategory := r.FormValue("categories")
// 	fmt.Println(FormCategory)

// 	fmt.Println(filter)

// 	if filter == "messages" && FormCategory == "" {
// 		topics, err = filterService.FilterTopic(int(user.ID))
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		if topics == nil {
// 			fmt.Println("utilisateur sans topics")
// 		}
// 		for i := range topics {
// 			categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
// 			if err != nil {
// 				log.Println("❌ error fetching categories:", err)
// 				http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	} else if filter == "" && FormCategory != "" {
// 		topics, err = filterService.FilterByCategorie(FormCategory)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		if topics == nil {
// 			fmt.Println("aucun topic pour cette categories")
// 		}
// 		categories, err = categoryService.GetAllCategories()
// 		if err != nil {
// 			log.Println("❌ error fetching categories:", err)
// 			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 			return
// 		}
// 	} else if filter == "messages" && FormCategory != "" {
// 		topics, err = filterService.FilterByCategorieAndUserId(FormCategory, int(user.ID))
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		if topics == nil {
// 			fmt.Println("aucun topic pour cette categories")
// 		}
// 		categories, err = categoryService.GetAllCategories()
// 		if err != nil {
// 			log.Println("❌ error fetching categories:", err)
// 			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 			return
// 		}
// 	} else if filter == "likes" && FormCategory == "" {
// 		topics, err = filterService.GetLikedTopics(user.ID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		if topics == nil {
// 			fmt.Println("utilisateur sans topics")
// 		}
// 		for i := range topics {
// 			categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
// 			if err != nil {
// 				log.Println("❌ error fetching categories:", err)
// 				http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	} else {
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	}
// 	for i := range topics {
// 		topics[i].Categories, err = categoryService.GetCategoriesByTopicID(topics[i].ID)
// 		if err != nil {
// 			log.Println("❌ error fetching categories:", err)
// 			http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// 	/*categories, err = categoryService.GetAllCategories()
// 	if err != nil {
// 		log.Println("❌ error fetching categories:", err)
// 		http.Error(w, "❌ error fetching categories", http.StatusInternalServerError)
// 		return
// 	}*/

// 	datas := Datas{
// 		Topics:     topics,
// 		IsLoggedIn: isLoggedIn,
// 		Categories: categories,
// 	}
// 	log.Printf("Nombre de topics: %d\n", len(topics))
// 	RenderTemplate(w, "home.html", datas)
// }
