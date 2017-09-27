/**** Functions for drawing on the lcd ****/
package main

import (
	"fmt"
	"image"
	"image/draw"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/xdsopl/framebuffer/src/framebuffer"
)

var fb draw.Image
var context *freetype.Context
var face font.Face
var fontSize float64

func initLCD() {
	var err error
	fb, err = framebuffer.Open(FramebufferDev)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	context = freetype.NewContext()
	context.SetDst(fb)
	context.SetDPI(DPI)
	context.SetHinting(font.HintingFull)
	context.SetClip(fb.Bounds())
}

func clearLCD() {
	if fb == nil {
		return
	}

	draw.Draw(fb, fb.Bounds(), image.Black, image.ZP, draw.Src)
}

func setFont(size float64, font []byte) {
	if fb == nil {
		return
	}

	context.SetFontSize(size)
	parsedFont, _ := freetype.ParseFont(font)
	context.SetFont(parsedFont)
	opts := truetype.Options{}
	opts.Size = size
	face = truetype.NewFace(parsedFont, &opts)
	fontSize = size
}

func nextLine(p fixed.Point26_6) (fixed.Point26_6) {
	return fixed.Point26_6{p.X, p.Y + fixed.I(int(fontSize * LineHeight))}
}

func drawStringCentered(s string, src image.Image, p fixed.Point26_6) (fixed.Point26_6) {
	if fb == nil {
		fmt.Printf("Error: No framebuffer to draw on\n")
		return p
	}

	// determine the width of rendered string
	stringWidth := 0
	for _, char := range(s) {
		charWidth, ok := face.GlyphAdvance(char)
		if !ok {
			fmt.Printf("Error: Font does not support character %s\n", string(char))
			return p
		}

		stringWidth += int(float64(charWidth) / 64)
	}

	// then draw the string
	context.SetSrc(src)
	p = fixed.Point26_6{fixed.I((fb.Bounds().Max.X - int(float64(stringWidth) * 1.03) + fb.Bounds().Min.X) / 2), p.Y}
	p, err := context.DrawString(s, p)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	return p
}

func drawMessage() {
	clearLCD()

	setFont(14, goregular.TTF)
	p := nextLine(freetype.Pt(fb.Bounds().Min.X, fb.Bounds().Min.Y))
	p = drawStringCentered("Lorem ipsum dolor sit amet", image.White, p)
	p = nextLine(p)
	p = drawStringCentered("Lorem ipsum dolor sit amet", image.White, p)
	p = nextLine(p)
	p = drawStringCentered("Lorem ipsum dolor sit amet", image.White, p)

	/*setFont(11, goregular.TTF)
	p, err := drawStringCentered("Go to", image.White, freetype.Pt(fb.Bounds().Min.X + 15, 0))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	setFont(13, gomedium.TTF)
	p, err = drawStringCentered(ServerBaseURL + "/" + deviceId, image.White, p)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	setFont(11, goregular.TTF)
	p, err = drawStringCentered("to keep the party goin", image.White, p)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}*/
}
