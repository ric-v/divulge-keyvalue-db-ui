package server

import (
	"encoding/json"
	"io"
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

	// create the new db file in the temp dir
	log.Println("creating new db file:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename)
	db, err := database.NewDB("temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+files.Filename, dbType)
	if err != nil {
		log.Println(err)
		ctx.Error("Error creating new file: "+err.Error(), fasthttp.StatusInternalServerError)
		os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
		return
	}

	dbSession := Session{dbKey, files.Filename, dbType, db}
	// store the db access in the session
	session.Store(dbKey, dbSession)

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Successfully uploaded and verified db", nil, &dbSession))
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

	dbSession := Session{dbKey, file, dbType, db}
	// store the db access in the session
	session.Store(dbKey, dbSession)

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Successfully created db", nil, &dbSession))
}

// loadFile is the handler for the POST /api/v1/load endpoint.
// Loads previously saved DB from local storage
func loadFile(ctx *fasthttp.RequestCtx) {

	// get the dbKey from header
	dbSession, err := handleDBSession(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Reconnected to db", nil, &dbSession))
}

// removeFile is the handler for the POST /api/v1/db/dbKey endpoint.
func removeFile(ctx *fasthttp.RequestCtx) {

	// get the dbKey from header
	dbSession, err := handleDBSession(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// close the db
	dbSession.DB.CloseDB()
	dbType := dbSession.DBType
	session.Delete(dbSession.DbKey)

	// remove the folder
	log.Println("removing folder:", "temp"+string(os.PathSeparator)+dbType+string(os.PathSeparator)+dbSession.DbKey)
	err = os.RemoveAll("temp" + string(os.PathSeparator) + dbType + string(os.PathSeparator) + dbSession.DbKey)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Cleared DB session and files", nil, &dbSession))
}

// downloadFile is the handler for the GET /api/v1/db/download/dbKey/file endpoint.
// Downloads the boltdb file to the UI.
func downloadFile(ctx *fasthttp.RequestCtx) {

	// get the dbKey from header
	dbSession, err := handleDBSession(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	var file = "temp" + string(os.PathSeparator) + dbSession.DBType + string(os.PathSeparator) + dbSession.DbKey + string(os.PathSeparator) + dbSession.FileName

	// check if db type is boltDB
	if dbSession.DBType == database.BOLT_DB {
		dbSession.DB.Conn().(*database.BoltDB).Sync()
		dbSession.DB.CloseDB()

		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/octet-stream")
		ctx.SetBody(data)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return

		// // copy file to temp dir
		// log.Println("copying file to temp dir:", "temp"+string(os.PathSeparator)+dbSession.DBType+string(os.PathSeparator)+dbSession.DbKey+string(os.PathSeparator)+dbSession.FileName)
		// err = copyFile("temp"+string(os.PathSeparator)+dbSession.DBType+string(os.PathSeparator)+dbSession.DbKey+string(os.PathSeparator)+dbSession.FileName, "temp"+string(os.PathSeparator)+dbSession.FileName+".temp")
		// if err != nil {
		// 	log.Println(err)
		// 	ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		// 	return
		// }
		// file = "temp" + string(os.PathSeparator) + dbSession.FileName + ".temp"
	}
	log.Println("file:", file)

	// return the file to the UI
	ctx.SendFile(file)
}

// copyFile is a helper function to copy a file from one location to another.
func copyFile(src, dst string) error {

	// open source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// open destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// copy file
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return err
}
