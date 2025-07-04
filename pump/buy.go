// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package pump

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Buys tokens from a bonding curve.
type BuyInstruction struct {
	Amount     *uint64
	MaxSolCost *uint64

	// [0] = [] global
	//
	// [1] = [WRITE] fee_recipient
	//
	// [2] = [] mint
	//
	// [3] = [WRITE] bonding_curve
	//
	// [4] = [WRITE] associated_bonding_curve
	//
	// [5] = [WRITE] associated_user
	//
	// [6] = [WRITE, SIGNER] user
	//
	// [7] = [] system_program
	//
	// [8] = [] token_program
	//
	// [9] = [WRITE] creator_vault
	//
	// [10] = [] event_authority
	//
	// [11] = [] program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewBuyInstructionBuilder creates a new `BuyInstruction` instruction builder.
func NewBuyInstructionBuilder() *BuyInstruction {
	nd := &BuyInstruction{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 12),
	}
	nd.AccountMetaSlice[7] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	nd.AccountMetaSlice[8] = ag_solanago.Meta(Addresses["TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"])
	return nd
}

// SetAmount sets the "amount" parameter.
func (inst *BuyInstruction) SetAmount(amount uint64) *BuyInstruction {
	inst.Amount = &amount
	return inst
}

// SetMaxSolCost sets the "max_sol_cost" parameter.
func (inst *BuyInstruction) SetMaxSolCost(max_sol_cost uint64) *BuyInstruction {
	inst.MaxSolCost = &max_sol_cost
	return inst
}

// SetGlobalAccount sets the "global" account.
func (inst *BuyInstruction) SetGlobalAccount(global ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(global)
	return inst
}

func (inst *BuyInstruction) findFindGlobalAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: global
	seeds = append(seeds, []byte{byte(0x67), byte(0x6c), byte(0x6f), byte(0x62), byte(0x61), byte(0x6c)})

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindGlobalAddressWithBumpSeed calculates Global account address with given seeds and a known bump seed.
func (inst *BuyInstruction) FindGlobalAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindGlobalAddress(bumpSeed)
	return
}

func (inst *BuyInstruction) MustFindGlobalAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindGlobalAddress finds Global account address with given seeds.
func (inst *BuyInstruction) FindGlobalAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindGlobalAddress(0)
	return
}

func (inst *BuyInstruction) MustFindGlobalAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindGlobalAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetGlobalAccount gets the "global" account.
func (inst *BuyInstruction) GetGlobalAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetFeeRecipientAccount sets the "fee_recipient" account.
func (inst *BuyInstruction) SetFeeRecipientAccount(feeRecipient ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(feeRecipient).WRITE()
	return inst
}

