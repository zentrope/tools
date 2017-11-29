# docker-proxy

A tiny server-utility to proxy web-requests (via `curl`, say, or even
a browser) to a [Docker REST API][api] generally only exposed via a
hard-to-use Unix domain socket.

## Quick start

Assuming you have a Golang development environment with `$GOPATH/bin`
on your shell's `$PATH`.

    $ go get -u github.com/zentrope/tools/cmd/docker-proxy
    $ docker-proxy

in another terminal:

    $ curl -v http://localhost:2375/version

and then you can go from there exporting the [Docker REST API][api].


[api]: https://docs.docker.com/engine/api/latest/

## Install

You can install the utility as part of your Go environment. On MacOS,
it looks like:

    $ go get -u github.com/zentrope/tools/cmd/docker-proxy

If you haven't got a Golang dev setup, then ... hm.

## Usage

You can just start the application simply:

    $ docker-proxy

which defaults to port `2375` and `/var/run/docker.sock` (used by the
Docker for Mac app).

The help output for `docker-proxy` is:

    Usage of docker-proxy:
      -port int
            Port for this server. (default 2375)
      -socket string
            Path to docker unix domain socket. (default "/var/run/docker.sock")

And you can use it with, say, `curl`:

    $ curl -v http://localhost:2375/version
    $ curl -v http://localhost:2375/v1.34/images/json

and so on.

## Exploration

Translating [Docker REST API][api] documentation into actual URLs can
be challenging.

Here's one way to help figure things out:

I suggest exporting the `DOCKER_HOST` environment variable
in a shell so that you can see what the URLs look like when using
the regular docker client:

    $ export DOCKER_HOST="tcp://localhost:2375"   ## bash
    $ set -gx DOCKER_HOST tcp://localhost:2375    ## fish

When you run something like:

    $ docker search --filter=is_official=true clojure

you can see `docker-proxy` output that looks like:

    $ /v1.34/images/search?filters=%7B%22is-official%22%3A%7B%22true%22%3Atrue%7D%7D&limit=25&term=clojure

This shows that the filter parameter translates to:

    filters={"is-official":{"true":true}}

which is neither intuitive, nor clear from the Docker REST API
docs. To invoke the same `docker` command using `curl` instead ends up
looking like:

    curl -G -XGET 'http://localhost:2375/images/search?term=clojure' \
      --data-urlencode 'filters={"is-official": {"true": true}}'

which isn't pretty, but there you go. Hopefully this is enough to
figure out how to write a client or a monitoring agent.


## License

Copyright © 2017-present Keith Irwin

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
