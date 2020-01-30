# IO模型

**linux环境下执行**

## 启动server

```
# 同步阻塞IO 
gcc blockServer.c -o blockServer
./blockServer
# IO多路复用select
gcc selectServer.c -o selectServer
./selectServer
# IO多路复用epoll
gcc epollServer.c -o epollServer
./epollServer
```

## 启动client

```
python3 client.py

```