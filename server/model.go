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
