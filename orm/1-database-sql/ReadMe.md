开发一个框架/库并不容易，详细的日志能够帮助我们快速地定位问题。因此，在写核心代码之前，我们先用几十行代码实现一个简单的 log 库。

这个简易的 log 库具备以下特性：

支持日志分级（Info、Error、Disabled 三级）。
不同层级日志显示时使用不同的颜色区分。
显示打印日志代码对应的文件名和行号。

[info ] 颜色为蓝色，[error] 为红色。
使用 log.Lshortfile 支持显示文件名和代码行号。
暴露 Error，Errorf，Info，Infof 4个方法。

这一部分的实现非常简单，三个层级声明为三个常量，通过控制 Output，来控制日志是否打印。
如果设置为 ErrorLevel，infoLog 的输出会被定向到 ioutil.Discard，即不打印该日志。

封装有 2 个目的，一是统一打印日志（包括 执行的SQL 语句和错误日志）。接下来呢，封装 Exec()、Query() 和 QueryRow() 三个原生方法。
二是执行完成后，清空 (s *Session).sql 和 (s *Session).sqlVars 两个变量。这样 Session 可以复用，开启一次会话，可以执行多次 SQL。

## 核心结构 Engine

Engine 的逻辑非常简单，最重要的方法是 NewEngine，NewEngine 主要做了两件事。

连接数据库，返回 *sql.DB。
调用 db.Ping()，检查数据库是否能够正常连接。
另外呢，提供了 Engine 提供了 NewSession() 方法，这样可以通过 Engine 实例创建会话，进而与数据库进行交互了。到这一步，整个 GeeORM 的框架雏形已经出来了。