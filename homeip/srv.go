package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vbatts/go-homeip/ipstore"
	"github.com/vbatts/go-httplog"
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
