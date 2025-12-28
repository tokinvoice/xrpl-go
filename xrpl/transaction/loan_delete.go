package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanDelete deletes a loan (XLS-66).
// The loan must be fully paid off or defaulted to be deleted.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanDelete",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanID": "...",
//	    "Fee": "10"
//	}
//
// ```
type LoanDelete struct {
	BaseTx
	// The ID of the loan to delete.
	LoanID types.Hash256
}

// TxType returns the type of the transaction (LoanDelete).
func (*LoanDelete) TxType() TxType {
	return LoanDeleteTx
}

// Flatten returns a flattened map of the LoanDelete transaction.
func (l *LoanDelete) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanDeleteTx.String()
	flattened["LoanID"] = l.LoanID.String()

	return flattened
}

// Validate validates the LoanDelete transaction.
func (l *LoanDelete) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanID == "" {
		return false, ErrInvalidLoanID
	}

	return true, nil
}

