package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/zuhrulumam/grpc-file-upload/pb"
)

var downloadCmd = &cobra.Command{
	Use:   "download [filename] [output_path]",
	Short: "Download a file from the gRPC server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_ = Download(args[0], args[1])
	},
}

func Download(filename, outputPath string) error {

	// request to server
	stream, err := client.Download(context.Background(), &pb.DownloadRequest{
		Filename: filename,
	})
	if err != nil {
		return err
	}
	// create file
	file, err := os.Create(outputPath + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	bar := progressbar.Default(-1, "Downloading")

	// write stream to created file
	for {
		// read stream
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// write stream to file created
		if _, err := writer.Write(resp.Chunks); err != nil {
			return err
		}

		_ = bar.Add(len(resp.Chunks))
	}

	writer.Flush()

	log.Println("Download completed.")

	return nil
}
