package main

import (
	"fmt"
	"os"

	"github.com/zentrope/tools/cmd/jdep/maven"
)

func resolve(repos []*maven.Repo, deps []*maven.Dependency, tab string) {

	if len(deps) == 0 {
		return
	}

	for _, dep := range deps {
		fmt.Printf("%s%s\n", tab, dep)
		resolve(repos, dep.GetDeps(repos), tab+" ")
	}
}

func main() {

	fmt.Printf("jdep (%s)\n\n", maven.GetLocalRepo())

	// These should be in a JSON file or something.

	repos := []*maven.Repo{
		maven.NewRepo("repo1.maven.org", "maven2"),
		maven.NewRepo("clojars.org", "repo"),
	}

	deps := []*maven.Dependency{
		maven.NewDep("ch.qos.logback", "logback-classic", "1.1.8"),
		maven.NewDep("com.mchange", "c3p0", "0.9.5.2"),
		maven.NewDep("http-kit", "http-kit", "2.3.0-alpha1"),
		maven.NewDep("integrant", "integrant", "0.1.5"),
		maven.NewDep("org.clojure", "clojure", "1.9.0-alpha14"),
		maven.NewDep("org.clojure", "core.async", "0.2.395"),
		maven.NewDep("org.clojure", "data.json", "0.2.6"),
		maven.NewDep("org.clojure", "java.jdbc", "0.7.0-alpha1"),
		maven.NewDep("org.clojure", "tools.logging", "0.3.1"),
		maven.NewDep("org.postgresql", "postgresql", "9.4.1212"),
	}

	resolve(repos, deps, "")

	os.Exit(0)
}
