package maven

import (
	"fmt"
	"strings"
)

type Dependency struct {
	GroupId    string       `xml:"groupId"`
	ArtifactId string       `xml:"artifactId"`
	Version    string       `xml:"version"`
	Scope      string       `xml:"scope"`
	Optional   bool         `xml:"optional"`
	Exclusions []Dependency `xml:"exclusions>exclusion"`
}

func NewDep(groupId, artifactId, version string) *Dependency {
	return &Dependency{groupId, artifactId, version, "", false, make([]Dependency, 0)}
}

func (d *Dependency) isRuntime() bool {
	return !d.Optional && d.Scope != "test" && d.Scope != "provided"
}

func (d *Dependency) RootName() string {
	if d.GroupId == d.ArtifactId {
		return d.GroupId
	} else {
		return d.GroupId + "/" + d.ArtifactId
	}
}

func (d *Dependency) String() string {
	excludes := make([]string, 0)
	for _, e := range d.Exclusions {
		excludes = append(excludes, "["+e.RootName()+"]")
	}

	clause := ""
	if len(excludes) != 0 {
		clause = fmt.Sprintf(" :excludes %v", excludes)
	}

	return fmt.Sprintf("[%s \"%s\"]%s",
		d.RootName(), d.Version, clause)
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
