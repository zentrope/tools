package maven

import "encoding/xml"

type Pom struct {
	Name         string       `xml:"name"`
	ArtifactId   string       `xml:"artifactId"`
	Packaging    string       `xml:"packaging"`
	Version      string       `xml:"version"` // If no version, use parents'
	Parent       Dependency   `xml:"parent"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
	Props        PropertyList `xml:"properties"`
}

type PropertyList struct {
	Property []Property `xml:",any"`
}

type Property struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",chardata"`
}

func (pom *Pom) HasParent() bool {
	return len(pom.Parent.GroupId) != 0
}

func (pom *Pom) GetParent() *Dependency {
	return &pom.Parent
}

func (pom *Pom) Properties() map[string]string {

	m := make(map[string]string)

	if len(pom.Version) != 0 {
		m["project.version"] = pom.Version
	}

	for _, p := range pom.Props.Property {
		m[p.XMLName.Local] = p.Value
	}

	return m
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
