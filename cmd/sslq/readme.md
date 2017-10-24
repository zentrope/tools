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

Basic usage (`ssql help`) is:

    USAGE: ssql COMMAND

    COMMANDS:
      help        - print usage help
      cert <host> - print the PEM encoded cert at host
      json <host> - print JSON metadata about the cert at host

This

## Install Binary

Assuming you've got `GOPATH` and `GOBIN` set up properly in your
`.bashrc` (or equivalent):

    export GOPATH="~/Go"
    export GOBIN="${GOPATH}/bin"
    export PATH="${PATH}:${GOBIN}

all you have to do to install this is:

    $ go install

and a binary will show up:

    $ ~/Go/bin/sslq

and be available on your path.

NOTE: Once built, you can copy this binary to other MacOS workstations
without having to install a Go development environment.

## Considerations

* I don't know how to print out the cert in the same way `openssl`
  does in an automated way.

* An `openssl` equivalent to this doesn't produce the same PEM files,
  but the PEM I get with this method looks just like the ones
  installed on my server. I'm going to assume it's basically correct.
