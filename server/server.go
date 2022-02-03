package server

import (
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func init() {

	// check if required directories exist
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		log.Println("creating ./temp directory")
		os.Mkdir("temp", 0755)
	}

}

func Serve(port string, debug bool) {

	if debug {

		os.Mkdir("logs", 0755)
		w, _ := os.OpenFile("logs/server.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetOutput(w)
	}

	// create a new router
	r := router.New()

	// /api/v1/ routes
	v1 := r.Group("/api/v1")
	{

		v1.POST("/hello", func(ctx *fasthttp.RequestCtx) {
			ctx.Success("text/plain", []byte("Hello, World!"))
		})

		// file upload handler for getting boltdb files
		v1.POST("/upload", uploadFile)

		// create new database file
		v1.POST("/new", newFile)

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

			// read the boltdb file
			dbroutes.GET("/{key}", getKeyValue)

			// update the boltdb file
			dbroutes.PUT("/{key}", updateKeyValue)

			// delete the boltdb file
			dbroutes.DELETE("/{key}", deleteKeyValue)
		}

		// /api/v1/bucket routes
		bucketroutes := v1.Group("/bucket")
		{

			// add new bucket from boltDB
			bucketroutes.POST("/", insertBucket)

			// list all buckets from boltDB
			bucketroutes.GET("/", listBuckets)

			// list all buckets from boltDB
			bucketroutes.GET("/{bucket}", getBucket)

			// remove a bucket from boltDB
			bucketroutes.DELETE("/", deleteBucket)
		}
	}

	// // serve files from the "./public" directory
	// // get the current directory
	// cwd, _ := os.Getwd()
	// log.Println("getting current working directory:", cwd)

	// // if cwd is not the root of the project, then we need to go up one level
	// if filepath.Base(cwd) == "cmd" {
	// 	cwd = filepath.Dir(cwd)
	// }
	// log.Println("updated cwd for setting file server:", cwd)

	// // set the file server to serve public files
	// r.ServeFiles("/{filepath:*}", cwd+"/public/")

	// serve the handlers on the router
	log.Fatal(fasthttp.ListenAndServe(":"+port, r.Handler))
}
