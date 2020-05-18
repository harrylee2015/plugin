package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/33cn/chain33/common/crypto"
)

//嗅探包，用于全网广播本验证节点的公钥与节点p2p ID之间的映射关系
type Detection struct {
	//节点P2P地址，用于p2p消息发送
	PeerID string `json:"peerid"`
	//PeerIp
	PeerIP string `json:"peerip"`
	//公钥地址，也就是节点的ID
	Address string `json:"address"`
	//PubKey    KeyText `json:"pub_key"`
	PubKeyBytes string `json:"pubkeybytes"`
	SignBytes   string `json:"signbytes,omitempty"`
	ExpireTime  int64  `json:"expiretime"`
}

//Sign detection签名
func (d *Detection) Sign(priv crypto.PrivKey) error {
	d.SignBytes = ""
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	sign := priv.Sign(data)
	d.SignBytes = hex.EncodeToString(sign.Bytes())
	return nil
}

func (d *Detection) CheckSign(validators []*Validator) bool {
	if d.SignBytes == "" {
		return false
	}
	copydt := *d
	copydt.SignBytes = ""
	data, err := json.Marshal(&copydt)
	if err != nil {
		return false
	}
	pubkeyBytes, err := hex.DecodeString(copydt.PubKeyBytes)
	if err != nil {
		return false
	}
	//寻找授权的节点
	for _, validator := range validators {
		if bytes.Equal(validator.PubKey, pubkeyBytes) {
			return checkSign(data, validator.PubKey, d.SignBytes)
		}
	}
	return false
}

//验证签名
func checkSign(data []byte, pubKeyBytes []byte, signStr string) bool {
	pubkey, err := PubKeyFromBytes(pubKeyBytes)
	if err != nil {
		return false
	}
	signature, err := SignatureFromString(signStr)
	if err != nil {
		return false
	}
	return pubkey.VerifyBytes(data, signature)
}
