package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/zuhrulumam/grpc-file-upload/pb"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [filepath]",
	Short: "Upload a file to the gRPC server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		err := Upload(filepath)
		if err != nil {
			fmt.Println("❌ Failed:", err)
		}
	},
}

func Upload(path string) error {
	// read file locally
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fi, _ := file.Stat()
	bar := progressbar.DefaultBytes(fi.Size(), "Uploading")

	// create stream upload
	stream, err := client.Upload(context.Background())
	if err != nil {
		return err
	}

	buf := make([]byte, chunkSize)
	firstChunk := true
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		req := &pb.UploadRequest{
			Chunks: buf[:n],
		}
		if firstChunk {
			req.Filename = filepath.Base(path)
			firstChunk = false
		}

		// send to stream with chunked file
		if err := stream.Send(req); err != nil {
			return err
		}

		_ = bar.Add(n)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Println("\n✅ Upload complete:", resp.Filename)
	}

	return nil
}