// GetFeeRecipientAccount gets the "fee_recipient" account.
func (inst *BuyInstruction) GetFeeRecipientAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetMintAccount sets the "mint" account.
func (inst *BuyInstruction) SetMintAccount(mint ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
func (inst *BuyInstruction) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetBondingCurveAccount sets the "bonding_curve" account.
func (inst *BuyInstruction) SetBondingCurveAccount(bondingCurve ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(bondingCurve).WRITE()
	return inst
}

func (inst *BuyInstruction) findFindBondingCurveAddress(mint ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: bonding-curve
	seeds = append(seeds, []byte{byte(0x62), byte(0x6f), byte(0x6e), byte(0x64), byte(0x69), byte(0x6e), byte(0x67), byte(0x2d), byte(0x63), byte(0x75), byte(0x72), byte(0x76), byte(0x65)})
	// path: mint
	seeds = append(seeds, mint.Bytes())

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindBondingCurveAddressWithBumpSeed calculates BondingCurve account address with given seeds and a known bump seed.
func (inst *BuyInstruction) FindBondingCurveAddressWithBumpSeed(mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindBondingCurveAddress(mint, bumpSeed)
	return
}

func (inst *BuyInstruction) MustFindBondingCurveAddressWithBumpSeed(mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindBondingCurveAddress(mint, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindBondingCurveAddress finds BondingCurve account address with given seeds.
func (inst *BuyInstruction) FindBondingCurveAddress(mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindBondingCurveAddress(mint, 0)
	return
}

func (inst *BuyInstruction) MustFindBondingCurveAddress(mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindBondingCurveAddress(mint, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetBondingCurveAccount gets the "bonding_curve" account.
func (inst *BuyInstruction) GetBondingCurveAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetAssociatedBondingCurveAccount sets the "associated_bonding_curve" account.
func (inst *BuyInstruction) SetAssociatedBondingCurveAccount(associatedBondingCurve ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(associatedBondingCurve).WRITE()
	return inst
}

func (inst *BuyInstruction) findFindAssociatedBondingCurveAddress(bondingCurve ag_solanago.PublicKey, mint ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// path: bondingCurve
	seeds = append(seeds, bondingCurve.Bytes())
	// const (raw): [6 221 246 225 215 101 161 147 217 203 225 70 206 235 121 172 28 180 133 237 95 91 55 145 58 140 245 133 126 255 0 169]
	seeds = append(seeds, []byte{byte(0x6), byte(0xdd), byte(0xf6), byte(0xe1), byte(0xd7), byte(0x65), byte(0xa1), byte(0x93), byte(0xd9), byte(0xcb), byte(0xe1), byte(0x46), byte(0xce), byte(0xeb), byte(0x79), byte(0xac), byte(0x1c), byte(0xb4), byte(0x85), byte(0xed), byte(0x5f), byte(0x5b), byte(0x37), byte(0x91), byte(0x3a), byte(0x8c), byte(0xf5), byte(0x85), byte(0x7e), byte(0xff), byte(0x0), byte(0xa9)})
	// path: mint
	seeds = append(seeds, mint.Bytes())

	programID := Addresses["ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"]

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, programID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, programID)
	}
	return
}

// FindAssociatedBondingCurveAddressWithBumpSeed calculates AssociatedBondingCurve account address with given seeds and a known bump seed.
func (inst *BuyInstruction) FindAssociatedBondingCurveAddressWithBumpSeed(bondingCurve ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindAssociatedBondingCurveAddress(bondingCurve, mint, bumpSeed)
	return
}

func (inst *BuyInstruction) MustFindAssociatedBondingCurveAddressWithBumpSeed(bondingCurve ag_solanago.PublicKey, mint ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindAssociatedBondingCurveAddress(bondingCurve, mint, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindAssociatedBondingCurveAddress finds AssociatedBondingCurve account address with given seeds.
func (inst *BuyInstruction) FindAssociatedBondingCurveAddress(bondingCurve ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindAssociatedBondingCurveAddress(bondingCurve, mint, 0)
	return
}

func (inst *BuyInstruction) MustFindAssociatedBondingCurveAddress(bondingCurve ag_solanago.PublicKey, mint ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindAssociatedBondingCurveAddress(bondingCurve, mint, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetAssociatedBondingCurveAccount gets the "associated_bonding_curve" account.
func (inst *BuyInstruction) GetAssociatedBondingCurveAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetAssociatedUserAccount sets the "associated_user" account.
func (inst *BuyInstruction) SetAssociatedUserAccount(associatedUser ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(associatedUser).WRITE()
	return inst
}

// GetAssociatedUserAccount gets the "associated_user" account.
func (inst *BuyInstruction) GetAssociatedUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetUserAccount sets the "user" account.
func (inst *BuyInstruction) SetUserAccount(user ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(user).WRITE().SIGNER()
	return inst
}

// GetUserAccount gets the "user" account.
func (inst *BuyInstruction) GetUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSystemProgramAccount sets the "system_program" account.
func (inst *BuyInstruction) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "system_program" account.
func (inst *BuyInstruction) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetTokenProgramAccount sets the "token_program" account.
func (inst *BuyInstruction) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "token_program" account.
func (inst *BuyInstruction) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetCreatorVaultAccount sets the "creator_vault" account.
func (inst *BuyInstruction) SetCreatorVaultAccount(creatorVault ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(creatorVault).WRITE()
	return inst
}

func (inst *BuyInstruction) findFindCreatorVaultAddress(creatorVault ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: creator-vault
	seeds = append(seeds, []byte{byte(0x63), byte(0x72), byte(0x65), byte(0x61), byte(0x74), byte(0x6f), byte(0x72), byte(0x2d), byte(0x76), byte(0x61), byte(0x75), byte(0x6c), byte(0x74)})
	// path: creatorVault
	seeds = append(seeds, creatorVault.Bytes())

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindCreatorVaultAddressWithBumpSeed calculates CreatorVault account address with given seeds and a known bump seed.
func (inst *BuyInstruction) FindCreatorVaultAddressWithBumpSeed(creatorVault ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindCreatorVaultAddress(creatorVault, bumpSeed)
	return
}

func (inst *BuyInstruction) MustFindCreatorVaultAddressWithBumpSeed(creatorVault ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindCreatorVaultAddress(creatorVault, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindCreatorVaultAddress finds CreatorVault account address with given seeds.
func (inst *BuyInstruction) FindCreatorVaultAddress(creatorVault ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindCreatorVaultAddress(creatorVault, 0)
	return
}

func (inst *BuyInstruction) MustFindCreatorVaultAddress(creatorVault ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindCreatorVaultAddress(creatorVault, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetCreatorVaultAccount gets the "creator_vault" account.
func (inst *BuyInstruction) GetCreatorVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetEventAuthorityAccount sets the "event_authority" account.
func (inst *BuyInstruction) SetEventAuthorityAccount(eventAuthority ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(eventAuthority)
	return inst
}

func (inst *BuyInstruction) findFindEventAuthorityAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: __event_authority
	seeds = append(seeds, []byte{byte(0x5f), byte(0x5f), byte(0x65), byte(0x76), byte(0x65), byte(0x6e), byte(0x74), byte(0x5f), byte(0x61), byte(0x75), byte(0x74), byte(0x68), byte(0x6f), byte(0x72), byte(0x69), byte(0x74), byte(0x79)})

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindEventAuthorityAddressWithBumpSeed calculates EventAuthority account address with given seeds and a known bump seed.
func (inst *BuyInstruction) FindEventAuthorityAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindEventAuthorityAddress(bumpSeed)
	return
}

func (inst *BuyInstruction) MustFindEventAuthorityAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindEventAuthorityAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindEventAuthorityAddress finds EventAuthority account address with given seeds.
func (inst *BuyInstruction) FindEventAuthorityAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindEventAuthorityAddress(0)
	return
}

func (inst *BuyInstruction) MustFindEventAuthorityAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindEventAuthorityAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetEventAuthorityAccount gets the "event_authority" account.
func (inst *BuyInstruction) GetEventAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetProgramAccount sets the "program" account.
func (inst *BuyInstruction) SetProgramAccount(program ag_solanago.PublicKey) *BuyInstruction {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(program)
	return inst
}

// GetProgramAccount gets the "program" account.
func (inst *BuyInstruction) GetProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

func (inst BuyInstruction) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_Buy,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst BuyInstruction) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *BuyInstruction) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if inst.MaxSolCost == nil {
			return errors.New("MaxSolCost parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Global is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.FeeRecipient is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.BondingCurve is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.AssociatedBondingCurve is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.AssociatedUser is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.User is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.CreatorVault is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.EventAuthority is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.Program is not set")
		}
	}
	return nil
}

func (inst *BuyInstruction) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Buy")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=2]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("      Amount", *inst.Amount))
						paramsBranch.Child(ag_format.Param("  MaxSolCost", *inst.MaxSolCost))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=12]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                  global", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("           fee_recipient", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("                    mint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("           bonding_curve", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("associated_bonding_curve", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("         associated_user", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("                    user", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("          system_program", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("           token_program", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("           creator_vault", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("         event_authority", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("                 program", inst.AccountMetaSlice.Get(11)))
					})
				})
		})
}

func (obj BuyInstruction) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	// Serialize `MaxSolCost` param:
	err = encoder.Encode(obj.MaxSolCost)
	if err != nil {
		return err
	}
	return nil
}
func (obj *BuyInstruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Amount`:
	err = decoder.Decode(&obj.Amount)
	if err != nil {
		return err
	}
	// Deserialize `MaxSolCost`:
	err = decoder.Decode(&obj.MaxSolCost)
	if err != nil {
		return err
	}
	return nil
}

// NewBuyInstruction declares a new Buy instruction with the provided parameters and accounts.
func NewBuyInstruction(
	// Parameters:
	amount uint64,
	max_sol_cost uint64,
	// Accounts:
	global ag_solanago.PublicKey,
	feeRecipient ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	bondingCurve ag_solanago.PublicKey,
	associatedBondingCurve ag_solanago.PublicKey,
	associatedUser ag_solanago.PublicKey,
	user ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	creatorVault ag_solanago.PublicKey,
	eventAuthority ag_solanago.PublicKey,
	program ag_solanago.PublicKey) *BuyInstruction {
	return NewBuyInstructionBuilder().
		SetAmount(amount).
		SetMaxSolCost(max_sol_cost).
		SetGlobalAccount(global).
		SetFeeRecipientAccount(feeRecipient).
		SetMintAccount(mint).
		SetBondingCurveAccount(bondingCurve).
		SetAssociatedBondingCurveAccount(associatedBondingCurve).
		SetAssociatedUserAccount(associatedUser).
		SetUserAccount(user).
		SetSystemProgramAccount(systemProgram).
		SetTokenProgramAccount(tokenProgram).
		SetCreatorVaultAccount(creatorVault).
		SetEventAuthorityAccount(eventAuthority).
		SetProgramAccount(program)
}
