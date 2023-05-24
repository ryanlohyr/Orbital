# API Gateway

An API Gateway that accepts HTTP requests encoded in JSON format and uses the Generic-Call feature of Kitex to translate these requests into Thrift binary format requests. The API Gateway will then forward the request to one of the backend RPC servers discovered from the registry center. 

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [License](#license)

## Installation

### Prerequisites
- Go 1.15 or above
- Kong API Gateway
- Consul
- Docker
- curl

### Steps
1. Clone the repository: `git clone https://github.com/your/repository.git`
2. Install dependencies: `go mod download`
3. Configure Kong API Gateway and Consul according to your environment.
4. Set up the necessary configuration files and environment variables.
5. Build the API Gateway: `go build -o api-gateway`
6. Start the API Gateway: `curl -Ls https://get.konghq.com/quickstart | bash`

## Usage

### Configuration
Before running the API Gateway, make sure to configure the following:

- Kong API Gateway: Set up routing, transformations, plugins, and load balancing configurations in Kong.
- Consul: Configure Consul as the service registry center and register the backend RPC servers.


### How to Run
To start the Kong API Gateway:

1. Ensure Kong and Consul are running and properly configured.
2. Execute the following command: `curl -Ls https://get.konghq.com/quickstart | bash`

### Example Requests
Once the API Gateway is running, you can send HTTP requests to it. Here are some example requests:

1. Request: POST http://localhost:8001/services
   - Request Body: url='http://mockbin.org', name=example_service
   - Description: create a service for the user to connect to.
   - Response: 201 response header will be returned if service was created.

2. Request: POST http://localhost:8001/services/example_service/routes
   - Request Body: 'paths[]=/mock', name=example_route
   - Description: Configure a new route on the /mock path to direct traffic 
to the example_service service created earlier.
   - Response: 201 response header will be returned if service was created.

3. Request: GET http://localhost:8000/mock/requests
    - Description: Proxy a request through Kong Gateway to the /requests resource
    - Response: Response of the initial service will be returned.

## Configuration

### API Gateway Configuration
The API Gateway can be configured through environment variables or a configuration file. The following options can be customized:

- Port: Specify the port on which the API Gateway listens.
- Kong URL: Set the URL of the Kong Admin API for integration.
- Consul URL: Configure the URL of the Consul service registry center.

### Service Registry Configuration
Ensure Consul is properly configured and running. The backend RPC servers should register themselves with Consul upon startup, providing their network locations and metadata.


## License

This project is licensed under the [MIT License](LICENSE).

## References

- [Kong API Gateway Documentation](https://docs.konghq.com/)
- [Consul Documentation](https://www.consul.io/docs)