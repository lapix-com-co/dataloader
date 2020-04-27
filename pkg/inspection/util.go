package inspection

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// CheckTypes takes multiple paths, loads the packages and checks if the given types
// exists within the given path.
func CheckTypes(where []string, typesName []string) ([]bool, error) {
	var results = make([]bool, len(typesName))

	if len(typesName) == 0 {
		return nil, errors.New("the type can not be empty")
	}

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, where...)
	if err != nil {
		return nil, fmt.Errorf("could not load the package: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("found %d packages", len(pkgs))
	}

	for index, typeName := range typesName {
		typeImpl := pkgs[0].Types.Scope().Lookup(typeName)
		results[index] = typeImpl != nil
	}

	return results, nil
}

// Names return the packages names from the given paths.
func Names(where []string) ([]string, error) {
	var names []string

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, where...)
	if err != nil {
		return nil, fmt.Errorf("could not load the package: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("found %d packages", len(pkgs))
	}

	for _, pkg := range pkgs {
		names = append(names, pkg.Name)
	}

	return names, nil
}

// Methods returns the interface's methods names.
func Methods(where []string, interfaceName string) ([]string, error) {
	var methods []string

	if interfaceName == "" {
		return nil, errors.New("the iterface name can not be empty")
	}

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, where...)
	if err != nil {
		return nil, fmt.Errorf("could not load the package: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("found %d packages", len(pkgs))
	}

	for _, file := range pkgs[0].Syntax {
		for _, decl := range file.Decls {
			d, ok := decl.(*ast.GenDecl)
			if !ok || d.Tok != token.TYPE {
				continue
			}

			for _, spec := range d.Specs {
				tSpec := spec.(*ast.TypeSpec)
				if tSpec.Name.Name == interfaceName {
					ast.Inspect(tSpec, func(n ast.Node) bool {
						if n == nil {
							return true
						}

						tSpec, ok := n.(*ast.TypeSpec)
						if !ok {
							return true
						}

						tInterface, ok := tSpec.Type.(*ast.InterfaceType)
						if !ok {
							return true
						}

						for _, methodList := range tInterface.Methods.List {
							methods = append(methods, methodList.Names[0].Name)
						}

						return false
					})
				}
			}
		}
	}

	return methods, nil
}
