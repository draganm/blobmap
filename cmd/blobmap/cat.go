package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/draganm/blobmap"
	"github.com/klauspost/compress/zstd"
	"github.com/urfave/cli/v2"
)

func cat() *cli.Command {
	cfg := struct {
		Zstd bool
		Gzip bool
	}{}
	return &cli.Command{
		Name:      "cat",
		Usage:     "cat a value from a blobmap file",
		Args:      true,
		UsageText: `<blobmap file> <key>  shows the value for the given key`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "zstd",
				Aliases:     []string{"z"},
				Usage:       "decompress the value using zstd",
				Value:       false,
				Destination: &cfg.Zstd,
			},
			&cli.BoolFlag{
				Name:        "gzip",
				Aliases:     []string{"g"},
				Usage:       "decompress the value using gzip",
				Value:       false,
				Destination: &cfg.Gzip,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 1 {
				return fmt.Errorf("no blobmap file provided")
			}
			blobmapFile := c.Args().First()

			if c.Args().Len() < 2 {
				return fmt.Errorf("no key provided")
			}

			keyString := c.Args().Get(1)

			key, err := strconv.ParseUint(keyString, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid key: %w", err)
			}

			blobmap, err := blobmap.Open(blobmapFile)
			if err != nil {
				return fmt.Errorf("failed to open blobmap file: %w", err)
			}

			value, err := blobmap.Read(key)
			if err != nil {
				return fmt.Errorf("failed to get value for the key %d: %w", key, err)
			}

			ln := len(value)
			fmt.Printf("Value length: %d\n", ln)

			if cfg.Zstd {
				r, err := zstd.NewReader(bytes.NewReader(value))
				if err != nil {
					return fmt.Errorf("failed to decompress value for the key %d: %w", key, err)
				}
				defer r.Close()
				_, err = io.Copy(os.Stdout, r)
				if err != nil {
					return fmt.Errorf("failed to write decompressed value to stdout for the key %d: %w", key, err)
				}
			}

			if cfg.Gzip {
				gr, err := gzip.NewReader(bytes.NewReader(value))
				if err != nil {
					return fmt.Errorf("failed to decompress value for the key %d: %w", key, err)
				}
				defer gr.Close()
				_, err = io.Copy(os.Stdout, gr)
				if err != nil {
					return fmt.Errorf("failed to write decompressed value to stdout for the key %d: %w", key, err)
				}
			}

			_, err = io.Copy(os.Stdout, bytes.NewReader(value))
			if err != nil {
				return fmt.Errorf("failed to write value to stdout for the key %d: %w", key, err)
			}

			return nil
		},
	}
}
