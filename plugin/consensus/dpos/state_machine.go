// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dpos

import (
	"bytes"
	"encoding/json"
	"github.com/outbrain/golib/math"
	"time"

	"github.com/33cn/chain33/common/crypto"
	dpostype "github.com/33cn/plugin/plugin/consensus/dpos/types"
	"github.com/33cn/chain33/common"
)
var (
	InitStateType = 1
	VotingStateType = 2
	VotedStateType = 3
	WaitNotifyStateType = 4

	StateTypeMapping = map[int] string { InitStateType: "InitState", VotingStateType: "VotingState", VotedStateType: "VotedState", WaitNotifyStateType: "WaitNotifyState"}
)


type DPosTask struct {
	nodeId int64
	cycleStart int64
	cycleStop int64
	periodStart int64
	periodStop int64
	blockStart int64
	blockStop int64
}

func DecideTaskByTime(now int64) (task DPosTask){
	task.nodeId = now % dposCycle / dposPeriod

	task.cycleStart = now - now % dposCycle
	task.cycleStop = task.cycleStart + dposCycle - 1

	task.periodStart = task.cycleStart + task.nodeId * dposBlockInterval * dposContinueBlockNum
	task.periodStop = task.periodStart + dposPeriod - 1

	task.blockStart = task.periodStart + now % dposCycle % dposPeriod / dposBlockInterval * dposBlockInterval
	task.blockStop = task.blockStart + dposBlockInterval - 1

	return task
}


// DposState is the base class of dpos state machine, it defines some interfaces.
type DposState interface {
	timeOut(cs *ConsensusState)
	sendVote(cs *ConsensusState, vote * dpostype.DPosVote)
	recvVote(cs *ConsensusState, vote * dpostype.DPosVote)
	sendVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply)
	recvVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply)
	sendNotify(cs *ConsensusState, notify * dpostype.DPosNotify)
	recvNotify(cs *ConsensusState, notify * dpostype.DPosNotify)
}


// InitState is the initial state of dpos state machine
type InitState struct {

}

func (init *InitState)timeOut(cs *ConsensusState) {
	//if available noes  < 2/3, don't change the state to voting.
	connections := cs.client.node.peerSet.Size()
	validators := cs.validatorMgr.Validators.Size()
	if connections == 0 || connections < (validators * 2 / 3 -1) {
		dposlog.Error("InitState timeout but available nodes less than 2/3,waiting for more connections", "connections", connections, "validators", validators)
		cs.ClearVotes()

		//设定超时时间，超时后再检查链接数量
		cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
	} else {
		//获得当前高度
		height := cs.client.GetCurrentHeight()
		now := time.Now().Unix()
		if cs.lastMyVote != nil && math.AbsInt64(now - cs.lastMyVote.VoteItem.PeriodStop) <= 1 {
			now += 2
		}
		//计算当前时间，属于哪一个周期，应该哪一个节点出块，应该出块的高度
		task := DecideTaskByTime(now)

		addr, validator := cs.validatorMgr.Validators.GetByIndex(int(task.nodeId))
		if addr == nil && validator == nil {
			dposlog.Error("Address and Validator is nil", "node index", task.nodeId, "now", now, "cycle", dposCycle, "period", dposPeriod)
			//cs.SetState(InitStateObj)
			cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
			return
		}

		//生成vote， 对于vote进行签名
		voteItem := &dpostype.VoteItem{
			VotedNodeAddress: addr,
			VotedNodeIndex: int32(task.nodeId),
			CycleStart: task.cycleStart,
			CycleStop: task.cycleStop,
			PeriodStart: task.periodStart,
			PeriodStop: task.periodStop,
			Height: height + 1,
		}

		encode, err := json.Marshal(voteItem)
		if err != nil {
			panic("Marshal vote failed.")
			cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
			return
		}

		voteItem.VoteId = crypto.Ripemd160(encode)

		index := -1
		for i := 0; i < cs.validatorMgr.Validators.Size(); i ++ {
			if bytes.Compare(cs.validatorMgr.Validators.Validators[i].Address, cs.privValidator.GetAddress()) == 0 {
				index = i
				break
			}
		}

		if index == -1 {
			panic("This node's address is not exist in Validators.")
		}

		vote := &dpostype.Vote{DPosVote: &dpostype.DPosVote{
				VoteItem:         voteItem,
				VoteTimestamp:    now,
				VoterNodeAddress: cs.privValidator.GetAddress(),
				VoterNodeIndex:   int32(index),
			},
		}

		if err := cs.privValidator.SignVote(cs.validatorMgr.ChainID, vote); err != nil {
			dposlog.Error("SignVote failed", "vote", vote.String())
			//cs.SetState(InitStateObj)
			cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
			return
		}

		vote2 := *vote.DPosVote

		cs.AddVotes(&vote2)
		cs.SetMyVote(&vote2)
		dposlog.Info("Available nodes equal or more than 2/3,change state to VotingState", "connections", connections, "validators", validators)
		cs.SetState(VotingStateObj)
		dposlog.Info("Change state.", "from", "InitState", "to", "VotingState")
		//通过node发送p2p消息到其他节点
		cs.dposState.sendVote(cs, vote.DPosVote)

		cs.scheduleDPosTimeout(time.Duration(timeoutVoting) * time.Millisecond, VotingStateType)
		//处理之前缓存的投票信息
		for i := 0; i < len(cs.cachedVotes); i++ {
			cs.dposState.recvVote(cs, cs.cachedVotes[i])
		}
		cs.ClearCachedVotes()
	}
}

