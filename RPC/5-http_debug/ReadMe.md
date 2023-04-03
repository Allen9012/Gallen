Web 开发中，我们经常使用 HTTP 协议中的 HEAD、GET、POST 等方式发送请求，等待响应。但 RPC 的消息格式与标准的 HTTP 协议并不兼容，在这种情况下，就需要一个协议的转换过程。HTTP 协议的 CONNECT 方法恰好提供了这个能力，CONNECT 一般用于代理服务。

假设浏览器与服务器之间的 HTTPS 通信都是加密的，浏览器通过代理服务器发起 HTTPS 请求时，
由于请求的站点地址和端口号都是加密保存在 HTTPS 请求报文头中的，代理服务器如何知道往哪里发送请求呢？
为了解决这个问题，浏览器通过 HTTP 明文形式向代理服务器发送一个 CONNECT 请求告诉代理服务器目标地址和端口，
代理服务器接收到这个请求后，会在对应端口与目标站点建立一个 TCP 连接，连接建立成功后返回 HTTP 200 状态码告诉浏览器与该站点的加密通道已经完成。
接下来代理服务器仅需透传浏览器和服务器之间的加密数据包即可，代理服务器无需解析 HTTPS 报文。

服务端支持 HTTP 协议
那通信过程应该是这样的：

客户端向 RPC 服务器发送 CONNECT 请求
1
CONNECT 10.0.0.1:9999/_geerpc_ HTTP/1.0
RPC 服务器返回 HTTP 200 状态码表示连接建立。
1
HTTP/1.0 200 Connected to Gee RPC
客户端使用创建好的连接发送 RPC 报文，先发送 Option，再发送 N 个请求报文，服务端处理 RPC 请求并响应。
在 server.go 中新增如下的方法：


defaultDebugPath 是为后续 DEBUG 页面预留的地址。

也就是说，只需要实现接口 Handler 即可作为一个 HTTP Handler 处理 HTTP 请求。接口 Handler 只定义了一个方法 ServeHTTP，实现该方法即可。

## 客户端支持 HTTP 协议
服务端已经能够接受 CONNECT 请求，并返回了 200 状态码 HTTP/1.0 200 Connected to Gee RPC，客户端要做的，发起 CONNECT 请求，检查返回状态码即可成功建立连接。

通过 HTTP CONNECT 请求建立连接之后，后续的通信过程就交给 NewClient 了。

为了简化调用，提供了一个统一入口 XDial

## 实现简单的 DEBUG 页面
支持 HTTP 协议的好处在于，RPC 服务仅仅使用了监听端口的 /_geerpc 路径，在其他路径上我们可以提供诸如日志、统计等更为丰富的功能。接下来我们在 /debug/geerpc 上展示服务的调用统计视图。