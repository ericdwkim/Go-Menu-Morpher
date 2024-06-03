package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"menu-morpher-golang/models"
	"net/http"
	"net/url"
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
	var htmlIndex = `<html><body><a href="/login">Google Log In</a></body></html>`
	_, err := fmt.Fprint(w, htmlIndex)
	if err != nil {
		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	authCodeURL := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authCodeURL, http.StatusTemporaryRedirect)
}

func readAndCloseResponse(body io.ReadCloser) ([]byte, error) {
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("Failed to close body: %v", err)
		}
	}()
	return io.ReadAll(body)
}

func exchangeToken(code string) (*oauth2.Token, error) {
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// todo: refactor syntax
func fetchAccounts(client *http.Client) (*models.Accounts, error) {
	accountsResponse, err := client.Get("https://mybusinessaccountmanagement.googleapis.com/v1/accounts")
	if err != nil {
		http.Error(w, "Failed to get accounts: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	body, err := readAndCloseResponse(accountsResponse.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Accounts response body: %s", body)

	var accountsData models.Accounts
	if err := json.Unmarshal(body, &accountsData); err != nil {
		return nil, err
	}
	if len(accountsData.Accounts) == 0 {
		http.Error(w, "No accounts found", http.StatusInternalServerError)
	}

	return &accountsData, nil
}

// todo: testing required
func getAccountId(accountsData models.Accounts) (accountId string) {
	accountID := accountsData.Accounts[0].Name
	log.Printf("Found account: %s", accountID)
	return accountId
}

// todo: refactor
func fetchLocations(client *http.Client, accountID string) (*models.Locations, error) {
	baseURL := fmt.Sprintf("https://mybusinessbusinessinformation.googleapis.com/v1/%s/locations", accountID)

	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	query.Set("readMask", "name")
	reqURL.RawQuery = query.Encode()

	locationsResponse, err := client.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	body, err := readAndCloseResponse(locationsResponse.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Locations response body: %s", body)

	var locationsData models.Locations
	if err := json.Unmarshal(body, &locationsData); err != nil {
		return nil, err
	}

	return &locationsData, nil

}

// todo: ultimate objective
//func fetchMenus(client *http.Client) (*models.Menus, error) {
//
//	resp, err := client.Get("https://mybusiness.googleapis.com/v4/")
//}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in the request", http.StatusBadRequest)
		return
	}

	token, err := exchangeToken(code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), token)

	accountsData, err := fetchAccounts(client)
	//if err != nil {
	//	http.Error(w, "Failed to get accounts: "+err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//if len(accountsData.Accounts) == 0 {
	//	http.Error(w, "No accounts found", http.StatusInternalServerError)
	//	return
	//}
	// todo: clean up
	accountID := accountsData.Accounts[0].Name
	log.Printf("Found account: %s", accountID)

	locationsData, err := fetchLocations(client, accountID)
	if err != nil {
		http.Error(w, "Failed to get locations: "+err.Error(), http.StatusInternalServerError)
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

	locationId := locationsData.Locations[0].Name
	log.Printf("Location ID: '%s' | Account ID '%s'", locationId, accountID)
	return
}
