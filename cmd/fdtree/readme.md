# fdtree

Print out a directory tree with each path preceded by a date so you
can sort to find the most recently changed objects.

You can do this with certain versions of `find` (not on macOS) but,
you know, why not just write what you want?


## Install

Standard `go` practice:

    $ go get -u github.com/zentrope/tools/cmd/fdtree

Or:

    $ go get github.com/zentrope/tools
    $ cd $GOPATH/github.com/zentrope/tools/cmd/fdtree
    $ go install

Either approach should install it in `$GOPATH/bin` which is hopefully
on your regular shell path.

## Usage

To use the app, just:

    $ fdtree

and it'll run based on the directory you're in. You can also supply a
directory as a parameter:

    $ fdtree $HOME/Source

to run it against a directory of your choice.

That's it!


## Example

Run a reverse sort on macOS/homebrew nginx config directory:

    $ fdtree /usr/local/etc/nginx | sort -r
    2017-06-09 10:37:33 - /usr/local/etc/nginx
    2017-04-28 13:37:57 - /usr/local/etc/nginx/win-utf
    2017-04-28 13:37:57 - /usr/local/etc/nginx/uwsgi_params.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/uwsgi_params
    2017-04-28 13:37:57 - /usr/local/etc/nginx/scgi_params.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/scgi_params
    2017-04-28 13:37:57 - /usr/local/etc/nginx/nginx.conf.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/mime.types.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/mime.types
    2017-04-28 13:37:57 - /usr/local/etc/nginx/koi-win
    2017-04-28 13:37:57 - /usr/local/etc/nginx/koi-utf
    2017-04-28 13:37:57 - /usr/local/etc/nginx/fastcgi_params.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/fastcgi_params
    2017-04-28 13:37:57 - /usr/local/etc/nginx/fastcgi.conf.default
    2017-04-28 13:37:57 - /usr/local/etc/nginx/fastcgi.conf
    2016-06-21 17:47:11 - /usr/local/etc/nginx/nginx.conf
    2016-05-28 23:10:41 - /usr/local/etc/nginx/nginx.conf~

## License

Copyright (c) 2017-2018 Keith Irwin

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
