package transaction

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanSet_TxType(t *testing.T) {
	tx := &LoanSet{}
	assert.Equal(t, tx.TxType(), LoanSetTx)
}

func TestLoanSet_Flatten(t *testing.T) {
	testcases := []struct {
		name     string
		tx       *LoanSet
		expected map[string]interface{}
	}{
		{
			name: "pass - empty",
			tx:   &LoanSet{},
			expected: map[string]interface{}{
				"TransactionType":    LoanSetTx.String(),
				"LoanBrokerID":       "",
				"PrincipalRequested": "",
			},
		},
		{
			name: "pass - complete",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:            "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					Fee:                1000000,
					Sequence:           1,
					LastLedgerSequence: 3000000,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "100000",
				Counterparty:       "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
				InterestRate:       5000,
				PaymentInterval:    2592000,
			},
			expected: map[string]interface{}{
				"TransactionType":    LoanSetTx.String(),
				"Account":            "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
				"Fee":                "1000000",
				"Sequence":           uint32(1),
				"LastLedgerSequence": uint32(3000000),
				"LoanBrokerID":       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				"PrincipalRequested": "100000",
				"Counterparty":       "rNZ9m6AP9K7z3EVg6GhPMx36V4QmZKeWds",
				"InterestRate":       uint32(5000),
				"PaymentInterval":    uint32(2592000),
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert.Equal(t, testcase.tx.Flatten(), testcase.expected)
		})
	}
}

func TestLoanSet_Validate(t *testing.T) {
	testcases := []struct {
		name     string
		tx       *LoanSet
		expected error
	}{
		{
			name: "fail - base tx invalid",
			tx: &LoanSet{
				BaseTx: BaseTx{
					TransactionType: LoanSetTx,
				},
			},
			expected: ErrInvalidAccount,
		},
		{
			name: "fail - LoanBrokerID required",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				PrincipalRequested: "100000",
			},
			expected: errors.New("LoanSet: LoanBrokerID is required"),
		},
		{
			name: "fail - PrincipalRequested required",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID: "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
			},
			expected: errors.New("LoanSet: PrincipalRequested is required"),
		},
		{
			name: "fail - LoanBrokerID invalid",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F43",
				PrincipalRequested: "100000",
			},
			expected: errors.New("LoanSet: LoanBrokerID must be 64 characters hexadecimal string"),
		},
		{
			name: "fail - PrincipalRequested invalid",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "invalid",
			},
			expected: errors.New("LoanSet: PrincipalRequested must be a valid XRPL number"),
		},
		{
			name: "fail - Data too long",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "100000",
				Data:               "A" + strings.Repeat("B", 512),
			},
			expected: errors.New("LoanSet: Data must be a valid non-empty hex string up to 512 characters"),
		},
		{
			name: "fail - OverpaymentFee too high",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "100000",
				OverpaymentFee:     100001,
			},
			expected: errors.New("LoanSet: OverpaymentFee must be between 0 and 100000 inclusive"),
		},
		{
			name: "fail - PaymentInterval too low",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "100000",
				PaymentInterval:    59,
			},
			expected: errors.New("LoanSet: PaymentInterval must be greater than or equal to 60"),
		},
		{
			name: "pass - complete",
			tx: &LoanSet{
				BaseTx: BaseTx{
					Account:         "rHLLL3Z7uBLK49yZcMaj8FAP7DU12Nw5A5",
					TransactionType: LoanSetTx,
				},
				LoanBrokerID:       "B91CD2033E73E0DD17AF043FBD458CE7D996850A83DCED23FB122A3BFAA7F430",
				PrincipalRequested: "100000",
				InterestRate:       5000,
				PaymentInterval:    2592000,
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

func TestLoanSet_Flags(t *testing.T) {
	tests := []struct {
		name     string
		setter   func(*LoanSet)
		expected uint32
	}{
		{
			name: "pass - SetLoanOverpaymentFlag",
			setter: func(ls *LoanSet) {
				ls.SetLoanOverpaymentFlag()
			},
			expected: tfLoanOverpayment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoanSet{}
			tt.setter(ls)
			if ls.Flags != tt.expected {
				t.Errorf("Expected LoanSet Flags to be %d, got %d", tt.expected, ls.Flags)
			}
		})
	}
}
