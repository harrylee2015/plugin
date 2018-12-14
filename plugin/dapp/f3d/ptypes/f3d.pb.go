// Code generated by protoc-gen-go. DO NOT EDIT.
// source: f3d.proto

/*
Package types is a generated protocol buffer package.

It is generated from these files:
	f3d.proto
	rpc.proto

It has these top-level messages:
	RoundInfo
	RoundsInfo
	KeyInfo
	F3DAction
	F3DStart
	F3DLuckyDraw
	F3DBuyKey
	QueryF3DByRound
	QueryF3DLastRound
	QueryF3DListByRound
	QueryBuyRecordByRoundAndAddr
	QueryKeyCountByRoundAndAddr
	QueryAddrInfo
	F3DRecord
	F3DStartRound
	F3DDrawRound
	F3DBuyRecord
	ReplyF3DList
	ReplyF3D
	ReplyBuyRecord
	ReplyKey
	ReplyKeyCount
	ReplyAddrInfoList
	AddrInfo
	ReceiptF3D
	Config
	GameStartReq
	GameDrawReq
	GameBuyKeysReq
	KeyInfoQueryReq
	RoundInfoQueryReq
*/
package types

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RoundInfo struct {
	// 游戏轮次
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	// 本轮游戏开始事件
	BeginTime int64 `protobuf:"varint,2,opt,name=beginTime" json:"beginTime,omitempty"`
	// 本轮游戏结束时间
	EndTime int64 `protobuf:"varint,3,opt,name=endTime" json:"endTime,omitempty"`
	// 本轮游戏目前为止最后一把钥匙持有人（游戏开奖时，这个就是中奖人）
	LastOwner string `protobuf:"bytes,4,opt,name=lastOwner" json:"lastOwner,omitempty"`
	// 最后一把钥匙的购买时间
	LastKeyTime int64 `protobuf:"varint,5,opt,name=lastKeyTime" json:"lastKeyTime,omitempty"`
	// 最后一把钥匙的价格
	LastKeyPrice float32 `protobuf:"fixed32,6,opt,name=lastKeyPrice" json:"lastKeyPrice,omitempty"`
	// 本轮游戏奖金池总额
	BonusPool float32 `protobuf:"fixed32,7,opt,name=bonusPool" json:"bonusPool,omitempty"`
	// 本轮游戏参与地址数
	UserCount int64 `protobuf:"varint,8,opt,name=userCount" json:"userCount,omitempty"`
	// 本轮游戏募集到的key个数
	KeyCount int64 `protobuf:"varint,9,opt,name=keyCount" json:"keyCount,omitempty"`
	// 距离开奖剩余时间
	RemainTime int64 `protobuf:"varint,10,opt,name=remainTime" json:"remainTime,omitempty"`
	UpdateTime int64 `protobuf:"varint,11,opt,name=updateTime" json:"updateTime,omitempty"`
}

func (m *RoundInfo) Reset()                    { *m = RoundInfo{} }
func (m *RoundInfo) String() string            { return proto.CompactTextString(m) }
func (*RoundInfo) ProtoMessage()               {}
func (*RoundInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RoundInfo) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *RoundInfo) GetBeginTime() int64 {
	if m != nil {
		return m.BeginTime
	}
	return 0
}

func (m *RoundInfo) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *RoundInfo) GetLastOwner() string {
	if m != nil {
		return m.LastOwner
	}
	return ""
}

func (m *RoundInfo) GetLastKeyTime() int64 {
	if m != nil {
		return m.LastKeyTime
	}
	return 0
}

func (m *RoundInfo) GetLastKeyPrice() float32 {
	if m != nil {
		return m.LastKeyPrice
	}
	return 0
}

func (m *RoundInfo) GetBonusPool() float32 {
	if m != nil {
		return m.BonusPool
	}
	return 0
}

func (m *RoundInfo) GetUserCount() int64 {
	if m != nil {
		return m.UserCount
	}
	return 0
}

func (m *RoundInfo) GetKeyCount() int64 {
	if m != nil {
		return m.KeyCount
	}
	return 0
}

func (m *RoundInfo) GetRemainTime() int64 {
	if m != nil {
		return m.RemainTime
	}
	return 0
}

func (m *RoundInfo) GetUpdateTime() int64 {
	if m != nil {
		return m.UpdateTime
	}
	return 0
}

type RoundsInfo struct {
	RoundsInfo []*RoundInfo `protobuf:"bytes,1,rep,name=roundsInfo" json:"roundsInfo,omitempty"`
}

func (m *RoundsInfo) Reset()                    { *m = RoundsInfo{} }
func (m *RoundsInfo) String() string            { return proto.CompactTextString(m) }
func (*RoundsInfo) ProtoMessage()               {}
func (*RoundsInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RoundsInfo) GetRoundsInfo() []*RoundInfo {
	if m != nil {
		return m.RoundsInfo
	}
	return nil
}

