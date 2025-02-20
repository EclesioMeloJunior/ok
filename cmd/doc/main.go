package doc

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/parser"
	"github.com/elliotchance/ok/util"
)

type Command struct{}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Description is shown in "ok -help".
func (*Command) Description() string {
	return "view documentation for a package"
}

// Run is the entry point for the "ok test" command.
func (*Command) Run(args []string) {
	if len(args) == 0 {
		args = []string{"."}
	}

	for _, arg := range args {
		packageName := util.PackageNameFromPath("", arg)

		fileNames, err := util.GetAllOKFilesInPath(arg, false)
		check(err)

		// Step 2: Parse all files.
		docs := map[string]string{}
		var funcs []*ast.Func
		var constantNames []string
		constants := map[string]*ast.Literal{}

		for _, fileName := range fileNames {
			data, err := ioutil.ReadFile(fileName)
			check(err)

			p := parser.ParseString(string(data), fileName)

			for _, fn := range p.File.Funcs {
				funcs = append(funcs, fn)
			}

			for name, c := range p.Constants {
				constants[name] = c
				constantNames = append(constantNames, name)
			}

			for _, comment := range p.File.Comments {
				if comment.Func == "" {
					continue
				}

				docs[comment.Func] = comment.String()
			}
		}

		// Sort by function name and output docs.
		sort.Slice(funcs, func(i, j int) bool {
			return funcs[i].Name < funcs[j].Name
		})
		sort.Strings(constantNames)

		fmt.Println("#", packageName)
		fmt.Println()

		// Constants should appear at the top (before functions).
		for _, name := range constantNames {
			if !util.IsPublic(name) {
				continue
			}

			fmt.Printf("- [%s %s](#constants)\n", name, constants[name].Kind)
		}

		if len(constants) > 0 {
			fmt.Println()
		}

		for _, fn := range funcs {
			if !util.IsPublic(fn.Name) {
				continue
			}

			fmt.Printf("- [%s](#%s)\n", fn.String(), fn.Name)
		}
		fmt.Println()

		if len(constants) > 0 {
			fmt.Println("## Constants")
			fmt.Println()

			for _, name := range constantNames {
				if !util.IsPublic(name) {
					continue
				}

				fmt.Println("```")
				fmt.Printf("%s = %s\n", name, constants[name].Value)
				fmt.Println("```")
				fmt.Println()
			}
		}

		for _, fn := range funcs {
			if !util.IsPublic(fn.Name) {
				continue
			}

			fmt.Println("##", fn.Name)
			fmt.Println()

			fmt.Println("```")
			fmt.Println(fn.String())
			fmt.Println("```")
			fmt.Println()

			if doc, ok := docs[fn.Name]; ok {
				fmt.Println(doc)
			} else {
				fmt.Println("No documentation.")
			}
			fmt.Println()
		}
	}
}
