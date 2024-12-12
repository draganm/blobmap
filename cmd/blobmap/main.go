package main

import (
	"fmt"

	"github.com/draganm/blobmap"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "blobmap",
		Usage: "cli tool for inspecting blobmap files",
		Commands: []*cli.Command{
			cat(),
		},
		UsageText: `<blobmap file>  shows the information about the blobmap file`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("no blobmap file provided")
			}
			blobmapFile := c.Args().First()
			blobmap, err := blobmap.Open(blobmapFile)
			if err != nil {
				return fmt.Errorf("failed to open blobmap file: %w", err)
			}

			defer blobmap.Close()

			fmt.Printf("File: %s\n", blobmapFile)
			fmt.Printf("First key: %d\n", blobmap.FirstKey())
			fmt.Printf("Last key: %d\n", blobmap.LastKey())
			fmt.Printf("Key count: %d\n", blobmap.LastKey()-blobmap.FirstKey()+1)

			return nil
		},
	}
	app.RunAndExitOnError()
}
