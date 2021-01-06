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

### To start only external db for working with local machine Editor/IDE:
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

### gRPC Services Endpoints Invoked:

#### Sanity test Client:
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
And then run the following:
```shell
export GODEBUG=x509ignoreCN=0 
go test -v -count=1 -covermode=count -coverpkg=./userauthn,./usermgmt,./mongodbhealth,./interestcal -coverprofile cover.out ./... && go tool cover -func cover.out
go test -v -count=1 -covermode=count -coverpkg=./userauthn,./usermgmt,./mongodbhealth,./interestcal -coverprofile cover.out ./... && go tool cover -html cover.out
```
Coverage Result for covered packages:  
**total:	(statements)	87.6%**  

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
to confirm is something else is already running

# Version
v3.0