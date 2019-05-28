### 准备

```
zk安装：see https://github.com/caijinlin/centos-dockerfile/blob/master/soft/zookeeper/build.sh
zk启动: see https://github.com/caijinlin/centos-dockerfile/blob/master/soft/zookeeper/control.sh 
```

执行命令后会启动三个zk server, 单机zk 伪集群


### 运行demo

```
go get github.com/samuel/go-zookeeper/zk
go run main.go
```
