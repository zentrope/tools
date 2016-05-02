# bsdpkg

**TODO**

Since this is a new language, pick small problems to solve before
stitching them all together.

- [x] how to traverse directories for file info
- [x] sha265 a file
- [x] deal with and use command-line options (flag pkg)
- [x] read in a JSON manifest file
- [x] interpolate file chsecksums into manifest
- [ ] interpolate pre/post de/install scripts into manifest
- [ ] how to print out a compact version of the manifest (omitempty?)
- [ ] add files to a tar
- [ ] gzip compress a tar
- [ ] can you just assert  `/usr/local/*` paths in the archive?

If I can get this working, packaging scripts can focus mostly on
moving files into an appropriate hierarchy, then invoking this.

So, this doesn't have to be complete enough to compete with the true
pkg-create command. Just be something good enough for Clojure (or
similar) apps.

**OUT OF SCOPE**

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
