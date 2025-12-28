package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// VaultWithdraw withdraws assets from a vault by redeeming share tokens (XLS-65).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultWithdraw",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "500"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type VaultWithdraw struct {
	BaseTx
	// The ID of the vault to withdraw from.
	VaultID types.Hash256
	// The amount to withdraw. Can be XRP, IOU, or MPT.
	Amount types.CurrencyAmount
	// Optional destination account for the withdrawn assets.
	Destination *types.Address `json:",omitempty"`
}

// TxType returns the type of the transaction (VaultWithdraw).
func (*VaultWithdraw) TxType() TxType {
	return VaultWithdrawTx
}

// Flatten returns a flattened map of the VaultWithdraw transaction.
func (v *VaultWithdraw) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultWithdrawTx.String()
	flattened["VaultID"] = v.VaultID.String()
	flattened["Amount"] = v.Amount.Flatten()

	if v.Destination != nil {
		flattened["Destination"] = v.Destination.String()
	}

	return flattened
}

// Validate validates the VaultWithdraw transaction.
func (v *VaultWithdraw) Validate() (bool, error) {
	_, err := v.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if v.VaultID == "" {
		return false, ErrInvalidVaultID
	}

	if v.Amount == nil {
		return false, ErrInvalidAmount
	}

	return true, nil
}

