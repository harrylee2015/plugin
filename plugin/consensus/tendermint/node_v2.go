// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tendermint

import (
	"encoding/hex"
	"fmt"
	"github.com/33cn/chain33/types"
	ttypes "github.com/33cn/plugin/plugin/consensus/tendermint/types"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/33cn/chain33/common/crypto"
)

// Node struct
type NodeV2 struct {
	privKey crypto.PrivKey
	Network string
	Version string
	ID      ID
	Address string
	peerSet *PeerSetV2
	//嗅探包信息
	receiveDetectMsgChannel  <-chan *ttypes.Detection
	sendMsgToP2pChannel      chan<- *types.ConsensusMsg
	receiveMsgFromP2pChannel <-chan *types.ConsensusMsg
	receiveMsgFromPeerSet    chan *types.ConsensusMsg

	state            *ConsensusState
	broadcastChannel chan MsgInfo
	started          uint32 // atomic
	stopped          uint32 // atomic
	quit             chan struct{}
}

// NewNode method
func NewNodeV2(detectChannel <-chan *ttypes.Detection, sendMsgToP2pChannel chan<- *types.ConsensusMsg, receiveMsgFromP2pChannel <-chan *types.ConsensusMsg, privKey crypto.PrivKey, network string, version string, state *ConsensusState) *NodeV2 {
	address := GenAddressByPubKey(privKey.PubKey())
	node := &NodeV2{
		peerSet:                  NewPeerSetV2(),
		receiveDetectMsgChannel:  detectChannel,
		sendMsgToP2pChannel:      sendMsgToP2pChannel,
		receiveMsgFromP2pChannel: receiveMsgFromP2pChannel,
		receiveMsgFromPeerSet:    make(chan *types.ConsensusMsg),
		privKey:                  privKey,
		Network:                  network,
		Version:                  version,
		ID:                       ID(hex.EncodeToString(address)),
		broadcastChannel:         make(chan MsgInfo, maxSendQueueSize),
		state:                    state,
		quit:                     make(chan struct{}),
	}
	state.SetOurID(node.ID)
	state.SetBroadcastChannel(node.broadcastChannel)
	return node
}

// Start node
func (node *NodeV2) Start() {
	if atomic.CompareAndSwapUint32(&node.started, 0, 1) {

		//寻找，并维护可用的共识节点网络
		go node.ListenPeerSet()
		go node.StartConsensusRoutine()
		go node.BroadcastRoutine()
		go node.Relay()
	}
}

// Stop ...
func (node *NodeV2) Stop() {
	atomic.CompareAndSwapUint32(&node.stopped, 0, 1)
	if node.quit != nil {
		close(node.quit)
	}
	// Stop peers
	for _, peer := range node.peerSet.List() {
		peer.Stop()
		node.peerSet.Remove(peer)
	}
	//stop consensus
	node.state.Stop()
}

// IsRunning ...
func (node *NodeV2) IsRunning() bool {
	return atomic.LoadUint32(&node.started) == 1 && atomic.LoadUint32(&node.stopped) == 0
}

