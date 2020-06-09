// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tendermint

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"net"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	ttypes "github.com/33cn/plugin/plugin/consensus/tendermint/types"
	tmtypes "github.com/33cn/plugin/plugin/dapp/valnode/types"
	"github.com/golang/protobuf/proto"
)

// Peer interface
type PeerV2 interface {
	ID() ID
	RemotePeerID() string
	IsOutbound() bool
	IsPersistent() bool

	Send(msg MsgInfo) bool
	TrySend(msg MsgInfo) bool

	Stop()

	SetTransferChannel(chan MsgInfo)
	ReceiveMsg(msg *types.ConsensusMsg)
	SendMsg(msg *types.ConsensusMsg)

	SetExpireTime(expireTime int64)
	CheckExpire() bool
}

// PeerConnState struct
type PeerConnStateV2 struct {
	mtx sync.Mutex
	id  ID
	ip  net.IP
	ttypes.PeerRoundState
}

type peerConnV2 struct {
	outbound             bool
	receiveConMsgChannel chan *types.ConsensusMsg
	sendConMsgChannel    chan<- *types.ConsensusMsg
	persistent           bool
	pubKey               crypto.PubKey
	ip                   net.IP
	id                   ID
	peerID               string

	sendQueue     chan MsgInfo
	sendQueueSize int32
	pongChannel   chan struct{}

	started uint32 //atomic
	stopped uint32 // atomic

	quitSend   chan struct{}
	quitUpdate chan struct{}
	quitBeat   chan struct{}
	waitQuit   sync.WaitGroup

	transferChannel chan MsgInfo

	sendBuffer []byte

	onPeerError func(PeerV2, interface{})

	myState *ConsensusState

	state            *PeerConnState
	updateStateQueue chan MsgInfo
	heartbeatQueue   chan proto.Message
	//有效期，原子锁,用于控制expireTime并发读写可能出现的资源竞争
	expireTime *int64
}

// PeerSet struct
type PeerSetV2 struct {
	mtx    sync.Mutex
	lookup map[ID]*peerSetItemV2
	list   []PeerV2
}

type peerSetItemV2 struct {
	peer  PeerV2
	index int
}

// NewPeerSet method
func NewPeerSetV2() *PeerSetV2 {
	return &PeerSetV2{
		lookup: make(map[ID]*peerSetItemV2),
		list:   make([]PeerV2, 0, 50),
	}
}

// Add adds the peer to the PeerSet.
// It returns an error carrying the reason, if the peer is already present.
func (ps *PeerSetV2) Add(peer PeerV2) error {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.lookup[peer.ID()] != nil {
		return fmt.Errorf("Duplicate peer ID %v", peer.ID())
	}

	index := len(ps.list)
	// Appending is safe even with other goroutines
	// iterating over the ps.list slice.
	ps.list = append(ps.list, peer)
	ps.lookup[peer.ID()] = &peerSetItemV2{peer, index}
	return nil
}

// Has returns true if the set contains the peer referred to by this
// peerKey, otherwise false.
func (ps *PeerSetV2) Has(peerKey ID) bool {
	ps.mtx.Lock()
	_, ok := ps.lookup[peerKey]
	ps.mtx.Unlock()
	return ok
}

func (ps *PeerSetV2) GetPeer(peerKey ID) PeerV2 {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	itemV2, ok := ps.lookup[peerKey]
	if ok {
		return itemV2.peer
	}
	return nil
}

// Size of list
func (ps *PeerSetV2) Size() int {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return len(ps.list)
}

// List returns the threadsafe list of peers.
func (ps *PeerSetV2) List() []PeerV2 {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return ps.list
}

// Remove discards peer by its Key, if the peer was previously memoized.
func (ps *PeerSetV2) Remove(peer PeerV2) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	item := ps.lookup[peer.ID()]
	if item == nil {
		return
	}
	peer.Stop()
	index := item.index
	// Create a new copy of the list but with one less item.
	// (we must copy because we'll be mutating the list).
	newList := make([]PeerV2, len(ps.list)-1)
	copy(newList, ps.list)
	// If it's the last peer, that's an easy special case.
	if index == len(ps.list)-1 {
		ps.list = newList
		delete(ps.lookup, peer.ID())
		return
	}

	// Replace the popped item with the last item in the old list.
	lastPeer := ps.list[len(ps.list)-1]
	lastPeerKey := lastPeer.ID()
	lastPeerItem := ps.lookup[lastPeerKey]
	newList[index] = lastPeer
	lastPeerItem.index = index
	ps.list = newList
	delete(ps.lookup, peer.ID())
}

//-------------------------peer connection--------------------------------
func (pc *peerConnV2) ID() ID {
	if len(pc.id) != 0 {
		return pc.id
	}
	address := GenAddressByPubKey(pc.pubKey)
	pc.id = ID(hex.EncodeToString(address))
	return pc.id
}

func (pc *peerConnV2) RemotePeerID() string {
	return pc.peerID
}

func (pc *peerConnV2) ReceiveMsg(msg *types.ConsensusMsg) {
	pc.receiveConMsgChannel <- msg
}
func (pc *peerConnV2) SendMsg(msg *types.ConsensusMsg) {
	pc.sendConMsgChannel <- msg
}

func (pc *peerConnV2) SetTransferChannel(transferChannel chan MsgInfo) {
	pc.transferChannel = transferChannel
}

func (pc *peerConnV2) SetExpireTime(expireTime int64) {
	atomic.StoreInt64(pc.expireTime,expireTime)
}
func (pc *peerConnV2) CheckExpire() bool {
	val:=atomic.LoadInt64(pc.expireTime)
	return val <= time.Now().Unix()
}
func (pc *peerConnV2) CloseConn() {
	close(pc.receiveConMsgChannel)
	close(pc.sendConMsgChannel)
}

func (pc *peerConnV2) IsOutbound() bool {
	return pc.outbound
}

func (pc *peerConnV2) IsPersistent() bool {
	return pc.persistent
}

