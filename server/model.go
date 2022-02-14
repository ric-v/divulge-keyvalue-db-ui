package server

import "github.com/ric-v/divulge-keyvalue-db-ui/database"

type apiResponse struct {
	DBKey    string          `json:"dbkey"`
	FileName string          `json:"filename"`
	DBType   string          `json:"dbtype"`
	Message  string          `json:"message"`
	Data     interface{}     `json:"data"`
	Error    []errorResponse `json:"error"`
}

// Datagrid - Datagrid struct for data visualization in UI
type Datagrid struct {
	Columns      []Columns `json:"columns"`
	Rows         []Rows    `json:"rows"`
	InitialState InitState `json:"initialState"`
}

// Columns - Columns for the datagrid header secion
type Columns struct {
	Field      string `json:"field"`
	HeaderName string `json:"headerName"`
	Flex       int    `json:"flex"`
	Editable   bool   `json:"editable"`
	Hide       bool   `json:"hide"`
}

// Rows - Rows for the datagrid body section
type Rows struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InitState - Initial state for the datagrid view
type InitState struct {
	Columns InitColumns `json:"columns"`
}

// InitColumns - Initial columns for the datagrid view
type InitColumns struct {
	ColumnVisibilityModel InitColumnVisibilityModel `json:"columnVisibilityModel"`
}

// InitColumnVisibilityModel - Initial column visibility model for the datagrid view
type InitColumnVisibilityModel struct {
	Id bool `json:"id"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server - for managing db sessions with unique dbKey for each open session from UI
type Session struct {
	dbKey    string
	FileName string
	DBType   string
	DB       database.DB
}

// NewEntry - for creating new entry in the database
type NewEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