type KeyInfo struct {
	// 游戏轮次  (是由系统合约填写后存储）
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	// 本次购买key的价格 (是由系统合约填写后存储）
	KeyPrice float32 `protobuf:"fixed32,2,opt,name=keyPrice" json:"keyPrice,omitempty"`
	// 用户本次买的key的数量
	KeyNum int64 `protobuf:"varint,3,opt,name=keyNum" json:"keyNum,omitempty"`
	// 用户地址 (是由系统合约填写后存储）
	Addr string `protobuf:"bytes,4,opt,name=addr" json:"addr,omitempty"`
	// 交易确认存储时间（被打包的时间）
	BuyKeyTime int64 `protobuf:"varint,5,opt,name=buyKeyTime" json:"buyKeyTime,omitempty"`
	// 买票的txHash
	BuyKeyTxHash string `protobuf:"bytes,6,opt,name=buyKeyTxHash" json:"buyKeyTxHash,omitempty"`
	Index        int64  `protobuf:"varint,7,opt,name=index" json:"index,omitempty"`
}

func (m *KeyInfo) Reset()                    { *m = KeyInfo{} }
func (m *KeyInfo) String() string            { return proto.CompactTextString(m) }
func (*KeyInfo) ProtoMessage()               {}
func (*KeyInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *KeyInfo) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *KeyInfo) GetKeyPrice() float32 {
	if m != nil {
		return m.KeyPrice
	}
	return 0
}

func (m *KeyInfo) GetKeyNum() int64 {
	if m != nil {
		return m.KeyNum
	}
	return 0
}

func (m *KeyInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *KeyInfo) GetBuyKeyTime() int64 {
	if m != nil {
		return m.BuyKeyTime
	}
	return 0
}

func (m *KeyInfo) GetBuyKeyTxHash() string {
	if m != nil {
		return m.BuyKeyTxHash
	}
	return ""
}

func (m *KeyInfo) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

// message for execs.f3d
type F3DAction struct {
	// Types that are valid to be assigned to Value:
	//	*F3DAction_Start
	//	*F3DAction_Draw
	//	*F3DAction_Buy
	Value isF3DAction_Value `protobuf_oneof:"value"`
	Ty    int32             `protobuf:"varint,4,opt,name=ty" json:"ty,omitempty"`
}

func (m *F3DAction) Reset()                    { *m = F3DAction{} }
func (m *F3DAction) String() string            { return proto.CompactTextString(m) }
func (*F3DAction) ProtoMessage()               {}
func (*F3DAction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isF3DAction_Value interface {
	isF3DAction_Value()
}

type F3DAction_Start struct {
	Start *F3DStart `protobuf:"bytes,1,opt,name=start,oneof"`
}
type F3DAction_Draw struct {
	Draw *F3DLuckyDraw `protobuf:"bytes,2,opt,name=draw,oneof"`
}
type F3DAction_Buy struct {
	Buy *F3DBuyKey `protobuf:"bytes,3,opt,name=buy,oneof"`
}

func (*F3DAction_Start) isF3DAction_Value() {}
func (*F3DAction_Draw) isF3DAction_Value()  {}
func (*F3DAction_Buy) isF3DAction_Value()   {}

func (m *F3DAction) GetValue() isF3DAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *F3DAction) GetStart() *F3DStart {
	if x, ok := m.GetValue().(*F3DAction_Start); ok {
		return x.Start
	}
	return nil
}

func (m *F3DAction) GetDraw() *F3DLuckyDraw {
	if x, ok := m.GetValue().(*F3DAction_Draw); ok {
		return x.Draw
	}
	return nil
}

func (m *F3DAction) GetBuy() *F3DBuyKey {
	if x, ok := m.GetValue().(*F3DAction_Buy); ok {
		return x.Buy
	}
	return nil
}

func (m *F3DAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*F3DAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _F3DAction_OneofMarshaler, _F3DAction_OneofUnmarshaler, _F3DAction_OneofSizer, []interface{}{
		(*F3DAction_Start)(nil),
		(*F3DAction_Draw)(nil),
		(*F3DAction_Buy)(nil),
	}
}

func _F3DAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*F3DAction)
	// value
	switch x := m.Value.(type) {
	case *F3DAction_Start:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Start); err != nil {
			return err
		}
	case *F3DAction_Draw:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Draw); err != nil {
			return err
		}
	case *F3DAction_Buy:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Buy); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("F3DAction.Value has unexpected type %T", x)
	}
	return nil
}

func _F3DAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*F3DAction)
	switch tag {
	case 1: // value.start
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(F3DStart)
		err := b.DecodeMessage(msg)
		m.Value = &F3DAction_Start{msg}
		return true, err
	case 2: // value.draw
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(F3DLuckyDraw)
		err := b.DecodeMessage(msg)
		m.Value = &F3DAction_Draw{msg}
		return true, err
	case 3: // value.buy
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(F3DBuyKey)
		err := b.DecodeMessage(msg)
		m.Value = &F3DAction_Buy{msg}
		return true, err
	default:
		return false, nil
	}
}

