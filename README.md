# kudejen

This application simple application which serves a gRPC server that interacts with a running kubernetes cluster. As of now the gRPC server provides a simple API for provisioning postgres databases, allowing clients to perform CRUD (Create, Read, Update, Delete) operations. This setup leverages a running Kubernetes.

## Installation

there is already a Makefile containing all the necessary commands.

## Makefile

#### Targets

**build:** This target builds the application. It depends on the clean and generate targets.

**clean:** This target cleans up any generated files and directories.
generate: This target generates Go code from the Protobuf definition files.

**compile:** This target compiles the application Go code to a binary.

**run:** This target runs the compiled application binary.

**test:** This target runs the application tests using go test.

**docker-build:** This target builds a Docker image for the application.

**docker-build-debug:** This target builds a Docker image with debug symbols for the application.

**docker-run:** This target runs the application in a detached Docker container.

**docker-run-attach:** This target runs the application in an attached Docker container, allowing you to interact with the terminal inside the container.

**docker-stop:** This target stops the running Docker container.

So for building the image while you are on the root folder.
```bash
make docker-build   
```
For Running The application.
```bash
make docker-run
```
This will open port 8080 for grpc and port 8081 for http on your local machine. if you want to change this ports you can achieve that in config file located src/internal/config/config.yml, just be sure that you have to change the port binding in docker-run command as well.

the docker image uses a non-root user ro run the application.

**** there are some vulnerabilities and recommendation which can be accessed using the below command 
```bash
docker scout cves local://kudejen:latest
## Recommended fixes

  Base image is  alpine:3 

  Name            │  3                                                                         
  Digest          │  sha256:647a509e17306d117943c7f91de542d7fb048133f59c72dff7893cffd1836e11   
  Vulnerabilities │    1C     0H     0M     0L                                                 
  Pushed          │ 3 weeks ago                                                                
  Size            │ 4.1 MB                                                                     
  Packages        │ 17                                                                         
  OS              │ 3.20.1                                                                     

                                                                                                                                                                                               
  │ The base image is also available under the supported tag(s)  3.20 ,  3.20.1 ,  latest . If you want to display recommendations specifically for a different tag, please re-run the command  
  │ using the  --tag  flag.   
```
**"However, it falls outside the scope of our current project."**

Kubernetes communication requires the kube.config file, located in "src/internal/config/kube.config". Using this method is not the best method for connecting to Kubernetes.

Using a kubeconfig file is the easiest and most common approach.
KUBECONFIG can be set as an environment variable or loaded explicitly.
However, it is less secure (environment is variable) and offers limited control.

Kubernetes communication requires the kube.config file, located in "src/internal/config/kube.config". Using this method is not the best method for connecting to Kubernetes.
Using a kubeconfig file is the easiest and most common approach.
KUBECONFIG can be set as an environment variable or loaded explicitly.
However, it is less secure (environment is variable) and offers limited control.

there other better alternative but since it needs more assumption on the cluster that the application is going to run on.**"However, it falls outside the scope of our current project."**

**In-cluster config (recommended in cluster):** Leverages service account, most secure within cluster.

**Service tokens:** Short-lived tokens for tasks without service account (manage lifecycle, less secure).

**RBAC registration:** Fine-grained access control (more config, secure).

if you want to get your current configuration from .kube for example you can simply do so like below.(**only for test purposes**)
```bash
kubectl config view --minify --raw > kube.config
```

## Usage

when the application is up and running 2 port will be available on your designated domain(localhost). one for http server (:8081/) and another for grpc server (:/8080/)

there are some approach to create api documentation for grpc servers but **"However, it falls outside the scope of our current project."**

## gRPC API Documentation Best Practices

| **Approach** | **Description** | **Pros** | **Cons** | **Best Use Case** |
|---|---|---|---|---|
| Protobuf Definition Language | Core definition language for gRPC APIs | * Built-in, human-readable * Central source of truth | * Less visual compared to OpenAPI | * Simple and clear API |
| gRPC Service Comments | Inline comments within Protobuf definitions | * Easy to maintain, contextually relevant | * Limited formatting options | * Document specific elements |
| Third-party Tools (like Swagger for Protobuf) | Generate API documentation in OpenAPI-like formats | * Familiar format for REST API users * Can generate code documentation | * Requires additional tools * Potential inconsistency with Protobuf definitions | * Detailed documentation or RESTful API view |


## Available Endpoints 

#### gRPC
all gRPC endpoint requires a bearer token. as of now and only for dake of this project it is hardcoded to "i am not a hacker!" and later on it should be change to a proper Oauth/OIDC verification!


the input validation is manual right now but input should thoroughly validate gRPC request data types, formats, ranges, and business rules to ensure data integrity and prevent security vulnerabilities.
for that https://github.com/bufbuild/protovalidate can be a good option and there is already a middleware for grpc server as well.
```python
/create

Fields:
    name (string 3-50)
    user (string 3-50)
    password (string 3-50)
    databaseName (string 3-50)

this will creates a postgres database - besides the avaivle parameters all other values for creating the deployment is hardcoded. the namespace will be default
```

```python
/update

Fields:
    name (string 3-50)
    replicas (int32 1-3)

this will change the number of replicas
```

```python
/delete

Fields:
    name (string 3-50)

this will delete the deployment
```

```python
/check

Fields:
    service (string 3-50)

check the status of the deployment
```

#### HTTP


```python
/healthz

Fields:

health endpoint
```

```python
/metrics

Fields:

prometheus metrics endpoint
```

## Tools

for sending requests and utilizing the grpc server you can use Postman. just upload the proto file and you can send the requests.

![postman](/assets/img1.png "postman")

if you want to connect the postgres database which has been created, don't forget to port-forward the 5432 port.

```bash
kubectl port-forward -n default svc/service-name 5432:5432
```

## Contributing
Pull requests are welcome. feel free to add comments and improvement request and I will change them ASAP.


## License
[MIT](https://choosealicense.com/licenses/mit/)
