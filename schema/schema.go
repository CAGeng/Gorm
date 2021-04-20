package schema
//这部分实现 ORM 框架中最为核心的转换——对象(object)和表(table)的转换
import (
	"Gorm/dialect"
	"go/ast"
	"reflect"
)
//field表示数据库的一列
type Field struct {
	Name string
	Type string
	Tag string
}

//schema表示数据库的一张表
type Schema struct {
	Model interface{}
	Name string
	Fields []*Field
	FieldNames []string
	fieldMap map[string]*Field
}

func (schema *Schema)GetField(name string) *Field{
	return schema.fieldMap[name]
}

// Values return the values of dest's member variables
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

type ITableName interface {
	TableName() string
}

//将任意的对象解析为Schema函数
//TypeOf() 和 ValueOf() 是 reflect 包最为基本也是最重要的 2 个方法，
//分别用来返回入参的类型和值。因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例。
func Parse(dest interface{}, d dialect.Dialect) *Schema{
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string] *Field),
	}

	for i := 0; i < modelType.NumField(); i++{
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name){
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//p.Tag 即额外的约束条件
			if v,ok := p.Tag.Lookup("Gorm");ok{
				field.Tag = v
			}
			schema.Fields = append(schema.Fields,field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field

		}
	}
	return schema
}
