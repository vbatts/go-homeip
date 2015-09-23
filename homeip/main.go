package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "homeip"
	app.Usage = "tokenized fetching for setting your homeip"
	app.Authors = []cli.Author{{"Vincent Batts", "vbatts@hashbangbash.com"}}
	app.Commands = []cli.Command{
		{
			Name:    "daemon",
			Aliases: []string{"d"},
			Usage:   "serve the http daemon",
			Action:  daemonCmd,
			Flags:   daemonFlags,
		},
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "call out to a configured homeip daemon",
			Action:  clientCmd,
			Flags:   clientFlags,
		},
	}
	app.Run(os.Args)
}
