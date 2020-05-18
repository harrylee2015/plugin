package types

import (
	"encoding/hex"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initCryptoImpl() error {
	cr, err := crypto.New(types.GetSignName("", types.ED25519))
	if err != nil {
		return err
	}
	ConsensusCrypto = cr
	return nil
}
func TestDetection_Sign(t *testing.T) {
	err := initCryptoImpl()
	assert.Nil(t, err)
	doc, err := GenesisDocFromFile("../genesis.json")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(doc.Validators))
	privValidator := LoadOrGenPrivValidatorFS("../priv_validator.json")
	pubkey, err := PubKeyFromString(doc.Validators[0].PubKey.Data)
	assert.Nil(t, err)
	detection := &Detection{PubKeyBytes: hex.EncodeToString(pubkey.Bytes()), Address: "xxxx", PeerID: "node1"}
	err = detection.Sign(privValidator.PrivKey)
	assert.Nil(t, err)
	ok := detection.CheckSign([]*Validator{&Validator{PubKey: pubkey.Bytes()}})
	assert.True(t, ok)
}
