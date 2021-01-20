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
### To start all services with TLS:
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

### Shutdown
```shell
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml down
docker-compose -f ./deploy/compose/docker-compose.grpc.server.yml down
```

# Kubernetes Deployment
(for Better control; For Local Setup tested with Docker Desktop latest version with Kubernetes Enabled)

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
docker build -t rsachdeva/illuminatingdeposits.grpc.server:v1.3.01 -f ./build/Dockerfile.grpc.server .  
docker build -t rsachdeva/illuminatingdeposits.dbindexes:v1.3.01 -f ./build/Dockerfile.dbindexes .  

docker push rsachdeva/illuminatingdeposits.grpc.server:v1.4.0
docker push rsachdeva/illuminatingdeposits.dbindexes:v1.4.0
``` 

### Quick deploy for all resources
```shell
kubectl apply -f deploy/kubernetes/.
```
If status for ```kubectl get pod -l job-name=seed | grep "Completed"```
shows completed for seed pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/seed.yaml
```

### Detailed - Step by Step

##### Start mongodb service

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mongodb -f ./deploy/kubernetes/mongodb/helm-values.yaml bitnami/mongodb
```
To see values applied see:
```shell
helm template mongodb -f ./deploy/kubernetes/mongodb/helm-values.yaml bitnami/mongodb
```


#### Then Migrate and set up seed data manually for more control initially:
First should see in logs
database system is ready to accept connections
```kubectl logs pod/postgres-deposits-0```
And then execute migration/seed data for manual control when getting started:
```shell
kubectl apply -f deploy/kubernetes/seed.yaml
```
And if status for ```kubectl get pod -l job-name=seed | grep "Completed"```
shows completed for seed pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/seed.yaml
```
To connect external tool with postgres to see database internals use:
Use a connection string similar to:
jdbc:postgresql://127.0.0.1:30007/postgres
If still an issue you can try
kubectl port-forward service/postgres 5432:postgres
Now can easily connect using
jdbc:postgresql://localhost:5432/postgres

#### Illuminating deposists gRPC server in Kubernetes!
```shell
kubectl apply -f deploy/kubernetes/grpc-server.yaml
```
And see logs using
```kubectl logs -l app=grpcserversvc -f```



### Remove all resources / Shutdown

```shell
kubectl delete -f ./deploy/kubernetes/.
helm uninstall ingress-nginx
```

# Sanity test Client - gRPC Services Endpoints Invoked Externally:
The server side DEPOSITS_GRPC_SERVICE_TLS should be consistent and set for client also.
Uncomment any request function if not desired.

```shell 
export GODEBUG=x509ignoreCN=0
go run ./cmd/sanitytestclient
```

# TLS files
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run -v $PWD/conf/tls:/tls tlscert:v0.1
``` 

To see openssl version being used in Docker:
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

# Editor/IDE development without docker/docker compose/kubernetes as described above
Execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml up
```
And then run following:
```shell
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_DB_HOST=127.0.0.1
go run ./tools/dbindexescli  (only once)
go run ./cmd/server
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
to confirm is something else is already running

# Version
v1.3.01