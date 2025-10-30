package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if len(clientID) == 0 {
		log.Fatal("Set GOOGLE_CLIENT_ID env var")
	}
	fmt.Println(clientID)
	redirectURI := "http://localhost:8086/api/google/callback/"

	scope := "openid email profile"
	state, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	fmt.Println(state)

	http.SetCookie(w, &http.Cookie{
		Name:     "state-Google",
		Value:    state.String(),
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	})
	log.Print(state)
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		clientID, redirectURI, scope, state,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("state-Google")
	if err != nil {
		// Au lieu de renvoyer une erreur, on crée un nouveau state
		newState, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Could not generate new state", http.StatusInternalServerError)
			return
		}

		stateCookie = &http.Cookie{
			Name:     "state-Google",
			Value:    newState.String(),
			Path:     "/",
			MaxAge:   int(time.Hour.Seconds()),
			Secure:   r.TLS != nil,
			HttpOnly: true,
		}
		http.SetCookie(w, stateCookie)
	}

	// Récupération du code d'autorisation
	code := r.URL.Query().Get("code")
	log.Printf("Received authorization code: %s", code)
	fmt.Println("google client id : " + os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("google client secret : " + os.Getenv("GOOGLE_CLIENT_SECRET"))

	// Échange du code contre un access_token
	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"redirect_uri":  {"http://localhost:8086/api/google/callback/"},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		http.Error(w, "Unable to exchange code", http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	// Décoder la réponse du token
	var tokenData struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		http.Error(w, "Unable to parse token response", http.StatusInternalServerError)
		return
	}
	log.Printf("Access Token: %s", tokenData.AccessToken) // Affichage du token

	// 🔹 Récupération des informations utilisateur Google
	userInfoBytes, err := getGoogleUserInfo(tokenData.AccessToken)
	if err != nil {
		log.Println("Google user info error:", err)
		http.Error(w, "Unable to get user info", http.StatusInternalServerError)
		return
	}

	// 🔹 Décoder les informations utilisateur dans une structure Go
	var user struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.Unmarshal(userInfoBytes, &user); err != nil {
		http.Error(w, "Unable to parse user info", http.StatusInternalServerError)
		return
	}

	err = authService.Google(user.Name, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie, err := r.Cookie("state-Google")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = authService.AuthToken(cookie.Value, user.Name, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGoogleUserInfo(accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
