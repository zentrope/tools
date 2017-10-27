# sslq

A tiny utility to print out an SSL cert in JSON and PEM or Java
Properties formats. The SSL cert can be loaded from a local PEM file,
or from an Internet host.

## Rationale

This is an initial (and likely end-of-the-line) utility for pulling
down information about SSL certs such that one could monitor them,
looking for changes over time.

## Usage

To print a representation of the certificate, use the following
command pattern:

    $ ssql amazon.com [text|cert|pem|json]

or:

    $ ssql ~/Desktop/cert.pem [text|cert|pem|json]

Where output format defaults to `text` but also supports `cert`,
`pem`, and `json`:

* **sslq amazon.com cert** or **sslq cert.pem pem**<br/> Display the certificate in
  the typical PEM format (the rows of base64 characters):

        -----BEGIN CERTIFICATE-----
        MIIG0zCCBbugAwIBAgIQKC6Ws2t21thSRu27MbIMmDANBgkqhkiG9w0BAQsFADB+
        MQswCQYDVQQGEwJVUzEdMBsGA1UEChMUU3ltYW50ZWMgQ29ycG9yYXRpb24xHzAd
        ...
        uxXJgLRy8z637agLgFpbusEan/jEDHos82JptHuIxaj8QEyOH1PgjEgRmX9YCaM5
        n9MOaiOpkBG7S/a+podi2l70IkZROvU=
        -----END CERTIFICATE-----

* **sslq amazon.com json** or **sslq cert.pem json**<br/> Display the
  certificate as a <small>JSON</small> document:

        // Lots of stuff removed from this example
        {
          "Version": 3,
          "SerialNumber": 53411022063429438665395896543651957912,
          "Issuer": {
            "Country": [ "US" ],
            "Organization": [ "Symantec Corporation" ],
            "OrganizationalUnit": [ "Symantec Trust Network" ],
            "CommonName": "Symantec Class 3 Secure Server CA - G4"
          },
          "Subject": {
            "Country": [ "US" ],
            "Organization": [ "Amazon.com, Inc." ],
            "Locality": [ "Seattle" ],
            "Province": [ "Washington" ],
            "CommonName": "www.amazon.com",
          "NotBefore": "2017-09-20T00:00:00Z",
          "NotAfter": "2018-09-21T23:59:59Z",
          "DNSNames": [
            "amazon.com",
            "amzn.com",
            "uedata.amazon.com"
          ]
        }

    The <small>JSON</small> format also contains a base64 encoded
    version of the complete certificate, not shown here.

* **sslq amazon.com text** or **sslq cert.pem text**<br/> Display the
  certificate as rows of text, using a [Java Properties][jp] format.

        # Same as the JSON version; same things removed.
        cert.version                        = 3
        cert.serial.number                  = 53411022063429438665395896543651957912
        cert.issuer.common.name             = Symantec Class 3 Secure Server CA - G4
        cert.issuer.country                 = US
        cert.issuer.organization            = Symantec Corporation
        cert.issuer.organizational.unit     = Symantec Trust Network
        cert.subject.organization           = Amazon.com, Inc.
        cert.subject.common.name            = www.amazon.com
        cert.subject.country                = US
        cert.not.valid.before               = 2017-09-20T00:00:00Z
        cert.not.valid.after                = 2018-09-21T23:59:59Z
        cert.dns.names                      = amazon.com, amzn.com, uedata.amazon.com...
        cert.verified                       = true

    In the case where the cert is _not_ verified, you'll see something
    like:

        cert.verified                       = false, x509: certificate signed by unknown authority

    The verification stuff isn't a part of any of the other formats.

The text version is especially good for [diffing][diff] the certificate over
time.

[jp]: https://en.wikipedia.org/wiki/.properties
[diff]: https://en.wikipedia.org/wiki/Diff_utility

## Help

The utility is a typical unix-ish command line application with regard
to a `help` parameter:


```text

USAGE: ssql hostname|file.pem [text|cert|pem|json]

FORMATS:
  cert | pem     - PEM base64-encoded format
  json           - JSON format
  text (default) - key/value text (like Java properties)
```

Hopefully this is reasonably self explanatory. If you do something the
utility doesn't understand, you're likely to see the usage
information, too.

## Install Binary

Assuming you've got `GOPATH` and `GOBIN` set up properly in your
`.bashrc` (or equivalent):

    export GOPATH="~/Go"
    export GOBIN="${GOPATH}/bin"
    export PATH="${PATH}:${GOBIN}

all you have to do to install this is:

    $ go install

and a binary will show up:

    $ $GOPATH/bin/sslq

and be available on your `$PATH`.

NOTE: Once built, you can copy this binary to other MacOS workstations
without having to install a Go development environment.

## Considerations

* Might be nice to turn those ASN.1 identifiers into actual text.

## License

Copyright (c) 2017 Keith Irwin

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see
[http://www.gnu.org/licenses/](http://www.gnu.org/licenses/).
