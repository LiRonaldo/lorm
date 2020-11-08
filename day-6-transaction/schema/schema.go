package schema

import (
	"go/ast"
	"lorm/day-6-transaction/dialect"
	"reflect"
)

/**
表名(table name) —— 结构体名(struct name)
字段名和字段类型 —— 成员变量和类型。
额外的约束条件(例如非空、主键等) —— 成员变量的Tag（Go 语言通过 Tag 实现，Java、Python 等语言通过注解实现）
*/
type Field struct {
	Name string
	Type string
	Tag  string
}
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.FieldMap[name]
}
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	//相当于java的反射，获得某个类。可以通过方法获得某个类的属性。
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("lorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.FieldMap[p.Name] = field
		}
	}
	return schema
}

//将元素平铺，比如user 变成 u1，("Tom", 18)
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValues := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValues.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
