package main

import (
	//"github.com/cfdrake/go-gdbm"
	"go-gdbm" // this is a symlink to my clone, at ~/src/go-gdbm
	"os"
	//"sync"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
)

// Global variable to access the database
var db *gdbm.Database

// use this channel to wrap the db r/w
var c = make(chan int)

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
	if (strings.HasPrefix(r.URL.Path, "/ip")) {
		Route_Ip(w,r)
	} else if (r.URL.Path == "/") {
		Route_Root(w,r)
	} else {
		LogRequest(r, 404)
		http.Error(w,"Not Found", 404)
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
	if (r.Method == "GET") {
		// read from database
		chunks := strings.Split(r.URL.Path, "/")
		if (len(chunks) > 2) {
			if (db.Exists(chunks[2])) {
				go func(key string) {
					ip, err := db.Fetch(key)
					if (err != nil) {
						fmt.Printf("%s\n", err)
					}
					fmt.Fprintf(w, "%s\n", ip)
					c <- 1 // send a signal
				}(chunks[2])
				<- c // wait for channel to clear
			} else {
				http.Error(w,"No Such Host", 218)
			}
		} else {
			fmt.Fprintf(w, "no hostname\n")
		}
	} else if (r.Method == "PUT") {
		// write to database
		chunks := strings.Split(r.URL.Path, "/")
		if (len(chunks) > 2) {
			ip := RealIP(r)
			if (strings.Contains(ip,":")) {
				ip_chunks := strings.Split(ip,":")
				ip = ip_chunks[0]
			}
			if (db.Exists(chunks[2])) {
				go func(key string, val string) {
					err := db.Replace(key,val)
					if (err != nil) {
						fmt.Printf("%s\n", err)
					}
					c <- 1 // send a signal
				}(chunks[2],ip)
				<- c // wait for channel to clear
			} else {
				go func(key string, val string) {
					err := db.Insert(chunks[2], ip)
					if (err != nil) {
						fmt.Printf("%s\n", err)
					}
					c <- 1 // send a signal
				}(chunks[2],ip)
				<- c // wait for channel to clear
			}
			fmt.Fprintf(w,"%s\n", ip)
		}
	}
}

// Open `filename` as "rw" if it already exists, otherwise "c" to create it.
// If any errors, return in `err`.
func OpenDB(filename string) (db *gdbm.Database, err error) {
	var (
		f_flags string
	)

	f_flags = "c"

	exists, err := FileExists(filename)
	if (exists) {
		f_flags = "w"
	} else if (err != nil) {
		return db, err
	}

	db, err = gdbm.Open(filename, f_flags)
	if err != nil {
		return db, err
	}

	return db, nil
}

// Simple check for whather the file exists.
// TODO: be more robust, i.e. permissions etc.
func FileExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if (err != nil) {
		return false, nil
	}
	return true, nil
}

func main() {
	var (
		err error
		ip string
		port string
		db_file string
	)

	ip_var := flag.String("ip","0.0.0.0","Set the IP address to serve")
	port_var := flag.String("port","8080","Set the port to serve")
	db_file_var := flag.String("db","/tmp/ips.db","Use the following database filename")

	flag.Parse()

	ip = *ip_var
	port = *port_var
	db_file = *db_file_var


	db, err = OpenDB(db_file)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("%s - Starting the app on %s:%s ...\n", time.Now(), ip, port)

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


