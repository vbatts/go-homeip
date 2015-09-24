package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"
)

var defaultClientConfig = path.Join(os.Getenv("HOME"), ".config/homeip.conf")

type clientConfig struct {
	Hostname string
	Remote   string
	Context  string
}

type contextConfig struct {
	Context string `xml:"value,attr"`
}

type remoteConfig struct {
	URL string `xml:"value,attr"`
}

type hostnameConfig struct {
	Name string `xml:"value,attr"`
}

func clientCmd(c *cli.Context) {
	if c.Bool("genconfig") {
		cc := clientConfig{
			Remote:   fmt.Sprintf("http://%s:%s/", defaultDaemonBind, defaultDaemonPort),
			Context:  uuid.New(),
			Hostname: os.Getenv("HOSTNAME"),
		}

		output, err := xml.MarshalIndent(cc, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(output))
		os.Exit(0)
	}

	var config clientConfig
	if _, err := os.Stat(c.String("file")); err == nil {
		buf, err := ioutil.ReadFile(c.String("file"))
		if err != nil {
			log.Fatalf("failed reading %q: %s", c.String("file"), err)
		}
		if err := xml.Unmarshal(buf, &config); err != nil {
			log.Fatalf("failed to read config at %q: %s", c.String("file"), err)
		}
	}
	if c.String("context") != "" {
		config.Context = c.String("context")
	}
	if c.String("remote") != "" {
		config.Remote = c.String("remote")
	}
	if c.String("hostname") != "" {
		config.Hostname = c.String("hostname")
	}
	url, err := url.Parse(config.Remote)
	if err != nil {
		log.Fatalf("%q is not a valid URL: %s", config.Remote, err)
	}
	//fmt.Printf("%#v\n", config)

	if config.Hostname == "" {
		log.Fatalf("hostname is not set")
	}
	url, err = DefaultRouter.Get("getTokenIp").Host(url.Host).URL("host", config.Hostname, "token", config.Context)
	if err != nil {
		log.Fatal(err)
	}

	// Set the ip for this hostname
	if c.Bool("put") {
		resp, err := http.Post(url.String(), "", nil)
		if err != nil {
			log.Printf("%#v", resp)
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("%s failed with %s", url.String(), resp.Status)
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		println(strings.TrimSpace(string(buf)))
		return
	}

	// Fetch the ip for this hostname
	resp, err := http.Get(url.String())
	if err != nil {
		log.Printf("%#v", resp)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%s failed with %s", url.String(), resp.Status)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println(strings.TrimSpace(string(buf)))
}

var clientFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "genconfig",
		Usage: "print a generic config to stdout",
	},
	cli.StringFlag{
		Name:  "file",
		Usage: "configuration file to use",
		Value: defaultClientConfig,
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
