package doc

import (
	"errors"
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
	stru map[string]*DocStruct
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

	//todo
	return pkg
}

func (pkg *Pkg) GetStru(name string) *DocStruct {
	s, ok := pkg.SetStru().stru[name]
	if !ok {
		return nil
	}
	return s
}
