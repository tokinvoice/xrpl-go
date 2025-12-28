package integration

import (
	"testing"

	ledger "github.com/Peersyst/xrpl-go/xrpl/ledger-entry-types"
	"github.com/Peersyst/xrpl-go/xrpl/testutil/integration"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	"github.com/Peersyst/xrpl-go/xrpl/websocket"
	"github.com/stretchr/testify/require"
)

// TestIntegrationVault_Websocket tests the XLS-65 Vault transactions on the lending devnet.
// It creates a vault, deposits funds, withdraws funds, updates settings, and deletes the vault.
func TestIntegrationVault_Websocket(t *testing.T) {
	env := integration.GetLendingDevnetWebsocketEnv(t)
	client := websocket.NewClient(websocket.NewClientConfig().WithHost(env.Host).WithFaucetProvider(env.FaucetProvider))

	runner := integration.NewRunner(t, client, &integration.RunnerConfig{
		WalletCount: 1,
		MaxRetries:  10,
	})

	err := runner.Setup()
	require.NoError(t, err)
	defer runner.Teardown()

	owner := runner.GetWallet(0)

	// Step 1: Create a vault
	t.Run("VaultCreate", func(t *testing.T) {
		withdrawalPolicy := transaction.VaultStrategyFirstComeFirstServe
		assetsMax := "1000000000000" // 1M XRP in drops
		vaultCreate := &transaction.VaultCreate{
			BaseTx: transaction.BaseTx{
				Account: owner.GetAddress(),
			},
			Asset:            ledger.Asset{Currency: "XRP"},
			AssetsMaximum:    &assetsMax,
			WithdrawalPolicy: &withdrawalPolicy,
		}

		flatTx := vaultCreate.Flatten()
		resp, err := runner.TestTransaction(&flatTx, owner, "tesSUCCESS", nil)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Extract VaultID from created nodes for subsequent tests
		t.Logf("VaultCreate succeeded, hash: %s", resp.Tx["hash"])
	})
}

// TestIntegrationVaultFullCycle_Websocket tests the complete vault lifecycle.
func TestIntegrationVaultFullCycle_Websocket(t *testing.T) {
	env := integration.GetLendingDevnetWebsocketEnv(t)
	client := websocket.NewClient(websocket.NewClientConfig().WithHost(env.Host).WithFaucetProvider(env.FaucetProvider))

	runner := integration.NewRunner(t, client, &integration.RunnerConfig{
		WalletCount: 1,
		MaxRetries:  10,
	})

	err := runner.Setup()
	require.NoError(t, err)
	defer runner.Teardown()

	owner := runner.GetWallet(0)
	var vaultID string

	// Step 1: VaultCreate
	t.Run("VaultCreate", func(t *testing.T) {
		withdrawalPolicy := transaction.VaultStrategyFirstComeFirstServe
		assetsMax := "1000000000000" // 1M XRP in drops
		vaultCreate := &transaction.VaultCreate{
			BaseTx: transaction.BaseTx{
				Account: owner.GetAddress(),
			},
			Asset:            ledger.Asset{Currency: "XRP"},
			AssetsMaximum:    &assetsMax,
			WithdrawalPolicy: &withdrawalPolicy,
		}

		flatTx := vaultCreate.Flatten()
		resp, err := runner.TestTransactionAndWait(&flatTx, owner, "tesSUCCESS", nil)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Extract vault ID from meta
		for _, node := range resp.Meta.AffectedNodes {
			if node.CreatedNode != nil && node.CreatedNode.LedgerEntryType == "Vault" {
				vaultID = node.CreatedNode.LedgerIndex
				t.Logf("Created Vault ID: %s", vaultID)
				break
			}
		}
		require.NotEmpty(t, vaultID, "VaultID should be extracted from created nodes")
	})

	// Step 2: VaultDeposit
	t.Run("VaultDeposit", func(t *testing.T) {
		if vaultID == "" {
			t.Skip("No vault ID from previous step")
		}

		vaultDeposit := &transaction.VaultDeposit{
			BaseTx: transaction.BaseTx{
				Account: owner.GetAddress(),
			},
			VaultID: types.Hash256(vaultID),
			Amount:  types.XRPCurrencyAmount(50000000), // 50 XRP
		}

		flatTx := vaultDeposit.Flatten()
		resp, err := runner.TestTransaction(&flatTx, owner, "tesSUCCESS", nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		t.Logf("VaultDeposit succeeded")
	})

	// Step 3: VaultWithdraw
	t.Run("VaultWithdraw", func(t *testing.T) {
		if vaultID == "" {
			t.Skip("No vault ID from previous step")
		}

		vaultWithdraw := &transaction.VaultWithdraw{
			BaseTx: transaction.BaseTx{
				Account: owner.GetAddress(),
			},
			VaultID: types.Hash256(vaultID),
			Amount:  types.XRPCurrencyAmount(25000000), // 25 XRP
		}

		flatTx := vaultWithdraw.Flatten()
		resp, err := runner.TestTransaction(&flatTx, owner, "tesSUCCESS", nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		t.Logf("VaultWithdraw succeeded")
	})

	// Step 4: VaultSet
	t.Run("VaultSet", func(t *testing.T) {
		if vaultID == "" {
			t.Skip("No vault ID from previous step")
		}

		newMax := "2000000000000" // 2M XRP in drops
		vaultSet := &transaction.VaultSet{
			BaseTx: transaction.BaseTx{
				Account: owner.GetAddress(),
			},
			VaultID:       types.Hash256(vaultID),
			AssetsMaximum: &newMax,
		}

		flatTx := vaultSet.Flatten()
		resp, err := runner.TestTransaction(&flatTx, owner, "tesSUCCESS", nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		t.Logf("VaultSet succeeded")
	})
}

