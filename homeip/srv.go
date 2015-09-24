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
	DefaultRouter.HandleFunc("/ip", getIp)
	DefaultRouter.HandleFunc("/ip/{host}", getIpHost).Methods("GET")
	DefaultRouter.HandleFunc("/ip/{host}", updateIpHost).Methods("POST", "PUT")
	DefaultRouter.HandleFunc("/ip/{host}", deleteIpHost).Methods("DELETE")
	DefaultRouter.HandleFunc("/ip/{host}/token/{token}", getIpHostToken).Methods("GET").Name("getTokenIp")
	DefaultRouter.HandleFunc("/ip/{host}/token/{token}", updateIpHostToken).Methods("POST", "PUT").Name("putTokenIp")
	DefaultRouter.HandleFunc("/ip/{host}/token/{token}", deleteIpHostToken).Methods("DELETE")
	DefaultRouter.HandleFunc("/token", routeToken)
	DefaultRouter.HandleFunc("/", routeRoot)
}

// the "/" route
func routeRoot(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	fmt.Fprintf(w, "Hello World!\n\n")
}

// provide a random UUID on GET for use with /ip/
func routeToken(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	fmt.Fprintf(w, "%s", uuid.New())
}

// all things "/ip" (including GET, PUT, etc.)
func getIp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", httplog.RealIP(r))
}

func getIpHost(w http.ResponseWriter, r *http.Request) {
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
func updateIpHost(w http.ResponseWriter, r *http.Request) {
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
func deleteIpHost(w http.ResponseWriter, r *http.Request) {
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
func getIpHostToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	token := vars["token"]
	if host != "" && token != "" {
		if ok, _ := ipstore.HostExistsToken(host, token); ok {
			ip, err := ipstore.GetHostIpToken(host, token)
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
func updateIpHostToken(w http.ResponseWriter, r *http.Request) {
	// write to database
	vars := mux.Vars(r)
	host := vars["host"]
	token := vars["token"]
	if host != "" && token != "" {
		ip := httplog.RealIP(r)
		go func(hostname, ip string) {
			err := ipstore.SetHostIpToken(hostname, ip, token)
			if err != nil {
				log.Println(err)
			}
		}(host, ip)

		fmt.Fprintf(w, "%s\n", ip)
	}
}
func deleteIpHostToken(w http.ResponseWriter, r *http.Request) {
	httplog.LogRequest(r, 200)
	// delete from database
	vars := mux.Vars(r)
	host := vars["host"]
	token := vars["token"]
	if host != "" && token != "" {
		go func(hostname string) {
			err := ipstore.DropHostIpToken(hostname, token)
			if err != nil {
				log.Println(err)
			}
		}(host)

		ip := httplog.RealIP(r)
		fmt.Fprintf(w, "Deleted: %s [last - %s]\n", host, ip)
	}
}
