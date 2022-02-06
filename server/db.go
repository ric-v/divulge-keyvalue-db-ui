package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	boltdb "github.com/ric-v/golang-key-value-db-browser/bolt-db"
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
	accessKey := uuid.New().String()

	// make new folder
	log.Println("making new folder", "temp"+string(os.PathSeparator)+accessKey)
	os.Mkdir("temp"+string(os.PathSeparator)+accessKey, 0777)

	// save the uploaded file in the temp dir
	log.Println("saving file to dir: ", "temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+files.Filename)
	err = fasthttp.SaveMultipartFile(files, "temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+files.Filename)
	if err != nil {
		log.Println(err)
		ctx.Error("Error getting file: "+err.Error(), fasthttp.StatusBadRequest)
		os.RemoveAll("temp" + string(os.PathSeparator) + accessKey)
		return
	}

	// create the new boltdb file in the temp dir
	log.Println("creating new boltdb file:", "temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+files.Filename)
	db, err := database.NewDB("temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+files.Filename, dbType)
	if err != nil {
		log.Println(err)
		ctx.Error("Error creating new file: "+err.Error(), fasthttp.StatusInternalServerError)
		os.RemoveAll("temp" + string(os.PathSeparator) + accessKey)
		return
	}

	// store the db access in the session
	session.Store(accessKey, Session{accessKey, files.Filename, dbType, db})

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
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

	// generate new accesskey
	accessKey := uuid.New().String()

	// make new folder
	log.Println("making new folder", "temp"+string(os.PathSeparator)+accessKey)
	os.Mkdir("temp"+string(os.PathSeparator)+accessKey, 0777)

	// switch on db type
	switch dbType {

	case database.BOLT_DB:

		// create the new boltdb file in the temp dir
		log.Println("creating new boltdb file:", "temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+file)
		db, err := database.NewDB("temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+file, dbType)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			os.RemoveAll("temp" + string(os.PathSeparator) + accessKey)
			return
		}
		defer db.CloseDB()

		// store the db access in the session
		session.Store(accessKey, Session{accessKey, file, dbType, db})

	case database.BUNT_DB:

		// create the new buntdb file in the temp dir
		log.Println("creating new buntdb file:", "temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+file)
		db, err := database.NewDB("temp"+string(os.PathSeparator)+accessKey+string(os.PathSeparator)+file, dbType)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			os.RemoveAll("temp" + string(os.PathSeparator) + accessKey)
			return
		}

		// store the db access in the session
		session.Store(accessKey, Session{accessKey, file, dbType, db})
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully created boltdb file: " + accessKey,
		Error:    nil,
	})
}

// listKeyValue is the handler for the POST /api/v1/db/accesskey/file endpoint.
// Opens the boltdb file and returns the file key-value paid for rendering in UI.
func listKeyValue(ctx *fasthttp.RequestCtx) {

	var data []database.KeyValuePair
	// get the accesskey from params
	accessKey := string(ctx.QueryArgs().Peek("accesskey"))
	log.Println("accesskey:", accessKey)

	// get the accesskey from params
	file := string(ctx.QueryArgs().Peek("file"))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	// switch on db type
	switch dbType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// load the db from user session
		userSession, valid := session.Load(accessKey)
		if !valid {
			log.Println("invalid accesskey")
			ctx.Error("invalid accesskey", fasthttp.StatusBadRequest)
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
		userSession, valid := session.Load(accessKey)
		if !valid {
			log.Println("invalid accesskey")
			ctx.Error("invalid accesskey", fasthttp.StatusBadRequest)
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
				Hide:       false,
			}, {
				Field:      "key",
				HeaderName: "KEY",
				Flex:       1,
				Editable:   false,
				Hide:       false,
			}, {
				Field:      "value",
				HeaderName: "VALUE",
				Flex:       3,
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
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully opened boltdb file: " + accessKey,
		Data:     datagrid,
		Error:    nil,
	})
}

func getKeyValue(ctx *fasthttp.RequestCtx) {

}

// removeFile is the handler for the POST /api/v1/db/accesskey endpoint.
func removeFile(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.QueryArgs().Peek("accesskey"))
	log.Println("accesskey:", accessKey)

	// load the db from user session
	userSession, valid := session.Load(accessKey)
	if !valid {
		log.Println("invalid accesskey")
		ctx.Error("invalid accesskey", fasthttp.StatusBadRequest)
		return
	}
	userSession.(Session).DB.CloseDB()
	session.Delete(accessKey)

	// remove the folder
	log.Println("removing folder:", "temp"+string(os.PathSeparator)+accessKey)
	os.RemoveAll("temp" + string(os.PathSeparator) + accessKey)

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:   accessKey,
		Message: "Successfully closed boltdb file: " + accessKey,
		Error:   nil,
	})
}

