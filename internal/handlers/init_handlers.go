package handlers

import (
	"fmt"
	"forum/internal/domain"
	"forum/internal/middleware"
	"html/template"
	"net/http"
	"path/filepath"
)

var userService domain.UserService
var topicPostService domain.TopicPostService
var categoryService domain.CategoryService
var reactionService domain.ReactionService
var filterService domain.FilterService
var templates *template.Template
var authService domain.AuthService
var middlewareService middleware.AuthMiddleware

func InitHandlers(us domain.UserService, tps domain.TopicPostService, cs domain.CategoryService, rs domain.ReactionService, flt domain.FilterService, auth domain.AuthService, mw middleware.AuthMiddleware) {
	userService = us
	topicPostService = tps
	categoryService = cs
	reactionService = rs
	filterService = flt
	authService = auth
	middlewareService = mw

	// Précharger tous les templates une seule fois
	var err error
	templates, err = template.ParseGlob(filepath.Join("internal", "templates", "*.html"))
	if err != nil {
		panic("❌ error parsing templates: " + err.Error())
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		fmt.Printf("template error: %v\n", err)
		//http.Error(w, "❌ internal error: "+err.Error(), http.StatusInternalServerError)
		w.Write([]byte("❌ internal error: " + err.Error()))
	}
}
