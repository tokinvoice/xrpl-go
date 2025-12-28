package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// VaultDeposit deposits assets into a vault and receives share tokens (XLS-65).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "VaultDeposit",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "1000"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type VaultDeposit struct {
	BaseTx
	// The ID of the vault to deposit into.
	VaultID types.Hash256
	// The amount to deposit. Can be XRP, IOU, or MPT.
	Amount types.CurrencyAmount
}

// TxType returns the type of the transaction (VaultDeposit).
func (*VaultDeposit) TxType() TxType {
	return VaultDepositTx
}

// Flatten returns a flattened map of the VaultDeposit transaction.
func (v *VaultDeposit) Flatten() FlatTransaction {
	flattened := v.BaseTx.Flatten()

	flattened["TransactionType"] = VaultDepositTx.String()
	flattened["VaultID"] = v.VaultID.String()
	flattened["Amount"] = v.Amount.Flatten()

	return flattened
}

// Validate validates the VaultDeposit transaction.
func (v *VaultDeposit) Validate() (bool, error) {
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

