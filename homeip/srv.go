package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/vbatts/go-homeip/ipstore"
	"github.com/vbatts/go-httplog"
)

var DefaultRouter = mux.NewRouter()

func init() {
	DefaultRouter.HandleFunc("/ip", GetIp)
	DefaultRouter.HandleFunc("/ip/{host}", GetIpHost).Methods("GET")
	DefaultRouter.HandleFunc("/ip/{host}", UpdateIpHost).Methods("POST", "PUT")
	DefaultRouter.HandleFunc("/ip/{host}", DeleteIpHost).Methods("DELETE")
	DefaultRouter.HandleFunc("/token", RouteToken)
	DefaultRouter.HandleFunc("/", RouteRoot)
}

// the "/" route
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	fmt.Fprintf(w, "Hello World!\n\n")
}

// provide a random UUID on GET for use with /ip/
func RouteToken(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	fmt.Fprintf(w, "%s", uuid.New())
}

// all things "/ip" (including GET, PUT, etc.)
func GetIp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", httplog.RealIP(r))
}

func GetIpHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host != "" {
		if ok, _ := ipstore.HostExists(host); ok {
			ip, err := ipstore.GetHostIp(host)
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
}
func UpdateIpHost(w http.ResponseWriter, r *http.Request) {
	// write to database
	vars := mux.Vars(r)
	host := vars["host"]
	if host != "" {
		ip := httplog.RealIP(r)
		go func(hostname, ip string) {
			err := ipstore.SetHostIp(hostname, ip)
			if err != nil {
				log.Println(err)
			}
		}(host, ip)

		fmt.Fprintf(w, "%s\n", ip)
	}
}
func DeleteIpHost(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	// delete from database
	vars := mux.Vars(r)
	host := vars["host"]
	if host != "" {
		go func(hostname string) {
			err := ipstore.DropHostIp(hostname)
			if err != nil {
				log.Println(err)
			}
		}(host)

		ip := httplog.RealIP(r)
		fmt.Fprintf(w, "Deleted: %s [last - %s]\n", host, ip)
	}
}
