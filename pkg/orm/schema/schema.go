package schema

import (
	"github.com/peterouob/golang_template/pkg/orm/dialect"
	"go/ast"
	"reflect"
)

// Field
//
//		type User struct{
//				Name string `orm:"prim"`
//				Age int
//		}
//	 transport this struct to schema like:
//	  CREATE TABLE `User` (`Name` text PRIMARY KEY, `Age` integer);
//
// /*
type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model     interface{}
	Name      string
	Fields    []*Field
	FieldName []string
	fieldMap  map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 這邊使用Indirect是因為假設使用Elem是因為有可能傳入struct
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	// NumField will be panic if modelType.Kind() is not a struct
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 非嵌入式struct且為public才要去將struct轉換成schema
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldName = append(schema.FieldName, field.Name)
			schema.fieldMap[field.Name] = field
		}
	}
	return schema
}
