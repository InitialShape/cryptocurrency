# Cryptocurrency

[![Build Status](https://travis-ci.org/InitialShape/cryptocurrency.svg?branch=master)](https://travis-ci.org/InitialShape/cryptocurrency)

This is a pet project of the author InitialShape. It's about creating an
actually working cryptocurrency and blockchain from scratch.
It's written in Go, not the first language of the author.

It's mostly a learning project, but shall be launched still at one point in the
future.

## Project status

A few things are working already. There is:

- A working p2p network that can bootstrap, ping, register and send
transactions around
- An HTTP API that allows for sending transactions, minted blocks and
that allows miners to be implemented
- An almost working transaction evaluation pipeline that checks the validity of
transactions within blocks
- A miner
- A tool for creating and sending transactions

A few things are still needing to be taken care of:

- The mining algorithm and verification is currently very primitive. Better ones
would be nice
- The p2p network is currently in a pet-status state. It can ping and send
transactions around, but it's not even close from being reliable
- A fork-choice algorithm has to be implemented
- A better transaction creation tool for sending transfers needs to be
implemented
- Incoming data like transactions and blocks' schema should be evaluated using
some JSO-schema-ish library
- The Nakamoto consensus algorithm is still outstanding for completion.

If anyone would like to help with this project, feel free to create issues and
help out. Let's make this a real cryptocurrency!

## Installation

```bash
# go in your golang Path
git clone https://github.com/InitialShape/cryptocurrency
dep ensure
cd github.com/InitialShape/cryptocurrency
# Generate a wallet.txt file in ./
go run main.go --generate_keys
# go run main.go <dbname> <port TCP> <port http>
# for example
go run main.go db 1234 8000

# to mine (have the full node running)
# go run cmd/miner/miner.go <full node http> <nr of processes>
go run cmd/miner/miner.go http://localhost:8000 4

# to send transactions to the mempool
go run cmd/txtool/txtool.go
```
