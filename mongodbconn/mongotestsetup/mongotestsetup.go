package mongotestsetup

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

// Basaed on Dockertest
func Connect(ctx context.Context, timeoutSec int) (*mongo.Client, *dockertest.Pool, *dockertest.Resource) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "4.4.2", nil)
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}
	var mt *mongo.Client
	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		log.Printf("uri is %v", uri)
		mt, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(time.Duration(timeoutSec)*time.Second))
		if err != nil {
			log.Fatal(err)
		}

		// Ping the primary
		// err = mt.Ping(ctx, readpref.Primary())
		// log.Printf("Pinging response is %v", err)
		// return err
		return mt.Ping(ctx, readpref.Primary())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return mt, pool, resource

}
