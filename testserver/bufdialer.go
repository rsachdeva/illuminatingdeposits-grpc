package testserver

import (
	"context"
	"net"

	"google.golang.org/grpc/test/bufconn"
)

func GetBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
