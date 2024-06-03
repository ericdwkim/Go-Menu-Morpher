package models

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
