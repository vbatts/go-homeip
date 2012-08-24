package main

import (
	//"github.com/cfdrake/go-gdbm"
	"go-gdbm" // this is a symlink to my clone, at ~/src/go-gdbm
	"os"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
)

// Global variable to access the database
var db *gdbm.Database

func LogHeaders(r *http.Request) {
	fmt.Printf("HEADERS:\n")
	for k, v := range r.Header {
		fmt.Printf("\t%s\n", k)
		for i, _ := range v {
			fmt.Printf("\t\t%s\n", v[i])
		}
	}
}

func LogRequest(r *http.Request, code int) {
	var (
		addr       string
		user_agent string
	)

	user_agent = ""
	addr = r.RemoteAddr

	for k, v := range r.Header {
		if k == "User-Agent" {
			user_agent = strings.Join(v, " ")
		}
		if k == "X-Forwarded-For" {
			addr = strings.Join(v, " ")
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
	LogHeaders(r)
	if (r.Method == "GET") {
		// read from database
		fmt.Fprintf(w, "out")
	} else if (r.Method == "PUT") {
		// write to database
		fmt.Fprintf(w, "in")
	}
}

// Open `filename` as "rw" if it already exists, otherwise "c" to create it.
// If any errors, return in `err`.
func OpenDB(filename string) (db *gdbm.Database, err error) {
	var (
		f_flags string
	)

	f_flags = "rw"

	exists, err := FileExists(filename)
	if (exists == false) {
		f_flags = "c"
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

	ip = *flag.String("ip","0.0.0.0","Set the IP address to serve")
	port = *flag.String("port","8080","Set the port to serve")
	db_file = *flag.String("db","/tmp/ips.db","Use the following database filename")

	flag.Parse()

	// TODO: make this global, such that the request handlers can access it
	db, err = OpenDB(db_file)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db_map, err := db.ToMap()
	if err != nil {
		println(err)
	}
	fmt.Printf("%d\n", len(db_map))

	for k, _ := range db_map {
		fmt.Printf("db_map[%s] = %s\n", k, db_map[k])
	}

	fmt.Printf("%T\n", db_map)
	fmt.Printf("%v\n", db_map)

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


