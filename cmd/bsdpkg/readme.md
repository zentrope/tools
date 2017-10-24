# bsdpkg

**notes**

Since this is a new language, pick small problems to solve before
stitching them all together.

- [x] how to traverse directories for file info
- [x] sha265 a file
- [x] deal with and use command-line options (flag pkg)
- [x] read in a JSON manifest file
- [x] interpolate file chsecksums into manifest
- [x] how to print out a compact version of the manifest (omitempty?)
- [x] interpolate pre/post de/install scripts into manifest
- [x] add files to a tar
- [x] gzip compress a tar
- [x] assert  `/usr/local/*` paths in the archive? (Yes)
- [x] flags to set app version and pkg number
- [x] move this to cmd/bsdpkg/readme.md
- [ ] tests

If I can get this working, packaging scripts can focus mostly on
moving files into an appropriate hierarchy (see `kmake`), then
invoking this.

So, this doesn't have to be complete enough to compete with the true
pkg-create command. Just be something good enough for Clojure (or
similar) apps.

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
