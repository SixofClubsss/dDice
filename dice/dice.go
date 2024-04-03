package dice

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/blang/semver/v4"
	"github.com/civilware/Gnomon/structures"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
)

type die struct {
	pip  [6]*canvas.Image
	back *canvas.Image
	cont *fyne.Container
	sync.RWMutex
}

type settings struct {
	dice dreams.AssetSelect
}

var DICESCID = "fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd"

var logger = structures.Logger.WithFields(logrus.Fields{})

var version = semver.MustParse("0.0.0-dev.1")

var chipStack map[uint64]*fyne.StaticResource
var Settings settings

var pathDice = filepath.Join(dreams.GetDir(), "datashards", "assets", "dice")

// Get current dice package version
func Version() semver.Version {
	return version
}

func DreamsMenuIntro() (entries map[string][]string) {
	entries = map[string][]string{
		"dDice": {
			"Decentralized dice games",
			"dDice7"},

		"dDice7": {
			"Similar to craps, dDice7 is a multiplayer game with a variety of betting options",
			"Accepts dReams and DERO",
			"Click the 'How to play' button in the Dice tab for detailed game instructions"},
	}

	return
}

// Return prop bet string of i
func propBetText(i int) string {
	switch i {
	case 0:
		return "Under 7"
	case 1:
		return "Over 7"
	case 2:
		return "Any 7"
	case 3:
		return "Any crap"
	case 4:
		return "Ace deuce"
	case 5:
		return "Yo"
	case 6:
		return "Aces"
	case 7:
		return "Midnight"
	default:
		return "Error"
	}
}

// Return place bet string of i
func placeBetText(i int) string {
	switch i {
	case 0:
		return "4"
	case 1:
		return "5"
	case 2:
		return "6"
	case 3:
		return "8"
	case 4:
		return "9"
	case 5:
		return "10"
	case 6:
		return "Field"
	default:
		return "Error"
	}
}

// Create new dice pair, requires 6 pip resources and a back image
func createDicePair(pips [6]*fyne.StaticResource, back *canvas.Image) (d1 die, d2 die) {
	for i, p := range pips {
		d1.pip[i] = canvas.NewImageFromResource(p)
		d1.pip[i].SetMinSize(fyne.NewSize(75, 75))

		d2.pip[i] = canvas.NewImageFromResource(p)
		d2.pip[i].SetMinSize(fyne.NewSize(75, 75))
	}

	d1.back = back
	d1.back.SetMinSize(fyne.NewSize(75, 75))

	d1.cont = container.NewCenter(back, d1.pip[0])
	d1.cont.Move(fyne.NewPos(512, 424))

	d2.back = back
	d2.back.SetMinSize(fyne.NewSize(75, 75))

	d2.cont = container.NewCenter(back, d1.pip[0])
	d2.cont.Move(fyne.NewPos(602, 424))

	return
}

// Return a chip stack image and any remaining chips for overflow
func chipStackImage(p int, amt uint64) (img *canvas.Image, rem uint64) {
	if roll.min == 0 {
		return
	}

	position := fyne.NewPos(0, 0)
	switch p {
	case 0:
		position = fyne.NewPos(457, 205)
	case 1:
		position = fyne.NewPos(542, 205)
	case 2:
		position = fyne.NewPos(631, 205)
	case 3:
		position = fyne.NewPos(720, 205)
	case 4:
		position = fyne.NewPos(809, 205)
	case 5:
		position = fyne.NewPos(897, 205)
	case 6:
		position = fyne.NewPos(258, 205)
	}

	single := amt / roll.min
	ten := amt / (roll.min * 10)
	hun := amt / (roll.min * 100)

	if single < 1 {
		img = canvas.NewImageFromImage(nil)
	} else if single < 10 {
		img = canvas.NewImageFromResource(chipStack[single])
	} else if ten < 10 {
		rem = single - (ten * 10)
		img = canvas.NewImageFromResource(chipStack[ten+9])
	} else if hun < 10 {
		rem = single - (hun * 100)
		img = canvas.NewImageFromResource(chipStack[hun+18])
	} else {
		img = canvas.NewImageFromResource(chipStack[27])
		rem = single - 2700
	}

	img.Resize(fyne.NewSize(60, 65))
	img.Move(position)

	return
}

