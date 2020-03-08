package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "github.com/dannyhinshaw/go-crawler/pb_crawler"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Init gRPC client
	client := pb.NewCrawlerClient(conn)

	// Define CLI
	var url string
	var cmdPrint = &cobra.Command{
		Use:   "crawl -start [url to crawl]",
		Short: "Crawl a website by url",
		Long: `crawl will start a crawler at the given url and find all links belonging to that domain.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			fmt.Println("crawling", url)
			r, err := client.CrawlerStart(context.Background(), &pb.StartRequest{Url: url})
			if err != nil {
				log.Fatalf("could not start crawler: %v", err)
			}
			log.Printf("response: %v", r)
		},
	}

	var cmdEcho = &cobra.Command{
		Use:   "crawl -stop [url to stop crawling]",
		Short: "Stop crawling a website by url",
		Long: `crawl -stop will stop a currently running crawler for a given url.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	cmdTimes.Flags().StringVar(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdPrint, cmdEcho)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}
