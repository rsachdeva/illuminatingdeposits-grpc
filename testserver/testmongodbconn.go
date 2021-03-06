// Provides mongodb connection for tests only
package testserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Based on Dockertest
func MongoDbConnect(ctx context.Context, timeoutSec int) (*mongo.Client, *dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// for singlle test if predetrmined port for easier db ui connection; manual clean up required; canot run in parallel restriction
	// as new docker container with mongodb for each connection; parallel allows faster test execution
	// opts := dockertest.RunOptions{
	// 	Repository:   "mongo",
	// 	Tag:          "4.4.2-bionic",
	// 	ExposedPorts: []string{"27017"},
	// 	PortBindings: map[docker.Port][]docker.PortBinding{
	// 		"27017": {
	// 			{HostIP: "127.0.0.1", HostPort: "27017"},
	// 		},
	// 	},
	// }
	// resource, err := pool.RunWithOptions(&opts)

	resource, err := pool.Run("mongo", "4.4.2-bionic", nil)
	if resource == nil {
		log.Fatalf("Resource issue in dockertest: %v", err)
	}
	log.Println("resource hostport is", resource.GetHostPort("27017/tcp"))
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}
	var mt *mongo.Client
	// pool.MaxWait = 3 * time.Second
	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		log.Printf("mongodb uri is %v", uri)
		mt, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(time.Duration(timeoutSec)*time.Second))
		if err != nil {
			log.Fatal(err)
		}
		err = mt.Ping(ctx, readpref.Primary())
		log.Printf("err when pinging for mongodb connection test setup %v", err)
		return err
	}); err != nil {
		log.Fatalf("could not connect to docker for mongodb connection : %s", err)
	}
	return mt, pool, resource
}
