package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//-----------------------------------------------------------------------------

type PkgDep struct {
	Version string `json:"version"`
	Origin  string `json:"origin"`
}

type PkgScripts struct {
	PreInstall    string `json:"pre-install"`
	PostInstall   string `json:"post-install"`
	PostDeInstall string `json:"post-deinstall"`
	PreDeInstall  string `json:"pre-deinstall"`
}

type Manifest struct {
	Maintainer string            `json:"maintainer"`
	Desc       string            `json:"desc"`
	Www        string            `json:"www"`
	Name       string            `json:"name"`
	Arch       string            `json:"arch"`
	Flatsize   int64             `json:"flatsize"`
	Prefix     string            `json:"prefix"`
	Comment    string            `json:"comment"`
	Origin     string            `json:"origin"`
	Version    string            `json:"version"`
	Scripts    *PkgScripts       `json:"scripts,omitempty"`
	Files      map[string]string `json:"files,omitempty"`
	Deps       map[string]PkgDep `json:"deps,omitempty"`
}

func newManifest(ctx context) *Manifest {
	m := Manifest{}
	contents, err := ioutil.ReadFile(ctx.manifestFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(contents, &m); err != nil {
		panic(err)
	}
	m.Version = fmt.Sprintf("%s_%s", ctx.version, ctx.pnum)
	m.Scripts = nil // remove this fakery if ok to keep in +COMPACT
	return &m
}

func (m *Manifest) archiveName() string {
	return fmt.Sprintf("%s-%s", m.Name, m.Version)
}

func (m *Manifest) addArtifact(artifact Artifact) {
	m.Files[artifact.path] = artifact.hash
	m.Flatsize += artifact.info.Size()
}

func (m *Manifest) addArtifacts(as []Artifact) {
	for _, a := range as {
		m.addArtifact(a)
	}
}

func (m *Manifest) addScripts(metaDir string) {

	findScript := func(name string) string {
		path := path.Join(metaDir, name)
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return ""
		} else {
			return string(contents)
		}
	}

	m.Scripts = &PkgScripts{
		PreInstall:    findScript("pre-install"),
		PostInstall:   findScript("post-install"),
		PostDeInstall: findScript("post-deinstall"),
		PreDeInstall:  findScript("pre-deinstall"),
	}
}

func (m *Manifest) manifest() string {
	pp, _ := json.MarshalIndent(m, "", "  ")
	return string(pp)
}

//-----------------------------------------------------------------------------

type Artifact struct {
	path     string      // Path in package, e.g., /usr/local/etc
	realPath string      // Real path to file on system
	info     os.FileInfo // All the file info
	hash     string      // sha256 hash as per BSD manifest spec
}

func flatSize(artifacts []Artifact) int64 {
	var size int64
	for _, a := range artifacts {
		size += a.info.Size()
	}
	return size
}

// Get all the artifacts at the root directory, return a slice for
// File structs.
func findArtifacts(root string) ([]Artifact, error) {

	var artifacts []Artifact

	rootDir := filepath.Clean(root)

	shave := func(path string) string {
		return strings.TrimPrefix(path, rootDir)
	}

	gather := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Error walking fs:", err)
				return nil
			}
			hash := "1$" + fmt.Sprintf("%0x", sha256.Sum256(contents))
			artifacts = append(artifacts, Artifact{shave(path), path, info, hash})
		}
		return nil
	}

	err := filepath.Walk(rootDir, gather)
	if err != nil {
		return artifacts, err
	}

	return artifacts, nil
}

//-----------------------------------------------------------------------------
// Archiving
//-----------------------------------------------------------------------------

func compress(filename string) error {

	archiveName := filename + ".gz"

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}

	writer, err := os.Create(archiveName)
	if err != nil {
		return err
	}

	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = archiveName
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)

	return err
}

func tarball(m *Manifest, c *Manifest, artifacts []Artifact, filename string) string {

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	tarball := tar.NewWriter(file)

	writeMeta := func(name string, data []byte) {
		header := &tar.Header{
			Name:    name,
			Mode:    0644,
			Size:    int64(len(data)),
			ModTime: time.Now(),
		}

		tarball.WriteHeader(header)
		tarball.Write(data)
	}

	// Write the +MANIFEST

	writeMeta("+COMPACT_MANIFEST", []byte(c.manifest()))
	writeMeta("+MANIFEST", []byte(m.manifest()))

	for _, artifact := range artifacts {

		header := &tar.Header{
			Name:    artifact.path,
			Mode:    int64(artifact.info.Mode()),
			Size:    artifact.info.Size(),
			ModTime: artifact.info.ModTime(),
		}

		if err != nil {
			panic(err)
		}

		tarball.WriteHeader(header)

		data, err := os.Open(artifact.realPath)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(tarball, data)

		if err != nil {
			panic(err)
		}

		defer data.Close()
	}

	defer func() {
		tarball.Flush()
		tarball.Close()
		file.Close()
	}()

	return filename
}

func tgz(manifest, compact *Manifest, artifacts []Artifact, dir string) error {
	tarName := path.Join(dir, manifest.archiveName()+".tar")
	gzpName := tarName + ".gz"
	tgzName := path.Join(dir, manifest.archiveName()+".tgz")

	tarball(manifest, compact, artifacts, tarName)
	err := compress(tarName)
	if err != nil {
		return err
	}

	if err := os.Remove(tarName); err != nil {
		return err
	}

	if err := os.Rename(gzpName, tgzName); err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------
// Main
//-----------------------------------------------------------------------------

type context struct {
	stageDir     string
	manifestFile string
	metaDir      string
	version      string
	pnum         string
}

func newContext() context {
	context := context{}

	flag.StringVar(&context.manifestFile, "manifest",
		"./meta/manifest.json", "Package manifest file.")

	flag.StringVar(&context.stageDir, "stage", "./stage",
		"Location of staged files.")

	flag.StringVar(&context.metaDir, "meta", "./meta",
		"Location of pre/post de/install scripts.")

	flag.StringVar(&context.version, "version", "0.1.0",
		"Version of the app being packaged.")

	flag.StringVar(&context.pnum, "pnum", "1",
		"Package number.")

	flag.Parse()
	return context
}

func main() {

	ctx := newContext()

	artifacts, err := findArtifacts(ctx.stageDir)
	if err != nil {
		log.Fatal(err)
	}

	compact := newManifest(ctx)

	// TODO: why not a single constructor?
	manifest := newManifest(ctx)
	manifest.addArtifacts(artifacts)
	manifest.addScripts(ctx.metaDir)

	if err := tgz(manifest, compact, artifacts, "."); err != nil {
		panic(err)
	}
}
