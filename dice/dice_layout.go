package dice

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

var D dwidget.ContainerStack

var die1, die2 die

var logRoll = widget.NewMultiLineEntry()
var logPlaced = widget.NewLabel("")

// Layout all Dice objects for app
func LayoutAll(d *dreams.AppObject) fyne.CanvasObject {
	// Set left label to update stats
	D.Left.Label = widget.NewLabel("")
	D.Left.SetUpdate(func() string {
		return fmt.Sprintf("Total Rolls: %d      Min Bet is %s DERO, Max Bet is %s       House: (%s DERO)  (%s dReams)", roll.total, rpc.FromAtomic(roll.min, 5), rpc.FromAtomic(roll.max, 5), roll.balance, roll.dbalance)
	})

	// Set right label to update wallet info
	D.Right.Label = widget.NewLabel("")
	D.Right.SetUpdate(dreams.SetBalanceLabelText)

	// Create dice objects
	die1, die2 = createDicePair(
		[6]*fyne.StaticResource{
			resourceDice1Png,
			resourceDice2Png,
			resourceDice3Png,
			resourceDice4Png,
			resourceDice5Png,
			resourceDice6Png},
		canvas.NewImageFromResource(resourceDice0Png))

	// Initialize chip stack images
	setDefaultChips()
	for i := 0; i < 21; i++ {
		roll.stack = append(roll.stack, canvas.NewImageFromImage(nil))
	}

	// Bet options
	opts := []string{
		"Under 7  (1:1)",
		"Over 7  (1:1)",
		"Any 7  (4:1)",
		"Any crap  (7:1)",
		"Ace deuce  (15:1)",
		"Yo  (15:1)",
		"Aces  (30:1)",
		"Midnight  (30:1)",
	}

	// Roll TX log
	logRoll.Disable()

	// Bet select
	bet := widget.NewSelect(opts, nil)

	// Currency select
	currency := widget.NewSelect([]string{"DERO", "dReams"}, nil)

	// Bet amount entry
	entry := dwidget.NewAmountEntry("", 0.1, 5)
	entry.AllowFloat = true
	entry.Validator = func(s string) (err error) {
		var f float64
		if f, err = strconv.ParseFloat(s, 64); err == nil {
			div := uint64(1)
			if currency.Selected == "dReams" {
				div = 300
			}

			u := rpc.ToAtomic(f, 5)
			if u < roll.min*div {
				return fmt.Errorf("below min bet")
			}

			if u > roll.max*div {
				return fmt.Errorf("above max bet")
			}

			return
		}

		entry.SetText("0.0")
		return
	}
	entry.SetText("0.0")

	currency.OnChanged = func(s string) {
		div := uint64(1)
		if s == "dReams" {
			div = 300
		}
		entry.SetText(rpc.FromAtomic(roll.min*div, 5))
	}

	// Label to display result
	resultsTop := canvas.NewText("", color.White)
	resultsTop.Move(fyne.NewPos(563, 470))
	resultsTop.Alignment = fyne.TextAlignCenter

	// Label to display further result info
	resultsBottom := canvas.NewText("", color.White)
	resultsBottom.Move(fyne.NewPos(563, 490))
	resultsBottom.Alignment = fyne.TextAlignCenter

	// Proposition bet button (Roll)
	btnProp := widget.NewButton("Roll", nil)
	btnProp.Importance = widget.HighImportance
	btnProp.OnTapped = func() {
		if bet.SelectedIndex() >= 0 {
			if roll.min == 0 || roll.max == 0 {
				dialog.NewInformation("Dice", "Could not get bet limits", d.Window).Show()
				return
			}

			if currency.SelectedIndex() < 0 {
				currency.FocusGained()
				info := dialog.NewInformation("Dice", "Select a currency", d.Window)
				info.SetOnClosed(currency.FocusLost)
				info.Show()
				return
			}

			div := uint64(1)
			if currency.Selected == "dReams" {
				div = 300
			}

			amt := rpc.ToAtomic(entry.Text, 5)
			if amt < roll.min*div {
				entry.FocusGained()
				info := dialog.NewInformation("Dice", fmt.Sprintf("Below minimum %s bet amount", currency.Selected), d.Window)
				info.SetOnClosed(entry.FocusLost)
				info.Show()
				return
			}

			if amt > roll.max*div {
				entry.FocusGained()
				info := dialog.NewInformation("Dice", fmt.Sprintf("Above max %s bet amount", currency.Selected), d.Window)
				info.SetOnClosed(entry.FocusLost)
				info.Show()
				return
			}

			if entry.Validate() != nil {
				entry.FocusGained()
				info := dialog.NewInformation("Dice", "Amount error", d.Window)
				info.SetOnClosed(entry.FocusLost)
				info.Show()
				return
			}

			dialog.NewConfirm("Roll", fmt.Sprintf("Confirm %s %s %s bet", rpc.FromAtomic(amt, 5), currency.Selected, propBetText(bet.SelectedIndex())), func(b bool) {
				if b {
					if tx := RollDice(amt, bet.SelectedIndex(), currency.Selected); tx != "" {
						go func() {
							D.Actions.Hide()
							roll.found = false
							resultsTop.Text = "Wait for block..."
							resultsTop.Refresh()
							resultsBottom.Text = fmt.Sprintf("Bet [%s] (%s %s)", propBetText(bet.SelectedIndex()), rpc.FromAtomic(amt, 5), currency.Selected)
							resultsBottom.Refresh()
							go menu.ShowTxDialog("Roll", "Dice", tx, 2*time.Second, d.Window)

							i := 0
							for !roll.found && i < 100 {
								go die1.roll(20, 50*time.Millisecond)
								go die2.roll(20, 50*time.Millisecond)
								getRoll(rpc.Float64Type(entry.Text), tx, currency.Selected)
								time.Sleep(time.Second)
								i++
							}

							if roll.found {
								if !d.IsWindows() {
									d.Notification("dReams - Dice", roll.rolled)
								}
								resultsTop.Text = roll.rolled
								resultsTop.Refresh()
								resultsBottom.Text = roll.result
								resultsBottom.Refresh()
								die1.land(roll.die1)
								die2.land(roll.die2)
							} else {
								resultsTop.Text = "Error, could not find TX"
								resultsTop.Refresh()
							}
							D.Actions.Show()
						}()
					}
				}
			}, d.Window).Show()
		} else {
			bet.FocusGained()
			info := dialog.NewInformation("Dice", "Select a bet", d.Window)
			info.SetOnClosed(bet.FocusLost)
			info.Show()
		}
	}

	// Set btn text to currently selected bet
	bet.OnChanged = func(s string) {
		split := strings.Split(s, "  ")
		btnProp.Text = split[0]
		btnProp.Refresh()
	}

	// Function to create table place buttons
	btnFunc := func(i int) *widget.Button {
		btn := widget.NewButton("", nil)
		btn.OnTapped = func() {
			if currency.SelectedIndex() < 0 {
				currency.FocusGained()
				info := dialog.NewInformation("Dice", "Select a currency", d.Window)
				info.SetOnClosed(currency.FocusLost)
				info.Show()
				return
			}

			text := placeBetText(i)
			lab := widget.NewLabel(fmt.Sprintf("Place on %s", text))
			lab.Alignment = fyne.TextAlignCenter

			div := float64(1)
			if currency.Selected == "dReams" {
				div = 300
			}

			ent := widget.NewEntry()
			ent.Disable()

			sli := widget.NewSlider(float64(roll.min)*div, float64(roll.max)*div)
			sli.Step = float64(roll.min)
			sli.OnChanged = func(f float64) {
				ent.SetText(fmt.Sprintf("%s %s", rpc.FromAtomic(sli.Value, 5), currency.Selected))
			}
			sli.SetValue(sli.Min)

			c := container.NewVBox(lab, ent, sli)
			dialog.NewCustomConfirm("Place", "Confirm", "Cancel", c, func(b bool) {
				if b {
					if tx := PlaceBet(uint64(sli.Value), i, currency.Selected); tx != "" {
						go func() {
							D.Actions.Hide()
							resultsTop.Text = "Wait for block..."
							resultsTop.Refresh()
							resultsBottom.Text = ""
							resultsBottom.Refresh()
							go menu.ShowTxDialog("Roll", "", tx, 2*time.Second, d.Window)
							rpc.ConfirmTx(tx, "Dice", 45)
							resultsTop.Text = fmt.Sprintf("Bet placed on [%s]", text)
							resultsTop.Refresh()
							resultsBottom.Text = ""
							resultsBottom.Refresh()
							D.Actions.Show()
						}()
					} else {
						go menu.ShowMessageDialog("Dice", "TX error, check logs", 3*time.Second, d.Window)
					}
				}
			}, d.Window).Show()
		}

		return btn
	}

	// Initialize place buttons
	var btnPlace []*widget.Button
	for i := 0; i < 7; i++ {
		btnPlace = append(btnPlace, btnFunc(i))
	}

	// Place box
	btnPlace[0].Move(fyne.NewPos(444, 273))
	btnPlace[1].Move(fyne.NewPos(531, 273))
	btnPlace[2].Move(fyne.NewPos(618, 273))
	btnPlace[3].Move(fyne.NewPos(705, 273))
	btnPlace[4].Move(fyne.NewPos(792, 273))
	btnPlace[5].Move(fyne.NewPos(879, 273))
	// Field
	btnPlace[6].Move(fyne.NewPos(123, 273))

	// Inside or Outside button OnPressed func
	inOutFunc := func(out bool) func() {
		return func() {
			if currency.SelectedIndex() < 0 {
				currency.FocusGained()
				info := dialog.NewInformation("Dice", "Select a currency", d.Window)
				info.SetOnClosed(currency.FocusLost)
				info.Show()
				return
			}

			in := float64(4)
			text := "Inside"
			details := "(5, 6, 8 and 9)"
			if out {
				details = "(4 and 10)"
				text = "Outside"
				in = 2
			}

			lab := widget.NewLabel(fmt.Sprintf("Place %s\n\n%s", text, details))
			lab.Alignment = fyne.TextAlignCenter

			ent := widget.NewEntry()
			ent.Disable()

			div := float64(1)
			if currency.Selected == "dReams" {
				div = 300
			}

			sli := widget.NewSlider((float64(roll.min)*div)*in, (float64(roll.max)*div)*in)
			sli.Step = float64(roll.min) * div
			sli.OnChanged = func(f float64) {
				ent.SetText(fmt.Sprintf("%s %s", rpc.FromAtomic(sli.Value, 5), currency.Selected))
			}
			sli.SetValue(sli.Min)

			c := container.NewVBox(lab, ent, sli)
			dialog.NewCustomConfirm("Place", "Confirm", "Cancel", c, func(b bool) {
				if b {
					if tx := InsideOutside(uint64(sli.Value), out, currency.Selected); tx != "" {
						go func() {
							D.Actions.Hide()
							roll.found = false
							resultsTop.Text = "Wait for block..."
							resultsTop.Refresh()
							resultsBottom.Text = ""
							resultsBottom.Refresh()
							go menu.ShowTxDialog(text, "", tx, 2*time.Second, d.Window)
							rpc.ConfirmTx(tx, "Dice", 45)
							resultsTop.Text = fmt.Sprintf("[%s] bet placed", text)
							resultsTop.Refresh()
							resultsBottom.Text = ""
							resultsBottom.Refresh()
							D.Actions.Show()
						}()
					} else {
						go menu.ShowMessageDialog("Dice", "TX error, check logs", 3*time.Second, d.Window)
					}
				}
			}, d.Window).Show()
		}
	}

	// Inside and outside place buttons
	bntInside := widget.NewButton("", nil)
	bntInside.Move(fyne.NewPos(123, 368))
	bntInside.Resize(fyne.NewSize(156, 45))
	bntInside.OnTapped = inOutFunc(false)

	bntOutside := widget.NewButton("", nil)
	bntOutside.Move(fyne.NewPos(284, 368))
	bntOutside.Resize(fyne.NewSize(156, 45))
	bntOutside.OnTapped = inOutFunc(true)

	table := canvas.NewImageFromResource(resourceDiceTablePng)
	table.Resize(d.GetMaxSize(1100, 600))
	table.Move(fyne.NewPos(5, 38))

	// Layout dReams container stack
	D.Back = *container.NewWithoutLayout(table)

	D.Front = *container.NewWithoutLayout(
		die1.cont,
		die2.cont,
		resultsTop,
		resultsBottom)
	for _, s := range roll.stack {
		D.Front.Add(s)
	}
	for i, s := range btnPlace {
		if i == 6 {
			s.Resize(fyne.NewSize(316, 91))
		} else {
			s.Resize(fyne.NewSize(82, 91))
		}

		D.Front.Add(s)
	}
	D.Front.Add(bntInside)
	D.Front.Add(bntOutside)

	// Search object
	var searched string
	search_entry := widget.NewEntry()
	search_entry.SetPlaceHolder("TXID:")
	search_button := widget.NewButtonWithIcon("", dreams.FyneIcon("search"), func() {
		txid := search_entry.Text
		if len(txid) == 64 {
			if txid != searched {
				searched = txid
				D.Actions.Hide()
				resultsTop.Text = "Searching..."
				resultsTop.Refresh()
				resultsBottom.Text = ""
				resultsBottom.Refresh()
				roll.found = false
				getRoll(0, search_entry.Text, "")
				if roll.found {
					resultsTop.Text = roll.rolled
					resultsTop.Refresh()
					resultsBottom.Text = roll.result
					resultsBottom.Refresh()
					die1.land(roll.die1)
					die2.land(roll.die2)
				} else {
					searched = ""
					dialog.NewInformation("Search", "No results found", d.Window).Show()
					resultsTop.Text = ""
					resultsTop.Refresh()
					resultsBottom.Text = ""
					resultsBottom.Refresh()
				}

				D.Actions.Show()
			}

			return
		}

		dialog.NewInformation("Search", "Not a valid TXID", d.Window).Show()
	})

	// How to play help button
	btnHelp := widget.NewButton("How to play", func() { layoutHelp(d) })
	btnHelp.Importance = widget.LowImportance

	// Game odds help button
	btnOdds := widget.NewButton("Odds", func() { dialog.NewCustom("Odds", "Done", layoutOdds(), d.Window).Show() })

	// Clear placed bet button
	btnClear := widget.NewButton("Clear Bets", nil)
	btnClear.Importance = widget.HighImportance
	btnClear.OnTapped = func() {
		found, bets := getBets()
		if !found {
			dialog.NewInformation("Dice", "You don't have any bets on the table", d.Window).Show()
			return
		}

		dialog.NewConfirm("Clear Bets", "Clear your bets off the table?", func(b bool) {
			if b {
				if tx := Clear(bets[0]); tx != "" {
					go func() {
						D.Actions.Hide()
						resultsTop.Text = "Wait for block..."
						resultsTop.Refresh()
						resultsBottom.Text = ""
						resultsBottom.Refresh()
						go menu.ShowTxDialog("Roll", "", tx, 2*time.Second, d.Window)
						rpc.ConfirmTx(tx, "Dice", 45)
						resultsTop.Text = "Cleared your placed bets"
						resultsTop.Refresh()
						resultsBottom.Text = ""
						resultsBottom.Refresh()
						D.Actions.Show()
					}()
				} else {
					go menu.ShowMessageDialog("Dice", "TX error, check logs", 3*time.Second, d.Window)
				}
			}
		}, d.Window).Show()
	}

	D.Actions = *container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(
			container.NewVBox(layout.NewSpacer(), logPlaced, container.NewHBox(btnClear)),
			layout.NewSpacer(),
			container.NewVBox(layout.NewSpacer(), btnOdds),
			layout.NewSpacer(),
			container.NewVBox(layout.NewSpacer(), btnHelp),
			layout.NewSpacer(),
			container.NewBorder(
				nil,
				container.NewHBox(
					layout.NewSpacer(),
					container.NewVBox(layout.NewSpacer(), container.NewBorder(nil, nil, nil, search_button, container.NewStack(dwidget.NewSpacer(560, 0), search_entry))),
					container.NewVBox(layout.NewSpacer(), bet, currency, entry, btnProp)),
				nil,
				nil),
			dwidget.NewSpacer(250, 0)))

	D.Actions.Hide()

	D.DApp = container.NewStack(
		container.NewVBox(
			dwidget.LabelColor(container.NewHBox(D.Left.Label, layout.NewSpacer(), D.Right.Label))),
		&D.Back,
		&D.Front,
		container.NewAdaptiveGrid(3,
			layout.NewSpacer(),
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), container.NewStack(container.NewBorder(dwidget.NewSpacer(250, 30), nil, nil, nil, logRoll)))),
		&D.Actions)

	// Main process routine
	go func() {
		var synced bool
		time.Sleep(3 * time.Second)
		for {
			select {
			case <-d.Receive():
				if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
					roll.min = 0
					D.Actions.Hide()
					resultsTop.Text = ""
					resultsTop.Refresh()
					resultsBottom.Text = ""
					resultsBottom.Refresh()
					synced = false
					d.WorkDone()
					continue
				}

				if !synced && gnomes.Scan(d.IsConfiguring()) {
					logger.Println("[Dice] Syncing")
					getLastRoll(&die1, &die2)
					synced = true
					D.Actions.Show()
				}

				getStats()
				d.WorkDone()
			case <-d.CloseDapp():
				logger.Println("[Dice] Done")
				return
			}
		}
	}()

	return D.DApp
}

