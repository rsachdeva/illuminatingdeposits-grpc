# Illuminating Deposits - gRPC

###### All commands should be executed from the root directory (illuminatingdeposits-grpc) of the project
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# Features include:
- Golang (Go) gRPC Service Methods with protobuf for Messages
- TLS for all requests
- Integration and Unit tests; run in parallel using dockertest for faster feedback (more in progress)
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
- Concurrency support for computing Banks Delta included for I/O processing
- Sanity test client included each deployment settings
  With this Sanity test client, you will be able to:
  - get status of Mongo DB
  - add a new user
  - edit json file that is read to make input request to gRPC service
  - JWT generation for Authentication
  - JWT Authentication for Interest Delta Calculations for each deposit; each bank with all deposits and all banks
    Quickly confirms Sanity check for set up with Kubernetes/Docker.
    There are also separate Integration and Unit tests.
- Dockering and using it for both Docker Compose and Kubernetes
- Docker compose deployment for development
- Kubernetes Deployment with Ingress; Helm; Mongodb internal replication setup
- Running from Editor/IDE directly included
- Log Based Message Broker using Kafka for Each Interest Calculation Event with Docker Compose and Kubernetes 
- Tracing enabled using Zipkin for Observability

#### gRPC pre setup  (in case desired to generate already included protobuf code)
This should not be needed when deploying; it is there for reference to how protobuf
and related code was generated.

```
brew install protobuf
``` 
Enable module mode (or just execute next command from any directory outside of project having go.mod)
Reference: [QuickStart](https://grpc.io/docs/languages/go/quickstart/)
```shell
brew install protobuf
protoc --version  # Ensure compiler version is 3+
# install Go plugins for the protocol compiler protoc
# before go 1.16
export GO111MODULE=on  
go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc 
# for go 1.16
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Now go to this project root directory and run following scripts for generating protobuf related code for the services:
```shell
sh ./generateinterestcalservice.sh
sh ./generatemongodbhealthservice.sh
sh ./generateusermgmtservice.sh 
```

# Docker Compose Deployment

### TLS files
```shell
export DEPOSITS_GRPC_SERVICE_ADDRESS=localhost
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run --env DEPOSITS_GRPC_SERVICE_ADDRESS=$DEPOSITS_GRPC_SERVICE_ADDRESS -v $PWD/conf/tls:/tls tlscert:v0.1
```

### Start mongodb and tracing service
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up
```

### Start message broker log streaming service
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.kafka.yml up
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

### Server Tracing
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)

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
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up
```
And then run following (./tools/dbindexescli is actually needed only once):
```shell
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_DB_HOST=127.0.0.1
export DEPOSITS_TRACE_URL=http://localhost:9411/api/v2/spans
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
### Server Tracing
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)

# Kubernetes Deployment
(for Better control; For Local Setup tested with Docker Desktop latest version with Kubernetes Enabled)

### TLS files
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

### Installing Kafka
Using helm to install kafka
```shell
brew install helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm repo list
```
and then use
```shell
helm install kafka -f ./deploy/kubernetes/kafka/helm-values.yaml bitnami/kafka
```
to install ingress controller
To see logs for kafka:
```shell
kubectl logs pod/kafka-grpc-0 -f
```

This installation message also allows to test installation of kafka quickly:
(kafka-console-producer.sh and kafka-console-consumer.sh can be run inside the kafka-grpc-client pod created below)
```shell
Kafka can be accessed by consumers via port 9092 on the following DNS name from within your cluster:

    kafka-grpc.default.svc.cluster.local

Each Kafka broker can be accessed by producers via port 9092 on the following DNS name(s) from within your cluster:

    kafka-grpc-0.kafka-grpc-headless.default.svc.cluster.local:9092

To create a pod that you can use as a Kafka client run the following commands:

    kubectl run kafka-grpc-client --restart='Never' --image docker.io/bitnami/kafka:2.8.1-debian-10-r0 --namespace default --command -- sleep infinity
    kubectl exec --tty -i kafka-grpc-client --namespace default -- bash

    PRODUCER:
        kafka-console-producer.sh \

            --broker-list kafka-grpc-0.kafka-grpc-headless.default.svc.cluster.local:9092 \
            --topic test

    CONSUMER:
        kafka-console-consumer.sh \

            --bootstrap-server kafka-grpc.default.svc.cluster.local:9092 \
            --topic test \
            --from-beginning
```

### Make docker images and Push Images to Docker Hub

