注册中心的位置如上图所示。注册中心的好处在于，客户端和服务端都只需要感知注册中心的存在，而无需感知对方的存在。更具体一些：

服务端启动后，向注册中心发送注册消息，注册中心得知该服务已经启动，处于可用状态。一般来说，服务端还需要定期向注册中心发送心跳，证明自己还活着。
客户端向注册中心询问，当前哪天服务是可用的，注册中心将可用的服务列表返回客户端。
客户端根据注册中心得到的服务列表，选择其中一个发起调用。
如果没有注册中心，就像 GeeRPC 第六天实现的一样，客户端需要硬编码服务端的地址，而且没有机制保证服务端是否处于可用状态。当然注册中心的功能还有很多，比如配置的动态同步、通知机制等。比较常用的注册中心有 etcd、zookeeper、consul，一般比较出名的微服务或者 RPC 框架，这些主流的注册中心都是支持的。

主流的注册中心 etcd、zookeeper 等功能强大，与这类注册中心的对接代码量是比较大的，需要实现的接口很多。GeeRPC 选择自己实现一个简单的支持心跳保活的注册中心。

GeeRegistry 的代码独立放置在子目录 registry 中。

首先定义 GeeRegistry 结构体，默认超时时间设置为 5 min，也就是说，任何注册的服务超过 5 min，即视为不可用状态。

为 GeeRegistry 实现添加服务实例和返回服务列表的方法。

putServer：添加服务实例，如果服务已经存在，则更新 start。
aliveServers：返回可用的服务列表，如果存在超时的服务，则删除。

为了实现上的简单，GeeRegistry 采用 HTTP 协议提供服务，且所有的有用信息都承载在 HTTP Header 中。

Get：返回所有可用的服务列表，通过自定义字段 X-Geerpc-Servers 承载。
Post：添加服务实例或发送心跳，通过自定义字段 X-Geerpc-Server 承载。

另外，提供 Heartbeat 方法，便于服务启动时定时向注册中心发送心跳，默认周期比注册中心设置的过期时间少 1 min。

GellenRegistryDiscovery 嵌套了 MultiServersDiscovery，很多能力可以复用。
registry 即注册中心的地址
timeout 服务列表的过期时间
lastUpdate 是代表最后从注册中心更新服务列表的时间，默认 10s 过期，即 10s 之后，需要从注册中心更新新的列表。
实现 Update 和 Refresh 方法，超时重新获取的逻辑在 Refresh 中实现：