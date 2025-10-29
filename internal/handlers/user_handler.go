package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"forum/internal/domain"
	"net/http"
	"time"
)

/*MARK: Logout
 */
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "home.html", nil)
		return
	}

	// Récupération des cookies
	cookie, err := r.Cookie("session_token")
	cookieStateGoogle, errStateGoogle := r.Cookie("state-Google")
	cookieStateGithub, errStateGithub := r.Cookie("state-Github")

	// Gestion du logout session standard
	if err == nil && cookie != nil && cookie.Value != "" {
		if err := userService.Logout(cookie.Value); err != nil {
			fmt.Println("Error logging out session:", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})
	}

	// Gestion du logout Google
	if errStateGoogle == nil && cookieStateGoogle != nil && cookieStateGoogle.Value != "" {
		if err := userService.Logout(cookieStateGoogle.Value); err != nil {
			fmt.Println("Error logging out Google:", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "state-Google",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})
	}

	// Gestion du logout Github
	if errStateGithub == nil && cookieStateGithub != nil && cookieStateGithub.Value != "" {
		if err := userService.Logout(cookieStateGithub.Value); err != nil {
			fmt.Println("Error logging out Github:", err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "state-Github",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*MARK: Authenticate
 */
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Affiche le formulaire login via home.html
		RenderTemplate(w, "home.html", nil)
		return
	}
	// Méthode POST : récupération du formulaire
	if err := r.ParseForm(); err != nil {
		http.Error(w, "❌ invalid form", http.StatusBadRequest)
		fmt.Println("❌ error processing form")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	_, err := userService.Authenticate(email, password)
	if err != nil {
		topics, _ := topicPostService.GetAllTopics()
		data := Datas{
			Error:      err.Error(),
			Email:      email,
			Topics:     topics,
			IsLoggedIn: false,
		}
		RenderTemplate(w, "home.html", data)
		return
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println("❌ error generating token ")
	}

	sessionToken := base64.URLEncoding.EncodeToString(bytes)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})

	userService.TokenLogIn(sessionToken, email)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*MARK: Register
 */
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "register.html", nil)
		return
	}
	// Méthode POST : récupération du formulaire
	if err := r.ParseForm(); err != nil {
		http.Error(w, "❌ error invalid form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := userService.Register(username, email, password)
	if err != nil {
		if err.Error() == "❌ Email déjà enregistré" {

			registerUser := domain.User{
				Error:    err.Error(),
				Username: username,
			}

			RenderTemplate(w, "register.html", registerUser)

		} else if err.Error() == "❌ Nom d'utilisateur déjà utilisé" {
			registerUser := domain.User{
				Error: "❌ Nom d'utilisateur déjà utilisé",
				Email: email,
			}

			RenderTemplate(w, "register.html", registerUser)
		} else {
			registerUser := domain.User{
				Error: "❌ error",
			}

			RenderTemplate(w, "register.html", registerUser)

		}
		return
	}

	// Enregistrement réussi, redirection vers la page d’accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

/*MARK: Session
 */
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