func (pc *peerConnV2) Send(msg MsgInfo) bool {
	if !pc.IsRunning() {
		return false
	}
	select {
	case pc.sendQueue <- msg:
		atomic.AddInt32(&pc.sendQueueSize, 1)
		return true
	case <-time.After(defaultSendTimeout):
		tendermintlog.Error("send msg timeout", "peerip", msg.PeerIP, "msg", msg.Msg)
		return false
	}
}

func (pc *peerConnV2) TrySend(msg MsgInfo) bool {
	if !pc.IsRunning() {
		return false
	}
	select {
	case pc.sendQueue <- msg:
		atomic.AddInt32(&pc.sendQueueSize, 1)
		return true
	default:
		return false
	}
}

// PickSendVote picks a vote and sends it to the peer.
// Returns true if vote was sent.
func (pc *peerConnV2) PickSendVote(votes ttypes.VoteSetReader) bool {
	if vote, ok := pc.state.PickVoteToSend(votes); ok {
		msg := MsgInfo{TypeID: ttypes.VoteID, Msg: vote.Vote, PeerID: pc.id, PeerIP: pc.ip.String()}
		tendermintlog.Debug("Sending vote message", "vote", msg)
		if pc.Send(msg) {
			pc.state.SetHasVote(vote)
			return true
		}
		return false
	}
	return false
}

func (pc *peerConnV2) IsRunning() bool {
	return atomic.LoadUint32(&pc.started) == 1 && atomic.LoadUint32(&pc.stopped) == 0
}

func (pc *peerConnV2) Start() error {
	if atomic.CompareAndSwapUint32(&pc.started, 0, 1) {
		if atomic.LoadUint32(&pc.stopped) == 1 {
			tendermintlog.Error("peerConn already stoped can not start", "peerIP", pc.ip.String())
			return nil
		}
		pc.pongChannel = make(chan struct{})
		pc.sendQueue = make(chan MsgInfo, maxSendQueueSize)
		pc.sendBuffer = make([]byte, 0, MaxMsgPacketPayloadSize)
		pc.quitSend = make(chan struct{})
		pc.quitUpdate = make(chan struct{})
		pc.quitBeat = make(chan struct{})
		pc.state = &PeerConnState{ip: pc.ip, PeerRoundState: ttypes.PeerRoundState{
			Round:              -1,
			ProposalPOLRound:   -1,
			LastCommitRound:    -1,
			CatchupCommitRound: -1,
		}}
		pc.updateStateQueue = make(chan MsgInfo, maxSendQueueSize)
		pc.heartbeatQueue = make(chan proto.Message, 100)
		pc.waitQuit.Add(5) //heartbeatRoutine, updateStateRoutine,gossipDataRoutine,gossipVotesRoutine,queryMaj23Routine

		go pc.sendRoutine()
		go pc.recvRoutine()
		go pc.updateStateRoutine()
		go pc.heartbeatRoutine()

		go pc.gossipDataRoutine()
		go pc.gossipVotesRoutine()
		go pc.queryMaj23Routine()

	}
	return nil
}

func (pc *peerConnV2) Stop() {
	if atomic.CompareAndSwapUint32(&pc.stopped, 0, 1) {
		pc.quitSend <- struct{}{}
		pc.quitUpdate <- struct{}{}
		pc.quitBeat <- struct{}{}

		pc.waitQuit.Wait()
		tendermintlog.Info("peerConn stop waitQuit", "peerIP", pc.ip.String())

		close(pc.sendQueue)
		close(pc.receiveConMsgChannel)
		pc.sendQueue = nil
		pc.transferChannel = nil
		tendermintlog.Info("peerConn stop finish", "peerIP", pc.ip.String())
	}
}

// Catch panics, usually caused by remote disconnects.
func (pc *peerConnV2) _recover() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		err := StackError{r, stack}
		pc.stopForError(err)
	}
}

func (pc *peerConnV2) stopForError(r interface{}) {
	tendermintlog.Error("peerConn recovered panic", "error", r, "peer", pc.ip.String())
	if pc.onPeerError != nil {
		pc.onPeerError(pc, r)
	} else {
		pc.Stop()
	}
}
func (pc *peerConnV2) sendRoutine() {
	defer pc._recover()
FOR_LOOP:
	for {
		select {
		case <-pc.quitSend:
			break FOR_LOOP
		case msg := <-pc.sendQueue:
			bytes, err := proto.Marshal(msg.Msg)
			if err != nil {
				tendermintlog.Error("peerConn sendroutine marshal data failed", "error", err)
				pc.stopForError(err)
				break FOR_LOOP
			}
			packet := &msgPacket{TypeID: msg.TypeID, Bytes: bytes}
			data, err := json.Marshal(packet)
			if err != nil {
				tendermintlog.Error("peerConn sendroutine marshal packet data failed", "error", err)
				pc.stopForError(err)
				break FOR_LOOP
			}
			//组装consensusMsg,设置目的地址
			conMsg := &types.ConsensusMsg{Data: data, MsgType: ttypes.ConsentType, ToPeerID: pc.peerID}
			pc.SendMsg(conMsg)

		case _, ok := <-pc.pongChannel:
			if ok {
				tendermintlog.Debug("Send Pong")
				packet := &msgPacket{TypeID: ttypes.PacketTypePong}
				data, err := json.Marshal(packet)
				if err != nil {
					tendermintlog.Error("peerConn sendroutine marshal packet data failed", "error", err)
					pc.stopForError(err)
					break FOR_LOOP
				}
				conMsg := &types.ConsensusMsg{Data: data, MsgType: ttypes.ConsentType}
				pc.SendMsg(conMsg)
			} else {
				pc.pongChannel = nil
			}
		}
	}
	tendermintlog.Info("peerConn stop sendRoutine", "peerIP", pc.ip.String())
}

