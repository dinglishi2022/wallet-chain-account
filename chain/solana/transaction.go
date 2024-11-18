package solana

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// GetAccountInfo retrieves account information for a given token account
func GetAccountInfo(sdkClient *rpc.Client, tokenAccount solana.PublicKey) (*rpc.GetAccountInfoResult, error) {
	accountInfo, err := sdkClient.GetAccountInfo(context.Background(), tokenAccount)
	if err != nil {
		log.Error("Failed to get account info", "err", err)
		return nil, err
	}
	return accountInfo, nil
}

// GetTokenSupply retrieves the token supply for a given mint public key
func GetTokenSupply(sdkClient *rpc.Client, mintPubkey solana.PublicKey) (*rpc.GetTokenSupplyResult, error) {
	tokenInfo, err := sdkClient.GetTokenSupply(context.Background(), mintPubkey, rpc.CommitmentFinalized)
	if err != nil {
		log.Error("Failed to get token supply", "err", err)
		return nil, err
	}
	return tokenInfo, nil
}
