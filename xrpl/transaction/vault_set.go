package transaction

import (
	"strconv"

	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
)

// VaultSet modifies an existing vault's parameters (XLS-65).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultSet",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "AssetsMaximum": "2000000",
//	    "Fee": "10"
//	}
//
// ```
type VaultSet struct {
	BaseTx
	// The ID of the vault to modify.
	VaultID types.Hash256
	// Optional metadata for the vault.
	Data *string `json:",omitempty"`
	// Maximum amount of assets the vault can hold.
	AssetsMaximum *uint64 `json:",omitempty"`
	// Optional domain ID for permissioned domains.
	DomainID *types.Hash256 `json:",omitempty"`
}

// TxType returns the type of the transaction (VaultSet).
func (*VaultSet) TxType() TxType {
	return VaultSetTx
}

// Flatten returns a flattened map of the VaultSet transaction.
func (v *VaultSet) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultSetTx.String()
	flattened["VaultID"] = v.VaultID.String()

	if v.Data != nil {
		flattened["Data"] = *v.Data
	}

	if v.AssetsMaximum != nil {
		flattened["AssetsMaximum"] = strconv.FormatUint(*v.AssetsMaximum, 10)
	}

	if v.DomainID != nil {
		flattened["DomainID"] = v.DomainID.String()
	}

	return flattened
}

// Validate validates the VaultSet transaction.
func (v *VaultSet) Validate() (bool, error) {
	_, err := v.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if v.VaultID == "" {
		return false, ErrInvalidVaultID
	}

	return true, nil
}

