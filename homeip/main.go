package main

import (
	"flag"
	"fmt"
	"github.com/vbatts/go-homeip/ipstore"
	"github.com/vbatts/go-httplog"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// The primary route handler.
// setup this way, to have a little more flexibility in the URL.Path matching
func Route_FigureItOut(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/ip/") {
		Route_Ip(w, r)
	} else if r.URL.Path == "/" {
		Route_Root(w, r)
	} else {
		httplog.LogRequest(r, 404)
		http.Error(w, "Not Found", 404)
	}
}

// the "/" route
func Route_Root(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	fmt.Fprintf(w, "Hello World!\n\n")
}

// all things "/ip" (including GET, PUT, etc.)
func Route_Ip(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	//httplog.LogHeaders(r)
	if r.Method == "GET" {
		// read from database
		chunks := strings.Split(r.URL.Path, "/")
		if len(chunks) > 2 {
			if ok, _ := ipstore.HostExists(chunks[2]); ok {
				ip, err := ipstore.GetHostIp(chunks[2])
				if err != nil {
					log.Println(err)
				}
				fmt.Fprintf(w, "%s\n", ip)
			} else {
				http.Error(w, "No Such Host", 218)
			}
		} else {
			fmt.Fprintf(w, "no hostname\n")
		}
	} else if r.Method == "PUT" || r.Method == "POST" {
		// write to database
		chunks := strings.Split(r.URL.Path, "/")
		if len(chunks) > 2 {
			ip := httplog.RealIP(r)
			go func(hostname, ip string) {
				err := ipstore.SetHostIp(hostname, ip)
				if err != nil {
					log.Println(err)
				}
			}(chunks[2], ip)

			fmt.Fprintf(w, "%s\n", ip)
		}
	} else if r.Method == "DELETE" {
		// delete from database
		chunks := strings.Split(r.URL.Path, "/")
		if len(chunks) > 2 {
			go func(hostname string) {
				err := ipstore.DropHostIp(hostname)
				if err != nil {
					log.Println(err)
				}
			}(chunks[2])

			ip := httplog.RealIP(r)
			fmt.Fprintf(w, "Deleted: %s [last - %s]\n", chunks[2], ip)
		}
	}
}

// Simple check for whather the file exists.
// TODO: be more robust, i.e. permissions etc.
func FileExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func main() {
	var (
		err     error
		ip      string
		port    string
		db_file string = "/tmp/ips.db"
	)

	ip_var := flag.String("ip", "0.0.0.0", "Set the IP address to serve")
	port_var := flag.String("port", "8080", "Set the port to serve")
	db_file_var := flag.String("db", "/tmp/ips.db", "Use the following database filename")

	flag.Parse()

	ip = *ip_var
	port = *port_var
	db_file = *db_file_var

	if err = ipstore.InitFilename(db_file); err != nil {
		panic(err)
	}

	log.Printf("Starting the app on %s:%s ...", ip, port)
	log.Printf("Writing to database: %s", db_file)

	// Frame up the server with our catch-all Handler
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", ip, port),
		Handler:        http.HandlerFunc(Route_FigureItOut),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
