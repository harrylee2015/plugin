package types

import "github.com/33cn/chain33/types"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerFlucky)
	// init executor type
	types.RegistorExecutor(FluckyX, NewType())
	types.RegisterDappFork(FluckyX, "Enable", 0)
}

// FluckyType 执行器基类结构体
type FluckyType struct {
	types.ExecTypeBase
}

// NewType 创建执行器类型
func NewType() *FluckyType {
	c := &FluckyType{}
	c.SetChild(c)
	return c
}

// GetPayload 获取flucky action
func (b *FluckyType) GetPayload() types.Message {
	return &FluckyAction{}
}

// GetName 获取执行器名称
func (b *FluckyType) GetName() string {
	return FluckyX
}

// GetLogMap 获取log的映射对应关系
func (b *FluckyType) GetLogMap() map[int64]*types.LogInfo {
	return logInfo
}

// GetTypeMap 根据action的name获取type
func (b *FluckyType) GetTypeMap() map[string]int32 {
	return actionName
}
