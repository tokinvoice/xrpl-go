package ledger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoan_EntryType(t *testing.T) {
	entry := &Loan{}
	assert.Equal(t, LoanEntry, entry.EntryType())
}

func TestLoan_SetLsfLoanDefault(t *testing.T) {
	loan := &Loan{}
	assert.False(t, loan.HasLsfLoanDefault())
	loan.SetLsfLoanDefault()
	assert.True(t, loan.HasLsfLoanDefault())
	assert.Equal(t, uint32(0x00010000), loan.Flags)
}

func TestLoan_SetLsfLoanImpaired(t *testing.T) {
	loan := &Loan{}
	assert.False(t, loan.HasLsfLoanImpaired())
	loan.SetLsfLoanImpaired()
	assert.True(t, loan.HasLsfLoanImpaired())
	assert.Equal(t, uint32(0x00020000), loan.Flags)
}

func TestLoan_SetLsfLoanOverpayment(t *testing.T) {
	loan := &Loan{}
	assert.False(t, loan.HasLsfLoanOverpayment())
	loan.SetLsfLoanOverpayment()
	assert.True(t, loan.HasLsfLoanOverpayment())
	assert.Equal(t, uint32(0x00040000), loan.Flags)
}

