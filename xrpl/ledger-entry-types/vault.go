package ledger

import "github.com/Peersyst/xrpl-go/xrpl/transaction/types"

// Vault ledger entry flags
const (
	// lsfVaultPrivate indicates the vault is private.
	lsfVaultPrivate uint32 = 0x00010000 // 65536
)

// Vault represents a single asset vault ledger entry (XLS-65).
// A Vault holds a single asset type and issues share tokens to depositors.
//
// Example:
//
//	{
//	    "LedgerEntryType": "Vault",
//	    "Flags": 0,
//	    "Sequence": 1,
//	    "OwnerNode": "0",
//	    "Owner": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
//	    "Asset": {
//	        "currency": "USD",
//	        "issuer": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
//	    },
//	    "AssetsTotal": "1000000",
//	    "AssetsAvailable": "1000000",
//	    "LossUnrealized": "0",
//	    "MPTokenIssuanceID": "00000001A407AF5856CCF3C42619DAA925813FC955C72983",
//	    "WithdrawalPolicy": 0,
//	    "PreviousTxnID": "...",
//	    "PreviousTxnLgrSeq": 123456
//	}
type Vault struct {
	// The unique ID for this ledger entry.
	// In JSON, this field is represented with different names depending on the context and API method.
	Index types.Hash256 `json:"index,omitempty"`
	// The type of ledger entry. Always "Vault" for this type.
	LedgerEntryType EntryType
	// Set of bit-flags for this ledger entry.
	Flags uint32
	// The sequence number used to identify this vault.
	Sequence uint32
	// A hint indicating which page of the owner directory links to this entry.
	OwnerNode uint64
	// The account that owns this vault.
	Owner types.Address
	// The account associated with this vault (same as Owner for vault operations).
	Account types.Address
	// The asset held by this vault. In JSON, this is an object with currency and issuer fields.
	Asset Asset
	// The total amount of assets in the vault, including those lent out.
	AssetsTotal string `json:",omitempty"`
	// The amount of assets currently available for withdrawal.
	AssetsAvailable string `json:",omitempty"`
	// Unrealized losses from lending activities.
	LossUnrealized string `json:",omitempty"`
	// The MPT Issuance ID for the share tokens issued by this vault.
	MPTokenIssuanceID types.Hash192
	// The withdrawal policy strategy code for this vault; see XLS-65 docs for valid values.
	WithdrawalPolicy uint8
	// The maximum amount of assets the vault can hold.
	AssetsMaximum string `json:",omitempty"`
	// Optional metadata for this vault.
	Data string `json:",omitempty"`
	// The identifying hash of the transaction that most recently modified this entry.
	PreviousTxnID types.Hash256 `json:",omitempty"`
	// The index of the ledger that contains the transaction that most recently modified this entry.
	PreviousTxnLgrSeq uint32 `json:",omitempty"`
}

// EntryType returns the type of the ledger entry.
func (*Vault) EntryType() EntryType {
	return VaultEntry
}

// SetLsfVaultPrivate sets the lsfVaultPrivate flag.
func (v *Vault) SetLsfVaultPrivate() {
	v.Flags |= lsfVaultPrivate
}

// HasLsfVaultPrivate returns true if the lsfVaultPrivate flag is set.
func (v *Vault) HasLsfVaultPrivate() bool {
	return v.Flags&lsfVaultPrivate != 0
}