func (init *InitState)sendVote(cs *ConsensusState, vote * dpostype.DPosVote) {
	dposlog.Info("InitState does not support sendVote,so do nothing")
}

func (init *InitState)recvVote(cs *ConsensusState, vote * dpostype.DPosVote) {
	dposlog.Info("InitState recvVote ,add it and will handle it later.")
	//cs.AddVotes(vote)
	cs.CacheVotes(vote)
}

func (init *InitState)sendVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply){
	dposlog.Info("InitState don't support sendVoteReply,so do nothing")
}

func (init *InitState)recvVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply){
	dposlog.Info("InitState recv Vote reply,ignore it.")
}

func (init *InitState)sendNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("InitState does not support sendNotify,so do nothing")
}

func (init *InitState)recvNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("InitState recvNotify")

	//zzh:需要增加对Notify的处理，可以考虑记录已经确认过的出快记录
	cs.SetNotify(notify)
}

// VotingState is the voting state of dpos state machine until timeout or get an agreement by votes.
type VotingState struct {

}

func (voting *VotingState)timeOut(cs *ConsensusState) {
	dposlog.Info("VotingState timeout but don't get an agreement. change state to InitState")

	//清理掉之前的选票记录，从初始状态重新开始
	cs.ClearVotes()
	cs.ClearCachedVotes()
	cs.SetState(InitStateObj)
	dposlog.Info("Change state because of timeOut.", "from", "VotingState", "to", "InitState")

	//由于连接多数情况下正常，快速触发InitState的超时处理
	cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
}

func (voting *VotingState)sendVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	cs.broadcastChannel <- MsgInfo{TypeID: dpostype.VoteID, Msg: vote, PeerID: cs.ourID, PeerIP: ""}
}

func (voting *VotingState)recvVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	dposlog.Info("VotingState get a vote", "VotedNodeIndex", vote.VoteItem.VotedNodeIndex,
		                                            "VotedNodeAddress", common.ToHex(vote.VoteItem.VotedNodeAddress),
		                                            "CycleStart", vote.VoteItem.CycleStart,
		                                            "CycleStop", vote.VoteItem.CycleStop,
		                                            "PeriodStart", vote.VoteItem.PeriodStart,
		                                            "periodStop", vote.VoteItem.PeriodStop,
		                                            "Height", vote.VoteItem.Height,
		                                            "VoteId", common.ToHex(vote.VoteItem.VoteId),
		                                            "VoteTimestamp", vote.VoteTimestamp,
		                                            "VoterNodeIndex", vote.VoterNodeIndex,
		                                            "VoterNodeAddress", common.ToHex(vote.VoterNodeAddress),
		                                            "Signature", common.ToHex(vote.Signature),
		                                            "localNodeIndex", cs.client.ValidatorIndex(),"now", time.Now().Unix())

	if !cs.VerifyVote(vote) {
		dposlog.Info("VotingState verify vote failed")
		return
	}

	cs.AddVotes(vote)

	result, voteItem := cs.CheckVotes()

	if result == voteSuccess {
		dposlog.Info("VotingState get 2/3 result", "final vote:", voteItem.String())
		dposlog.Info("VotingState change state to VotedState")
		//切换状态
		cs.SetState(VotedStateObj)
		dposlog.Info("Change state because of check votes successfully.", "from", "VotingState", "to", "VotedState")

		cs.SetCurrentVote(voteItem)

		//1s后检查是否出块，是否需要重新投票
		cs.scheduleDPosTimeout(time.Millisecond * 500, VotedStateType)
	} else if result == continueToVote {
		dposlog.Info("VotingState get a vote, but don't get an agreement,waiting for new votes...")
	} else {
		dposlog.Info("VotingState get a vote, but don't get an agreement,vote fail,abort voting")
		//清理掉之前的选票记录，从初始状态重新开始
		cs.ClearVotes()
		cs.SetState(InitStateObj)
		dposlog.Info("Change state because of vote failed.", "from", "VotingState", "to", "InitState")
		cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
	}
}

