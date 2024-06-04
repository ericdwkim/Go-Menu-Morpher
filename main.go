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

func fetchAccounts(client *http.Client) (*models.Accounts, error) {
	accountsResponse, err := client.Get("https://mybusinessaccountmanagement.googleapis.com/v1/accounts")
	if err != nil {
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
		return nil, fmt.Errorf("No accounts found")
	}

	return &accountsData, nil
}

func _getAccountId(accountsData *models.Accounts) (string, error) {
	if len(accountsData.Accounts) > 0 {
		var accountId = accountsData.Accounts[0].Name
		log.Printf("Found account: %s", accountId)
		return accountId, nil
	}
	return "", fmt.Errorf("No accounts found")
}

func getAccountId(client *http.Client) (string, error) {

	accountsData, err := fetchAccounts(client)
	if err != nil {
		return "", err
	}
	accountId, err := _getAccountId(accountsData)
	if err != nil {
		return "", err
	}
	return accountId, nil
}

func fetchLocations(client *http.Client, accountId string) (*models.Locations, error) {
	baseURL := fmt.Sprintf("https://mybusinessbusinessinformation.googleapis.com/v1/%s/locations", accountId)

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
	if len(locationsData.Locations) == 0 {
		return nil, fmt.Errorf("No locations found")
	}

	return &locationsData, nil
}

func _getLocationId(locationsData *models.Locations) (string, error) {
	if len(locationsData.Locations) > 0 {
		var locationId = locationsData.Locations[0].Name
		log.Printf("Found location: %s", locationId)
		return locationId, nil
	}
	return "", fmt.Errorf("No locations found")
}

func getLocationId(client *http.Client, accountId string) (string, error) {
	locationsData, err := fetchLocations(client, accountId)
	if err != nil {
		return "", err
	}
	locationId, err := _getLocationId(locationsData)
	if err != nil {
		return "", err
	}
	return locationId, nil
}

func getMenus(client *http.Client, accountId string, locationId string) (*models.Menus, error) {
	getFoodMenusUrl := fmt.Sprintf("https://mybusinessinformation.googleapis.com/v1/%s/locations/%s/foodMenus", accountId, locationId)
	resp, err := client.Get(getFoodMenusUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get food menus: %s", resp.Status)
	}

	body, err := readAndCloseResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	// Save response body of menus as stdout
	err = os.WriteFile("menu.json", body, 0644)
	if err != nil {
		return nil, err
	}

	var menusData models.Menus
	if err := json.Unmarshal(body, &menusData); err != nil {
		return nil, err
	}

	if len(menusData.Menus) == 0 {
		return nil, fmt.Errorf("No food menus found")

	}
	return &menusData, nil
}

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

	accountId, _ := getAccountId(client)

	locationId, _ := getLocationId(client, accountId)

	menus, err := getMenus(client, accountId, locationId)
	if err != nil {
		log.Printf("Failed to get menus: %v", err)
		http.Error(w, "Failed to get menus: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "Login Completed!\nLocation ID: '%s'\nAccount ID '%s'\n", locationId, accountId)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}

}
