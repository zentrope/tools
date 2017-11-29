//
// Copyright © 2017-present Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type dialerFunc func(proto, addr string) (net.Conn, error)

func makeDialer(socketPath string) dialerFunc {
	return func(proto, addr string) (conn net.Conn, err error) {
		return net.Dial("unix", socketPath)
	}
}

func makeProxy(socketPath string) http.HandlerFunc {

	transport := &http.Transport{
		Dial: makeDialer(socketPath),
	}

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		client := &http.Client{Transport: transport}
		resp, err := client.Get("http://docker-to-unix.proxy" + r.URL.String())
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, fmt.Errorf("%v", err))
			return
		}

		for k := range resp.Header {
			if k == "Content-Length" {
				continue
			}
			w.Header().Set(k, resp.Header.Get(k))
		}

		// Set CORS so even browsers can talk to this
		w.Header().Set("Access-Control-Allow-Origin", "*")

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, fmt.Errorf("%v", err))
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}

func main() {

	var httpPort int
	var socketPath string

	flag.IntVar(&httpPort, "port", 2375, "Port for this server.")
	flag.StringVar(&socketPath, "socket", "/var/run/docker.sock", "Path to docker unix domain socket.")
	flag.Parse()

	log.Println("Use '-h', '-help' or '--help' to view docker-proxy options.")
	log.Printf("Delegating requests from http://localhost:%v to '%v'.", httpPort, socketPath)
	log.Printf(" - export DOCKER_HOST='tcp://localhost:%v' to spy on docker commands.", httpPort)

	address := fmt.Sprintf("0.0.0.0:%v", httpPort)
	http.HandleFunc("/", makeProxy(socketPath))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
