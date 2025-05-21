package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zuhrulumam/grpc-file-upload/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "grpccli",
	Short: "gRPC File Upload/Download CLI",
	Long:  `A CLI for uploading and downloading files using gRPC.`,
}

const chunkSize = 1024 * 32 // 32KB

var (
	client pb.FileServiceClient
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
}

func Execute() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	client = pb.NewFileServiceClient(conn)

	cobra.CheckErr(rootCmd.Execute())
}
