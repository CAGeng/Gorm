package dialect
import "reflect"
//ORM 框架往往需要兼容多种数据库，因此我们需要将差异的这一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。这部分代码称之为 dialect.
type Dialect interface {
	//用于将 Go 语言的类型转换为该数据库的数据类型。
	DataTypeOf(typ reflect.Value) string
	//返回某个表是否存在的 SQL 语句，参数是表名(table)。
	TableExistSQL(tableName string) (string,[]interface{})
}

var dialectMap = map[string]Dialect{}

func RegisterDialect(name string,dialect Dialect){
	dialectMap[name] = dialect
}

func GetDialect(name string)(dialect Dialect, ok bool){
	dialect, ok = dialectMap[name]
	return
}