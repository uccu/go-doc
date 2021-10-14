package doc

import (
	"errors"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

var pkgs map[string]*Pkg

func init() {
	pkgs = make(map[string]*Pkg)
}

type Pkg struct {
	Name string
	pkg  *packages.Package
	pkgs map[string]*Pkg
	stru map[string]*TypeSpec
}

func GetPkg(dir string) *Pkg {
	dir = strings.Trim(dir, "\"")
	if pkg, ok := pkgs[dir]; ok {
		return pkg
	}

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgss, err := packages.Load(cfg, dir)
	if err != nil || len(pkgss) == 0 {
		panic(errors.New("package is not exist: " + dir))
	}
	pkg := &Pkg{
		pkg:  pkgss[0],
		Name: pkgss[0].Syntax[0].Name.Name,
	}
	pkgs[dir] = pkg
	return pkg
}

func (pkg *Pkg) SetPkgs() *Pkg {
	if pkg.pkgs != nil {
		return pkg
	}
	pkg.pkgs = make(map[string]*Pkg)
	for _, f := range pkg.pkg.Syntax {
		for _, p := range f.Imports {
			mpkg := GetPkg(p.Path.Value)
			name := mpkg.Name
			if p.Name != nil {
				name = p.Name.Name
			}
			pkg.pkgs[name] = mpkg
		}
	}
	return pkg
}

func (pkg *Pkg) GetPkg(name string) *Pkg {
	p, ok := pkg.SetPkgs().pkgs[name]
	if !ok {
		return nil
	}
	return p
}

func (pkg *Pkg) SetStru() *Pkg {
	if pkg.stru != nil {
		return pkg
	}

	pkg.stru = make(map[string]*TypeSpec)
	for _, f := range pkg.pkg.Syntax {
		for _, p := range f.Scope.Objects {
			if p.Kind != ast.Typ {
				continue
			}
			typeSpec, _ := p.Decl.(*ast.TypeSpec)
			pkg.stru[typeSpec.Name.Name] = ParseTypeSpec(typeSpec, pkg)
		}
	}
	return pkg
}

func (pkg *Pkg) GetStru(name string) *TypeSpec {
	s, ok := pkg.SetStru().stru[name]
	if !ok || s == nil {
		return nil
	}
	s.Name = name
	return s
}

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
