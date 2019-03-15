/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package executor

import "fmt"

func calcGold5GBuyRound(round int64, addr string, index int64) []byte {
	key := fmt.Sprintf("LODB-gold5g-buy:%010d:%s:%018d", round, addr, index)
	return []byte(key)
}

func calcGold5GBuyPrefix(round int64, addr string) []byte {
	key := fmt.Sprintf("LODB-gold5g-buy:%010d:%s:", round, addr)
	return []byte(key)
}

func calcGold5GAddrRound(round int64, addr string) []byte {
	key := fmt.Sprintf("LODB-gold5g-addrInfos:%010d:%s", round, addr)
	return []byte(key)
}

func calcGold5GAddrPrefix(round int64) []byte {
	key := fmt.Sprintf("LODB-gold5g-addrInfos:%010d:", round)
	return []byte(key)
}

func calcGold5GStartRound(round int64) []byte {
	key := fmt.Sprintf("LODB-gold5g-start:%010d", round)
	return []byte(key)
}

func calcGold5GStartPrefix() []byte {
	key := fmt.Sprintf("LODB-gold5g-start:")
	return []byte(key)
}

func calcGold5GDrawRound(round int64) []byte {
	key := fmt.Sprintf("LODB-gold5g-draw:%010d", round)
	return []byte(key)
}

func calcGold5GDrawPrefix() []byte {
	key := fmt.Sprintf("LODB-gold5g-draw:")
	return []byte(key)
}
