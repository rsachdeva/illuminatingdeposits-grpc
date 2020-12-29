package mongodbtestconn

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ExampleConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mt, pool, resource := Connect(ctx, 10)
	fmt.Println("Connected mongo client Pinged again check, errors:", mt.Ping(ctx, readpref.Primary()))

	// When you're done, kill and remove the container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	fmt.Println("Purged")
	// Output:
	// Connected mongo client Pinged again check, errors: <nil>
	// Purged
}

// For temp usage if want to test can connect through db ui to see data
// func ExampleConnect_nocleanup() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	mt, _, _ := Connect(ctx, 10)
// 	fmt.Println("Connected mongo client Pinged again check, errors:", mt.Ping(ctx, readpref.Primary()))
// 	fmt.Println("Not purged; manually delete, errors:", mt.Ping(ctx, readpref.Primary()))
//
// 	// Output:
// 	// Connected mongo client Pinged again check, errors: <nil>
// 	// Not purged; manually delete, errors: <nil>
// }