func (pc *peerConnV2) recvRoutine() {
	defer pc._recover()
FOR_LOOP:
	for {
		select {
		case data := <-pc.receiveConMsgChannel:
			var pkt msgPacket
			err := json.Unmarshal(data.Data, &pkt)
			if err != nil {
				tendermintlog.Error("Connection failed @ recvRoutine (unmarsha packet)", "conn", pc, "err", err)
				pc.stopForError(err)
				break FOR_LOOP
			}
			if pkt.TypeID == ttypes.PacketTypePong {
				tendermintlog.Info("Receive Pong")
			} else if pkt.TypeID == ttypes.PacketTypePing {
				tendermintlog.Info("Receive Ping")
				pc.pongChannel <- struct{}{}
			} else {
				if v, ok := ttypes.MsgMap[pkt.TypeID]; ok {
					realMsg := reflect.New(v).Interface()
					err := proto.Unmarshal(pkt.Bytes, realMsg.(proto.Message))
					if err != nil {
						tendermintlog.Error("peerConn recvRoutine Unmarshal data failed", "err", err)
						continue
					}
					if pc.transferChannel != nil && (pkt.TypeID == ttypes.ProposalID || pkt.TypeID == ttypes.VoteID ||
						pkt.TypeID == ttypes.ProposalBlockID) {
						pc.transferChannel <- MsgInfo{pkt.TypeID, realMsg.(proto.Message), pc.ID(), pc.ip.String()}
						if pkt.TypeID == ttypes.ProposalID {
							proposal := realMsg.(*tmtypes.Proposal)
							tendermintlog.Debug("Receiving proposal", "proposal-height", proposal.Height, "peerip", pc.ip.String())
							pc.state.SetHasProposal(proposal)
						} else if pkt.TypeID == ttypes.VoteID {
							vote := &ttypes.Vote{Vote: realMsg.(*tmtypes.Vote)}
							tendermintlog.Debug("Receiving vote", "vote-height", vote.Height, "peerip", pc.ip.String())
							pc.state.SetHasVote(vote)
						} else if pkt.TypeID == ttypes.ProposalBlockID {
							block := &ttypes.TendermintBlock{TendermintBlock: realMsg.(*tmtypes.TendermintBlock)}
							tendermintlog.Debug("Receiving proposal block", "block-height", block.Header.Height, "peerip", pc.ip.String())
							pc.state.SetHasProposalBlock(block)
						}
					} else if pkt.TypeID == ttypes.ProposalHeartbeatID {
						pc.heartbeatQueue <- realMsg.(*tmtypes.Heartbeat)
					} else {
						pc.updateStateQueue <- MsgInfo{pkt.TypeID, realMsg.(proto.Message), pc.ID(), pc.ip.String()}
					}
				} else {
					err := fmt.Errorf("Unknown message type %v", pkt.TypeID)
					tendermintlog.Error("Connection failed @ recvRoutine", "conn", pc, "err", err)
					pc.stopForError(err)
					break FOR_LOOP
				}
			}
		}

	}

	close(pc.pongChannel)
	close(pc.heartbeatQueue)
	close(pc.updateStateQueue)
	tendermintlog.Info("peerConn stop recvRoutine", "peerIP", pc.ip.String())
}

func (pc *peerConnV2) updateStateRoutine() {
FOR_LOOP:
	for {
		select {
		case <-pc.quitUpdate:
			pc.waitQuit.Done()
			break FOR_LOOP
		case msg := <-pc.updateStateQueue:
			typeID := msg.TypeID
			if typeID == ttypes.NewRoundStepID {
				pc.state.ApplyNewRoundStepMessage(msg.Msg.(*tmtypes.NewRoundStepMsg))
			} else if typeID == ttypes.ValidBlockID {
				pc.state.ApplyValidBlockMessage(msg.Msg.(*tmtypes.ValidBlockMsg))
			} else if typeID == ttypes.HasVoteID {
				pc.state.ApplyHasVoteMessage(msg.Msg.(*tmtypes.HasVoteMsg))
			} else if typeID == ttypes.VoteSetMaj23ID {
				tmp := msg.Msg.(*tmtypes.VoteSetMaj23Msg)
				tendermintlog.Debug("updateStateRoutine", "VoteSetMaj23Msg", tmp)
				pc.myState.SetPeerMaj23(tmp.Height, int(tmp.Round), byte(tmp.Type), pc.id, tmp.BlockID)
				var myVotes *ttypes.BitArray
				switch byte(tmp.Type) {
				case ttypes.VoteTypePrevote:
					myVotes = pc.myState.GetPrevotesState(tmp.Height, int(tmp.Round), tmp.BlockID)
				case ttypes.VoteTypePrecommit:
					myVotes = pc.myState.GetPrecommitsState(tmp.Height, int(tmp.Round), tmp.BlockID)
				default:
					tendermintlog.Error("Bad VoteSetBitsMessage field Type", "type", byte(tmp.Type))
					return
				}
				if myVotes != nil && myVotes.TendermintBitArray != nil {
					voteSetBitMsg := &tmtypes.VoteSetBitsMsg{
						Height:  tmp.Height,
						Round:   tmp.Round,
						Type:    tmp.Type,
						BlockID: tmp.BlockID,
						Votes:   myVotes.TendermintBitArray,
					}
					pc.sendQueue <- MsgInfo{TypeID: ttypes.VoteSetBitsID, Msg: voteSetBitMsg, PeerID: pc.id, PeerIP: pc.ip.String()}
				}

			} else if typeID == ttypes.ProposalPOLID {
				pc.state.ApplyProposalPOLMessage(msg.Msg.(*tmtypes.ProposalPOLMsg))
			} else if typeID == ttypes.VoteSetBitsID {
				tmp := msg.Msg.(*tmtypes.VoteSetBitsMsg)
				if pc.myState.Height == tmp.Height {
					var myVotes *ttypes.BitArray
					switch byte(tmp.Type) {
					case ttypes.VoteTypePrevote:
						myVotes = pc.myState.GetPrevotesState(tmp.Height, int(tmp.Round), tmp.BlockID)
					case ttypes.VoteTypePrecommit:
						myVotes = pc.myState.GetPrecommitsState(tmp.Height, int(tmp.Round), tmp.BlockID)
					default:
						tendermintlog.Error("Bad VoteSetBitsMessage field Type", "type", byte(tmp.Type))
						return
					}
					pc.state.ApplyVoteSetBitsMessage(tmp, myVotes)
				} else {
					pc.state.ApplyVoteSetBitsMessage(tmp, nil)
				}
			} else {
				tendermintlog.Error("Unknown message type in updateStateRoutine", "msg", msg)
			}
		}
	}
	tendermintlog.Info("peerConn stop updateStateRoutine", "peerIP", pc.ip.String())
}

