对象关系映射（Object Relational Mapping，简称ORM）是通过使用描述对象和数据库之间映射的元数据，将面向对象语言程序中的对象自动持久化到关系数据库中。


## 那对象和数据库是如何映射的呢？

数据库	面向对象的编程语言
表(table)	类(class/struct)
记录(record, row)	对象 (object)
字段(field, column)	对象属性(attribute)
```go
package main
type User struct {
Name string
Age  int
}

orm.CreateTable(&User{})
orm.Save(&User{"Tom", 18})
var users []User
orm.Find(&users)
```



CreateTable 方法需要从参数 &User{} 得到对应的结构体的名称 User 作为表名，成员变量 Name, Age 作为列名，同时还需要知道成员变量对应的类型。
Save 方法则需要知道每个成员变量的值。
Find 方法仅从传入的空切片 &[]User，得到对应的结构体名也就是表名 User，并从数据库中取到所有的记录，将其转换成 User 对象，添加到切片中。


这就面临了一个很重要的问题：如何根据任意类型的指针，得到其对应的结构体的信息。这涉及到了 Go 语言的反射机制(reflect)，通过反射，可以获取到对象对应的结构体名称，成员变量、方法等信息，例如：

typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type()
fmt.Println(typ.Name()) // Account

for i := 0; i < typ.NumField(); i++ {
field := typ.Field(i)
fmt.Println(field.Name) // Username Password
}

reflect.ValueOf() 获取指针对应的反射值。
reflect.Indirect() 获取指针指向的对象的反射值。
(reflect.Type).Name() 返回类名(字符串)。
(reflect.Type).Field(i) 获取第 i 个成员变量。

## 除了对象和表结构/记录的映射以外，设计 ORM 框架还需要关注什么问题呢？

1）MySQL，PostgreSQL，SQLite 等数据库的 SQL 语句是有区别的，ORM 框架如何在开发者不感知的情况下适配多种数据库？

2）如何对象的字段发生改变，数据库表结构能够自动更新，即是否支持数据库自动迁移(migrate)？

3）数据库支持的功能很多，例如事务(transaction)，ORM 框架能实现哪些？


表的创建、删除、迁移。
记录的增删查改，查询条件的链式操作。
单一主键的设置(primary key)。
钩子(在创建/更新/删除/查找之前或之后)
事务(transaction)。