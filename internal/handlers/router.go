package handlers

import (
	"forum/internal/domain"
	"forum/internal/middleware"
	"net/http"
)

func Router(userService domain.UserService, topicPostService domain.TopicPostService, categoryService domain.CategoryService, reactionService domain.ReactionService, filterService domain.FilterService, authService domain.AuthService, middleware middleware.AuthMiddleware) http.Handler {

	InitHandlers(userService, topicPostService, categoryService, reactionService, filterService, authService, middleware)

	mux := http.NewServeMux()

	// Routes:
	mux.Handle("/", middleware.Handle(http.HandlerFunc(HomeHandler)))
	mux.Handle("/thread", middleware.Handle(http.HandlerFunc(ThreadHandler)))
	mux.HandleFunc("/login", AuthenticateHandler)
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/register", RegisterHandler)
	mux.Handle("/create-topic", middleware.Handle(http.HandlerFunc(CreateTopicHandler)))
	mux.HandleFunc("/topic", TopicHandler)
	mux.Handle("/add-post", middleware.Handle(http.HandlerFunc(AddPostHandler)))
	mux.Handle("/react", middleware.Handle(http.HandlerFunc(ReactHandler)))
	mux.Handle("/remove-reaction", middleware.Handle(http.HandlerFunc(RemoveReactionHandler)))
	mux.Handle("/filter", middleware.Handle(http.HandlerFunc(FilterTopicByUser)))

	mux.HandleFunc("/api/login-gh/", githubLoginHandler)
	mux.HandleFunc("/api/github/callback/", githubCallbackHandler)

	mux.HandleFunc("/api/login-google/", GoogleLoginHandler)
	mux.HandleFunc("/api/google/callback/", GoogleCallbackHandler)

	fs := http.FileServer(http.Dir("internal/templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	return mux
}
