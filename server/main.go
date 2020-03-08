package main

import (
	"context"
	"log"
	"net"
	"server/crawler"

	pb "github.com/dannyhinshaw/go-crawler/pb_crawler"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedCrawlerServer
}

// CrawlerStart - implements pb_crawler.StartCrawler
func (s *server) CrawlerStart(ctx context.Context, in *pb.StartRequest) (*pb.ControlResponse, error) {
	url := in.GetUrl()
	log.Printf("Received: %v", url)

	crawler.StartCrawl(url)
	return &pb.ControlResponse{Started: true}, nil
}

// CrawlerStop - implements pb_crawler.CrawlerStop
func (s *server) CrawlerStop(ctx context.Context, in *pb.StopRequest) (*pb.ControlResponse, error) {
	return &pb.ControlResponse{Started: false}, nil
}

// CrawlerList - implements pb_crawler.CrawlerList
func (s *server) CrawlerList(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	return &pb.ListResponse{Urls: ""}, nil
}

func main() {
	log.Println("starting grpc server...")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCrawlerServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