func _F3DAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*F3DAction)
	// value
	switch x := m.Value.(type) {
	case *F3DAction_Start:
		s := proto.Size(x.Start)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *F3DAction_Draw:
		s := proto.Size(x.Draw)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *F3DAction_Buy:
		s := proto.Size(x.Buy)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type F3DStart struct {
	// 轮次，这个填不填不重要，合约里面会自动校验的
	Round int64 `protobuf:"varint,1,opt,name=Round" json:"Round,omitempty"`
}

func (m *F3DStart) Reset()                    { *m = F3DStart{} }
func (m *F3DStart) String() string            { return proto.CompactTextString(m) }
func (*F3DStart) ProtoMessage()               {}
func (*F3DStart) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *F3DStart) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type F3DLuckyDraw struct {
	// 轮次，这个填不填不重要，合约里面会自动校验的
	Round int64 `protobuf:"varint,1,opt,name=Round" json:"Round,omitempty"`
}

func (m *F3DLuckyDraw) Reset()                    { *m = F3DLuckyDraw{} }
func (m *F3DLuckyDraw) String() string            { return proto.CompactTextString(m) }
func (*F3DLuckyDraw) ProtoMessage()               {}
func (*F3DLuckyDraw) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *F3DLuckyDraw) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type F3DBuyKey struct {
	// 用户本次买的key的数量
	KeyNum int64 `protobuf:"varint,3,opt,name=keyNum" json:"keyNum,omitempty"`
}

func (m *F3DBuyKey) Reset()                    { *m = F3DBuyKey{} }
func (m *F3DBuyKey) String() string            { return proto.CompactTextString(m) }
func (*F3DBuyKey) ProtoMessage()               {}
func (*F3DBuyKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *F3DBuyKey) GetKeyNum() int64 {
	if m != nil {
		return m.KeyNum
	}
	return 0
}

// 查询f3d 游戏信息,这里面其实包含了key的最新价格信息
type QueryF3DByRound struct {
	// 轮次，默认查询最新的
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
}

func (m *QueryF3DByRound) Reset()                    { *m = QueryF3DByRound{} }
func (m *QueryF3DByRound) String() string            { return proto.CompactTextString(m) }
func (*QueryF3DByRound) ProtoMessage()               {}
func (*QueryF3DByRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *QueryF3DByRound) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type QueryF3DLastRound struct {
}

func (m *QueryF3DLastRound) Reset()                    { *m = QueryF3DLastRound{} }
func (m *QueryF3DLastRound) String() string            { return proto.CompactTextString(m) }
func (*QueryF3DLastRound) ProtoMessage()               {}
func (*QueryF3DLastRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type QueryF3DListByRound struct {
	// 轮次，默认查询最新的
	StartRound int64 `protobuf:"varint,1,opt,name=startRound" json:"startRound,omitempty"`
	// 单页返回多少条记录，默认返回10条，单次最多返回50条
	Count int32 `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
	// 0降序，1升序，默认降序
	Direction int32 `protobuf:"varint,5,opt,name=direction" json:"direction,omitempty"`
}

func (m *QueryF3DListByRound) Reset()                    { *m = QueryF3DListByRound{} }
func (m *QueryF3DListByRound) String() string            { return proto.CompactTextString(m) }
func (*QueryF3DListByRound) ProtoMessage()               {}
func (*QueryF3DListByRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *QueryF3DListByRound) GetStartRound() int64 {
	if m != nil {
		return m.StartRound
	}
	return 0
}

func (m *QueryF3DListByRound) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *QueryF3DListByRound) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

// key 信息查询
type QueryBuyRecordByRoundAndAddr struct {
	// 轮次,必填参数
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	// 用户地址
	Addr  string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
	Index int64  `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *QueryBuyRecordByRoundAndAddr) Reset()                    { *m = QueryBuyRecordByRoundAndAddr{} }
func (m *QueryBuyRecordByRoundAndAddr) String() string            { return proto.CompactTextString(m) }
func (*QueryBuyRecordByRoundAndAddr) ProtoMessage()               {}
func (*QueryBuyRecordByRoundAndAddr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *QueryBuyRecordByRoundAndAddr) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *QueryBuyRecordByRoundAndAddr) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *QueryBuyRecordByRoundAndAddr) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

// 用户key数量查询
type QueryKeyCountByRoundAndAddr struct {
	// 轮次,必填参数
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	// 用户地址
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *QueryKeyCountByRoundAndAddr) Reset()                    { *m = QueryKeyCountByRoundAndAddr{} }
func (m *QueryKeyCountByRoundAndAddr) String() string            { return proto.CompactTextString(m) }
func (*QueryKeyCountByRoundAndAddr) ProtoMessage()               {}
func (*QueryKeyCountByRoundAndAddr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *QueryKeyCountByRoundAndAddr) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *QueryKeyCountByRoundAndAddr) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type QueryAddrInfo struct {
	Round     int64  `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	Addr      string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
	Count     int32  `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
	Direction int32  `protobuf:"varint,4,opt,name=direction" json:"direction,omitempty"`
}

func (m *QueryAddrInfo) Reset()                    { *m = QueryAddrInfo{} }
func (m *QueryAddrInfo) String() string            { return proto.CompactTextString(m) }
func (*QueryAddrInfo) ProtoMessage()               {}
func (*QueryAddrInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *QueryAddrInfo) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *QueryAddrInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *QueryAddrInfo) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *QueryAddrInfo) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

type F3DRecord struct {
	// 用户地址
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	// index
	Index int64 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	// round
	Round int64 `protobuf:"varint,3,opt,name=round" json:"round,omitempty"`
}

