package server

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
)

var session sync.Map

// listKeyValue is the handler for the POST /api/v1/bucket/ endpoint.
// Opens the boltdb file and returns the file key-value paid for rendering in UI.
func listKeyValue(ctx *fasthttp.RequestCtx) {

	var data []database.KeyValuePair
	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbkey"))
	log.Println("dbKey:", dbKey)

	// load the db from user session
	userSession, valid := session.Load(dbKey)
	if !valid {
		var err error
		// try restoring user session
		if userSession, err = restoreSession(dbKey); err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}
	sessionInfo := userSession.(Session)

	// switch on db type
	switch sessionInfo.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// load the db from user session
		userSession, valid := session.Load(dbKey)
		if !valid {
			log.Println("invalid dbKey")
			ctx.Error("invalid dbKey", fasthttp.StatusBadRequest)
			return
		}
		db := userSession.(Session).DB

		// open view on the boltdb file
		views, err := db.List()
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		log.Println("views:", views)
		data = views

	case database.BUNT_DB:

		// load the db from user session
		userSession, valid := session.Load(dbKey)
		if !valid {
			log.Println("invalid dbKey")
			ctx.Error("invalid dbKey", fasthttp.StatusBadRequest)
			return
		}
		db := userSession.(Session).DB

		// open view on the boltdb file
		views, err := db.List()
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		log.Println("views:", views)
		data = views
	}

	// init new datagrid object
	var datagrid = Datagrid{
		Columns: []Columns{
			{
				Field:      "id",
				HeaderName: "#",
				Flex:       1,
				Hide:       false,
			}, {
				Field:      "key",
				HeaderName: "KEY",
				Flex:       2,
				Editable:   false,
				Hide:       false,
			}, {
				Field:      "value",
				HeaderName: "VALUE",
				Flex:       14,
				Editable:   true,
				Hide:       false,
			},
		},
		Rows: []Rows{},
		InitialState: InitState{
			Columns: InitColumns{
				ColumnVisibilityModel: InitColumnVisibilityModel{
					Id: false,
				},
			},
		},
	}

	// loop through the data and create a datagrid
	for i, kv := range data {
		datagrid.Rows = append(
			datagrid.Rows,
			Rows{
				ID:    fmt.Sprint(i + 1),
				Key:   kv.Key,
				Value: kv.Value,
			},
		)
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully opened boltdb file: " + dbKey,
		Data:     datagrid,
		Error:    nil,
	})
}

// insertKeyValue is the handler for the POST /api/v1/db/dbKey/file endpoint.
func insertKeyValue(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbkey"))
	log.Println("dbKey:", dbKey)

	// load the db from user session
	userSession, valid := session.Load(dbKey)
	if !valid {
		var err error
		// try restoring user session
		if userSession, err = restoreSession(dbKey); err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}
	sessionInfo := userSession.(Session)

	// get the value from payload
	var data NewEntry
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// add new entry to DB
	err = sessionInfo.DB.Add(data.Key, data.Value)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully added key-value pair to boltdb file: " + dbKey,
		Error:    nil,
	})
}

// deleteKeyValue is the handler for the POST /api/v1/db/dbKey/file/key endpoint.
// Removes a key from the boltdb file.
func deleteKeyValue(ctx *fasthttp.RequestCtx) {

	type deleteKeys struct {
		Keys []string `json:"keys"`
	}

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbkey"))
	log.Println("dbKey:", dbKey)

	// load the db from user session
	userSession, valid := session.Load(dbKey)
	if !valid {
		var err error
		// try restoring user session
		if userSession, err = restoreSession(dbKey); err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}
	sessionInfo := userSession.(Session)

	// get the value from payload
	var keys deleteKeys
	err := json.Unmarshal(ctx.PostBody(), &keys)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	log.Println("keys:", keys)

	switch sessionInfo.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

	}

	// for each selected keys delete from DB
	for _, key := range keys.Keys {

		// delete key from DB
		err = sessionInfo.DB.Delete(key)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully deleted keys from DB",
		Error:    nil,
	})
}

// updateKeyValue is the handler for the POST /api/v1/db/dbKey/file/key endpoint.
// Updates a key in the boltdb file.
func updateKeyValue(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbkey"))
	log.Println("dbKey:", dbKey)

	// get the key from params
	key := string(ctx.UserValue("key").(string))
	log.Println("key:", key)

	// load the db from user session
	userSession, valid := session.Load(dbKey)
	if !valid {
		var err error
		// try restoring user session
		if userSession, err = restoreSession(dbKey); err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}
	sessionInfo := userSession.(Session)

	// get the value from payload
	var data NewEntry
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	log.Println("data:", data)

	switch sessionInfo.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

	}

	// add new entry to DB
	err = sessionInfo.DB.Add(key, data.Value)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully updated key: " + key,
		Error:    nil,
	})
}
