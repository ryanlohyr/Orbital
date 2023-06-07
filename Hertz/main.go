package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

//run command in command line: go build -o hertz_demo && ./hertz_demo
//to test endpoint: curl http://127.0.0.1:8888/ping

func makeThriftCall(IDLPath string, response []byte, requestURL string, ctx context.Context) (interface{},  error)  {
	p, err := generic.NewThriftFileProvider(IDLPath)
	if err != nil {
		fmt.Println("error creating thrift file provider")
		return 0, err
	}
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		return 0, errors.New(("error creating thrift generic"))
	} 
	cli, err := genericclient.NewClient("Hello", g)


	if err != nil {
		return 0, errors.New(("invalid client name"))
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(response))
	req.Header.Set("token", "1")
	if err != nil {
		fmt.Println("error constructing req")
		return 0, err
	}

	customReq, err := generic.FromHTTPRequest(req)

	if err != nil {
		fmt.Println("error creating custom req")
		return 0, err
	}

	//error 1 : function lookup failed, path=/login
	//error 2 : function lookup failed, no root with method=GET

	fmt.Println(customReq)
	

	resp, err := cli.GenericCall(ctx, "", customReq)

	if err != nil {
		fmt.Println("error making generic call")
		return 0,err
	}

	fmt.Println("generic call successful")
	fmt.Println(resp)
	
    return resp, nil
}

func main() {
	//we can use server.WithHostPorts as it returns a type of config.option
	h := server.Default(server.WithHostPorts("0.0.0.0:8881")) 

	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "hello from ryan"})
	})

	h.GET("/get", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "get")
		// c.AbortWithStatus(300)
	})

	h.POST("/post", func(ctx context.Context, c *app.RequestContext) {
        //url to send request to??
		var requestUrl string = "https://127.0.0.1:8888"
		var IDLPath string = "./hello.thrift"
		var jsonData map[string]interface{}
        
		response := c.GetRawData() //returns data in an array of bytes

		err := json.Unmarshal(response, &jsonData)

		if err != nil {
			fmt.Println("Error:", err)
			c.String(consts.StatusBadRequest, "bad post request")
			return
		}

		dataValue, ok := jsonData["Message"].(string)
        fmt.Println("message is " + dataValue)
		if !ok {
			c.String(consts.StatusBadRequest, "data provided was not a string")
			return
		}
		
		responseFromRPC, err := makeThriftCall(IDLPath, response, requestUrl,ctx)

		
        if(err != nil){
			fmt.Println(err)
			c.String(consts.StatusBadRequest, "error in thrift call ")
			return
        }

		fmt.Println(responseFromRPC)

		stringResponseFromRPC, ok := responseFromRPC.(string)
		if(!ok){
			c.String(consts.StatusBadRequest, "not string")
			return
		}
		
		//returns the string as a response, also acts as a return statement
		c.String(consts.StatusOK, stringResponseFromRPC)
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