func (pc *peerConnV2) heartbeatRoutine() {
FOR_LOOP:
	for {
		select {
		case <-pc.quitBeat:
			pc.waitQuit.Done()
			break FOR_LOOP
		case heartbeat := <-pc.heartbeatQueue:
			msg := heartbeat.(*tmtypes.Heartbeat)
			tendermintlog.Debug("Received proposal heartbeat message",
				"height", msg.Height, "round", msg.Round, "sequence", msg.Sequence,
				"valIdx", msg.ValidatorIndex, "valAddr", msg.ValidatorAddress)
		}
	}
	tendermintlog.Info("peerConn stop heartbeatRoutine", "peerIP", pc.ip.String())
}

func (pc *peerConnV2) gossipDataRoutine() {
OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !pc.IsRunning() {
			pc.waitQuit.Done()
			tendermintlog.Info("peerConn stop gossipDataRoutine", "peerIP", pc.ip.String())
			return
		}

		rs := pc.myState.GetRoundState()
		prs := pc.state

		// If the peer is on a previous height, help catch up.
		if (0 < prs.Height) && (prs.Height < rs.Height) {
			if prs.ProposalBlockHash == nil || prs.ProposalBlock {
				time.Sleep(pc.myState.PeerGossipSleep())
				continue OUTER_LOOP
			}
			tendermintlog.Info("help catch up", "peerip", pc.ip.String(), "selfHeight", rs.Height, "peerHeight", prs.Height)
			proposalBlock := pc.myState.client.LoadProposalBlock(prs.Height)
			newBlock := &ttypes.TendermintBlock{TendermintBlock: proposalBlock}
			if proposalBlock == nil {
				tendermintlog.Error("Fail to load propsal block", "selfHeight", rs.Height,
					"blockstoreHeight", pc.myState.client.GetCurrentHeight())
				time.Sleep(pc.myState.PeerGossipSleep())
				continue OUTER_LOOP
			} else if !bytes.Equal(newBlock.Hash(), prs.ProposalBlockHash) {
				tendermintlog.Error("Peer ProposalBlockHash mismatch", "ProposalBlockHash", fmt.Sprintf("%X", prs.ProposalBlockHash),
					"newBlockHash", fmt.Sprintf("%X", newBlock.Hash()))
				time.Sleep(pc.myState.PeerGossipSleep())
				continue OUTER_LOOP
			}
			msg := MsgInfo{TypeID: ttypes.ProposalBlockID, Msg: proposalBlock, PeerID: pc.id, PeerIP: pc.ip.String()}
			tendermintlog.Info("Sending block for catchup", "peerip", pc.ip.String(), "block(H/R)",
				fmt.Sprintf("%v/%v", proposalBlock.Header.Height, proposalBlock.Header.Round))
			if pc.Send(msg) {
				prs.SetHasProposalBlock(newBlock)
			}
			continue OUTER_LOOP
		}

		// If height and round don't match, sleep.
		if (rs.Height != prs.Height) || (rs.Round != prs.Round) {
			time.Sleep(pc.myState.PeerGossipSleep())
			continue OUTER_LOOP
		}

		// By here, height and round match.
		// Proposal block parts were already matched and sent if any were wanted.
		// (These can match on hash so the round doesn't matter)
		// Now consider sending other things, like the Proposal itself.

		// Send Proposal && ProposalPOL BitArray?
		if rs.Proposal != nil && !prs.Proposal {
			// Proposal: share the proposal metadata with peer.
			{
				msg := MsgInfo{TypeID: ttypes.ProposalID, Msg: rs.Proposal, PeerID: pc.id, PeerIP: pc.ip.String()}
				tendermintlog.Debug(fmt.Sprintf("Sending proposal. Self state: %v/%v/%v", rs.Height, rs.Round, rs.Step),
					"peerip", pc.ip.String(), "proposal-height", rs.Proposal.Height, "proposal-round", rs.Proposal.Round)
				if pc.Send(msg) {
					prs.SetHasProposal(rs.Proposal)
				}
			}
			// ProposalPOL: lets peer know which POL votes we have so far.
			// Peer must receive ttypes.ProposalMessage first.
			// rs.Proposal was validated, so rs.Proposal.POLRound <= rs.Round,
			// so we definitely have rs.Votes.Prevotes(rs.Proposal.POLRound).
			if 0 <= rs.Proposal.POLRound {
				msg := MsgInfo{TypeID: ttypes.ProposalPOLID, Msg: &tmtypes.ProposalPOLMsg{
					Height:           rs.Height,
					ProposalPOLRound: rs.Proposal.POLRound,
					ProposalPOL:      rs.Votes.Prevotes(int(rs.Proposal.POLRound)).BitArray().TendermintBitArray,
				}, PeerID: pc.id, PeerIP: pc.ip.String()}
				tendermintlog.Debug("Sending POL", "height", prs.Height, "round", prs.Round)
				pc.Send(msg)
			}
			continue OUTER_LOOP
		}

		// Send proposal block
		if rs.Proposal != nil && prs.ProposalBlockHash != nil && bytes.Equal(rs.Proposal.Blockhash, prs.ProposalBlockHash) {
			if rs.ProposalBlock != nil && !prs.ProposalBlock {
				msg := MsgInfo{TypeID: ttypes.ProposalBlockID, Msg: rs.ProposalBlock.TendermintBlock, PeerID: pc.id, PeerIP: pc.ip.String()}
				tendermintlog.Debug(fmt.Sprintf("Sending proposal block. Self state: %v/%v/%v", rs.Height, rs.Round, rs.Step),
					"peerip", pc.ip.String(), "block-height", rs.ProposalBlock.Header.Height, "block-round", rs.ProposalBlock.Header.Round)
				if pc.Send(msg) {
					prs.SetHasProposalBlock(rs.ProposalBlock)
				}
				continue OUTER_LOOP
			}
		}

		// Nothing to do. Sleep.
		time.Sleep(pc.myState.PeerGossipSleep())
		continue OUTER_LOOP
	}
}

