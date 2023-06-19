TinyRPC 是基于Go语言标准库 net/rpc 扩展的远程过程调用框架，它具有以下特性：

基于TCP传输层协议
支持多种压缩格式：gzip、snappy、zlib；
基于二进制的 Protocol Buffer 序列化协议：具有协议编码小及高扩展性和跨平台性；
支持生成工具：TinyRPC提供的 protoc-gen-tinyrpc 插件可以帮助开发者快速定义自己的服务；
支持自定义序列化器
TinyRPC 的源代码仅有一千行左右，通过学习 TinyRPC ，开发者可以得到以下收获：

代码简洁规范
涵盖大多数 Go 语言基础用法和高级特性
单元测试编写技巧
TCP流中处理数据包的技巧
RPC框架的设计理念

## 基于TCP

在TinyRPC中，请求消息由TinyRPC客户端的应用程序发出，在TCP的字节流中，请求消息分为三部分：

+ 由可变长量编码的 **uint 类型**用来标识请求头的长度；
+ 基于自定义协议编码的请求头部信息
+ 基于 **Protocol Buffer** 协议编码的请求体，见图所示：

Request：

![image-20230618141255192](https://cloudmage.oss-cn-shanghai.aliyuncs.com/img/202306181412300.png)

1. Headerlength
2. Header
3. Body

Response：

+ 由可变长量编码的 **uint 类型**用来标识响应头的长度；
+ 基于自定义协议编码的响应头部信息
+ 基于 **Protocol Buffer** 协议编码的响应体，见图所示：

![image-20230618141353651](https://cloudmage.oss-cn-shanghai.aliyuncs.com/img/202306181413825.png)

+ 其中ID为RPC调用的序号，以便在并发调用时，客户端根据响应的ID序号来判断RPC的调用结果；
+ Error message为调用时发生错误的消息，若该内容为空则表示未出现RPC调用错误；
+ 在请求I/O流中，请求体（Request Body）表示RPC的参数内容；而在响应I/O流中，响应体（Response Body）则表示RPC调用的结果，这些Body在TinyRPC中均采用 **Protocol Buffer** 协议编码。

可以看到比GeeRPC的设计更加负责和精细一点

### 实现方法

marshal

1. 加锁
2. 用idx记录下标
3. 按照头的定义写入协议

unmarshal

1. 加锁
2. recover异常
3. 读取协议内容写入结构体

### header/pool.go

为了减少创建请求头部对象 **RequestHeader** 和响应头部对象 **ResponseHeader** 的次数**，**我们通过为这两个结构体建立对象池，以便可以进行复用。

同时我们为 *RequestHeader* 和 *ResponseHeader* 都实现了ResetHeader方法，当每次使用完这些对象时，我们调用ResetHeader让结构体内容初始化，随后再把它们丢回对象池里。

> 1. **减少内存分配次数**：对象池允许重用已经分配的请求头部和响应头部对象，避免频繁地创建新的对象。这可以减少内存分配次数，降低内存管理的开销。
> 2. **提高性能**：由于对象池避免了创建新对象的开销，可以显著提高系统的性能。在高并发场景下，减少对象创建和销毁的操作可以提高系统的吞吐量和响应速度。
> 3. **避免垃圾回收压力**：通过复用对象，可以减少不必要的垃圾对象产生，从而减轻垃圾回收器的压力。这有助于系统的稳定性和可预测性。
> 4. **代码简洁性**：通过对象池和ResetHeader方法，可以将对象的初始化和复用逻辑封装在内部，使代码更加简洁和可读。开发人员可以专注于业务逻辑，而无需频繁处理对象的创建和销毁。

## IO操作

TinyRPC的IO操作函数在[codec/io.go](https://link.zhihu.com/?target=https%3A//github.com/zehuamama/tinyrpc/blob/main/codec/io.go)中，其中 *sendFrame* 函数会向IO中写入**uvarint**类型的 *size* ，表示要发送数据的长度，随后将该字节slice类型的数据 *data* 写入IO流中。

+ 若写入数据的长度为 0 ，此时*sendFrame* 函数会向IO流写入uvarint类型的 0 值；
+ 若写入数据的长度大于 0 ，此时*sendFrame* 函数会向IO流写入uvarint类型的 *len(data)* 值，随后将该字节串的数据 data 写入IO流中。

*recvFrame* 函数与*sendFrame* 函数类似，首先会向IO中读入**uvarint**类型的 *size* ，表示要接收数据的长度，随后将该从IO流中读取该 *size* 长度字节串。

注意，由于 codec 层会传入一个**bufio**类型的结构体，**bufio**类型实现了有缓冲的IO操作，以便减少IO在用户态与内核态拷贝的次数。

若 *recvFrame* 函数从IO流读取**uvarint**类型的 *size* 值大于0，随后 *recvFrame* 将该从IO流中读取该 *size* 长度字节串。
