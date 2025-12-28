package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanPay makes a payment on a loan (XLS-66).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanPay",
//	    "Account": "rBorrower...",
//	    "LoanID": "...",
//	    "Amount": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	        "value": "1000"
//	    },
//	    "Fee": "10"
//	}
//
// ```
type LoanPay struct {
	BaseTx
	// The ID of the loan to pay.
	LoanID types.Hash256
	// The amount to pay. Can be XRP, IOU, or MPT.
	Amount types.CurrencyAmount
}

// TxType returns the type of the transaction (LoanPay).
func (*LoanPay) TxType() TxType {
	return LoanPayTx
}

// Flatten returns a flattened map of the LoanPay transaction.
func (l *LoanPay) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanPayTx.String()
	flattened["LoanID"] = l.LoanID.String()
	flattened["Amount"] = l.Amount.Flatten()

	return flattened
}

// Validate validates the LoanPay transaction.
func (l *LoanPay) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanID == "" {
		return false, ErrInvalidLoanID
	}

	if l.Amount == nil {
		return false, ErrInvalidAmount
	}

	return true, nil
}