func (m *F3DRecord) Reset()                    { *m = F3DRecord{} }
func (m *F3DRecord) String() string            { return proto.CompactTextString(m) }
func (*F3DRecord) ProtoMessage()               {}
func (*F3DRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *F3DRecord) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *F3DRecord) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *F3DRecord) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type F3DStartRound struct {
	// round
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
}

func (m *F3DStartRound) Reset()                    { *m = F3DStartRound{} }
func (m *F3DStartRound) String() string            { return proto.CompactTextString(m) }
func (*F3DStartRound) ProtoMessage()               {}
func (*F3DStartRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *F3DStartRound) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type F3DDrawRound struct {
	// round
	Round int64 `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
}

func (m *F3DDrawRound) Reset()                    { *m = F3DDrawRound{} }
func (m *F3DDrawRound) String() string            { return proto.CompactTextString(m) }
func (*F3DDrawRound) ProtoMessage()               {}
func (*F3DDrawRound) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *F3DDrawRound) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

type F3DBuyRecord struct {
	Round int64  `protobuf:"varint,1,opt,name=round" json:"round,omitempty"`
	Addr  string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
	Index int64  `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *F3DBuyRecord) Reset()                    { *m = F3DBuyRecord{} }
func (m *F3DBuyRecord) String() string            { return proto.CompactTextString(m) }
func (*F3DBuyRecord) ProtoMessage()               {}
func (*F3DBuyRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *F3DBuyRecord) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *F3DBuyRecord) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *F3DBuyRecord) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

// f3d round查询返回数据
type ReplyF3DList struct {
	Rounds []*RoundInfo `protobuf:"bytes,1,rep,name=rounds" json:"rounds,omitempty"`
}

func (m *ReplyF3DList) Reset()                    { *m = ReplyF3DList{} }
func (m *ReplyF3DList) String() string            { return proto.CompactTextString(m) }
func (*ReplyF3DList) ProtoMessage()               {}
func (*ReplyF3DList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *ReplyF3DList) GetRounds() []*RoundInfo {
	if m != nil {
		return m.Rounds
	}
	return nil
}

type ReplyF3D struct {
	Round *RoundInfo `protobuf:"bytes,1,opt,name=round" json:"round,omitempty"`
}

func (m *ReplyF3D) Reset()                    { *m = ReplyF3D{} }
func (m *ReplyF3D) String() string            { return proto.CompactTextString(m) }
func (*ReplyF3D) ProtoMessage()               {}
func (*ReplyF3D) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *ReplyF3D) GetRound() *RoundInfo {
	if m != nil {
		return m.Round
	}
	return nil
}

// 用户查询买的key信息返回数据
type ReplyBuyRecord struct {
	RecordList []*KeyInfo `protobuf:"bytes,1,rep,name=recordList" json:"recordList,omitempty"`
}

func (m *ReplyBuyRecord) Reset()                    { *m = ReplyBuyRecord{} }
func (m *ReplyBuyRecord) String() string            { return proto.CompactTextString(m) }
func (*ReplyBuyRecord) ProtoMessage()               {}
func (*ReplyBuyRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *ReplyBuyRecord) GetRecordList() []*KeyInfo {
	if m != nil {
		return m.RecordList
	}
	return nil
}

type ReplyKey struct {
	Key *KeyInfo `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *ReplyKey) Reset()                    { *m = ReplyKey{} }
func (m *ReplyKey) String() string            { return proto.CompactTextString(m) }
func (*ReplyKey) ProtoMessage()               {}
func (*ReplyKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *ReplyKey) GetKey() *KeyInfo {
	if m != nil {
		return m.Key
	}
	return nil
}

type ReplyKeyCount struct {
	Count int64 `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
}

func (m *ReplyKeyCount) Reset()                    { *m = ReplyKeyCount{} }
func (m *ReplyKeyCount) String() string            { return proto.CompactTextString(m) }
func (*ReplyKeyCount) ProtoMessage()               {}
func (*ReplyKeyCount) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

func (m *ReplyKeyCount) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type ReplyAddrInfoList struct {
	AddrInfoList []*AddrInfo `protobuf:"bytes,1,rep,name=addrInfoList" json:"addrInfoList,omitempty"`
}

