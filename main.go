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
func GetDefaultConfig(defaultProofChannelEventId string) *Config {
	return &Config{
		// Database: DatabaseConfig{
		// 	SslMode:  "disable",
		// 	MaxConns: 200,
		// },
		Nostr: NostrConfig{
			Url:                 "wss://nos.lol",
			ProofChannelEventId: defaultProofChannelEventId,
		},
	}

}

func main() {
	var relay *nostr.Relay = nil
	var conf *Config = nil
	defaultProofChannelEventId := "6cde92f2a057368e4871d5dac0f830e09f399bad666a2070e4c4c6c40d235667"
	cmd := &cli.Command{
		Name:      "Agora Watchdog CLI",
		Version:   "v0.1.0",
		Copyright: "(c) 2025 ZKorum SAS",
		Usage:     "verify Agora moderation history",
		UsageText: "agora-watchdog -c ./watchdog.yml listen",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Location of the config file - default to ./watchdog.yml",
				Required: false,
				Value:    "", // will only load default values
				Sources: cli.ValueSourceChain{
					Chain: []cli.ValueSource{
						cli.EnvVar("WATCHDOG_CONFIG"),
						cli.File("./watchdog.yml"),
					},
				},
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Get config struct and default variables

			conf = GetDefaultConfig(defaultProofChannelEventId)

			// this will mutate the `conf` variable:
			_, err := snakelet.InitAndLoad(conf, "./watchdog.yml")
			if err != nil {
				fmt.Println("Unable to init and load config: %w", err)
				return nil, err
			}
			relay, err = nostr.RelayConnect(ctx, conf.Nostr.Url)
			if err != nil {
				fmt.Println("Unable to connect to Nostr relay: %w", err)
				return nil, err
			}
			fmt.Printf("Connected to %s\n", conf.Nostr.Url)
			return nil, nil
		},
		Commands: []*cli.Command{

			{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "listen to proofs, verify and store them in a local PostgreSQL database",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "proof-channel-event-id",
						Usage:    "NIP 40 Event ID for the Proof Channel, used by Agora to broadcast data in the Nostr relay.",
						Required: false,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					proofChannelEventId := cmd.String("proof-channel-event-id")
					var actualProofChannelEventId string
					if proofChannelEventId != "" {
						actualProofChannelEventId = proofChannelEventId
					} else {
						actualProofChannelEventId = conf.Nostr.ProofChannelEventId
					}
					tagMap := map[string][]string{
						"e": {actualProofChannelEventId},
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

					fmt.Printf("Listening for proof events in channel id '%s'...\n", actualProofChannelEventId)
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
