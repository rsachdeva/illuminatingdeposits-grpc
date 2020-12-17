# Plugin style is the right way with go_out currently

#protoc -I=api/interestcalpb --go_out=plugins=grpc:third_party/interestcalpb interestcalservice.proto

# using https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc
# https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
api/usermgmtpb/usermgmtservice.proto