// Over flow chip stack image placed behind when first image is at max stack
func overflowStackImage(p int, rem uint64, back bool) (img *canvas.Image, r uint64) {
	if roll.min == 0 || rem == 0 {
		return canvas.NewImageFromImage(nil), 0
	}

	off := float32(0)
	if back {
		off = 2.5
	}

	position2 := fyne.NewPos(0, 0)
	switch p {
	case 0:
		position2 = fyne.NewPos(465+off*4, 180-off*10)
	case 1:
		position2 = fyne.NewPos(550+off*4, 180-off*10)
	case 2:
		position2 = fyne.NewPos(639+off*4, 180-off*10)
	case 3:
		position2 = fyne.NewPos(728+off*4, 180-off*10)
	case 4:
		position2 = fyne.NewPos(817+off*4, 180-off*10)
	case 5:
		position2 = fyne.NewPos(905+off*4, 180-off*10)
	case 6:
		position2 = fyne.NewPos(266+off*4, 180-off*10)
	}

	ten := rem / 10
	hun := rem / 100

	if rem < 1 {
		img = canvas.NewImageFromImage(nil)
	} else if rem < 10 {
		img = canvas.NewImageFromResource(chipStack[rem])
	} else if ten < 10 {
		img = canvas.NewImageFromResource(chipStack[ten+9])
	} else if hun < 10 {
		img = canvas.NewImageFromResource(chipStack[hun+18])
	} else {
		r = rem - 2700
		img = canvas.NewImageFromResource(chipStack[27])
	}

	img.Resize(fyne.NewSize(60, 65))
	img.Move(position2)

	return
}

// Dice roll animation
func (d *die) roll(f int, t time.Duration) {
	d.Lock()
	defer d.Unlock()
	for i := 0; i < f; i++ {
		rand.NewSource(time.Now().UnixNano())
		d.cont.Objects[1] = d.pip[rand.Intn(6)]
		d.cont.Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(60, 60))
		d.cont.Objects[1].Refresh()
		time.Sleep(t)
	}
}

// Land a die on selected value i
func (d *die) land(i int) {
	if i > 5 {
		return
	}

	d.Lock()
	defer d.Unlock()
	d.cont.Objects[1] = d.pip[i]
	d.cont.Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(60, 60))
	d.cont.Objects[1].Refresh()
}