```shell
docker rmi rsachdeva/illuminatingdeposits.grpc.server:v1.6.0
docker rmi rsachdeva/illuminatingdeposits.dbindexes:v1.6.0

docker build -t rsachdeva/illuminatingdeposits.grpc.server:v1.6.0 -f ./build/Dockerfile.grpc.server .  
docker build -t rsachdeva/illuminatingdeposits.dbindexes:v1.6.0 -f ./build/Dockerfile.dbindexes .  

docker push rsachdeva/illuminatingdeposits.grpc.server:v1.6.0
docker push rsachdeva/illuminatingdeposits.dbindexes:v1.6.0
```
The v1.6.0 should match to version being used in Kubernetes resources (grpc-server.yaml; dbindexes.yaml).

### Quick deploy for all Kubernetes resources ( more Detailed Kubernetes set up - Step by Step is below; Above kubernetes steps are common)
TLS File set up should have been installed using steps above along with installing using help ingress-nginx and kafka

MongoDB set up with internal replication has to be done once unless persistent volumes are deleted or
number of replicas are changed:
```shell
kubectl apply -f deploy/kubernetes/mongodb/.
```
Once all mongo pods are running, we have to do this once: ( unless we change replication pods)

#### Mongodb internal replication setup
Inside the mongo-0 pod we open a mongo shell:
```shell
kubectl exec -it mongo-0 -- mongo
```
Then inside shell:
```shell
rs.initiate({
      _id: "rs0",
      members: [
         { _id: 0, host : "mongo-0.mongo.default.svc.cluster.local:27017" },
         {_id: 1, host: "mongo-1.mongo.default.svc.cluster.local:27017"},
         {_id: 2, host: "mongo-2.mongo.default.svc.cluster.local:27017"}
      ]
   }
)
#  wait till you get 1 primary
rs.status()
```
Allows connecting MongoDB UI using NodePort at 30010 from outside cluster locally to view data.
We only need to set secrets once after tls files have been generated; This is using dry run:
```shell
# only so secret can be applied with all resources as once; easier to delete with all resources
kubectl delete secret illuminatingdeposits-grpc-secret-tls
kubectl create --dry-run=client secret tls illuminatingdeposits-grpc-secret-tls --key conf/tls/serverkeyto.pem --cert conf/tls/servercrtto.pem -o yaml > ./deploy/kubernetes/tls-secret-ingress.yaml
```
Now lets depoly the grpc application related resources:
```shell
# in case docker image not built -- refer 'Make docker images and Push Images to Docker Hub' above
kubectl apply -f deploy/kubernetes/.
```
If status for ```kubectl get pod -l job-name=dbindexes | grep "Completed"```
shows completed for dbindexes pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/dbindexes.yaml
```
And see logs using
```kubectl logs -l app=grpcserversvc -f```

### Sanity test Client - gRPC Services Endpoints Invoked Externally:
The server side DEPOSITS_GRPC_SERVICE_TLS should be consistent and set for client also.
Uncomment any request function if not desired.
As stated before, All commands should be executed from the root directory (illuminatingdeposits-grpc) of the project
```shell
export GODEBUG=x509ignoreCN=0
export DEPOSITS_GRPC_SERVICE_TLS=true
export DEPOSITS_GRPC_SERVICE_ADDRESS=grpcserversvc.127.0.0.1.nip.io
go run ./cmd/sanitytestclient
```

### Server Tracing
Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io](http://zipkin.127.0.0.1.nip.io)

### Detailed Kubernetes set up - Step by Step
See 'Quick deploy for all Kubernetes resources' above. If any issues, you can follow detailed steps.
TLS File set up and Ingress controller should have been installed using steps above.
##### Start mongodb service

```shell
kubectl apply -f deploy/kubernetes/mongodb/mongodb.yaml
```
See above 'Mongodb internal replication setup' in Quick deploy.

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
```shell
kubectl logs -l app=grpcserversvc -f
```
  
### Remove all resources / Shutdown

```shell
kubectl delete -f ./deploy/kubernetes/.
kubectl delete -f ./deploy/kubernetes/mongodb/.
helm uninstall ingress-nginx
helm uninstall kafka
```
This does not remove prestient volume.
In case you are sure, you don't want to keep data only then
```shell
kubectl get pvc
kubectl delete pvc -l app=mongo 
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
**total:	(statements)	92.5%**  
<p align="center">
<img src="./coverageresults.png" alt="Illuminating Deposits gRPC Test Coverage" title="lluminating Deposits gRPC Test Coverage" />
</p>

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
For running single test for debugging purpose:
```shell
go test -v -count=1 ./mongodbhealth -run TestServiceServer_HealthOk
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
Otherwise you will get error as:
```
transport: Error while dialing dial tcp
```

# Version
v1.6.0