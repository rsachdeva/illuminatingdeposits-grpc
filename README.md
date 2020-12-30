# Illuminating Deposits - gRPC

#### gRPC pre setup  (in case desired to generate already included protobuf code)

```
brew install protobuf
``` 
Enable module mode (or just execute next command from any directory outside of project having go.mod)
```
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

### To start only external db and trace service for working with local machine Editor/IDE:
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

# Running Integration/Unit tests
```shell 
docker pull mongo:4.4.2-bionic (only once as tests use this image; so faster)
``` 
and then do:
```shell
go test -v ./... 
```
The -count=1 is mainly to not use caching and can be added as follows if needed for 
any go test command:
```shell 
go test -v -count=1 ./...
```
For coverage for per package:
```shell
go test -cover ./...
```
To see coverage stats in more detail:
```shell 
go test -v -coverprofile cover.out ./...
```
To see overall total coverage
```shell 
go tool cover -func=cover.out
```
or
```shell 
go tool cover -func cover.out
```
To see in the browser the covered parts:
```shell 
go tool cover -html cover.out
```
Single package:
```shell 
go test -v -coverprofile cover.out ./mongodbhealth && go tool cover -func cover.out
```
To run tests with coverage and see reports in excluding certain packages not needed:
In grep -v means "invert the match" in grep, in other words, return all non-matching lines
```shell 
go test -v -count=1 -coverprofile cover.out $(go list ./... | grep -v mongodbconn | grep -v pb) && go tool cover -func cover.out
go test -v -count=1 -coverprofile cover.out $(go list ./... | grep -v mongodbconn) && go tool cover -html cover.out
```
See Editor specifcs to see Covered Parts in the Editor.
Docker containers are mostly auto removed.
There could be a container to inspect data.
In case any docker containers still running after tests:

```shell 
docker stop $(docker ps -qa)
docker rm -f $(docker ps -qa)
```
And if mongodb not connecting for tests: (reference: https://www.xspdf.com/help/52284027.html)
```shell 
docker volume rm $(docker volume ls -qf dangling=true)
```
# Version
v2.35