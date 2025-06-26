package dsl

import (
	"fmt"
	"math/big"

	"github.com/fntz/pumpfun-dsl-go/pump"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computeBudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
)

const maxSlippage uint = 100
const minSlippage uint = 0

type Builder struct {
	client      Client
	createSetup *NewTokenSetup
	buySetup    *BuySetup
	sellSetup   *SellSetup
	logger      DSLLogger
}

func NewBuilder(client Client, logger DSLLogger) *Builder {
	return &Builder{
		client: client,
		logger: logger,
	}
}

func (b *Builder) Create(setup *NewTokenSetup) *Builder {
	b.createSetup = setup
	return b
}

func (b *Builder) Buy(setup *BuySetup) *Builder {
	b.buySetup = setup
	return b
}

func (b *Builder) Sell(setup *SellSetup) *Builder {
	b.sellSetup = setup
	return b
}

func (b *Builder) CreateTx() (*solana.Transaction, error) {
	allInstructions, err := b.GetInstructions()
	if err != nil {
		return nil, err
	}

	lastHash, err := b.client.GetLatestBlockhash()
	if err != nil {
		return nil, err
	}

	var userPubKey solana.PublicKey
	if b.createSetup != nil {
		userPubKey = b.createSetup.UserPubKey
	} else if b.buySetup != nil {
		userPubKey = b.buySetup.UserPubKey
	} else if b.sellSetup != nil {
		userPubKey = b.sellSetup.UserPubKey
	} else {
		return nil, fmt.Errorf("no setup provided")
	}

	tx, err := solana.NewTransaction(
		allInstructions,
		*lastHash,
		solana.TransactionPayer(userPubKey),
	)

	return tx, err
}

func (b *Builder) GetInstructions() ([]solana.Instruction, error) {
	allInstructions := []solana.Instruction{}

	b.logger.Info("create transaction")

	if b.createSetup != nil {
		b.logger.Info("found create setup")
		instructions, err := b.fromCreateSetup(b.createSetup)
		if err != nil {
			return nil, err
		}
		allInstructions = append(allInstructions, instructions...)
	}

	if b.buySetup != nil {
		b.logger.Info("found buy setup")
		buySetup := b.buySetup
		if !buySetup.Filled() {
			if b.createSetup == nil {
				return nil, fmt.Errorf("buy setup is not filled and no create setup provided")
			}
			buySetup = b.createSetup.ToBuySetup(buySetup.AmountLamports, buySetup.Slippage)
		}
		instructions, err := b.fromBuySetup(buySetup, b.createSetup != nil && b.createSetup.Mint == buySetup.Mint)
		if err != nil {
			return nil, err
		}
		allInstructions = append(allInstructions, instructions...)
	}

	if b.sellSetup != nil {
		b.logger.Info("found sell setup")
		instructions, err := b.fromSellSetup(b.sellSetup, b.createSetup != nil && b.createSetup.Mint == b.sellSetup.Mint)
		if err != nil {
			return nil, err
		}
		allInstructions = append(allInstructions, instructions...)
	}

	if len(allInstructions) == 0 {
		b.logger.Error("no setup provided", nil)
		return nil, fmt.Errorf("no setup provided")
	}

	return allInstructions, nil
}

func (b *Builder) fromCreateSetup(mint *NewTokenSetup) ([]solana.Instruction, error) {
	instructions := []solana.Instruction{}
	bd, err := GetBondingCurveSetup(mint.Mint)
	if err != nil {
		return nil, err
	}
	b.logger.Info(fmt.Sprintf("found mint: %s", mint.Mint))
	metadata, _, err := solana.FindTokenMetadataAddress(mint.Mint)
	if err != nil {
		return nil, fmt.Errorf("can't find token metadata address: %w", err)
	}
	createInstruction := pump.NewCreateInstruction(
		mint.Name,
		mint.Symbol,
		mint.Uri,
		mint.UserPubKey,
		mint.Mint,
		pumpFunMintAuthority,
		bd.BondingCurve,
		bd.AssociatedBondingCurve,
		globalPumpFunAddress,
		solana.TokenMetadataProgramID,
		metadata,
		mint.UserPubKey,
		system.ProgramID,
		token.ProgramID,
		associatedtokenaccount.ProgramID,
		solana.SysVarRentPubkey,
		pumpFunEventAuthority,
		pump.ProgramID,
	)

	cuLimitInst := computeBudget.NewSetComputeUnitLimitInstruction(pumpFunCoumuteLimit)
	cupInst, err := b.client.GetCUPriceInstructions(mint.UserPubKey)
	if err != nil {
		return nil, err
	}

	instructions = append(instructions, cuLimitInst.Build())
	instructions = append(instructions, cupInst.Build())
	instructions = append(instructions, createInstruction.Build())

	return instructions, nil
}

