package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// VaultDelete deletes an existing vault (XLS-65).
// The vault must be empty (no assets and no outstanding shares) to be deleted.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultDelete",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "Fee": "10"
//	}
//
// ```
type VaultDelete struct {
	BaseTx
	// The ID of the vault to delete.
	VaultID types.Hash256
}

// TxType returns the type of the transaction (VaultDelete).
func (*VaultDelete) TxType() TxType {
	return VaultDeleteTx
}

// Flatten returns a flattened map of the VaultDelete transaction.
func (v *VaultDelete) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultDeleteTx.String()
	flattened["VaultID"] = v.VaultID.String()

	return flattened
}

// Validate validates the VaultDelete transaction.
func (v *VaultDelete) Validate() (bool, error) {
	_, err := v.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if v.VaultID == "" {
		return false, ErrInvalidVaultID
	}

	return true, nil
}

