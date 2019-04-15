package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	_ "net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	addrs := []string{"127.0.0.1:6379", "127.0.0.1:6380"}
	masterAddr := getMasterAddr(addrs)

	fmt.Println(masterAddr)
	// pool := newPool(addrs) // read pool
	pool := newPool([]string{masterAddr}) // write pool

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		c := pool.Get()

		defer c.Close() // put or close

		_, er := c.Do("SET", "username", "nick")
		if er != nil {
			fmt.Fprintf(w, "redis set failed:"+er.Error())
		}
		username, err := redis.String(c.Do("GET", "username"))
		if err != nil {
			fmt.Fprintf(w, "redis get failed:"+err.Error())
		} else {
			fmt.Fprintf(w, "Got username:"+username)
		}
	})

	server := http.Server{
		Addr: ":3333",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
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