func (pc *peerConnV2) gossipVotesRoutine() {
	// Simple hack to throttle logs upon sleep.
	var sleeping = 0

OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !pc.IsRunning() {
			pc.waitQuit.Done()
			tendermintlog.Info("peerConn stop gossipVotesRoutine", "peerIP", pc.ip.String())
			return
		}

		rs := pc.myState.GetRoundState()
		prs := pc.state

		switch sleeping {
		case 1: // First sleep
			sleeping = 2
		case 2: // No more sleep
			sleeping = 0
		}

		// If height matches, then send LastCommit, Prevotes, Precommits.
		if rs.Height == prs.Height {
			if pc.gossipVotesForHeight(rs, &prs.PeerRoundState) {
				continue OUTER_LOOP
			}
		}

		// Special catchup logic.
		// If peer is lagging by height 1, send LastCommit.
		if prs.Height != 0 && rs.Height == prs.Height+1 {
			if pc.PickSendVote(rs.LastCommit) {
				tendermintlog.Debug("Picked rs.LastCommit to send", "peerip", pc.ip.String(), "height", prs.Height)
				continue OUTER_LOOP
			}
		}

		// Catchup logic
		// If peer is lagging by more than 1, send Commit.
		if prs.Height != 0 && rs.Height >= prs.Height+2 {
			// Load the block commit for prs.Height,
			// which contains precommit signatures for prs.Height.
			commit := pc.myState.client.LoadBlockCommit(prs.Height + 1)
			commitObj := &ttypes.Commit{TendermintCommit: commit}
			if pc.PickSendVote(commitObj) {
				tendermintlog.Info("Picked Catchup commit to send",
					"commit(H/R)", fmt.Sprintf("%v/%v", commitObj.Height(), commitObj.Round()),
					"BitArray", commitObj.BitArray().String(),
					"peerip", pc.ip.String(), "height", prs.Height)
				continue OUTER_LOOP
			}
		}

		if sleeping == 0 {
			// We sent nothing. Sleep...
			sleeping = 1
			tendermintlog.Debug("No votes to send, sleeping", "peerip", pc.ip.String(), "rs.Height", rs.Height, "prs.Height", prs.Height,
				"localPV", rs.Votes.Prevotes(rs.Round).BitArray(), "peerPV", prs.Prevotes,
				"localPC", rs.Votes.Precommits(rs.Round).BitArray(), "peerPC", prs.Precommits)
		} else if sleeping == 2 {
			// Continued sleep...
			sleeping = 1
		}

		time.Sleep(pc.myState.PeerGossipSleep())
		continue OUTER_LOOP
	}
}

func (pc *peerConnV2) gossipVotesForHeight(rs *ttypes.RoundState, prs *ttypes.PeerRoundState) bool {
	// If there are lastCommits to send...
	if prs.Step == ttypes.RoundStepNewHeight {
		if pc.PickSendVote(rs.LastCommit) {
			tendermintlog.Debug("Picked rs.LastCommit to send", "peerip", pc.ip.String(),
				"peer(H/R)", fmt.Sprintf("%v/%v", prs.Height, prs.Round))
			return true
		}
	}
	// If there are POL prevotes to send...
	if prs.Step <= ttypes.RoundStepPropose && prs.Round != -1 && prs.Round <= rs.Round && prs.ProposalPOLRound != -1 {
		if polPrevotes := rs.Votes.Prevotes(prs.ProposalPOLRound); polPrevotes != nil {
			if pc.PickSendVote(polPrevotes) {
				tendermintlog.Debug("Picked rs.Prevotes(prs.ProposalPOLRound) to send",
					"peerip", pc.ip.String(), "peer(H/R)", fmt.Sprintf("%v/%v", prs.Height, prs.Round),
					"POLRound", prs.ProposalPOLRound)
				return true
			}
		}
	}
	// If there are prevotes to send...
	if prs.Step <= ttypes.RoundStepPrevoteWait && prs.Round != -1 && prs.Round <= rs.Round {
		if pc.PickSendVote(rs.Votes.Prevotes(prs.Round)) {
			tendermintlog.Debug("Picked rs.Prevotes(prs.Round) to send",
				"peerip", pc.ip.String(), "peer(H/R)", fmt.Sprintf("%v/%v", prs.Height, prs.Round))
			return true
		}
	}
	// If there are precommits to send...
	if prs.Step <= ttypes.RoundStepPrecommitWait && prs.Round != -1 && prs.Round <= rs.Round {
		if pc.PickSendVote(rs.Votes.Precommits(prs.Round)) {
			tendermintlog.Debug("Picked rs.Precommits(prs.Round) to send",
				"peerip", pc.ip.String(), "peer(H/R)", fmt.Sprintf("%v/%v", prs.Height, prs.Round))
			return true
		}
	}
	// If there are prevotes to send...Needed because of validBlock mechanism
	if prs.Round != -1 && prs.Round <= rs.Round {
		if pc.PickSendVote(rs.Votes.Prevotes(prs.Round)) {
			tendermintlog.Debug("Picked rs.Prevotes(prs.Round) to send",
				"peerip", pc.ip.String(), "peer(H/R)", fmt.Sprintf("%v/%v", prs.Height, prs.Round))
			return true
		}
	}
	// If there are POLPrevotes to send...
	if prs.ProposalPOLRound != -1 {
		if polPrevotes := rs.Votes.Prevotes(prs.ProposalPOLRound); polPrevotes != nil {
			if pc.PickSendVote(polPrevotes) {
				tendermintlog.Debug("Picked rs.Prevotes(prs.ProposalPOLRound) to send",
					"peerip", pc.ip.String(), "round", prs.ProposalPOLRound)
				return true
			}
		}
	}
	return false
}

