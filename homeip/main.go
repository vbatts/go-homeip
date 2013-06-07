package main

import (
	"flag"
	"fmt"
	"github.com/vbatts/go-homeip/ipstore"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// for debugging request headers
func LogHeaders(r *http.Request) {
	fmt.Printf("HEADERS:\n")
	for k, v := range r.Header {
		fmt.Printf("\t%s\n", k)
		for i, _ := range v {
			fmt.Printf("\t\t%s\n", v[i])
		}
	}
}

func RealIP(r *http.Request) (ip string) {
	ip = r.RemoteAddr

  port_pos := strings.LastIndex(ip,":")
  if port_pos != -1 {
    ip = ip[0:port_pos]
  }

	for k, v := range r.Header {
		if k == "X-Forwarded-For" {
			ip = strings.Join(v, " ")
		}
	}

	return ip
}

// Make an access.log type output
func LogRequest(r *http.Request, code int) {
	var (
		addr       string
		user_agent string
	)

	user_agent = ""
	addr = RealIP(r)

	for k, v := range r.Header {
		if k == "User-Agent" {
			user_agent = strings.Join(v, " ")
		}
	}

	fmt.Printf("%s - - [%s] \"%s %s\" \"%s\" %d %d\n",
		addr,
		time.Now(),
		r.Method,
		r.URL.Path,
		user_agent,
		code,
		r.ContentLength)
}

// The primary route handler.
// setup this way, to have a little more flexibility in the URL.Path matching
func Route_FigureItOut(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/ip/") {
		Route_Ip(w, r)
	} else if r.URL.Path == "/" {
		Route_Root(w, r)
	} else {
		LogRequest(r, 404)
		http.Error(w, "Not Found", 404)
	}
}

// the "/" route
func Route_Root(w http.ResponseWriter, r *http.Request) {
	LogRequest(r, 200)
	fmt.Fprintf(w, "Hello World!\n\n")
}

// all things "/ip" (including GET, PUT, etc.)
func Route_Ip(w http.ResponseWriter, r *http.Request) {
	LogRequest(r, 200)
	//LogHeaders(r)
	if r.Method == "GET" {
		// read from database
		chunks := strings.Split(r.URL.Path, "/")
		if len(chunks) > 2 {
			if ok, _ := ipstore.HostExists(chunks[2]); ok {
				go func(ip string) {
					ip, err := ipstore.GetHostIp(ip)
					if err != nil {
						log.Println(err)
					}
					fmt.Fprintf(w, "%s\n", ip)
				}(chunks[2])
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
			ip := RealIP(r)
			go func(hostname, ip string) {
				err := ipstore.SetHostIp(hostname, ip)
				if err != nil {
					log.Println(err)
				}
			}(chunks[2], ip)

			fmt.Fprintf(w, "%s\n", ip)
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
