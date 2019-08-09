package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/caijinlin/learning-pratice/thrift/lib/request"
	"github.com/caijinlin/learning-pratice/thrift/lib/service"
)

func main() {
	transport, err := thrift.NewTServerSocket(":9898")
	if err != nil {
		panic(err)
	}

	handler := &Handle{}
	processor := service.NewServiceProcessor(handler)

	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		transportFactory,
		protocolFactory,
	)

	fmt.Println("正在监听...")

	if err := server.Serve(); err != nil {
		panic(err)
	}
}

/**
* 接口实现，做项目的时候单独放一个文件夹
 */
type Handle struct {
}

func (*Handle) SayHello(req *request.SayHelloRequest) (*service.Response, error) {
	fmt.Printf("message from client: %v\n", req.GetLogid())

	res := &service.Response{
		ErrCode: 0,
		ErrMsg:  "success",
	}

	return res, nil
}

func (*Handle) GetUser(req *request.GetUserRequest) (*service.Response, error) {
	fmt.Printf("recv request from client: %v\n", req.GetLogid())

	res := &service.Response{
		ErrCode: 0,
		ErrMsg:  "success",
		Data: map[string]string{
			"uid":  "1",
			"name": "cai",
		},
	}

	return res, nil
}
