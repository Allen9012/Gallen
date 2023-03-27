对 net/rpc 而言，一个函数需要能够被远程调用，需要满足如下五个条件：

the method’s type is exported.
the method is exported.
the method has two arguments, both exported (or builtin) types.
the method’s second argument is a pointer.
the method has return type error.

为了支持异步调用，Call 结构体中添加了一个字段 Done，Done 的类型是 chan *Call，当调用结束时，会调用 call.done() 通知调用方。

client/client.go

Client 的字段比较复杂：

1. cc 是消息的编解码器，和服务端类似，用来序列化将要发送出去的请求，以及反序列化接收到的响应。 
2. sending 是一个互斥锁，和服务端类似，为了保证请求的有序发送，即防止出现多个请求报文混淆。 
3. header 是每个请求的消息头，header 只有在请求发送时才需要，而请求发送是互斥的，因此每个客户端只需要一个，声明在 Client 结构体中可以复用。 
4. seq 用于给发送的请求编号，每个请求拥有唯一编号。 
5. pending 存储未处理完的请求，键是编号，值是 Call 实例。 
6. closing 和 shutdown 任意一个值置为 true，则表示 Client 处于不可用的状态，但有些许的差别， 
7. closing 是用户主动关闭的，即调用 Close 方法，而 shutdown 置为 true 一般是有错误发生。

registerCall：将参数 call 添加到 client.pending 中，并更新 client.seq。
removeCall：根据 seq，从 client.pending 中移除对应的 call，并返回。
terminateCalls：服务端或客户端发生错误时调用，将 shutdown 设置为 true，且将错误信息通知所有 pending 状态的 call。
对一个客户端端来说，接收响应、发送请求是最重要的 2 个功能。那么首先实现接收功能，接收到的响应有三种情况：

call 不存在，可能是请求没有发送完整，或者因为其他原因被取消，但是服务端仍旧处理了。
call 存在，但服务端处理出错，即 h.Error 不为空。
call 存在，服务端处理正常，那么需要从 body 中读取 Reply 的值。

第一层锁sending锁定了正在发送请求的锁，防止在关闭客户端的过程中仍有请求正在发送，从而保证所有正在发送的请求完成之前不会关闭客户端。

第二层锁mu锁定了Client结构体的pending和shutdown字段。通过加锁操作，它可以防止其他goroutine在关闭客户端时修改pending和shutdown字段，从而保证线程安全性。因为在关闭客户端时，需要遍历pending字段，对其中的所有Call结构体设置Error字段并完成它们的done方法。而在这个过程中，其他goroutine可能会访问pending字段，并且对它进行修改，从而导致数据不一致的情况发生。

因此，通过在terminateCalls函数中加上两层锁，可以保证在关闭客户端时，所有请求都能够被安全地处理，从而避免数据竞争和其他并发问题。

func (client *Client) receive() 
创建 Client 实例时，首先需要完成一开始的协议交换，即发送 Option 信息给服务端。协商好消息的编解码方式之后，再创建一个子协程调用 receive() 接收响应。

此时，GeeRPC 客户端已经具备了完整的创建连接和接收响应的能力了，最后还需要实现发送请求的能力。

Go 和 Call 是客户端暴露给用户的两个 RPC 服务调用接口，Go 是一个异步接口，返回 call 实例。
Call 是对 Go 的封装，阻塞 call.Done，等待响应返回，是一个同步接口。
至此，一个支持异步和并发的 GeeRPC 客户端已经完成。