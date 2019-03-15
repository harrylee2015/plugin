/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package types

import "fmt"

// some errors definition
var (
	ErrGold5GStartRound       = fmt.Errorf("%s", "There's still one round left,you cann't start next round Gold5G!")
	ErrGold5GManageAddr       = fmt.Errorf("%s", "You don't have permission to start Gold5G game.")
	ErrGold5GManageBuyTicket  = fmt.Errorf("%s", "You are manager,you don't have permission to buy ticket")
	ErrGold5GBuyTicket        = fmt.Errorf("%s", "the Gold5G is not start a new round!")
	ErrGold5GBuyTicketTimeOut = fmt.Errorf("%s", "The rest of the time is over,you can't buy any more tickets!")
	ErrGold5GDrawRound        = fmt.Errorf("%s", "There's not Gold5G round to draw!")
	ErrGold5GDrawRemainTime   = fmt.Errorf("%s", "There is time reamining,you can't draw the round game!")
	ErrGold5GDrawRepeat       = fmt.Errorf("%s", "You can't repeat draw!")
)
