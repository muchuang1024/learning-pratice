## 编写thrift开发四部曲

### 1. 编写thrift idl 文件

```
# 这里我把一个服务拆成3个thrift文件了，完全可以放在一起，只用一个thrift文件
serivce.thrift
request.thrift
types.thrift
```

### 2. 基于idl生成thfit lib文件

```
# -out参数指定生成文件输出目录
thrift -gen go -out lib  service.thrift
thrift -gen go -out lib  request.thrift
thrift -gen go -out lib  types.thrift
```

### 3. 编写服务器端程序

```
go run server.go
```

### 4. 编写客户端程序

```
go run client.go
```

## 注意

> 如果thrift文件相互include，对于go而言，编译后的lib文件，需要修改import路径为绝对路径，所以建议使用一个thrift文件