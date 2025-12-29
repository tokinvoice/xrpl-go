package transaction

import (
	"strconv"

	ledger "github.com/Peersyst/xrpl-go/xrpl/ledger-entry-types"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
)

// VaultCreate transaction flags
const (
	// tfVaultPrivate indicates the vault is private.
	tfVaultPrivate uint32 = 0x00010000 // 65536
	// tfVaultShareNonTransferable indicates vault shares cannot be transferred.
	tfVaultShareNonTransferable uint32 = 0x00020000 // 131072
)

// Withdrawal policy constants
const (
	// VaultStrategyFirstComeFirstServe is the first-come-first-serve withdrawal policy.
	VaultStrategyFirstComeFirstServe uint8 = 1
)

// VaultCreate creates a new vault for holding a single asset type (XLS-65).
// The vault issues share tokens (MPT) to depositors.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultCreate",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "Asset": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type VaultCreate struct {
	BaseTx
	// The asset this vault will hold. In JSON, this is an object with currency and issuer fields.
	Asset ledger.Asset
	// Optional metadata for the vault.
	Data *string `json:",omitempty"`
	// Maximum amount of assets the vault can hold.
	AssetsMaximum *uint64 `json:",omitempty"`
	// Optional metadata for the MPT share tokens.
	MPTokenMetadata *string `json:",omitempty"`
	// The withdrawal policy for this vault.
	WithdrawalPolicy *uint8 `json:",omitempty"`
	// Optional domain ID for permissioned domains.
	DomainID *types.Hash256 `json:",omitempty"`
}

// TxType returns the type of the transaction (VaultCreate).
func (*VaultCreate) TxType() TxType {
	return VaultCreateTx
}

// Flatten returns a flattened map of the VaultCreate transaction.
func (v *VaultCreate) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultCreateTx.String()
	flattened["Asset"] = v.Asset.Flatten()

	if v.Data != nil {
		flattened["Data"] = *v.Data
	}

	if v.AssetsMaximum != nil {
		flattened["AssetsMaximum"] = strconv.FormatUint(*v.AssetsMaximum, 10)
	}

	if v.MPTokenMetadata != nil {
		flattened["MPTokenMetadata"] = *v.MPTokenMetadata
	}

	if v.WithdrawalPolicy != nil {
		flattened["WithdrawalPolicy"] = *v.WithdrawalPolicy
	}

	if v.DomainID != nil {
		flattened["DomainID"] = v.DomainID.String()
	}

	return flattened
}

// SetTfVaultPrivate sets the tfVaultPrivate flag.
func (v *VaultCreate) SetTfVaultPrivate() {
	v.Flags |= tfVaultPrivate
}

// SetTfVaultShareNonTransferable sets the tfVaultShareNonTransferable flag.
func (v *VaultCreate) SetTfVaultShareNonTransferable() {
	v.Flags |= tfVaultShareNonTransferable
}

// Validate validates the VaultCreate transaction.
func (v *VaultCreate) Validate() (bool, error) {
	_, err := v.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if v.Asset.Currency == "" {
		return false, ErrInvalidAsset
	}

	return true, nil
}

