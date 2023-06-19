/**
  @author: Allen
  @since: 2023/5/14
  @desc: //TODO
**/
package clause

import (
	"strings"
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Type is the type of Clause
type Type int

// Support types for Clause
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

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
// Build generate the final SQL and SQLVars
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
