package integration

import (
	"os"
	"testing"
)

const (
	// IntegrationEnvVar is the environment variable name used to specify the integration test environment.
	IntegrationEnvVar = "INTEGRATION"
)

// GetWebsocketEnv returns the integration environment.
// If the environment is not set, it skips the tests.
// This function is intended to be used in tests that need to run against a specific environment.
// Run it before creating the runner to retrieve the environment host and faucet provider.
func GetWebsocketEnv(t *testing.T) Env {
	if _, ok := IntegrationWebsocketEnvs[EnvKey(os.Getenv(IntegrationEnvVar))]; !ok {
		t.Skip("skipping integration tests")
	}

	return IntegrationWebsocketEnvs[EnvKey(os.Getenv(IntegrationEnvVar))]
}

// GetRPCEnv returns the integration environment.
// If the environment is not set, it skips the tests.
// This function is intended to be used in tests that need to run against a specific environment.
// Run it before creating the runner to retrieve the environment host and faucet provider.
func GetRPCEnv(t *testing.T) Env {
	if _, ok := IntegrationRPCEnvs[EnvKey(os.Getenv(IntegrationEnvVar))]; !ok {
		t.Skip("skipping integration tests")
	}

	return IntegrationRPCEnvs[EnvKey(os.Getenv(IntegrationEnvVar))]
}

// GetLendingDevnetWebsocketEnv returns the lending devnet websocket environment.
// If the environment is not set to "lendingdevnet", it skips the tests.
// This function is intended for XLS-65/XLS-66 tests that require the lending devnet.
func GetLendingDevnetWebsocketEnv(t *testing.T) Env {
	if os.Getenv(IntegrationEnvVar) != string(LendingDevnetEnv) {
		t.Skip("skipping lending devnet integration tests (set INTEGRATION=lendingdevnet to run)")
	}

	return IntegrationWebsocketEnvs[LendingDevnetEnv]
}

// GetLendingDevnetRPCEnv returns the lending devnet RPC environment.
// If the environment is not set to "lendingdevnet", it skips the tests.
// This function is intended for XLS-65/XLS-66 tests that require the lending devnet.
func GetLendingDevnetRPCEnv(t *testing.T) Env {
	if os.Getenv(IntegrationEnvVar) != string(LendingDevnetEnv) {
		t.Skip("skipping lending devnet integration tests (set INTEGRATION=lendingdevnet to run)")
	}

	return IntegrationRPCEnvs[LendingDevnetEnv]
}
