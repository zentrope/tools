package main

import (
	"fmt"
	"os"

	"github.com/zentrope/tools/cmd/mave/maven"
)

type Tree struct {
	Dependency *maven.Dependency
	Children   []*Tree
}

func NewTree(dep *maven.Dependency) *Tree {
	return &Tree{dep, make([]*Tree, 0)}
}

func resolve(resolver *maven.Resolver, node *Tree) *Tree {

	deps := resolver.GetDeps(node.Dependency)

	for i, _ := range deps {
		child := resolve(resolver, NewTree(deps[i]))
		node.Children = append(node.Children, child)
	}

	return node
}

func _prune(seen map[string]bool, node *Tree) *Tree {

	seen[node.Dependency.RootName()] = true
	for _, e := range node.Dependency.Exclusions {
		seen[e.RootName()] = true
	}

	unseen := make([]*Tree, 0)
	for i, _ := range node.Children {
		child := node.Children[i]
		key := child.Dependency.RootName()
		if !seen[key] {
			unseen = append(unseen, child)
			seen[key] = true
		}
	}

	newChildren := make([]*Tree, 0)

	for i, _ := range unseen {
		child := unseen[i]
		newChildren = append(newChildren, _prune(seen, child))
	}

	node.Children = newChildren

	return node
}

func prune(node *Tree) *Tree {
	return _prune(make(map[string]bool, 0), node)
}

// IO

func _walk(node *Tree, tab string) {
	fmt.Printf("%s%s\n", tab, node.Dependency)
	for _, dep := range node.Children {
		_walk(dep, tab+"  ")
	}
}

func walk(node *Tree) {
	_walk(node, " ")
}

func main() {

	resolver := maven.NewResolver(
		maven.NewRepo("repo1.maven.org", "maven2"),
		maven.NewRepo("clojars.org", "repo"),
	)

	logback := maven.NewDep("ch.qos.logback", "logback-classic", "1.1.8")
	pg := maven.NewDep("org.postgresql", "postgresql", "9.4.1212")
	datomic := maven.NewDep("com.datomic", "datomic-free", "0.9.5554")

	deps := []*maven.Dependency{
		logback,
		maven.NewDep("com.mchange", "c3p0", "0.9.5.2"),
		maven.NewDep("http-kit", "http-kit", "2.3.0-alpha1"),
		maven.NewDep("integrant", "integrant", "0.2.0"),
		maven.NewDep("org.clojure", "clojure", "1.9.0-alpha14"),
		maven.NewDep("org.clojure", "core.async", "0.2.395"),
		maven.NewDep("org.clojure", "data.json", "0.2.6"),
		maven.NewDep("org.clojure", "java.jdbc", "0.7.0-alpha1"),
		maven.NewDep("org.clojure", "tools.logging", "0.3.1"),
		pg,
		datomic,
	}

	// TODO: Gather all deps, _then_ prune afterwards
	for _, d := range deps {
		a := resolve(resolver, NewTree(d))
		b := prune(a)
		walk(b)
	}

	os.Exit(0)
}
