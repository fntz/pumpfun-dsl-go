package dsl

import (
	"fmt"

	"github.com/fntz/pumpfun-dsl-go/pump"
	"github.com/gagliardetto/solana-go"
)

type BondingCurve struct {
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
}

func GetBondingCurveSetup(mint solana.PublicKey) (*BondingCurve, error) {
	seeds := [][]byte{
		[]byte("bonding-curve"),
		mint.Bytes(),
	}
	bondingCurve, _, err := solana.FindProgramAddress(seeds, pump.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive bonding curve address: %w", err)
	}

	associatedBondingCurve, _, err := solana.FindAssociatedTokenAddress(
		bondingCurve,
		mint,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive associated bonding curve address: %w", err)
	}
	return &BondingCurve{
		BondingCurve:           bondingCurve,
		AssociatedBondingCurve: associatedBondingCurve,
	}, nil
}
