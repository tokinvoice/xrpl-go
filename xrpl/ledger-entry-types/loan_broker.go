package ledger

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// LoanBroker represents a loan broker ledger entry (XLS-66).
// A LoanBroker facilitates lending by connecting vault assets with borrowers.
//
// Example:
//
//	{
//	    "LedgerEntryType": "LoanBroker",
//	    "Flags": 0,
//	    "Sequence": 1,
//	    "LoanSequence": 0,
//	    "OwnerNode": "0",
//	    "VaultNode": "0",
//	    "VaultID": "...",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "Owner": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "DebtMaximum": "1000000000",
//	    "CoverRateMinimum": 150,
//	    "CoverRateLiquidation": 120,
//	    "PreviousTxnID": "...",
//	    "PreviousTxnLgrSeq": 123456
//	}
type LoanBroker struct {
	// The unique ID for this ledger entry.
	// In JSON, this field is represented with different names depending on the context and API method.
	Index types.Hash256 `json:"index,omitempty"`
	// The type of ledger entry. Always "LoanBroker" for this type.
	LedgerEntryType EntryType
	// Set of bit-flags for this ledger entry.
	Flags uint32
	// The sequence number used to identify this loan broker.
	Sequence uint32
	// The sequence number for the next loan created by this broker.
	LoanSequence uint32
	// A hint indicating which page of the owner directory links to this entry.
	OwnerNode uint64
	// A hint indicating which page of the vault directory links to this entry.
	VaultNode uint64
	// The ID of the vault this broker is connected to.
	VaultID types.Hash256
	// The account associated with this loan broker.
	Account types.Address
	// The account that owns this loan broker.
	Owner types.Address
	// The number of objects owned by this loan broker.
	OwnerCount uint32 `json:",omitempty"`
	// The total amount of debt currently outstanding from loans issued by this broker.
	DebtTotal string `json:",omitempty"`
	// The maximum amount of debt this broker can have outstanding.
	DebtMaximum string
	// The amount of cover collateral currently available.
	CoverAvailable string `json:",omitempty"`
	// The minimum collateral coverage rate required for loans (as a percentage * 1000).
	CoverRateMinimum uint32 `json:",omitempty"`
	// The collateral rate at which loans become eligible for liquidation (as a percentage * 1000).
	CoverRateLiquidation uint32 `json:",omitempty"`
	// The identifying hash of the transaction that most recently modified this entry.
	PreviousTxnID types.Hash256 `json:",omitempty"`
	// The index of the ledger that contains the transaction that most recently modified this entry.
	PreviousTxnLgrSeq uint32 `json:",omitempty"`
}

// EntryType returns the type of the ledger entry.
func (*LoanBroker) EntryType() EntryType {
	return LoanBrokerEntry
}

