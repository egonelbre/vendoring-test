package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()

	paths := flag.Args()

	targets := make(map[string]bool)
	for _, p := range paths {
		abspath, err := filepath.Abs(p)
		check(err)

		targets[abspath] = true
		fmt.Println(abspath)
	}

	os.Chdir("..")
	components, err := ioutil.ReadDir("cmd")
	check(err)

	for _, c := range components {
		if dependsOn("./cmd/"+c.Name(), targets) {
			fmt.Println("component changed", c.Name())
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func dependsOn(pkgname string, targets map[string]bool) bool {
	queue, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps,
	}, pkgname)
	check(err)

	queued := map[*packages.Package]bool{}
	for _, p := range queue {
		queued[p] = true
	}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		for _, imp := range p.Imports {
			if !queued[imp] {
				queued[imp] = true
				queue = append(queue, imp)
			}
		}

		for _, f := range p.GoFiles {
			if targets[f] {
				return true
			}
		}

		for _, f := range p.OtherFiles {
			if targets[f] {
				return true
			}
		}

	}

	return false
}
