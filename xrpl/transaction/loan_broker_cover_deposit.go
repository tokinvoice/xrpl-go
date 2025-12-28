package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanBrokerCoverDeposit deposits cover (collateral) into a loan broker (XLS-66).
// Cover is used to back loans and absorb losses.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanBrokerCoverDeposit",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanBrokerID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "10000"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type LoanBrokerCoverDeposit struct {
	BaseTx
	// The ID of the loan broker to deposit cover into.
	LoanBrokerID types.Hash256
	// The amount of cover to deposit. Can be XRP, IOU, or MPT.
	Amount types.CurrencyAmount
}

// TxType returns the type of the transaction (LoanBrokerCoverDeposit).
func (*LoanBrokerCoverDeposit) TxType() TxType {
	return LoanBrokerCoverDepositTx
}

// Flatten returns a flattened map of the LoanBrokerCoverDeposit transaction.
func (l *LoanBrokerCoverDeposit) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanBrokerCoverDepositTx.String()
	flattened["LoanBrokerID"] = l.LoanBrokerID.String()
	flattened["Amount"] = l.Amount.Flatten()

	return flattened
}

// Validate validates the LoanBrokerCoverDeposit transaction.
func (l *LoanBrokerCoverDeposit) Validate() (bool, error) {
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

