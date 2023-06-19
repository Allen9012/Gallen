/**
  @author: Allen
  @since: 2023/4/30
  @desc: //TODO
**/
package session

import (
	"errors"
	"github/Allen9012/geeorm/clause"
	"reflect"
)

//
// Insert
//  @Description: Insert one or more records into database.
// 1) 多次调用 clause.Set() 构造好每一个子句。
// 2) 调用一次 clause.Build() 按照传入的顺序构造出最终的 SQL 语句。
// 3) 构造完成后，调用 Raw().Exec() 方法执行。
//  @receiver s
//  @param values
//  @return int64
//  @return error
//
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		// insert into table fields
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	// insert into tables values
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

//
// Find
//  @Description: Find records by providing query conditions.
// 1) 通过传入的切片类型的指针，创建一个反射值，然后调用 reflect.Append() 方法在切片尾部追加新的反射值。
// 2) 通过反射值的 FieldByName() 方法获取到每一个字段的反射值，然后调用 Addr().Interface() 方法得到字段值的指针。
// 3) 将每一个字段值的指针追加到 values 中。
// 4) 调用 rows.Scan() 方法将该行的值依次赋值给 values 中的每一个字段值。
// 5) 将该行记录追加到切片中。
// 6) 返回 rows.Close() 的结果。
//  @receiver s
//  @param values
//  @return error
//
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	// select table field
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}

//
// Update
//  @Description: Update records by providing query conditions.
//  @receiver s
//  @param values
//  @return error
// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	// 如果是map类型不用处理，否则需要将kv转换为map
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	// update table set field1=?, field2=? where ...
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete
//  @Description: Delete records with where clause
//  @receiver s
//  @return int64
//  @return error
//
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	// delete from $tableName where ...
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

//
// Count
//  @Description: Count records with where clause
//  @receiver s
//  @return int64
//  @return error
//
func (s *Session) Count() (int64, error) {
	// select count(*) from $tableName where ...
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

//
// Limit
//  @Description: Limit adds limit condition to clause
//  @receiver s
//  @param num
//  @return *Session
//
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

//
// Where
//  @Description: Where adds where condition to clause
//  @receiver s
//  @param desc
//  @param args
//  @return *Session
//
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

//
// OrderBy
//  @Description: OrderBy adds order by condition to clause
//  @receiver s
//  @param desc
//  @return *Session
//
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

//
// First
//  @Description: First finds the first record that match given conditions, then assign it to value pointed by dest.
//	实现原理：根据传入的类型，利用反射构造切片，调用 Limit(1) 限制返回的行数，调用 Find 方法获取到查询结果。
//  @receiver s
//  @param value
//  @return error
//
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
