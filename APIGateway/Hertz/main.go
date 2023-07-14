package main

import (
	// "bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

type ctxKey int

const (
	ctxConsistentKey ctxKey = iota
)


/**
 *
 *Initializes a generic client using the given generic type and returns a client instance.
 * This function initializes a generic client by configuring various parameters such as server configuration,
 * client configuration, load balancing algorithm, resolver, and RPC timeout. It returns a client instance
 * along with an error, if any.
 * @param g The generic type to be used for the client.
 * @return The initialized client instance and an error, if any.
 *
**/
func initialiseClient(g generic.Generic) (genericclient.Client, error) {
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
	}

	resolvercli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	//NewWeightedBalancer creates a loadbalancer using weighted-round-robin algorithm.
	lb := loadbalance.NewWeightedBalancer()

	//client specifies the endpoint for the rpc backend
	cli, err := genericclient.NewClient(
		"Hello",
		g,
		client.WithHostPorts("0.0.0.0:8888", "0.0.0.0:8889"),
		client.WithLoadBalancer(lb),
		client.WithResolver(resolver.NewNacosResolver(resolvercli)),
		client.WithRPCTimeout(time.Second*3),
	)

	return cli, err
}

/**
 * Makes a Thrift call to the specified endpoint.
 *
 * @param IDLPath The path to the Thrift IDL file.
 * @param response The response message to be sent in the request body.
 * @param requestURL The URL of the request.
 * @param ctx The context for the request.
 * @return The response from the Thrift call.
 * @return An error if there was an issue with the Thrift call.
 */
func makeThriftCall(IDLPath string, response string, ctx context.Context) (interface{}, error) {
	var jsonData map[string]interface{}

	p, err := generic.NewThriftFileProvider(IDLPath)

	if err != nil {
		fmt.Println("error creating thrift file provider")
		return 0, err
	}

	g, err := generic.JSONThriftGeneric(p)

	if err != nil {
		return 0, errors.New(("error creating thrift generic"))
	}

	cli, err := initialiseClient(g)

	if err != nil {
		return 0, errors.New(("error creating client"))
	}

	//inputs message sent by client
	message := fmt.Sprintf("{\"Msg\": \"%s\"}", response)

	ctx = context.WithValue(ctx, ctxConsistentKey, "my key0")
	
	resp, err := cli.GenericCall(ctx, "ExampleMethod", message)

	if err != nil {
		fmt.Println(err)
		return 0, errors.New(("error making rpc call to server "))
	}

	str, ok := resp.(string)

	if !ok {
		return 0, errors.New(("not a string"))
	}

	//converts JSON string into JSON object
	json.Unmarshal([]byte((str)), &jsonData)

	return jsonData, nil

}

func main() {

	h := server.Default(server.WithHostPorts("0.0.0.0:8881"))

	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "hello from ryan"})
	})

	h.GET("/get", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "get")
	})

	h.POST("/post", func(ctx context.Context, c *app.RequestContext) {

		//url to send request to
		var IDLPath string = "./hello.thrift"
		var jsonData map[string]interface{}

		//returns data in an array of bytes
		response := c.GetRawData()

		//converts the array of bytes into array format and loads it into jsonData
		err := json.Unmarshal(response, &jsonData)

		if err != nil {
			fmt.Println("Error:", err)
			c.String(consts.StatusBadRequest, "bad post request")
			return
		}

		//whatever the key value is,  has to be consistent with backend

		//in this case key must be set as 'text'
		dataValue, ok := jsonData["text"]
		if !ok {
			//error handling
			c.String(consts.StatusBadRequest, `key provided has to be called "text" `)
			return
		}

		//ensures that data is a string
		stringValue, ok := dataValue.(string)

		//request validation
		if !ok {
			//error handling
			c.String(consts.StatusBadRequest, `value has to be string `)
			return
		}

		//converts the response to thrift binary format
		responseFromRPC, err := makeThriftCall(IDLPath, stringValue, ctx)

		if err != nil {
			fmt.Println(err)
			c.String(consts.StatusBadRequest, "error in thrift call ")
			return
		}

		fmt.Println("Post Request successful")

		c.JSON(consts.StatusOK, responseFromRPC)

	})

	h.PUT("/put", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "put")
	})
	h.DELETE("/delete", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "delete")
	})
	h.PATCH("/patch", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "patch")
	})

	//spin runs the application
	h.Spin()
}

//converts json into []bytes
//jsonBytes := []byte(`{"data":"helloworld"}`)
