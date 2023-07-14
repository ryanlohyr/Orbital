package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// Parse IDL with Local Files
	p, err := generic.NewThriftFileProvider("./hello.thrift")
	if err != nil {
		panic(err)
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	// the nacos client config
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		// more ...
	}
	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	// r, err := registry.NewDefaultNacosRegistry()
	// if err != nil {
	//     panic(err)
	// }
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	svr0 := genericserver.NewServer(
		new(GenericServiceImpl),
		g,
		server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}),
	)

	svr1 := genericserver.NewServer(
		new(GenericServiceImpl2),
		g,
		server.WithServiceAddr(&net.TCPAddr{Port: 8889}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}),
	)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := svr0.Run(); err != nil {
			log.Println("server0 stopped with error:", err)
		} else {
			log.Println("server0 stopped")
		}
	}()
	go func() {
		defer wg.Done()
		if err := svr1.Run(); err != nil {
			log.Println("server1 stopped with error:", err)
		} else {
			log.Println("server1 stopped")
		}
	}()
	wg.Wait()
}

type GenericServiceImpl struct {
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Recv in server 1: %v\n", m)
	return "{\"Msg\": \"Post request recieved\"}", nil
}

type GenericServiceImpl2 struct {
}

func (g *GenericServiceImpl2) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Recv in server 2: %v\n", m)
	return "{\"Msg\": \"Post request recieved\"}", nil
}
