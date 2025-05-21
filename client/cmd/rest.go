package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/zuhrulumam/grpc-file-upload/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var restCmd = &cobra.Command{
	Use:   "rest start",
	Short: "Start rest server",
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

func Serve() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := pb.RegisterFileServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(120*1024*1024),
			grpc.MaxCallSendMsgSize(120*1024*1024),
		)})
	if err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

	log.Println("🚀 REST gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
