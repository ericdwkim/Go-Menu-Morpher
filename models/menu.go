package models

// Accounts struct
type Account struct {
	Name string `json:"name"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

// Locations struct
type Location struct {
	Name string `json:"name"`
}

type Locations struct {
	Locations []Location `json:"locations"`
}

// Menu structs
type Price struct {
	CurrencyCode string `json:"currencyCode"`
	Units        int    `json:"units"`
	Nanos        int    `json:"nanos"`
}

type Attributes struct {
	Price Price `json:"price"`
}

type Label struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
}

type MenuItem struct {
	Labels     []Label    `json:"labels"`
	Attributes Attributes `json:"attributes"`
}

type MenuCategory struct {
	Labels []Label    `json:"labels"`
	Items  []MenuItem `json:"items"`
}

type Menu struct {
	Categories []MenuCategory `json:"sections"`
}

type Menus struct {
	Menus []Menu `json:"menus"`
}
