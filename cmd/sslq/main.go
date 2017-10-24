//
// Copyright (C) 2017 Keith Irwin
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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func getCert(host string) (*x509.Certificate, error) {
	hostname := fmt.Sprintf("%v:443", host)

	dialer := &net.Dialer{Timeout: time.Second * 3}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostname, &tls.Config{})

	if err != nil {
		log.Fatal("failed to connect: " + err.Error())
	}

	defer conn.Close()

	state := conn.ConnectionState()
	certs := state.PeerCertificates

	if len(certs) == 0 {
		return nil, errors.New("no certs found")
	}

	return certs[0], nil
}

func pemDump(host string) {
	certificate, err := getCert(host)
	if err != nil {
		log.Fatal(err)
	}

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})
	fmt.Print(string(buf))
}

func jsonDump(host string) {
	certificate, err := getCert(host)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")

	if err := enc.Encode(certificate); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", buf.String())
}

func usage(errorMsg string) {
	if errorMsg != "" {
		fmt.Printf("ERROR: %v\n\n", errorMsg)
	}
	fmt.Printf("USAGE: ssql COMMAND\n\n")
	fmt.Println("COMMANDS:")
	fmt.Println("  help        - print usage help")
	fmt.Println("  cert <host> - print the host's PEM encoded cert")
	fmt.Println("  json <host> - print the host cert's metadata as JSON")
}

func main() {

	if len(os.Args) >= 2 && os.Args[1] == "help" {
		usage("")
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		usage("<host> required.")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		usage("Invalid number of parameters.")
		os.Exit(1)
	}

	command := os.Args[1]
	host := os.Args[2]

	switch command {

	case "cert":
		pemDump(host)

	case "json":
		jsonDump(host)

	default:
		usage("Unrecognized command: " + command)
	}
}
