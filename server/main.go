package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/zuhrulumam/grpc-file-upload/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const chunkSize = 1024 * 32 // 32KB

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
}

func (fs *fileServiceServer) Upload(stream pb.FileService_UploadServer) error {
	var (
		file     *os.File
		filename string
		writer   *bufio.Writer
	)

	for {
		// read stream
		// if end of file return with message success
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// create file if not yet created
		if file == nil {
			filename = req.Filename
			file, err = os.Create("./uploads/" + filename)
			if err != nil {
				return err
			}

			writer = bufio.NewWriter(file)

			defer file.Close()
		}

		// write chunks into file
		if _, err := writer.Write(req.GetChunks()); err != nil {
			return err
		}
	}

	writer.Flush()

	return stream.SendAndClose(&pb.UploadResponse{
		Filename: filename,
		Success:  true,
	})

}

func (fs *fileServiceServer) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	file, err := os.Open("./uploads/" + req.Filename)
	if err != nil {
		return err
	}

	// read file per chunksize (32kb)
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// send to client
		_ = stream.Send(&pb.DownloadResponse{
			Chunks: buf[:n],
		})

	}

	return nil
}

func (fs *fileServiceServer) UploadSimple(ctx context.Context, req *pb.UploadRequestSimple) (*pb.UploadResponse, error) {
	// Example implementation — adjust as needed
	fmt.Printf("Received UploadSimple: %s\n", req.Filename)

	err := os.WriteFile("uploads/"+req.Filename, req.File, 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	return &pb.UploadResponse{
		Filename: req.Filename,
		Success:  true,
	}, nil
}

func (fs *fileServiceServer) DownloadSimple(ctx context.Context, req *pb.DownloadRequestSimple) (*pb.DownloadResponse, error) {
	// Example implementation — adjust as needed
	fmt.Printf("Received downloadSimple: %s\n", req.Filename)

	file, err := os.Open("uploads/" + req.Filename)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file: %v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadResponse{
		Chunks: data,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("e", err)
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(120*1024*1024),
		grpc.MaxSendMsgSize(120*1024*1024),
	)
	pb.RegisterFileServiceServer(s, &fileServiceServer{})

	log.Println("listening tcp:50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalln("err", err)
	}

}
