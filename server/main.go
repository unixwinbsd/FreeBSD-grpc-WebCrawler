package main

import (
	"context"
	"encoding/json"
	pb "github.com/unixwinbsd/FreeBSD-grpc-WebCrawler/pb_crawler"
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
	Urls: map[string]map[string]struct{}{},
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
	log.Println("Getting site trees...")

	var allTrees [][]c.Node
	urls := crawler.Urls
	for host := range urls {
		i := 0
		currUrls := urls[host]
		keys := make([]string, len(currUrls))
		for url := range currUrls {
			keys[i] = url
			i++
		}

		allTrees = append(allTrees, c.BuildTree(keys))
	}

	b, err := json.Marshal(map[string][][]c.Node{"trees": allTrees})
	if err != nil {
		panic(err)
	}

	treeStr := string(b)
	return &pb.ListResponse{Tree: treeStr}, nil
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
