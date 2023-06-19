/**
  @author: Allen
  @since: 2023/4/22
  @desc: //TODO
**/
package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	// mysql没有bool
	//case reflect.Bool:
	//	return "bool"
	case reflect.Int8, reflect.Uint8:
		return "TINYINT"
	case reflect.Int16, reflect.Uint16:
		return "SMALLINT"
	case reflect.Int, reflect.Int32,
		reflect.Uint, reflect.Uint32, reflect.Uintptr:
		return "INT"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "varchar"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT TABLE_NAME from information_schema.tables where TABLE_SCHEMA   ='geeorm' and TABLE_NAME = ?", args
}
