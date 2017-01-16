package maven

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//-----------------------------------------------------------------------------

type Repo struct {
	Host string
	Path string
}

func NewRepo(host, path string) *Repo {
	return &Repo{host, path}
}

func (repo *Repo) getPom(pomUri string) (*Pom, error) {

	local, error := getPomFromFile(pomUri)
	if error == nil {
		return local, nil
	}

	data, err := repo.getPomFromNet(pomUri)
	if err == nil {
		return data, nil
	}

	return &Pom{}, errors.New("Not found")
}

func (repo *Repo) getPomFromNet(pomUri string) (*Pom, error) {

	uri := fmt.Sprintf("https://%s/%s/%s", repo.Host, repo.Path, pomUri)

	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if err != nil {
		return &Pom{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Pom{}, err
	}

	return unmarshalPom(data)
}

//-----------------------------------------------------------------------------

// A given pom might have a parent (which might have a parent) each of
// which might contain properties used by descendent poms.

type Pom struct {
	Name         string       `xml:"name"`
	ArtifactId   string       `xml:"artifactId"`
	Packaging    string       `xml:"packaging"`
	Version      string       `xml:"version"` // If no version, use parents'
	Parent       Dependency   `xml:"parent,omitempty"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

func (pom *Pom) HasParent() bool {
	return len(pom.Parent.GroupId) != 0
}

func (pom *Pom) Deps() []*Dependency {

	deps := pom.Dependencies
	results := []*Dependency{}

	for i, dep := range deps {
		if dep.isRuntime() {
			results = append(results, &deps[i])
		}
	}

	return results
}

//-----------------------------------------------------------------------------

type Dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
	Optional   bool   `xml:"optional"`
}

func NewDep(groupId, artifactId, version string) *Dependency {
	return &Dependency{groupId, artifactId, version, "", false}
}

func (d *Dependency) isRuntime() bool {
	return !d.Optional && d.Scope != "test" && d.Scope != "provided"
}

func (d *Dependency) String() string {
	return fmt.Sprintf("[%s/%s \"%s\"]", d.GroupId, d.ArtifactId, d.Version)
}

func (d *Dependency) JarName() string {
	return fmt.Sprintf("%s-%s.jar", d.ArtifactId, d.Version)
}

func (d *Dependency) PomName() string {
	return fmt.Sprintf("%s-%s.pom", d.ArtifactId, d.Version)
}

func (d *Dependency) Path() string {
	aname := strings.Replace(d.GroupId, ".", "/", -1)
	return fmt.Sprintf("%s/%s/%s/%s", aname, d.ArtifactId, d.Version, d.PomName())
}

func (d *Dependency) GetPom(repos []*Repo) (*Pom, error) {
	pomUri := d.Path()

	local, error := getPomFromFile(pomUri)
	if error == nil {
		return local, nil
	}

	for _, repo := range repos {
		data, err := repo.getPomFromNet(pomUri)
		if err == nil {
			return data, nil
		}
	}

	return &Pom{}, errors.New("Not found")
}

func (d *Dependency) GetDeps(repos []*Repo) []*Dependency {
	pom, err := d.GetPom(repos)
	if err != nil {
		var empty []*Dependency
		return empty
	}

	return pom.Deps()
}

//-----------------------------------------------------------------------------

func GetLocalRepo() string {
	home := os.Getenv("HOME")
	return home + "/.m2/repository"
}

func unmarshalPom(data []byte) (*Pom, error) {
	var pom Pom
	error := xml.Unmarshal(data, &pom)
	return &pom, error
}

func getPomFromFile(pomUri string) (*Pom, error) {
	localCache := GetLocalRepo() + "/" + pomUri
	data, err := ioutil.ReadFile(localCache)

	if err != nil {
		return &Pom{}, err
	}

	return unmarshalPom(data)
}
