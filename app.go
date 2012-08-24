package main

import (
	//"github.com/cfdrake/go-gdbm"
	"os"
	"fmt"
	"go-gdbm" // this is a symlink to my clone, at ~/src/go-gdbm
	"net/http"
	"time"
	"strings"
)

func LogHeaders(r *http.Request) {
	fmt.Printf("HEADERS:\n")
	for k, v := range r.Header {
		fmt.Printf("\t%s\n", k)
		for i, _ := range v {
			fmt.Printf("\t\t%s\n", v[i])
		}
	}
}

func LogRequest(r *http.Request) {
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
		200,
		r.ContentLength)
}

func rt_hello(w http.ResponseWriter, r *http.Request) {
	LogRequest(r)
	//LogHeaders(r)
	fmt.Fprintf(w, "Hello World!\n\n")
}

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
	db, err := OpenDB("/tmp/ips.db")
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
}
