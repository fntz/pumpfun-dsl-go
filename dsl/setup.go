package dsl

import (
	"github.com/gagliardetto/solana-go"
)

type NewTokenSetup struct {
	Name       string
	Symbol     string
	Uri        string
	UserPubKey solana.PublicKey
	Mint       solana.PublicKey
}

func NewToken(name string, symbol string, uri string, mint solana.PublicKey, userPubKey solana.PublicKey) *NewTokenSetup {
	return &NewTokenSetup{
		Name:       name,
		Symbol:     symbol,
		Uri:        uri,
		UserPubKey: userPubKey,
		Mint:       mint,
	}
}

func (n *NewTokenSetup) ToBuySetup(amountLamports uint64, slippage uint) *BuySetup {
	return NewBuyToken(amountLamports, slippage, n.UserPubKey, n.Mint)
}

type BuySetup struct {
	AmountLamports uint64
	Slippage       uint
	UserPubKey     solana.PublicKey
	Mint           solana.PublicKey
}

func (b *BuySetup) Filled() bool {
	return b.UserPubKey != (solana.PublicKey{}) && b.Mint != (solana.PublicKey{})
}

// will be filled by the NewTokenSetup
func NewBuySetup(amountLamports uint64, slippage uint) *BuySetup {
	return &BuySetup{
		AmountLamports: amountLamports,
		Slippage:       slippage,
	}
}

func NewBuyToken(amountLamports uint64, slippage uint, userPubKey solana.PublicKey, mint solana.PublicKey) *BuySetup {
	return &BuySetup{
		AmountLamports: amountLamports,
		Slippage:       slippage,
		UserPubKey:     userPubKey,
		Mint:           mint,
	}
}

func (b *BuySetup) ToPercent() float64 {
	return toPercent(b.Slippage)
}

type SellSetup struct {
	Tokens     uint64
	Slippage   uint
	All        bool
	UserPubKey solana.PublicKey
	Mint       solana.PublicKey
}

func NewSellSetup(tokens uint64, slippage uint, all bool, userPubKey solana.PublicKey, mint solana.PublicKey) *SellSetup {
	return &SellSetup{
		Tokens:     tokens,
		Slippage:   slippage,
		All:        all,
		UserPubKey: userPubKey,
		Mint:       mint,
	}
}

func (s *SellSetup) ToPercent() float64 {
	return toPercent(s.Slippage)
}

func toPercent(slippage uint) float64 {
	return float64(slippage) / 100
}
