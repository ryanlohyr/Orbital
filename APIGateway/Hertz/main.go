package main

import (
	// "bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ctxKey int

const (
	ctxConsistentKey ctxKey = iota
	serviceRegistryIP = "http://127.0.0.1:8848"

)

/**
 * Retrieves a list of service hosts from a service registry.
 *
 * @param hosts               The name of the service to query for hosts.
 * @param serviceRegistryIP   The IP address of the service registry.
 *
 * @return                    A map containing the retrieved service hosts.
 *                            The keys represent the host names, and the values
 *                            are interface{} values containing the host information.
 */
func getServiceHosts(hosts string, serviceRegistryIP string) (map[string]interface{}){
		route := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s",serviceRegistryIP,hosts)
		response, err := http.Get(route)
		var jsonData map[string]interface{}
		if err != nil {
			fmt.Printf("Error making GET request: %s\n", err)
			return jsonData
		}
		// Read the response body
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			return jsonData
		}

		err = json.Unmarshal(body, &jsonData)

		if err != nil {
			fmt.Printf("Error converting response body: %s\n", err)
			return jsonData
		}

		fmt.Println("Get Host request successful")
		return jsonData
}


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
func initialiseClient(g generic.Generic,serviceName string) (genericclient.Client, error) {
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

	hostList,ok := getServiceHosts(serviceName,serviceRegistryIP)["hosts"].([]interface{})
	if(!ok){
		return nil, errors.New(("error converting into list"))
	}

	if len(hostList) == 0{
		return nil, errors.New(("service name not found"))
	}


	//client specifies the endpoint for the rpc backend
	cli, err := genericclient.NewClient(
		serviceName,
		g,
		//we dont need to specify port names anymore as we are now using service discovery
		// client.WithHostPorts("0.0.0.0:8888", "0.0.0.0:8889"),
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
func makeThriftCall(IDLPath string, response string,serviceName string,methodName string, ctx context.Context) (interface{}, error) {
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

	cli, err := initialiseClient(g ,serviceName)

	if err != nil {
		return nil, err
	}

	//inputs message sent by client
	message := fmt.Sprintf("{\"Msg\": \"%s\"}", response)

	ctx = context.WithValue(ctx, ctxConsistentKey, "my key0")
	var resp interface{}
	resp, err = cli.GenericCall(ctx, methodName, message)
	
	//TO DO: create rpc call like in client/hello/main.go in OrbitalTest
	// if(serviceName == "Hello"){
	// 	resp, err = cli.GenericCall(ctx, methodName, message)
	// }else{ //since we only have two services, if serviceName is echo it will fall here
	// 	resp, err = cli.Echo(ctx, methodName, message)
	// }

	if err != nil {
		fmt.Println(err)
		return 0, err
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

	h.GET("/getServiceHosts/:hosts", func(ctx context.Context, c *app.RequestContext) {
		hosts := c.Param("hosts")
		
		c.JSON(consts.StatusOK, getServiceHosts(hosts,serviceRegistryIP))
	})

	h.POST("/post/:serviceName/:methodName", func(ctx context.Context, c *app.RequestContext) {

		serviceName := c.Param("serviceName")
		
		methodName := c.Param("methodName")

		var jsonData map[string]interface{}

		thriftDirectory := fmt.Sprintf("./thriftFiles/%s.thrift",serviceName)

		var IDLPath string = thriftDirectory

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
		dataValue, ok := jsonData["data"]

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
		responseFromRPC, err := makeThriftCall(IDLPath, stringValue,serviceName,methodName, ctx)

		if err != nil {
			fmt.Println("suo")
			fmt.Println(err)
			fmt.Println("hello")
			// c.Error(err)
			c.String(consts.StatusBadRequest, err.Error())
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
