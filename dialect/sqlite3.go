package dialect

import (
	"fmt"
	"reflect"
	"time"
)

//提供对sqlite3的支持，结构sqlite3要实现接口dialect的所有方法
type sqlite3 struct {}

//下面这句用来检查结构sqlite3是否实现接口dialect的所有方法
var _ Dialect = (*sqlite3)(nil)

func init(){
	RegisterDialect("sqlite3",&sqlite3{})
}

func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		//时间
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))

}

func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}  //巧妙的语法
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
