package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"
)

type clientConfig struct {
	Hosts  []hostConfig `xml:"hosts"`
	Remote remoteConfig `xml:"remote"`
}

type remoteConfig struct {
	URL string `xml:"url,attr"`
}

type hostConfig struct {
	Context  string `xml:"context,attr"`
	Hostname string `xml:"hostname,attr"`
}

func clientCmd(c *cli.Context) {
	if c.Bool("genconfig") {
		cc := new(clientConfig)
		cc.Remote.URL = fmt.Sprintf("http://%s:%s/", defaultDaemonBind, defaultDaemonPort)
		cc.Hosts = []hostConfig{
			{
				Context:  uuid.New(),
				Hostname: os.Getenv("HOSTNAME"),
			},
		}
		output, err := xml.MarshalIndent(cc, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(output))
		os.Exit(0)
	}
}

var clientFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "genconfig",
		Usage: "print a generic config to stdout",
	},
	cli.StringFlag{
		Name:  "file",
		Usage: "configuration file to use",
		Value: path.Join(os.Getenv("HOME"), ".config/homeip.conf"),
	},
	cli.StringFlag{
		Name:  "context",
		Usage: "context to use for hostname (setting this overrrides the config file)",
	},
	cli.StringFlag{
		Name:  "hostname",
		Usage: "hostname to use (setting this overrrides the config file)",
	},
	cli.StringFlag{
		Name:  "remote",
		Usage: "remote daemon to call to (setting this overrrides the config file)",
	},
	cli.BoolFlag{
		Name:  "put",
		Usage: "set the hostname on the configured remote daemon",
	},
}
