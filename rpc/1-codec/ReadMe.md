# 消息的序列化与反序列化
客户端发送的请求包括服务名 Arith，方法名 Multiply，参数 args 三个，服务端的响应包括错误 error，返回值 reply 2 个。
我们将请求和响应中的参数和返回值抽象为 body，剩余的信息放在 header 中，那么就可以抽象出数据结构 Header：

codec/codec.go
ServiceMethod 是服务名和方法名，通常与 Go 语言中的结构体和方法相映射。
Seq 是请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求。
Error 是错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中。

我们将和消息编解码相关的代码都放到 codec 子目录中，在此之前，还需要在根目录下使用 go mod init geerpc 初始化项目，方便后续子 package 之间的引用。

进一步，抽象出对消息体进行编解码的接口 Codec，抽象出接口是为了实现不同的 Codec

紧接着，抽象出 Codec 的构造函数，客户端和服务端可以通过 Codec 的 Type 得到构造函数，从而创建 Codec 实例。
这部分代码和工厂模式类似，与工厂模式不同的是，返回的是构造函数，而非实例。

我们定义了 2 种 Codec，Gob 和 Json，但是实际代码中只实现了 Gob 一种，事实上，2 者的实现非常接近，甚至只需要把 gob 换成 json 即可。

首先定义 GobCodec 结构体，这个结构体由四部分构成，conn 是由构建函数传入，通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例，
dec 和 enc 对应 gob 的 Decoder 和 Encoder，buf 是为了防止阻塞而创建的带缓冲的 Writer，一般这么做能提升性能。

接着实现 ReadHeader、ReadBody、Write 和 Close 方法。

# 通信过程

客户端与服务端的通信需要协商一些内容，例如 HTTP 报文，分为 header 和 body 2 部分，body 的格式和长度通过 header 中的 Content-Type 和 Content-Length 指定，
服务端通过解析 header 就能够知道如何从 body 中读取需要的信息。
对于 RPC 协议来说，这部分协商是需要自主设计的。
为了提升性能，一般在报文的最开始会规划固定的字节，来协商相关的信息。
比如第1个字节用来表示序列化方式，第2个字节表示压缩方式，
第3-6字节表示 header 的长度，7-10 字节表示 body 的长度。

对于 GeeRPC 来说，目前需要协商的唯一一项内容是消息的编解码方式。我们将这部分信息，放到结构体 Option 中承载。
目前，已经进入到服务端的实现阶段了。

codec/server.go
一般来说，涉及协议协商的这部分信息，需要设计固定的字节来传输的。
但是为了实现上更简单，GeeRPC 客户端固定采用 JSON 编码 Option，
后续的 header 和 body 的编码方式由 Option 中的 CodeType 指定，
服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType 解码剩余的内容。
即报文将以这样的形式发送：
```text
Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
```