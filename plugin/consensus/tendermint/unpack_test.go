package tendermint

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPack(t *testing.T) {

	packet1 := &msgPacket{Bytes: []byte("123456abcd"), TypeID: byte(0x01)}
	data, err := json.Marshal(&packet1)
	assert.Equal(t, nil, err)
	var packet2 msgPacket
	err = json.Unmarshal(data, &packet2)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("123456abcd"), packet2.Bytes)

}
