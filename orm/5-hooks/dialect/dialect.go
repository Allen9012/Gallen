/**
  @author: Allen
  @since: 2023/4/22
  @desc: //TODO
**/
package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf 用于将 Go 语言的类型转换为该数据库的数据类型。
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL 返回某个表是否存在的 SQL 语句，参数是表名(table)。
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect register a dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect Get a dialect from map
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
