package ledger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanBroker_EntryType(t *testing.T) {
	entry := &LoanBroker{}
	assert.Equal(t, LoanBrokerEntry, entry.EntryType())
}

