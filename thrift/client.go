package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/caijinlin/learning-pratice/thrift/lib/request"
	"github.com/caijinlin/learning-pratice/thrift/lib/service"

	"net"
	"os"
)

func main() {
	host := "127.0.0.1"
	port := "9898"
	trans, err := thrift.NewTSocket(net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client := service.NewServiceClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}
	defer trans.Close()

	req := &request.GetUserRequest{Logid: "1", UID: 1}
	res, err := client.GetUser(req)
	if err != nil {
		fmt.Println("GetUser failed:", err)
		return
	}
	fmt.Println(res)
}
