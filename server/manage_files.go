package server

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
)

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

	// create the new db file in the temp dir
	log.Println("creating new db file:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file)
	db, err := database.NewDB("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file, dbType)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
		return
	}

	// store the db access in the session
	session.Store(dbKey, Session{dbKey, file, dbType, db})

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
