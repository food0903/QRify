package services

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func addTextBelow(img image.Image, text string) (image.Image, error) {
	qrBounds := img.Bounds()
	textHeight := 20
	gap := 20
	newHeight := qrBounds.Dy() + gap + textHeight
	newImg := image.NewRGBA(image.Rect(0, 0, qrBounds.Dx(), newHeight))

	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, qrBounds, img, image.Point{}, draw.Over)

	col := color.Black
	point := fixed.Point26_6{
		X: fixed.I((qrBounds.Dx() - len(text)*7) / 2),
		Y: fixed.I(qrBounds.Dy() + gap + 15),
	}
	d := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
	return newImg, nil
}
