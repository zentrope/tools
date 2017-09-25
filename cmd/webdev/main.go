package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func blockUntilShutdownThenDo(fn func()) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill, os.Interrupt, syscall.SIGTERM,
		syscall.SIGKILL, syscall.SIGHUP)
	v := <-sigChan
	log.Printf("Signal: %v\n", v)
	fn()
}

type system struct {
	wwwroot string
	host    string
	files   http.Handler
}

func newSystem(port, wwwroot string) *system {
	return &system{
		wwwroot: wwwroot,
		host:    ":" + port,
		files:   http.FileServer(http.Dir(wwwroot)),
	}
}

func (sys *system) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s %s", r.Method, r.URL.Path)

	rootDir, _ := filepath.Abs(sys.wwwroot)
	path, _ := filepath.Abs(filepath.Join(rootDir, r.URL.Path))

	// If the path doesn't exist, serve up the root path
	// so that browser path manipulator stuff will work.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		r.URL.Path = "/"
	}
	sys.files.ServeHTTP(w, r)

}

func (sys *system) start() {
	log.Println("Starting system.")

	server := http.Server{
		Addr:    sys.host,
		Handler: sys,
	}

	log.Fatal(server.ListenAndServe())
}

func (sys *system) stop() {
	log.Println("Stopping system.")
}

func main() {
	log.Printf("HTTP static content server for dev.")

	var port string
	var docroot string

	flag.StringVar(&port, "port", "3000", "Port for this server.")
	flag.StringVar(&docroot, "docroot", ".", "Document root.")
	flag.Parse()

	system := newSystem(port, docroot)
	go system.start()

	blockUntilShutdownThenDo(system.stop)
}