func (b *Builder) fromBuySetup(buySetup *BuySetup, isFromNewTokenSetup bool) ([]solana.Instruction, error) {
	if buySetup.Slippage < minSlippage || buySetup.Slippage > maxSlippage {
		return nil, fmt.Errorf("slippage must be in range [%d, %d]", minSlippage, maxSlippage)
	}

	instructions := []solana.Instruction{}

	if !isFromNewTokenSetup {
		limitInstr := computeBudget.NewSetComputeUnitLimitInstruction(uint32(250000))
		priceInstr := computeBudget.NewSetComputeUnitPriceInstruction(100000)

		instructions = append(instructions, limitInstr.Build())
		instructions = append(instructions, priceInstr.Build())
	}

	bd, err := GetBondingCurveSetup(buySetup.Mint)
	if err != nil {
		return nil, err
	}

	assocAddress, _, err := solana.FindAssociatedTokenAddress(buySetup.UserPubKey, buySetup.Mint)
	if err != nil {
		return nil, fmt.Errorf("can't find associated token address: %w", err)
	}

	exist, err := b.client.AccountExist(assocAddress)
	if err != nil {
		return nil, err
	}

	if !exist {
		ataInstr, err := associatedtokenaccount.NewCreateInstruction(buySetup.UserPubKey, buySetup.UserPubKey, buySetup.Mint).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("can't create associated token account: %w", err)
		}
		instructions = append(instructions, ataInstr)
	}

	bondingCurve, err := b.client.FetchBoundingCurve(bd.BondingCurve)
	if err != nil {
		if !isFromNewTokenSetup {
			return nil, fmt.Errorf("can't fetch bonding curve: %w", err)
		} else {
			b.logger.Info("no bonding curve found (is new token setup), using default")
			bondingCurve = &FetchBoundingCurveResponse{
				RealTokenReserves:    big.NewInt(0),
				VirtualTokenReserves: big.NewInt(1_000_000_000),
				VirtualSolReserves:   big.NewInt(30),
			}
		}
	}

	buyResult := bondingCurve.CalculateBuyAmount(buySetup.AmountLamports, buySetup.ToPercent())

	b.logger.Info(fmt.Sprintf("buy result: %+v", buyResult))

	var creator solana.PublicKey
	if isFromNewTokenSetup {
		creator = buySetup.UserPubKey
	} else {
		creator, err = b.client.GetAccountInfo(bd.BondingCurve)
		if err != nil {
			return nil, fmt.Errorf("can't get account info: %w", err)
		}
	}

	creatorVault, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("creator-vault"),
			creator.Bytes(),
		},
		pump.ProgramID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't find creator vault address: %w", err)
	}

	buyInstr := pump.NewBuyInstruction(
		buyResult.TokensToBuy.Uint64(),
		buyResult.SolCost.Uint64(),
		globalPumpFunAddress,
		pumpFunFeeRecipient,
		buySetup.Mint,
		bd.BondingCurve,
		bd.AssociatedBondingCurve,
		assocAddress,
		buySetup.UserPubKey,
		system.ProgramID,
		token.ProgramID,
		creatorVault,
		pumpFunEventAuthority,
		pump.ProgramID,
	)

	buy, err := buyInstr.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("can't validate and build buy instruction: %w", err)
	}

	instructions = append(instructions, buy)

	return instructions, nil
}

func (b *Builder) fromSellSetup(sellSetup *SellSetup, isFromNewTokenSetup bool) ([]solana.Instruction, error) {
	if sellSetup.Slippage < minSlippage || sellSetup.Slippage > maxSlippage {
		return nil, fmt.Errorf("slippage must be in range [%d, %d]", minSlippage, maxSlippage)
	}

	instructions := []solana.Instruction{}

	culInst := computeBudget.NewSetComputeUnitLimitInstruction(uint32(250000))
	cupInst := computeBudget.NewSetComputeUnitPriceInstruction(uint64(10000))
	instructions = append(instructions, culInst.Build())
	instructions = append(instructions, cupInst.Build())

	assocAddress, _, err := solana.FindAssociatedTokenAddress(sellSetup.UserPubKey, sellSetup.Mint)
	if err != nil {
		return nil, fmt.Errorf("can't find associated token address: %w", err)
	}

	sellTokens := sellSetup.Tokens

	if sellSetup.All {
		currentBalance, err := b.client.GetTokenAccountBalance(sellSetup.Mint, assocAddress)
		if err != nil {
			return nil, fmt.Errorf("can't get current balance: %w", err)
		}
		if currentBalance == 0 {
			return nil, fmt.Errorf("no tokens to sell")
		}

		sellTokens = uint64(currentBalance)
	}

	bd, err := GetBondingCurveSetup(sellSetup.Mint)
	if err != nil {
		return nil, err
	}

	boundingCurve, err := b.client.FetchBoundingCurve(bd.BondingCurve)
	if err != nil {
		return nil, fmt.Errorf("can't fetch bonding curve: %w", err)
	}

	var creator solana.PublicKey
	if isFromNewTokenSetup {
		creator = sellSetup.UserPubKey
	} else {
		creator, err = b.client.GetAccountInfo(bd.BondingCurve)
		if err != nil {
			return nil, fmt.Errorf("can't get account info: %w", err)
		}
	}

	creatorVault, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("creator-vault"),
			creator.Bytes(),
		},
		pump.ProgramID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't find creator vault address: %w", err)
	}

	quote := boundingCurve.CalculateSellQuote(sellTokens, sellSetup.ToPercent())

	sellInstr := pump.NewSellInstruction(
		sellTokens,
		quote.Uint64(),
		globalPumpFunAddress,
		pumpFunFeeRecipient,
		sellSetup.Mint,
		bd.BondingCurve,
		bd.AssociatedBondingCurve,
		assocAddress,
		sellSetup.UserPubKey,
		system.ProgramID,
		creatorVault,
		token.ProgramID,
		pumpFunEventAuthority,
		pump.ProgramID,
	)

	sell, err := sellInstr.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("can't validate and build sell instruction: %w", err)
	}

	instructions = append(instructions, sell)

	return instructions, nil
}