// Place chip stack images with overflow
func placeChipStack() {
	var p4, p5, p6, p8, p9, p10, field, amtPlaced uint64
	for i := 0; i <= 30; i++ {
		if _, b := gnomon.GetSCIDValuesByKey(DICESCID, uint64(i)); b != nil {
			if _, amt := gnomon.GetSCIDValuesByKey(DICESCID, fmt.Sprintf("b_%damt", i)); amt != nil {
				div := uint64(1)
				if _, tkn := gnomon.GetSCIDValuesByKey(DICESCID, fmt.Sprintf("b_%dt", i)); tkn != nil {
					div = 300
				}

				switch b[0] {
				case 0:
					p4 = p4 + (amt[0] / div)
				case 1:
					p5 = p5 + (amt[0] / div)
				case 2:
					p6 = p6 + (amt[0] / div)
				case 3:
					p8 = p8 + (amt[0] / div)
				case 4:
					p9 = p9 + (amt[0] / div)
				case 5:
					p10 = p10 + (amt[0] / div)
				case 6:
					field = field + (amt[0] / div)
				default:
					// Nothing
				}

				if addr, _ := gnomon.GetSCIDValuesByKey(DICESCID, fmt.Sprintf("b_%d", i)); addr != nil {
					if addr[0] == rpc.Wallet.Address {
						amtPlaced = amtPlaced + (amt[0] / div)
					}
				}
			}
		}
	}

	logPlaced.SetText(fmt.Sprintf("%s placed", rpc.FromAtomic(amtPlaced, 5)))

	// Front row 16-21
	// Mid row   10-15
	// Back row  4-9
	// Field     22-24

	if p4 > 0 {
		img, rem := chipStackImage(0, p4)
		D.Front.Objects[16] = img
		img2, r := overflowStackImage(0, rem, false)
		D.Front.Objects[10] = img2
		img3, _ := overflowStackImage(0, r, true)
		D.Front.Objects[4] = img3
	} else {
		D.Front.Objects[4] = canvas.NewImageFromImage(nil)
		D.Front.Objects[10] = canvas.NewImageFromImage(nil)
		D.Front.Objects[16] = canvas.NewImageFromImage(nil)
	}

	if p5 > 0 {
		img, rem := chipStackImage(1, p5)
		D.Front.Objects[17] = img
		img2, r := overflowStackImage(1, rem, false)
		D.Front.Objects[11] = img2
		img3, _ := overflowStackImage(1, r, true)
		D.Front.Objects[5] = img3
	} else {
		D.Front.Objects[5] = canvas.NewImageFromImage(nil)
		D.Front.Objects[11] = canvas.NewImageFromImage(nil)
		D.Front.Objects[17] = canvas.NewImageFromImage(nil)
	}

	if p6 > 0 {
		img, rem := chipStackImage(2, p6)
		D.Front.Objects[18] = img
		img2, r := overflowStackImage(2, rem, false)
		D.Front.Objects[12] = img2
		img3, _ := overflowStackImage(2, r, true)
		D.Front.Objects[6] = img3
	} else {
		D.Front.Objects[6] = canvas.NewImageFromImage(nil)
		D.Front.Objects[12] = canvas.NewImageFromImage(nil)
		D.Front.Objects[18] = canvas.NewImageFromImage(nil)
	}

	if p8 > 0 {
		img, rem := chipStackImage(3, p8)
		D.Front.Objects[19] = img
		img2, r := overflowStackImage(3, rem, false)
		D.Front.Objects[13] = img2
		img3, _ := overflowStackImage(3, r, true)
		D.Front.Objects[7] = img3
	} else {
		D.Front.Objects[7] = canvas.NewImageFromImage(nil)
		D.Front.Objects[13] = canvas.NewImageFromImage(nil)
		D.Front.Objects[19] = canvas.NewImageFromImage(nil)
	}

	if p9 > 0 {
		img, rem := chipStackImage(4, p9)
		D.Front.Objects[20] = img
		img2, r := overflowStackImage(4, rem, false)
		D.Front.Objects[14] = img2
		img3, _ := overflowStackImage(4, r, true)
		D.Front.Objects[8] = img3
	} else {
		D.Front.Objects[8] = canvas.NewImageFromImage(nil)
		D.Front.Objects[14] = canvas.NewImageFromImage(nil)
		D.Front.Objects[20] = canvas.NewImageFromImage(nil)
	}

	if p10 > 0 {
		img, rem := chipStackImage(5, p10)
		D.Front.Objects[21] = img
		img2, r := overflowStackImage(5, rem, false)
		D.Front.Objects[15] = img2
		img3, _ := overflowStackImage(5, r, true)
		D.Front.Objects[9] = img3
	} else {
		D.Front.Objects[9] = canvas.NewImageFromImage(nil)
		D.Front.Objects[15] = canvas.NewImageFromImage(nil)
		D.Front.Objects[21] = canvas.NewImageFromImage(nil)
	}

	if field > 0 {
		img, rem := chipStackImage(6, field)
		D.Front.Objects[24] = img
		img2, r := overflowStackImage(6, rem, false)
		D.Front.Objects[23] = img2
		img3, _ := overflowStackImage(6, r, true)
		D.Front.Objects[22] = img3
	} else {
		D.Front.Objects[22] = canvas.NewImageFromImage(nil)
		D.Front.Objects[23] = canvas.NewImageFromImage(nil)
		D.Front.Objects[24] = canvas.NewImageFromImage(nil)
	}
}

// Set default chip stack images
func setDefaultChips() {
	if chipStack == nil {
		chipStack = make(map[uint64]*fyne.StaticResource)
	}

	chipStack = map[uint64]*fyne.StaticResource{
		0:  nil,
		1:  resourceStack0Png,
		2:  resourceStack1Png,
		3:  resourceStack2Png,
		4:  resourceStack3Png,
		5:  resourceStack4Png,
		6:  resourceStack5Png,
		7:  resourceStack6Png,
		8:  resourceStack7Png,
		9:  resourceStack8Png,
		10: resourceStack9Png,
		11: resourceStack10Png,
		12: resourceStack11Png,
		13: resourceStack12Png,
		14: resourceStack13Png,
		15: resourceStack14Png,
		16: resourceStack15Png,
		17: resourceStack16Png,
		18: resourceStack17Png,
		19: resourceStack18Png,
		20: resourceStack19Png,
		21: resourceStack20Png,
		22: resourceStack21Png,
		23: resourceStack22Png,
		24: resourceStack23Png,
		25: resourceStack24Png,
		26: resourceStack25Png,
		27: resourceStack26Png,
	}
}

