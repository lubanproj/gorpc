## gorpc

### Installation
Before installing gorpc, you need to install go and configure the go environment

In gopath mode, you simply execute the following command to install

```
go get -u -v github.com/lubanproj/gorpc
```
In the go modules mode, you simply import the package "github.com/lubanproj/gorpc", which automatically downloads the dependency when you execute the go [build|run|test]

### Quick Start
```
git clone https://github.com/lubanproj/gorpc.git
cd gorpc/examples/helloworld
# start server
go run server/server.go
# start client，start another terminal and execute
go run client/client.go
```
### Documentation
[Examples](https://github.com/lubanproj/gorpc/tree/master/examples).
[FAQ](https://github.com/lubanproj/gorpc/wiki/FAQ)
### Features
- **High performance**, performance far exceeds that of the grpc, refer to [Performance](#Performance)
- Support **reflection**, **code generation** two ways to call
- All plug-ins are **configurable**
- **Multi-protocol** support, currently supported tcp, extensible and customizable
- Support **interceptor**, extensible and customizable
- Support **connection pooling**, extensible and customizable
- Support **service discovery**, provides the **consul** plug-in implementation, extensible and customizable
- Support **load balancing**, provides the **random**、**round robin**、**weighted round robin**、**consistent hash** implementation, extensible and customizable
- Support **opentracing**, provides the **jaeger** plug-in implementation, extensible and customizable
- Support for **a variety of serialization methods**, the default reflection calls using **protobuf** serialization,  code generation calls using **msgpack** serialization，support custom serialization
- More features are being supported ......

### <span id="Performance">Performance</span>
**Environment** :
- CPU : Intel(R) Xeon(R) Gold 61xx CPU @2.44GHz
- CPU cores : 8
- Memory : 16G
- Disk : 540G

**Result** :
The results were performed using [**gorpc-benchmark**](https://github.com/lubanproj/gorpc-benchmark) ，which performed three times to get the maximum

**gorpc** :
```
git clone https://github.com/lubanproj/gorpc-benchmark.git
cd gorpc-benchmark
# start gorpc server
go run server.go
# start gorpc-benchmark client，start another terminal and execute
go run client.go -concurrency=100 -total=1000000
```
The performance test results are as follows : 
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
Test grpc with the same machine :
```
git clone https://github.com/lubanproj/gorpc-benchmark.git
cd gorpc-benchmark/grpc
# 运行 gorpc server
go run server.go
# start gorpc-benchmark client, start another terminal and execute 
go run client.go -concurrency=100 -total=1000000
```
The performance test results are as follows : 
```
> go run client.go -concurrency=100 -total=1000000
2020/02/29 15:46:14 client.go:77: [INFO] took 17169 ms for 1000000 requests
2020/02/29 15:46:14 client.go:78: [INFO] sent     requests      : 1000000
2020/02/29 15:46:14 client.go:79: [INFO] received requests      : 1000000
2020/02/29 15:46:14 client.go:80: [INFO] received requests succ : 1000000
2020/02/29 15:46:14 client.go:81: [INFO] received requests fail : 0
2020/02/29 15:46:14 client.go:82: [INFO] throughput  (TPS)      : 58244
```

### Contributing
[Contributors](https://github.com/lubanproj/gorpc/graphs/contributors)

Welcome to submit issues|requirements|PRs and so on


