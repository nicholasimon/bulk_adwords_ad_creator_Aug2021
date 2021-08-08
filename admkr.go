package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var ███████████████████████████████

	//build ads
	kwlen, hdlen, desc1len, desc2len, dispurllen int
	maxhd                                        = 30
	maxdesc                                      = 90
	maxdisp                                      = 15

	// menu
	creatadson                                                               bool
	loadedkw, loadedhd, loadeddesc1, loadeddesc2, loadeddisp, loadeddest     bool
	choosegeneric, genbuy, genrent, gensale, gensign                         bool
	buttonselect                                                             int
	messagetxt                                                               = make([]msg, 20)
	messagecount                                                             int
	busnameinc, busnameall, busnamerand, contactinc, contactall, contactrand bool

	dispasis, headasis, keywasis, desc2asis, desc1asis = true, true, true, true, true

	keywup, keywtitle, keywlower, keywsent, headup, headtitle, headlower, headsent, desc1up, desc1title, desc1lower, desc1sent, desc2up, desc2title, desc2lower, desc2sent, dispup, disptitle, displower, dispsent bool

	//words
	keywords        = make([]words, 10000)
	headlines       = make([]words, 10000)
	description1    = make([]words, 10000)
	description2    = make([]words, 10000)
	displayurls     = make([]words, 10000)
	destinationurls = make([]words, 10000)
	adscomplete     = make([]ads, 10000)
	//files
	dropfileon bool
	filename   []string
	// fonts
	fontui, font1 *rl.Font

	//imgs
	cursor = rl.NewRectangle(0, 0, 12, 12)
	// core
	options, paused, scanlines, pixelnoise, ghosting                                bool
	optionselect, tilesize, centerblok, drawblok, drawbloknext, draww, drawh, drawa int
	mouseblok                                                                       int
	mousepos                                                                        rl.Vector2
	centerlines, grid, debug, fadeblinkon, fadeblink2on                             bool
	monw, monh                                                                      int
	fps                                                                             = 30
	framecount                                                                      int
	imgs                                                                            rl.Texture2D
	camera, camerabackg                                                             rl.Camera2D
	fadeblink                                                                       = float32(0.2)
	fadeblink2                                                                      = float32(0.1)
	onoff1, onoff2, onoff3, onoff6, onoff10, onoff15, onoff30, onoff60              bool
)

type ads struct {
	activ                                    bool
	headline, desc1, desc2, dispurl, desturl string
}
type words struct {
	activ   bool
	text    string
	textlen int
}
type msg struct {
	activ bool
	txt   string
}

/*
AdWords
	Headline 1	30 characters
	Headline 2	30 characters
	Headline 3	30 characters
	Description 1	90 characters
	Description 2	90 characters
	Path (2)	15 characters each

*/

