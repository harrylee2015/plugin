package types

import "fmt"

const (
	validRangeNum = 2

	minRangeNumIndex = 0
	maxRangeNumIndex = 1
)

var (
	// 平台地址
	platFormAddr = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

	// 奖池最大额度
	maxBonus = float32(100000)
	// 奖池最小额度
	minBonus = float32(500)

	// 超过最大额度，奖池转出资金
	bonusToPlatform = float32(50000)
	// 小于最小额度，奖池转入资金
	platformToBonus = float32(500)

	// 模数
	modNum = int64(10000)

	rewardRule  RewardRule
	betTimeRule BetTimesRule
)

type RewardRule struct {
	rules []*Rule
}

type BetTimesRule struct {
	rules []*BetTime
}

type RangeInfo struct {
	minRange int64
	maxRange int64
}

func SetConfig(config *FluckyCfg) {
	if config.GetPlatformAddr() != "" {
		platFormAddr = config.GetPlatformAddr()
	}

	if config.GetMaxBonus() > 0 && config.GetMinBonus() >= 0 && config.GetMaxBonus() > config.GetMinBonus() {
		maxBonus = config.GetMaxBonus()
		minBonus = config.GetMinBonus()
	}

	if config.GetBonusToPlatform() >= 0 && config.GetBonusToPlatform() <= config.GetMaxBonus() &&
		config.GetMaxBonus()-config.GetBonusToPlatform() > config.GetMinBonus() {
		bonusToPlatform = config.GetBonusToPlatform()
	}

	if config.GetPlatformToBonus() >= 0 && config.GetMinBonus()+config.GetPlatformToBonus() < config.GetMaxBonus() {
		platformToBonus = config.GetPlatformToBonus()
	}

	if config.GetModNum() > 0 {
		modNum = config.GetModNum()
	}

	checkRewardRule(config)
	checkBetTImeRule(config)
}

func checkRewardRule(config *FluckyCfg) {
	if config.GetRewardRule() == nil {
		useDefaultRewardRule()
		return
	}
	for _, rule := range config.GetRewardRule() {
		if !bValidRule(rule) {
			useDefaultRewardRule()
			return
		}
		rewardRule.rules = append(rewardRule.rules, rule)
	}

	if err := CheckRangeInfo(rewardRule.rules); err != nil {
		useDefaultRewardRule()
		return
	}
}

func checkBetTImeRule(config *FluckyCfg) {
	if config.GetBetTimeRule() == nil {
		useDefaultBetTimeRule()
		return
	}
	for _, rule := range config.GetBetTimeRule() {
		if !bValidRule(rule) {
			useDefaultBetTimeRule()
			return
		}
		betTimeRule.rules = append(betTimeRule.rules, rule)
	}
}

func useDefaultRewardRule() {
	rewardRule.rules = make([]*Rule, 0)
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{0, 5000}, Type: 0, Percent: 0.1})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{5001, 8000}, Type: 0, Percent: 0.5})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{8001, 9000}, Type: 0, Percent: 1})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9001, 9500}, Type: 1, Percent: 0.005})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9501, 9900}, Type: 1, Percent: 0.01})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9901, 9950}, Type: 1, Percent: 0.05})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9951, 9990}, Type: 1, Percent: 0.1})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9991, 9998}, Type: 1, Percent: 0.2})
	rewardRule.rules = append(rewardRule.rules, &Rule{Range: []int64{9999}, Type: 1, Percent: 0.5})
}

func useDefaultBetTimeRule() {
	betTimeRule.rules = make([]*BetTime, 0)
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 1, Times: 1})
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 2, Times: 3})
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 3, Times: 5})
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 4, Times: 8})
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 5, Times: 12})
	betTimeRule.rules = append(betTimeRule.rules, &BetTime{Amount: 10, Times: 30})

}

func GetPlatFormAddr() string {
	return platFormAddr
}

func GetMaxBonus() float32 {
	return maxBonus
}

