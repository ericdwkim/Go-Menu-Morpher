package main

import (
	"context"
	"encoding/json"
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
	// Check if code is empty string
	if code == "" {
		http.Error(w, "No code in the request", http.StatusBadRequest)
		return
	}
	// Exchange auth code for token
	token, err := oauth2Config.Exchange(context.Background(), code)

	// Error handling for token exchange
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Create HTTP client authenticated with obtained token
	client := oauth2Config.Client(context.Background(), token)

	// Fetch accounts
	accountsResponse, err := client.Get("https://mybusiness.googleapis.com/v4/accounts")
	if err != nil {
		http.Error(w, "Failed to get accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer accountsResponse.Body.Close()

	var accountsData struct {
		Accounts []struct {
			Name string `json:"name"`
		} `json:"accounts"`
	}

	if err := json.NewDecoder(accountsResponse.Body).Decode(&accountsData); err != nil {
		http.Error(w, "Failed to decode accounts response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(accountsData.Accounts) == 0 {
		http.Error(w, "No accounts found", http.StatusInternalServerError)
		return
	}

	accountID := accountsData.Accounts[0].Name // Use the first account ID

	// Fetch locations for the account
	locationsResponse, err := client.Get(fmt.Sprintf("https://mybusiness.googleapis.com/v4/%s/locations", accountID))
	if err != nil {
		http.Error(w, "Failed to get locations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(locationsResponse.Body)

	var locationsData struct {
		Locations []struct {
			Name string `json:"name"`
		} `json:"locations"`
	}

	if err := json.NewDecoder(locationsResponse.Body).Decode(&locationsData); err != nil {
		http.Error(w, "Failed to decode locations response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(locationsData.Locations) == 0 {
		http.Error(w, "No locations found", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "Login Completed. Found locations: %v", locationsData.Locations)
	if err != nil {
		return
	}
}
