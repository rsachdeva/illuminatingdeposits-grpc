package mongodbhealth_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn/mongotestsetup"
)

func TestGetMongoDBHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _, _ = mongotestsetup.Connect(ctx, 10)
}
