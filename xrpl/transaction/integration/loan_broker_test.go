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

// TestIntegrationLoanBroker_Websocket tests the XLS-66 LoanBroker transactions on the lending devnet.
// It creates a vault, attaches a loan broker, deposits cover, withdraws cover, and deletes the broker.
func TestIntegrationLoanBroker_Websocket(t *testing.T) {
	env := integration.GetLendingDevnetWebsocketEnv(t)
	client := websocket.NewClient(websocket.NewClientConfig().WithHost(env.Host).WithFaucetProvider(env.FaucetProvider))

	runner := integration.NewRunner(t, client, &integration.RunnerConfig{
		WalletCount: 1,
		MaxRetries:  10,
	})

	err := runner.Setup()
	require.NoError(t, err)
	defer runner.Teardown()

	lender := runner.GetWallet(0)
	var vaultID string
	var loanBrokerID string

	// Step 1: Create a vault (prerequisite for LoanBroker)
	t.Run("VaultCreate", func(t *testing.T) {
		withdrawalPolicy := transaction.VaultStrategyFirstComeFirstServe
		assetsMax := "1000000000000" // 1M XRP
		vaultCreate := &transaction.VaultCreate{
			BaseTx: transaction.BaseTx{
				Account: lender.GetAddress(),
			},
			Asset:            ledger.Asset{Currency: "XRP"},
			AssetsMaximum:    &assetsMax,
			WithdrawalPolicy: &withdrawalPolicy,
		}

		flatTx := vaultCreate.Flatten()
		resp, err := runner.TestTransactionAndWait(&flatTx, lender, "tesSUCCESS", nil)
		require.NoError(t, err)

		// Extract vault ID from meta
		for _, node := range resp.Meta.AffectedNodes {
			if node.CreatedNode != nil && node.CreatedNode.LedgerEntryType == "Vault" {
				vaultID = node.CreatedNode.LedgerIndex
				t.Logf("Created Vault ID: %s", vaultID)
				break
			}
		}
		require.NotEmpty(t, vaultID)
	})

	// Step 2: Deposit into vault
	t.Run("VaultDeposit", func(t *testing.T) {
		if vaultID == "" {
			t.Skip("No vault ID")
		}

		vaultDeposit := &transaction.VaultDeposit{
			BaseTx: transaction.BaseTx{
				Account: lender.GetAddress(),
			},
			VaultID: types.Hash256(vaultID),
			Amount:  types.XRPCurrencyAmount(50000000), // 50 XRP
		}

		flatTx := vaultDeposit.Flatten()
		_, err := runner.TestTransaction(&flatTx, lender, "tesSUCCESS", nil)
		require.NoError(t, err)
		t.Logf("Deposited 50 XRP into vault")
	})

	// Step 3: LoanBrokerSet - Create loan broker
	t.Run("LoanBrokerSet", func(t *testing.T) {
		if vaultID == "" {
			t.Skip("No vault ID")
		}

		debtMax := "100000000000" // 100k XRP in drops
		mgmtFee := uint32(100)    // 0.1%
		coverMin := uint32(5000)  // 5%
		coverLiq := uint32(2500)  // 2.5%
		loanBrokerSet := &transaction.LoanBrokerSet{
			BaseTx: transaction.BaseTx{
				Account: lender.GetAddress(),
			},
			VaultID:              types.Hash256(vaultID),
			ManagementFeeRate:    &mgmtFee,
			DebtMaximum:          &debtMax,
			CoverRateMinimum:     &coverMin,
			CoverRateLiquidation: &coverLiq,
		}

		flatTx := loanBrokerSet.Flatten()
		resp, err := runner.TestTransactionAndWait(&flatTx, lender, "tesSUCCESS", nil)
		require.NoError(t, err)

		// Extract LoanBroker ID from meta
		for _, node := range resp.Meta.AffectedNodes {
			if node.CreatedNode != nil && node.CreatedNode.LedgerEntryType == "LoanBroker" {
				loanBrokerID = node.CreatedNode.LedgerIndex
				t.Logf("Created LoanBroker ID: %s", loanBrokerID)
				break
			}
		}
		require.NotEmpty(t, loanBrokerID)
	})

	// Step 4: LoanBrokerCoverDeposit - Deposit first-loss capital
	t.Run("LoanBrokerCoverDeposit", func(t *testing.T) {
		if loanBrokerID == "" {
			t.Skip("No LoanBroker ID")
		}

		coverDeposit := &transaction.LoanBrokerCoverDeposit{
			BaseTx: transaction.BaseTx{
				Account: lender.GetAddress(),
			},
			LoanBrokerID: types.Hash256(loanBrokerID),
			Amount:       types.XRPCurrencyAmount(10000000), // 10 XRP
		}

		flatTx := coverDeposit.Flatten()
		_, err := runner.TestTransaction(&flatTx, lender, "tesSUCCESS", nil)
		require.NoError(t, err)
		t.Logf("Deposited 10 XRP as first-loss capital")
	})

	// Step 5: LoanBrokerCoverWithdraw - Withdraw some first-loss capital
	t.Run("LoanBrokerCoverWithdraw", func(t *testing.T) {
		if loanBrokerID == "" {
			t.Skip("No LoanBroker ID")
		}

		coverWithdraw := &transaction.LoanBrokerCoverWithdraw{
			BaseTx: transaction.BaseTx{
				Account: lender.GetAddress(),
			},
			LoanBrokerID: types.Hash256(loanBrokerID),
			Amount:       types.XRPCurrencyAmount(5000000), // 5 XRP
		}

		flatTx := coverWithdraw.Flatten()
		_, err := runner.TestTransaction(&flatTx, lender, "tesSUCCESS", nil)
		require.NoError(t, err)
		t.Logf("Withdrew 5 XRP first-loss capital")
	})
}

