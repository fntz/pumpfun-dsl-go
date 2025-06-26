package dsl

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"

	"github.com/fntz/pumpfun-dsl-go/pump"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computeBudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	rpc "github.com/gagliardetto/solana-go/rpc"
)

type Client interface {
	GetCUPriceInstructions(userPubKey solana.PublicKey) (*computeBudget.SetComputeUnitPrice, error)
	GetLatestBlockhash() (*solana.Hash, error)
	AccountExist(pubkey solana.PublicKey) (bool, error)
	FetchBoundingCurve(boundingCurve solana.PublicKey) (*FetchBoundingCurveResponse, error)
	GetTokenAccountBalance(token, user solana.PublicKey) (int, error)
	GetAccountInfo(bd solana.PublicKey) (solana.PublicKey, error)
}

type DefaultClient struct {
	rpcClient *rpc.Client
}

var _ Client = (*DefaultClient)(nil)

func NewDefaultClient(rpcClient *rpc.Client) *DefaultClient {
	return &DefaultClient{
		rpcClient: rpcClient,
	}
}

type BondingCurveData struct {
	Padding              [8]byte // 8 байт заполнения
	VirtualTokenReserves uint64
	VirtualSolReserves   uint64
	RealTokenReserves    uint64
	RealSolReserves      uint64
	TokenTotalSupply     uint64
	Complete             bool
	Creator              solana.PublicKey // 32-байтный публичный ключ
}

func ParseBondingCurveAccount(data []byte) (*BondingCurveData, error) {
	if len(data) != 8+8*5+1+32 { // 8 (padding) + 8*5 (uint64) + 1 (bool) + 32 (pubkey)
		return nil, fmt.Errorf("invalid data length: expected %d, got %d", 8+8*5+1+32, len(data))
	}

	var parsed BondingCurveData

	// Чтение полей
	copy(parsed.Padding[:], data[:8])
	offset := 8

	parsed.VirtualTokenReserves = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	parsed.VirtualSolReserves = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	parsed.RealTokenReserves = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	parsed.RealSolReserves = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	parsed.TokenTotalSupply = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8

	parsed.Complete = data[offset] != 0
	offset += 1

	copy(parsed.Creator[:], data[offset:offset+32])

	return &parsed, nil
}

// https://github.com/AL-THE-BOT-FATHER/pump_fun_py/blob/main/pump_fun_py/coin_data.py#L20
func (c *DefaultClient) GetAccountInfo(bd solana.PublicKey) (solana.PublicKey, error) {
	accountInfo, err := c.rpcClient.GetAccountInfoWithOpts(context.TODO(), bd, &rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})
	if err != nil || accountInfo.Value == nil {
		return solana.PublicKey{}, fmt.Errorf("FBCD: failed to get account info: %w", err)
	}

	data := accountInfo.Value.Data.GetBinary()
	result, err := ParseBondingCurveAccount(data)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("FBCD: failed to parse account info: %w", err)
	}

	return result.Creator, nil
}

func (c *DefaultClient) GetCUPriceInstructions(userPubKey solana.PublicKey) (*computeBudget.SetComputeUnitPrice, error) {
	out, err := c.rpcClient.GetRecentPrioritizationFees(context.TODO(),
		solana.PublicKeySlice{userPubKey,
			pump.ProgramID,
			pumpFunMintAuthority,
			globalPumpFunAddress,
			solana.TokenMetadataProgramID,
			system.ProgramID, token.ProgramID, associatedtokenaccount.ProgramID, solana.SysVarRentPubkey, pumpFunEventAuthority})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent prioritization fees: %w", err)
	}
	var median uint64
	length := uint64(len(out))
	for _, fee := range out {
		median = fee.PrioritizationFee
	}
	median /= length
	cupInst := computeBudget.NewSetComputeUnitPriceInstruction(median)
	return cupInst, nil
}

