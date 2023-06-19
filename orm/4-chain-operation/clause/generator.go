/**
  @author: Allen
  @since: 2023/4/29
  @desc: //TODO
**/
package clause

import (
	"fmt"
	"github/Allen9012/geeorm/log"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})
type Type int

// 保存实际语句Type 的 map
var generators map[Type]generator

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	sql.WriteString("VALUES ")
	var vars []interface{}
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		// (?)
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

//
// _update
//  @Description: 拼接update语句
//  @param values
//  @return string
//  @return []interface{}
//
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	// 提取参数
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	// eg: user = ?, age = ?
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	// UPDATE $tableName SET $field1 = ?, $field2 = ?, ...
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

//
// _delete
//  @Description: 拼接delete语句
//  @param values
//  @return string
//  @return []interface{}
//
func _delete(values ...interface{}) (string, []interface{}) {
	// DELETE FROM $tableName
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

//
// _count
//  @Description: 拼接count语句
//  @param values
//  @return string
//  @return []interface{}
//
func _count(values ...interface{}) (string, []interface{}) {
	// SELECT COUNT(*) FROM $tableName
	return _select(values[0], []string{"count(*)"})
}

//
// Set
//  @Description: Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句。
//  @receiver c
//  @param name
//  @param vars
//
func (c *Clause) Set(name Type, vars ...interface{}) {
	// 需要注意先得初始化map
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

//
// Build
//  @Description: Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。
//  @receiver c
//  @param orders
//  @return string
//  @return []interface{}
//
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		} else {
			log.Error("no such order, order: ", order, ok)
		}
	}
	return strings.Join(sqls, " "), vars
}
