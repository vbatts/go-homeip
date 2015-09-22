package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vbatts/go-homeip/ipstore"
)

var (
	ip      = flag.String("ip", "0.0.0.0", "Set the IP address to serve")
	port    = flag.String("port", "8080", "Set the port to serve")
	db_file = flag.String("db", "/tmp/ips.db", "Use the following database filename")
)

func main() {
	flag.Parse()

	if err := ipstore.InitFilename(*db_file); err != nil {
		panic(err)
	}

	log.Printf("Starting the app on %s:%s ...", *ip, *port)
	log.Printf("Writing to database: %s", *db_file)

	// Frame up the server with our catch-all Handler
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", *ip, *port),
		Handler:        DefaultRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
