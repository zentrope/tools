package main

import (
	"fmt"
	"os"

	"github.com/zentrope/tools/cmd/jdep/maven"
)

func resolve(resolver *maven.Resolver, deps []*maven.Dependency, tab string) {

	if len(deps) == 0 {
		return
	}

	for _, dep := range deps {
		fmt.Printf("%s%s\n", tab, dep)
		resolve(resolver, resolver.GetDeps(dep), tab+"  ")
	}
}

func main() {

	resolver := maven.NewResolver(
		maven.NewRepo("repo1.maven.org", "maven2"),
		maven.NewRepo("clojars.org", "repo"),
	)

	fmt.Printf("jdep (%s)\n\n", resolver.LocalCache)

	logback := maven.NewDep("ch.qos.logback", "logback-classic", "1.1.8")
	pg := maven.NewDep("org.postgresql", "postgresql", "9.4.1212")
	datomic := maven.NewDep("com.datomic", "datomic-free", "0.9.5544")

	deps := []*maven.Dependency{
		logback,
		maven.NewDep("com.mchange", "c3p0", "0.9.5.2"),
		maven.NewDep("http-kit", "http-kit", "2.3.0-alpha1"),
		maven.NewDep("integrant", "integrant", "0.1.5"),
		maven.NewDep("org.clojure", "clojure", "1.9.0-alpha14"),
		maven.NewDep("org.clojure", "core.async", "0.2.395"),
		maven.NewDep("org.clojure", "data.json", "0.2.6"),
		maven.NewDep("org.clojure", "java.jdbc", "0.7.0-alpha1"),
		maven.NewDep("org.clojure", "tools.logging", "0.3.1"),
		pg,
		datomic,
	}

	resolve(resolver, deps, "  ")

	// maven.pprops(resolver.GetProperties(artemis))

	os.Exit(0)
}
