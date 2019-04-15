package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	_ "net"
	"net/http"
	"os"
	"time"
)

func main() {

	pool := newPool("127.0.0.1:6380")

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

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   6,
		IdleTimeout: 1 * time.Second,
		Wait:        false,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}
