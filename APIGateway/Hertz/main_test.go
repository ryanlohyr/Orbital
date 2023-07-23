package main

import (
    "testing"
    "fmt"
    "github.com/cloudwego/kitex/pkg/generic"
)




func TestValidServiceHosts(t *testing.T){
    //available services are "TravelService" and "add"
    hostList,ok := getServiceHosts("TravelService",serviceRegistryIP)["hosts"].([]interface{})
    if(!ok){
        t.Fatalf("List was not returned")
	}
    if len(hostList) == 0{
		t.Fatalf("Should not have returned an empty list")
	}

    hostList,ok = getServiceHosts("add",serviceRegistryIP)["hosts"].([]interface{})
    if(!ok){
        t.Fatalf("List was not returned")
	}
    if len(hostList) == 0{
		t.Fatalf("Should not have returned an empty list")
	}

}

func TestEmptyServiceHosts(t *testing.T){
    hostList,ok := getServiceHosts("invalidServiceName",serviceRegistryIP)["hosts"].([]interface{})
    if(!ok){
        t.Fatalf("List was not returned")
	}
    if len(hostList) != 0{
		t.Fatalf("Should return an empty list")
	}
}

func TestInitialisingGenericClient(t *testing.T){

    thriftDirectory := fmt.Sprintf("./thriftFiles/%s.thrift","TravelService")

    var IDLPath string = thriftDirectory

    p, err := generic.NewThriftFileProvider(IDLPath)

	if err != nil {
        fmt.Println("Error:", err)
		t.Fatalf("Should return an empty list")
	}

	g, err := generic.JSONThriftGeneric(p)

    if err != nil {
		t.Fatalf("Should return an empty list")
	}
    
    cli, err := initialiseClient(g ,"TravelService")
    if err != nil {
		t.Fatalf("Should return an empty list")
	}

    fmt.Println(cli)
}

func TestInitialisingGenericClientWithInvalidIDL(t *testing.T){

    thriftDirectory := fmt.Sprintf("./thriftFiles/%s.thrift","invalidIDLName")

    var IDLPath string = thriftDirectory

    p, err := generic.NewThriftFileProvider(IDLPath)

	if err == nil {
        fmt.Println("Error:", err)
		t.Fatalf("Should return an error as its an invalidIDL")
	}

    fmt.Println(p)
}




