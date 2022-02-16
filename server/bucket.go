package server

import (
	"encoding/json"
	"log"

	"github.com/ric-v/divulge-keyvalue-db-ui/database"
	"github.com/valyala/fasthttp"
)

// addBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
// Adds a new bucket to the open DB file.
func addBucket(ctx *fasthttp.RequestCtx) {

	var bucket string
	// get the dbKey from header
	dbSession, err := sessionHandler(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	switch dbSession.DBType {

	case database.BOLT_DB:
		// get the DB type from params
		bucket = string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// add the bucket to the db
		err := dbSession.DB.Conn().(*database.BoltDB).AddBucket(bucket)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

	default:
		log.Println("DB type not supported:", dbSession.DBType)
		ctx.Error("DB type not supported: "+dbSession.DBType, fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Successfully added bucket: "+bucket, nil, &dbSession))
}

// listBucket godoc - loads the list of buckets in a boltdb file
func listBuckets(ctx *fasthttp.RequestCtx) {

	var buckets []string

	// get the dbKey from header
	dbSession, err := sessionHandler(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	switch dbSession.DBType {

	case database.BOLT_DB:
		var err error
		// get the list of buckets from the db
		buckets, err = dbSession.DB.Conn().(*database.BoltDB).ListBuckets()
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("", buckets, &dbSession))
}

// setBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
func setBucket(ctx *fasthttp.RequestCtx) {

	var bucket string
	// get the dbKey from header
	dbSession, err := sessionHandler(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	switch dbSession.DBType {

	case database.BOLT_DB:
		// get the DB type from params
		bucket = string(ctx.UserValue("bucket").(string))
		log.Println("bucket:", bucket)

		// set the bucket in the db
		dbSession.DB.Conn().(*database.BoltDB).SetBucket(bucket)

	default:
		log.Println("DB type not supported:", dbSession.DBType)
		ctx.Error("DB type not supported: "+dbSession.DBType, fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Successfully applied default bucket: "+bucket, nil, &dbSession))
}

// deleteBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
// Removes a bucket from the open DB file.
func deleteBucket(ctx *fasthttp.RequestCtx) {

	var bucket string
	// get the dbKey from header
	dbSession, err := sessionHandler(ctx)
	if err != nil {
		log.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	switch dbSession.DBType {

	case database.BOLT_DB:

		// get the DB type from params
		bucket = string(ctx.QueryArgs().Peek("bucket"))
		log.Println("bucket:", bucket)

		// delete the bucket from the db
		err := dbSession.DB.Conn().(*database.BoltDB).DeleteBucket(bucket)
		if err != nil {
			log.Println(err)
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

	default:
		log.Println("DB type not supported:", dbSession.DBType)
		ctx.Error("DB type not supported: "+dbSession.DBType, fasthttp.StatusInternalServerError)
		return
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(generateResponse("Successfully removed bucket: "+bucket, nil, &dbSession))
}
