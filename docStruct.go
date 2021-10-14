package doc

import (
	"errors"
	"go/ast"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/uccu/go-stringify"
	"golang.org/x/tools/go/packages"
)

type Type int

const (
	NilType Type = 0 + iota
	BoolType
	IntType
	UintType
	FloatType
	StringType
	InterfaceType
	StructType
	SliceType
	MapType
	TypeType
	CustomType
)

var nameKinds = map[string]reflect.Kind{
	"bool":    reflect.Bool,
	"int":     reflect.Int,
	"int8":    reflect.Int8,
	"int16":   reflect.Int16,
	"int32":   reflect.Int32,
	"int64":   reflect.Int64,
	"uint":    reflect.Uint,
	"uint8":   reflect.Uint8,
	"uint16":  reflect.Uint16,
	"uint32":  reflect.Uint32,
	"uint64":  reflect.Uint64,
	"uintptr": reflect.Uintptr,
	"float32": reflect.Float32,
	"float64": reflect.Float64,
	"string":  reflect.String,
	"byte":    reflect.Uint8,
	"rune":    reflect.Int32,
}

var nameTypes = map[string]Type{
	"bool":    BoolType,
	"int":     IntType,
	"int8":    IntType,
	"int16":   IntType,
	"int32":   IntType,
	"int64":   IntType,
	"uint":    UintType,
	"uint8":   UintType,
	"uint16":  UintType,
	"uint32":  UintType,
	"uint64":  UintType,
	"uintptr": UintType,
	"float32": FloatType,
	"float64": FloatType,
	"string":  StringType,
	"byte":    UintType,
	"rune":    IntType,
}

type TypeSpec struct {
	TypeName string
	Kind     reflect.Kind
	Type     Type
	Value    []*TypeSpec
	Tags     *structtag.Tags
	Name     string
	pkg      *Pkg
	Doc      []string
	Comment  string
}

var NilTypeSpec = &TypeSpec{
	Kind: reflect.Invalid,
	Type: NilType,
}

type Struct interface {
	GetStru(string) *TypeSpec
}

func (ts *TypeSpec) GetStru(name string) *TypeSpec {
	if ts.Type == StructType {
		for _, field := range ts.Value {
			if field.Name == name {
				return field
			}
		}
	}
	return nil
}

func ParseTypeSpec(t *ast.TypeSpec, pkg *Pkg) *TypeSpec {
	ts := ParseType(t.Type, pkg)
	if ts == nil {
		return ts
	}
	ts.Name = t.Name.Name
	if t.Comment != nil {
		ts.Comment = strings.Trim(t.Comment.List[0].Text, "/ ")
	}
	if t.Doc != nil {
		list := []string{}
		for _, c := range t.Doc.List {
			list = append(list, strings.Trim(c.Text, "/ "))
		}
		ts.Doc = list
	}
	return ts
}

func parseField(f *ast.Field, pkg *Pkg) *TypeSpec {
	ts := ParseType(f.Type, pkg)
	if ts == nil {
		return ts
	}
	if f.Comment != nil {
		ts.Comment = strings.Trim(f.Comment.List[0].Text, "/ ")
	}
	if f.Doc != nil {
		list := []string{}
		for _, c := range f.Doc.List {
			list = append(list, strings.Trim(c.Text, "/ "))
		}
		ts.Doc = list
	}
	if f.Tag != nil {
		var err error
		ts.Tags, err = structtag.Parse(strings.Trim(f.Tag.Value, "`"))
		if err != nil {
			return nil
		}
	}

	ts.Name = ts.TypeName
	if len(f.Names) > 0 {
		ts.Name = f.Names[0].Name
	}

	return ts
}

func ParseType(t ast.Expr, pkg *Pkg) *TypeSpec {
	reflectType := reflect.TypeOf(t).Elem()

	typeSpec := &TypeSpec{
		pkg:      pkg,
		TypeName: reflectType.Name(),
		Kind:     reflectType.Kind(),
	}

	switch typeSpec.TypeName {
	case "StructType":
		s, ok := t.(*ast.StructType)
		if !ok {
			return nil
		}

		list := []*TypeSpec{}

		for _, f := range s.Fields.List {
			field := parseField(f, pkg)
			if field == nil {
				continue
			}
			list = append(list, field)
		}

		typeSpec.Type = StructType
		typeSpec.Value = list
		return typeSpec
	case "ArrayType":
		arr, ok := t.(*ast.ArrayType)
		if !ok {
			return nil
		}
		typeSpec.Type = SliceType

		field := ParseType(arr.Elt, pkg)
		if field == nil {
			return nil
		}
		typeSpec.Value = []*TypeSpec{field}
		typeSpec.TypeName = "array"
		return typeSpec
	case "MapType":
		arr, ok := t.(*ast.MapType)
		if !ok {
			return nil
		}
		typeSpec.Type = MapType
		key := ParseType(arr.Key, pkg)
		val := ParseType(arr.Value, pkg)
		if key == nil || val == nil {
			return nil
		}
		typeSpec.Value = []*TypeSpec{key, val}
		typeSpec.TypeName = "map"
		return typeSpec
	case "InterfaceType":
		_, ok := t.(*ast.InterfaceType)
		if !ok {
			return nil
		}
		typeSpec.TypeName = "any"
		typeSpec.Type = InterfaceType
		return typeSpec
	case "StarExpr":
		return ParseType(t.(*ast.StarExpr).X, pkg)
	case "Time":
		typeSpec.Type = CustomType
		typeSpec.Kind = reflect.Struct
		return typeSpec

	case "Ident":
		ident, ok := t.(*ast.Ident)
		if !ok {
			return nil
		}

		if ident.Obj == nil {
			typeSpec.TypeName = ident.Name
			kind, ok := nameKinds[typeSpec.TypeName]
			if !ok {
				return nil
			}
			typeSpec.Kind = kind
			typeSpec.Type = nameTypes[typeSpec.TypeName]
			return typeSpec
		}

		if ident.Obj.Kind != 3 {
			return nil
		}

		ts, ok := ident.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			return nil
		}

		typeSpec.TypeName = ts.Name.Name
		typeSpec.Type = TypeType
		return typeSpec
	}

	return nil

}

func parseTypeType(typeName string, pkg *Pkg) *TypeSpec {

	paths := stringify.ToStringSlice(typeName, ".")
	if len(paths) == 0 {
		return NilTypeSpec
	}

	ts := pkg.GetStru(paths[0])
	if len(paths) == 1 {
		if ts == nil {
			return NilTypeSpec
		}
		return ts
	}

	pkg = pkg.GetPkg(paths[0])
	if pkg == nil {
		return NilTypeSpec
	}

	ts = pkg.GetStru(paths[1])
	if ts == nil {
		return NilTypeSpec
	}
	return ts
}

func GetType(dir string) *Pkg {
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
