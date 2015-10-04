package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	File string `envconfig:"optional"`
	Addr string `envconfig:"optional"`
}

var config Config
var r *regexp.Regexp = regexp.MustCompile("(.*)(:\\d*)$")

type ipHandler struct {
}

func (h *ipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := "output"
	if len(config.File) > 0 {
		filename = config.File
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Println("[err] Can't open the file:", err.Error())
	}

	_, err = f.WriteString(removePort(r.RemoteAddr))
	if err != nil {
		log.Println("[err] Can't write the file:", err.Error())
	}

	err = f.Close()
	if err != nil {
		log.Println("[err] Can't close the file:", err.Error())
	}

	log.Println(removePort(r.RemoteAddr))
}

func removePort(str string) string {
	host, _, _ := net.SplitHostPort(str)
	return host
}

func main() {
	if err := envconfig.Init(&config); err != nil {
		log.Println("Can't start:", err.Error())
	}

	log.Println("Configuration:")
	log.Println(config)

	mux := http.NewServeMux()
	mux.Handle("/", &ipHandler{})

	addr := config.Addr
	if len(addr) == 0 {
		addr = ":9004"
	}
	http.ListenAndServe(addr, mux)
}
