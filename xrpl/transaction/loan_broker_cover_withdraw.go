package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanBrokerCoverWithdraw withdraws cover (collateral) from a loan broker (XLS-66).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanBrokerCoverWithdraw",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanBrokerID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "5000"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type LoanBrokerCoverWithdraw struct {
	BaseTx
	// The ID of the loan broker to withdraw cover from.
	LoanBrokerID types.Hash256
	// The amount of cover to withdraw. Can be XRP, IOU, or MPT.
	Amount types.CurrencyAmount
	// Optional destination account for the withdrawn cover.
	Destination *types.Address `json:",omitempty"`
}

// TxType returns the type of the transaction (LoanBrokerCoverWithdraw).
func (*LoanBrokerCoverWithdraw) TxType() TxType {
	return LoanBrokerCoverWithdrawTx
}

// Flatten returns a flattened map of the LoanBrokerCoverWithdraw transaction.
func (l *LoanBrokerCoverWithdraw) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanBrokerCoverWithdrawTx.String()
	flattened["LoanBrokerID"] = l.LoanBrokerID.String()
	flattened["Amount"] = l.Amount.Flatten()

	if l.Destination != nil {
		flattened["Destination"] = l.Destination.String()
	}

	return flattened
}

// Validate validates the LoanBrokerCoverWithdraw transaction.
func (l *LoanBrokerCoverWithdraw) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanBrokerID == "" {
		return false, ErrInvalidLoanBrokerID
	}

	if l.Amount == nil {
		return false, ErrInvalidAmount
	}

	return true, nil
}

