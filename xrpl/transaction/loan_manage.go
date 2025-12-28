package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanManage transaction flags
const (
	// tfLoanDefault marks the loan as defaulted.
	tfLoanDefault uint32 = 0x00010000 // 65536
	// tfLoanImpair marks the loan as impaired.
	tfLoanImpair uint32 = 0x00020000 // 131072
	// tfLoanUnimpair removes the impaired status from the loan.
	tfLoanUnimpair uint32 = 0x00040000 // 262144
)

// LoanManage manages the state of a loan (XLS-66).
// Used to mark loans as defaulted, impaired, or to remove impairment.
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanManage",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "LoanID": "...",
//	    "Flags": 65536,
//	    "Fee": "10"
//	}
//
// ```
type LoanManage struct {
	BaseTx
	// The ID of the loan to manage.
	LoanID types.Hash256
}

// TxType returns the type of the transaction (LoanManage).
func (*LoanManage) TxType() TxType {
	return LoanManageTx
}

// Flatten returns a flattened map of the LoanManage transaction.
func (l *LoanManage) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanManageTx.String()
	flattened["LoanID"] = l.LoanID.String()

	return flattened
}

// SetTfLoanDefault sets the tfLoanDefault flag.
func (l *LoanManage) SetTfLoanDefault() {
	l.Flags |= tfLoanDefault
}

// SetTfLoanImpair sets the tfLoanImpair flag.
func (l *LoanManage) SetTfLoanImpair() {
	l.Flags |= tfLoanImpair
}

// SetTfLoanUnimpair sets the tfLoanUnimpair flag.
func (l *LoanManage) SetTfLoanUnimpair() {
	l.Flags |= tfLoanUnimpair
}

// Validate validates the LoanManage transaction.
func (l *LoanManage) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanID == "" {
		return false, ErrInvalidLoanID
	}

	return true, nil
}

