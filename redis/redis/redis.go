package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"strings"
	"time"
)

type RedisClient struct {
	masterAddr string
	pool       *redis.Pool
	readPool   *redis.Pool
}

type RedisCommand struct {
	sflags string // Flags as string representation, one char per flag.
}

func NewClient(addrs []string) *RedisClient {
	masterAddr := getMasterAddr(addrs)
	client := &RedisClient{
		masterAddr: masterAddr,
		readPool:   newPool(addrs),
		pool:       newPool([]string{masterAddr}),
	}

	return client
}

func (client *RedisClient) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	var conn redis.Conn
	// 读写分离
	if isWriteCommand(commandName) {
		conn = client.pool.Get()
	} else {
		conn = client.readPool.Get()
	}
	defer conn.Close()
	reply, err = conn.Do(commandName, args...)
	return
}

func newPool(addrs []string) *redis.Pool {
	var rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &redis.Pool{
		MaxIdle:     30,
		MaxActive:   60,
		IdleTimeout: 240 * time.Second,
		Wait:        false,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			// 随机数实现负载均衡
			index := rand_gen.Intn(len(addrs))
			return redis.Dial("tcp", addrs[index])
		},
	}
}

func getMasterAddr(addrs []string) string {
	masterAddr := ""
	for _, addr := range addrs {
		conn, _ := redis.Dial("tcp", addr)
		res, err := redis.String(conn.Do("INFO", "replication"))
		if err == nil {
			sres := strings.Split(res, "\r\n")

			for _, s := range sres {
				si := strings.Split(s, ":")
				if si[0] == "master_host" {
					masterAddr = si[1]
				}
				if si[0] == "master_port" {
					masterAddr = fmt.Sprintf("%s:%s", masterAddr, si[1])
				}
			}

			if masterAddr != "" {
				return masterAddr
			}
		}
	}

	return masterAddr
}

