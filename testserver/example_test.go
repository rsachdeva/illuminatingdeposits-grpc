package testserver_test

import (
	"context"
	"fmt"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/testserver"
)

func ExampleMongoDbConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, pool, resource := testserver.MongoDbConnect(ctx, 10)
	err := pool.Purge(resource)
	fmt.Println("ExampleMongoDbConnect err is", err)
	// Output:
	// ExampleMongoDbConnect err is <nil>
}
