package binarycodec

import (
	"testing"
)

func TestVaultCreateEncode(t *testing.T) {
	tx := map[string]any{
		"TransactionType": "VaultCreate",
		"Account":         "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		"Asset": map[string]any{
			"currency": "XRP",
		},
		"AssetsMaximum":    "1000000000000",
		"WithdrawalPolicy": uint8(1),
		"Fee":              "10",
		"Sequence":         uint32(1),
	}

	encoded, err := Encode(tx)
	if err != nil {
		t.Fatalf("Error encoding VaultCreate: %v", err)
	}
	t.Logf("Encoded: %s", encoded)

	// Now decode it back
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding VaultCreate: %v", err)
	}
	t.Logf("Decoded: %+v", decoded)
}

