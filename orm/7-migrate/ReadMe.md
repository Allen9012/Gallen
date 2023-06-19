数据库 Migrate 一直是数据库运维人员最为头痛的问题，如果仅仅是一张表增删字段还比较容易，那如果涉及到外键等复杂的关联关系，数据库的迁移就会变得非常困难。

GeeORM 的 Migrate 操作仅针对最为简单的场景，即支持字段的新增与删除，不支持字段类型变更。

1.1 新增字段
ALTER TABLE table_name ADD COLUMN col_name, col_type;
1.2 删除字段
对于 SQLite 来说，删除字段并不像新增字段那么容易，一个比较可行的方法需要执行下列几个步骤：
第一步：从 old_table 中挑选需要保留的字段到 new_table 中。
第二步：删除 old_table。
第三步：重命名 new_table 为 old_table。

2 GeeORM 实现 Migrate
按照原生的 SQL 命令，利用之前实现的事务，在 geeorm.go 中实现 Migrate 方法。

difference 用来计算前后两个字段切片的差集。新表 - 旧表 = 新增字段，旧表 - 新表 = 删除字段。
使用 ALTER 语句新增字段。
使用创建新表并重命名的方式删除字段。