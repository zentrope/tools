package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Dep struct {
	Version string `json:"version"`
	Origin  string `json:"origin"`
}

type Script struct {
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
	Scripts    *Script           `json:"scripts,omitempty"`
	Files      map[string]string `json:"files,omitempty"`
	Deps       map[string]Dep    `json:"deps,omitempty"`
}

func findManifest(path string) *Manifest {
	m := Manifest{}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(contents, &m); err != nil {
		panic(err)
	}
	m.Scripts = nil // remove this fakery if ok to keep in +COMPACT
	return &m
}

func (m *Manifest) addArtifact(artifact Artifact) {
	m.Files[artifact.path] = artifact.hash
	m.Flatsize += artifact.info.Size()
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

type context struct {
	stageDir     string
	manifestFile string
	metaDir      string
}

func newContext() context {
	context := context{}

	flag.StringVar(&context.manifestFile, "manifest",
		"./meta/manifest.json", "Package manifest file.")

	flag.StringVar(&context.stageDir, "stage", "./stage",
		"Location of staged files.")

	flag.StringVar(&context.metaDir, "meta", "./meta",
		"Location of pre/post de/install scripts.")

	flag.Parse()
	return context
}

func main() {

	ctx := newContext()

	manifest := findManifest(ctx.manifestFile)

	artifacts, err := findArtifacts(ctx.stageDir)
	if err != nil {
		log.Fatal(err)
	}

	compact, _ := json.MarshalIndent(manifest, "", "  ")
	fmt.Println(string(compact))

	for _, artifact := range artifacts {
		manifest.addArtifact(artifact)
	}

	manifest.Scripts = &Script{
		PostInstall:   "post-install",
		PreInstall:    "pre-install",
		PostDeInstall: "post-deinstall",
		PreDeInstall:  "pre-deinstall",
	}

	pp, _ := json.MarshalIndent(manifest, "", "  ")
	fmt.Println(string(pp))
}
