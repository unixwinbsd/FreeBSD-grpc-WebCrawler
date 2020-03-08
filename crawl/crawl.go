package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/dannyhinshaw/go-crawler/pb_crawler"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	var url string

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Init gRPC client
	client := pb.NewCrawlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "start",
				Usage:       "start crawling url",
				Destination: &url,
			},
			&cli.StringFlag{
				Name:  "stop",
				Usage: "stop crawling url",
			},
			&cli.StringFlag{
				Name:  "list",
				Usage: "list site tree",
			},
		},
		Action: func(c *cli.Context) error {
			log.Printf("c:: %v", c)
			if url == "" {
				fmt.Println("You must specify a url to crawl.")
			} else {
				fmt.Println("Crawling", url)
				r, err := client.CrawlerStart(ctx, &pb.StartRequest{Url: url})
				if err != nil {
					log.Fatalf("could not start crawler: %v", err)
				}
				log.Printf("response: %v", r)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