func (pc *peerConnV2) queryMaj23Routine() {
OUTER_LOOP:
	for {
		// Manage disconnects from self or peer.
		if !pc.IsRunning() {
			pc.waitQuit.Done()
			tendermintlog.Info("peerConn stop queryMaj23Routine", "peerIP", pc.ip.String())
			return
		}

		// Maybe send Height/Round/Prevotes
		{
			rs := pc.myState.GetRoundState()
			prs := pc.state
			if rs.Height == prs.Height {
				if maj23, ok := rs.Votes.Prevotes(prs.Round).TwoThirdsMajority(); ok {
					msg := MsgInfo{TypeID: ttypes.VoteSetMaj23ID, Msg: &tmtypes.VoteSetMaj23Msg{
						Height:  prs.Height,
						Round:   int32(prs.Round),
						Type:    int32(ttypes.VoteTypePrevote),
						BlockID: &maj23,
					}, PeerID: pc.id, PeerIP: pc.ip.String(),
					}
					pc.TrySend(msg)
					time.Sleep(pc.myState.PeerQueryMaj23Sleep())
				}
			}
		}

		// Maybe send Height/Round/Precommits
		{
			rs := pc.myState.GetRoundState()
			prs := pc.state.GetRoundState()
			if rs.Height == prs.Height {
				if maj23, ok := rs.Votes.Precommits(prs.Round).TwoThirdsMajority(); ok {
					msg := MsgInfo{TypeID: ttypes.VoteSetMaj23ID, Msg: &tmtypes.VoteSetMaj23Msg{
						Height:  prs.Height,
						Round:   int32(prs.Round),
						Type:    int32(ttypes.VoteTypePrecommit),
						BlockID: &maj23,
					}, PeerID: pc.id, PeerIP: pc.ip.String(),
					}
					pc.TrySend(msg)
					time.Sleep(pc.myState.PeerQueryMaj23Sleep())
				}
			}
		}

		// Maybe send Height/Round/ProposalPOL
		{
			rs := pc.myState.GetRoundState()
			prs := pc.state.GetRoundState()
			if rs.Height == prs.Height && prs.ProposalPOLRound >= 0 {
				if maj23, ok := rs.Votes.Prevotes(prs.ProposalPOLRound).TwoThirdsMajority(); ok {
					msg := MsgInfo{TypeID: ttypes.VoteSetMaj23ID, Msg: &tmtypes.VoteSetMaj23Msg{
						Height:  prs.Height,
						Round:   int32(prs.ProposalPOLRound),
						Type:    int32(ttypes.VoteTypePrevote),
						BlockID: &maj23,
					}, PeerID: pc.id, PeerIP: pc.ip.String(),
					}
					pc.TrySend(msg)
					time.Sleep(pc.myState.PeerQueryMaj23Sleep())
				}
			}
		}

		// Little point sending LastCommitRound/LastCommit,
		// These are fleeting and non-blocking.

		// Maybe send Height/CatchupCommitRound/CatchupCommit.
		{
			prs := pc.state.GetRoundState()
			if prs.CatchupCommitRound != -1 && 0 < prs.Height && prs.Height <= pc.myState.client.csStore.LoadStateHeight() {
				commit := pc.myState.LoadCommit(prs.Height)
				commitTmp := ttypes.Commit{TendermintCommit: commit}
				msg := MsgInfo{TypeID: ttypes.VoteSetMaj23ID, Msg: &tmtypes.VoteSetMaj23Msg{
					Height:  prs.Height,
					Round:   int32(commitTmp.Round()),
					Type:    int32(ttypes.VoteTypePrecommit),
					BlockID: commit.BlockID,
				}, PeerID: pc.id, PeerIP: pc.ip.String(),
				}
				pc.TrySend(msg)
				time.Sleep(pc.myState.PeerQueryMaj23Sleep())
			}
		}

		time.Sleep(pc.myState.PeerQueryMaj23Sleep())

		continue OUTER_LOOP
	}
}

// GetRoundState returns an atomic snapshot of the PeerRoundState.
// There's no point in mutating it since it won't change PeerState.
func (ps *PeerConnStateV2) GetRoundState() *ttypes.PeerRoundState {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	prs := ps.PeerRoundState // copy
	return &prs
}

// GetHeight returns an atomic snapshot of the PeerRoundState's height
// used by the mempool to ensure peers are caught up before broadcasting new txs
func (ps *PeerConnStateV2) GetHeight() int64 {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return ps.PeerRoundState.Height
}

// SetHasProposal sets the given proposal as known for the peer.
func (ps *PeerConnStateV2) SetHasProposal(proposal *tmtypes.Proposal) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.Height != proposal.Height || ps.Round != int(proposal.Round) {
		return
	}
	if ps.Proposal {
		return
	}
	tendermintlog.Debug("Peer set proposal", "peerip", ps.ip.String(),
		"peer-state", fmt.Sprintf("%v/%v/%v", ps.Height, ps.Round, ps.Step),
		"proposal(H/R/Hash)", fmt.Sprintf("%v/%v/%X", proposal.Height, proposal.Round, proposal.Blockhash))
	ps.Proposal = true

	ps.ProposalBlockHash = proposal.Blockhash
	ps.ProposalPOLRound = int(proposal.POLRound)
	ps.ProposalPOL = nil // Nil until ttypes.ProposalPOLMessage received.
}

