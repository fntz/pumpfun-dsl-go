#!/bin/bash

rm -rf solana-anchor-go
rm -rf pump-fun.json

wget https://raw.githubusercontent.com/rckprtr/pumpdotfun-sdk/refs/heads/main/src/IDL/pump-fun.json 

git clone https://github.com/fragmetric-labs/solana-anchor-go

cd solana-anchor-go

go mod tidy

go build .

cd ..

./solana-anchor-go/solana-anchor-go -dst pump -src pump-fun.json

rm -rf pump-fun.json
rm -rf solana-anchor-go