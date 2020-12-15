# Plugin style is the right way with go_out currently

#protoc -I=api/interestcal --go_out=plugins=grpc:third_party/interestcalpb interestcalservice.proto


protoc --go_out=third_party/interestcalpb --go_opt=paths=source_relative \
 --go-grpc_out=third_party/interestcalpb --go-grpc_opt=paths=source_relative \
 api/interestcal/intrestcalservice.proto

