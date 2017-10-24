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
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

func textDump(cert *x509.Certificate) {

	// I suspect this might be more interesting if I figured out
	// reflection.

	p := "cert.%-30v = %v\n"

	var spf func(x interface{}) string
	spf = func(x interface{}) string {
		switch t := x.(type) {
		case []string:
			return strings.Join(x.([]string), ", ")
		case []int:
			vals := make([]string, 0)
			for _, v := range x.([]int) {
				vals = append(vals, fmt.Sprintf("%v", v))
			}
			return spf(vals)
		case []asn1.ObjectIdentifier:
			vals := make([]string, 0)
			for _, v := range x.([]asn1.ObjectIdentifier) {
				vals = append(vals, fmt.Sprintf("%v", v))
			}
			return spf(vals)
		case []byte:
			return base64.StdEncoding.EncodeToString(x.([]byte))
		case time.Time:
			return x.(time.Time).Format(time.RFC3339)
		default:
			return fmt.Sprintf("%v", t)
		}
	}

	pp := func(prop string, value interface{}) {
		fmt.Printf(p, prop, spf(value))
	}

	ppATV := func(prefix string, tvs []pkix.AttributeTypeAndValue) {
		for i, tv := range tvs {
			prop := fmt.Sprintf("%v.%v", prefix, i)
			value := spf([]string{spf(tv.Type), spf(tv.Value)})
			pp(prop, value)
		}
	}

	ppkix := func(prefix string, pn pkix.Name) {
		props := map[string]interface{}{
			".country":             pn.Country,
			".organization":        pn.Organization,
			".organizational.unit": pn.OrganizationalUnit,
			".street.address":      pn.StreetAddress,
			".postal.code":         pn.PostalCode,
			".serial.number":       pn.SerialNumber,
			".common.name":         pn.CommonName,
		}

		for k, v := range props {
			pp(prefix+k, v)
		}
		ppATV(prefix+".names", pn.Names)
		ppATV(prefix+".names.extra", pn.ExtraNames)
	}

	pp("version", cert.Version)
	pp("serial.number", cert.SerialNumber)
	ppkix("issuer", cert.Issuer)
	ppkix("subject", cert.Subject)
	pp("notbefore", cert.NotBefore)
	pp("notafter", cert.NotAfter)
	pp("keyusage", cert.KeyUsage)
	pp("signature", cert.Signature)
	pp("signature.algorithm", cert.SignatureAlgorithm)
	pp("basic.contraints.valid", cert.BasicConstraintsValid)
	pp("is.ca", cert.IsCA)
	pp("max.path.len", cert.MaxPathLen)
	pp("max.path.len.zero", cert.MaxPathLenZero)
	pp("subject.key.id", cert.SubjectKeyId)
	pp("authority.key.id", cert.AuthorityKeyId)
	pp("ocsp.server", cert.OCSPServer)
	pp("issuing.certificate.url", cert.IssuingCertificateURL)
	pp("dns.names", cert.DNSNames)
	pp("email.addresses", cert.EmailAddresses)
	pp("ip.addresses", cert.IPAddresses)
	pp("dns.domains.permitted.critical", cert.PermittedDNSDomainsCritical)
	pp("dns.domains.permitted", cert.PermittedDNSDomains)
	pp("dns.domains.excluded", cert.ExcludedDNSDomains)
	pp("crl.distribution.points", cert.CRLDistributionPoints)
	pp("policy.identifiers", cert.PolicyIdentifiers)
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

	certificate, err := getCert(host)
	if err != nil {
		log.Fatal(err)
	}

	switch command {

	case "cert":
		pemDump(certificate)

	case "json":
		jsonDump(certificate)

	case "text":
		textDump(certificate)

	default:
		usage("Unrecognized command: " + command)
	}
}
