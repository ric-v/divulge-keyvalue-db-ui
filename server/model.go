package server

import "github.com/ric-v/divulge-keyvalue-db-ui/database"

type apiResponse struct {
	DBKey    string          `json:"dbkey"`
	FileName string          `json:"filename"`
	DBType   string          `json:"dbtype"`
	Message  string          `json:"message"`
	Data     Datagrid        `json:"data"`
	Error    []errorResponse `json:"error"`
}

type Datagrid struct {
	Columns      []Columns `json:"columns"`
	Rows         []Rows    `json:"rows"`
	InitialState InitState `json:"initialState"`
}

type Columns struct {
	Field      string `json:"field"`
	HeaderName string `json:"headerName"`
	Hide       bool   `json:"hide"`
}

type Rows struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InitState struct {
	Columns InitColumns `json:"columns"`
}

type InitColumns struct {
	ColumnVisibilityModel InitColumnVisibilityModel `json:"columnVisibilityModel"`
}

type InitColumnVisibilityModel struct {
	Id bool `json:"id"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Session struct {
	AccessKey string
	FileName  string
	DBType    string
	DB        database.DB
}
