package doc

import (
	"bufio"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/uccu/go-stringify"
)

var pkgs map[string]*Pkg
var pkgMod string

func init() {
	pkgs = make(map[string]*Pkg)
}

type Pkg struct {
	Dir  string
	Name string
	pkg  *ast.Package
	pkgs map[string]*Pkg
	stru map[string]*TypeSpec
}

func setMod() {
	if pkgMod == "" {
		base, _ := os.Getwd()
		f, err := os.Open(base + "/go.mod")
		if err != nil {
			panic(err)
		}
		rd := bufio.NewScanner(f)
		if rd.Scan() {
			pkgMod = string(([]byte(rd.Text()))[7:])
		} else {
			panic(errors.New("no mod file"))
		}
		f.Close()
	}
}

func GetPkg(pkgName string) *Pkg {

	setMod()

	if strings.Index(pkgName, pkgMod) != 0 {
		return nil
	}

	pkgName = strings.Replace(pkgName, pkgMod, "", 1)
	base, _ := os.Getwd()
	dir := base + pkgName

	if pkg, ok := pkgs[dir]; ok {
		return pkg
	}
	pkgMap, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	for name, pkg := range pkgMap {
		slp := stringify.ToStringSlice(name, "_")
		if slp[len(slp)-1] == "test" {
			continue
		}

		pkg := &Pkg{
			pkg:  pkg,
			Name: pkg.Name,
		}

		pkgs[dir] = pkg
		return pkg
	}

	return nil
}

func (pkg *Pkg) SetPkgs() *Pkg {
	if pkg.pkgs != nil {
		return pkg
	}
	pkg.pkgs = make(map[string]*Pkg)
	for _, f := range pkg.pkg.Files {
		for _, p := range f.Imports {
			pkgName := strings.Trim(p.Path.Value, "\"")
			mpkg := GetPkg(pkgName)
			if mpkg == nil {
				continue
			}
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
	for _, f := range pkg.pkg.Files {
		for _, p := range f.Scope.Objects {
			if p.Kind != ast.Typ {
				continue
			}
			typeSpec, _ := p.Decl.(*ast.TypeSpec)
			pkg.stru[typeSpec.Name.Name] = ParseTypeSpec(typeSpec, pkg)
			if pkg.stru[typeSpec.Name.Name] == nil {
				continue
			}
			pkg.stru[typeSpec.Name.Name].Name = typeSpec.Name.Name
		}
	}
	return pkg
}

func (pkg *Pkg) GetStru(name string) *TypeSpec {
	s, ok := pkg.SetStru().stru[name]
	if !ok || s == nil {
		return nil
	}
	return s
}

func GetApis(pacakges ...string) []*DocApi {
	apis := []*DocApi{}
	for _, p := range pacakges {
		pkg := GetPkg(p)
		if pkg == nil {
			continue
		}
		for _, f := range pkg.pkg.Files {
			for _, f := range f.Decls {
				funcDecl, ok := f.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if funcDecl.Doc == nil {
					continue
				}
				api := NewDocApi(funcDecl.Doc, pkg)
				if api != nil {
					apis = append(apis, api)
				}
			}
		}
	}
	return apis
}