func (c *DefaultClient) AccountExist(pubkey solana.PublicKey) (bool, error) {
	_, err := c.rpcClient.GetAccountInfo(context.TODO(), pubkey)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (c *DefaultClient) GetLatestBlockhash() (*solana.Hash, error) {
	recent, err := c.rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("error while getting recent block hash: %w", err)
	}

	return &recent.Value.Blockhash, nil
}

func (c *DefaultClient) GetTokenAccountBalance(token, user solana.PublicKey) (int, error) {
	tokenAccounts, err := c.rpcClient.GetTokenAccountBalance(
		context.TODO(),
		user,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return 0, fmt.Errorf("can't get amount of token in balance: %w", err)
	}
	amount, err := strconv.Atoi(tokenAccounts.Value.Amount)
	if err != nil {
		return 0, fmt.Errorf("can't convert token amount to integer: %w", err)
	}

	return amount, nil
}

func (c *DefaultClient) FetchBoundingCurve(boundingCurve solana.PublicKey) (*FetchBoundingCurveResponse, error) {
	accountInfo, err := c.rpcClient.GetAccountInfoWithOpts(context.TODO(), boundingCurve, &rpc.GetAccountInfoOpts{Encoding: solana.EncodingBase64, Commitment: rpc.CommitmentProcessed})
	if err != nil || accountInfo.Value == nil {
		return nil, fmt.Errorf("FBCD: failed to get account info: %w", err)
	}

	data := accountInfo.Value.Data.GetBinary()
	if len(data) < 24 {
		return nil, fmt.Errorf("FBCD: insufficient data length")
	}

	// Decode the bonding curve data assuming it follows little-endian format
	realTokenReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[0:8]))
	virtualTokenReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[8:16]))
	virtualSolReserves := big.NewInt(0).SetUint64(binary.LittleEndian.Uint64(data[16:24]))

	fmt.Println("realTokenReserves: ", realTokenReserves)
	fmt.Println("virtualTokenReserves: ", virtualTokenReserves)
	fmt.Println("virtualSolReserves: ", virtualSolReserves)

	return &FetchBoundingCurveResponse{
		RealTokenReserves:    realTokenReserves,
		VirtualTokenReserves: virtualTokenReserves,
		VirtualSolReserves:   virtualSolReserves,
	}, nil
}

type FetchBoundingCurveResponse struct {
	RealTokenReserves    *big.Int
	VirtualTokenReserves *big.Int
	VirtualSolReserves   *big.Int
}

type BuyResult struct {
	TokensToBuy   *big.Int
	AvgPriceInSol *big.Float
	SolCost       *big.Int
}

func (f *FetchBoundingCurveResponse) CalculateBuyAmount(lamportsIn uint64, percent float64) *BuyResult {
	virtSolDecimal := new(big.Float).SetInt64(f.VirtualSolReserves.Int64())
	virtTokensDecimal := new(big.Float).SetInt64(f.VirtualTokenReserves.Int64())
	solInput := new(big.Float).SetInt64(int64(lamportsIn))
	solInput = solInput.Quo(solInput, big.NewFloat(float64(solana.LAMPORTS_PER_SOL)))
	slippageTolDecimal := new(big.Float).SetFloat64(percent)

	// ───────── commission ───────────────────────────────────────
	solNet := new(big.Float).Quo(solInput, big.NewFloat(1+PUMP_FEE_PCT))

	// ───────── invariant before ───────────────────────────────────────
	k := new(big.Float).Mul(virtSolDecimal, virtTokensDecimal)

	// ───────── how many tokens will be bought ──────────────────────────
	xBefore := virtSolDecimal
	xAfter := new(big.Float).Add(virtSolDecimal, solNet)

	tokensBefore := new(big.Float).Quo(k, xBefore)
	tokensAfter := new(big.Float).Quo(k, xAfter)
	tokensBought := new(big.Float).Sub(tokensBefore, tokensAfter)

	// ───────── calculate average price and min-receive ─────────────────────────
	avgPrice := new(big.Float).Quo(solNet, tokensBought)
	minTokensOut := new(big.Float).Mul(tokensBought, new(big.Float).Sub(big.NewFloat(1), slippageTolDecimal))

	avgPriceInt, _ := avgPrice.Mul(avgPrice, big.NewFloat(float64(solana.LAMPORTS_PER_SOL))).Int(nil)
	minTokensOutInt, _ := minTokensOut.Int(nil)

	return &BuyResult{
		TokensToBuy:   minTokensOutInt,
		AvgPriceInSol: avgPrice,
		SolCost:       avgPriceInt,
	}
}

