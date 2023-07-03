An API Gateway that accepts HTTP requests encoded in JSON format and uses the Generic-Call feature of Kitex to translate these requests into Thrift binary format requests. The API Gateway will then forward the request to one of the backend RPC servers discovered from the registry center. 

 ## Table of Contents
 - [Installation](#installation)
 - [Usage](#usage)
 - [Configuration](#configuration)
 - [License](#license)

 ## Installation

 ### Prerequisites
 - Go 1.15 or above
 - Docker
 - curl/postman (to simulate sending the post request from the client)
 - Thriftgo
 - Kitex


 ## Usage

 ### Configuration
 Before running the API Gateway, make sure to configure the following:

 - Make sure that GO111MODULE is set to on.
 - Install Kitex: go install github.com/cloudwego/kitex/tool/cmd/kitex@latest.
 - Install thriftgo: go install github.com/cloudwego/ thriftgo@latest.


 ### How to Run
 To test the API Gateway:

 ### Steps
 1. Clone the repository: `git clone https://github.com/ryanlohyr/Orbital.git`
 2. Install dependencies: `go mod download`
 3. Create two seperate terminals/command prompt. One to start the API Gateway and one to start the RPC Backend server
 4. For the first terminal, change your directory to be in the same directory as APIGateway/Hertz/main.go. Then run the api gateway with the command `go build -o hertz_demo && ./hertz_demo`
 5. For the second terminal, change your directory to be in the same directory as RPCBackend/server/main.go. Then run the api gateway with the command `go main.go`
 6. If both the API Gateways and RPC backend server are up and running you can send a curl request of  `curl -X POST -H "Content-Type: application/json" -d '{"text":"sup there"}' http://127.0.0.1:8881/post to test the API Gatway. 

 ## License

 This project is licensed under the [MIT License](LICENSE).

 ## References

 - [Kitex Documentation](https://www.cloudwego.io/docs/kitex )
 - [Hertz Documentation](https://www.cloudwego.io/docs/hertz/)