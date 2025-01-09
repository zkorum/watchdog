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

[ ] Run the server to listen to proofs to the Nostr protocol
[ ] Verify proofs received
[ ] Store proofs in a PGSQL database locally
[ ] Refine the proof semantics definition, grouping proofs by actions
[ ] Add support for Waku
[ ] Allow for receiving payloads as well
[ ] Allow broadcasting payloads via CLI commands
[ ] Generalize semantics outside of the Agora-specific context - moving Agora-specific stuff into configuration files or a specific usage of an external library
[ ] Move the underlying abstracted code to https://github.com/zkorum/racine and just reuse this by wrapping it with Agora-specific configuration

## License

© 2024 ZKorum SAS <hello@zkorum.com>.

This repo is licensed under the [MPLv2 License](./COPYING).
