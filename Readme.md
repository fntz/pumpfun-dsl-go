
# pumpfun-dsl-go

create init/buy/sell transaction in go, without pain.

motivation: I want to create a library to quickly check the possibilities of creatring tokens. Part of the code is taken from [this repo](https://github.com/prdsrm/pumpdotfun-go-sdk), but I didn't use that library directly because I wanted to:

- Decouple the methods from requiring an rpcClient to be passed from above

- Provide a simple interface focused only on making transactions"


# examples:

### before

```go
import (
    "github.com/fntz/pumpfun-dsl-go/dsl"
)

userWallet = ...

rpcClient := rpc.New(RPC_URL)

solanaClient := dsl.NewDefaultClient(rpcClient)

// choose logger
zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
logger := dsl.NewZeroLogger(&zl)
// or 
logger := dsl.NewFMTLogger()
// or implement own

builder := dsl.NewBuilder(solanaClient, logger)

```

### create & buy

```go

// make metadata
request := dsl.NewTokenRequest{
    FilePath:    IMAGE_PATH,
    Name:        "Test",
    Symbol:      "TEST",
    Description: "Test",
}
httpClient := &http.Client{}
response, err := dsl.Upload(httpClient, request)
if err != nil {
    log.Fatal(err)
}


// mint and buy in the same time
mint := solana.NewWallet()

tx, err := builder.
    Create(
        dsl.NewToken(request.Name, request.Symbol, response.MetadataUri, mint.PublicKey(), userWallet.PublicKey()),
    ).
    Buy(
        dsl.NewBuySetup(1000000, 90),
    ).
    CreateTx()
if err != nil {
    log.Fatal(err)
}

// fmt.Println("mint token: ", mint.PublicKey())

tx.Sign(
    func(key solana.PublicKey) *solana.PrivateKey {
        if userWallet.PublicKey().Equals(key) {
            return &userWallet.PrivateKey
        }
        if mint.PublicKey().Equals(key) {
            return &mint.PrivateKey
        }
        return nil
    },
)

sig, err := rpcClient.SendTransaction(ctx, tx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("tx sent: ", sig.String())

```

### sell

```go


toSell := solana.MustPublicKeyFromBase58("...")
	
tx, err := builder.Sell(
   dsl.NewSellSetup(1000000, 90, true, userWallet.PublicKey(), toSell),
).CreateTx()

if err != nil {
    log.Fatal(err)
}


tx.Sign(
    func(key solana.PublicKey) *solana.PrivateKey {
        if userWallet.PublicKey().Equals(key) {
            return &userWallet.PrivateKey
        }
        
        return nil
    },
)

sig, err := rpcClient.SendTransaction(ctx, tx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("tx sent: ", sig.String())

```

### get instructions

```go
	ins, err : = builder.
        Create(
			dsl.NewToken(request.Name, request.Symbol, response.MetadataUri, mint.PublicKey(), userWallet.PublicKey()),
		).
		Buy(
			dsl.NewBuySetup(1000000, 90),
		).
        GetInstructions()

	if err != nil {
		log.Fatal(err)
	}

    fmt.Println(ins)

```