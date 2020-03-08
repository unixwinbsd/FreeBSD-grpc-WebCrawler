package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var url string

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
			if url == "" {
				fmt.Println("You must specify a url to crawl.")
			} else {
				fmt.Println("Crawling", url)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