// SetHasProposalBlock sets the given proposal block as known for the peer.
func (ps *PeerConnStateV2) SetHasProposalBlock(block *ttypes.TendermintBlock) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.Height != block.Header.Height || ps.Round != int(block.Header.Round) {
		return
	}
	if ps.ProposalBlock {
		return
	}
	tendermintlog.Debug("Peer set proposal block", "peerip", ps.ip.String(),
		"peer-state", fmt.Sprintf("%v/%v/%v", ps.Height, ps.Round, ps.Step),
		"block(H/R)", fmt.Sprintf("%v/%v", block.Header.Height, block.Header.Round))
	ps.ProposalBlock = true
}

// PickVoteToSend picks a vote to send to the peer.
// Returns true if a vote was picked.
// NOTE: `votes` must be the correct Size() for the Height().
func (ps *PeerConnStateV2) PickVoteToSend(votes ttypes.VoteSetReader) (vote *ttypes.Vote, ok bool) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if votes.Size() == 0 {
		return nil, false
	}

	height, round, voteType, size := votes.Height(), votes.Round(), votes.Type(), votes.Size()

	// Lazily set data using 'votes'.
	if votes.IsCommit() {
		ps.ensureCatchupCommitRound(height, round, size)
	}
	ps.ensureVoteBitArrays(height, size)

	psVotes := ps.getVoteBitArray(height, round, voteType)
	if psVotes == nil {
		return nil, false // Not something worth sending
	}

	if index, ok := votes.BitArray().Sub(psVotes).PickRandom(); ok {
		tendermintlog.Debug("PickVoteToSend", "peer(H/R)", fmt.Sprintf("%v/%v", ps.Height, ps.Round),
			"vote(H/R)", fmt.Sprintf("%v/%v", height, round), "type", voteType, "selfVotes", votes.BitArray().String(),
			"peerVotes", psVotes.String(), "peerip", ps.ip.String())
		return votes.GetByIndex(index), true
	}
	return nil, false
}

func (ps *PeerConnStateV2) getVoteBitArray(height int64, round int, voteType byte) *ttypes.BitArray {
	if !ttypes.IsVoteTypeValid(voteType) {
		return nil
	}

	if ps.Height == height {
		if ps.Round == round {
			switch voteType {
			case ttypes.VoteTypePrevote:
				return ps.Prevotes
			case ttypes.VoteTypePrecommit:
				return ps.Precommits
			}
		}
		if ps.CatchupCommitRound == round {
			switch voteType {
			case ttypes.VoteTypePrevote:
				return nil
			case ttypes.VoteTypePrecommit:
				return ps.CatchupCommit
			}
		}
		if ps.ProposalPOLRound == round {
			switch voteType {
			case ttypes.VoteTypePrevote:
				return ps.ProposalPOL
			case ttypes.VoteTypePrecommit:
				return nil
			}
		}
		return nil
	}
	if ps.Height == height+1 {
		if ps.LastCommitRound == round {
			switch voteType {
			case ttypes.VoteTypePrevote:
				return nil
			case ttypes.VoteTypePrecommit:
				return ps.LastCommit
			}
		}
		return nil
	}
	return nil
}

// 'round': A round for which we have a +2/3 commit.
func (ps *PeerConnStateV2) ensureCatchupCommitRound(height int64, round int, numValidators int) {
	if ps.Height != height {
		return
	}
	/*
		NOTE: This is wrong, 'round' could change.
		e.g. if orig round is not the same as block LastCommit round.
		if ps.CatchupCommitRound != -1 && ps.CatchupCommitRound != round {
			ttypes.PanicSanity(ttypes.Fmt("Conflicting CatchupCommitRound. Height: %v, Orig: %v, New: %v", height, ps.CatchupCommitRound, round))
		}
	*/
	if ps.CatchupCommitRound == round {
		return // Nothing to do!
	}
	tendermintlog.Debug("ensureCatchupCommitRound", "height", height, "round", round, "ps.CatchupCommitRound", ps.CatchupCommitRound,
		"ps.Round", ps.Round, "peerip", ps.ip.String())
	ps.CatchupCommitRound = round
	if round == ps.Round {
		ps.CatchupCommit = ps.Precommits
	} else {
		ps.CatchupCommit = ttypes.NewBitArray(numValidators)
	}
}

// EnsureVoteBitArrays ensures the bit-arrays have been allocated for tracking
// what votes this peer has received.
// NOTE: It's important to make sure that numValidators actually matches
// what the node sees as the number of validators for height.
func (ps *PeerConnStateV2) EnsureVoteBitArrays(height int64, numValidators int) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	ps.ensureVoteBitArrays(height, numValidators)
}

func (ps *PeerConnStateV2) ensureVoteBitArrays(height int64, numValidators int) {
	if ps.Height == height {
		if ps.Prevotes == nil {
			ps.Prevotes = ttypes.NewBitArray(numValidators)
		}
		if ps.Precommits == nil {
			ps.Precommits = ttypes.NewBitArray(numValidators)
		}
		if ps.CatchupCommit == nil {
			ps.CatchupCommit = ttypes.NewBitArray(numValidators)
		}
		if ps.ProposalPOL == nil {
			ps.ProposalPOL = ttypes.NewBitArray(numValidators)
		}
	} else if ps.Height == height+1 {
		if ps.LastCommit == nil {
			ps.LastCommit = ttypes.NewBitArray(numValidators)
		}
	}
}

// SetHasVote sets the given vote as known by the peer
func (ps *PeerConnStateV2) SetHasVote(vote *ttypes.Vote) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.setHasVote(vote.Height, int(vote.Round), byte(vote.Type), int(vote.ValidatorIndex))
}

func (ps *PeerConnStateV2) setHasVote(height int64, round int, voteType byte, index int) {
	// NOTE: some may be nil BitArrays -> no side effects.
	psVotes := ps.getVoteBitArray(height, round, voteType)
	tendermintlog.Debug("setHasVote before", "height", height, "psVotes", psVotes.String(), "peerip", ps.ip.String())
	if psVotes != nil {
		psVotes.SetIndex(index, true)
	}
	tendermintlog.Debug("setHasVote after", "height", height, "index", index, "type", voteType,
		"peerVotes", psVotes.String(), "peerip", ps.ip.String())
}

