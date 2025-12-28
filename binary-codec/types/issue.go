package types

import (
	"encoding/hex"
	"errors"
	"strings"

	addresscodec "github.com/Peersyst/xrpl-go/address-codec"
	"github.com/Peersyst/xrpl-go/binary-codec/types/interfaces"
)

const (
	// MPTIssuanceIDBytesLength is the number of bytes for an MPT issuance ID.
	MPTIssuanceIDBytesLength = 24
	// CurrencyBytesLength is the number of bytes for a currency code.
	CurrencyBytesLength = 20
	// AccountIDBytesLength is the number of bytes for an account ID.
	AccountIDBytesLength = 20
	// MPTSequenceBytesLength is the number of bytes for an MPT sequence.
	MPTSequenceBytesLength = 4
)

var (
	// ErrInvalidIssueObject is returned when the JSON object is not a valid Issue.
	// ErrInvalidIssueObject is returned when the JSON object is not a valid Issue.
	ErrInvalidIssueObject = errors.New("invalid issue object")
	// ErrInvalidCurrency is returned when the currency field is missing or invalid in the Issue JSON.
	ErrInvalidCurrency = errors.New("invalid currency")
	// ErrInvalidIssuer is returned when the issuer field is missing or invalid in the Issue JSON.
	ErrInvalidIssuer = errors.New("invalid issuer")
	// ErrMissingIssueLengthOption is returned when no length option is provided to Issue.ToJSON.
	ErrMissingIssueLengthOption = errors.New("missing length option for Issue.ToJSON")
	// XRPBytes is the serialized byte representation for native XRP (zero-value currency issuer).
	XRPBytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// NoAccountBytes is the serialized byte representation for the "no account" placeholder used in MPT.
	// This is 0x0000000000000000000000000000000000000001
	NoAccountBytes = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
)

// Issue represents an XRPL Issue, which is essentially an AccountID.
// It is used to identify the issuer of a currency in the XRPL.
// The FromJson method converts a classic address string to an AccountID byte slice.
// The ToJson method converts an AccountID byte slice back to a classic address string.
// This type is crucial for handling currency issuers in XRPL transactions and ledger entries.
type Issue struct {
	length int
}

// FromJSON parses a classic address string and returns the corresponding AccountID byte slice.
// It uses the addresscodec package to decode the classic address.
// If the input is not a valid classic address, it returns an error.
func (i *Issue) FromJSON(json any) ([]byte, error) {
	if !i.isIssueObject(json) {
		return nil, ErrInvalidIssueObject
	}

	mapObj, ok := json.(map[string]any)
	if !ok {
		return nil, ErrInvalidIssueObject
	}

	currency, ok := mapObj["currency"]
	if !ok {
		mptIssuanceID, ok := mapObj["mpt_issuance_id"].(string)
		if !ok {
			return nil, ErrInvalidCurrency
		}

		mptIssuanceIDBytes, err := hex.DecodeString(mptIssuanceID)
		if err != nil {
			return nil, err
		}

		i.length = MPTIssuanceIDBytesLength

		return mptIssuanceIDBytes, nil
	}

	currencyCodec := &Currency{}

	currencyBytes, err := currencyCodec.FromJSON(currency)
	if err != nil {
		return nil, err
	}

	issuer, ok := mapObj["issuer"]
	if issuerString, okstring := issuer.(string); ok && okstring {
		_, issuerBytes, err := addresscodec.DecodeClassicAddressToAccountID(issuerString)
		if err != nil {
			return nil, err
		}

		return append(currencyBytes, issuerBytes...), nil
	}

	return currencyBytes, nil
}

// ToJSON converts serialized Issue bytes back to a JSON object.
// The Issue type is self-describing:
// - XRP: 20 bytes of zeros
// - IOU: 20 bytes currency + 20 bytes issuer account
// - MPT: 20 bytes account + 20 bytes NO_ACCOUNT + 4 bytes sequence
// The opts parameter is optional and can be used to specify the length hint.
func (i *Issue) ToJSON(p interfaces.BinaryParser, opts ...int) (any, error) {
	currencyCodec := &Currency{}

	// If a length hint is provided and it's MPT length, handle MPT directly
	if len(opts) > 0 && opts[0] == MPTIssuanceIDBytesLength {
		return i.parseMPTIssue(p)
	}

	// Read the first 20 bytes (currency or account for MPT)
	currencyOrAccount, err := p.ReadBytes(CurrencyBytesLength)
	if err != nil {
		return nil, err
	}

	// Check if it's XRP (all zeros)
	isXRP := true
	for _, b := range currencyOrAccount {
		if b != 0 {
			isXRP = false
			break
		}
	}

	if isXRP {
		return map[string]any{
			"currency": "XRP",
		}, nil
	}

	// Read the next 20 bytes (issuer account or NO_ACCOUNT for MPT)
	issuerOrNoAccount, err := p.ReadBytes(AccountIDBytesLength)
	if err != nil {
		return nil, err
	}

	// Check if it's NO_ACCOUNT (MPT indicator)
	isNoAccount := true
	for idx, b := range issuerOrNoAccount {
		if b != NoAccountBytes[idx] {
			isNoAccount = false
			break
		}
	}

	if isNoAccount {
		// This is an MPT - read the 4-byte sequence
		sequence, err := p.ReadBytes(MPTSequenceBytesLength)
		if err != nil {
			return nil, err
		}

		// Convert from little-endian to big-endian for the mpt_issuance_id
		// The mpt_issuance_id format is: sequence (4 bytes big-endian) + account (20 bytes)
		sequenceBE := []byte{sequence[3], sequence[2], sequence[1], sequence[0]}
		mptIssuanceID := append(sequenceBE, currencyOrAccount...)

		return map[string]any{
			"mpt_issuance_id": strings.ToUpper(hex.EncodeToString(mptIssuanceID)),
		}, nil
	}

	// This is an IOU - convert currency bytes to string
	currencyStr, err := currencyCodec.bytesToCurrencyString(currencyOrAccount)
	if err != nil {
		return nil, err
	}

	address, err := addresscodec.Encode(issuerOrNoAccount, []byte{addresscodec.AccountAddressPrefix}, addresscodec.AccountAddressLength)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"currency": currencyStr,
		"issuer":   address,
	}, nil
}

// parseMPTIssue parses an MPT issue from the binary parser.
func (i *Issue) parseMPTIssue(p interfaces.BinaryParser) (any, error) {
	b, err := p.ReadBytes(MPTIssuanceIDBytesLength)
	if err != nil {
		return nil, err
	}

	id := hex.EncodeToString(b)

	return map[string]any{
		"mpt_issuance_id": strings.ToUpper(id),
	}, nil
}

func (i *Issue) isIssueObject(obj any) bool {
	mapObj, ok := obj.(map[string]any)
	if !ok {
		return false
	}

	nKeys := len(mapObj)

	_, okMptIssuanceID := mapObj["mpt_issuance_id"]
	if nKeys == 1 && okMptIssuanceID {
		return true
	}

	_, okCurrency := mapObj["currency"]
	if nKeys == 1 && okCurrency {
		return true
	}

	_, okIssuer := mapObj["issuer"]
	if nKeys == 2 && okCurrency && okIssuer {
		return true
	}

	return false
}
