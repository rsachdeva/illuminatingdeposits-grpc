// Provides mongodb connection for tests only
package mongodbtestconn

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

	resource, err := pool.Run("mongo", "4.4.2-bionic", nil)
	log.Println("resource is", resource.Container.Config)
	log.Println("resource hostport is", resource.GetHostPort("27017/tcp"))
	log.Println("err ir ", err)
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}
	var mt *mongo.Client
	//pool.MaxWait = 3 * time.Second
	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		log.Printf("uri is %v", uri)
		mt, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(time.Duration(timeoutSec)*time.Second))
		if err != nil {
			log.Fatal(err)
		}
		err = mt.Ping(ctx, readpref.Primary())
		log.Printf("err when pinging for connection setup %v", err)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return mt, pool, resource

}
