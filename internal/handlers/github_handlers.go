package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/internal/domain"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {

	GithubClientID := os.Getenv("GITHUB_CLIENT_ID")

	//state, err := randString(16)
	state, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	c := &http.Cookie{
		Name:     "state-Github",
		Value:    state.String(),
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
	log.Println(state)
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&state=%s", GithubClientID, state)
	http.Redirect(w, r, redirectURL, 301)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	GithubClientID := os.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	// Step 2: Users are redirected back to your site by GitHub
	//
	// The user is authenticated w/ GitHub by this point, and GH provides us
	// a temporary code we can exchange for an access token using the app's
	// full credentials.
	//
	// Start by checking the state returned by GitHub matches what
	// we've stored in the cookie.
	state, err := r.Cookie("state-Github")
	if err != nil {
		// Si pas de cookie state, renvoyer vers la page de login
		http.Redirect(w, r, "/api/login-gh/", http.StatusTemporaryRedirect)
		return
	}

	// Vérifier que le state dans l'URL correspond au cookie
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// We use the code, alongside our client ID and secret to ask GH for an
	// access token to the API.
	code := r.URL.Query().Get("code")
	requestBodyMap := map[string]string{
		"client_id":     GithubClientID,
		"client_secret": GithubClientSecret,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "unable to connect to access_token endpoint", http.StatusInternalServerError)
		return
	}
	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	var ghresp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	json.Unmarshal(respbody, &ghresp)
	userInfo := getGitHubUserInfo(ghresp.AccessToken)

	w.Header().Set("Content-type", "application/json")
	//http.Redirect(w, r, "/", 301)

	var user domain.GitHubUser
	if err := json.Unmarshal(userInfo, &user); err != nil {
		fmt.Println("Erreur de parsing:", err)
		return
	}
	fmt.Println("Login:", user.Login)
	fmt.Println("ID:", user.ID)
	fmt.Println("Email:", user.Email)

	err = authService.GitHub(user.Login, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie, err := r.Cookie("state-Github")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = authService.AuthToken(cookie.Value, user.Login, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/", 301)
}

func getGitHubUserInfo(accessToken string) []byte {
	// Query the GH API for user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	respbody, _ := io.ReadAll(resp.Body)
	return respbody
}
