#!/bin/bash

# rm -rf solana-anchor-go
# rm -rf pump-fun.json

# wget https://raw.githubusercontent.com/pump-fun/pump-public-docs/refs/heads/main/idl/pump.json

#git clone https://github.com/fragmetric-labs/solana-anchor-go

cd solana-anchor-go

go mod tidy

go build .

cd ..

./solana-anchor-go/solana-anchor-go -dst pump -src pump.json

# rm -rf *.json
# rm -rf solana-anchor-go
