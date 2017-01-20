package maven

import (
	"fmt"
	"strings"
)

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
