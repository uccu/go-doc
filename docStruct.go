package doc

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/uccu/go-stringify"
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

type TypeSpecWithKey struct {
	Key  string
	Tags *structtag.Tags
	*TypeSpec
}

func (s *TypeSpecWithKey) isPublic() bool {
	if s.Key == "" {
		return true
	}
	strArry := []byte(s.Key)
	return strArry[0] >= 65 && strArry[0] <= 90
}

type TypeSpec struct {
	Name     string
	TypeName string
	Kind     reflect.Kind
	Type     Type
	Value    []*TypeSpecWithKey
	pkg      *Pkg
	file     string
	Doc      []string
	Comment  string
}

var NilTypeSpec = &TypeSpec{
	Kind: reflect.Invalid,
	Type: NilType,
}
var NilTypeSpecWithKey = &TypeSpecWithKey{
	TypeSpec: NilTypeSpec,
}

type Struct interface {
	GetStru(string) *TypeSpec
}

func (ts *TypeSpec) GetStru(name string) *TypeSpec {
	if ts.Type == StructType {
		for _, field := range ts.Value {
			if field.TypeName == name {
				return field.TypeSpec
			}
		}
	}
	return nil
}

func ParseTypeSpec(t *ast.TypeSpec, pkg *Pkg, file string) *TypeSpec {
	ts := ParseType(t.Type, pkg, file)
	if ts == nil {
		return ts
	}
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

func parseField(f *ast.Field, pkg *Pkg, file string) *TypeSpec {
	ts := ParseType(f.Type, pkg, file)
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
	return ts
}

func ParseType(t ast.Expr, pkg *Pkg, file string) *TypeSpec {
	reflectType := reflect.TypeOf(t).Elem()

	typeSpec := &TypeSpec{
		file: file,
		pkg:  pkg,
		Kind: reflectType.Kind(),
	}

	reflectTypeName := reflectType.Name()

	switch reflectTypeName {
	case "StructType":
		s, ok := t.(*ast.StructType)
		if !ok {
			return nil
		}

		list := []*TypeSpecWithKey{}

		for _, f := range s.Fields.List {
			field := parseField(f, pkg, file)
			if field == nil {
				continue
			}
			var key string
			if f.Names != nil {
				key = f.Names[0].Name
			}
			t := &TypeSpecWithKey{Key: key, TypeSpec: field}

			if !t.isPublic() {
				continue
			}

			if f.Tag != nil {
				var err error
				t.Tags, err = structtag.Parse(strings.Trim(f.Tag.Value, "`"))
				if err != nil {
					return nil
				}
			}
			list = append(list, t)
		}

		typeSpec.TypeName = "object"
		typeSpec.Type = StructType
		typeSpec.Value = list
		return typeSpec
	case "ArrayType":
		arr, ok := t.(*ast.ArrayType)
		if !ok {
			return nil
		}
		typeSpec.Type = SliceType

		field := ParseType(arr.Elt, pkg, file)
		if field == nil {
			return nil
		}
		typeSpec.Value = []*TypeSpecWithKey{{TypeSpec: field}}
		typeSpec.TypeName = "array"
		return typeSpec
	case "MapType":
		arr, ok := t.(*ast.MapType)
		if !ok {
			return nil
		}
		typeSpec.Type = MapType
		key := ParseType(arr.Key, pkg, file)
		val := ParseType(arr.Value, pkg, file)
		if key == nil || val == nil {
			return nil
		}
		typeSpec.Value = []*TypeSpecWithKey{{TypeSpec: key}, {TypeSpec: val}}
		typeSpec.TypeName = "map"
		return typeSpec
	case "InterfaceType":
		_, ok := t.(*ast.InterfaceType)
		if !ok {
			return nil
		}
		typeSpec.Type = InterfaceType
		typeSpec.TypeName = "any"
		return typeSpec
	case "StarExpr":
		return ParseType(t.(*ast.StarExpr).X, pkg, file)
	case "SelectorExpr":
		expr := t.(*ast.SelectorExpr)
		ident := expr.X.(*ast.Ident)
		typeSpec.TypeName = ident.Name + "." + expr.Sel.Name
		typeSpec.Type = TypeType
		return typeSpec
	case "Ident":
		ident, ok := t.(*ast.Ident)
		if !ok {
			return nil
		}

		if ident.Obj == nil {
			typeSpec.TypeName = ident.Name
			kind, ok := nameKinds[ident.Name]
			if !ok {
				return nil
			}
			typeSpec.Kind = kind
			typeSpec.Type = nameTypes[ident.Name]
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

	typeSpec.Type = CustomType
	typeSpec.Kind = reflect.Struct
	typeSpec.TypeName = reflectTypeName
	return typeSpec

}

func parseTypeType(typeName string, pkg *Pkg, file string) *TypeSpecWithKey {

	paths := stringify.ToStringSlice(typeName, ".")
	if len(paths) == 0 {
		return nil
	}

	ts := pkg.GetStru(paths[0])
	if len(paths) == 1 {
		if ts == nil {
			return nil
		}
		return &TypeSpecWithKey{TypeSpec: ts}
	}

	pkg = pkg.GetPkg(file, paths[0])
	if pkg == nil {
		return nil
	}

	ts = pkg.GetStru(paths[1])
	if ts == nil {
		return nil
	}
	return &TypeSpecWithKey{TypeSpec: ts}
}