func insertKeyValue(ctx *fasthttp.RequestCtx) {

	type NewEntry struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	// get the value from payload
	var data NewEntry
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	fmt.Println("data:", data)

	switch dbType {

	case "boltdb":

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// open the boltdb file from temp dir
		db, err := boltdb.New("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		defer db.Close()

		// add the key-value pair to the boltdb file
		err = db.Add([]byte(bucket), []byte(data.Key), []byte(data.Value))
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully added key-value pair to boltdb file: " + accessKey,
		Error:    nil,
	})
}

// insertBucket is the handler for the POST /api/v1/db/bucket/accesskey/file endpoint.
// Adds a new bucket to the open DB file.
func insertBucket(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	switch dbType {

	case "boltdb":

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// open the boltdb file from temp dir
		db, err := boltdb.New("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		defer db.Close()

		// remove the bucket from the boltdb file
		err = db.AddBucket([]byte(bucket))
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully added key-value pair to boltdb file: " + accessKey,
		Error:    nil,
	})
}

// deleteBucket is the handler for the POST /api/v1/db/bucket/accesskey/file endpoint.
// Removes a bucket from the open DB file.
func deleteBucket(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	switch dbType {

	case "boltdb":

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// open the boltdb file from temp dir
		db, err := boltdb.New("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		defer db.Close()

		// remove the bucket from the boltdb file
		err = db.RemoveBucket([]byte(bucket))
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully removed bucket from boltdb file: " + accessKey,
		Error:    nil,
	})
}

// deleteKeyValue is the handler for the POST /api/v1/db/accesskey/file/key endpoint.
// Removes a key from the boltdb file.
func deleteKeyValue(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the key from params
	key := string(ctx.UserValue("key").(string))
	key, _ = url.QueryUnescape(key)
	log.Println("key:", key)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	switch dbType {

	case "boltdb":

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// open the boltdb file from temp dir
		db, err := boltdb.New("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		defer db.Close()

		// delete the key from the boltdb file
		err = db.Delete([]byte(bucket), []byte(key))
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully deleted key: " + key,
		Error:    nil,
	})
}

// updateKeyValue is the handler for the POST /api/v1/db/accesskey/file/key endpoint.
// Updates a key in the boltdb file.
func updateKeyValue(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// get the key from params
	key := string(ctx.UserValue("key").(string))
	key, _ = url.QueryUnescape(key)
	log.Println("key:", key)

	// get the DB type from params
	dbType := string(ctx.QueryArgs().Peek("dbtype"))
	log.Println("dbtype:", dbType)

	// get the value from payload
	var value string
	err := json.Unmarshal(ctx.PostBody(), &value)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	log.Println("value:", value)

	switch dbType {

	case "boltdb":

		// get the DB type from params
		bucket := string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// open the boltdb file from temp dir
		db, err := boltdb.New("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		defer db.Close()

		// update the key in the boltdb file
		err = db.Update([]byte(bucket), []byte(key), []byte(value))
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    accessKey,
		FileName: file,
		DBType:   dbType,
		Message:  "Successfully updated key: " + key,
		Error:    nil,
	})
}

// downloadFile is the handler for the GET /api/v1/db/download/accesskey/file endpoint.
// Downloads the boltdb file to the UI.
func downloadFile(ctx *fasthttp.RequestCtx) {

	// get the accesskey from params
	accessKey := string(ctx.UserValue("accesskey").(string))
	log.Println("accesskey:", accessKey)

	// get the file from params
	file := string(ctx.UserValue("file").(string))
	file, _ = url.QueryUnescape(file)
	log.Println("file:", file)

	// return the file to the UI
	ctx.SendFile("temp" + string(os.PathSeparator) + accessKey + string(os.PathSeparator) + file)
}