// StartConsensusRoutine if peers reached the threshold start consensus routine
func (node *NodeV2) StartConsensusRoutine() {
	for {
		//TODO:the peer count need be optimized
		if node.peerSet.Size() > 0 {
			node.state.Start()
			tendermintlog.Debug("===================StartConsensusRoutine=====================")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// BroadcastRoutine receive to broadcast
func (node *NodeV2) BroadcastRoutine() {
	for {
		msg, ok := <-node.broadcastChannel
		if !ok {
			tendermintlog.Debug("broadcastChannel closed")
			return
		}
		node.Broadcast(msg)
	}
}

//中继器，用于中转从p2p模块发来得消息和发往p2p模块
func (node *NodeV2) Relay() {
	go func() {
		for msg := range node.receiveMsgFromP2pChannel {
			if node.peerSet.Has(ID(msg.FromPeerID)) {
				node.peerSet.GetPeer(ID(msg.FromPeerID)).ReceiveMsg(msg)
			} else {
				tendermintlog.Error("not found peerconn", "peerID", msg.FromPeerID)
			}
		}
	}()
	//从peerSet接收消息发送给p2p
	go func() {
		for msg := range node.receiveMsgFromPeerSet {
			msg.FromPeerID = string(node.ID)
			node.sendMsgToP2pChannel <- msg
		}
	}()
}

func (node *NodeV2) connectComming(inConn net.Conn) {
	maxPeers := maxNumPeers
	if maxPeers <= node.peerSet.Size() {
		tendermintlog.Debug("Ignoring inbound connection: already have enough peers", "address", inConn.RemoteAddr().String(), "numPeers", node.peerSet.Size(), "max", maxPeers)
		return
	}
}

func (node *NodeV2) stopAndRemovePeer(peer PeerV2, reason interface{}) {
	node.peerSet.Remove(peer)
	peer.Stop()
}

//根据嗅探消息,维护PeerSet集合
func (node *NodeV2) ListenPeerSet() {
	//优雅读取channel数据
	for detect := range node.receiveDetectMsgChannel {
		//验证签名,嗅探包有效时间控制在两秒内
		tendermintlog.Debug("node.state valdators", "num", len(node.state.Validators.Validators))
		if (detect.ExpireTime+2) >= time.Now().Unix() && detect.CheckSign(node.state.Validators.Validators) {
			//添加共识地址与p2p节点ID之间得映射关系
			tendermintlog.Debug("I receive detection", "peer_ID", detect.PeerID, "address", detect.Address)
			//不存在需要添加
			if string(node.ID) != detect.Address && !node.peerSet.Has(ID(detect.Address)) {
				node.addPeer(newPeerConnV2(node.receiveMsgFromPeerSet, node.state, ID(detect.Address), detect.PeerID, detect.PeerIP))
				tendermintlog.Debug("I have  add peerconn", "peer_ID", detect.PeerID, "address", detect.Address)
			} else {
				tendermintlog.Debug("Ignoring inbound connection: already have enough peers", "address", detect.Address, "numPeers", node.peerSet.Size())
			}
		}

	}

}

// addPeer checks the given peer's validity, performs a handshake, and adds the
// peer to the switch and to all registered reactors.
// NOTE: This performs a blocking handshake before the peer is added.
// NOTE: If error is returned, caller is responsible for calling peer.CloseConn()
func (node *NodeV2) addPeer(pc *peerConnV2) error {
	// Avoid self
	if node.ID == pc.id {
		return fmt.Errorf("Connect to self: %v", node.ID)
	}

	// Avoid duplicate
	if node.peerSet.Has(pc.id) {
		return fmt.Errorf("Duplicate peer ID %v", pc.id)
	}

	// All good. Start peer
	if node.IsRunning() {
		tendermintlog.Info("start peer", "peer", pc.id)
		pc.SetTransferChannel(node.state.peerMsgQueue)
		if err := node.startInitPeer(pc); err != nil {
			return err
		}
	}

	if err := node.peerSet.Add(pc); err != nil {
		return err
	}

	tendermintlog.Info("Added peer", "peer", pc.id)
	return nil
}

// Broadcast to peers in set
func (node *NodeV2) Broadcast(msg MsgInfo) chan bool {
	successChan := make(chan bool, len(node.peerSet.List()))
	tendermintlog.Debug("Broadcast", "msgtype", msg.TypeID)
	var wg sync.WaitGroup
	for _, peer := range node.peerSet.List() {
		wg.Add(1)
		go func(peer PeerV2) {
			defer wg.Done()
			success := peer.Send(msg)
			successChan <- success
		}(peer)
	}
	go func() {
		wg.Wait()
		close(successChan)
	}()
	return successChan
}

func (node *NodeV2) startInitPeer(peer *peerConnV2) error {
	err := peer.Start() // spawn send/recv routines
	if err != nil {
		// Should never happen
		tendermintlog.Error("Error starting peer", "peer", peer, "err", err)
		return err
	}

	return nil
}

// FilterConnByAddr TODO:can make fileter by addr
func (node *NodeV2) FilterConnByAddr(addr net.Addr) error {
	return nil
}

// CompatibleWith one node by nodeInfo
func (node *NodeV2) CompatibleWith(other NodeInfo) error {
	iMajor, iMinor, _, iErr := splitVersion(node.Version)
	oMajor, oMinor, _, oErr := splitVersion(other.Version)

	// if our own version number is not formatted right, we messed up
	if iErr != nil {
		return iErr
	}

	// version number must be formatted correctly ("x.x.x")
	if oErr != nil {
		return oErr
	}

	// major version must match
	if iMajor != oMajor {
		return fmt.Errorf("Peer is on a different major version. Got %v, expected %v", oMajor, iMajor)
	}

	// minor version can differ
	if iMinor != oMinor {
		// ok
	}

	// nodes must be on the same network
	if node.Network != other.Network {
		return fmt.Errorf("Peer is on a different network. Got %v, expected %v", other.Network, node.Network)
	}

	return nil
}

func newPeerConnV2(
	sendChan chan<- *types.ConsensusMsg,
	state *ConsensusState,
	id ID,
	peerID string,
	peerIP string,
) *peerConnV2 {
	// Only the information we already have

	return &peerConnV2{
		sendConMsgChannel:    sendChan,
		receiveConMsgChannel: make(chan *types.ConsensusMsg),
		myState:              state,
		id:                   id,
		peerID:               peerID,
		ip:                   net.ParseIP(peerIP),
	}
}
