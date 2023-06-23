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

## 压缩器

TinyRPC的压缩器代码部分很短，RawCompressor、GzipCompressor、SnappyCompressor、ZlibCompressor压缩器均实现了Compressor 接口

## 实现ClientCodec
其中 compressor  表示压缩类型，serializer 表示使用的序列化器，response 是响应的头部，mutex 是用于保护 pending 的互斥锁

## 实现serverCodec
TinyRPC在codec层还需要实现net/rpc的ServerCodec接口

ServerCodec 的接口和 ClientCodec 接口十分类似
其中 ServerCodec 接口包括写响应、读请求头部和读请求体，我们建立一个 serverCodec 的结构体用来实现 ServerCodec 接口

## TinyRPC的Server
TinyRPC的服务端非常简单，把标准库 net/rpc 的 Server 结构包装了一层，其中 ServeCodec 使用的是TinyRPC的编解码器

## TinyRPC的Client
TinyRPC的客户端也很简单，把标准库 net/rpc 的 Client 结构包s装了一层，其中 ClientCodec 使用的是TinyRPC的编解码器

## 压缩方式对比
Zlib、Gzip和Snappy都是常见的压缩算法和压缩格式，它们在压缩方式和性能方面有一些区别：

Zlib压缩方式：

Zlib是一个通用的数据压缩库，使用DEFLATE算法进行压缩。
Zlib压缩算法具有相对较高的压缩比，适用于需要高度压缩的数据，但压缩和解压缩速度可能较慢。
Gzip压缩方式：

Gzip是一种文件压缩格式，使用DEFLATE算法对数据进行压缩，并在压缩数据前添加了文件头和校验和等信息。
Gzip常用于文件压缩和网络传输，可以将多个文件打包成单个Gzip文件，并保留文件的元数据。
Gzip在压缩效率和压缩速度之间提供了一个良好的平衡，通常比Zlib稍微快一些。
Snappy压缩方式：

Snappy是一种快速的压缩/解压缩算法，旨在提供高速的数据压缩和解压缩性能。
Snappy的压缩比相对较低，但具有非常快速的压缩和解压缩速度，适用于对速度要求较高的场景，如实时数据传输和处理。
总体而言，Zlib和Gzip提供了更高的压缩比，适用于对压缩效率较为重视的情况。Snappy则注重压缩和解压缩的速度，适用于对速度要求较高的场景。选择使用哪种压缩方式取决于具体的应用需求和对压缩比和性能的权衡考虑。

## go-proto生成
执行脚本：
proto的生成
o install github.com/golang/protobuf/protoc-gen-go
go install github.com/zehuamama/tinyrpc/protoc-gen-tinyrpc
protoc --tinyrpc_out=. --go_out=. .\arith.proto