func isWriteCommand(cmd string) bool {
	var redisCommandTable = map[string]RedisCommand{
		"get":               {sflags: "rF"},
		"set":               {sflags: "wm"},
		"setnx":             {sflags: "wmF"},
		"setex":             {sflags: "wm"},
		"psetex":            {sflags: "wm"},
		"append":            {sflags: "wm"},
		"strlen":            {sflags: "rF"},
		"del":               {sflags: "w"},
		"exists":            {sflags: "rF"},
		"setbit":            {sflags: "wm"},
		"getbit":            {sflags: "rF"},
		"setrange":          {sflags: "wm"},
		"getrange":          {sflags: "r"},
		"substr":            {sflags: "r"},
		"incr":              {sflags: "wmF"},
		"decr":              {sflags: "wmF"},
		"mget":              {sflags: "r"},
		"rpush":             {sflags: "wmF"},
		"lpush":             {sflags: "wmF"},
		"rpushx":            {sflags: "wmF"},
		"lpushx":            {sflags: "wmF"},
		"linsert":           {sflags: "wm"},
		"rpop":              {sflags: "wF"},
		"lpop":              {sflags: "wF"},
		"brpop":             {sflags: "ws"},
		"brpoplpush":        {sflags: "wms"},
		"blpop":             {sflags: "ws"},
		"llen":              {sflags: "rF"},
		"lindex":            {sflags: "r"},
		"lset":              {sflags: "wm"},
		"lrange":            {sflags: "r"},
		"ltrim":             {sflags: "w"},
		"lrem":              {sflags: "w"},
		"rpoplpush":         {sflags: "wm"},
		"sadd":              {sflags: "wmF"},
		"srem":              {sflags: "wF"},
		"smove":             {sflags: "wF"},
		"sismember":         {sflags: "rF"},
		"scard":             {sflags: "rF"},
		"spop":              {sflags: "wRsF"},
		"srandmember":       {sflags: "rR"},
		"sinter":            {sflags: "rS"},
		"sinterstore":       {sflags: "wm"},
		"sunion":            {sflags: "rS"},
		"sunionstore":       {sflags: "wm"},
		"sdiff":             {sflags: "rS"},
		"sdiffstore":        {sflags: "wm"},
		"smembers":          {sflags: "rS"},
		"sscan":             {sflags: "rR"},
		"zadd":              {sflags: "wmF"},
		"zincrby":           {sflags: "wmF"},
		"zrem":              {sflags: "wF"},
		"zremrangebyscore":  {sflags: "w"},
		"zremrangebyrank":   {sflags: "w"},
		"zremrangebylex":    {sflags: "w"},
		"zunionstore":       {sflags: "wm"},
		"zinterstore":       {sflags: "wm"},
		"zrange":            {sflags: "r"},
		"zrangebyscore":     {sflags: "r"},
		"zrevrangebyscore":  {sflags: "r"},
		"zrangebylex":       {sflags: "r"},
		"zrevrangebylex":    {sflags: "r"},
		"zcount":            {sflags: "rF"},
		"zlexcount":         {sflags: "rF"},
		"zrevrange":         {sflags: "r"},
		"zcard":             {sflags: "rF"},
		"zscore":            {sflags: "rF"},
		"zrank":             {sflags: "rF"},
		"zrevrank":          {sflags: "rF"},
		"zscan":             {sflags: "rR"},
		"hset":              {sflags: "wmF"},
		"hsetnx":            {sflags: "wmF"},
		"hget":              {sflags: "rF"},
		"hmset":             {sflags: "wm"},
		"hmget":             {sflags: "r"},
		"hincrby":           {sflags: "wmF"},
		"hincrbyfloat":      {sflags: "wmF"},
		"hdel":              {sflags: "wF"},
		"hlen":              {sflags: "rF"},
		"hstrlen":           {sflags: "rF"},
		"hkeys":             {sflags: "rS"},
		"hvals":             {sflags: "rS"},
		"hgetall":           {sflags: "r"},
		"hexists":           {sflags: "rF"},
		"hscan":             {sflags: "rR"},
		"incrby":            {sflags: "wmF"},
		"decrby":            {sflags: "wmF"},
		"incrbyfloat":       {sflags: "wmF"},
		"getset":            {sflags: "wm"},
		"mset":              {sflags: "wm"},
		"msetnx":            {sflags: "wm"},
		"randomkey":         {sflags: "rR"},
		"select":            {sflags: "rlF"},
		"move":              {sflags: "wF"},
		"rename":            {sflags: "w"},
		"renamenx":          {sflags: "wF"},
		"expire":            {sflags: "wF"},
		"expireat":          {sflags: "wF"},
		"pexpire":           {sflags: "wF"},
		"pexpireat":         {sflags: "wF"},
		"keys":              {sflags: "rS"},
		"scan":              {sflags: "rR"},
		"dbsize":            {sflags: "rF"},
		"auth":              {sflags: "rsltF"},
		"ping":              {sflags: "rtF"},
		"echo":              {sflags: "rF"},
		"save":              {sflags: "ars"},
		"bgsave":            {sflags: "ar"},
		"bgrewriteaof":      {sflags: "ar"},
		"shutdown":          {sflags: "arlt"},
		"lastsave":          {sflags: "rRF"},
		"type":              {sflags: "rF"},
		"multi":             {sflags: "rsF"},
		"exec":              {sflags: "sM"},
		"discard":           {sflags: "rsF"},
		"sync":              {sflags: "ars"},
		"psync":             {sflags: "ars"},
		"replconf":          {sflags: "arslt"},
		"flushdb":           {sflags: "w"},
		"flushall":          {sflags: "w"},
		"sort":              {sflags: "wm"},
		"info":              {sflags: "rlt"},
		"monitor":           {sflags: "ars"},
		"ttl":               {sflags: "rF"},
		"pttl":              {sflags: "rF"},
		"persist":           {sflags: "wF"},
		"slaveof":           {sflags: "ast"},
		"role":              {sflags: "lst"},
		"debug":             {sflags: "as"},
		"config":            {sflags: "art"},
		"subscribe":         {sflags: "rpslt"},
		"unsubscribe":       {sflags: "rpslt"},
		"psubscribe":        {sflags: "rpslt"},
		"punsubscribe":      {sflags: "rpslt"},
		"publish":           {sflags: "pltrF"},
		"pubsub":            {sflags: "pltrR"},
		"watch":             {sflags: "rsF"},
		"unwatch":           {sflags: "rsF"},
		"cluster":           {sflags: "ar"},
		"restore":           {sflags: "wm"},
		"restore-asking":    {sflags: "wmk"},
		"migrate":           {sflags: "w"},
		"asking":            {sflags: "r"},
		"readonly":          {sflags: "rF"},
		"readwrite":         {sflags: "rF"},
		"dump":              {sflags: "r"},
		"object":            {sflags: "r"},
		"client":            {sflags: "rs"},
		"eval":              {sflags: "s"},
		"evalsha":           {sflags: "s"},
		"slowlog":           {sflags: "r"},
		"script":            {sflags: "rs"},
		"time":              {sflags: "rRF"},
		"bitop":             {sflags: "wm"},
		"bitcount":          {sflags: "r"},
		"bitpos":            {sflags: "r"},
		"wait":              {sflags: "rs"},
		"command":           {sflags: "rlt"},
		"geoadd":            {sflags: "wm"},
		"georadius":         {sflags: "r"},
		"georadiusbymember": {sflags: "r"},
		"geohash":           {sflags: "r"},
		"geopos":            {sflags: "r"},
		"geodist":           {sflags: "r"},
		"pfselftest":        {sflags: "r"},
		"pfadd":             {sflags: "wmF"},
		"pfcount":           {sflags: "r"},
		"pfmerge":           {sflags: "wm"},
		"pfdebug":           {sflags: "w"},
		"latency":           {sflags: "arslt"},
	}

	rc, ok := redisCommandTable[strings.ToLower(cmd)]
	if !ok {
		return false
	}
	return strings.Contains(rc.sflags, "w")
}
