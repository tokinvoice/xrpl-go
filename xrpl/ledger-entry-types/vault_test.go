package ledger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVault_EntryType(t *testing.T) {
	entry := &Vault{}
	assert.Equal(t, VaultEntry, entry.EntryType())
}

func TestVault_SetLsfVaultPrivate(t *testing.T) {
	vault := &Vault{}
	assert.False(t, vault.HasLsfVaultPrivate())
	vault.SetLsfVaultPrivate()
	assert.True(t, vault.HasLsfVaultPrivate())
	assert.Equal(t, uint32(0x00010000), vault.Flags)
}

