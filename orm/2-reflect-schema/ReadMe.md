SQL 语句中的类型和 Go 语言中的类型是不同的，例如Go 语言中的 int、int8、int16 等类型均对应 SQLite 中的 integer 类型。因此实现 ORM 映射的第一步，需要思考如何将 Go 语言的类型映射为数据库中的类型。
同时，不同数据库支持的数据类型也是有差异的，即使功能相同，在 SQL 语句的表达上也可能有差异。ORM 框架往往需要兼容多种数据库，因此我们需要将差异的这一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。这部分代码称之为 dialect。
在根目录下新建文件夹 dialect，并在 dialect 文件夹下新建文件 dialect.go，抽象出各个数据库差异的部分。

当然，不同数据库之间的差异远远不止这两个地方，随着 ORM 框架功能的增多，dialect 的实现也会逐渐丰富起来，同时框架的其他部分不会受到影响。
同时，声明了 RegisterDialect 和 GetDialect 两个方法用于注册和获取 dialect 实例。如果新增加对某个数据库的支持，那么调用 RegisterDialect 即可注册到全局。
接下来，在dialect 目录下新建文件 sqlite3.go 增加对 SQLite 的支持。

1. Dialect

Dialect 实现了一些特定的 SQL 语句的转换，接下来我们将要实现 ORM 框架中最为核心的转换——对象(object)和表(table)的转换。给定一个任意的对象，转换为关系型数据库中的表结构。

在数据库中创建一张表需要哪些要素呢？

表名(table name) —— 结构体名(struct name)
字段名和字段类型 —— 成员变量和类型。
额外的约束条件(例如非空、主键等) —— 成员变量的Tag（Go 语言通过 Tag 实现，Java、Python 等语言通过注解实现）

2. Schema

Field 包含 3 个成员变量，字段名 Name、类型 Type、和约束条件 Tag
Schema 主要包含被映射的对象 Model、表名 Name 和字段 Fields。
FieldNames 包含所有的字段名(列名)，fieldMap 记录字段名和 Field 的映射关系，方便之后直接使用，无需遍历 Fields。

TypeOf() 和 ValueOf() 是 reflect 包最为基本也是最重要的 2 个方法，分别用来返回入参的类型和值。因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例。
modelType.Name() 获取到结构体的名称作为表名。
NumField() 获取实例的字段的个数，然后通过下标获取到特定字段 p := modelType.Field(i)。
p.Name 即字段名，p.Type 即字段类型，通过 (Dialect).DataTypeOf() 转换为数据库的字段类型，p.Tag 即额外的约束条件。

3. Session 的核心功能是与数据库进行交互。因此，我们将数据库表的增/删操作实现在子包 session 中。在此之前，Session 的结构需要做一些调整。
   Session 成员变量新增 dialect 和 refTable
   构造函数 New 的参数改为 2 个，db 和 dialect。
   在文件夹 session 下新建 table.go 用于放置操作数据库表相关的代码。

Model() 方法用于给 refTable 赋值。解析操作是比较耗时的，因此将解析的结果保存在成员变量 refTable 中，即使 Model() 被调用多次，如果传入的结构体名称不发生变化，则不会更新 refTable 的值。
RefTable() 方法返回 refTable 的值，如果 refTable 未被赋值，则打印错误日志。

因为 Session 构造函数增加了对 dialect 的依赖，Engine 需要作一些细微的调整。

NewEngine 创建 Engine 实例时，获取 driver 对应的 dialect。
NewSession 创建 Session 实例时，传递 dialect 给构造函数 New。
至此，第二天的内容已经完成了，总结一下今天的成果：

1）为适配不同的数据库，映射数据类型和特定的 SQL 语句，创建 Dialect 层屏蔽数据库差异。
2）设计 Schema，利用反射(reflect)完成结构体和数据库表结构的映射，包括表名、字段名、字段类型、字段 tag 等。
3）构造创建(create)、删除(drop)、存在性(table exists) 的 SQL 语句完成数据库表的基本操作。