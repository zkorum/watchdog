# Agora Watchdog Server

> [!WARNING]
> This repository is under heavy development (pre-alpha)

The Agora Watchdog Server is an Agora-specific implementation of one of the component of the [Racine Protocol](https://github.com/zkorum/racine).

Another aspect of the Racine Protocol is data generation on the first place, and especially proofs of origin with integration with the Rarimo Protocol. This is still integrated within the Agora product for now.

This repo contains:
- a CLI to:
    - post directly to Agora either via the Waku Protocol or via Nostr
    - record proofs (and potentially payloads) that are broadcast to Nostr or Waku
    - API to automatically audit Agora

## Roadmap

- [x] Run the server to listen to proofs from the Nostr protocol
- [ ] Verify proofs received
- [ ] Store proofs in a PGSQL database locally
- [ ] Refine the proof semantics definition, grouping proofs by actions
- [ ] Group proofs logically in PostgreSQL DB based on above semantics
- [ ] Add support to use Agora API to audit its moderation
- [ ] Add support for Waku
- [ ] Allow for receiving payloads as well
- [ ] Allow broadcasting payloads via CLI commands
- [ ] Allow to use the node to request specific data from _any_ storage (starting from Agora centralized DB, Nostr and Waku)
- [ ] Generalize semantics outside of the Agora-specific context - moving Agora-specific stuff into configuration files or a specific usage of an external library
- [ ] Move the underlying abstracted code to https://github.com/zkorum/racine and just reuse this by wrapping it with Agora-specific configuration

## Usage

→ go run main.go --help
NAME:
   Agora Watchdog CLI - verify Agora moderation history

USAGE:
   agora-watchdog -c ./watchdog.yml listen

VERSION:
   v0.1.0

COMMANDS:
   listen, l  listen to proofs, verify and store them in a local PostgreSQL database
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Location of the config file - default to ./watchdog.yml [$WATCHDOG_CONFIG]
   --help, -h                show help
   --version, -v             print the version

COPYRIGHT:
   (c) 2025 ZKorum SAS

## Acknowledgements

This project has received funding from the European Union's Horizon Europe 2020 research and innovation program through the [NGI TRUSTCHAIN](https://trustchain.ngi.eu/) program under cascade funding agreement No. 101093274.
The grant funded the project's ongoing ideation process, along with the development of the initial CLI structure and the implementation of the listen command.

## License

© 2024 ZKorum SAS <hello@zkorum.com>.

This repo is licensed under the [MPLv2 License](./COPYING).
