package dice

import (
	"fyne.io/fyne/v2/canvas"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/cryptography/crypto"
	dero "github.com/deroproject/derohe/rpc"
)

type dice7 struct {
	found    bool
	rolled   string
	result   string
	balance  string
	dbalance string
	die1     int
	die2     int
	last     int64
	total    uint64
	min      uint64
	max      uint64
	stack    []*canvas.Image
}

var roll dice7

// Roll the dice with proposition bet
func RollDice(amt uint64, b int, token string) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Roll"},
		dero.Argument{Name: "bet", DataType: "U", Value: b},
	}

	scid := crypto.ZEROHASH
	if token == "dReams" {
		scid = crypto.HashHexToHash(rpc.DreamsSCID)
	}

	t1 := dero.Transfer{
		SCID:        scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(DICESCID, "[Dice]", args, t, rpc.HighLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     DICESCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	client, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()
	if err := client.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Dice] Roll: %s", err)
		return
	}

	// if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
	// 	rpc.PrintError("[Baccarat] Bet: %s", err)
	// 	return
	// }

	rpc.PrintLog("[Dice] Roll TX: %s", txid)

	return txid.TXID
}

// Place a single place bet
func PlaceBet(amt uint64, p int, token string) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Place"},
		dero.Argument{Name: "p", DataType: "U", Value: p},
	}

	scid := crypto.ZEROHASH
	if token == "dReams" {
		scid = crypto.HashHexToHash(rpc.DreamsSCID)
	}

	t1 := dero.Transfer{
		SCID:        scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(DICESCID, "[Dice]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     DICESCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	client, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()
	if err := client.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Dice] Place: %s", err)
		return
	}

	// if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
	// 	rpc.PrintError("[Baccarat] Bet: %s", err)
	// 	return
	// }

	rpc.PrintLog("[Dice] Place TX: %s", txid)

	return txid.TXID
}

// Place a inside or outside bet
func InsideOutside(amt uint64, out bool, token string) (tx string) {
	entrypoint := "Inside"
	if out {
		entrypoint = "Outside"
	}

	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: entrypoint}}

	scid := crypto.ZEROHASH
	if token == "dReams" {
		scid = crypto.HashHexToHash(rpc.DreamsSCID)
	}

	t1 := dero.Transfer{
		SCID:        scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(DICESCID, "[Dice]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     DICESCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	client, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()
	if err := client.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Dice] %s: %s", entrypoint, err)
		return
	}

	// if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
	// 	rpc.PrintError("[Baccarat] Bet: %s", err)
	// 	return
	// }

	rpc.PrintLog("[Dice] %s TX: %s", entrypoint, txid)

	return txid.TXID
}

// Withdraws any current place bets a wallet has
func Clear(b uint64) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Clear"},
		dero.Argument{Name: "b", DataType: "U", Value: b},
	}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(DICESCID, "[Dice]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     DICESCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	client, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()
	if err := client.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Dice] Clear: %s", err)
		return
	}

	// if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
	// 	rpc.PrintError("[Baccarat] Bet: %s", err)
	// 	return
	// }

	rpc.PrintLog("[Dice] Clear TX: %s", txid)

	return txid.TXID
}
