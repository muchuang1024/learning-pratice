### 准备

```
sh install.sh
```

执行命令后会启动两个redis server, master指向6379端口，slave指向6380端口，三个sentinel，监控redis master


### 登录redis

1. 登录master

```
redis-cli -p 6379
```

2. 登录slave

```
redis-cli -p 6380
```

### 登录sentinel

```
redis-cli -p 26379
redis-cli -p 26380
redis-cli -p 26381

```

info 查看监听redis

```
sentinel_masters:1
sentinel_tilt:0
sentinel_running_scripts:0
sentinel_scripts_queue_length:0
sentinel_simulate_failure_flags:0
master0:name=localmaster,status=ok,address=127.0.0.1:6379,slaves=1,sentinels=3
```

### 模拟主数据库故障

终端关闭数据库

```
/usr/local/bin/redis-cli -p 6379 shutdown  
```

查看sentinel日志

```
tail -f /var/log/redis-sentinel-*
```

追加了如下内容

```
=> /var/log/redis-sentinel-26381.log <==
24331:X 27 Mar 2019 18:33:55.121 # +promoted-slave slave 127.0.0.1:6380 127.0.0.1 6380 @ localmaster 127.0.0.1 6379
24331:X 27 Mar 2019 18:33:55.121 # +failover-state-reconf-slaves master localmaster 127.0.0.1 6379
24331:X 27 Mar 2019 18:33:55.191 * +slave-reconf-sent slave 172.18.0.1:6380 172.18.0.1 6380 @ localmaster 127.0.0.1 6379
24331:X 27 Mar 2019 18:33:55.191 # +failover-end master localmaster 127.0.0.1 6379
24331:X 27 Mar 2019 18:33:55.193 # +switch-master localmaster 127.0.0.1 6379 127.0.0.1 6380
24331:X 27 Mar 2019 18:33:55.193 * +slave slave 172.18.0.1:6380 172.18.0.1 6380 @ localmaster 127.0.0.1 6380
24331:X 27 Mar 2019 18:33:55.193 * +slave slave 127.0.0.1:6379 127.0.0.1 6379 @ localmaster 127.0.0.1 6380

==> /var/log/redis-sentinel-26379.log <==
24327:X 27 Mar 2019 18:33:55.194 # +config-update-from sentinel b7d901f1775f63057bad25fdd8117f008ba25532 127.0.0.1 26381 @ localmaster 127.0.0.1 6379
24327:X 27 Mar 2019 18:33:55.195 # +switch-master localmaster 127.0.0.1 6379 127.0.0.1 6380

==> /var/log/redis-sentinel-26380.log <==
24329:X 27 Mar 2019 18:33:55.194 # +config-update-from sentinel b7d901f1775f63057bad25fdd8117f008ba25532 127.0.0.1 26381 @ localmaster 127.0.0.1 6379
24329:X 27 Mar 2019 18:33:55.195 # +switch-master localmaster 127.0.0.1 6379 127.0.0.1 6380
```

### 恢复故障数据库

```
/usr/local/bin/redis-server ./conf/redis_6379.conf
```

查看sentinel日志，故障数据库加入slaves了

```
==> /var/log/redis-sentinel-26381.log <==
24331:X 28 Mar 2019 10:36:28.922 * +convert-to-slave slave 127.0.0.1:6379 127.0.0.1 6379 @ localmaster 127.0.0.1 6380
```
