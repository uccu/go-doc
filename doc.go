package doc

import "go/ast"

// todo
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
