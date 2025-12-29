package transaction

import (
	"strconv"

	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
)

// LoanBrokerSet creates or modifies a LoanBroker object (XLS-66).
// A LoanBroker manages loans from a vault.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanBrokerSet",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "VaultID": "...",
//	    "DebtMaximum": "1000000",
//	    "CoverRateMinimum": 15000,
//	    "CoverRateLiquidation": 12000,
//	    "Fee": "10"
//	}
//
// ```
type LoanBrokerSet struct {
	BaseTx
	// The ID of the vault to create a loan broker for (required for creation).
	VaultID types.Hash256
	// The ID of an existing loan broker to modify (required for modification).
	LoanBrokerID *types.Hash256 `json:",omitempty"`
	// Optional metadata for the loan broker.
	Data *string `json:",omitempty"`
	// Fee rate charged for managing loans (basis points).
	ManagementFeeRate *uint32 `json:",omitempty"`
	// Maximum total debt allowed for this loan broker.
	DebtMaximum *uint64 `json:",omitempty"`
	// Minimum cover rate required (basis points).
	CoverRateMinimum *uint32 `json:",omitempty"`
	// Cover rate at which liquidation can occur (basis points).
	CoverRateLiquidation *uint32 `json:",omitempty"`
}

// TxType returns the type of the transaction (LoanBrokerSet).
func (*LoanBrokerSet) TxType() TxType {
	return LoanBrokerSetTx
}

// Flatten returns a flattened map of the LoanBrokerSet transaction.
func (l *LoanBrokerSet) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanBrokerSetTx.String()
	flattened["VaultID"] = l.VaultID.String()

	if l.LoanBrokerID != nil {
		flattened["LoanBrokerID"] = l.LoanBrokerID.String()
	}

	if l.Data != nil {
		flattened["Data"] = *l.Data
	}

	if l.ManagementFeeRate != nil {
		flattened["ManagementFeeRate"] = *l.ManagementFeeRate
	}

	if l.DebtMaximum != nil {
		flattened["DebtMaximum"] = strconv.FormatUint(*l.DebtMaximum, 10)
	}

	if l.CoverRateMinimum != nil {
		flattened["CoverRateMinimum"] = *l.CoverRateMinimum
	}

	if l.CoverRateLiquidation != nil {
		flattened["CoverRateLiquidation"] = *l.CoverRateLiquidation
	}

	return flattened
}

// Validate validates the LoanBrokerSet transaction.
func (l *LoanBrokerSet) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.VaultID == "" {
		return false, ErrInvalidVaultID
	}

	return true, nil
}

