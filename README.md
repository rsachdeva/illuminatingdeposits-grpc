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


### All commands should be executed from the root directory (illuminatingdeposits-grpc) of the project
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# gRPC API with protobuf for Messages
# Docker Compose Deployment

# Start mongodb
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-only.yml up

### Then set up mongodb indexes
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.dbindexes.yml up --build
````

### To start all services without TLS:
Make sure DEPOSITS_WEB_SERVICE_SERVER_TLS=false in docker-compose.grpc.server.yml
### To start all services with TLS:
Make sure DEPOSITS_WEB_SERVICE_SERVER_TLS=true in docker-compose.grpc.server.yml
### And then execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.grpc.server.yml up --build

### Logs of running services (in a separate terminal):
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
export DEPOSITS_WEB_SERVICE_SERVER_TLS=true
export DEPOSITS_GRPS_DB_HOST=127.0.0.1
go run ./tools/dbindexescli  (only once)
go run ./cmd/server
```

### gRPC Services Endpoints Invoked:

#### Sanity test Client:
See    
cmd/sanitytestclient/main.go  
In main function -- change tls setting to true or false.
The server side main function should also be consistent with tls setting.
Uncomment any desired function request.
Make sure to make email unique to avoid error.

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

And enter version check
```shell
openssl version
```

### Troubleshooting
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
v0.97