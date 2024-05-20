package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
)

// Declare pointer variable
var (
	oauth2Config *oauth2.Config
)

// Initialize variable instance to Struct
func init() {
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/business.manage"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = "<html><body><a href=\"/login\">Google Log In</a></body></html>"
	_, err := fmt.Fprint(w, htmlIndex)
	if err != nil {
		log.Println(err)
		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// Exchange auth code for token
	token, err := oauth2Config.Exchange(context.Background(), code)

	// Error handling for token exchange
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create HTTP client authenticated with obtained token
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://mybusiness.googleapis.com/v4/accounts")

	// Error handling for account info response request
	if err != nil {
		http.Error(w, "Failed to get accounts: "+err.Error(), http.StatusInternalServerError)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

}
