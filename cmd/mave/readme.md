# mave

Resolve java dependencies using something that starts/stops fast and
can compile to a static binary. In other words, Go.

## To Do

- [ ] Design a project file in json or edn or toml or yaml.
- [ ] Add commands (classpath, tree, etc)
- [ ] Document rationale
- [ ] Filter out redundant deps for a vector of dep peers.
- [ ] Cache jars.
- [ ] Cache sha1 or md5 jar and pom checksums.
- [ ] Actually validate based on checksum.
- [x] Rename go-maven or something (jdeps is an existing JDK app).
- [x] Exclusions.
- [x] Build a tree (rather than just recursively print).
- [x] Resolve versions and props.
