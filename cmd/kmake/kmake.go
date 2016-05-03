package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zentrope/tools/lib/fs"
)

type Project struct {
	Project     string            `json:"project"`
	Version     string            `json:"version"`
	Package     string            `json:"package"`
	ProjectDir  string            `json:"project-dir"`
	BuildCmd    string            `json:"build-cmd"`
	BuildTarget string            `json:"build-target"`
	Artifacts   map[string]string `json:"artifacts"`
}

func (p *Project) projectDir() string {
	home := os.Getenv("HOME")
	dir := p.ProjectDir
	dir = strings.Replace(dir, "~", home, 1)
	return dir
}

func (p *Project) stageDir() string {
	return filepath.Clean("./stage")
}

func (p *Project) build() error {
	fmt.Printf("\nbuilding...\n")

	dir := p.projectDir()

	//name := fmt.Sprintf("%s-%s_%s", p.Project, p.Version, p.Package)
	fmt.Printf(" - '%s' in '%s' using '%s %s'...\n", p.Project, dir, p.BuildCmd, p.BuildTarget)
	cmd := exec.Command(p.BuildCmd, p.BuildTarget)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *Project) stage() error {
	fmt.Printf("\nstaging...\n")
	dir := p.projectDir()
	stage := p.stageDir()

	// Resolve a potentially globbed path into the first matching
	// element and place it in the project directory.
	resolve := func(pattern string) string {
		path := filepath.Join(dir, pattern)
		matches, _ := filepath.Glob(path)
		if len(matches) == 0 {
			return pattern
		}
		return matches[0]
	}

	// map[source]destination
	for k, v := range p.Artifacts {

		sourceFile := resolve(k)
		targetFile := filepath.Join(stage, resolve(v))

		targetDir := filepath.Dir(targetFile)
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			return err
		}

		if err, ok := fs.IsDir(sourceFile); ok {
			fs.CopyDir(sourceFile, targetFile)
		} else {
			fmt.Printf(" - copy: %s -> %s\n", sourceFile, targetFile)
			fs.CopyFile(sourceFile, targetFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	// TODO: -target <build target>
	// TODO: -stage-dir
	// TODO: -f <make file>
	// TODO: -v verbose

	// TODO: clean command -- with detritus in make.json?
	// TODO: build command

	contents, err := ioutil.ReadFile("make.json")
	if err != nil {
		panic(err)
	}

	var project Project
	if err := json.Unmarshal(contents, &project); err != nil {
		panic(err)
	}

	err = project.build()
	if err != nil {
		panic(err)
	}

	err = project.stage()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
