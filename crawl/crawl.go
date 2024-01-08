package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/unixwinbsd/FreeBSD-grpc-WebCrawler/pb_crawler"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type Node struct {
	Name     string `json:"name"`
	Children []Node `json:"children"`
}

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

	// Create CLI
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "start",
				Usage:       "start crawling url",
				Destination: &url,
			},
			&cli.StringFlag{
				Name:        "stop",
				Usage:       "stop crawling url",
				Destination: &url,
			},
			&cli.BoolFlag{
				Name:  "list",
				Usage: "list site trees",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {

			// Parse command flags
			set := &flag.FlagSet{}
			nc := cli.NewContext(c.App, set, c)
			flags := nc.FlagNames()
			nFlags := len(flags)
			singleFlag := nFlags == 1

			// Parsing the command with default
			targetFlag := ""
			if singleFlag {
				targetFlag = flags[0]
			}

			// List Tree Command
			if singleFlag && targetFlag == "list" {
				r, err := client.ListTree(ctx, &pb.ListRequest{})
				if err != nil {
					log.Fatalf("could not get tree: %v", err)
				}

				var tree map[string][][]Node
				json.Unmarshal([]byte(r.Tree), &tree)
				j, err := json.MarshalIndent(tree, "", "  ")
				if err != nil {
					log.Println(err)
					return nil
				}

				// Print out the JSON tree to terminal
				log.Printf("response: %v", string(j))
				return nil
			}

			// Start & Stop require url
			if url == "" {
				return errors.New("you must specify a url to crawl")
			}

			// Stop command
			if singleFlag && targetFlag == "stop" {
				r, err := client.CrawlerStop(ctx, &pb.StopRequest{})
				if err != nil {
					log.Fatalf("could not stop crawler: %v", err)
				}

				log.Printf("response: %v", r)
				return nil
			}

			// Start command
			if singleFlag && targetFlag == "start" {
				fmt.Println("crawling", url)
				r, err := client.CrawlerStart(ctx, &pb.StartRequest{Url: url})
				if err != nil {
					log.Fatalf("could not start crawler: %v", err)
				}

				log.Printf("response: %v", r)
				return nil
			}

			return errors.New("you may only pass one flag to crawl")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
