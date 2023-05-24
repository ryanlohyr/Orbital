//To run the file type: go run <filename>.go

package main //creates a executable file, has to be same as <filename>

import "fmt" //library to format strings and print lines

func main() { //entry point of the file, has to be called main, every file has to have main
	fmt.Println("Hello, ni de dog")
}


//flow of project (admin)

//1. create a service for the user to connect to 

// curl -i -s -X POST http://localhost:8001/services \
//   --data name=example_service \
//   --data url='http://mockbin.org'


//2. 




//flow of project (user)

//1. send http request to gateway
//2. 



// curl -X GET http://localhost:8000/mock/requests