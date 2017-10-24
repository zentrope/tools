# sslq

A tiny utility to print out an SSL cert in JSON and PEM formats.

## Rationale

This is an initial (and likely end-of-the-line) utility for pulling
down information about SSL certs such that I could monitor them over
time, looking for changes.

## Usage

To print a PEM version of the cert:

    $ ssql cert google.com


To print a JSON version of the cert:

    $ ssql json google.com

Basic usage:

    $ ssql help

    USAGE: ssql COMMAND

    COMMANDS:
      help        - print usage help
      cert <host> - print the host's PEM encoded cert
      json <host> - print the host cert's metadata as JSON

Hopefully this is reasonably self explanatory.

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

and be available on your path.

NOTE: Once built, you can copy this binary to other MacOS workstations
without having to install a Go development environment.

## Considerations

* I don't know how to print out the cert in the same way `openssl`
  does in an automated way.

* An `openssl` equivalent to this doesn't produce the same PEM files,
  but the PEM I get with this method looks just like the ones
  installed on my server. I'm going to assume it's basically correct.

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
