package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	pb "github.com/peterouob/golang_template/api/protobuf"
	"google.golang.org/grpc"
)

type videoServer struct {
	pb.UnimplementedVideoServer
	storageDir string
}

func (s *videoServer) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	log.Printf("Received UploadVideo: video_data length=%d, filename=%s", len(req.VideoData), req.FileName)
	if len(req.VideoData) == 0 {
		return nil, fmt.Errorf("no video data")
	}
	videoID := uuid.New().String()
	filePath := filepath.Join(s.storageDir, videoID+"_"+req.FileName)
	err := os.WriteFile(filePath, req.VideoData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save video: %v", err)
	}
	log.Printf("Saved video %s", filePath)
	return &pb.UploadVideoResponse{VideoId: videoID}, nil
}

func (s *videoServer) GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
	files, err := os.ReadDir(s.storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage: %v", err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), req.VideoId+"_") {
			data, err := os.ReadFile(filepath.Join(s.storageDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read video: %v", err)
			}
			return &pb.GetVideoResponse{
				VideoData: data,
				FileName:  strings.TrimPrefix(file.Name(), req.VideoId+"_"),
			}, nil
		}
	}
	return nil, fmt.Errorf("video not found")
}

func main() {
	storageDir := "./videos"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Fatalf("failed to create storage directory: %v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterVideoServer(s, &videoServer{storageDir: storageDir})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
