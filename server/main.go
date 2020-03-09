package main

import (
	"context"
	pb "github.com/dannyhinshaw/go-crawler/pb_crawler"
	"google.golang.org/grpc"
	"log"
	"net"
	c "server/crawler"
)

const port = ":50051"

type server struct {
	pb.UnimplementedCrawlerServer
}

var crawler = &c.Crawler{
	Run:  true,
	Urls: map[string]struct{}{},
}

// CrawlerStart - implements pb_crawler.StartCrawler
func (s *server) CrawlerStart(ctx context.Context, in *pb.StartRequest) (*pb.ControlResponse, error) {
	url := in.GetUrl()
	log.Printf("Starting crawler: %v", url)

	// Start the crawler
	go crawler.Crawl(url)
	return &pb.ControlResponse{Started: true}, nil
}

// CrawlerStop - implements pb_crawler.CrawlerStop
func (s *server) CrawlerStop(ctx context.Context, in *pb.StopRequest) (*pb.ControlResponse, error) {
	log.Println("Stopping crawler...")
	crawler.Run = false
	return &pb.ControlResponse{Started: false}, nil
}

// CrawlerList - implements pb_crawler.ListTree
func (s *server) ListTree(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	log.Println("Getting site tree...")
	urls := crawler.Urls
	keys := make([]string, len(urls))

	i := 0
	for k := range urls {
		keys[i] = k
		i++
	}

	return &pb.ListResponse{Tree: c.BuildTree(keys)}, nil
}

func main() {
	log.Println("starting gRPC server...")
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