func (voting *VotingState)sendVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply){
	dposlog.Info("VotingState don't support sendVoteReply,so do nothing")
	//cs.broadcastChannel <- MsgInfo{TypeID: dpostype.VoteReplyID, Msg: reply, PeerID: cs.ourID, PeerIP: ""}
}

func (voting *VotingState)recvVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply){
	dposlog.Info("VotingState recv Vote reply")
	voting.recvVote(cs, reply.Vote)
}

func (voting *VotingState)sendNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("VotingState does not support sendNotify,so do nothing")
}

func (voting *VotingState)recvNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("VotingState does not support recvNotify,so do nothing")
}

// VotingState is the voting state of dpos state machine until timeout or get an agreement by votes.
type VotedState struct {

}

func (voted *VotedState)timeOut(cs *ConsensusState) {
	now := time.Now().Unix()
	block := cs.client.GetCurrentBlock()
	task := DecideTaskByTime(now)

	if bytes.Equal(cs.privValidator.GetAddress(), cs.currentVote.VotedNodeAddress) {
		//当前节点为出块节点
		if now >= cs.currentVote.PeriodStop {
			//当前时间超过了节点切换时间，需要进行重新投票
			dposlog.Info("VotedState timeOut over periodStop.", "periodStop", cs.currentVote.PeriodStop)
			//当前时间超过了节点切换时间，需要进行重新投票
			notify := &dpostype.Notify{
					DPosNotify: &dpostype.DPosNotify{
						Vote:              cs.currentVote,
						HeightStop:        block.Height,
						HashStop:          block.Hash(),
						NotifyTimestamp:   now,
						NotifyNodeAddress: cs.privValidator.GetAddress(),
						NotifyNodeIndex:   int32(cs.privValidatorIndex),
					},
			}

			dposlog.Info("Send notify.", "HeightStop", notify.HeightStop, "HashStop", common.ToHex(notify.HashStop))

			if err := cs.privValidator.SignNotify(cs.validatorMgr.ChainID, notify); err != nil {
				dposlog.Error("SignNotify failed", "notify", notify.String())
				cs.SaveVote()
				cs.SaveMyVote()
				cs.ClearVotes()

				cs.SetState(InitStateObj)
				dposlog.Info("Change state because of time.", "from", "VotedState", "to", "InitState")
				cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections)*time.Millisecond, InitStateType)
				return
			}

			cs.SaveVote()
			cs.SaveMyVote()
			cs.SaveNotify()

			notify2 := *notify
			cs.SetNotify(notify2.DPosNotify)
			cs.dposState.sendNotify(cs, notify.DPosNotify)
			cs.ClearVotes()
			cs.SetState(InitStateObj)
			dposlog.Info("Change state because of time.", "from", "VotedState", "to", "InitState")
			cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections)*time.Millisecond, InitStateType)
			return
		} else {
			//如果区块未同步，则等待；如果区块已同步，则进行后续正常出块的判断和处理。
			if block.Height + 1 < cs.currentVote.Height{
				dposlog.Info("VotedState timeOut but block is not sync,wait...", "localHeight", block.Height, "vote height", cs.currentVote.Height)
				cs.scheduleDPosTimeout(time.Second*1, VotedStateType)
				return
			}

			//当前时间未到节点切换时间，则继续进行出块判断
			if block.BlockTime >= task.blockStop {
				//已出块，或者时间落后了。
				dposlog.Info("VotedState timeOut but block already is generated.", "blocktime", block.BlockTime, "blockStop", task.blockStop, "now", now)
				cs.scheduleDPosTimeout(time.Second * 1, VotedStateType)

				return
			} else if block.BlockTime < task.blockStart {
				//本出块周期尚未出块，则进行出块
				if task.blockStop-now <= 1 {
					dposlog.Info("Create new block.", "height", block.Height + 1)

					cs.client.SetBlockTime(task.blockStop)
					cs.client.CreateBlock()
					cs.scheduleDPosTimeout(time.Millisecond * 500, VotedStateType)
					return
				} else {
					dposlog.Info("Wait time to create block near blockStop.", "waittime", task.blockStop-now-1)
					//cs.scheduleDPosTimeout(time.Second * time.Duration(task.blockStop - now - 1), VotedStateType)
					cs.scheduleDPosTimeout(time.Millisecond * 500, VotedStateType)
					return
				}

			} else {
				//本周期已经出块
				dposlog.Info("Wait to next block cycle.", "waittime", task.blockStop-now+1)

				//cs.scheduleDPosTimeout(time.Second * time.Duration(task.blockStop-now+1), VotedStateType)
				cs.scheduleDPosTimeout(time.Millisecond * 500, VotedStateType)
				return
			}
		}
	} else{
		dposlog.Info("This node is not current owner.", "current owner index", cs.currentVote.VotedNodeIndex, "this node index", cs.client.ValidatorIndex())

		//非当前出块节点，如果到了切换出块节点的时间，则进行状态切换，进行投票
		if now >= cs.currentVote.PeriodStop {
			//当前时间超过了节点切换时间，需要进行重新投票
			cs.SaveVote()
			cs.SaveMyVote()
			cs.ClearVotes()
			cs.SetState(WaitNotifyStateObj)
			dposlog.Info("Change state because of time.", "from", "VotedState", "to", "WaitNotifyState")
			cs.scheduleDPosTimeout(time.Duration(timeoutWaitNotify) * time.Millisecond, WaitNotifyStateType)
			if cs.cachedNotify != nil {
				cs.dposState.recvNotify(cs, cs.cachedNotify)

			}
			return
		} else {
			//设置超时时间
			dposlog.Info("wait until change state.", "waittime", cs.currentVote.PeriodStop - now + 1)
			cs.scheduleDPosTimeout(time.Second * time.Duration(cs.currentVote.PeriodStop - now + 1), VotedStateType)
			return
		}
	}
}

