package faucet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
)

const (
	// LendingDevnetFaucetHost is the hostname for the XRPL Lending Devnet faucet service.
	LendingDevnetFaucetHost = "lend-faucet.devnet.rippletest.net"
	// LendingDevnetFaucetPath is the API path for account operations on the Lending Devnet faucet.
	LendingDevnetFaucetPath = "/accounts"
)

// LendingDevnetFaucetProvider implements the FaucetProvider interface for the XRPL Lending Devnet.
// It provides functionality to interact with the Lending Devnet faucet for funding wallets.
type LendingDevnetFaucetProvider struct {
	host        string // The hostname of the Lending Devnet faucet
	accountPath string // The API path for account-related operations
}

// NewLendingDevnetFaucetProvider creates and returns a new instance of LendingDevnetFaucetProvider
// with predefined Lending Devnet faucet host and account path.
func NewLendingDevnetFaucetProvider() *LendingDevnetFaucetProvider {
	return &LendingDevnetFaucetProvider{
		host:        LendingDevnetFaucetHost,
		accountPath: LendingDevnetFaucetPath,
	}
}

// FundWallet sends a request to the Lending Devnet faucet to fund the specified wallet address.
// It returns an error if the funding request fails.
func (fp *LendingDevnetFaucetProvider) FundWallet(address types.Address) error {
	url := fmt.Sprintf("https://%s%s", fp.host, fp.accountPath)
	payload := map[string]string{"destination": address.String()}
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return ErrMarshalPayload{
			Err: err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ErrCreateRequest{Err: err}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrSendRequest{Err: err}
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode{
			Code: resp.StatusCode,
		}
	}

	return nil
}

