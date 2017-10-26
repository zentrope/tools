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
		return nil, err
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

	ky := []string{"Digital Signature", "Content Commitment", "KeyEncipherment",
		"DataEnciphermetn", "KeyAgreement", "CertSign", "CRLSign", "EncipherOnly",
		"DecipherOnly"}

	eky := []string{"Any", "ServerAuth", "ClientAuth", "CodeSigning", "EmailProtection",
		"IPSECEndSystem", "IPSEC Tunnel", "IPSEC User", "TimeStamping", "OCSPSigning",
		"Microsoft Server Gated Crypto", "NetscapeServerGatedCrypto"}

	var spf func(x interface{}) string

	cjoin := func(a, b string) string {
		return strings.Trim(strings.Join([]string{a, b}, ", "), ", ")
	}

	spf = func(x interface{}) string {

		switch t := x.(type) {

		case string:
			return strings.TrimSpace(t)
		case x509.KeyUsage:
			return ky[t]
		case x509.ExtKeyUsage:
			return eky[t]
		case time.Time:
			return t.Format(time.RFC3339)

		case []byte:
			return base64.StdEncoding.EncodeToString(t)

		case []string:
			if len(t) == 0 {
				return ""
			}
			return cjoin(spf(t[0]), spf(t[1:]))

		case []net.IP:
			if len(t) == 0 {
				return ""
			}
			return cjoin(spf(t[0]), spf(t[1:]))

		case []x509.ExtKeyUsage:
			if len(t) == 0 {
				return ""
			}
			return cjoin(spf(t[0]), spf(t[1:]))

		case []asn1.ObjectIdentifier:
			if len(t) == 0 {
				return ""
			}
			return cjoin(spf(t[0]), spf(t[1:]))

		default:
			return fmt.Sprintf("%v", t)
		}
	}

	pattern := "cert.%-30v = %v\n"

	pp := func(prop string, value interface{}) {
		fmt.Printf(pattern, prop, spf(value))
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

	ppext := func(px string, exs []pkix.Extension) {
		if len(exs) == 0 {
			pp(px, "")
			return
		}
		for i, ex := range exs {
			prop := fmt.Sprintf("%v.%v.", px, i)
			pp(prop+"id", ex.Id)
			pp(prop+"critical", ex.Critical)
			pp(prop+"value", ex.Value)
		}
	}

	pp("version", cert.Version)
	pp("serial.number", cert.SerialNumber)
	ppkix("issuer", cert.Issuer)
	ppkix("subject", cert.Subject)
	pp("not.valid.before", cert.NotBefore)
	pp("not.valid.after", cert.NotAfter)
	pp("keyusage", cert.KeyUsage)

	ppext("extensions", cert.Extensions)
	ppext("extensions.extra", cert.ExtraExtensions)
	pp("unhandled.critical.extensions", cert.UnhandledCriticalExtensions)
	pp("extended.key.usages", cert.ExtKeyUsage)
	pp("extended.key.usages.unknown", cert.UnknownExtKeyUsage)

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

	certificate, err := getCert(host)
	if err != nil {
		log.Fatal(err)
	}

	switch format {

	case "cert", "pem":
		pemDump(certificate)

	case "json":
		jsonDump(certificate)

	case "text":
		textDump(certificate)

	default:
		usage("Unrecognized output format: '%v'.", format)
	}
}
