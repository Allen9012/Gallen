如果使用硬编码的方式来实现结构体与服务的映射，那么每暴露一个方法，就需要编写等量的代码。那有没有什么方式，能够将这个映射过程自动化呢？可以借助反射。

通过反射，我们能够非常容易地获取某个结构体的所有方法，并且能够通过方法，获取到该方法所有的参数类型与返回值。

通过反射实现 service


每一个 methodType 实例包含了一个方法的完整信息。包括

method：方法本身
ArgType：第一个参数的类型
ReplyType：第二个参数的类型
numCalls：后续统计方法调用次数时会用到
另外，我们还实现了 2 个方法 newArgv 和 newReplyv，用于创建对应类型的实例。newArgv 方法有一个小细节，指针类型和值类型创建实例的方式有细微区别。

service 的定义也是非常简洁的，name 即映射的结构体的名称，比如 T，比如 WaitGroup；typ 是结构体的类型；rcvr 即结构体的实例本身，保留 rcvr 是因为在调用时需要 rcvr 作为第 0 个参数；method 是 map 类型，存储映射的结构体的所有符合条件的方法。

接下来，完成构造函数 newService，入参是任意需要映射为服务的结构体实例。

registerMethods 过滤出了符合条件的方法：

两个导出或内置类型的入参（反射时为 3 个，第 0 个是自身，类似于 python 的 self，java 中的 this）
返回值有且只有 1 个，类型为 error

## 集成到服务端
通过反射结构体已经映射为服务，但请求的处理过程还没有完成。从接收到请求到回复还差以下几个步骤：第一步，根据入参类型，将请求的 body 反序列化；第二步，调用 service.call，完成方法调用；第三步，将 reply 序列化为字节流，构造响应报文，返回。

回到代码本身，补全之前在 server.go 中遗留的 2 个 TODO 任务 readRequest 和 handleRequest 即可。

在这之前，我们还需要为 Server 实现一个方法 Register。

配套实现 findService 方法，即通过 ServiceMethod 从 serviceMap 中找到对应的 service

findService 的实现看似比较繁琐，但是逻辑还是非常清晰的。因为 ServiceMethod 的构成是 “Service.Method”，因此先将其分割成 2 部分，第一部分是 Service 的名称，第二部分即方法名。现在 serviceMap 中找到对应的 service 实例，再从 service 实例的 method 中，找到对应的 methodType。

readRequest 方法中最重要的部分，即通过 newArgv() 和 newReplyv() 两个方法创建出两个入参实例，然后通过 cc.ReadBody() 将请求报文反序列化为第一个入参 argv，在这里同样需要注意 argv 可能是值类型，也可能是指针类型，所以处理方式有点差异。