// Switch dice images, downloads image files if none exists in pathDice
func getDice(url string) {
	var files []string
	Settings.dice.URL = url
	path := filepath.Join(pathDice, Settings.dice.Name, "dice1.png")
	if !dreams.FileExists(path, "Dice") {
		logger.Println("[Dice] Downloading " + Settings.dice.URL)
		if err := dreams.DownloadFile(url, path); err != nil {
			logger.Errorln("[getDice]", err)
			return
		}

		files = GetZip(Settings.dice.Name, Settings.dice.URL)
	} else {
		var err error
		files, err = filepath.Glob(filepath.Join(pathDice, Settings.dice.Name) + string(filepath.Separator) + "*.png")
		if err != nil {
			logger.Errorln("[getDice]", err)
			return
		}
	}

	if len(files) < 7 {
		logger.Errorln("[getDice] Invalid number of dice asset files")
		return
	}

	var diceRes [7]*fyne.StaticResource
	for i := 0; i < 7; i++ {
		by, err := os.ReadFile(files[i])
		if err != nil {
			logger.Errorln("[getDice]", err)
			return
		}
		diceRes[i] = fyne.NewStaticResource(files[i], by)
	}

	die1, die2 = createDicePair(
		[6]*fyne.StaticResource{
			diceRes[1],
			diceRes[2],
			diceRes[3],
			diceRes[4],
			diceRes[5],
			diceRes[6]},
		canvas.NewImageFromResource(diceRes[0]))
}

// Handle zip files for packaged assets
func GetZip(name, assetPath string) (filenames []string) {
	err := os.MkdirAll(assetPath, os.ModePerm)
	if err != nil {
		logger.Errorln("[GetZip]", err)
		return
	}

	path := filepath.Join(assetPath, name+".zip")
	filenames, err = dreams.UnzipFile(path, strings.TrimSuffix(path, ".zip"))
	if err != nil {
		logger.Errorln("[GetZip]", err)
		return
	}

	logger.Debugln("[GetZip] Unzipped files:\n" + strings.Join(filenames, "\n"))

	return
}

// Dice dreams.AssetSelect
func DiceSelect(assets map[string]string) fyne.CanvasObject {
	var max *fyne.Container
	options := []string{"Light", "Dark"}
	icon := menu.AssetIcon(resourceDiceCirclePng.StaticContent, "", 60)
	Settings.dice.Select = widget.NewSelect(options, nil)
	Settings.dice.Select.SetSelectedIndex(0)
	Settings.dice.Select.OnChanged = func(s string) {
		switch Settings.dice.Select.SelectedIndex() {
		case -1:
			Settings.dice.Name = "light"
		case 0:
			Settings.dice.Name = "light"
		case 1:
			Settings.dice.Name = "dark"
		default:
			Settings.dice.Name = s
		}

		go func() {
			scid := assets[s]
			_, collection, _ := gnomes.GetAssetInfo(scid)
			if menu.IsDreamsNFACollection(collection) {
				getDice(gnomes.GetAssetUrl(0, scid))
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon(collection, s, gnomes.GetAssetUrl(1, scid), 60)
			} else {
				Settings.dice.URL = ""
				img := canvas.NewImageFromResource(resourceDiceCirclePng)
				img.SetMinSize(fyne.NewSize(60, 60))
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = img

				if s == "Dark" {
					die1, die2 = createDicePair(
						[6]*fyne.StaticResource{
							resourceDarkDice1Png,
							resourceDarkDice2Png,
							resourceDarkDice3Png,
							resourceDarkDice4Png,
							resourceDarkDice5Png,
							resourceDarkDice6Png},
						canvas.NewImageFromResource(resourceDarkDice0Png))
				} else {
					die1, die2 = createDicePair(
						[6]*fyne.StaticResource{
							resourceDice1Png,
							resourceDice2Png,
							resourceDice3Png,
							resourceDice4Png,
							resourceDice5Png,
							resourceDice6Png},
						canvas.NewImageFromResource(resourceDice0Png))
				}
			}

			D.Front.Objects[0] = die1.cont
			D.Front.Objects[1] = die2.cont
			D.Front.Refresh()

			roll.rolled = ""
			D.Front.Objects[2].(*canvas.Text).Text = roll.rolled
			D.Front.Objects[2].(*canvas.Text).Refresh()

			roll.result = ""
			D.Front.Objects[3].(*canvas.Text).Text = roll.result
			D.Front.Objects[3].(*canvas.Text).Refresh()
		}()
	}

	Settings.dice.Select.PlaceHolder = "Dice:"
	max = container.NewBorder(nil, nil, icon, nil, container.NewVBox(Settings.dice.Select))

	return max
}

// Add dice asset to AssetSelect options
func (s *settings) AddDice(add, check string) {
	s.dice.Add(add, check)
}

// Sort dice package AssetSelect options
func (s *settings) SortAssets() {
	sort.Strings(s.dice.Select.Options)

	ld := []string{"Light", "Dark"}
	s.dice.Select.Options = append(ld, s.dice.Select.Options...)
}

// Clear dice package AssetSelect options
func (s *settings) ClearAssets() {
	s.dice.Select.Options = []string{}
}