// ApplyNewRoundStepMessage updates the peer state for the new round.
func (ps *PeerConnStateV2) ApplyNewRoundStepMessage(msg *tmtypes.NewRoundStepMsg) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	// Ignore duplicates or decreases
	if CompareHRS(msg.Height, int(msg.Round), ttypes.RoundStepType(msg.Step), ps.Height, ps.Round, ps.Step) <= 0 {
		return
	}

	// Just remember these values.
	psHeight := ps.Height
	psRound := ps.Round
	//psStep := ps.Step
	psCatchupCommitRound := ps.CatchupCommitRound
	psCatchupCommit := ps.CatchupCommit
	psPrecommits := ps.Precommits

	startTime := time.Now().Add(-1 * time.Duration(msg.SecondsSinceStartTime) * time.Second)
	ps.Height = msg.Height
	ps.Round = int(msg.Round)
	ps.Step = ttypes.RoundStepType(msg.Step)
	ps.StartTime = startTime

	tendermintlog.Debug("ApplyNewRoundStepMessage", "peerip", ps.ip.String(),
		"peer(H/R)", fmt.Sprintf("%v/%v", psHeight, psRound),
		"msg(H/R/S)", fmt.Sprintf("%v/%v/%v", msg.Height, msg.Round, ps.Step))

	if psHeight != msg.Height || psRound != int(msg.Round) {
		tendermintlog.Debug("Reset Proposal, Prevotes, Precommits", "peerip", ps.ip.String(),
			"peer(H/R)", fmt.Sprintf("%v/%v", psHeight, psRound))
		ps.Proposal = false
		ps.ProposalBlock = false
		ps.ProposalBlockHash = nil
		ps.ProposalPOLRound = -1
		ps.ProposalPOL = nil
		// We'll update the BitArray capacity later.
		ps.Prevotes = nil
		ps.Precommits = nil
	}
	if psHeight == msg.Height && psRound != int(msg.Round) && int(msg.Round) == psCatchupCommitRound {
		// Peer caught up to CatchupCommitRound.
		// Preserve psCatchupCommit!
		// NOTE: We prefer to use prs.Precommits if
		// pr.Round matches pr.CatchupCommitRound.
		tendermintlog.Debug("Reset Precommits to CatchupCommit", "peerip", ps.ip.String(),
			"peer(H/R)", fmt.Sprintf("%v/%v", psHeight, psRound))
		ps.Precommits = psCatchupCommit
	}
	if psHeight != msg.Height {
		tendermintlog.Debug("Reset LastCommit, CatchupCommit", "peerip", ps.ip.String(),
			"peer(H/R)", fmt.Sprintf("%v/%v", psHeight, psRound))
		// Shift Precommits to LastCommit.
		if psHeight+1 == msg.Height && psRound == int(msg.LastCommitRound) {
			ps.LastCommitRound = int(msg.LastCommitRound)
			ps.LastCommit = psPrecommits
		} else {
			ps.LastCommitRound = int(msg.LastCommitRound)
			ps.LastCommit = nil
		}
		// We'll update the BitArray capacity later.
		ps.CatchupCommitRound = -1
		ps.CatchupCommit = nil
	}
}

// ApplyValidBlockMessage updates the peer state for the new valid block.
func (ps *PeerConnStateV2) ApplyValidBlockMessage(msg *tmtypes.ValidBlockMsg) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.Height != msg.Height {
		return
	}
	if ps.Round != int(msg.Round) && !msg.IsCommit {
		return
	}
	tendermintlog.Debug("ApplyValidBlockMessage", "peerip", ps.ip.String(),
		"peer(H/R)", fmt.Sprintf("%v/%v", ps.Height, ps.Round),
		"blockhash", fmt.Sprintf("%X", msg.Blockhash))

	ps.ProposalBlockHash = msg.Blockhash
}

// ApplyProposalPOLMessage updates the peer state for the new proposal POL.
func (ps *PeerConnStateV2) ApplyProposalPOLMessage(msg *tmtypes.ProposalPOLMsg) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.Height != msg.Height {
		return
	}
	if ps.ProposalPOLRound != int(msg.ProposalPOLRound) {
		return
	}

	// TODO: Merge onto existing ps.ProposalPOL?
	// We might have sent some prevotes in the meantime.
	ps.ProposalPOL = &ttypes.BitArray{TendermintBitArray: msg.ProposalPOL}
}

// ApplyHasVoteMessage updates the peer state for the new vote.
func (ps *PeerConnStateV2) ApplyHasVoteMessage(msg *tmtypes.HasVoteMsg) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.Height != msg.Height {
		return
	}

	tendermintlog.Debug("ApplyHasVoteMessage", "msg(H/R)", fmt.Sprintf("%v/%v", msg.Height, msg.Round),
		"peerip", ps.ip.String())
	ps.setHasVote(msg.Height, int(msg.Round), byte(msg.Type), int(msg.Index))
}

// ApplyVoteSetBitsMessage updates the peer state for the bit-array of votes
// it claims to have for the corresponding BlockID.
// `ourVotes` is a BitArray of votes we have for msg.BlockID
// NOTE: if ourVotes is nil (e.g. msg.Height < rs.Height),
// we conservatively overwrite ps's votes w/ msg.Votes.
func (ps *PeerConnStateV2) ApplyVoteSetBitsMessage(msg *tmtypes.VoteSetBitsMsg, ourVotes *ttypes.BitArray) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	votes := ps.getVoteBitArray(msg.Height, int(msg.Round), byte(msg.Type))
	if votes != nil {
		if ourVotes == nil {
			bitarray := &ttypes.BitArray{TendermintBitArray: msg.Votes}
			votes.Update(bitarray)
		} else {
			otherVotes := votes.Sub(ourVotes)
			bitarray := &ttypes.BitArray{TendermintBitArray: msg.Votes}
			hasVotes := otherVotes.Or(bitarray)
			votes.Update(hasVotes)
		}
	}
}
