## gorpc

一个简单，易用，高性能，可插拔的微服务框架

### 安装

在安装 gorpc 之前，您需要安装 go 并配置 go 环境。

gopath 模式下，您只需执行以下命令即可安装：

```
go get -u -v github.com/lubanproj/gorpc
```

go modules 模式下，您只需导入软件包 “ github.com/lubanproj/gorpc”，该软件包在执行 go [build | run | test] 时会自动下载依赖项。

### 快速开始

```
git clone https://github.com/lubanproj/gorpc.git
cd gorpc/examples/helloworld
# start server
go run server/server.go
# start client，start another terminal and execute
go run client/client.go
```

### 示例

发起一个服务调用只需要以下三个步骤：

1. 定义一个服务
2. server 发布服务
3. 使用一个客户端发起调用

**1.定义一个服务**

```
type Service struct {

}

type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg : "world",
	}

	return rsp, nil
}

```

**2. server 发布服务**

```
func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithSerializationType("msgpack"),
		gorpc.WithTimeout(time.Millisecond * 2000),
	}
	s := gorpc.NewServer(opts ...)
	if err := s.RegisterService("/helloworld.Greeter", new(helloworld.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}
```

**3. 使用一个客户端发起调用**

```
func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
		client.WithSerializationType("msgpack"),
	}
	c := client.DefaultClient
	req := &helloworld.HelloRequest{
		Msg: "hello",
	}
	rsp := &helloworld.HelloReply{}
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
	fmt.Println(rsp.Msg, err)
}
```

更多细节可以参考 [helloworld](https://github.com/lubanproj/gorpc/tree/master/examples/helloworld) 


### 文档

- [Examples](https://github.com/lubanproj/gorpc/tree/master/examples).
- [FAQ](https://github.com/lubanproj/gorpc/wiki/FAQ)

### 特性

- **高性能**，性能远远超过 grpc ，详情可以参考 [性能](#Performance)
- 支持 **反射**, **代码生成** 两种调用方式
- **可插拔** 所有插件都是**可插拔、支持业务自定义的**
- **多协议支持**，目前支持 tcp 和 udp，后续会支持更多协议
- 实现了**拦截器**，支持业务自己定义拦截器
- 实现了**连接池**，支持业务自定义连接池
- 支持**服务发现**，提供了基于 **consul** 的默认服务发现实现，支持业务自定义服务发现实现。
- 支持**负载均衡** ，提供了**随机、轮询、加权轮询、一致性哈希**等默认负载均衡实现，支持业务自定义负载均衡实现。
- 支持**分布式链路追踪**，遵循业界 opentracing 规范，提供了基于 **jaeger** 的分布式链路追踪默认实现，支持业务自定义。
- 支持**多种序列化方式**，框架默认采用 **protocol** 和 **msgpack** 序列化，用代码生成方式调用会使用 protocol 序列化。用反射方式调用会采用 msgpack 序列化，支持业务自定义序列化方式。
- 更多特性正在陆续支持中 ......

### <span id="Performance">性能</span>

**环境** :

- CPU : Intel(R) Xeon(R) Gold 61xx CPU @2.44GHz
- CPU cores : 8
- Memory : 16G
- Disk : 540G

**Result** :
使用 [**gorpc-benchmark**](https://github.com/lubanproj/gorpc-benchmark) 进行测试，测试三次取性能最大值

**gorpc** :

```
git clone https://github.com/lubanproj/gorpc-benchmark.git
cd gorpc-benchmark
# start gorpc server
go run server.go
# start gorpc-benchmark client，start another terminal and execute
go run client.go -concurrency=100 -total=1000000
```

性能测试结果如下：

```
> go run client.go -concurrency=100 -total=1000000
2020/02/29 15:56:57 client.go:71: [INFO] took 5214 ms for 1000000 requests
2020/02/29 15:56:57 client.go:72: [INFO] sent     requests      : 1000000
2020/02/29 15:56:57 client.go:73: [INFO] received requests      : 1000000
2020/02/29 15:56:57 client.go:74: [INFO] received requests succ : 1000000
2020/02/29 15:56:57 client.go:75: [INFO] received requests fail : 0
2020/02/29 15:56:57 client.go:76: [INFO] throughput  (TPS)      : 191791
```

**grpc** : 

在相同机器上进行 grpc 性能测试，如下：

```
git clone https://github.com/lubanproj/gorpc-benchmark.git
cd gorpc-benchmark/grpc
# run gorpc server
go run server.go
# run gorpc-benchmark client, start another terminal and execute 
go run client.go -concurrency=100 -total=1000000
```

性能测试结果如下：

```
> go run client.go -concurrency=100 -total=1000000
2020/02/29 15:46:14 client.go:77: [INFO] took 17169 ms for 1000000 requests
2020/02/29 15:46:14 client.go:78: [INFO] sent     requests      : 1000000
2020/02/29 15:46:14 client.go:79: [INFO] received requests      : 1000000
2020/02/29 15:46:14 client.go:80: [INFO] received requests succ : 1000000
2020/02/29 15:46:14 client.go:81: [INFO] received requests fail : 0
2020/02/29 15:46:14 client.go:82: [INFO] throughput  (TPS)      : 58244
```

### 贡献

[贡献者](https://github.com/lubanproj/gorpc/graphs/contributors)

如何进行贡献？

可以参考 [Contributing](