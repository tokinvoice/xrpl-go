package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanBrokerDelete deletes a LoanBroker object (XLS-66).
// The loan broker must have no outstanding loans to be deleted.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanBrokerDelete",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanBrokerID": "...",
//	    "Fee": "10"
//	}
//
// ```
type LoanBrokerDelete struct {
	BaseTx
	// The ID of the loan broker to delete.
	LoanBrokerID types.Hash256
}

// TxType returns the type of the transaction (LoanBrokerDelete).
func (*LoanBrokerDelete) TxType() TxType {
	return LoanBrokerDeleteTx
}

// Flatten returns a flattened map of the LoanBrokerDelete transaction.
func (l *LoanBrokerDelete) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanBrokerDeleteTx.String()
	flattened["LoanBrokerID"] = l.LoanBrokerID.String()

	return flattened
}

// Validate validates the LoanBrokerDelete transaction.
func (l *LoanBrokerDelete) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanBrokerID == "" {
		return false, ErrInvalidLoanBrokerID
	}

	return true, nil
}

