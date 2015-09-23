package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/vbatts/go-homeip/ipstore"
)

var (
	defaultDaemonBind = "0.0.0.0"
	defaultDaemonPort = "8080"
)

func daemonCmd(c *cli.Context) {
	if err := ipstore.InitFilename(c.String("db")); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting the app on %s:%s ...", c.String("ip"), c.String("port"))
	log.Printf("Writing to database: %s", c.String("db"))

	// Frame up the server with our catch-all Handler
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", c.String("ip"), c.String("port")),
		Handler:        DefaultRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

var daemonFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "ip",
		Value: defaultDaemonBind,
		Usage: "Set the IP address to serve",
	},
	cli.StringFlag{
		Name:  "port",
		Value: defaultDaemonPort,
		Usage: "Set the port to serve",
	},
	cli.StringFlag{
		Name:  "db",
		Value: "/tmp/ips.db",
		Usage: "Use the following database filename",
	},
}