func (m *ReplyAddrInfoList) Reset()                    { *m = ReplyAddrInfoList{} }
func (m *ReplyAddrInfoList) String() string            { return proto.CompactTextString(m) }
func (*ReplyAddrInfoList) ProtoMessage()               {}
func (*ReplyAddrInfoList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

func (m *ReplyAddrInfoList) GetAddrInfoList() []*AddrInfo {
	if m != nil {
		return m.AddrInfoList
	}
	return nil
}

type AddrInfo struct {
	Addr     string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	KeyNum   int64  `protobuf:"varint,2,opt,name=keyNum" json:"keyNum,omitempty"`
	BuyCount int64  `protobuf:"varint,3,opt,name=buyCount" json:"buyCount,omitempty"`
	Round    int64  `protobuf:"varint,4,opt,name=round" json:"round,omitempty"`
}

func (m *AddrInfo) Reset()                    { *m = AddrInfo{} }
func (m *AddrInfo) String() string            { return proto.CompactTextString(m) }
func (*AddrInfo) ProtoMessage()               {}
func (*AddrInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{23} }

func (m *AddrInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *AddrInfo) GetKeyNum() int64 {
	if m != nil {
		return m.KeyNum
	}
	return 0
}

func (m *AddrInfo) GetBuyCount() int64 {
	if m != nil {
		return m.BuyCount
	}
	return 0
}

func (m *AddrInfo) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

// 合约内部日志记录，待补全
type ReceiptF3D struct {
	Addr     string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	Round    int64  `protobuf:"varint,2,opt,name=round" json:"round,omitempty"`
	Index    int64  `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
	Action   int64  `protobuf:"varint,4,opt,name=action" json:"action,omitempty"`
	BuyCount int64  `protobuf:"varint,5,opt,name=buyCount" json:"buyCount,omitempty"`
	KeyNum   int64  `protobuf:"varint,6,opt,name=keyNum" json:"keyNum,omitempty"`
}

func (m *ReceiptF3D) Reset()                    { *m = ReceiptF3D{} }
func (m *ReceiptF3D) String() string            { return proto.CompactTextString(m) }
func (*ReceiptF3D) ProtoMessage()               {}
func (*ReceiptF3D) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{24} }

func (m *ReceiptF3D) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *ReceiptF3D) GetRound() int64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *ReceiptF3D) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ReceiptF3D) GetAction() int64 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *ReceiptF3D) GetBuyCount() int64 {
	if m != nil {
		return m.BuyCount
	}
	return 0
}

func (m *ReceiptF3D) GetKeyNum() int64 {
	if m != nil {
		return m.KeyNum
	}
	return 0
}

type Config struct {
	Manager        string  `protobuf:"bytes,1,opt,name=manager" json:"manager,omitempty"`
	Developer      string  `protobuf:"bytes,2,opt,name=developer" json:"developer,omitempty"`
	WinnerBonus    float32 `protobuf:"fixed32,3,opt,name=winnerBonus" json:"winnerBonus,omitempty"`
	KeyBonus       float32 `protobuf:"fixed32,4,opt,name=keyBonus" json:"keyBonus,omitempty"`
	PoolBonus      float32 `protobuf:"fixed32,5,opt,name=poolBonus" json:"poolBonus,omitempty"`
	DeveloperBonus float32 `protobuf:"fixed32,6,opt,name=developerBonus" json:"developerBonus,omitempty"`
	LifeTime       int64   `protobuf:"varint,7,opt,name=lifeTime" json:"lifeTime,omitempty"`
	KeyIncrTime    int64   `protobuf:"varint,8,opt,name=keyIncrTime" json:"keyIncrTime,omitempty"`
	MaxkeyIncrTime int64   `protobuf:"varint,9,opt,name=maxkeyIncrTime" json:"maxkeyIncrTime,omitempty"`
	NouserDecrTime int64   `protobuf:"varint,10,opt,name=nouserDecrTime" json:"nouserDecrTime,omitempty"`
	DrawDelayTime  int64   `protobuf:"varint,11,opt,name=drawDelayTime" json:"drawDelayTime,omitempty"`
	StartKeyPrice  float32 `protobuf:"fixed32,12,opt,name=startKeyPrice" json:"startKeyPrice,omitempty"`
	IncrKeyPrice   float32 `protobuf:"fixed32,13,opt,name=incrKeyPrice" json:"incrKeyPrice,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{25} }

func (m *Config) GetManager() string {
	if m != nil {
		return m.Manager
	}
	return ""
}

func (m *Config) GetDeveloper() string {
	if m != nil {
		return m.Developer
	}
	return ""
}

func (m *Config) GetWinnerBonus() float32 {
	if m != nil {
		return m.WinnerBonus
	}
	return 0
}

func (m *Config) GetKeyBonus() float32 {
	if m != nil {
		return m.KeyBonus
	}
	return 0
}

func (m *Config) GetPoolBonus() float32 {
	if m != nil {
		return m.PoolBonus
	}
	return 0
}

func (m *Config) GetDeveloperBonus() float32 {
	if m != nil {
		return m.DeveloperBonus
	}
	return 0
}

func (m *Config) GetLifeTime() int64 {
	if m != nil {
		return m.LifeTime
	}
	return 0
}

func (m *Config) GetKeyIncrTime() int64 {
	if m != nil {
		return m.KeyIncrTime
	}
	return 0
}

func (m *Config) GetMaxkeyIncrTime() int64 {
	if m != nil {
		return m.MaxkeyIncrTime
	}
	return 0
}

func (m *Config) GetNouserDecrTime() int64 {
	if m != nil {
		return m.NouserDecrTime
	}
	return 0
}

func (m *Config) GetDrawDelayTime() int64 {
	if m != nil {
		return m.DrawDelayTime
	}
	return 0
}

func (m *Config) GetStartKeyPrice() float32 {
	if m != nil {
		return m.StartKeyPrice
	}
	return 0
}

func (m *Config) GetIncrKeyPrice() float32 {
	if m != nil {
		return m.IncrKeyPrice
	}
	return 0
}

func init() {
	proto.RegisterType((*RoundInfo)(nil), "types.RoundInfo")
	proto.RegisterType((*RoundsInfo)(nil), "types.RoundsInfo")
	proto.RegisterType((*KeyInfo)(nil), "types.KeyInfo")
	proto.RegisterType((*F3DAction)(nil), "types.F3dAction")
	proto.RegisterType((*F3DStart)(nil), "types.F3dStart")
	proto.RegisterType((*F3DLuckyDraw)(nil), "types.F3dLuckyDraw")
	proto.RegisterType((*F3DBuyKey)(nil), "types.F3dBuyKey")
	proto.RegisterType((*QueryF3DByRound)(nil), "types.QueryF3dByRound")
	proto.RegisterType((*QueryF3DLastRound)(nil), "types.QueryF3dLastRound")
	proto.RegisterType((*QueryF3DListByRound)(nil), "types.QueryF3dListByRound")
	proto.RegisterType((*QueryBuyRecordByRoundAndAddr)(nil), "types.QueryBuyRecordByRoundAndAddr")
	proto.RegisterType((*QueryKeyCountByRoundAndAddr)(nil), "types.QueryKeyCountByRoundAndAddr")
	proto.RegisterType((*QueryAddrInfo)(nil), "types.QueryAddrInfo")
	proto.RegisterType((*F3DRecord)(nil), "types.F3dRecord")
	proto.RegisterType((*F3DStartRound)(nil), "types.F3dStartRound")
	proto.RegisterType((*F3DDrawRound)(nil), "types.F3dDrawRound")
	proto.RegisterType((*F3DBuyRecord)(nil), "types.F3dBuyRecord")
	proto.RegisterType((*ReplyF3DList)(nil), "types.ReplyF3dList")
	proto.RegisterType((*ReplyF3D)(nil), "types.ReplyF3d")
	proto.RegisterType((*ReplyBuyRecord)(nil), "types.ReplyBuyRecord")
	proto.RegisterType((*ReplyKey)(nil), "types.ReplyKey")
	proto.RegisterType((*ReplyKeyCount)(nil), "types.ReplyKeyCount")
	proto.RegisterType((*ReplyAddrInfoList)(nil), "types.ReplyAddrInfoList")
	proto.RegisterType((*AddrInfo)(nil), "types.AddrInfo")
	proto.RegisterType((*ReceiptF3D)(nil), "types.ReceiptF3d")
	proto.RegisterType((*Config)(nil), "types.Config")
}

func init() { proto.RegisterFile("f3d.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 964 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0x36, 0x49, 0x51, 0x96, 0x46, 0xb2, 0xd3, 0xac, 0x8b, 0x42, 0x48, 0x8d, 0x40, 0x60, 0xdd,
	0x44, 0x05, 0x0a, 0xa3, 0xb0, 0x2e, 0x3d, 0x15, 0xb5, 0x63, 0xb8, 0x2a, 0x1c, 0xa4, 0xe9, 0xb6,
	0xe7, 0x02, 0x94, 0xb8, 0x76, 0x08, 0x49, 0xa4, 0xb0, 0x24, 0x63, 0xf3, 0x4d, 0x7a, 0xea, 0xbb,
	0xb4, 0xef, 0xd0, 0xf7, 0x29, 0x66, 0xf6, 0x87, 0x4b, 0x41, 0xf2, 0xa1, 0xb9, 0x71, 0xbe, 0x99,
	0xe1, 0xfc, 0x7d, 0xbb, 0x3b, 0xd0, 0xbf, 0x9b, 0x26, 0xe7, 0x1b, 0x99, 0x97, 0x39, 0x0b, 0xcb,
	0x7a, 0x23, 0x8a, 0xe8, 0x5f, 0x1f, 0xfa, 0x3c, 0xaf, 0xb2, 0xe4, 0xe7, 0xec, 0x2e, 0x67, 0x9f,
	0x43, 0x28, 0x51, 0x18, 0x79, 0x63, 0x6f, 0x12, 0x70, 0x25, 0xb0, 0x53, 0xe8, 0xcf, 0xc5, 0x7d,
	0x9a, 0xfd, 0x9e, 0xae, 0xc5, 0xc8, 0x27, 0x4d, 0x03, 0xb0, 0x11, 0x1c, 0x8a, 0x2c, 0x21, 0x5d,
	0x40, 0x3a, 0x23, 0xa2, 0xdf, 0x2a, 0x2e, 0xca, 0x5f, 0x1e, 0x32, 0x21, 0x47, 0x9d, 0xb1, 0x37,
	0xe9, 0xf3, 0x06, 0x60, 0x63, 0x18, 0xa0, 0x70, 0x2b, 0x6a, 0xf2, 0x0d, 0xc9, 0xd7, 0x85, 0x58,
	0x04, 0x43, 0x2d, 0xbe, 0x97, 0xe9, 0x42, 0x8c, 0xba, 0x63, 0x6f, 0xe2, 0xf3, 0x16, 0x46, 0xb9,
	0xe5, 0x59, 0x55, 0xbc, 0xcf, 0xf3, 0xd5, 0xe8, 0x90, 0x0c, 0x1a, 0x00, 0xb5, 0x55, 0x21, 0xe4,
	0x9b, 0xbc, 0xca, 0xca, 0x51, 0x4f, 0x65, 0x6e, 0x01, 0xf6, 0x02, 0x7a, 0x4b, 0x51, 0x2b, 0x65,
	0x9f, 0x94, 0x56, 0x66, 0x2f, 0x01, 0xa4, 0x58, 0xc7, 0xba, 0x68, 0x20, 0xad, 0x83, 0xa0, 0xbe,
	0xda, 0x24, 0x71, 0x29, 0x48, 0x3f, 0x50, 0xfa, 0x06, 0x89, 0x7e, 0x00, 0xa0, 0xb6, 0x16, 0xd4,
	0xd7, 0xef, 0x00, 0xa4, 0x95, 0x46, 0xde, 0x38, 0x98, 0x0c, 0x2e, 0x3e, 0x3b, 0xa7, 0x09, 0x9c,
	0xdb, 0xee, 0x73, 0xc7, 0x26, 0xfa, 0xdb, 0x83, 0xc3, 0x5b, 0x51, 0x3f, 0x31, 0x15, 0x95, 0xbd,
	0xea, 0x8c, 0x4f, 0x85, 0x5b, 0x99, 0x7d, 0x01, 0xdd, 0xa5, 0xa8, 0xdf, 0x55, 0x6b, 0x3d, 0x12,
	0x2d, 0x31, 0x06, 0x9d, 0x38, 0x49, 0xcc, 0x30, 0xe8, 0x1b, 0x2b, 0x99, 0x57, 0x75, 0x7b, 0x0c,
	0x0e, 0x82, 0x53, 0xd0, 0xd2, 0xe3, 0x2c, 0x2e, 0x3e, 0xd0, 0x14, 0xfa, 0xbc, 0x85, 0x61, 0x86,
	0x69, 0x96, 0x88, 0x47, 0x9a, 0x40, 0xc0, 0x95, 0x10, 0xfd, 0xe5, 0x41, 0xff, 0x66, 0x9a, 0x5c,
	0x2e, 0xca, 0x34, 0xcf, 0xd8, 0x6b, 0x08, 0x8b, 0x32, 0x96, 0x25, 0x55, 0x31, 0xb8, 0x78, 0xa6,
	0xcb, 0xbf, 0x99, 0x26, 0xbf, 0x21, 0x3c, 0x3b, 0xe0, 0x4a, 0xcf, 0xbe, 0x81, 0x4e, 0x22, 0xe3,
	0x07, 0x2a, 0x6a, 0x70, 0x71, 0xd2, 0xd8, 0xbd, 0xad, 0x16, 0xcb, 0xfa, 0x5a, 0xc6, 0x0f, 0xb3,
	0x03, 0x4e, 0x26, 0xec, 0x0c, 0x82, 0x79, 0x55, 0x53, 0x91, 0x4d, 0x43, 0x6f, 0xa6, 0xc9, 0x15,
	0x25, 0x37, 0x3b, 0xe0, 0xa8, 0x66, 0xc7, 0xe0, 0x97, 0x35, 0xd5, 0x1c, 0x72, 0xbf, 0xac, 0xaf,
	0x0e, 0x21, 0xfc, 0x18, 0xaf, 0x2a, 0x11, 0x8d, 0xa1, 0x67, 0xc2, 0x63, 0x09, 0xdc, 0x6d, 0x32,
	0x09, 0xd1, 0x19, 0x0c, 0xdd, 0xc0, 0x7b, 0xac, 0xbe, 0xa2, 0x3a, 0x55, 0xd0, 0x7d, 0xbd, 0x8f,
	0x5e, 0xc3, 0xb3, 0x5f, 0x2b, 0x21, 0x6b, 0xb4, 0xac, 0xc9, 0x6f, 0xf7, 0x60, 0xa3, 0x13, 0x78,
	0x6e, 0x0c, 0xdf, 0xc6, 0x45, 0xa9, 0x42, 0xa4, 0x70, 0x62, 0xc1, 0xb4, 0x28, 0xcd, 0x1f, 0x5e,
	0x02, 0x50, 0xd3, 0xdc, 0xa4, 0x1c, 0x04, 0x23, 0x2c, 0x88, 0xdf, 0x3e, 0x55, 0xaf, 0x04, 0x3c,
	0x16, 0x49, 0x2a, 0x05, 0xcd, 0x85, 0x26, 0x1e, 0xf2, 0x06, 0x88, 0xfe, 0x80, 0x53, 0x0a, 0x75,
	0x55, 0xd5, 0x5c, 0x2c, 0x72, 0x69, 0xd2, 0xbd, 0xcc, 0x92, 0x4b, 0x24, 0xcc, 0x6e, 0x3a, 0x1a,
	0x6a, 0xf9, 0x0e, 0xb5, 0x2c, 0x2d, 0x02, 0x97, 0x16, 0x3f, 0xc1, 0x97, 0xf4, 0xff, 0x5b, 0x7d,
	0xd6, 0xfe, 0xef, 0xef, 0xa3, 0x35, 0x1c, 0xd1, 0x8f, 0xd0, 0xed, 0x89, 0x83, 0xb2, 0x27, 0x33,
	0xd5, 0x97, 0x60, 0x6f, 0x5f, 0x3a, 0xdb, 0x7d, 0xb9, 0xa5, 0x29, 0xab, 0x96, 0xd8, 0x9f, 0x7a,
	0xbb, 0xca, 0xf5, 0x9d, 0x72, 0x9b, 0xa4, 0x02, 0x77, 0xc8, 0x5f, 0xc3, 0x91, 0xa1, 0xde, 0x53,
	0x5c, 0x50, 0xfc, 0x43, 0xea, 0x3d, 0x65, 0xf5, 0x8e, 0xac, 0xec, 0xbc, 0x3e, 0x79, 0x42, 0xdf,
	0xc3, 0x90, 0x8b, 0xcd, 0xca, 0x90, 0x8d, 0x4d, 0xa0, 0xab, 0xae, 0xa6, 0xbd, 0x57, 0x97, 0xd6,
	0x47, 0x17, 0xd0, 0x33, 0x9e, 0xec, 0x95, 0x9b, 0xc5, 0x2e, 0x27, 0x9d, 0xfd, 0x8f, 0x70, 0x4c,
	0x3e, 0x4d, 0xfe, 0xe7, 0x78, 0xf9, 0xe2, 0x17, 0x46, 0xd7, 0x31, 0x8f, 0xb5, 0xbb, 0xbe, 0x14,
	0xb9, 0x63, 0x11, 0x7d, 0xab, 0xa3, 0xe2, 0xf1, 0x1b, 0x43, 0xb0, 0x14, 0xb5, 0x8e, 0xb9, 0xed,
	0x84, 0x2a, 0x6c, 0xbd, 0xb1, 0x56, 0x77, 0xbd, 0x25, 0x83, 0x6e, 0x17, 0x09, 0xd1, 0x0c, 0x9e,
	0x93, 0x99, 0x61, 0x17, 0x75, 0x62, 0x0a, 0xc3, 0xd8, 0x91, 0x75, 0x6e, 0xe6, 0x2e, 0x33, 0xa6,
	0xbc, 0x65, 0x14, 0x7d, 0x80, 0x9e, 0xa5, 0xe8, 0x2e, 0xde, 0x34, 0x37, 0x86, 0xdf, 0xba, 0xad,
	0x5f, 0x40, 0x6f, 0x5e, 0xe9, 0xf7, 0x49, 0xcd, 0xc7, 0xca, 0xcd, 0x88, 0x3b, 0x2e, 0x11, 0xfe,
	0xf4, 0x00, 0xb8, 0x58, 0x88, 0x74, 0x53, 0xe2, 0x04, 0xf6, 0x90, 0x54, 0x39, 0xfa, 0x2e, 0x37,
	0x76, 0xf2, 0x00, 0x13, 0x8b, 0x9b, 0xc3, 0x10, 0x70, 0x2d, 0xb5, 0x12, 0x0b, 0xb7, 0x12, 0x6b,
	0x8a, 0xe9, 0xb6, 0xae, 0xbf, 0x7f, 0x02, 0xe8, 0xbe, 0xc9, 0xb3, 0xbb, 0xf4, 0x1e, 0x37, 0x86,
	0x75, 0x9c, 0xc5, 0xf7, 0xc2, 0x64, 0x66, 0x44, 0x3a, 0x80, 0xe2, 0xa3, 0x58, 0xe5, 0x1b, 0x61,
	0x78, 0xda, 0x00, 0xb8, 0x31, 0x3c, 0xa4, 0x59, 0x26, 0xe4, 0x15, 0x3e, 0xf0, 0x94, 0xaa, 0xcf,
	0x5d, 0x48, 0xbf, 0x89, 0x4a, 0xdd, 0xb1, 0x6f, 0xa2, 0xd2, 0x9d, 0x42, 0x7f, 0x93, 0xe7, 0x2b,
	0xa5, 0x0c, 0xd5, 0xa6, 0x60, 0x01, 0xf6, 0x0a, 0x8e, 0x6d, 0x20, 0x65, 0xa2, 0xb6, 0x8d, 0x2d,
	0x14, 0x23, 0xac, 0xd2, 0x3b, 0xf5, 0xea, 0xab, 0xc7, 0xce, 0xca, 0x98, 0xdf, 0x12, 0x89, 0xb6,
	0x90, 0xa4, 0x56, 0xfb, 0x86, 0x0b, 0x61, 0x94, 0x75, 0xfc, 0xe8, 0x1a, 0xa9, 0xbd, 0x63, 0x0b,
	0x45, 0xbb, 0x2c, 0xc7, 0x45, 0xe5, 0x5a, 0x68, 0x3b, 0xb5, 0x81, 0x6c, 0xa1, 0xec, 0x0c, 0x8e,
	0xf0, 0x1d, 0xbc, 0x16, 0xab, 0xb8, 0x76, 0x16, 0x91, 0x36, 0x88, 0x56, 0xf4, 0x24, 0xd8, 0x45,
	0x6a, 0x48, 0xa5, 0xb5, 0x41, 0x7c, 0xe7, 0xd3, 0x6c, 0x21, 0xad, 0xd1, 0x91, 0xda, 0xb6, 0x5c,
	0x6c, 0xde, 0xa5, 0xdd, 0x71, 0xfa, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1c, 0x52, 0x7b, 0x9a,
	0x48, 0x0a, 0x00, 0x00,
}
