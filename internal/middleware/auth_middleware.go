package middleware

import (
	"context"
	"forum/internal/domain"
	"net/http"
)

type AuthMiddleware struct {
	userService domain.UserService
}

func NewAuthMiddleware(us domain.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: us}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var user *domain.User
		var err error

		// 1️⃣ — On essaie le cookie de session standard
		cookie, err := r.Cookie("session_token")
		if err == nil && cookie.Value != "" {
			user, err = m.userService.Home(cookie.Value)
		}

		// 2️⃣ — Si session_token absent ou invalide, on tente avec le cookie "state"
		if err != nil || user == nil {
			cookieState, errState := r.Cookie("state")
			if errState == nil && cookieState.Value != "" {
				user, err = m.userService.Home(cookieState.Value)
			}
		}
		//fmt.Printf("🧠 Middleware → user.ID=%d | user.Username=%s | token=%s\n", user.ID, user.Username, cookieState.Value)

		// 4️⃣ — Ajout du user dans le contexte
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
