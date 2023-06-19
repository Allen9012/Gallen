通过链式(chain)操作，支持查询条件(where, order by, limit 等)的叠加。
实现记录的更新(update)、删除(delete)和统计(count)功能。

1 支持 Update、Delete 和 Count

clause 负责构造 SQL 语句，如果需要增加对更新(update)、删除(delete)和统计(count)功能的支持，第一步自然是在 clause 中实现 update、delete 和 count 子句的生成器。
1.1 子句生成器
第一步：在原来的基础上，新增 UPDATE、DELETE、COUNT 三个 Type 类型的枚举值。

_update 设计入参是2个，第一个参数是表名(table)，第二个参数是 map 类型，表示待更新的键值对。
_delete 只有一个入参，即表名。
_count 只有一个入参，即表名，并复用了 _select 生成器。
1.2 Update方法
子句的 generator 已经准备好了，接下来和 Insert、Find 等方法一样，在 session/record.go 中按照一定顺序拼接 SQL 语句并调用就可以了。

Update 方法比较特别的一点在于，Update 接受 2 种入参，平铺开来的键值对和 map 类型的键值对。
因为 generator 接受的参数是 map 类型的键值对，因此 Update 方法会动态地判断传入参数的类型，如果是不是 map 类型，则会自动转换。

1.3 Delete方法
1.4 Count 方法
2 链式调用(chain)

目标是
s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
var users []User
s.Where("Age > 18").Limit(3).Find(&users)

3 First 只返回一条记录