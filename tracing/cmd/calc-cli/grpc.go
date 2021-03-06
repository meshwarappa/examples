package main

import (
	"fmt"
	"os"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	cli "goa.design/examples/basic/gen/grpc/cli/calc"
	"goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/grpc/middleware/xray"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
)

func doGRPC(scheme, host string, timeout int, debug bool) (goa.Endpoint, interface{}, error) {
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			// Mount the X-Ray and trace client middleware. Order is very important.
			xray.UnaryClient(host),
			middleware.UnaryClientTrace(),
		)),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("could not connect to gRPC server at %s: %v", host, err))
	}
	return cli.ParseEndpoint(conn)
}
