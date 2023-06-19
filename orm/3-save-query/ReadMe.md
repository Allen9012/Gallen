查询语句的构成：
SELECT col1, col2, ...
FROM table_name
WHERE [ conditions ]
GROUP BY col1
HAVING [ conditions ]
由于复杂所以需要单独生成
首先在 clause/generator.go 中实现各个子句的生成规则。

generator.go

Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句。
Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。

首先为 Session 添加成员变量 clause
INSERT 对应的 SQL 语句一般是这样的：
INSERT INTO table_name(col1, col2, col3, ...) VALUES
(A1, A2, A3, ...),
(B1, B2, B3, ...),
...

在 ORM 框架中期望 Insert 的调用方式如下：
s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
u1 := &User{Name: "Tom", Age: 18}
u2 := &User{Name: "Sam", Age: 25}
s.Insert(u1, u2, ...)

也就是说，我们还需要一个步骤，根据数据库中列的顺序，从对象中找到对应的值，按顺序平铺。即 u1、u2 转换为 ("Tom", 18), ("Same", 25) 这样的格式。

后续所有构造 SQL 语句的方式都将与 Insert 中构造 SQL 语句的方式一致。分两步：

1）多次调用 clause.Set() 构造好每一个子句。
2）调用一次 clause.Build() 按照传入的顺序构造出最终的 SQL 语句。
构造完成后，调用 Raw().Exec() 方法执行。

期望Find的调用方式：
s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
var users []User
s.Find(&users);

Find 功能的难点和 Insert 恰好反了过来。Insert 需要将已经存在的对象的每一个字段的值平铺开来，而 Find 则是需要根据平铺开的字段的值构造出对象。同样，也需要用到反射(reflect)。

Find 的代码实现比较复杂，主要分为以下几步：

destSlice.Type().Elem() 获取切片的单个元素的类型 destType，使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，映射出表结构 RefTable()。
2）根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
3）遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
4）调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段。
5）将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中。