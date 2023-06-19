/**
  @author: Allen
  @since: 2023/4/16
  @desc: //TODO
**/
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github/Allen9012/geeorm"
)

func main() {
	engine, _ := geeorm.NewEngine("mysql", "root:9012@tcp(localhost:3306)/geeorm")
	defer engine.Close()
	session := engine.NewSession()
	_, _ = session.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = session.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = session.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := session.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
