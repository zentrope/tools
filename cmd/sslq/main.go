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

	"github.com/zentrope/tools/lib"
)

func getCerts(host string) ([]*x509.Certificate, error) {
	hostname := fmt.Sprintf("%v:443", host)

	dialer := &net.Dialer{Timeout: time.Second * 3}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostname, &tls.Config{
		InsecureSkipVerify: true, // We want to seee bad certs, too.
	})

	if err != nil {
		fmt.Println("ERROR: Unable to establish TLS connection, conn ")
		return nil, err
	}

	defer conn.Close()

	state := conn.ConnectionState()
	certs := state.PeerCertificates

	if len(certs) == 0 {
		return nil, errors.New("no certs found")
	}

	return certs, nil
}

func getCert(host string) (*x509.Certificate, error) {
	certs, err := getCerts(host)
	if err != nil {
		return nil, err
	}
	return certs[0], nil
}

func pemDump(certificate *x509.Certificate) {

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})
	fmt.Print(string(buf))
}

func jsonDump(certificate *x509.Certificate) {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")

	if err := enc.Encode(certificate); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", buf.String())
}

func textDump(certs []*x509.Certificate) {
	props := lib.CertProperties(certs)
	for _, k := range props.Keys() {
		fmt.Printf("%-35v = %v\n", k, props.Get(k))
	}
}

func usage(errorMsg string, params ...interface{}) {
	if errorMsg != "" {
		msg := fmt.Sprintf(errorMsg, params...)
		fmt.Printf("ERROR: %v\n\n", msg)
	}
	fmt.Printf("USAGE: ssql hostname [text|cert|pem|json]\n\n")
	fmt.Println("FORMATS:")
	fmt.Println("  cert | pem     - PEM base64-encoded format")
	fmt.Println("  json           - JSON format")
	fmt.Println("  text (default) - key/value text (like Java properties)")
}

func getTargetOrExit(args []string) (string, string) {
	format := "text"

	if len(args) < 2 || args[1] == "help" {
		usage("")
		os.Exit(0)
	}

	if len(args) >= 3 {
		format = args[2]
	}

	host := args[1]
	return host, format
}

func main() {

	host, format := getTargetOrExit(os.Args)

	certs, err := getCerts(host)
	if err != nil {
		log.Fatal(err)
	}

	switch format {

	case "cert", "pem":
		pemDump(certs[0])

	case "json":
		jsonDump(certs[0])

	case "text":
		textDump(certs)

	default:
		usage("Unrecognized output format: '%v'.", format)
	}
}