func (voted *VotedState)sendVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply) {
	dposlog.Info("VotedState sendVoteReply")
	cs.broadcastChannel <- MsgInfo{TypeID: dpostype.VoteReplyID, Msg: reply, PeerID: cs.ourID, PeerIP: ""}
}

func (voted *VotedState)recvVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply) {
	dposlog.Info("VotedState recv Vote reply", "from index", reply.Vote.VoterNodeIndex, "local index", cs.privValidatorIndex)
	cs.AddVotes(reply.Vote)
}

func (voted *VotedState)sendVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	dposlog.Info("VotedState does not support sendVote,so do nothing")
}

func (voted *VotedState)recvVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	dposlog.Info("VotedState recv Vote, will reply it", "from index", vote.VoterNodeIndex, "local index", cs.privValidatorIndex)
	if cs.currentVote.PeriodStart >= vote.VoteItem.PeriodStart {
		vote2 := *cs.myVote
		reply := &dpostype.DPosVoteReply{Vote: &vote2}
		cs.dposState.sendVoteReply(cs, reply)
	} else {
		dposlog.Info("VotedState recv future Vote, will cache it")

		cs.CacheVotes(vote)
	}
}

func (voted *VotedState)sendNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	cs.broadcastChannel <- MsgInfo{TypeID: dpostype.NotifyID, Msg: notify, PeerID: cs.ourID, PeerIP: ""}
}

func (voted *VotedState)recvNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("VotedState recvNotify")

	if bytes.Equal(cs.privValidator.GetAddress(), cs.currentVote.VotedNodeAddress) {
		dposlog.Info("ignore recvNotify because this node is the owner now.")
		return
	} else {
		cs.CacheNotify(notify)
		cs.SaveVote()
		cs.SaveMyVote()
		cs.ClearVotes()
		cs.SetState(WaitNotifyStateObj)
		dposlog.Info("Change state because of recv notify.", "from", "VotedState", "to", "WaitNotifyState")
		cs.scheduleDPosTimeout(time.Duration(timeoutWaitNotify)*time.Millisecond, WaitNotifyStateType)
		if cs.cachedNotify != nil {
			cs.dposState.recvNotify(cs, cs.cachedNotify)

		}
		return
	}
}

