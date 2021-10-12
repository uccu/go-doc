package doc

import (
	"encoding/json"
	"go/ast"
	"os"
	"strings"
)

func GetApis(pacakges ...string) []*DocApi {

	apis := []*DocApi{}

	for _, p := range pacakges {
		pkg := GetPkg(p)
		for _, f := range pkg.pkg.Syntax {
			for _, f := range f.Scope.Objects {
				if f.Kind == ast.Fun {

					funcDecl, ok := f.Decl.(*ast.FuncDecl)
					if !ok {
						continue
					}
					api := NewDocApi(funcDecl.Doc, pkg)
					if api != nil {
						apis = append(apis, api)
					}
				}
			}
		}
	}
	return apis
}

func Export(dir string, pacakges ...string) error {

	dir = strings.TrimRight(dir, "/\\")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	apis := GetApis(pacakges...)

	js, err := json.Marshal(apis)
	if err != nil {
		return err
	}
	dir += "/doc.json"

	os.WriteFile(dir, js, os.ModePerm)
	return nil
}
