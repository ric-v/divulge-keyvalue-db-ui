package server

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fasthttp/router"
	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
)

// Serve godoc - starts the server
func Serve(port string, debug bool) {

	if debug {

		os.Mkdir("logs", 0755)
		w, _ := os.OpenFile("logs/server.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetOutput(w)
	}

	// create a new router
	r := router.New()
	r.HandleOPTIONS = true
	r.GlobalOPTIONS = func(ctx *fasthttp.RequestCtx) {
		// allow cors
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	// /api/v1/ routes
	v1 := r.Group("/api/v1")
	{

		// file upload handler for getting boltdb files
		v1.POST("/upload", uploadFile)

		// create new database file
		v1.POST("/new", newFile)

		// load existing file from browser cache
		v1.POST("/load", loadFile)

		// close the boltdb file and remove the entries
		v1.DELETE("/clear", removeFile)

		// download the updated file
		v1.GET("/download", downloadFile)

		// /api/v1/db routes
		dbroutes := v1.Group("/db")
		{

			// add a new key-value pair to the boltdb file
			dbroutes.POST("/", insertKeyValue)

			// read the boltdb file
			dbroutes.GET("/", listKeyValue)

			// update the boltdb file
			dbroutes.PUT("/{key}", updateKeyValue)

			// delete the boltdb file
			dbroutes.DELETE("/", deleteKeyValue)
		}

		// /api/v1/bucket routes
		bucketroutes := v1.Group("/bucket")
		{

			// add new bucket from boltDB
			bucketroutes.POST("/", addBucket)

			// list all buckets from boltDB
			bucketroutes.GET("/", listBuckets)

			// set the bucket for the boltDB
			bucketroutes.PUT("/{bucket}", setBucket)

			// remove a bucket from boltDB
			bucketroutes.DELETE("/", deleteBucket)
		}
	}

	// serve files from the "./public" directory
	// get the current directory
	cwd, _ := os.Getwd()
	log.Println("getting current working directory:", cwd)

	// if cwd is not the root of the project, then we need to go up one level
	if filepath.Base(cwd) == "cmd" {
		cwd = filepath.Dir(cwd)
	}
	log.Println("updated cwd for setting file server:", cwd)

	// set the file server to serve public files
	r.ServeFiles("/{filepath:*}", cwd+"/ui/build/")

	server := fasthttp.Server{
		Handler:            r.Handler,
		MaxRequestBodySize: 1024 * 1024 * 1024, // 1 GB
		LogAllErrors:       true,
	}
	// serve the handlers on the router
	log.Println("starting server on port:", port)
	log.Fatal(server.ListenAndServe(":" + port))
}

// handleDBSession godoc - loads the db key from header for db access
func handleDBSession(ctx *fasthttp.RequestCtx) (dbSession Session, err error) {

	// get the dbKey from params
	dbKey := string(ctx.Request.Header.Peek("Db-Key"))
	log.Println("dbKey:", dbKey)

	// load the db from user session
	userSession, valid := session.Load(dbKey)
	if !valid {

		var dbTypes []fs.FileInfo
		// get the file name under folder
		dbTypes, err = ioutil.ReadDir("temp" + string(os.PathSeparator))
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

					var dbConn database.DB
					dbConn, err = database.NewDB("temp"+string(os.PathSeparator)+dbType.Name()+string(os.PathSeparator)+dbKey+string(os.PathSeparator)+file, dbType.Name())
					if err != nil {
						log.Println(err)
						err = errors.New(err.Error())
						os.RemoveAll("temp" + string(os.PathSeparator) + dbKey)
						return
					}

					dbSession = Session{dbKey, file, dbType.Name(), dbConn}
					log.Println("userSession: ", dbSession)
					session.Store(dbKey, dbSession)
					return
				}
			}
		}

		// if user session is still nil, return error
		if (dbSession == Session{}) || dbSession.DB == nil {
			log.Println("invalid dbKey")
			err = errors.New("invalid dbKey")
			return
		}
	}

	var ok bool
	if dbSession, ok = userSession.(Session); !ok {
		err = errors.New("invalid session")
		return
	}
	return
}
