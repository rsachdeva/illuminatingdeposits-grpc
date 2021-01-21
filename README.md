# Illuminating Deposits - gRPC

#### gRPC pre setup  (in case desired to generate already included protobuf code)

```
brew install protobuf
``` 
Enable module mode (or just execute next command from any directory outside of project having go.mod)
Reference: https://grpc.io/docs/languages/go/quickstart/
```shell
brew install protobuf
protoc --version  # Ensure compiler version is 3+
# install Go plugins for the protocol compiler protoc
export GO111MODULE=on  
go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc 
```

Now go to this project root directory and run following scripts for generating protobuf related code for the services:
```
generateinterestcalservice.sh
generatemongodbhealthservice.sh
generateusermgmtservice.sh 
```

###### All commands should be executed from the root directory (illuminatingdeposits-grpc) of the project
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# Features include:
- Golang (Go) gRPC Service Methods with protobuf for Messages
- TLS for all requests
- Integration and Unit tests run in parallel
- Coverage Result for key packages
- MongoDB health check service
- User Management service with MongoDB for user creation
- JWT generation for Authentication 
- JWT Authentication for Interest Calculations
- 30daysInterest for a deposit is called Delta
- Delta is for 
     - each deposit 
     - each bank with all deposits
     - all banks!
- Sanity test client included
- Docker support 
- Docker compose deployment for development

# Docker Compose Deployment

# TLS files
```shell
export DEPOSITS_GRPC_SERVICE_ADDRESS=localhost
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run --env DEPOSITS_GRPC_SERVICE_ADDRESS=$DEPOSITS_GRPC_SERVICE_ADDRESS -v $PWD/conf/tls:/tls tlscert:v0.1
```

# Start mongodb
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml up 
```

### Then set up mongodb indexes
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.dbindexes.yml up --build
````

### To start all services without TLS:
Make sure DEPOSITS_GRPC_SERVICE_TLS=false in docker-compose.grpc.server.yml
### To start all services with TLS - this is the way to go in gRPC:
Make sure DEPOSITS_GRPC_SERVICE_TLS=true in docker-compose.grpc.server.yml
InterestCal service only supports TLS in gRPC as it carries access token.
This means users want to transmit security
information (e.g., OAuth2 token) which requires secure connection.Otherwise will see error as:
the credentials require transport level security (use grpc.WithTransportCredentials() to set).
### And then execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.grpc.server.yml up --build
```

### Logs of running services (in a separate terminal):
```shell
docker-compose -f ./deploy/compose/docker-compose.grpc.server.yml logs -f --tail 1  
``` 
### Sanity test Client - gRPC Services Endpoints Invoked Externally:
The server side DEPOSITS_GRPC_SERVICE_TLS should be consistent and set for client also.
Uncomment any request function if not desired.

```shell
export GODEBUG=x509ignoreCN=0
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_SERVICE_ADDRESS=localhost
go run ./cmd/sanitytestclient
```

### Shutdown
```shell
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml down
docker-compose -f ./deploy/compose/docker-compose.grpc.server.yml down
```
# Runing from Editor/IDE

### TLS files -same as in Docker compose
```shell
export DEPOSITS_GRPC_SERVICE_ADDRESS=localhost
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run --env DEPOSITS_GRPC_SERVICE_ADDRESS=$DEPOSITS_GRPC_SERVICE_ADDRESS -v $PWD/conf/tls:/tls tlscert:v0.1
```
### Start DB:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml up
```
And then run following (./tools/dbindexescli is actually needed only once):
```shell
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_DB_HOST=127.0.0.1
go run ./tools/dbindexescli
go run ./cmd/server
```

### Sanity test Client - gRPC Services Endpoints Invoked Externally:
The server side DEPOSITS_GRPC_SERVICE_TLS should be consistent and set for client also.
Uncomment any request function if not desired.

```shell
export GODEBUG=x509ignoreCN=0
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_SERVICE_ADDRESS=localhost
go run ./cmd/sanitytestclient
```

# Kubernetes Deployment
(for Better control; For Local Setup tested with Docker Desktop latest version with Kubernetes Enabled)

# TLS files
```shell
export DEPOSITS_GRPC_SERVICE_ADDRESS=grpcserversvc.127.0.0.1.nip.io
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run --env DEPOSITS_GRPC_SERVICE_ADDRESS=$DEPOSITS_GRPC_SERVICE_ADDRESS -v $PWD/conf/tls:/tls tlscert:v0.1
```
As a side note, For any troubleshooting, To see openssl version being used in Docker:
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run -ti -v $PWD/conf/tls:/tls tlscert:v0.1 sh
```
You get a prompt
/tls
Check version using command:
```shell
openssl version
```

### Installing Ingress controller
Using helm to install nginx ingress controller
```shell
brew install helm
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm repo list
```
and then use
```shell
helm install ingress-nginx -f ./deploy/kubernetes/nginx-ingress-controller/helm-values.yaml ingress-nginx/ingress-nginx
```
to install ingress controller
To see logs for nginx ingress controller:
```shell
kubectl logs -l app.kubernetes.io/name=ingress-nginx -f
```

### Make docker images and Push Images to Docker Hub

```shell
docker build -t rsachdeva/illuminatingdeposits.grpc.server:v1.4.0 -f ./build/Dockerfile.grpc.server .  
docker build -t rsachdeva/illuminatingdeposits.dbindexes:v1.4.0 -f ./build/Dockerfile.dbindexes .  

docker push rsachdeva/illuminatingdeposits.grpc.server:v1.4.0
docker push rsachdeva/illuminatingdeposits.dbindexes:v1.4.0
```

