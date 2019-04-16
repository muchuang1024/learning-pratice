## 代码目录
```
├── README.md # 操作指南
├── conf      # 配置文件
│   ├── redis_6379.conf
│   ├── redis_6380.conf
│   ├── sentinel_26379.conf
│   ├── sentinel_26380.conf
│   └── sentinel_26381.conf
├── install.sh # redis与sentinel启动脚本
├── main.go    # demo示例
└── redis
    └── redis.go # redis封装库（长连接与读写分离）
```


### 准备

```
sh install.sh
```

执行命令后会启动两个redis server, master指向6379端口，slave指向6380端口，三个sentinel，监控redis master


### 使用redis

登录master

```
redis-cli -p 6379
```

登录slave

```
redis-cli -p 6380
```

查看连接clients及设置

```
redis-cli -p 6379 info clients
redis-cli -p 6379 client list
redis-cli -p 6379 config get maxclients
```

### 使用sentinel

登录各个sentinel

```
redis-cli -p 26379
redis-cli -p 26380
redis-cli -p 26381

```

info 查看监听的redis master实例

```
sentinel_masters:1
sentinel_tilt:0
sentinel_running_scripts:0
sentinel_scripts_queue_length:0
sentinel_simulate_failure_flags:0
master0:name=localmaster,status=ok,address=127.0.0.1:6379,slaves=1,sentinels=3
```

### 模拟主数据库故障

关闭主数据库

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

### 模拟故障数据库恢复

启动故障数据库

```
/usr/local/bin/redis-server ./conf/redis_6379.conf
```

查看sentinel日志，故障数据库加入slaves了

```
==> /var/log/redis-sentinel-26381.log <==
24331:X 28 Mar 2019 10:36:28.922 * +convert-to-slave slave 127.0.0.1:6379 127.0.0.1 6379 @ localmaster 127.0.0.1 6380
```

### 长连接与读写分离demo

* 辅助库

```
https://github.com/gomodule/redigo
```

* 运行demo
```
go get github.com/gomodule/redigo
go run main.go
```

* 测试用例
```
# 模拟20个请求，10个并发 curl 请求
ab -c 10 -n 20 http://127.0.0.1:3333/
```

* 查看连接数
```
netstat -anp | grep 6380
```

* 核心代码解析

```
// get prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) get(ctx context.Context) (*poolConn, error) {

	// Handle limit for p.Wait == true.
	var waited time.Duration
	if p.Wait && p.MaxActive > 0 {
		p.lazyInit()

		// wait indicates if we believe it will block so its not 100% accurate
		// however for stats it should be good enough.
		wait := len(p.ch) == 0
		var start time.Time
		if wait {
			start = time.Now()
		}
		if ctx == nil {
			<-p.ch
		} else {
			select {
			case <-p.ch:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		if wait {
			waited = time.Since(start)
		}
	}

    // 防止高并发时两个请求用一个socket导致错误，加锁获取连接
	p.mu.Lock()

	if waited > 0 {
		p.waitCount++
		p.waitDuration += waited
	}

    // 移除过了IdleTimeout时间没有使用的连接
	// Prune stale connections at the back of the idle list.
	if p.IdleTimeout > 0 {
		n := p.idle.count
		for i := 0; i < n && p.idle.back != nil && p.idle.back.t.Add(p.IdleTimeout).Before(nowFunc()); i++ {
			pc := p.idle.back
			p.idle.popBack()
			p.mu.Unlock()
			pc.c.Close()
			p.mu.Lock()
			p.active--
		}
	}

    // 从最大长度为MaxIdle的idlelist弹出一个连接
	// Get idle connection from the front of idle list.
	for p.idle.front != nil {
		pc := p.idle.front
		p.idle.popFront()
		p.mu.Unlock()
		if (p.TestOnBorrow == nil || p.TestOnBorrow(pc.c, pc.t) == nil) &&
			(p.MaxConnLifetime == 0 || nowFunc().Sub(pc.created) < p.MaxConnLifetime) {
			return pc, nil
		}
		pc.c.Close()
		p.mu.Lock()
		p.active--
	}

    // 异常连接报错
	// Check for pool closed before dialing a new connection.
	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("redigo: get on closed pool")
	}

    // 该时刻超出设置的活跃连接数报错
	// Handle limit for p.Wait == false.
	if !p.Wait && p.MaxActive > 0 && p.active >= p.MaxActive {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}

	p.active++
	p.mu.Unlock()
	c, err := p.dial(ctx)
	if err != nil {
		c = nil
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		p.mu.Unlock()
	}
	return &poolConn{c: c, created: nowFunc()}, err
}
```
