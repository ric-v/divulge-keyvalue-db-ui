package server

import (
	"log"

	"github.com/valyala/fasthttp"
)

// listBucket godoc - loads the list of buckets in a boltdb file
func listBuckets(ctx *fasthttp.RequestCtx) {

	// get the dbKey from params
	dbKey := string(ctx.QueryArgs().Peek("dbKey"))
	log.Println("dbKey:", dbKey)

}

func getBucket(ctx *fasthttp.RequestCtx) {

}
