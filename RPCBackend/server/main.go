package main

import (
    "github.com/cloudwego/kitex/pkg/generic"
    "github.com/cloudwego/kitex/server/genericserver"
    // "github.com/cloudwego/kitex/server"
    // "github.com/cloudwego/kitex/pkg/rpcinfo"
    // "github.com/kitex-contrib/registry-nacos/registry"
	"fmt"
    // "net"
	"context"
)

func main() {
    // Parse IDL with Local Files
    // YOUR_IDL_PATH thrift file path,eg: ./idl/example.thrift
    p, err := generic.NewThriftFileProvider("./hello.thrift")
    // r, err := registry.NewDefaultNacosRegistry()
    if err != nil {
        panic(err)
    }

    g, err := generic.JSONThriftGeneric(p)
    if err != nil {
        panic(err)
    }
    svr := genericserver.NewServer(
        new(GenericServiceImpl), 
        g,
        // server.WithRegistry(r),
        // server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}),
		// server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}),

    )
    if err != nil {
        panic(err)
    }
    errr := svr.Run()
    if errr != nil {
        panic(err)
    }
    // resp is a JSON string
}

type GenericServiceImpl struct {
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
        // use jsoniter or other json parse sdk to assert request
        m := request.(string)
        fmt.Printf("Recv: %v\n", m)
        return  "{\"Msg\": \"Post request recieved\"}", nil
}

