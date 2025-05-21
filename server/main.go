package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/zuhrulumam/grpc-file-upload/pb"
	"google.golang.org/grpc"
)

const chunkSize = 1024 * 32 // 32KB

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
}

func (fs *fileServiceServer) Upload(stream pb.FileService_UploadServer) error {
	var (
		file     *os.File
		filename string
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

			defer file.Close()
		}

		// write chunks into file
		if _, err := file.Write(req.GetChunks()); err != nil {
			return err
		}
	}

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

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("e", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &fileServiceServer{})

	log.Println("listening tcp:50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalln("err", err)
	}

}
