package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanBrokerCoverClawback claws back cover from a loan broker (XLS-66).
// This can only be performed by the asset issuer.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanBrokerCoverClawback",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanBrokerID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "1000"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type LoanBrokerCoverClawback struct {
	BaseTx
	// The ID of the loan broker to claw back cover from.
	LoanBrokerID *types.Hash256 `json:",omitempty"`
	// The amount to claw back. If omitted, claws back all.
	Amount types.CurrencyAmount `json:",omitempty"`
}

// TxType returns the type of the transaction (LoanBrokerCoverClawback).
func (*LoanBrokerCoverClawback) TxType() TxType {
	return LoanBrokerCoverClawbackTx
}

// Flatten returns a flattened map of the LoanBrokerCoverClawback transaction.
func (l *LoanBrokerCoverClawback) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanBrokerCoverClawbackTx.String()

	if l.LoanBrokerID != nil {
		flattened["LoanBrokerID"] = l.LoanBrokerID.String()
	}

	if l.Amount != nil {
		flattened["Amount"] = l.Amount.Flatten()
	}

	return flattened
}

// Validate validates the LoanBrokerCoverClawback transaction.
func (l *LoanBrokerCoverClawback) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	return true, nil
}

