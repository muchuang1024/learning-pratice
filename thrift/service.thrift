include "./request.thrift"

typedef map<string, string> Data

struct Response {
    1:required i32 errCode; //错误码
    2:required string errMsg; //错误信息
    3:required Data data;
}

//定义服务
service Service {

    Response SayHello(
    	1:required request.SayHelloRequest reqParam
    )

    Response GetUser(
    	1:required request.GetUserRequest reqParam
    )
}

