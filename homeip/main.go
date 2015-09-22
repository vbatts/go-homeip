package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/vbatts/go-homeip/ipstore"
)

func main() {
	app := cli.NewApp()
	app.Name = "homeip"
	app.Usage = "tokenized fetching for setting your homeip"
	app.Commands = []cli.Command{
		{
			Name:    "daemon",
			Aliases: []string{"d"},
			Usage:   "serve the http daemon",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ip",
					Value: "0.0.0.0",
					Usage: "Set the IP address to serve",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "8080",
					Usage: "Set the port to serve",
				},
				cli.StringFlag{
					Name:  "db",
					Value: "/tmp/ips.db",
					Usage: "Use the following database filename",
				},
			},
			Action: func(c *cli.Context) {
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
			},
		},
	}
	app.Run(os.Args)
}