// VotingState is the voting state of dpos state machine until timeout or get an agreement by votes.
type WaitNofifyState struct {

}

func (wait *WaitNofifyState)timeOut(cs *ConsensusState) {
	//cs.clearVotes()
	cs.SetState(InitStateObj)
	dposlog.Info("Change state because of time.", "from", "WaitNofifyState", "to", "InitState")
	cs.scheduleDPosTimeout(time.Duration(timeoutCheckConnections) * time.Millisecond, InitStateType)
}

func (wait *WaitNofifyState)sendVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	dposlog.Info("WaitNofifyState does not support sendVote,so do nothing")
}

func (wait *WaitNofifyState)recvVote(cs *ConsensusState,vote * dpostype.DPosVote) {
	dposlog.Info("WaitNofifyState recvVote,store it.")
	//对于vote进行保存，在后续状态中进行处理。 保存的选票有先后，同一个节点发来的最新的选票被保存。
	//if !cs.VerifyVote(vote) {
	//	dposlog.Info("VotingState verify vote failed", "vote", vote.String())
	//	return
	//}

	//cs.AddVotes(vote)
	cs.CacheVotes(vote)
}

func (wait *WaitNofifyState)sendVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply) {
	dposlog.Info("WaitNofifyState does not support sendVoteReply,so do nothing")
}

func (wait *WaitNofifyState)recvVoteReply(cs *ConsensusState, reply * dpostype.DPosVoteReply){
	dposlog.Info("WaitNofifyState recv Vote reply,ignore it.")
}

func (wait *WaitNofifyState)sendNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("WaitNofifyState does not support sendNotify,so do nothing")
}

func (wait *WaitNofifyState)recvNotify(cs *ConsensusState, notify * dpostype.DPosNotify) {
	dposlog.Info("WaitNofifyState recvNotify")
	//记录Notify，校验区块，标记不可逆区块高度
	if !cs.VerifyNotify(notify) {
		dposlog.Info("VotedState verify notify failed")
		return
	}

	block := cs.client.GetCurrentBlock()
	if (block.Height > notify.HeightStop){
		dposlog.Info("Local block height is advanced than notify, discard it.", "localheight", block.Height, "notify", notify.String())
		return
	} else if block.Height == notify.HeightStop && bytes.Equal(block.Hash(), notify.HashStop){
		dposlog.Info("Local block height is sync with notify", "notify", notify.String())
	} else {
		dposlog.Info("Local block height is not sync with notify", "localheight", cs.client.GetCurrentHeight(), "notify", notify.String())
		hint := time.NewTicker(3 * time.Second)
		beg := time.Now()
		catchupFlag := false
		OuterLoop:
			for !catchupFlag {
				select {
				case <-hint.C:
					dposlog.Info("Still catching up max height......", "Height", cs.client.GetCurrentHeight(),"notifyHeight", notify.HeightStop, "cost", time.Since(beg))
					if cs.client.IsCaughtUp() &&  cs.client.GetCurrentHeight() >= notify.HeightStop {
						dposlog.Info("This node has caught up max height", "Height", cs.client.GetCurrentHeight(), "isHashSame", bytes.Equal(block.Hash(), notify.HashStop))
						break OuterLoop
					}

				default:
					if cs.client.IsCaughtUp() &&  cs.client.GetCurrentHeight() >= notify.HeightStop {
						dposlog.Info("This node has caught up max height", "Height", cs.client.GetCurrentHeight())
						break OuterLoop
					}
					time.Sleep(time.Second)
				}
			}
		hint.Stop()
	}

	cs.ClearCachedNotify()
	cs.SaveNotify()
	cs.SetNotify(notify)

	//cs.clearVotes()
	cs.SetState(InitStateObj)
	dposlog.Info("Change state because recv notify.", "from", "WaitNofifyState", "to", "InitState")
	cs.dposState.timeOut(cs)
	//cs.scheduleDPosTimeout(time.Second * 1, InitStateType)
}

var InitStateObj  = & InitState{}
var VotingStateObj  = & VotingState{}
var VotedStateObj  = & VotedState{}
var WaitNotifyStateObj = & WaitNofifyState{}