package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// VaultClawback claws back share tokens from a holder (XLS-65).
// This can only be performed by the vault owner.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultClawback",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "Holder": "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
//	    "Fee": "10"
//	}
//
// ```
type VaultClawback struct {
	BaseTx
	// The ID of the vault.
	VaultID types.Hash256
	// The account holding the share tokens to claw back.
	Holder types.Address
	// Optional amount to claw back. If omitted, claws back all shares.
	Amount types.CurrencyAmount `json:",omitempty"`
}

// TxType returns the type of the transaction (VaultClawback).
func (*VaultClawback) TxType() TxType {
	return VaultClawbackTx
}

// Flatten returns a flattened map of the VaultClawback transaction.
func (v *VaultClawback) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultClawbackTx.String()
	flattened["VaultID"] = v.VaultID.String()
	flattened["Holder"] = v.Holder.String()

	if v.Amount != nil {
		flattened["Amount"] = v.Amount.Flatten()
	}

	return flattened
}

// Validate validates the VaultClawback transaction.
func (v *VaultClawback) Validate() (bool, error) {
	_, err := v.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if v.VaultID == "" {
		return false, ErrInvalidVaultID
	}

	if v.Holder == "" {
		return false, ErrInvalidHolder
	}

	return true, nil
}

