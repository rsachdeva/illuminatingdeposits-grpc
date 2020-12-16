# Illuminating Deposits - gRPC

# gRPC pre setup 

```
brew install protobuf
``` 
Enable module mode (or just execute next command from any directory outside of project having go.mod)
```
export GO111MODULE=on 
go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc 
```

Now go to this project root directory.

# All commands should be executed from the root directory (illuminatingdeposits-grpc) of the project
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# gRPC API using protobuf for Messages

Proto file exists in api folder; accordingly run the following to generate grpc code for the service
```
run generateinterestcalservice.sh
go mod tidy  
```

# Version
v0.81