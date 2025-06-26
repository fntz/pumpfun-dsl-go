package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fntz/pumpfun-dsl-go/dsl"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	PrivateKey string `env:"PK"`
	RPC_URL    string `env:"RPC_URL"`
	WS_URL     string `env:"WS_URL"`
	IMAGE_PATH string `env:"IMAGE_PATH"`
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	// upload metadata to pump.fun
	request := dsl.NewTokenRequest{
		FilePath:    c.IMAGE_PATH,
		Name:        "Test",
		Symbol:      "TEST1",
		Description: "Test",
	}
	httpClient := &http.Client{}
	response, err := dsl.Upload(httpClient, request)
	if err != nil {
		log.Fatal(err)
	}

	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := dsl.NewZeroLogger(&zl)
	str, _ := json.Marshal(response)
	fmt.Println("metadata uploaded: ", string(str))

	rpcClient := rpc.New(c.RPC_URL)

	// create token
	solanaClient := dsl.NewDefaultClient(rpcClient)

	builder := dsl.NewBuilder(solanaClient, logger)

	userWallet, err := solana.WalletFromPrivateKeyBase58(c.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------ create and buy
	// toBuy := solana.MustPublicKeyFromBase58("2WhgQqn7jh5XRCdkhTarnk3zruaGkd3o176GsjhMoDur")
	mint := solana.NewWallet()

	// // 100000
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
	// toSell := solana.MustPublicKeyFromBase58("8wRAX1uqexbDEkKy6xukquC1yR7iqgwvDJo9LYCqVAyt")
	// tx, err := builder.Sell(dsl.NewSellSetup(1000000, 90, true, userWallet.PublicKey(), toSell)).CreateTx()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	// ------------------------------------------------------------
}