func raylib() { // MARK: raylib
	rl.InitWindow(monw, monh, "admkr")
	//rl.ToggleFullscreen()
	rl.SetExitKey(rl.KeyEnd) // key to end the game and close window
	// MARK: load images
	imgs = rl.LoadTexture("./data/imgs.png") // load images
	font1 = rl.LoadFont("./data/arial.ttf")
	fontui = rl.LoadFont("./data/mecha.png")

	rl.SetTargetFPS(fps)
	rl.HideCursor()
	//	rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		framecount++
		mousepos = rl.GetMousePosition()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//	drawnocameraback()
		rl.BeginMode2D(camera)

		drawlayers()

		rl.EndMode2D()
		drawnocamera()

		if debug {
			drawdebug()
		}
		update()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func drawlayers() { // MARK: drawlayers

}
func drawnocamera() { // MARK: drawnocamera

	//backg
	rl.DrawRectangle(0, 0, monw, monh, rl.RayWhite)
	rl.DrawRectangle(0, 0, monw/2, monh, blulitewindows())
	//	rl.DrawText("PREVIEW", monw/3+10, 10, 20, bluwindows())
	v2 := rl.NewVector2(float32(monw/2+10), 10)
	rl.DrawTextEx(*fontui, "preview", v2, float32(fontui.BaseSize)*1.0, 5, bluwindows())

	// draw menu buttons main screen
	drawmainmenu()
	drawmsgbox()
	drawpreview()
	//drop file overlay
	if dropfileon {
		drawdropfile()
	}
	// choose generic overlay
	if choosegeneric {
		drawgeneric()
	}
	// create ads overlay
	if creatadson {
		drawcreateads()
	}
	//mouse cursor
	v2 = rl.NewVector2(mousepos.X, mousepos.Y)
	rl.DrawTextureRec(imgs, cursor, v2, rl.White)

}
func drawpreview() {

	v2 := rl.NewVector2(float32(monw/2+40), 50)
	rl.DrawTextEx(*font1, "Adwords Headlines are up to 30", v2, float32(font1.BaseSize)*1.0, 5, adhdcolor())
	v2.Y += 50

	textrec := rl.NewRectangle(v2.X, v2.Y, float32(monw/2)-80, float32(monh-100))
	rl.DrawTextRecEx(*font1, "This is a description line, number one, which can be up to 90 characters in length!", textrec, float32(font1.BaseSize)*1.0, 5, true, adtxtcolor(), 0, 0, blulitewindows(), rl.White)

}
func drawmsgbox() { // MARK: drawmsgbox
	rl.DrawRectangle(20, monh-120, monw/3-40, 80, rl.White)

	count := 0

	y := monh - 55
	for a := len(messagetxt) - 1; a > -1; a-- {
		if messagetxt[a].activ {
			rl.DrawText(messagetxt[a].txt, 30, y, 10, rl.Black)
			count++
			y -= 15
		}

		if count == 4 {
			break
		}
	}

}
func drawdropfile() { // MARK: drawdropfile

	rl.DrawRectangle(0, 0, monw, monh, rl.White)

	textlen := rl.MeasureText("drag a .csv file here to load", 40)
	v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+20, float32(monh/2)-200)
	rl.DrawTextEx(*fontui, "drag a .csv file here to load", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())

	if buttonselect != 1 && buttonselect != 5 && buttonselect != 6 {
		textlen = rl.MeasureText("or click here to load generic words", 40)
		v2 = rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+20, float32(monh/2)+140)

		genericrec := rl.NewRectangle(float32(monw/2)-(float32(textlen)/2)-20, float32(monh/2)+138, float32(textlen), 60)

		if rl.CheckCollisionPointRec(mousepos, genericrec) {
			rl.DrawRectangleRec(genericrec, blulitewindows())
			rl.DrawTextEx(*fontui, "or click here to load generic words", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				choosegeneric = true
				dropfileon = false
				buttonselect = 0
			}

		} else {
			rl.DrawRectangleRec(genericrec, bludarkwindows())
			rl.DrawTextEx(*fontui, "or click here to load generic words", v2, float32(fontui.BaseSize)*3.0, 5, rl.White)
		}
	}

	rl.DrawRectangle((monw/2)-10, (monh/2)-100, 20, 200, blulitewindows())
	rl.DrawRectangle((monw/2)-100, (monh/2)-10, 200, 20, blulitewindows())

	closwinrec := rl.NewRectangle(float32(monw-50), 10, 25, 25)
	if rl.CheckCollisionPointRec(mousepos, closwinrec) {
		rl.DrawRectangleRec(closwinrec, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, rl.White)
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, rl.White)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			dropfileon = false
			buttonselect = 0
		}

	} else {
		rl.DrawRectangleRec(closwinrec, blulitewindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, bludarkwindows())
	}

	if rl.IsFileDropped() {
		filename = rl.GetDroppedFiles()
		parsecsv()
		buttonselect = 0
		dropfileon = false
		rl.ClearDroppedFiles()
	}

}
func drawmainmenu() { // MARK: drawmainmenu

	createrec := rl.NewRectangle(float32(monw/2)-150, float32(monh-120), 100, 40)
	createv2 := rl.NewVector2(createrec.X+25, createrec.Y+10)
	if rl.CheckCollisionPointRec(mousepos, createrec) {
		rl.DrawRectangleRec(createrec, rl.White)
		rl.DrawTextEx(*fontui, "create", createv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			creatadson = true
		}
	} else {
		rl.DrawRectangleRec(createrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "create", createv2, float32(fontui.BaseSize)*1.0, 5, rl.White)
	}

	x := 10
	y := 10
	v2 := rl.NewVector2(float32(x), float32(y))
	rl.DrawTextEx(*fontui, "ad details", v2, float32(fontui.BaseSize)*2.0, 5, bludarkwindows())

	y += 40

	recw := float32(175)
	rech := float32(30)
	infotxtv2 := rl.NewVector2(10, float32(monh-30))
	textlen := rl.MeasureText("keywords", 20)
	keyrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, keyrec) && !dropfileon {
		rl.DrawRectangleRec(keyrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "load keyword list csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 1
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(keyrec, bluwindows())
	}
	//keywords select boxes
	rl.DrawText("sentence", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec := rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if keywsent {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "sentence case - This is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if keywsent {
				keywsent = false
			} else {
				keywasis = false
				keywsent = true
				keywtitle = false
				keywlower = false
				keywup = false
			}
		}
	}

	rl.DrawText("title", x+int(recw)+115, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+90), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if keywtitle {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "title case - This Is An Example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if keywtitle {
				keywtitle = false
			} else {
				keywasis = false
				keywsent = false
				keywtitle = true
				keywlower = false
				keywup = false
			}
		}
	}

	rl.DrawText("upper", x+int(recw)+170, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+145), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if keywup {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "upper case - THIS IS AN EXAMPLE", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if keywup {
				keywup = false
			} else {
				keywasis = false
				keywsent = false
				keywtitle = false
				keywlower = false
				keywup = true
			}
		}
	}

	rl.DrawText("lower", x+int(recw)+235, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+210), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if keywlower {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "lower case - this is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if keywlower {
				keywlower = false
			} else {
				keywasis = false
				keywsent = false
				keywtitle = false
				keywlower = true
				keywup = false
			}
		}
	}

	rl.DrawText("as provided", x+int(recw)+295, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+270), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if keywasis {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "case will remain same as in the csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if keywasis {
				keywasis = false
				keywsent = true
			} else {
				keywasis = true
				keywsent = false
				keywtitle = false
				keywlower = false
				keywup = false
			}
		}
	}

	textx := x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "keywords", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)
	y += 40

	textlen = rl.MeasureText("headlines", 20)
	headrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, headrec) && !dropfileon {
		rl.DrawRectangleRec(headrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "ad headline", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 2
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(headrec, bluwindows())
	}
	//headlines select boxes
	rl.DrawText("sentence", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if headsent {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "sentence case - This is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if headsent {
				headsent = false
			} else {
				headasis = false
				headsent = true
				headtitle = false
				headlower = false
				headup = false
			}
		}
	}

	rl.DrawText("title", x+int(recw)+115, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+90), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if headtitle {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "title case - This Is An Example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if headtitle {
				headtitle = false
			} else {
				headasis = false
				headsent = false
				headtitle = true
				headlower = false
				headup = false
			}
		}
	}

	rl.DrawText("upper", x+int(recw)+170, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+145), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if headup {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "upper case - THIS IS AN EXAMPLE", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if headup {
				headup = false
			} else {
				headasis = false
				headsent = false
				headtitle = false
				headlower = false
				headup = true
			}
		}
	}

	rl.DrawText("lower", x+int(recw)+235, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+210), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if headlower {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "lower case - this is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if headlower {
				headlower = false
			} else {
				headasis = false
				headsent = false
				headtitle = false
				headlower = true
				headup = false
			}
		}
	}

	rl.DrawText("as provided", x+int(recw)+295, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+270), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if headasis {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "case will remain same as in the csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if headasis {
				headasis = false
				headsent = true
			} else {
				headasis = true
				headsent = false
				headtitle = false
				headlower = false
				headup = false
			}
		}
	}

	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "headlines", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)
	y += 40

	//desc1 select boxes
	rl.DrawText("sentence", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc1sent {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "sentence case - This is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc1sent {
				desc1sent = false
			} else {
				desc1asis = false
				desc1sent = true
				desc1title = false
				desc1lower = false
				desc1up = false
			}
		}
	}

	rl.DrawText("title", x+int(recw)+115, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+90), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc1title {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "title case - This Is An Example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc1title {
				desc1title = false
			} else {
				desc1asis = false
				desc1sent = false
				desc1title = true
				desc1lower = false
				desc1up = false
			}
		}
	}

	rl.DrawText("upper", x+int(recw)+170, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+145), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc1up {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "upper case - THIS IS AN EXAMPLE", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc1up {
				desc1up = false
			} else {
				desc1asis = false
				desc1sent = false
				desc1title = false
				desc1lower = false
				desc1up = true
			}
		}
	}

	rl.DrawText("lower", x+int(recw)+235, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+210), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc1lower {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "lower case - this is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc1lower {
				desc1lower = false
			} else {
				desc1asis = false
				desc1sent = false
				desc1title = false
				desc1lower = true
				desc1up = false
			}
		}
	}

	rl.DrawText("as provided", x+int(recw)+295, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+270), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc1asis {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "case will remain same as in the csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc1asis {
				desc1asis = false
				desc1sent = true
			} else {
				desc1asis = true
				desc1sent = false
				desc1title = false
				desc1lower = false
				desc1up = false
			}
		}
	}

	textlen = rl.MeasureText("description one", 20)
	desc1rec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, desc1rec) && !dropfileon {
		rl.DrawRectangleRec(desc1rec, bludarkwindows())
		rl.DrawTextEx(*fontui, "1st description line", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 3
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(desc1rec, bluwindows())
	}

	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "description one", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)
	y += 40

	//desc2 select boxes
	rl.DrawText("sentence", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc2sent {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "sentence case - This is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc2sent {
				desc2sent = false
			} else {
				desc2asis = false
				desc2sent = true
				desc2title = false
				desc2lower = false
				desc2up = false
			}
		}
	}

	rl.DrawText("title", x+int(recw)+115, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+90), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc2title {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "title case - This Is An Example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc2title {
				desc2title = false
			} else {
				desc2asis = false
				desc2sent = false
				desc2title = true
				desc2lower = false
				desc2up = false
			}
		}
	}

	rl.DrawText("upper", x+int(recw)+170, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+145), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc2up {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "upper case - THIS IS AN EXAMPLE", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc2up {
				desc2up = false
			} else {
				desc2asis = false
				desc2sent = false
				desc2title = false
				desc2lower = false
				desc2up = true
			}
		}
	}

	rl.DrawText("lower", x+int(recw)+235, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+210), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc2lower {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "lower case - this is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc2lower {
				desc2lower = false
			} else {
				desc2asis = false
				desc2sent = false
				desc2title = false
				desc2lower = true
				desc2up = false
			}
		}
	}

	rl.DrawText("as provided", x+int(recw)+295, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+270), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if desc2asis {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "case will remain same as in the csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if desc2asis {
				desc2asis = false
				desc2sent = true
			} else {
				desc2asis = true
				desc2sent = false
				desc2title = false
				desc2lower = false
				desc2up = false
			}
		}
	}

	textlen = rl.MeasureText("description two", 20)
	desc2rec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, desc2rec) && !dropfileon {
		rl.DrawRectangleRec(desc2rec, bludarkwindows())
		rl.DrawTextEx(*fontui, "2nd description line", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 4
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(desc2rec, bluwindows())
	}

	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "description two", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)
	y += 40

	//displayurl select boxes
	rl.DrawText("sentence", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if dispsent {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "sentence case - This is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if dispsent {
				dispsent = false
			} else {
				dispasis = false
				dispsent = true
				disptitle = false
				displower = false
				dispup = false
			}
		}
	}

	rl.DrawText("title", x+int(recw)+115, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+90), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if disptitle {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "title case - This Is An Example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if disptitle {
				disptitle = false
			} else {
				dispasis = false
				dispsent = false
				disptitle = true
				displower = false
				dispup = false
			}
		}
	}

	rl.DrawText("upper", x+int(recw)+170, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+145), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if dispup {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "upper case - THIS IS AN EXAMPLE", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if dispup {
				dispup = false
			} else {
				dispasis = false
				dispsent = false
				disptitle = false
				displower = false
				dispup = true
			}
		}
	}

	rl.DrawText("lower", x+int(recw)+235, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+210), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if displower {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "lower case - this is an example", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if displower {
				displower = false
			} else {
				dispasis = false
				dispsent = false
				disptitle = false
				displower = true
				dispup = false
			}
		}
	}

	rl.DrawText("as provided", x+int(recw)+295, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+270), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if dispasis {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "case will remain same as in the csv", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if dispasis {
				dispasis = false
				dispsent = true
			} else {
				dispasis = true
				dispsent = false
				disptitle = false
				displower = false
				dispup = false
			}
		}
	}

	textlen = rl.MeasureText("display url", 20)
	displurlrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, displurlrec) && !dropfileon {
		rl.DrawRectangleRec(displurlrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "visible ad url", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 5
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(displurlrec, bluwindows())
	}
	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))

	rl.DrawTextEx(*fontui, "display url", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)

	y += 40
	textlen = rl.MeasureText("destination url", 20)
	desturlrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, desturlrec) && !dropfileon {
		rl.DrawRectangleRec(desturlrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "hidden landing page url", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !dropfileon {
				buttonselect = 6
				dropfileon = true
			}
		}
	} else {
		rl.DrawRectangleRec(desturlrec, bluwindows())
	}
	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "destination url", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)

	y += 40
	v2 = rl.NewVector2(float32(x), float32(y))
	rl.DrawTextEx(*fontui, "extra information", v2, float32(fontui.BaseSize)*2.0, 5, bludarkwindows())

	y += 40
	textlen = rl.MeasureText("business name", 20)
	businessrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, businessrec) && !dropfileon {
		rl.DrawRectangleRec(businessrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "include a business name", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		}
	} else {
		rl.DrawRectangleRec(businessrec, bluwindows())
	}

	//business name select boxes
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if busnameinc {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "include business name in ads", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if busnameinc {
				busnameinc = false
			} else {
				busnameinc = true
				if !busnameall && !busnamerand {
					busnamerand = true
				}
			}
		}
	}
	rl.DrawText("include", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+80), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if busnameall {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "name will appear in all ads where possible", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if busnameall {
				busnameall = false
				busnamerand = true
			} else {
				busnameall = true
				busnamerand = false
			}
		}
	}
	rl.DrawText("all ads", x+int(recw)+105, y+10, 10, bludarkwindows())

	selectrec = rl.NewRectangle(float32(x+int(recw)+150), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if busnamerand {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "name will appear in random ads", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if busnamerand {
				busnamerand = false
				busnameall = true
			} else {
				busnamerand = true
				busnameall = false
			}
		}
	}
	rl.DrawText("random", x+int(recw)+175, y+10, 10, bludarkwindows())

	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "business name", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)

	y += 40
	textlen = rl.MeasureText("contact details", 20)
	contactrec := rl.NewRectangle(float32(x), float32(y), recw, rech)
	if rl.CheckCollisionPointRec(mousepos, contactrec) && !dropfileon {
		rl.DrawRectangleRec(contactrec, bludarkwindows())
		rl.DrawTextEx(*fontui, "include contact details", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		}
	} else {
		rl.DrawRectangleRec(contactrec, bluwindows())
	}
	//contact select boxes
	selectrec = rl.NewRectangle(float32(x+int(recw)+10), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if contactinc {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "include contact details in ads", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if contactinc {
				contactinc = false
			} else {
				contactinc = true
				if !contactall && !contactrand {
					contactrand = true
				}
			}
		}
	}
	rl.DrawText("include", x+int(recw)+35, y+10, 10, bludarkwindows())
	selectrec = rl.NewRectangle(float32(x+int(recw)+80), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if contactall {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "name will appear in all ads where possible", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if contactall {
				contactall = false
				contactrand = true
			} else {
				contactall = true
				contactrand = false
			}
		}
	}
	rl.DrawText("all ads", x+int(recw)+105, y+10, 10, bludarkwindows())

	selectrec = rl.NewRectangle(float32(x+int(recw)+150), float32(y+5), 20, 20)
	rl.DrawRectangleRec(selectrec, rl.White)
	if contactrand {
		rl.DrawRectangle(int(selectrec.X)+3, int(selectrec.Y)+3, 14, 14, bludarkwindows())
	}
	if rl.CheckCollisionPointRec(mousepos, selectrec) && !dropfileon && !choosegeneric {
		rl.DrawRectangleRec(selectrec, rl.Fade(bludarkwindows(), fadeblink))
		rl.DrawTextEx(*fontui, "name will appear in random ads", infotxtv2, float32(fontui.BaseSize)*1.0, 5, bludarkwindows())
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if contactrand {
				contactrand = false
				contactall = true
			} else {
				contactrand = true
				contactall = false
			}
		}
	}
	rl.DrawText("random", x+int(recw)+175, y+10, 10, bludarkwindows())

	textx = x + (textlen + 20)
	textx = textx / 2
	textx -= textlen / 2
	textx += 3
	v2 = rl.NewVector2(float32(x+textx), float32(y+5))
	rl.DrawTextEx(*fontui, "contact details", v2, float32(fontui.BaseSize)*1.0, 5, rl.White)
}
func drawgeneric() { // MARK: choosegeneric

	rl.DrawRectangle(0, 0, monw, monh, rl.White)

	textlen := rl.MeasureText("choose generic type below", 40)
	v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+20, float32(monh/2)-200)
	rl.DrawTextEx(*fontui, "choose generic type below", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())

	x := float32(monw / 4)
	y := float32((monh / 2) - 50)

	for a := 0; a < 4; a++ {

		buttonrec := rl.NewRectangle(x, y, 150, 40)
		if rl.CheckCollisionPointRec(mousepos, buttonrec) {
			rl.DrawRectangleRec(buttonrec, blulitewindows())
			switch a {
			case 0:
				rl.DrawText("buy", int(x+10), int(y+10), 20, bludarkwindows())
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					genbuy = true
					genrent = false
					gensale = false
					gensign = false
					choosegeneric = false
					messagetxt[messagecount].activ = true
					messagetxt[messagecount].txt = "buy generic words loaded"
					messagecount++
					loadedhd = true
					loadeddesc1 = true
					loadeddesc2 = true
				}
			case 1:
				rl.DrawText("rent", int(x+10), int(y+10), 20, bludarkwindows())
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					genbuy = false
					genrent = true
					gensale = false
					gensign = false
					choosegeneric = false
					messagetxt[messagecount].activ = true
					messagetxt[messagecount].txt = "rent generic words loaded"
					messagecount++
					loadedhd = true
					loadeddesc1 = true
					loadeddesc2 = true
				}
			case 2:
				rl.DrawText("on sale", int(x+10), int(y+10), 20, bludarkwindows())
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					genbuy = false
					genrent = false
					gensale = true
					gensign = false
					choosegeneric = false
					messagetxt[messagecount].activ = true
					messagetxt[messagecount].txt = "on sale generic words loaded"
					messagecount++
					loadedhd = true
					loadeddesc1 = true
					loadeddesc2 = true
				}
			case 3:
				rl.DrawText("sign up", int(x+10), int(y+10), 20, bludarkwindows())
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					genbuy = false
					genrent = false
					gensale = false
					gensign = true
					choosegeneric = false
					messagetxt[messagecount].activ = true
					messagetxt[messagecount].txt = "sign up generic words loaded"
					messagecount++
					loadedhd = true
					loadeddesc1 = true
					loadeddesc2 = true
				}
			}
		} else {
			rl.DrawRectangleRec(buttonrec, bludarkwindows())
			switch a {
			case 0:
				rl.DrawText("buy", int(x+10), int(y+10), 20, rl.White)
			case 1:
				rl.DrawText("rent", int(x+10), int(y+10), 20, rl.White)
			case 2:
				rl.DrawText("on sale", int(x+10), int(y+10), 20, rl.White)
			case 3:
				rl.DrawText("sign up", int(x+10), int(y+10), 20, rl.White)
			}
		}

		x += 170

	}

	closwinrec := rl.NewRectangle(float32(monw-50), 10, 25, 25)
	if rl.CheckCollisionPointRec(mousepos, closwinrec) {
		rl.DrawRectangleRec(closwinrec, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, rl.White)
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, rl.White)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			choosegeneric = false
			buttonselect = 0
		}

	} else {
		rl.DrawRectangleRec(closwinrec, blulitewindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, bludarkwindows())
	}

}
func drawcreateads() { // MARK: drawcreateads

	rl.DrawRectangle(0, 0, monw, monh, rl.White)

	if !loadedkw {
		textlen := rl.MeasureText("you need to load keywords in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+20, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load keywords in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else if !loadedhd {
		textlen := rl.MeasureText("you need to load headlines or use generics in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+40, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load headlines or use generics in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else if !loadeddesc1 {
		textlen := rl.MeasureText("you need to load description 1 or use generics in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+60, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load description 1 or use generics in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else if !loadeddesc2 {
		textlen := rl.MeasureText("you need to load description 2 or use generics in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+60, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load description 2 or use generics in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else if !loadeddisp {
		textlen := rl.MeasureText("you need to load display urls in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+60, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load display urls in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else if !loadeddest {
		textlen := rl.MeasureText("you need to load destination urls in .csv format", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+60, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "you need to load destination urls in .csv format", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("<< go back", int(closerec.X)+40, int(closerec.Y)+20, 20, rl.White)
		}
	} else {
		textlen := rl.MeasureText("a .csv will be saved in the admkr\ads\\ directory", 40)
		v2 := rl.NewVector2(float32(monw/2)-(float32(textlen)/2)+60, float32(monh/2)-200)
		rl.DrawTextEx(*fontui, "a .csv will be saved in the admkr\ads\\ directory", v2, float32(fontui.BaseSize)*3.0, 5, bludarkwindows())
		closerec := rl.NewRectangle(float32(monw/2)-100, float32(monh/2)-30, 200, 60)

		if rl.CheckCollisionPointRec(mousepos, closerec) {
			rl.DrawRectangleRec(closerec, blulitewindows())
			rl.DrawText("build ads!", int(closerec.X)+55, int(closerec.Y)+20, 20, bludarkwindows())
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				creatadson = false
				buildads()
			}
		} else {
			rl.DrawRectangleRec(closerec, bludarkwindows())
			rl.DrawText("build ads!", int(closerec.X)+55, int(closerec.Y)+20, 20, rl.White)
		}
	}

	closwinrec := rl.NewRectangle(float32(monw-50), 10, 25, 25)
	if rl.CheckCollisionPointRec(mousepos, closwinrec) {
		rl.DrawRectangleRec(closwinrec, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, rl.White)
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, rl.White)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			creatadson = false

		}

	} else {
		rl.DrawRectangleRec(closwinrec, blulitewindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+5, int(closwinrec.X)+20, int(closwinrec.Y)+20, bludarkwindows())
		rl.DrawLine(int(closwinrec.X)+5, int(closwinrec.Y)+20, int(closwinrec.X)+20, int(closwinrec.Y)+5, bludarkwindows())
	}

}

func buildads() { // MARK: drawcreateads

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	currentdir := filepath.Dir(ex)

	filename2 := "something"

	if genbuy {
		filename2 = currentdir + "./generics/head1.csv"
		// Open the file
		csvfile, err := os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r := csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count := 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			headlines[count].activ = true
			headlines[count].text = record[0]
			headlines[count].textlen = 0
			for range headlines[count].text {
				headlines[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc1.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description1[count].activ = true
			description1[count].text = record[0]
			description1[count].textlen = 0
			for range description1[count].text {
				description1[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc2.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description2[count].activ = true
			description2[count].text = record[0]
			description2[count].textlen = 0
			for range description2[count].text {
				description2[count].textlen++
			}
			count++
		}
	} else if genrent {
		filename2 = currentdir + "./generics/head1.csv"
		// Open the file
		csvfile, err := os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r := csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count := 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			headlines[count].activ = true
			headlines[count].text = record[0]
			headlines[count].textlen = 0
			for range headlines[count].text {
				headlines[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc1.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description1[count].activ = true
			description1[count].text = record[0]
			description1[count].textlen = 0
			for range description1[count].text {
				description1[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc2.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description2[count].activ = true
			description2[count].text = record[0]
			description2[count].textlen = 0
			for range description2[count].text {
				description2[count].textlen++
			}
			count++
		}

	} else if gensale {
		filename2 = currentdir + "./generics/head1.csv"
		// Open the file
		csvfile, err := os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r := csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count := 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			headlines[count].activ = true
			headlines[count].text = record[0]
			headlines[count].textlen = 0
			for range headlines[count].text {
				headlines[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc1.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description1[count].activ = true
			description1[count].text = record[0]
			description1[count].textlen = 0
			for range description1[count].text {
				description1[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc2.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description2[count].activ = true
			description2[count].text = record[0]
			description2[count].textlen = 0
			for range description2[count].text {
				description2[count].textlen++
			}
			count++
		}

	} else if gensign {
		filename2 = currentdir + "./generics/head1.csv"
		// Open the file
		csvfile, err := os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r := csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count := 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			headlines[count].activ = true
			headlines[count].text = record[0]
			headlines[count].textlen = 0
			for range headlines[count].text {
				headlines[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc1.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description1[count].activ = true
			description1[count].text = record[0]
			description1[count].textlen = 0
			for range description1[count].text {
				description1[count].textlen++
			}
			count++
		}
		filename2 = currentdir + "./generics/desc2.csv"
		// Open the file
		csvfile, err = os.Open(filename2)
		if err != nil {
			rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
		} else {
			messagetxt[messagecount].activ = true
			messagetxt[messagecount].txt = "loaded " + filename2
			messagecount++
			if messagecount == 19 {
				messagecount = 0
			}
		}
		// Parse the file
		r = csv.NewReader(csvfile)
		//r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		count = 0
		for {
			// Read each record from csv
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			description2[count].activ = true
			description2[count].text = record[0]
			description2[count].textlen = 0
			for range description2[count].text {
				description2[count].textlen++
			}
			count++
		}

	}

	kwlen = 0
	for a := 0; a < len(keywords); a++ {
		if keywords[a].activ {
			kwlen++
		}
	}
	hdlen = 0
	for a := 0; a < len(headlines); a++ {
		if headlines[a].activ {
			hdlen++
		}
	}
	desc1len = 0
	for a := 0; a < len(description1); a++ {
		if description1[a].activ {
			desc1len++
		}
	}
	desc2len = 0
	for a := 0; a < len(description2); a++ {
		if description2[a].activ {
			desc2len++
		}
	}
	dispurllen = 0
	for a := 0; a < len(displayurls); a++ {
		if displayurls[a].activ {
			dispurllen++
		}
	}

	//change kw case
	if !keywasis {
		if keywsent {
			for a := 0; a < kwlen; a++ {
				keywords[a].text = sentencecase(keywords[a].text)
			}
		} else if keywtitle {
			for a := 0; a < kwlen; a++ {
				keywords[a].text = (strings.Title(strings.ToLower(keywords[a].text)))
			}
		} else if keywlower {
			for a := 0; a < kwlen; a++ {
				keywords[a].text = strings.ToLower(keywords[a].text)
			}
		} else if keywup {
			for a := 0; a < kwlen; a++ {
				keywords[a].text = (strings.ToUpper(strings.ToLower(keywords[a].text)))
			}
		}
	}
	//change headline case
	if !headasis {
		if headsent {
			for a := 0; a < kwlen; a++ {
				headlines[a].text = sentencecase(headlines[a].text)
			}
		} else if headtitle {
			for a := 0; a < kwlen; a++ {
				headlines[a].text = (strings.Title(strings.ToLower(headlines[a].text)))
			}
		} else if headlower {
			for a := 0; a < kwlen; a++ {
				headlines[a].text = strings.ToLower(headlines[a].text)
			}
		} else if headup {
			for a := 0; a < kwlen; a++ {
				headlines[a].text = (strings.ToUpper(strings.ToLower(headlines[a].text)))
			}
		}
	}

	//change desc1 case
	if !desc1asis {
		if desc1sent {
			for a := 0; a < kwlen; a++ {
				description1[a].text = sentencecase(description1[a].text)
			}
		} else if desc1title {
			for a := 0; a < kwlen; a++ {
				description1[a].text = (strings.Title(strings.ToLower(description1[a].text)))
			}
		} else if desc1lower {
			for a := 0; a < kwlen; a++ {
				description1[a].text = strings.ToLower(description1[a].text)
			}
		} else if desc1up {
			for a := 0; a < kwlen; a++ {
				description1[a].text = (strings.ToUpper(strings.ToLower(description1[a].text)))
			}
		}
	}

	//change desc2 case
	if !desc2asis {
		if desc2sent {
			for a := 0; a < kwlen; a++ {
				description2[a].text = sentencecase(description2[a].text)
			}
		} else if desc2title {
			for a := 0; a < kwlen; a++ {
				description2[a].text = (strings.Title(strings.ToLower(description2[a].text)))
			}
		} else if desc2lower {
			for a := 0; a < kwlen; a++ {
				description2[a].text = strings.ToLower(description2[a].text)
			}
		} else if desc2up {
			for a := 0; a < kwlen; a++ {
				description2[a].text = (strings.ToUpper(strings.ToLower(description2[a].text)))
			}
		}
	}

	//change dispurl case
	if !dispasis {
		if dispsent {
			for a := 0; a < kwlen; a++ {
				displayurls[a].text = sentencecase(displayurls[a].text)
			}
		} else if disptitle {
			for a := 0; a < kwlen; a++ {
				displayurls[a].text = (strings.Title(strings.ToLower(displayurls[a].text)))
			}
		} else if displower {
			for a := 0; a < kwlen; a++ {
				displayurls[a].text = strings.ToLower(displayurls[a].text)
			}
		} else if dispup {
			for a := 0; a < kwlen; a++ {
				displayurls[a].text = (strings.ToUpper(strings.ToLower(displayurls[a].text)))
			}
		}
	}

	//create final ads
	for a := 0; a < kwlen; a++ {
		adscomplete[a].activ = true
		adscomplete[a].headline = headlines[rInt(0, hdlen)].text + " " + keywords[a].text
		adscomplete[a].desc1 = description1[rInt(0, desc1len)].text
		adscomplete[a].desc2 = description2[rInt(0, desc2len)].text
		adscomplete[a].dispurl = displayurls[rInt(0, dispurllen)].text
	}

	for a := 0; a < kwlen; a++ {
		fmt.Println(adscomplete[a].headline)
		fmt.Println(adscomplete[a].desc1)
		fmt.Println(adscomplete[a].desc2)
		fmt.Println(adscomplete[a].dispurl)
		fmt.Println()
	}

	// 1. Open the file
	adsfilename := "ads_" + randfilename()
	recordFile, err := os.Create("./ads/" + adsfilename + ".csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}

	// 2. Initialize the writer
	writer := csv.NewWriter(recordFile)
	for a := 0; a < kwlen; a++ {
		var csvData = [][]string{
			{adscomplete[a].headline, adscomplete[a].desc1, adscomplete[a].desc2, adscomplete[a].dispurl},
		}
		// 3. Write all the records
		err = writer.WriteAll(csvData) // returns error
		if err != nil {
			fmt.Println("An error encountered ::", err)
		}
	}

}
func update() { // MARK: update

	input()
	timers()

}

func parsecsv() { // MARK: parsecsv
	// Open the file
	csvfile, err := os.Open(filename[0])
	if err != nil {
		rl.DrawText("Cannot load file", monw/2-100, monh/2-100, 20, bludarkwindows())
	} else {
		messagetxt[messagecount].activ = true
		messagetxt[messagecount].txt = "loaded " + filename[0]
		messagecount++
		if messagecount == 19 {
			messagecount = 0
		}
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	count := 0
	for {
		// Read each record from csv
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		switch buttonselect {
		case 1:
			keywords[count].activ = true
			keywords[count].text = record[0]
			keywords[count].textlen = 0
			for range keywords[count].text {
				keywords[count].textlen++
			}
			count++
			loadedkw = true
		case 2:
			headlines[count].activ = true
			headlines[count].text = record[0]
			headlines[count].textlen = 0
			for range headlines[count].text {
				headlines[count].textlen++
			}
			count++
			loadedhd = true
		case 3:
			description1[count].activ = true
			description1[count].text = record[0]
			description1[count].textlen = 0
			for range description1[count].text {
				description1[count].textlen++
			}
			count++
			loadeddesc1 = true
		case 4:
			description2[count].activ = true
			description2[count].text = record[0]
			description2[count].textlen = 0
			for range description2[count].text {
				description2[count].textlen++
			}
			count++
			loadeddesc2 = true
		case 5:
			displayurls[count].activ = true
			displayurls[count].text = record[0]
			displayurls[count].textlen = 0
			for range displayurls[count].text {
				displayurls[count].textlen++
			}
			count++
			loadeddisp = true
		case 6:
			destinationurls[count].activ = true
			destinationurls[count].text = record[0]
			destinationurls[count].textlen = 0
			for range destinationurls[count].text {
				destinationurls[count].textlen++
			}
			count++
			loadeddest = true
		}
	}
}

// MARK: core  █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █
func setinitialvalues() { // MARK: setinitialvalues

	messagetxt[0].activ = true
	messagetxt[0].txt = "Messages will appear here as you work"
	messagecount = 1

}
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLogLevel(rl.LogError) // hides info window
	rl.InitWindow(monw, monh, "setres")
	setres(1280, 720)
	rl.CloseWindow()
	setinitialvalues()
	raylib()

}
func input() { // MARK: input

	if rl.IsKeyPressed(rl.KeyKpAdd) {
		camera.Zoom += 0.2
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {
		camera.Zoom -= 0.2
	}
	if rl.IsKeyPressed(rl.KeyKpDivide) {
		if centerlines {
			centerlines = false
		} else {
			centerlines = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debug {
			debug = false
		} else {
			debug = true
		}
	}

}
func drawdebug() { // MARK: DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG

	rl.DrawRectangle(monw-300, 0, 300, monh, rl.Fade(rl.Black, 0.8))
	textx := monw - 290
	textx2 := monw - 145
	texty := 10

	//	camerazoomtext := fmt.Sprintf("%g", camera.Zoom)
	//	playermovingtext := strconv.FormatBool(player.moving)

	dropfileontxt := strconv.FormatBool(dropfileon)
	monwtxt := strconv.Itoa(monw)

	rl.DrawText("monwtxt", textx, texty, 10, rl.White)
	rl.DrawText(monwtxt, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("dropfileon", textx, texty, 10, rl.White)
	rl.DrawText(dropfileontxt, textx2, texty, 10, rl.White)
	texty += 12

	// fps
	rl.DrawRectangle(monw-110, monh-110, 100, 40, rl.Black)
	rl.DrawFPS(monw-100, monh-100)

}
func timers() { // MARK: timers

	if framecount%1 == 0 {
		if onoff1 {
			onoff1 = false
		} else {
			onoff1 = true
		}
	}

	if framecount%2 == 0 {
		if onoff2 {
			onoff2 = false
		} else {
			onoff2 = true
		}
	}
	if framecount%3 == 0 {
		if onoff3 {
			onoff3 = false
		} else {
			onoff3 = true
		}
	}
	if framecount%6 == 0 {
		if onoff6 {
			onoff6 = false
		} else {
			onoff6 = true
		}
	}
	if framecount%10 == 0 {
		if onoff10 {
			onoff10 = false
		} else {
			onoff10 = true
		}
	}
	if framecount%15 == 0 {
		if onoff15 {
			onoff15 = false
		} else {
			onoff15 = true
		}
	}
	if framecount%30 == 0 {
		if onoff30 {
			onoff30 = false
		} else {
			onoff30 = true
		}
	}
	if framecount%60 == 0 {
		if onoff60 {
			onoff60 = false
		} else {
			onoff60 = true
		}
	}
	if fadeblinkon {
		if fadeblink > 0.2 {
			fadeblink -= 0.05
		} else {
			fadeblinkon = false
		}
	} else {
		if fadeblink < 0.6 {
			fadeblink += 0.05
		} else {
			fadeblinkon = true
		}
	}
	if onoff3 {
		if fadeblink2on {
			if fadeblink2 > 0.1 {
				fadeblink2 -= 0.01
			} else {
				fadeblink2on = false
			}
		} else {
			if fadeblink2 < 0.2 {
				fadeblink2 += 0.01
			} else {
				fadeblink2on = true
			}
		}
	}
}
func setres(w, h int) { // MARK: setres

	if w == 0 {

		monw = rl.GetMonitorWidth(0)
		monh = rl.GetMonitorHeight(0)
		camera.Zoom = 1.0
		camerabackg.Zoom = 1.0

		if monw >= 1600 {
			tilesize = 96

		} else if monw < 1600 {
			tilesize = 72

		}

	} else {
		monw = w
		monh = h
		camera.Zoom = 1.0
		camerabackg.Zoom = 1.0

		if monw >= 1600 {
			tilesize = 96

		} else if monw < 1600 {
			tilesize = 72
		}
	}

}

// MARK: colors
// https://www.rapidtables.com/web/color/RGB_Color.html
func adtxtcolor() rl.Color {
	color := rl.NewColor(71, 81, 86, 255)
	return color
}
func adhdcolor() rl.Color {
	color := rl.NewColor(26, 13, 171, 255)
	return color
}
func darkred() rl.Color {
	color := rl.NewColor(55, 0, 0, 255)
	return color
}
func semidarkred() rl.Color {
	color := rl.NewColor(70, 0, 0, 255)
	return color
}
func brightred() rl.Color {
	color := rl.NewColor(230, 0, 0, 255)
	return color
}
func randomgrey() rl.Color {
	color := rl.NewColor(uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(0, 255)))
	return color
}
func randombluelight() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 180)), uint8(rInt(120, 256)), uint8(rInt(120, 256)), 255)
	return color
}
func randombluedark() rl.Color {
	color := rl.NewColor(0, 0, uint8(rInt(120, 250)), 255)
	return color
}
func randomyellow() rl.Color {
	color := rl.NewColor(255, uint8(rInt(150, 256)), 0, 255)
	return color
}
func randomorange() rl.Color {
	color := rl.NewColor(uint8(rInt(250, 256)), uint8(rInt(60, 210)), 0, 255)
	return color
}
func randomred() rl.Color {
	color := rl.NewColor(uint8(rInt(128, 256)), uint8(rInt(0, 129)), uint8(rInt(0, 129)), 255)
	return color
}
func randomgreen() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 170)), uint8(rInt(100, 256)), uint8(rInt(0, 50)), 255)
	return color
}
func randomcolor() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 256)), uint8(rInt(0, 256)), uint8(rInt(0, 256)), 255)
	return color
}
func brightyellow() rl.Color {
	color := rl.NewColor(uint8(255), uint8(255), uint8(0), 255)
	return color
}
func brightbrown() rl.Color {
	color := rl.NewColor(uint8(218), uint8(165), uint8(32), 255)
	return color
}
func brightgrey() rl.Color {
	color := rl.NewColor(uint8(212), uint8(212), uint8(213), 255)
	return color
}
func bluwindows() rl.Color {
	color := rl.NewColor(uint8(30), uint8(123), uint8(213), 255)
	return color
}
func bludarkwindows() rl.Color {
	color := rl.NewColor(uint8(24), uint8(103), uint8(182), 255)
	return color
}
func blulitewindows() rl.Color {
	color := rl.NewColor(uint8(173), uint8(209), uint8(237), 255)
	return color
}

// random numbers
func rF32(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}

// sentence case
func sentencecase(str string) string {
	if len(str) == 0 {
		return ""
	}
	tmp := []rune(str)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

func randfilename() string {
	return time.Now().Format("20060102150405")
}
