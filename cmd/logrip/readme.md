# log → rip

Read various log file formats and attempt to parse them into a
cannonical or generic data representation something else might
leverage.

## Dev testing with included data

The data stored in `~/data` is compressed, so uncompress it and pipe
it through the filter to see what happens.

    $ cat /var/log/system.log | go run main.go

Use something like this when developing new grammars.

## To do

* A command-line parameter to indicate which grammar file to use.

* Define a grammar for a web log.

* Figure out how the grammer should add interpolated values (such as
  year) if it's not already available in the log.

* Add a map of field names to canonical representation to the grammar
  file so that we can programmatically pull out the correct parse
  results to fill in the normalized log fact record.

* Add a mechanism to cache compiled regular expresssions so they're
  not re-compiled per line.

* Maybe embed grammar files in the produced binary? Either that, or a
  parameter to point to a data repo location.

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
