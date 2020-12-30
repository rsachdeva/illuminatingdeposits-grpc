package mongodbtestconn_test

import (
	"context"
	"fmt"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongodbtestconn"
)

func ExampleConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, pool, resource := mongodbtestconn.Connect(ctx, 10)
	err := pool.Purge(resource)
	fmt.Println("err is", err)
	// Output:
	// err is <nil>
}
