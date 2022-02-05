package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ric-v/divulge-keyvalue-db-ui/server"
)

// Version is the current version of the application
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

	// check if required directories exist
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		log.Println("creating ./temp directory")
		os.Mkdir("temp", 0755)
	}
	// main service
	server.Serve(*port, *debug)
}
