# kmake

**notes, not a proper readme**

Figure out an interesting data structure that can represent arbitrary
projects such that the whole thing can be automated: take in the data
structure, produce a binary.

```clojure
{:project     "server"
 :version     "0.1.2" ;; overwritten by cmd line ops
 :package     "23"
 :project-dir "~/workspace/project"
 :build-cmd   "lein build"
 :manifest    "scripts/manifest.json"
 :artifacts   {"target/app*standalone.jar" "/usr/local/opt/project/app.jar"
               "scripts/project"           "/usr/local/etc/rc.d/project"
               "scripts/control.sh"        "/usr/local/opt/project/control.sh"
               "resources/logback.xml"     "/usr/local/opt/project/logback.xml"}}
```

Huh. Maybe something like that would work. Code would run in a work
dir and create tmp dirs, etc, etc.

Might even work with repackaging other code by making "build-cmd"
point to a script.

Also, use globs so that you could have an artifact like:

```clojure
{"untarred/stuff/*" "/usr/local/opt/project/server/"}
```

Maybe signal intention with the trailing `/`, meaning, copy all
matches into that directory. If no `/`, then take the first glob match
and copy that over.

If I can get all that to work, I can probably make a lein
plugin. Maybe.


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
