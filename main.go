package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nicobao/snakelet"
	"github.com/urfave/cli/v3"
)

type Config struct {
	// Database DatabaseConfig `mapstructure:"database"`
	Nostr NostrConfig `mapstructure:"nostr"`
}

// type DatabaseConfig struct {
// 	User     string `mapstructure:"user" validate:"required"`
// 	Password string `mapstructure:"password" validate:"required"`
// 	Host     string `mapstructure:"host" validate:"required"`
// 	Port     int    `mapstructure:"port" validate:"required"`
// 	Name     string `mapstructure:"name" validate:"required"`
// 	SslMode  string `mapstructure:"sslMode"`
// 	MaxConns int    `mapstructure:"maxConns"`
// }

type NostrConfig struct {
	Url                 string `mapstructure:"url"`
	ProofChannelEventId string `mapstructure:"proofChannelEventId"`
}

// only set default for non-required tag
func GetDefaultConfig() *Config {
	return &Config{
		// Database: DatabaseConfig{
		// 	SslMode:  "disable",
		// 	MaxConns: 200,
		// },
		Nostr: NostrConfig{
			Url:                 "wss://nos.lol",
			ProofChannelEventId: "6cde92f2a057368e4871d5dac0f830e09f399bad666a2070e4c4c6c40d235667",
		},
	}

}

func main() {
	// Get config struct and default variables
	conf := GetDefaultConfig()

	// this will mutate the `conf` variable:
	_, err := snakelet.InitAndLoad(conf, "./watchdog.yml")
	if err != nil {
		fmt.Println("Unable to init and load config: %w", err)
		os.Exit(1)
	}

	ctx := context.Background()
	relay, err := nostr.RelayConnect(ctx, conf.Nostr.Url)
	if err != nil {
		fmt.Println("Unable to connect to Nostr relay: %w", err)
		os.Exit(1)
	}

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "listen to proof and store them in local PostgreSQL database",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tagMap := map[string][]string{
						"e": {conf.Nostr.ProofChannelEventId},
					}
					// t := make(map[string][]string)
					// // making a "p" tag for the above public key.
					// // this filters for messages tagged with the user, mainly replies.
					//
					// t["e"] = []string{conf.Nostr.ProofChannelEventId}
					var filters nostr.Filters
					filters = []nostr.Filter{{
						Kinds: []int{nostr.KindChannelMessage},
						Tags:  tagMap,
						Limit: 1,
					}}

					sub, err := relay.Subscribe(ctx, filters)
					if err != nil {
						fmt.Println("Error while subscribing to Nostr: %w", err)
						return err
					}
					defer sub.Unsub()

					fmt.Println("Listening for proof events...")
					for ev := range sub.Events {
						// handle returned event.
						// channel will stay open until the ctx is cancelled (in this case, context timeout)
						// TODO: verify integrity of proofs
						// TODO: store to database
						fmt.Println(ev.String())
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
