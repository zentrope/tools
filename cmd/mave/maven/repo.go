package maven

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/charmap"
)

// Types

type Resolver struct {
	LocalCache string
	Repos      []*Repo
}

type Repo struct {
	Host string
	Path string
}

// Constructors

func NewResolver(repos ...*Repo) *Resolver {
	cacheDir := os.Getenv("HOME") + "/.m2/repository"
	return &Resolver{cacheDir, repos}
}

// Methods

func NewRepo(host, path string) *Repo {
	return &Repo{host, path}
}

func (resolver *Resolver) GetPom(uri string) (*Pom, error) {

	local, error := resolver.getPomFromFile(uri)
	if error == nil {
		return local, nil
	}

	for _, repo := range resolver.Repos {
		data, err := repo.getPomFromNet(resolver.LocalCache, uri)
		if err != nil {
			fmt.Printf("Unable to find data on net: %#v, %#v\n", repo, err)
			continue
		}

		fmt.Printf("got at %#v, so returning\n", repo)
		return data, nil

	}

	return &Pom{}, errors.New("Not found")
}

func (resolver *Resolver) GetDeps(dep *Dependency) []*Dependency {

	pom, err := resolver.GetPom(dep.Path())
	if err != nil {
		var empty []*Dependency
		return empty
	}

	props := resolver.GetProperties(dep)

	deps := make([]*Dependency, 0)

	for _, d := range pom.Deps() {
		if isProperty(d.Version) {
			d.Version = props[varName(d.Version)]
		} else if d.Version == "" {
			d.Version = props[d.ArtifactId]
		}

		deps = append(deps, d)
	}

	return deps
}

func (resolver *Resolver) GetProperties(dep *Dependency) map[string]string {
	as := resolver.ancestors(dep)
	return mergeProps(as)
}

// Implementation

func (resolver *Resolver) ancestors(dep *Dependency) []*Pom {

	results := make([]*Pom, 0)

	pom, err := resolver.GetPom(dep.Path())

	if err != nil {
		return results
	}

	results = append(results, pom)

	if pom.HasParent() {
		as := resolver.ancestors(pom.GetParent())
		return append(results, as...)
	} else {
		return results
	}
}

func mergeProps(poms []*Pom) map[string]string {

	result := make(map[string]string, 0)

	for i := len(poms) - 1; i >= 0; i-- {
		p := poms[i]

		props := p.Properties()

		for k, v := range props {
			if v != "" {
				result[k] = v
			}
		}
	}

	return result
}

func writeToCache(cacheDir, pomUri string, data []byte) {

	pomPath := cacheDir + "/" + pomUri
	pomDir := filepath.Dir(pomPath)

	fmt.Printf("* caching: %s\n", pomPath)

	if err := os.MkdirAll(pomDir, 0755); err != nil {
		fmt.Printf("* Unable to create parent dirs at [%s], reason [%s].\n", pomDir, err)
	}

	if err := ioutil.WriteFile(pomPath, data, 0644); err != nil {
		fmt.Printf("* Unable to cache POM at [%s], reason [%s].\n", pomPath, err)
	}

}

func (repo *Repo) getPomFromNet(cacheDir string, pomUri string) (*Pom, error) {

	// fmt.Printf("                                * not cached: %s\n", pomUri)

	uri := fmt.Sprintf("https://%s/%s/%s", repo.Host, repo.Path, pomUri)

	resp, err := http.Get(uri)

	defer resp.Body.Close()

	if err != nil {
		return &Pom{}, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("* unable to download %s [%s]\n", uri, resp.Status)
		return &Pom{}, errors.New(fmt.Sprintf("Unable to comply, status %s."))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Pom{}, err
	}

	writeToCache(cacheDir, pomUri, data)

	return unmarshalPom(data)
}

func (resolver *Resolver) getPomFromFile(pomUri string) (*Pom, error) {
	location := resolver.LocalCache + "/" + pomUri
	data, err := ioutil.ReadFile(location)

	if err != nil {
		return &Pom{}, err
	}

	return unmarshalPom(data)
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" {
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("Unknown charset: %s", charset)
}

func unmarshalPom(data []byte) (*Pom, error) {

	r := bytes.NewReader(data)

	var pom Pom
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = makeCharsetReader
	err := decoder.Decode(&pom)
	return &pom, err
}
