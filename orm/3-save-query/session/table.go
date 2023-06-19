/**
  @author: Allen
  @since: 2023/4/22
  @desc: //TODO
**/
package session

import (
	"fmt"
	"github/Allen9012/geeorm/log"
	"github/Allen9012/geeorm/schema"
	"reflect"
	"strings"
)

//
// Model
//  @Description: 用于给 refTable 赋值
//  @receiver s
//  @param value
//  @return *Session
//
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

//
// CreateTable
//  @Description: 创建表
//  @receiver s
//  @return error
//
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		if field.Tag == "PRIMARY KEY" {
			field.Tag = fmt.Sprintf("NOT NULL " + "PRIMARY KEY")
		}
		if field.Type == "varchar" {
			//todo 需要修改
			field.Type = fmt.Sprintf("varchar" + "(255)")
		}
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s  (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
