package main

import (
	"flag"
	"fmt"

	"github.com/ric-v/divulge-keyvalue-db-ui/server"
)

var Version = "undefined"

func main() {

	// flags
	showVersion := flag.Bool("version", false, "print version and exit")
	port := flag.String("port", "8080", "port to listen on")
	debug := flag.Bool("debug", true, "enable debug logging")
	flag.Parse()

	// version
	if *showVersion {
		fmt.Println("Version:", Version)
		return
	}

	// main service
	server.Serve(*port, *debug)
}