### Quick deploy for all resources
We only need to set secrets once after tls files have been generated
```shell
kubectl delete secret illuminatingdeposits-grpc-secret-tls
kubectl create --dry-run=client secret tls illuminatingdeposits-grpc-secret-tls --key conf/tls/serverkeyto.pem --cert conf/tls/servercrtto.pem -o yaml > ./deploy/kubernetes/tls-secret-ingress.yaml
```
Now lets depoly:
```shell
kubectl apply -f deploy/kubernetes/.
```
If status for ```kubectl get pod -l job-name=dbindexes | grep "Completed"```
shows completed for dbindexes pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/dbindexes.yaml
```

### Detailed - Step by Step

##### Start mongodb service

```shell
kubectl apply -f deploy/kubernetes/mongodb.yaml
```

#### Then Migrate and set up dbindexes data manually for more control initially:
First should see in logs
database system is ready to accept connections
```kubectl logs pod/mongodb-deposits-0```
And then execute migration/dbindexes data for manual control when getting started:
```shell
kubectl apply -f deploy/kubernetes/dbindexes.yaml
```
And if status for ```kubectl get pod -l job-name=dbindexes | grep "Completed"```
shows completed for dbindexes pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/dbindexes.yaml
```
To connect external tool with mongodb to see database internals use port 30010

#### Set up secret

```shell
kubectl delete secret illuminatingdeposits-grpc-secret-tls
kubectl create --dry-run=client secret tls illuminatingdeposits-grpc-secret-tls --key conf/tls/serverkeyto.pem --cert conf/tls/servercrtto.pem -o yaml > ./deploy/kubernetes/tls-secret-ingress.yaml
kubectl apply -f deploy/kubernetes/tls-secret-ingress.yaml
```

#### Illuminating deposists gRPC server in Kubernetes!
```shell
kubectl apply -f deploy/kubernetes/grpc-server.yaml
```
And see logs using
```kubectl logs -l app=grpcserversvc -f```

### Sanity test Client - gRPC Services Endpoints Invoked Externally:
The server side DEPOSITS_GRPC_SERVICE_TLS should be consistent and set for client also.
Uncomment any request function if not desired.

```shell
export GODEBUG=x509ignoreCN=0
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_SERVICE_ADDRESS=grpcserversvc.127.0.0.1.nip.io
go run ./cmd/sanitytestclient
```
With this Sanity test client, you will be able to:
- get status of Mongo DB
- add a new user
- JWT generation for Authentication
- JWT Authentication for Interest Delta Calculations for each deposit; each bank with all deposits and all banks
  Quickly confirms Sanity check for set up with Kubernetes/Docker.
  There are also separate Integration and Unit tests.
  
### Remove all resources / Shutdown

```shell
kubectl delete -f ./deploy/kubernetes/.
helm uninstall ingress-nginx
```

# Running Integration/Unit tests
Tests are designed to run in parallel with its own test server and docker based mongodb using dockertest.
To run all tests with coverages reports for focussed packages:
Run following only once as tests use this image; so faster:
```shell 
docker pull mongo:4.4.2-bionic 
``` 
And then run the following with coverages for key packages concerned:
```shell
export GODEBUG=x509ignoreCN=0 
go test -v -count=1 -covermode=count -coverpkg=./userauthn,./usermgmt,./mongodbhealth,./interestcal -coverprofile cover.out ./... && go tool cover -func cover.out
go test -v -count=1 -covermode=count -coverpkg=./userauthn,./usermgmt,./mongodbhealth,./interestcal -coverprofile cover.out ./... && go tool cover -html cover.out
```
Coverage Result for key packages: 
**total:	(statements)	91.9%**  

The -v is for Verbose output: log all tests as they are run. Search "FAIL:" in parallel test output here to see reason for failure
in case any test fails.
Just to run all easily with verbose ouput:
```shell
go test -v ./... 
```
The -count=1 is mainly to not use caching and can be added as follows if needed for
any go test command:
```shell 
go test -v -count=1 ./...
```
See Editor specifcs to see Covered Parts in the Editor.
#### Test Docker containers for Mongodb
Docker containers are mostly auto removed. This is done by passing true to testserver.InitGRPCServerBuffConn(ctx, t, false)
in your test.
If you want to examine mongodb data for a particular test, you can temporarily
set allowPurge as false in testserver.InitGRPCServerBuffConn(ctx, t, false) for your test.
Then after running specific failed test connect to mongo db in the docker container using any db ui.
As an example, if you want coverage on a specific package and run a single test in a package with verbose output:
```shell 
go test -v -count=1 -covermode=count -coverpkg=./usermgmt -coverprofile cover.out -run=TestServiceServer_CreateUser ./usermgmt/... && go tool cover -func cover.out
```
Any docker containers still running after tests should be manually removed:
```shell 
docker ps
docker stop $(docker ps -qa)
docker rm -f $(docker ps -qa)
```
And if mongodb not connecting for tests: (reference: https://www.xspdf.com/help/52284027.html)
```shell 
docker volume rm $(docker volume ls -qf dangling=true)
```

# Troubleshooting
If for any reason no connection is happening from client to server or client hangs or server start up issues:
Run 
```
ps aux | grep "go run" 
```
or
```
ps aux | grep "go_build" 
```
to confirm is something else is already running.
Make sure to follow above TLS set up according to Kubernetes deployment, Docker compose deployment or Running from Editor.
Make sure to follow Ingress controller installation for Kubernetes deployment.

# Version
v1.4.0