package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

/**
 * @brief Initializes and returns a generic.Generic instance from a Thrift definition file.
 * @param[in] thriftName The name of the Thrift definition file (without extension) located in the
 *                       "thriftFiles" directory.
 *
 * @return A generic.Generic instance initialized from the Thrift definition file on success.
 * @return An error on failure, such as if the Thrift file cannot be found or if there are errors
 *         while processing the Thrift file.
 */
func initialiseThriftGeneric(thriftName string)(generic.Generic,error){
	thriftDirectory := fmt.Sprintf("./thriftFiles/%s.thrift",thriftName)
	// Parse IDL with Local Files
	p, err := generic.NewThriftFileProvider(thriftDirectory)
	if err != nil {
		return nil,err
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return g,err
	}

	return g,nil
}


func main() {

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

	g_one, err := initialiseThriftGeneric("TravelService")
	if err != nil {
		panic(err)
	}

	g_two, err := initialiseThriftGeneric("ReviewService")
	if err != nil {
		panic(err)
	}
	
	svr0 := genericserver.NewServer(
		new(GenericServiceImpl),
		g_one,
		server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "TravelService"}),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}),
	)

	svr1 := genericserver.NewServer(
		new(GenericServiceImpl),
		g_one,
		server.WithServiceAddr(&net.TCPAddr{Port: 8889}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "TravelService"}),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}),
	)


	svr2 := genericserver.NewServer(
		new(GenericServiceImpl2),
		g_two,
		server.WithServiceAddr(&net.TCPAddr{Port: 8887}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "ReviewService"}),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}),
	)

	svr3 := genericserver.NewServer(
		new(GenericServiceImpl2),
		g_two,
		server.WithServiceAddr(&net.TCPAddr{Port: 8886}),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "ReviewService"}),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}),
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

	go func() {
		defer wg.Done()
		if err := svr2.Run(); err != nil {
			log.Println("server1 stopped with error:", err)
		} else {
			log.Println("server1 stopped")
		}
	}()

	go func() {
		defer wg.Done()
		if err := svr3.Run(); err != nil {
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
	fmt.Printf("Recv in server 2: %v\n", m)
	fmt.Printf("Method is %s",method)

    jsonData := []byte(m)
	var data map[string]interface{}

    // Unmarshal the JSON data into the map
    err = json.Unmarshal(jsonData, &data)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	var jsonResponse string 

	switch method {
	case "SendClientData":
		fmt.Println((data["Msg"]))
		jsonResponse = fmt.Sprintf("{\"Msg\": \"Post request recieved, the message sent was %s\",\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}",data["Msg"])
	case "RetrieveClientData":
		jsonResponse = fmt.Sprintf("{\"VisitedCountries\": [\"%s\",\"Singapore\",\"Malaysia\",\"Japan\"],\"Name\": \"Ryan\",\"userID\": %i,\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}","Taiwan",data["userID"])
	case "GetAllTravelDestinations":
		jsonResponse = fmt.Sprintf("{\"Destinations\": [\"%s\",\"Japan\",\"Sweden\",\"Netherlands\"],\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}","Myammar")
	default:
		jsonResponse = "{\"Msg\": \"Invalid Service name\"}"
	}
	return jsonResponse, nil
}

type GenericServiceImpl2 struct {
}

func (g *GenericServiceImpl2) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Recv in echo server : %v\n", m)
	fmt.Printf("Method is %s \n",method)

    jsonData := []byte(m)
	var data map[string]interface{}

    // Unmarshal the JSON data into the map
    err = json.Unmarshal(jsonData, &data)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	var jsonResponse string 

	switch method {
	case "sendReview":
		jsonResponse = fmt.Sprintf("{\"action\": \"%s was successfully uploaded\",\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}","Review Upload")
	case "editReview":
		jsonResponse = fmt.Sprintf("{\"action\": \"%s was successfully uploaded\",\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}","Review Edit")
	case "deleteReview":
		jsonResponse = fmt.Sprintf("{\"action\": \"%s was successfully uploaded\",\"BaseResp\":{\"StatusCode\":200,\"StatusMessage\":\"Success\"}}","Review Deletion")
	default:
		jsonResponse = "{\"Msg\": \"Invalid Service name\"}"
	}
	return jsonResponse, nil
	
}





