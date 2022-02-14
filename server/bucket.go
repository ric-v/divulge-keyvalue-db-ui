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

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbKey"))
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

		// add the bucket to the db
		err := sessionInfo.DB.Conn().(*database.BoltDB).AddBucket(bucket)
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
		Message:  "Successfully added bucket to boltdb file: " + dbKey,
		Error:    nil,
	})
}

// listBucket godoc - loads the list of buckets in a boltdb file
func listBuckets(ctx *fasthttp.RequestCtx) {

	var buckets []string
	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbKey"))
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

		var err error
		// get the list of buckets from the db
		buckets, err = sessionInfo.DB.Conn().(*database.BoltDB).ListBuckets()
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
		Data:     buckets,
		Error:    nil,
	})
}

// setBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
func setBucket(ctx *fasthttp.RequestCtx) {
	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbKey"))
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
		bucket := string(ctx.UserValue("bucket").(string))
		log.Println("bucket:", bucket)

		// set the bucket in the db
		sessionInfo.DB.Conn().(*database.BoltDB).SetBucket(bucket)
	}

	// return success message to UI
	json.NewEncoder(ctx).Encode(apiResponse{
		DBKey:    dbKey,
		FileName: sessionInfo.FileName,
		DBType:   sessionInfo.DBType,
		Message:  "Successfully applied default bucket: " + dbKey,
		Error:    nil,
	})
}

// deleteBucket is the handler for the POST /api/v1/db/bucket/dbKey/file endpoint.
// Removes a bucket from the open DB file.
func deleteBucket(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbKey"))
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

		// delete the bucket from the db
		err := sessionInfo.DB.Conn().(*database.BoltDB).DeleteBucket(bucket)
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
		Message:  "Successfully removed bucket from boltdb file: " + dbKey,
		Error:    nil,
	})
}
