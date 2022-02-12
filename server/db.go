package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
)

var session sync.Map

// uploadFile is the handler for the POST /api/v1/upload endpoint.
// Opens the boltdb file and returns the file handle.
func uploadFile(ctx *fasthttp.RequestCtx) {

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	// get the db file
	files, err := ctx.FormFile("file")
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	log.Println(files.Filename, files.Size)

	// save the file to temp dir
	dbKey := uuid.New().String()

	// make new folder
	log.Println("making new folder", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey)
	os.MkdirAll("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey, 0777)

	// save the uploaded file in the temp dir
	log.Println("saving file to dir: ", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename)
	err = fasthttp.SaveMultipartFile(files, "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename)
	if err != nil {
		log.Println(err)
		ctx.Error("Error getting file: "+err.Error(), fasthttp.StatusBadRequest)
		os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
		return
	}

	// create the new boltdb file in the temp dir
	log.Println("creating new boltdb file:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename)
	db, err := database.NewDB("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename, dbType)
	if err != nil {
		log.Println(err)
		ctx.Error("Error creating new file: "+err.Error(), fasthttp.StatusInternalServerError)
		os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
		return
	}

	// store the db access in the session
	session.Store(dbKey, Session{dbKey, files.Filename, dbType, db})

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: files.Filename,
		DBType:   dbType,
		Message:  "Successfully opened boltdb file",
		Error:    nil,
	})
}

// newFile is the handler for the POST /api/v1/new endpoint.
// Creates a new boltdb file.
func newFile(ctx *fasthttp.RequestCtx) {

	// get the file from params
	file := string(ctx.QueryArgs().Peek("file"))
	log.Println("file:", file)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	// generate new dbKey
	dbKey := uuid.New().String()

	// make new folder
	log.Println("making new folder", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey)
	os.MkdirAll("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey, 0777)

	// switch on db type
	switch dbType {

	case database.BOLT_DB:

		// create the new boltdb file in the temp dir
		log.Println("creating new boltdb file:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file)
		db, err := database.NewDB("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file, dbType)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
			return
		}
		defer db.CloseDB()

		// store the db access in the session
		session.Store(dbKey, Session{dbKey, file, dbType, db})

	case database.BUNT_DB:

		// create the new buntdb file in the temp dir
		log.Println("creating new buntdb file:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file)
		db, err := database.NewDB("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file, dbType)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
			return
		}

		// store the db access in the session
		session.Store(dbKey, Session{dbKey, file, dbType, db})
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully created boltdb file: " + dbKey,
		Error:    nil,
	})
}

// loadFile is the handler for the POST /api/v1/load endpoint.
// Loads previously saved DB from local storage
func loadFile(ctx *fasthttp.RequestCtx) {

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

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: userSession.(Session).FileName,
		DBType:   userSession.(Session).DBType,
		Message:  "Successfully created boltdb file: " + dbKey,
		Error:    nil,
	})
}

// listKeyValue is the handler for the POST /api/v1/db/dbKey/file endpoint.
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

func getKeyValue(ctx *fasthttp.RequestCtx) {

}

// removeFile is the handler for the POST /api/v1/db/dbKey endpoint.
func removeFile(ctx *fasthttp.RequestCtx) {

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
	log.Println("userSession:", userSession)
	userSession.(Session).DB.CloseDB()
	dbType := userSession.(Session).DBType
	session.Delete(dbKey)

	// remove the folder
	log.Println("removing folder:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey)
	err := os.RemoveAll("temp" + string(os.PathSeparator) + dbType + string(os.PathSeparator) + dbKey)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:   dbKey,
		Message: "Successfully closed boltdb file: " + dbKey,
		Error:   nil,
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

// insertBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
// Adds a new bucket to the open DB file.
func insertBucket(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.UserValue("dbKey").(string))
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

	switch sessionInfo.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

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

// deleteBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
// Removes a bucket from the open DB file.
func deleteBucket(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.UserValue("dbKey").(string))
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

	switch sessionInfo.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully removed bucket from boltdb file: " + dbKey,
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

// downloadFile is the handler for the GET /api/v1/db/download/dbKey/file endpoint.
// Downloads the boltdb file to the UI.
func downloadFile(ctx *fasthttp.RequestCtx) {

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
	sessionInfo.DB.CloseDB()

	// return the file to the UI
	ctx.SendFile("temp" + string(os.PathSeparator) + sessionInfo.DBType + string(os.PathSeparator) + dbKey + string(os.PathSeparator) + sessionInfo.FileName)
}

// restoreSession restores the user session from the boltdb / buntdb file if it exists on client
func restoreSession(dbKey string) (userSession Session, err error) {

	// get the file name under folder
	dbTypes, err := ioutil.ReadDir("temp" + string(os.PathSeparator))
	if err != nil {
		log.Println(err)
		err = errors.New("Error reading db folder: " + err.Error())
		return
	}
	log.Println("dbTypes:", dbTypes)

	// iterate over files
	for _, dbType := range dbTypes {

		var dbKeys []fs.FileInfo
		// get the file name under folder
		dbKeys, err = ioutil.ReadDir("temp" + string(os.PathSeparator) + dbType.Name() + string(os.PathSeparator))
		if err != nil {
			log.Println(err)
			err = errors.New("Error reading db folder: " + err.Error())
			return
		}
		log.Println("dbKeys:", dbKeys)

		// iterate over files
		for _, dbkey := range dbKeys {
			log.Println("dbkey:", dbkey.Name(), " | dbKey:", dbkey)

			if dbkey.Name() == dbKey {

				var files []fs.FileInfo
				// get the file name under folder
				files, err = ioutil.ReadDir("temp" + string(os.PathSeparator) + dbType.Name() + string(os.PathSeparator) + dbKey)
				if err != nil {
					log.Println(err)
					err = errors.New("Error reading db folder: " + err.Error())
					return
				}

				// get the file name
				file := files[0].Name()

				userSession.DB, err = database.NewDB("temp"+string(os.PathSeparator)+dbType.Name()+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file, dbType.Name())
				if err != nil {
					log.Println(err)
					err = errors.New(err.Error())
					os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
					return
				}

				userSession = Session{dbKey, file, dbType.Name(), userSession.DB}
				log.Println("userSession: ", userSession)
				session.Store(dbKey, userSession)
			}
		}
	}

	// if user session is still nil, return error
	if (userSession == Session{}) || userSession.DB == nil {
		log.Println("invalid dbKey")
		err = errors.New("invalid dbKey")
		return
	}
	return
}
