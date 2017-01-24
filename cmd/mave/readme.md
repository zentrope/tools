# jdeps

Resolve java dependencies using something that starts/stops fast and
can compile to a static binary. In other words, Go.

## to do

* [ ] rename go-maven or something (jdeps is an existing JDK app)
* [ ] exclusions
* [ ] build a tree (rather than just recursively print)
* [ ] filter out redundant deps
* [ ] cache jars
* [ ] cache sha1 or md5 jar and pom checksums
* [ ] actually validate based on checksum
* [x] resolve versions and props
