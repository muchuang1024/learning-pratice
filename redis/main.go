package main

// 线上服务可把redis移动到公共库，不用相对路径引用
import (
	"./redis"
	"fmt"
	redislib "github.com/gomodule/redigo/redis"
	"net/http"
	"os"
)

func main() {

	addrs := []string{"127.0.0.1:6379", "127.0.0.1:6380"}

	// 实现了长连接与读写分离的redis client
	redisClient := redis.NewClient(addrs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		_, er := redisClient.Do("set", "username", "nick")
		if er != nil {
			fmt.Fprintf(w, "redis set failed:"+er.Error())
		}
		username, err := redislib.String(redisClient.Do("get", "username"))
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
