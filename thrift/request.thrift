include "./types.thrift"

struct GetUserRequest {
	1:required string logid;
	2:required i32 uid;
}

struct SayHelloRequest {
	1:required string logid;
    2:required types.UserList userlist;
}