package dice

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/canvas"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
)

var gnomon = gnomes.NewGnomes()

// Scan for a roll result bt TXID
func getRoll(amt float64, tx string, curr string) {
	if gnomon.IsReady() && tx != "" {
		if _, rolls := gnomon.GetSCIDValuesByKey(DICESCID, "rolls"); rolls != nil {
			if _, display := gnomon.GetSCIDValuesByKey(DICESCID, "display"); display != nil {
				disp := display[0]
				if roll.total < display[0] {
					disp = roll.total
				}

				h := float64(1)
				if _, house := gnomon.GetSCIDValuesByKey(DICESCID, "display"); house != nil {
					h = (10000 - rpc.Float64Type(house[0])) / 10000
				}

				for i := roll.total; i > roll.total-disp; i-- {
					pre := "roll"
					index := strconv.Itoa(int(i))
					if tx_string, _ := gnomon.GetSCIDValuesByKey(DICESCID, pre+index); tx_string != nil {
						res := strings.Split(tx_string[0], "_")
						if len(res) < 4 {
							logger.Errorln("[Dice] wrong len")
							continue
						}

						txid := res[0]

						if txid == tx {
							roll.found = true
							win := "Lost"
							bet := rpc.IntType(res[1])
							die1 := rpc.IntType(res[2])
							die2 := rpc.IntType(res[3])
							roll.die1 = die1 - 1
							roll.die2 = die2 - 1
							rolled := die1 + die2

							if amt > 0 {
								if rolled == 7 {
									// Any seven
									if bet == 2 {
										win = fmt.Sprintf("Win %.5f %s", (amt*4)*h, curr)
									}
								} else if rolled > 7 {
									// Over
									if bet == 1 {
										win = fmt.Sprintf("Win %.5f %s", (amt*1)*h, curr)
									}

									// Yo
									if rolled == 11 && bet == 5 {
										win = fmt.Sprintf("Win %.5f %s", (amt*15)*h, curr)
									}

									// Midnight or Any crap
									if rolled == 12 && (bet == 7 || bet == 3) {
										win = fmt.Sprintf("Win %.5f %s", (amt*30)*h, curr)
									}
								} else {
									// Under
									if bet == 0 {
										win = fmt.Sprintf("Win %.5f %s", (amt*1)*h, curr)
									}

									// Aces or Any crap
									if rolled == 2 && (bet == 6 || bet == 3) {
										win = fmt.Sprintf("Win %.5f %s", (amt*30)*h, curr)
									}

									if rolled == 3 {
										switch bet {
										case 3:
											// Any crap
											win = fmt.Sprintf("Win %.5f %s", (amt*7)*h, curr)
										case 4:
											// Ace deuce
											win = fmt.Sprintf("Win %.5f %s", (amt*15)*h, curr)
										default:
											// nothing
										}
									}
								}
								roll.result = fmt.Sprintf("Bet [%s] (%s)", propBetText(bet), win)
							} else {
								roll.result = ""
							}

							roll.rolled = fmt.Sprintf("Roll# %s  -  [%d] (%d & %d)", index, rolled, die1, die2)

							return
						}
					}
				}
			}
		}
	}
}

// Get current table stats
func getStats() {
	if gnomon.IsReady() {
		var display uint64
		_, dis := gnomon.GetSCIDValuesByKey(DICESCID, "display")
		if dis != nil {
			display = dis[0]
		}

		_, rolls := gnomon.GetSCIDValuesByKey(DICESCID, "rolls")
		if rolls != nil {
			roll.total = rolls[0]
		}

		_, bal := gnomon.GetSCIDValuesByKey(DICESCID, "bal")
		if bal != nil {
			roll.balance = rpc.FromAtomic(bal[0], 5)
		}

		_, dbal := gnomon.GetSCIDValuesByKey(DICESCID, "balD")
		if dbal != nil {
			roll.dbalance = rpc.FromAtomic(dbal[0], 5)
		}

		_, min := gnomon.GetSCIDValuesByKey(DICESCID, "min")
		if min != nil {
			roll.min = min[0]
		}

		_, max := gnomon.GetSCIDValuesByKey(DICESCID, "max")
		if max != nil {
			roll.max = max[0]
		}

		D.Left.UpdateText()
		D.Right.UpdateText()

		placeChipStack()

		// Update TX roll log
		if gnomon.GetLastHeight() > roll.last {
			roll.last = gnomon.GetLastHeight()

			disp := display
			if roll.total < display {
				disp = roll.total
			}

			var results string
			for i := roll.total; i > roll.total-disp; i-- {
				if r, _ := gnomon.GetSCIDValuesByKey(DICESCID, "roll"+strconv.Itoa(int(i))); r != nil {
					split := strings.Split(r[0], "_")
					if len(split) > 3 {
						die1 := rpc.IntType(split[2])
						die2 := rpc.IntType(split[3])
						num := die1 + die2
						results = results + fmt.Sprintf("# %d - [%d]  (%d & %d)\n", i, num, die1, die2)
					}
				}
			}

			logRoll.SetText(results)
		}
	}
}

// Find if wallet has any active place bets and return those bets as []uint64
func getBets() (found bool, bets []uint64) {
	for i := uint64(0); i <= 30; i++ {
		if _, u := gnomon.GetSCIDValuesByKey(DICESCID, i); u != nil {
			if addr, _ := gnomon.GetSCIDValuesByKey(DICESCID, fmt.Sprintf("b_%d", i)); addr != nil {
				if addr[0] == rpc.Wallet.Address {
					found = true
					bets = append(bets, i)
				}
			}
		}
	}

	return
}

// Find the last roll that occurred and land the dice on it
func getLastRoll(d1, d2 *die) {
	if _, rolls := gnomon.GetSCIDValuesByKey(DICESCID, "rolls"); rolls != nil {
		if r, _ := gnomon.GetSCIDValuesByKey(DICESCID, "roll"+strconv.Itoa(int(rolls[0]))); r != nil {
			split := strings.Split(r[0], "_")
			if len(split) > 3 {
				n := int(rolls[0])
				die1 := rpc.IntType(split[2])
				die2 := rpc.IntType(split[3])
				num := die1 + die2
				d1.land(die1 - 1)
				d2.land(die2 - 1)
				D.Front.Objects[2].(*canvas.Text).Text = fmt.Sprintf("Roll# %d - [%d]  (%d & %d)", n, num, die1, die2)
			}
		}
	}
}
