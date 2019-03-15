/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package types

var (

	// 本游戏合约管理员地址（热钱包地址）
	gold5gManagerAddr = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

	// 本游戏合约平台开发者分成地址（冷钱包地址，由项目方提供）
	gold5gDeveloperAddr = "1Ft6SUwdnhzryBwPDEnchLeNHgvSAXPbAK"

	// 本游戏平台合伙人分红地址（冷钱包地址，由项目方提供）
	gold5gPartnerAddr = "1HHy3xZEtEYHjPJYkvWvR25CDPZNLkwZQ4"

	// 本游戏一等奖奖金锁定地址（冷钱包地址，由项目方提供）
	gold5gFirstPrizeAddr = "1GTqLXA8cjUKWDtXjax9Acs2itwbAQierf"

	// 本游戏二等奖奖金锁定地址（冷钱包地址，由项目方提供）
	gold5gSecondaryPrizeAddr = "1NuPawKogkafCncNVmpHmg3qEvCx8GcNWj"

	// 本游戏推广奖金锁定地址(热钱包地址，每轮结束后，由后台统计，系统进行统一派发）
	gold5gPromotionAddr = "19MJmA7GcE1NfMwdGqgLJioBjVbzQnVYvR"

	// 本游戏合约护盘基金地址（冷钱包地址，由项目方提供）
	gold5gSupportFundAddr = "1CtY2o6AAJv8ae7B4ptADRVnkALvPERYoJ"

	// unassigned地址（冷钱包地址，由项目方提供）
	gold5gUnassignedAddr = "1GaPfaU8dSAdipzUtWevZq8fMwvKWHnVFc"

	//一等奖（倒数100名分成百分比）
	gold5gBonusFirstPrize = float64(0.3)

	//二等奖 (倒数101-5000名参与者分红百分比)
	gold5gBonusSecondaryPrize = float64(0.5)

	// 滚动到下期奖金池百分比
	gold5gBonusPool = float64(0.2)

	// 护盘基金
	gold5gBonusSupportFund = float64(1)

	// 合伙人分红
	gold5gBonusPartner = float64(2)

	// 平台运营及开发者费用
	gold5gBonusDeveloper = float64(3)

	// 奖金池
	gold5gBonusBonusPool = float64(4)

	// 一级推荐+二级+见点奖（一层是0.35，根据会员等级分5层，10层，20层）
	gold5gBonusPromotion = float64(10)

	// 本游戏一轮运行的最长周期（单位：秒）
	gold5gTimeLife = int64(86400)

	// 购买一票延长时间（单位：秒）
	gold5gTimeTicket = int64(60)

	//ticket价格
	gold5gTicketPrice = float64(100)

	tokenSymbol = ""
)

func SetConfig(config *Config) {
	//// manager 地址
	//managerAddr := config.GetManager()
	//if validAddr(managerAddr) {
	//	f3dManagerAddr = managerAddr
	//}
	//
	//// developer 地址
	//developerAddr := config.GetDeveloper()
	//if validAddr(developerAddr) {
	//	f3dDeveloperAddr = developerAddr
	//}
	//
	//// 赢家获取的奖金百分比
	//winnerBonus := config.GetWinnerBonus()
	//// 用户持有key分红百分比
	//keyBonus := config.GetKeyBonus()
	//// 滚动到下期奖金池百分比
	//poolBonus := config.GetPoolBonus()
	//// 平台运营及开发者费用百分比
	//developBonus := config.GetDeveloperBonus()
	//
	//if validSum(winnerBonus, keyBonus, poolBonus, developBonus) {
	//	f3dBonusWinner = winnerBonus
	//	f3dBonusKey = keyBonus
	//	f3dBonusPool = poolBonus
	//	f3dBonusDeveloper = developBonus
	//}
	//
	//// 本游戏一轮运行的最长周期（单位：秒）
	//lifeTime := config.GetLifeTime()
	//if validTime(lifeTime) {
	//	f3dTimeLife = lifeTime
	//}
	//
	//// 一把钥匙延长的游戏时间（单位：秒）
	//keyTime := config.GetKeyIncrTime()
	//if validTime(keyTime) {
	//	f3dTimeKey = keyTime
	//}
	//
	//// 一次购买钥匙最多延长的游戏时间（单位：秒）
	//keyMaxTime := config.GetMaxkeyIncrTime()
	//if validTime(keyMaxTime) {
	//	f3dTimeMaxkey = keyMaxTime
	//}
	//
	//// 开奖延迟时间
	//delayTime := config.GetDrawDelayTime()
	//if validTime(delayTime) {
	//	f3dTimeDeplay = delayTime
	//}
	//
	//// 钥匙涨价幅度（下一个人购买钥匙时在上一把钥匙基础上浮动幅度百分比），范围1-100
	//keyPriceIncr := config.GetIncrKeyPrice()
	//if validPercent(keyPriceIncr) {
	//	f3dKeyPriceIncr = keyPriceIncr
	//}
	//
	//// start Key price  0.1 token
	//keyStartPrice := config.GetStartKeyPrice()
	//if keyStartPrice > 0 {
	//	f3dKeyPriceStart = keyStartPrice
	//}
}
func GetTokenSymbol() string {
	return tokenSymbol
}

func Getgold5gManagerAddr() string {
	return gold5gManagerAddr
}

func Getgold5gDeveloperAddr() string {
	return gold5gDeveloperAddr
}

func Getgold5gPartnerAddr() string {
	return gold5gPartnerAddr
}

func Getgold5gFirstPrizeAddr() string {
	return gold5gFirstPrizeAddr
}

func Getgold5gSecondaryPrizeAddr() string {
	return gold5gSecondaryPrizeAddr
}
func Getgold5gUnassignedAddr() string {
	return gold5gUnassignedAddr
}

func Getgold5gPromotionAddr() string {
	return gold5gPromotionAddr
}

func Getgold5gSupportFundAddr() string {
	return gold5gSupportFundAddr
}

func Getgold5gBonusFirstPrize() float64 {
	return gold5gBonusFirstPrize
}
func Getgold5gBonusSecondaryPrize() float64 {
	return gold5gBonusSecondaryPrize
}
func Getgold5gBonusPool() float64 {
	return gold5gBonusPool
}
func Getgold5gBonusSupportFund() float64 {
	return gold5gBonusSupportFund
}
func Getgold5gBonusPartner() float64 {
	return gold5gBonusPartner
}
func Getgold5gBonusDeveloper() float64 {
	return gold5gBonusDeveloper
}
func Getgold5gBonusBonusPool() float64 {
	return gold5gBonusBonusPool
}

func Getgold5gBonusPromotion() float64 {
	return gold5gBonusPromotion
}

func Getgold5gTimeLife() int64 {
	return gold5gTimeLife
}

func Getgold5gTimeTicket() int64 {
	return gold5gTimeTicket
}

func Getgold5gTicketPrice() float64 {
	return gold5gTicketPrice
}

//
//func validAddr(addr string) bool {
//	if addr != "" && len(addr) == 64 {
//		return true
//	}
//	return false
//}
//
//func validPercent(percent float32) bool {
//	if percent > 0 && percent < 1 {
//		return true
//	}
//	return false
//}
//
//func validTime(time int64) bool {
//	if time > 0 {
//		return true
//	}
//	return false
//}
//
//func validSum(vals ...float32) bool {
//	sum := float32(0)
//	for _, val := range vals {
//		if validPercent(val) {
//			sum += val
//		} else {
//			return false
//		}
//	}
//
//	if sum == 1 {
//		return true
//	} else {
//		return false
//	}
//}
