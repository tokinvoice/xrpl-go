package transaction

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanSet transaction flags
const (
	// tfLoanOverpayment allows the loan to be overpaid.
	tfLoanOverpayment uint32 = 0x00010000 // 65536
)

// CounterpartySignature contains the counterparty's signature for a loan.
type CounterpartySignature struct {
	SigningPubKey *string        `json:",omitempty"`
	TxnSignature  *string        `json:",omitempty"`
	Signers       []types.Signer `json:",omitempty"`
}

// Flatten returns a flattened map of the CounterpartySignature.
func (c *CounterpartySignature) Flatten() map[string]interface{} {
	flat := make(map[string]interface{})
	if c.SigningPubKey != nil {
		flat["SigningPubKey"] = *c.SigningPubKey
	}
	if c.TxnSignature != nil {
		flat["TxnSignature"] = *c.TxnSignature
	}
	if len(c.Signers) > 0 {
		signers := make([]interface{}, len(c.Signers))
		for i, s := range c.Signers {
			signers[i] = s.Flatten()
		}
		flat["Signers"] = signers
	}
	return flat
}

// LoanSet creates a new loan from a loan broker (XLS-66).
// Requires signatures from both the borrower and the loan broker owner (counterparty).
//
// Example:
//
// ```json
//
//	{
//	    "TransactionType": "LoanSet",
//	    "Account": "rBorrower...",
//	    "LoanBrokerID": "...",
//	    "PrincipalRequested": "100000",
//	    "Counterparty": "rLender...",
//	    "InterestRate": 500,
//	    "PaymentTotal": 12,
//	    "PaymentInterval": 2592000,
//	    "Fee": "10"
//	}
//
// ```
type LoanSet struct {
	BaseTx
	// The ID of the loan broker to request a loan from.
	LoanBrokerID types.Hash256
	// The amount of principal requested for the loan.
	PrincipalRequested string
	// The counterparty's signature (loan broker owner).
	CounterpartySignature *CounterpartySignature `json:",omitempty"`
	// The counterparty account (loan broker owner).
	Counterparty *types.Address `json:",omitempty"`
	// Optional metadata for the loan.
	Data *string `json:",omitempty"`
	// One-time fee charged when the loan is originated.
	LoanOriginationFee *string `json:",omitempty"`
	// Ongoing fee charged for loan servicing.
	LoanServiceFee *string `json:",omitempty"`
	// Fee charged for late payments.
	LatePaymentFee *string `json:",omitempty"`
	// Fee charged when closing the loan early.
	ClosePaymentFee *string `json:",omitempty"`
	// Fee rate for overpayment (basis points).
	OverpaymentFee *uint32 `json:",omitempty"`
	// Interest rate (basis points).
	InterestRate *uint32 `json:",omitempty"`
	// Interest rate for late payments (basis points).
	LateInterestRate *uint32 `json:",omitempty"`
	// Interest rate for early close (basis points).
	CloseInterestRate *uint32 `json:",omitempty"`
	// Interest rate for overpayment (basis points).
	OverpaymentInterestRate *uint32 `json:",omitempty"`
	// Total number of payments for the loan.
	PaymentTotal *uint32 `json:",omitempty"`
	// Interval between payments in seconds.
	PaymentInterval *uint32 `json:",omitempty"`
	// Grace period in seconds after payment due date.
	GracePeriod *uint32 `json:",omitempty"`
}

// TxType returns the type of the transaction (LoanSet).
func (*LoanSet) TxType() TxType {
	return LoanSetTx
}

// Flatten returns a flattened map of the LoanSet transaction.
func (l *LoanSet) Flatten() FlatTransaction {
	flattened := l.BaseTx.Flatten()

	flattened["TransactionType"] = LoanSetTx.String()
	flattened["LoanBrokerID"] = l.LoanBrokerID.String()
	flattened["PrincipalRequested"] = l.PrincipalRequested

	if l.CounterpartySignature != nil {
		flattened["CounterpartySignature"] = l.CounterpartySignature.Flatten()
	}
	if l.Counterparty != nil {
		flattened["Counterparty"] = l.Counterparty.String()
	}
	if l.Data != nil {
		flattened["Data"] = *l.Data
	}
	if l.LoanOriginationFee != nil {
		flattened["LoanOriginationFee"] = *l.LoanOriginationFee
	}
	if l.LoanServiceFee != nil {
		flattened["LoanServiceFee"] = *l.LoanServiceFee
	}
	if l.LatePaymentFee != nil {
		flattened["LatePaymentFee"] = *l.LatePaymentFee
	}
	if l.ClosePaymentFee != nil {
		flattened["ClosePaymentFee"] = *l.ClosePaymentFee
	}
	if l.OverpaymentFee != nil {
		flattened["OverpaymentFee"] = *l.OverpaymentFee
	}
	if l.InterestRate != nil {
		flattened["InterestRate"] = *l.InterestRate
	}
	if l.LateInterestRate != nil {
		flattened["LateInterestRate"] = *l.LateInterestRate
	}
	if l.CloseInterestRate != nil {
		flattened["CloseInterestRate"] = *l.CloseInterestRate
	}
	if l.OverpaymentInterestRate != nil {
		flattened["OverpaymentInterestRate"] = *l.OverpaymentInterestRate
	}
	if l.PaymentTotal != nil {
		flattened["PaymentTotal"] = *l.PaymentTotal
	}
	if l.PaymentInterval != nil {
		flattened["PaymentInterval"] = *l.PaymentInterval
	}
	if l.GracePeriod != nil {
		flattened["GracePeriod"] = *l.GracePeriod
	}

	return flattened
}

// SetTfLoanOverpayment sets the tfLoanOverpayment flag.
func (l *LoanSet) SetTfLoanOverpayment() {
	l.Flags |= tfLoanOverpayment
}

// Validate validates the LoanSet transaction.
func (l *LoanSet) Validate() (bool, error) {
	_, err := l.BaseTx.Validate()
	if err != nil {
		return false, err
	}

	if l.LoanBrokerID == "" {
		return false, ErrInvalidLoanBrokerID
	}

	if l.PrincipalRequested == "" {
		return false, ErrInvalidPrincipalRequested
	}

	return true, nil
}