/*

func buyTokens(virtSol, virtTokens float64, lamportsIn int64, slippageTol float64) (map[string]float64, error) {
	// ───────── подготовка чисел ─────────────────────────────────────────────
	virtSolDecimal := new(big.Float).SetFloat64(virtSol)
	virtTokensDecimal := new(big.Float).SetFloat64(virtTokens)
	solInput := new(big.Float).SetInt64(lamportsIn)
	solInput.Quo(solInput, big.NewFloat(LAMPORTS_PER_SOL))
	slippageTolDecimal := new(big.Float).SetFloat64(slippageTol)

	// ───────── Комиссия 1 % (сверху) ───────────────────────────────────────
	solNet := new(big.Float).Quo(solInput, big.NewFloat(1+PUMP_FEE_PCT)) // попадает в пул

	// ───────── k-инвариант до сделки ───────────────────────────────────────
	k := new(big.Float).Mul(virtSolDecimal, virtTokensDecimal)

	// ───────── Сколько токенов получит покупатель ──────────────────────────
	xBefore := virtSolDecimal
	xAfter := new(big.Float).Add(virtSolDecimal, solNet)

	tokensBefore := new(big.Float).Quo(k, xBefore)
	tokensAfter := new(big.Float).Quo(k, xAfter)
	tokensBought := new(big.Float).Sub(tokensBefore, tokensAfter)

	// ───────── Обновление пулов (flat-fee) ─────────────────────────────────
	virtSolNew := new(big.Float).Sub(xAfter, big.NewFloat(PUMP_FLAT_FEE)) // вычитаем 0.015 SOL
	virtTokensNew := new(big.Float).Sub(virtTokensDecimal, tokensBought)

	// ───────── Вычислим среднюю цену и min-receive ─────────────────────────
	avgPrice := new(big.Float).Quo(solNet, tokensBought)
	minTokensOut := new(big.Float).Mul(tokensBought, new(big.Float).Sub(big.NewFloat(1), slippageTolDecimal))

	tb, _ := tokensBought.Float64()
	mt, _ := minTokensOut.Float64()
	avg, _ := avgPrice.Float64()
	vsn, _ := virtSolNew.Float64()
	vtn, _ := virtTokensNew.Float64()

	result := map[string]float64{
		"tokens_bought":   tb,
		"min_tokens_out":  mt,
		"avg_price_SOL":   avg,
		"virt_sol_new":    vsn,
		"virt_tokens_new": vtn,
	}

	return result, nil
}
*/

func (f *FetchBoundingCurveResponse) CalculateSellQuote(tokens uint64, percent float64) *big.Int {
	amount := big.NewInt(int64(tokens))

	// Clone bonding curve data to avoid mutations
	virtualSolReserves := new(big.Int).Set(f.VirtualSolReserves)
	virtualTokenReserves := new(big.Int).Set(f.VirtualTokenReserves)

	// Compute the new virtual reserves
	x := new(big.Int).Mul(virtualSolReserves, amount)
	y := new(big.Int).Add(virtualTokenReserves, amount)
	a := new(big.Int).Div(x, y)
	percentageMultiplier := big.NewFloat(percent)
	sol := new(big.Float).SetInt(a)
	number := new(big.Float).Mul(sol, percentageMultiplier)
	final, _ := number.Int(nil)

	return final
}
