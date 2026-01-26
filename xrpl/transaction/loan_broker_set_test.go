package transaction

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanBrokerSet_TxType(t *testing.T) {
	tx := &LoanBrokerSet{}
	assert.Equal(t, tx.TxType(), LoanBrokerSetTx)
}

func TestLoanBrokerSet_Flatten(t *testing.T) {
	testcases := []struct {
		name     string
		tx       *LoanBrokerSet
		expected map[string]interface{}
	}{
		{
			name: "pass - empty",
			tx:   &LoanBrokerSet{},
			expected: map[string]interface{}{
				"TransactionType": LoanBrokerSetTx.String(),
				"VaultID":         "",
			},
		},
		{
			name: "pass - complete",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:            "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					Fee:                1000000,
					Sequence:           1,
					LastLedgerSequence: 3000000,
				},
				VaultID:           "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				ManagementFeeRate: func() *uint32 { v := uint32(1000); return &v }(),
				CoverRateMinimum:  func() *uint32 { v := uint32(5000); return &v }(),
			},
			expected: map[string]interface{}{
				"TransactionType":    LoanBrokerSetTx.String(),
				"Account":            "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
				"Fee":                "1000000",
				"Sequence":           uint32(1),
				"LastLedgerSequence": uint32(3000000),
				"VaultID":            "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				"ManagementFeeRate":  uint32(1000),
				"CoverRateMinimum":   uint32(5000),
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert.Equal(t, testcase.tx.Flatten(), testcase.expected)
		})
	}
}

func TestLoanBrokerSet_Validate(t *testing.T) {
	testcases := []struct {
		name     string
		tx       *LoanBrokerSet
		expected error
	}{
		{
			name: "fail - base tx invalid",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					TransactionType: LoanBrokerSetTx,
				},
			},
			expected: ErrInvalidAccount,
		},
		{
			name: "fail - VaultID required",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID: "",
			},
			expected: errors.New("LoanBrokerSet: VaultID is required"),
		},
		{
			name: "fail - VaultID invalid",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID: "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F43",
			},
			expected: errors.New("LoanBrokerSet: VaultID must be 64 characters hexadecimal string"),
		},
		{
			name: "fail - Data too long",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID: "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				Data:    func() *string { v := "A" + strings.Repeat("B", 512); return &v }(),
			},
			expected: errors.New("LoanBrokerSet: Data must be a valid non-empty hex string up to 512 characters"),
		},
		{
			name: "fail - ManagementFeeRate too high",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID:           "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				ManagementFeeRate: func() *uint32 { v := uint32(10001); return &v }(),
			},
			expected: errors.New("LoanBrokerSet: ManagementFeeRate must be between 0 and 10000 inclusive"),
		},
		{
			name: "fail - CoverRateMinimum too high",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID:          "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				CoverRateMinimum: func() *uint32 { v := uint32(100001); return &v }(),
			},
			expected: errors.New("LoanBrokerSet: CoverRateMinimum must be between 0 and 100000 inclusive"),
		},
		{
			name: "pass - complete",
			tx: &LoanBrokerSet{
				BaseTx: BaseTx{
					Account:         "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
					TransactionType: LoanBrokerSetTx,
				},
				VaultID:           "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				ManagementFeeRate: func() *uint32 { v := uint32(1000); return &v }(),
				CoverRateMinimum:  func() *uint32 { v := uint32(5000); return &v }(),
			},
			expected: nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			ok, err := testcase.tx.Validate()
			assert.Equal(t, ok, testcase.expected == nil)
			if testcase.expected != nil {
				assert.Contains(t, err.Error(), testcase.expected.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