func GetMinBonus() float32 {
	return minBonus
}

func GetBonusToPlatform() float32 {
	return bonusToPlatform
}

func GetPlatformToBonus() float32 {
	return platformToBonus
}

func GetModeNum() int64 {
	return modNum
}

func GetRewardRule(num int64) *Rule {
	for _, rule := range rewardRule.rules {
		if num >= getMinRangeNum(rule) && num <= getMaxRangeNum(rule) {
			return rule
		}
	}
	return nil
}

func GetBetTimesRule(amount int64) *BetTime {
	for _, rule := range betTimeRule.rules {
		if amount == rule.GetAmount() {
			return rule
		}
	}
	return nil
}

func bValidRule(rule interface{}) (bRet bool) {
	if rewardRule, ok := rule.(*Rule); ok {
		if bValidRewardRule(rewardRule) {
			bRet = true
		}
	} else if betTimesRule, ok := rule.(*BetTime); ok {
		if bValidBetTimesRule(betTimesRule) {
			bRet = true
		}
	}
	return
}

func bValidRewardRule(rule *Rule) (bRet bool) {
	if bValidPercent(rule.GetPercent()) && bValidRange(rule.GetRange()) && bValidType(rule.GetType()) {
		bRet = true
	}
	return
}

func bValidBetTimesRule(rule *BetTime) (bRet bool) {
	if bValidNum(rule.GetAmount()) && bValidNum(rule.GetTimes()) {
		bRet = true
	}
	return
}

func bValidRange(scope []int64) (bRet bool) {
	if (len(scope) == validRangeNum && scope[minRangeNumIndex] <= scope[maxRangeNumIndex]) || (len(scope) == 1 && scope[minRangeNumIndex] > 0) {
		bRet = true
	}
	return
}

func bValidType(t int32) (bRet bool) {
	if t == 0 || t == 1 {
		bRet = true
	}
	return
}

func bValidPercent(percent float32) (bRet bool) {
	if percent > 0 && percent <= 1 {
		bRet = true
	}
	return
}

func bValidNum(amount int64) (bRet bool) {
	if amount > 0 {
		bRet = true
	}
	return
}

func CheckRangeInfo(rules []*Rule) (err error) {
	var orderRangeInfo []*RangeInfo
	for _, rule := range rules {
		var info RangeInfo
		info.minRange = rule.Range[minRangeNumIndex]
		if len(rule.GetRange()) == 1 {
			info.maxRange = rule.Range[minRangeNumIndex]
		} else {
			info.maxRange = rule.Range[maxRangeNumIndex]
		}
		orderRangeInfo = append(orderRangeInfo, &info)
	}

	sort(orderRangeInfo)

	for i := 0; i < len(orderRangeInfo)-1; i++ {
		if orderRangeInfo[i].maxRange != orderRangeInfo[i+1].minRange-1 {
			return fmt.Errorf("Wrong Range Info...")
		}
	}

	if orderRangeInfo[0].minRange != 0 || orderRangeInfo[len(orderRangeInfo)-1].maxRange != GetModeNum()-1 {
		return fmt.Errorf("Wrong Range Info...")
	}
	return nil
}

func sort(rangeInfo []*RangeInfo) {
	// 数据量不大，冒泡就可以
	len := len(rangeInfo)
	for i := 0; i < len-1; i++ {
		for j := i + 1; j < len; j++ {
			if rangeInfo[i].minRange > rangeInfo[j].minRange {
				var info RangeInfo
				info = *rangeInfo[j]
				rangeInfo[j] = rangeInfo[i]
				rangeInfo[i] = &info
			}
		}
	}
	return
}

func getMinRangeNum(rule *Rule) int64 {
	return rule.Range[minRangeNumIndex]
}

func getMaxRangeNum(rule *Rule) int64 {
	if len(rule.GetRange()) == 1 {
		return rule.Range[minRangeNumIndex]
	}
	return rule.Range[maxRangeNumIndex]
}
