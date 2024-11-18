package solana

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

func setup() (chain.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func Test_GetSupportChains(t *testing.T) {
	adaptor := ChainAdaptor{}

	req := &account.SupportChainsRequest{
		Chain:   ChainName,
		Network: "mainnet",
	}

	resp, err := adaptor.GetSupportChains(req)

	if err != nil {
		t.Errorf("GetSupportChains failed with error: %v", err)
	}
	fmt.Printf("resp: %s\n", resp)

	if resp.Code != common2.ReturnCode_SUCCESS {
		t.Errorf("Expected success code, got %v", resp.Code)
	}

	if !resp.Support {
		t.Error("Expected Support to be true")
	}
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	const (
		validPublicKey        = "7e376c64c64e88054b7a2d25dc716f45551d2f796ddc9e7be405e49c522b887c"
		validPublicKeyAddress = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidPublicKey      = "invalid_hex"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: validPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, validPublicKeyAddress, resp.Address)
		assert.Equal(t, "convert address success", resp.Msg)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Converted Address: %s", resp.Address)
	})

	t.Run("Empty Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: "",
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid Public Key Format", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: invalidPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	const (
		validAddress   = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidAddress = "invalid_address"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: validAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, true, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Invalid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: invalidAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Empty Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: "",
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	assert.NoError(t, err, "failed to initialize sol solclient")

	adaptor := &ChainAdaptor{
		solCli: solClient,
	}

	t.Run("Valid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 300944802,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		assert.Equal(t, int64(300944802), resp.Height)

		assert.NotNil(t, resp.Transactions)
		if len(resp.Transactions) > 0 {
			tx := resp.Transactions[0]
			assert.NotEmpty(t, tx.Hash)

			t.Logf("Transaction Hash: %s", tx.Hash)
			t.Logf("From Address: %s", tx.From)
			t.Logf("To Address: %s", tx.To)
			t.Logf("Amount: %s", tx.Amount)
		}

		t.Logf("Block Height: %d", resp.Height)
		t.Logf("Block Hash: %s", resp.Hash)
		t.Logf("Transaction Count: %d", len(resp.Transactions))

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Zero Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 0,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		t.Logf("Genesis Block Hash: %s", resp.Hash)

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Invalid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 999999999999,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.NotEmpty(t, resp.Msg)

		t.Logf("Error Message: %s", resp.Msg)
	})

}

// tx, err := solana.NewTransaction(
// []solana.Instruction{
// system.NewTransferInstruction(
// value,
// fromPubkey,
// toPubkey,
// ).Build(),
// },
// solana.HashFromBytes(binary.BigEndian.AppendUint64(make([]byte, 24), data.Nonce)),
// solana.TransactionPayer(fromPubkey),
func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("CreateUnSignTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.UnSignTx)
}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("TestChainAdaptor_BuildSignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.SignedTx)
}

func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Signature: "4PDQnHWBNSD12Lv2ZL7ZvAN6qxNryKL5rY69q69P2Pa2CDkYy35PFs4uny55pAr8VNL6BfxQaxiuDFYt1hJtqt7C5ZVo6K8zpagrqLdBBxAKiuvE4DEw8gksmGPJaahPGd1RNHbVicMRm5gKPk8vHQ1sJbEwJCM7nCcgYGbQxFvUrCc23BMmmGL2c5KvrQYru7qEfWU8kR8UcZMRmgWDJ7MGp5p85LTW655o6qZ3ohEzmXVf1mSGj7ENxCEeQE6FvhHxH4Xny2Whex85hVVGzztcrBk3NBsgY5SGuJw5q2sPQQZo4x7bNJJuAmySMxZwmM36saEZZPxJPh",
	})
	if err != nil {
		log.Error("TestChainAdaptor_VerifySignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
}
func createTestBase64Tx() string {

	testTx := TxStructure{
		Nonce:       "BDNPPBF4UHF76iMXcsvQXeBr9K7bEyo92GdpWBfMoxb3",
		FromAddress: "7YcpSkLK7gnSJ4JpysHR9BQgwe2gfffRQmMxHDbNf5ve",
		ToAddress:   "EUVrmoaKaSsHNkMFw7mVARR522wwH41BFRMha3WC8gha",
		Value:       "0.001",
		//ContractAddress: "So11111111111111111111111111111111111111112", //5VzPuctbhMdqZBpxgxHCyH41sSckqPEKZ7qxbdgMN29Fbvmnpy3x6GcmUFxFw98oy3LcEEVCxwdr4gyQwcboSW6C
		ContractAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", //3L64aQvAmdhbaZJFdWXTSjLgmH1GwBhNE8eezqCFAHRvj9a76bwXoarivTSjzAJLiJ48CxtZ5Zke3djnfhuckKs
		Signature:       "61999700779d0d8be05b30de4838384f7397e2a7dc85ddc8974f6dd1d51109a2d82c9e1728efa95e7975389f200a3181b58b074a7b501acbc6bf1bd8956ad807",
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
