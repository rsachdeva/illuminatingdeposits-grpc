package main

import "github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"

func main() {
	ctx, mt := mongodbconn.ConnectMongoDB()
	_ = mt.Database("depositsmongodb")
	mongodbconn.DisconnectMongodb(mt, ctx)
}
