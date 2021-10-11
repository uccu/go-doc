package doc

import (
	"reflect"

	"github.com/fatih/structtag"
)

type Struct interface {
	GetStru(string) *DocStruct
}

type DocStruct struct {
	Name      string
	FieldList []*Field
}

type Field struct {
	Tags     *structtag.Tags
	TypeName string
	Kind     reflect.Kind
	Name     string
	Value    interface{}
}

func GetDocStruct() *DocStruct {
	return &DocStruct{}
}

func (doc *DocStruct) GetStru(name string) *DocStruct {
	for _, field := range doc.FieldList {
		if field.Name == name && field.Kind == reflect.Struct {
			return field.Value.(*DocStruct)
		}
	}
	return nil
}

// todo
