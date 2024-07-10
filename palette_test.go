package ledgrid

import (
	"image/color"
	"testing"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"
)

const (
    NumCols          = 10
    NumRows          = 15
	FieldWidth       = 100
	FieldHeight      = 100
	CornerRadius     = 20
	Margin           = 10.0
	ColorPaddingHori = 40.0
	ColorPaddingVert = 60.0
	FieldPadding     = 0.0
	FontSize         = 24.0
)

var (
	Font            = fonts.GoBold
	PaletteFileList = []string{
		"colourlovers.json",
		"palettes.json",
	}
)

func TestReadJsonData(t *testing.T) {
	for _, fileName := range PaletteFileList {
		t.Logf("Reading palette file '%s'", fileName)
		ReadJsonData(fileName)
	}
}

// func TestDrawPalette(t *testing.T) {
// 	face := fonts.NewFace(Font, FontSize)
// 	jsonFile := "colourlovers.json"
// 	palList := ReadJsonData(jsonFile)
// 	t.Logf("Number of colors in %s: %d", jsonFile, len(palList))
// 	gc := gg.NewContext(2*Margin+10*(5*FieldWidth+4*FieldPadding)+9*ColorPaddingHori,
// 		2*Margin+10*FieldHeight+10*ColorPaddingVert)
// 	gc.SetFontFace(face)
// 	gc.SetFillColor(color.White)
// 	gc.Clear()
// 	for i, pal := range palList {
// 		col, row := i/10, i%10
// 		x := Margin + float64(col)*(5*FieldWidth+4*FieldPadding+ColorPaddingHori)
// 		y := Margin + float64(row)*(FieldHeight+ColorPaddingVert)
// 		for j, color := range pal.Colors {
// 			x := x + float64(j)*(FieldWidth+FieldPadding)
// 			gc.SetFillColor(color)
// 			gc.SetStrokeWidth(0.0)
// 			gc.DrawRoundedRectangle(x, y, FieldWidth, FieldHeight, CornerRadius)
// 			gc.Fill()
// 		}
// 		gc.SetStrokeColor(color.Black)
// 		gc.DrawStringAnchored(pal.Name, x, y+FieldHeight+Margin/2, 0.0, 1.0)
// 	}
// 	gc.SavePNG("palette_test.png")
// }

func TestPaletteSamples(t *testing.T) {
	face := fonts.NewFace(Font, FontSize)
	gc := gg.NewContext(2*Margin+NumCols*(5*FieldWidth+4*FieldPadding)+(NumCols-1)*ColorPaddingHori,
		2*Margin+NumRows*FieldHeight+NumRows*ColorPaddingVert)
	gc.SetFontFace(face)
	gc.SetFillColor(color.White)
	gc.Clear()
	for i, name := range PaletteNames {
		col, row := i/NumRows, i%NumRows
		pal := PaletteMap[name]
		x := Margin + float64(col)*(5*FieldWidth+4*FieldPadding+ColorPaddingHori)
		y := Margin + float64(row)*(FieldHeight+ColorPaddingVert)
		for n := range 101 {
			t := float64(n) / float64(100)
			x := x + float64(n)*(5.0)
			gc.SetFillColor(pal.Color(t))
			gc.SetStrokeWidth(0.0)
			gc.DrawRectangle(x, y, 5.0, FieldHeight)
			gc.Fill()
		}
		gc.SetStrokeColor(color.Black)
		gc.DrawStringAnchored(pal.Name(), x, y+FieldHeight+Margin/2, 0.0, 1.0)
	}
	gc.SavePNG("palette_test.png")
}