// Layout game odds objects
func layoutOdds() fyne.CanvasObject {
	imgs := []*canvas.Image{
		canvas.NewImageFromResource(resourceDiceAcesPng),
		canvas.NewImageFromResource(resourceDiceAceDeucePng),
		canvas.NewImageFromResource(resourceDiceYoPng),
		canvas.NewImageFromResource(resourceDiceMidnightPng),
	}

	// Proposition bets
	propNames := []string{
		"Under 7 (1:1)",
		"Over 7 (1:1)",
		"Any 7 (4:1)",
		"Any crap (7:1)",
		"Aces (30:1)",
		"Ace Deuce (15:1)",
		"Yo (15:1)",
		"Midnight (30:1)",
	}

	propForm := []*widget.FormItem{}
	propForm = append(propForm, widget.NewFormItem("One Roll Bets", widget.NewLabel("(Proposition Bet)")))

	for i, n := range propNames {
		propForm = append(propForm, widget.NewFormItem(n, dwidget.NewSpacer(135, 50)))
		if i < 4 {
			switch n {
			case "Under 7 (1:1)":
				propForm[i+1].Widget = widget.NewLabel("Roll under 7")
			case "Over 7 (1:1)":
				propForm[i+1].Widget = widget.NewLabel("Roll over 7")
			case "Any 7 (4:1)":
				propForm[i+1].Widget = widget.NewLabel("Roll any 7")
			case "Any crap (7:1)":
				propForm[i+1].Widget = widget.NewLabel("Roll a 2, 3 or 12")
			}
		} else {
			imgs[i-4].SetMinSize(fyne.NewSize(135, 50))
			propForm[i+1].Widget = container.NewCenter(imgs[i-4])
		}
	}

	// Place bets
	placeForm := []*widget.FormItem{}
	placeForm = append(placeForm, widget.NewFormItem("Multi Roll Bets", widget.NewLabel("(Place Bet)")))
	placeForm = append(placeForm, widget.NewFormItem("Place 4 (9:5)", widget.NewLabel("Roll a 4 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Place 5 (7:5)", widget.NewLabel("Roll a 5 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Place 6 (7:6)", widget.NewLabel("Roll a 6 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Place 8 (7:6)", widget.NewLabel("Roll a 8 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Place 9 (7:5)", widget.NewLabel("Roll a 9 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Place 10 (9:5)", widget.NewLabel("Roll a 10 before 7")))
	placeForm = append(placeForm, widget.NewFormItem("Inside (Place pay)", widget.NewLabel("Place on 5, 6, 8 and 9")))
	placeForm = append(placeForm, widget.NewFormItem("Outside (Place pay)", widget.NewLabel("Place on 4 and 10")))
	placeForm = append(placeForm, widget.NewFormItem("Field", widget.NewLabel("Roll a 3, 4, 9, 10 or 11 for (1:1)\nRoll a 2 or 12 for (2:1)\nRolling a 5, 6, 7, 8 loses")))

	return container.NewHBox(
		container.NewVBox(widget.NewForm(propForm...)),
		container.NewVBox(widget.NewForm(placeForm...)))
}

// Layout how to play dialog
func layoutHelp(d *dreams.AppObject) {
	text := `dDice7 is multiplayer dice game, similar to craps

	There are two types of bets in dDice7

	- Proposition bet 
	(one roll bet, where you win or lose on that single roll) 

	- Place bet 
	(multi roll bet, where the bet stays on the table and can win until a specific loosing roll occurs)

	All players of the game play at the same table

	Players can either make place bets and wait for others to roll, or roll themselves
	- Click on the table to make a place bet 
	(4, 5, Field, Inside ect...)

	To roll, select a proposition bet from the drop down and enter amount

	The outcome of each roll transaction is valid to all place bets on the table

	Chip stacks on the table show the cumulative funds on that bet currently

	Players can clear their individual place bets from the table at any time

	Click the Odds button to view odds and roll details
	`

	dialog.NewInformation("How to Play", text, d.Window).Show()
}
