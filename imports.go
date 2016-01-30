package rebuilder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type importWalker struct {
	imports map[string]struct{}
}

func GetImportsForDir(dirname string) ([]string, error) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirname, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse directory %s. %v", dirname, err)
	}

	imports := make(map[string]struct{})
	// We're looking for what pkgName depends on.
	for _, pkg := range pkgs {
		iw := &importWalker{imports: imports}
		ast.Walk(iw, pkg)
	}

	return keys(imports), nil
}

func (iw *importWalker) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {
	case *ast.ImportSpec:
		// name := ""
		path := ""
		// if node.Name != nil {
		// 	name = node.Name.Name
		// }
		if node.Path != nil {
			path = node.Path.Value
			path = strings.Trim(path, `"`)
		}
		// log.Printf("name=%s, path=%s", name, path)
		// If path changes we need to rebuild iw.pkg
		iw.imports[path] = struct{}{}
		return nil
	}
	return iw
}

func keys(m map[string]struct{}) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
