package maven

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Types

type Repo struct {
	Host string
	Path string
}

type Resolver struct {
	LocalCache string
	Repos      []*Repo
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
		data, err := repo.getPomFromNet(uri)
		if err == nil {
			return data, nil
		}
	}

	return &Pom{}, errors.New("Not found")
}

func (resolver *Resolver) GetDeps(dep *Dependency) []*Dependency {
	pom, err := resolver.GetPom(dep.Path())
	if err != nil {
		var empty []*Dependency
		return empty
	}

	return pom.Deps()
}

// Implementation

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

func (resolver *Resolver) getPomFromFile(pomUri string) (*Pom, error) {
	location := resolver.LocalCache + "/" + pomUri
	data, err := ioutil.ReadFile(location)

	if err != nil {
		return &Pom{}, err
	}

	return unmarshalPom(data)
}

func unmarshalPom(data []byte) (*Pom, error) {
	var pom Pom
	error := xml.Unmarshal(data, &pom)
	return &pom, error
}
