# webdev

A simple HTTP server for serving static "front-end" HTML, CSS, and JavaScript.

## Rationale

Say you're developing a web application but are _not_ using the elaborate node/npm/yarn/jspm/gulp build chain stuff...

When building a single-page-web application you might want to use URL routing and history management (such as react-router-dom) so that users can bookmark `/reports` or `/account/id/2` kinds of things. These routes are synthetic in that they don't represent actual directories on a server, but, rather, code paths through your single-page-app. In most cases, such as with `nginx`, you'd use a [`try_files`][tf] directive to substitute `index.html` for not-found resources.

You can't do this sort of thing with applications like [`darkhttpd`][dh] (available via homebrew) or [`python -m SimpleHTTPServer`][ps], but you can with this app.

## Install Binary

Assuming you've got `GOPATH` and `GOBIN` set up properly in your `.bashrc` (or equivalent):

    export GOPATH="~/Go"
    export GOBIN="${GOPATH}/bin"
    export PATH="${PATH}:${GOBIN}

all you have to do to install this is:

    $ go install

and a binary will show up:

    $ ~/Go/bin/webdev

and be available on your path.

NOTE: Once built, you can copy this binary to with macos (or whatever) workstations without having to install a Go development environment.

## Usage

Simple:

    $ webdev -docroot . -port 3000

Those are default values, so you can omit them if you're fine with port 3000 serving files in your current directory.

## Future / API Proxy

I should add a proxy feature so that the app can delegate to an API server or two. Not sure how to do it in any kind of flexible way, but, probably, for me, using a param instead of a hostname would work:

    const apiUrl = window.location.href + "?proxy-tag=api"

or even better yet, setting a request header:

    fetch("/api/stuff", { headers: {"Proxy-Tag", "api"} }).then(...)

The `webdev` app can use the `proxy-tag` value to figure out which reverse-proxied service to delegate to.

The advantages of this system:

* Can be ignored in production.
* Doesn't require static configuration.
* Doesn't rely on any specific DNS convention.
* Avoids CORS issues
* Doesn't require the backend service to adjust self-referential URLs

Hm. Worth trying.

[tf]: http://nginx.org/en/docs/http/ngx_http_core_module.html#try_files
[dh]: https://unix4lyfe.org/darkhttpd/
[ps]: http://2ality.com/2014/06/simple-http-server.html
