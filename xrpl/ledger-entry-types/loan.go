package ledger

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// Loan ledger entry flags
const (
	// lsfLoanDefault indicates the loan is in default.
	lsfLoanDefault uint32 = 0x00010000 // 65536
	// lsfLoanImpaired indicates the loan is impaired.
	lsfLoanImpaired uint32 = 0x00020000 // 131072
	// lsfLoanOverpayment indicates overpayment is allowed.
	lsfLoanOverpayment uint32 = 0x00040000 // 262144
)

// Loan represents a loan ledger entry (XLS-66).
// A Loan tracks an individual loan between a borrower and a loan broker.
//
// Example:
//
//	{
//	    "LedgerEntryType": "Loan",
//	    "Flags": 0,
//	    "LoanSequence": 1,
//	    "OwnerNode": "0",
//	    "LoanBrokerNode": "0",
//	    "LoanBrokerID": "...",
//	    "Borrower": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "StartDate": 1234567890,
//	    "PaymentInterval": 86400,
//	    "GracePeriod": 604800,
//	    "NextPaymentDueDate": 1234654290,
//	    "PaymentRemaining": 12,
//	    "PrincipalOutstanding": "100000000",
//	    "PeriodicPayment": "10000000",
//	    "TotalValueOutstanding": "120000000",
//	    "PreviousTxnID": "...",
//	    "PreviousTxnLgrSeq": 123456
//	}
type Loan struct {
	// The unique ID for this ledger entry.
	// In JSON, this field is represented with different names depending on the context and API method.
	Index types.Hash256 `json:"index,omitempty"`
	// The type of ledger entry. Always "Loan" for this type.
	LedgerEntryType EntryType
	// Set of bit-flags for this ledger entry.
	Flags uint32
	// The sequence number for this loan within the broker.
	LoanSequence uint32
	// A hint indicating which page of the owner directory links to this entry.
	OwnerNode uint64
	// A hint indicating which page of the loan broker directory links to this entry.
	LoanBrokerNode uint64
	// The ID of the loan broker managing this loan.
	LoanBrokerID types.Hash256
	// The account that borrowed funds.
	Borrower types.Address
	// One-time fee charged when the loan is originated.
	LoanOriginationFee string `json:",omitempty"`
	// Periodic fee charged for servicing the loan.
	LoanServiceFee string `json:",omitempty"`
	// Fee charged for late payments.
	LatePaymentFee string `json:",omitempty"`
	// Fee charged for early loan closure.
	ClosePaymentFee string `json:",omitempty"`
	// Fee charged for overpayment.
	OverpaymentFee uint32 `json:",omitempty"`
	// The interest rate for regular payments (as a percentage * 10000).
	InterestRate uint32 `json:",omitempty"`
	// The interest rate applied to late payments (as a percentage * 10000).
	LateInterestRate uint32 `json:",omitempty"`
	// The interest rate applied when closing early (as a percentage * 10000).
	CloseInterestRate uint32 `json:",omitempty"`
	// The interest rate applied to overpayments (as a percentage * 10000).
	OverpaymentInterestRate uint32 `json:",omitempty"`
	// The start date of the loan (Ripple time).
	StartDate uint32
	// The interval between payments in seconds.
	PaymentInterval uint32
	// The grace period for late payments in seconds.
	GracePeriod uint32
	// The date of the previous payment (Ripple time).
	PreviousPaymentDate uint32 `json:",omitempty"`
	// The date the next payment is due (Ripple time).
	NextPaymentDueDate uint32
	// The number of payments remaining.
	PaymentRemaining uint32
	// The outstanding principal amount.
	PrincipalOutstanding string
	// The periodic payment amount.
	PeriodicPayment string
	// The total value outstanding including interest and fees.
	TotalValueOutstanding string
	// The identifying hash of the transaction that most recently modified this entry.
	PreviousTxnID types.Hash256 `json:",omitempty"`
	// The index of the ledger that contains the transaction that most recently modified this entry.
	PreviousTxnLgrSeq uint32 `json:",omitempty"`
}

// EntryType returns the type of the ledger entry.
func (*Loan) EntryType() EntryType {
	return LoanEntry
}

// SetLsfLoanDefault sets the lsfLoanDefault flag.
func (l *Loan) SetLsfLoanDefault() {
	l.Flags |= lsfLoanDefault
}

// HasLsfLoanDefault returns true if the lsfLoanDefault flag is set.
func (l *Loan) HasLsfLoanDefault() bool {
	return l.Flags&lsfLoanDefault != 0
}

// SetLsfLoanImpaired sets the lsfLoanImpaired flag.
func (l *Loan) SetLsfLoanImpaired() {
	l.Flags |= lsfLoanImpaired
}

// HasLsfLoanImpaired returns true if the lsfLoanImpaired flag is set.
func (l *Loan) HasLsfLoanImpaired() bool {
	return l.Flags&lsfLoanImpaired != 0
}

// SetLsfLoanOverpayment sets the lsfLoanOverpayment flag.
func (l *Loan) SetLsfLoanOverpayment() {
	l.Flags |= lsfLoanOverpayment
}

// HasLsfLoanOverpayment returns true if the lsfLoanOverpayment flag is set.
func (l *Loan) HasLsfLoanOverpayment() bool {
	return l.Flags&lsfLoanOverpayment != 0
}

