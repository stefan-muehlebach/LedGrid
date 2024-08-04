package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

const (
	defWidth  = 40
	defHeight = 10
	defHost   = "raspi-3"
	defPort   = 5333
)

var (
	width, height int
	gridSize      image.Point
	refreshRate        = 30 * time.Millisecond
	backAlpha          = 1.0
	animCtrl *AnimationController
)

//----------------------------------------------------------------------------

func RegularPolygonTest(c *Canvas) {
	posList := []geom.Point{
		ConvertPos(geom.Point{-5.5, 4.5}),
		ConvertPos(geom.Point{44.5, 4.5}),
	}
	posCenter := ConvertPos(geom.Point{19.5, 4.5})
	smallSize := ConvertSize(geom.Point{9.0, 9.0})
	largeSize := ConvertSize(geom.Point{80.0, 80.0})

	polyList := make([]*RegularPolygon, 9)

	aSeq := NewSequence()
	for n := 3; n <= 6; n++ {
		col := colornames.RandColor()
		polyList[n] = NewRegularPolygon(n, posList[n%2], smallSize, col)
		c.Add(polyList[n])
		dur := 2*time.Second + rand.N(time.Second)
		sign := []float64{+1.0, -1.0}[n%2]
		angle := sign * 2 * math.Pi
		animPos := NewPositionAnimation(&polyList[n].Pos, posCenter, dur)
		animAngle := NewFloatAnimation(&polyList[n].Angle, angle, dur)
		animSize := NewSizeAnimation(&polyList[n].Size, largeSize, 4*time.Second)
		animSize.Cont = true
		animFade := NewColorAnimation(&polyList[n].BorderColor, colornames.Black, 4*time.Second)
		animFade.Cont = true
		// animFade := NewAlphaAnimation(&polyList[n].FillColor.A, 0x00, time.Second)

		aGrpIn := NewGroup(animPos, animAngle)
		// aGrp.duration = dur
		aGrpOut := NewGroup(animSize, animFade)
		aObjSeq := NewSequence(aGrpIn, aGrpOut)
		aSeq.Add(aObjSeq)
	}
	aSeq.RepeatCount = AnimationRepeatForever
	aSeq.Start()
}

func GroupTest(ctrl *Canvas) {
	rPos1 := ConvertPos(geom.Point{4.5, 4.5})
	rPos2 := ConvertPos(geom.Point{float64(width) - 1.5, 4.5})
	rSize1 := ConvertSize(geom.Point{7.0, 7.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})
	rColor1 := colornames.SkyBlue
	rColor2 := colornames.GreenYellow

	r := NewRectangle(rPos1, rSize1, rColor1)
	ctrl.Add(r)

	aPos := NewPositionAnimation(&r.Pos, rPos2, time.Second)
	aPos.AutoReverse = true
	aSize := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aSize.AutoReverse = true
	aColor := NewColorAnimation(&r.BorderColor, rColor2, time.Second)
	aColor.AutoReverse = true
	aAngle := NewFloatAnimation(&r.Angle, math.Pi, time.Second)
	aAngle.AutoReverse = true

	aGroup := NewGroup(aPos, aSize, aColor, aAngle)
	aGroup.RepeatCount = AnimationRepeatForever

	aGroup.Start()
}

func SequenceTest(ctrl *Canvas) {
	rPos := ConvertPos(geom.NewPointIMG(gridSize).Mul(0.5).SubXY(0.5, 0.5))
	rSize1 := ConvertSize(geom.NewPointIMG(gridSize).SubXY(1, 1))
	rSize2 := ConvertSize(geom.Point{5.0, 3.0})

	r := NewRectangle(rPos, rSize1, colornames.SkyBlue)
	ctrl.Add(r)

	aSize1 := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aColor1 := NewColorAnimation(&r.BorderColor, colornames.OrangeRed, time.Second/2)
	aColor1.AutoReverse = true
	aColor2 := NewColorAnimation(&r.BorderColor, colornames.Crimson, time.Second/2)
	aColor2.AutoReverse = true
	aColor3 := NewColorAnimation(&r.BorderColor, colornames.Coral, time.Second/2)
	aColor3.AutoReverse = true
	aColor4 := NewColorAnimation(&r.BorderColor, colornames.FireBrick, time.Second/2)
	aSize2 := NewSizeAnimation(&r.Size, rSize1, time.Second)
	aSize2.Cont = true
	aColor5 := NewColorAnimation(&r.BorderColor, colornames.SkyBlue, time.Second)
	aColor5.Cont = true

	aSeq := NewSequence(aSize1, aColor1, aColor2, aColor3, aColor4, aSize2, aColor5)
	aSeq.RepeatCount = AnimationRepeatForever
	aSeq.Start()
}

func TimelineTest(ctrl *Canvas) {
	r1Pos := ConvertPos(geom.Point{6.5, float64(height-1) / 2.0})
	r1Size := ConvertSize(geom.Point{9.0, 5.0})
	r2Pos := ConvertPos(geom.Point{float64(width-1) / 2.0, float64(height-1) / 2.0})
	r2Size := ConvertSize(geom.Point{11.0, 7.0})
	r3Pos := ConvertPos(geom.Point{float64(width) - 7.5, float64(height-1) / 2.0})
	r3Size := ConvertSize(geom.Point{9.0, 5.0})

	r1 := NewRectangle(r1Pos, r1Size, colornames.GreenYellow)
	r2 := NewRectangle(r2Pos, r2Size, colornames.Gold)
	r3 := NewRectangle(r3Pos, r3Size, colornames.SkyBlue)
	ctrl.Add(r1, r2, r3)

	aAngle1 := NewFloatAnimation(&r1.Angle, math.Pi, time.Second)
	aAngle2 := NewFloatAnimation(&r1.Angle, 0.0, time.Second)
	aAngle2.Cont = true

	aColor1 := NewColorAnimation(&r1.BorderColor, colornames.OrangeRed, 200*time.Millisecond)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 3
	aColor2 := NewColorAnimation(&r1.BorderColor, colornames.Purple, 500*time.Millisecond)
	aColor3 := NewColorAnimation(&r1.BorderColor, colornames.GreenYellow, 500*time.Millisecond)
	aColor3.Cont = true

	aPos1 := NewPositionAnimation(&r1.Pos, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
	aPos1.AutoReverse = true

	aAngle3 := NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
	aAngle4 := NewFloatAnimation(&r3.Angle, 0.0, time.Second)
	aAngle4.Cont = true

	aColor4 := NewColorAnimation(&r3.BorderColor, colornames.DarkOrange, 200*time.Millisecond)
	aColor4.AutoReverse = true
	aColor4.RepeatCount = 3
	aColor5 := NewColorAnimation(&r3.BorderColor, colornames.Purple, 500*time.Millisecond)
	aColor6 := NewColorAnimation(&r3.BorderColor, colornames.SkyBlue, 500*time.Millisecond)
	aColor6.Cont = true

	aPos2 := NewPositionAnimation(&r3.Pos, r2.Pos.AddXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
	aPos2.AutoReverse = true

	aColor7 := NewColorAnimation(&r2.BorderColor, colornames.Cornsilk, 500*time.Millisecond)
	aColor7.AutoReverse = true
	aBorder := NewFloatAnimation(&r2.BorderWidth, ConvertLen(3.0), 500*time.Millisecond)
	aBorder.AutoReverse = true

	tl := NewTimeline(5 * time.Second)
	tl.RepeatCount = AnimationRepeatForever

	// Timeline positions for the first rectangle
	tl.Add(300*time.Millisecond, aColor1)
	tl.Add(1800*time.Millisecond, aAngle1)
	tl.Add(2300*time.Millisecond, aColor2, aPos1)
	tl.Add(2900*time.Millisecond, aAngle2)
	tl.Add(3400*time.Millisecond, aColor3)

	// Timeline positions for the second rectangle
	tl.Add(500*time.Millisecond, aColor4)
	tl.Add(2000*time.Millisecond, aAngle3)
	tl.Add(2500*time.Millisecond, aColor5, aPos2)
	tl.Add(3100*time.Millisecond, aAngle4)
	tl.Add(3600*time.Millisecond, aColor6)

	tl.Add(2400*time.Millisecond, aColor7, aBorder)

	tl.Start()
}

func PathTest(ctrl *Canvas) {
	duration := 4 * time.Second
	pathA := FullCirclePathA
	pathB := FullCirclePathB

	pos1 := ConvertPos(geom.Point{1.5, float64(height-1) / 2.0})
	pos2 := ConvertPos(geom.Point{float64(width-1) / 2.0, 1.5})
	pos3 := ConvertPos(geom.Point{float64(width-1) - 1.5, float64(height-1) / 2.0})
	pos4 := ConvertPos(geom.Point{float64(width-1) / 2.0, float64(height-1) - 1.5})
	cSize := ConvertSize(geom.Point{3.0, 3.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)
	ctrl.Add(c1, c2, c3, c4)

	c1Path := NewPathAnimation(&c1.Pos, pathB, ConvertSize(geom.Point{float64(width - 4), 6.0}), duration)
	c1Path.AutoReverse = true
	c3Path := NewPathAnimation(&c3.Pos, pathB, ConvertSize(geom.Point{-float64(width - 4), -6.0}), duration)
	c3Path.AutoReverse = true

	c2Path := NewPathAnimation(&c2.Pos, pathA, ConvertSize(geom.Point{float64(width) / 3.0, 6.0}), duration)
	c2Path.AutoReverse = true
	c4Path := NewPathAnimation(&c4.Pos, pathA, ConvertSize(geom.Point{-float64(width) / 3.0, -6.0}), duration)
	c4Path.AutoReverse = true

	aGrp := NewGroup(c1Path, c2Path, c3Path, c4Path)
	aGrp.RepeatCount = AnimationRepeatForever
	aGrp.Start()
}

func PolygonPathTest(ctrl *Canvas) {

	cPos := ConvertPos(geom.Point{1, 1})
	cSize := ConvertSize(geom.Point{2, 2})

	polyPath1 := NewPolygonPath(
		ConvertPos(geom.Point{1, 1}),
		ConvertPos(geom.Point{float64(width) - 2, 1}),
		ConvertPos(geom.Point{float64(width) - 2, float64(height) - 2}),
		ConvertPos(geom.Point{1, float64(height) - 2}),

		ConvertPos(geom.Point{1, 2}),
		ConvertPos(geom.Point{float64(width) - 3, 2}),
		ConvertPos(geom.Point{float64(width) - 3, float64(height) - 3}),
		ConvertPos(geom.Point{2, float64(height) - 3}),

		ConvertPos(geom.Point{2, 3}),
		ConvertPos(geom.Point{float64(width) - 4, 3}),
		ConvertPos(geom.Point{float64(width) - 4, float64(height) - 4}),
		ConvertPos(geom.Point{3, float64(height) - 4}),

		ConvertPos(geom.Point{3, 4}),
		ConvertPos(geom.Point{float64(width) - 5, 4}),
		ConvertPos(geom.Point{float64(width) - 5, float64(height) - 5}),
		ConvertPos(geom.Point{4, float64(height) - 5}),
	)

	polyPath2 := NewPolygonPath(
		ConvertPos(geom.Point{1, 1}),
		ConvertPos(geom.Point{4, 8}),
		ConvertPos(geom.Point{7, 2}),
		ConvertPos(geom.Point{10, 7}),
		ConvertPos(geom.Point{13, 3}),
		ConvertPos(geom.Point{16, 6}),
		ConvertPos(geom.Point{19, 4}),
		ConvertPos(geom.Point{22, 5}),
	)

	c1 := NewEllipse(cPos, cSize, colornames.GreenYellow)
	ctrl.Add(c1)

	aPath1 := NewPolyPathAnimation(&c1.Pos, polyPath1, 7*time.Second)
	aPath1.AutoReverse = true

	aPath2 := NewPolyPathAnimation(&c1.Pos, polyPath2, 7*time.Second)
	aPath2.AutoReverse = true

	seq := NewSequence(aPath1, aPath2)
	seq.RepeatCount = AnimationRepeatForever

	seq.Start()
}

func RandomWalk(ctrl *Canvas) {
	rect := geom.Rectangle{Min: ConvertPos(geom.Point{1.0, 1.0}), Max: ConvertPos(geom.Point{float64(width) - 1.0, float64(height) - 1.0})}
	pos1 := ConvertPos(geom.Point{1.0, 1.0})
	pos2 := ConvertPos(geom.Point{18.0, 8.0})
	size1 := ConvertSize(geom.Point{2.0, 2.0})
	size2 := ConvertSize(geom.Point{4.0, 4.0})

	c1 := NewEllipse(pos1, size1, colornames.SkyBlue)
	c2 := NewEllipse(pos2, size2, colornames.GreenYellow)
	ctrl.Add(c1, c2)

	aPos1 := NewPositionAnimation(&c1.Pos, geom.Point{}, 1300*time.Millisecond)
	aPos1.Cont = true
	aPos1.ValFunc = RandPointTrunc(rect, 1.0)
	aPos1.RepeatCount = AnimationRepeatForever

	aPos2 := NewPositionAnimation(&c2.Pos, geom.Point{}, 901*time.Millisecond)
	aPos2.Cont = true
	aPos2.ValFunc = RandPoint(rect)
	aPos2.RepeatCount = AnimationRepeatForever

	aPos1.Start()
	aPos2.Start()
}

//----------------------------------------------------------------------------

type BouncingEllipse struct {
	Ellipse
	Vel, Acc geom.Point
	Field    geom.Rectangle
}

func NewBouncingEllipse(pos, size geom.Point, col ledgrid.LedColor) *BouncingEllipse {
	b := &BouncingEllipse{}
	b.Ellipse = *NewEllipse(pos, size, col)
	b.Vel = geom.Point{}
	b.Acc = geom.Point{}
	return b
}

func (b *BouncingEllipse) Update(pit time.Time) bool {
	deltaVel := b.Acc.Mul(0.3)
	b.Vel = b.Vel.Add(deltaVel)
	b.Pos = b.Pos.Add(b.Vel)
	if b.Pos.X < b.Field.Min.X || b.Pos.X >= b.Field.Max.X {
		b.Vel.X = -b.Vel.X
	}
	if b.Pos.Y < b.Field.Min.Y || b.Pos.Y >= b.Field.Max.Y {
		b.Vel.Y = -b.Vel.Y
	}
	return true
}

func (b *BouncingEllipse) Duration() time.Duration {
	return time.Duration(0)
}
func (b *BouncingEllipse) SetDuration(dur time.Duration) {}
func (b *BouncingEllipse) Start()                        {}
func (b *BouncingEllipse) Stop()                         {}
func (b *BouncingEllipse) Continue()                     {}
func (b *BouncingEllipse) IsStopped() bool {
	return false
}

func BounceAround(c *Canvas) {
	pos1 := ConvertPos(geom.Point{2.0, 2.0})
	pos2 := ConvertPos(geom.Point{37.0, 7.0})
	size := ConvertSize(geom.Point{4.0, 4.0})
	vel1 := geom.Point{1.5, 0.75}
	vel2 := geom.Point{-3.5, -2.5}
	// acc := ConvertSize(geom.Point{0, 0.1})

	obj1 := NewBouncingEllipse(pos1, size, colornames.GreenYellow)
	obj1.Vel = vel1
	obj1.Field = geom.NewRectangleIMG(c.canvas.Bounds())
	obj2 := NewBouncingEllipse(pos2, size, colornames.LightSeaGreen)
	obj2.Vel = vel2
	obj2.Field = geom.NewRectangleIMG(c.canvas.Bounds())

	c.Add(obj1, obj2)
	animCtrl.AddAnim(obj1, obj2)
}

//----------------------------------------------------------------------------

func Piiiiixels(ctrl *Canvas) {
	dPosX := ConvertSize(geom.Point{2.0, 0.0})
	dPosY := ConvertSize(geom.Point{0.0, 2.0})
	p1Pos1 := ConvertPos(geom.Point{1.0, 1.0})

	aGrp := NewGroup()
	aGrp.RepeatCount = AnimationRepeatForever
	for i := range 5 {
		for j := range 5 {
			palName := ledgrid.PaletteNames[5*i+j]
			pos := p1Pos1.Add(dPosX.Mul(float64(j)).Add(dPosY.Mul(float64(i))))
			pix := NewPixel(pos, colornames.OrangeRed)
			ctrl.Add(pix)

			pixColor := NewPaletteAnimation(&pix.Color, ledgrid.PaletteMap[palName], 4*time.Second)
			aGrp.Add(pixColor)
		}
	}
	aGrp.Start()
}

func CirclingCircles(ctrl *Canvas) {
	pos1 := ConvertPos(geom.Point{1.0, 1.0})
	pos2 := ConvertPos(geom.Point{10.0, 8.0})
	pos3 := ConvertPos(geom.Point{19.0, 1.0})
	pos4 := ConvertPos(geom.Point{28.0, 8.0})
	pos5 := ConvertPos(geom.Point{37.0, 1.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)
	c5 := NewEllipse(pos5, cSize, colornames.YellowGreen)

	stepRD := ConvertSize(geom.Point{9.0, 7.0})
	stepLU := stepRD.Neg()
	stepRU := ConvertSize(geom.Point{9.0, -7.0})
	stepLD := stepRU.Neg()

	c1Path1 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, stepRD, time.Second)
	c1Path1.Cont = true
	c1Path2 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, stepRU, time.Second)
	c1Path2.Cont = true
	c1Path3 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, stepLD, time.Second)
	c1Path3.Cont = true
	c1Path4 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, stepLU, time.Second)
	c1Path4.Cont = true

	c2Path1 := NewPathAnimation(&c2.Pos, QuarterCirclePathA, stepLU, time.Second)
	c2Path1.Cont = true
	c2Path2 := NewPathAnimation(&c2.Pos, QuarterCirclePathA, stepRD, time.Second)
	c2Path2.Cont = true

	c3Path1 := NewPathAnimation(&c3.Pos, QuarterCirclePathA, stepLD, time.Second)
	c3Path1.Cont = true
	c3Path2 := NewPathAnimation(&c3.Pos, QuarterCirclePathA, stepRU, time.Second)
	c3Path2.Cont = true

	c4Path1 := NewPathAnimation(&c4.Pos, QuarterCirclePathA, stepLU, time.Second)
	c4Path1.Cont = true
	c4Path2 := NewPathAnimation(&c4.Pos, QuarterCirclePathA, stepRD, time.Second)
	c4Path2.Cont = true

	c5Path1 := NewPathAnimation(&c5.Pos, QuarterCirclePathA, stepLD, time.Second)
	c5Path1.Cont = true
	c5Path2 := NewPathAnimation(&c5.Pos, QuarterCirclePathA, stepRU, time.Second)
	c5Path2.Cont = true

	aGrp1 := NewGroup(c1Path1, c2Path1)
	aGrp2 := NewGroup(c1Path2, c3Path1)
	aGrp3 := NewGroup(c1Path1, c4Path1)
	aGrp4 := NewGroup(c1Path2, c5Path1)
	aGrp5 := NewGroup(c1Path3, c5Path2)
	aGrp6 := NewGroup(c1Path4, c4Path2)
	aGrp7 := NewGroup(c1Path3, c3Path2)
	aGrp8 := NewGroup(c1Path4, c2Path2)
	aSeq := NewSequence(aGrp1, aGrp2, aGrp3, aGrp4, aGrp5, aGrp6, aGrp7, aGrp8)
	aSeq.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1, c2, c3, c4, c5)
	aSeq.Start()
}

func ChasingCircles(ctrl *Canvas) {
	c1Pos1 := ConvertPos(geom.Point{36.5, 4.5})
	c1Size1 := ConvertSize(geom.Point{10.0, 10.0})
	c1Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c1PosSize := ConvertSize(geom.Point{-34.0, -5.0})
	c2Pos := ConvertPos(geom.Point{2.5, 4.5})
	c2Size1 := ConvertSize(geom.Point{5.0, 5.0})
	c2Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c2PosSize := ConvertSize(geom.Point{34.0, 7.0})

	aGrp := NewGroup()

	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.LedColorModel.Convert(colornames.DeepSkyBlue).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Lime).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Teal).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.SkyBlue).(ledgrid.LedColor),
	)

	c1 := NewEllipse(c1Pos1, c1Size1, colornames.Gold)

	c1pos := NewPathAnimation(&c1.Pos, FullCirclePathB, c1PosSize, 4*time.Second)
	c1pos.RepeatCount = AnimationRepeatForever
	c1pos.Curve = AnimationLinear

	c1size := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1size.AutoReverse = true
	c1size.RepeatCount = AnimationRepeatForever

	c1bcolor := NewColorAnimation(&c1.BorderColor, colornames.OrangeRed, time.Second)
	c1bcolor.AutoReverse = true
	c1bcolor.RepeatCount = AnimationRepeatForever

	aGrp.Add(c1pos, c1size, c1bcolor)

	c2 := NewEllipse(c2Pos, c2Size1, colornames.Lime)

	c2pos := NewPathAnimation(&c2.Pos, FullCirclePathB, c2PosSize, 4*time.Second)
	c2pos.RepeatCount = AnimationRepeatForever
	c2pos.Curve = AnimationLinear

	c2size := NewSizeAnimation(&c2.Size, c2Size2, time.Second)
	c2size.AutoReverse = true
	c2size.RepeatCount = AnimationRepeatForever

	c2color := NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
	c2color.RepeatCount = AnimationRepeatForever
	c2color.Curve = AnimationLinear

	aGrp.Add(c2pos, c2size, c2color)

	ctrl.Add(c2, c1)
	aGrp.Start()
}

func CircleAnimation(ctrl *Canvas) {
	c1Pos1 := ConvertPos(geom.Point{1.5, 4.5})
	c1Pos3 := ConvertPos(geom.Point{37.5, 4.5})

	c1Size1 := ConvertSize(geom.Point{3.0, 3.0})
	c1Size2 := ConvertSize(geom.Point{9.0, 9.0})

	c1 := NewEllipse(c1Pos1, c1Size1, colornames.OrangeRed)

	c1pos := NewPositionAnimation(&c1.Pos, c1Pos3, 2*time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = AnimationRepeatForever

	c1radius := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = AnimationRepeatForever

	c1color := NewColorAnimation(&c1.BorderColor, colornames.Gold, 4*time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
}

func PushingRectangles(ctrl *Canvas) {
	rSize1 := ConvertSize(geom.Point{float64(width) - 3.0, 1.0})
	rSize2 := ConvertSize(geom.Point{1.0, float64(height) - 1.0})

	r1Pos1 := ConvertPos(geom.Point{0.5, float64(height-1) / 2.0})
	r1Pos2 := ConvertPos(geom.Point{0.5 + float64(width-3)/2.0, float64(height-1) / 2.0})

	r2Pos1 := ConvertPos(geom.Point{float64(width-1) - 0.5, float64(height-1) / 2.0})
	r2Pos2 := ConvertPos(geom.Point{float64(width-1) - 0.5 - float64(width-3)/2.0, float64(height-1) / 2.0})
	duration := 2 * time.Second

	r1 := NewRectangle(r1Pos1, rSize2, colornames.Crimson)

	a1Pos := NewPositionAnimation(&r1.Pos, r1Pos2, duration)
	a1Pos.AutoReverse = true
	a1Pos.RepeatCount = AnimationRepeatForever

	a1Size := NewSizeAnimation(&r1.Size, rSize1, duration)
	a1Size.AutoReverse = true
	a1Size.RepeatCount = AnimationRepeatForever

	a1Color := NewColorAnimation(&r1.BorderColor, colornames.GreenYellow, duration)
	a1Color.AutoReverse = true
	a1Color.RepeatCount = AnimationRepeatForever

	r2 := NewRectangle(r2Pos2, rSize1, colornames.SkyBlue)

	a2Pos := NewPositionAnimation(&r2.Pos, r2Pos1, duration)
	a2Pos.AutoReverse = true
	a2Pos.RepeatCount = AnimationRepeatForever

	a2Size := NewSizeAnimation(&r2.Size, rSize2, duration)
	a2Size.AutoReverse = true
	a2Size.RepeatCount = AnimationRepeatForever

	a2Color := NewColorAnimation(&r2.BorderColor, colornames.Crimson, duration)
	a2Color.AutoReverse = true
	a2Color.RepeatCount = AnimationRepeatForever

	ctrl.Add(r1, r2)
	a1Pos.Start()
	a1Size.Start()
	a1Color.Start()
	a2Pos.Start()
	a2Size.Start()
	a2Color.Start()
}

func pixIdx(x, y int) int {
	return y*width + x
}

func pixCoord(idx int) (x, y int) {
	return idx % width, idx / width
}

func GlowingPixels(ctrl *Canvas) {
	var pixField []*Pixel

	pixField = make([]*Pixel, width*height)

	for idx, _ := range pixField {
		x, y := pixCoord(idx)
		pos := ConvertPos(geom.Point{float64(x), float64(y)})
		col := colornames.DimGray.Interpolate(colornames.DarkGray, rand.Float64())

		pix := NewPixel(pos, col)
		pixField[idx] = pix
		ctrl.Add(pix)

		dur := time.Second + rand.N(400*time.Millisecond)
		aAlpha := NewAlphaAnimation(&pix.Color.A, 127, dur)
		aAlpha.AutoReverse = true
		aAlpha.RepeatCount = AnimationRepeatForever
		aAlpha.Start()
	}

	aGrpPurple := NewGroup()
	aGrpYellow := NewGroup()
	aGrpGreen := NewGroup()

	for _, pix := range pixField {
		aColor := NewColorAnimation(&pix.Color, colornames.MediumPurple.Interpolate(colornames.Fuchsia, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpPurple.Add(aColor)

		aColor = NewColorAnimation(&pix.Color, colornames.Gold.Interpolate(colornames.Khaki, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpYellow.Add(aColor)

		aColor = NewColorAnimation(&pix.Color, colornames.GreenYellow.Interpolate(colornames.LightSeaGreen, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpGreen.Add(aColor)
	}

	aTimel := NewTimeline(40 * time.Second)
	aTimel.Add(10*time.Second, aGrpPurple)
	aTimel.Add(20*time.Second, aGrpYellow)
	aTimel.Add(30*time.Second, aGrpGreen)
	aTimel.RepeatCount = AnimationRepeatForever

	aTimel.Start()
}

var (
	pts = []geom.Point{
		ConvertPos(geom.Point{0, 0}),
		ConvertPos(geom.Point{0, float64(height)}),
		ConvertPos(geom.Point{float64(width), float64(height)}),
		ConvertPos(geom.Point{float64(width), 0}),
	}
	lastP0 = 0
)

func randPoint() geom.Point {
	p0 := rand.IntN(len(pts))
	for p0 == lastP0 {
		p0 = rand.IntN(len(pts))
	}
	lastP0 = p0
	p1 := (p0 + 1) % len(pts)

	return pts[p0].Interpolate(pts[p1], rand.Float64())
}

func MovingText(c *Canvas) {

	t1 := NewText(randPoint(), "Beni", colornames.LightSeaGreen)
	t2 := NewText(randPoint(), "Stefan", colornames.YellowGreen)
	c.Add(t1, t2)

	aPos1 := NewPositionAnimation(&t1.Pos, geom.Point{}, 4*time.Second)
	aPos1.ValFunc = randPoint
	aPos1.RepeatCount = AnimationRepeatForever
	aPos1.Cont = true

	aPos2 := NewPositionAnimation(&t2.Pos, geom.Point{}, 3*time.Second)
	aPos2.ValFunc = randPoint
	aPos2.RepeatCount = AnimationRepeatForever
	aPos2.Cont = true

	aAngle1 := NewFloatAnimation(&t1.Angle, 0.0, 3*time.Second)
	aAngle1.ValFunc = RandFloat(math.Pi/2.0, math.Pi)
	aAngle1.AutoReverse = true
	aAngle1.RepeatCount = AnimationRepeatForever

	aAngle2 := NewFloatAnimation(&t2.Angle, 0.0, 4*time.Second)
	aAngle2.ValFunc = RandFloat(math.Pi/6.0, math.Pi/2.0)
	aAngle2.AutoReverse = true
	aAngle2.RepeatCount = AnimationRepeatForever

	aAngle1.Start()
	aAngle2.Start()
	aPos1.Start()
	aPos2.Start()
}

func FlyingImages(c *Canvas) {
	pos1 := ConvertPos(geom.Point{-7.5, 4.5})
	pos2 := ConvertPos(geom.Point{14.5, 4.5})
	i4 := NewImageFromFile(pos1, "images/lorem.png")
	i4.Size = ConvertSize(geom.Point{16.0, 3.0})
	c.Add(i4)

	aPos1 := NewPositionAnimation(&i4.Pos, pos2, 4*time.Second)
	aSize1 := NewSizeAnimation(&i4.Size, ConvertSize(geom.Point{633.0, 118.0}), 4*time.Second)
	aPos2 := NewPositionAnimation(&i4.Pos, pos2.Add(geom.Point{250.0, 0.0}), 6*time.Second)
	aPos2.Cont = true
	aSize2 := NewSizeAnimation(&i4.Size, ConvertSize(geom.Point{16.0, 3.0}), 4*time.Second)
	aSize2.Cont = true
	aPos3 := NewPositionAnimation(&i4.Pos, pos1, 3*time.Second)
	aPos3.Cont = true

	aSeq := NewSequence(aPos1, aSize1, aPos2, aSize2)
	aSeq.Start()
}

func CameraTest(c *Canvas) {
	pos := ConvertPos(geom.Point{float64(width) / 2.0, float64(height) / 2.0})
	size := ConvertSize(geom.Point{float64(width), float64(height)})

	cam := NewCamera(pos, size)
	c.Add(cam)

	cam.Start()
}

func BlinkenAnimation(c *Canvas) {
	posA := ConvertPos(geom.Point{0.5, 0.5})
	posB := ConvertPos(geom.Point{5.5, 5.5})

	bml := ReadBlinkenFile("marioWalkRight.bml")
	img := bml.Image(0)
	img.Pos = posA

	aPos := NewPositionAnimation(&img.Pos, posB, 3*time.Second)
	aPos.AutoReverse = true
	aPos.RepeatCount = AnimationRepeatForever

	c.Add(img)
	aPos.Start()
}

//-----------------------------------------------------------------------------
// Folgende Animationen sind mit dem wesentlich schnelleren aber graphisch
// etwas eingeschraenkteren Grid-Typ erstellt.

func GlowingGridPixels(g *Grid) {
	aGrpPurple := NewGroup()
	aGrpYellow := NewGroup()
	aGrpGreen := NewGroup()
	aGrpGrey := NewGroup()

	for y := range g.ledGrid.Rect.Dy() {
		for x := range g.ledGrid.Rect.Dx() {
			pos := image.Point{x, y}
			t := rand.Float64()
			col := (colornames.DimGray.Dark(0.2)).Interpolate((colornames.DarkGrey.Dark(0.2)), t)
			pix := NewGridPixel(pos, col)
			g.Add(pix)

			dur := time.Second + rand.N(time.Second)
			aAlpha := NewAlphaAnimation(&pix.Color.A, 200, dur)
			aAlpha.AutoReverse = true
			aAlpha.RepeatCount = AnimationRepeatForever
			aAlpha.Start()

			aColor := NewColorAnimation(&pix.Color, col, 1*time.Second)
			aColor.Cont = true
			aGrpGrey.Add(aColor)

			aColor = NewColorAnimation(&pix.Color, colornames.MediumPurple.Interpolate(colornames.Fuchsia, t), 5*time.Second)
			aColor.Cont = true
			// aColor.AutoReverse = true
			aGrpPurple.Add(aColor)

			aColor = NewColorAnimation(&pix.Color, colornames.Gold.Interpolate(colornames.Khaki, t), 5*time.Second)
			aColor.Cont = true
			aGrpYellow.Add(aColor)

			aColor = NewColorAnimation(&pix.Color, colornames.GreenYellow.Interpolate(colornames.LightSeaGreen, t), 5*time.Second)
			aColor.Cont = true
			aGrpGreen.Add(aColor)
		}
	}

	txt1 := NewGridText(fixed.P(width/2, height/2), colornames.GreenYellow.Alpha(0.0), "REDEN")
	aTxt1 := NewAlphaAnimation(&txt1.Color.A, 255, 2*time.Second)
	aTxt1.AutoReverse = true
	txt2 := NewGridText(fixed.P(width/2, height/2), colornames.DarkViolet.Alpha(0.0), "DENKEN")
	aTxt2 := NewAlphaAnimation(&txt2.Color.A, 255, 2*time.Second)
	aTxt2.AutoReverse = true
	txt3 := NewGridText(fixed.P(width/2, height/2), colornames.OrangeRed.Alpha(0.0), "LACHEN")
	aTxt3 := NewAlphaAnimation(&txt3.Color.A, 255, 2*time.Second)
	aTxt3.AutoReverse = true
	g.Add(txt1, txt2, txt3)

	aTimel := NewTimeline(42 * time.Second)
	aTimel.Add(7*time.Second, aGrpPurple)
	aTimel.Add(12*time.Second, aTxt1)
	aTimel.Add(13*time.Second, aGrpGrey)

	aTimel.Add(22*time.Second, aGrpYellow)
	aTimel.Add(27*time.Second, aTxt2)
	aTimel.Add(28*time.Second, aGrpGrey)

	aTimel.Add(35*time.Second, aGrpGreen)
	aTimel.Add(40*time.Second, aTxt3)
	aTimel.Add(41*time.Second, aGrpGrey)
	aTimel.RepeatCount = AnimationRepeatForever

	aTimel.Start()
}

func RandomGridPixels(g *Grid) {
	for y := range g.ledGrid.Rect.Dy() {
		for x := range g.ledGrid.Rect.Dx() {
			pos := image.Pt(x, y)
			colorGrp1 := colornames.ColorGroup(x/3) % colornames.NumColorGroups
			colorGrp2 := (colorGrp1 + 1) % colornames.NumColorGroups
			col := colornames.RandGroupColor(colorGrp1)
			pix := NewGridPixel(pos, col)
			g.Add(pix)
			dur := time.Second
			aColor := NewColorAnimation(&pix.Color, colornames.RandGroupColor(colorGrp2), dur)
			aColor.AutoReverse = true
			aColor.RepeatCount = AnimationRepeatForever
			aColor.Start()
		}
	}
}

func TextOnGrid(g *Grid) {
	basePt := fixed.P(0, 5)
	baseColor1 := colornames.SkyBlue

	txt1 := NewGridText(basePt, baseColor1, "STEFAN")
	g.Add(txt1)

	aPos := NewFixedPosAnimation(&txt1.Pos, fixed.P(25, 5), 3*time.Second)
	aPos.AutoReverse = true
	aPos.RepeatCount = AnimationRepeatForever
	aPos.Start()
}

func WalkingPixelOnGrid(g *Grid) {
	pos := image.Point{0, 0}
	col := colornames.GreenYellow
	pix := NewGridPixel(pos, col)
	g.Add(pix)

	go func() {
		idx := 0
		for {
			col := idx % width
			row := idx / width
			pix.Pos = image.Point{col, row}
			time.Sleep(time.Second / 5)
			idx++
		}
	}()
}

func ImagesOnGrid(g *Grid) {
	pos := image.Point{5, 2}
	size := image.Point{10, 10}

	img := NewGridImage(pos, size)
	for row := range size.Y {
		for col := range size.X {
			img.Img.SetRGBA(col, row, color.RGBA{0x8f, 0x8f, 0x8f, 0xff})
		}
	}

	g.Add(img)
}

//----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

type canvasSceneRecord struct {
	name string
	fnc  func(canvas *Canvas)
}

type gridSceneRecord struct {
	name string
	fnc  func(grid *Grid)
}

func main() {
	var host string
	var port uint
	var input string
	var ch byte
	var sceneId int
	var runInteractive bool
	var pixCtrl ledgrid.PixelClient

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.StringVar(&input, "scene", input, "play one single scene (no menu)")
	flag.BoolVar(&doLog, "log", doLog, "enable logging")
	flag.Parse()

	gridSize = image.Point{width, height}

	if len(input) > 0 {
		runInteractive = false
		ch = input[0]
	} else {
		runInteractive = true
	}

	canvasSceneList := []canvasSceneRecord{
		{"Group test", GroupTest},
		{"Sequence test", SequenceTest},
		{"Timeline test", TimelineTest},
		{"Path test", PathTest},
		{"PolygonPath Test", PolygonPathTest},
		{"Random walk", RandomWalk},
		{"Let's bounce around", BounceAround},
		{"Piiiiixels", Piiiiixels},
		{"Regular Polygon test", RegularPolygonTest},
		{"Circling Circles", CirclingCircles},
		{"Chasing Circles", ChasingCircles},
		{"Circle Animation", CircleAnimation},
		{"Pushing Rectangles", PushingRectangles},
		{"Glowing Pixels", GlowingPixels},
		{"Moving Text", MovingText},
		{"Flying images", FlyingImages},
		{"Hidden or visible camera", CameraTest},
		{"Animation from a BlinkenLight file", BlinkenAnimation},
	}

	gridSceneList := []gridSceneRecord{
		{"Glowing Pixels", GlowingGridPixels},
		{"Random Pixels", RandomGridPixels},
		{"Text on a grid", TextOnGrid},
		{"Walking pixel", WalkingPixelOnGrid},
		{"Draw images on a grid", ImagesOnGrid},
	}

	pixCtrl = ledgrid.NewNetPixelClient(host, port)

	ledGrid := ledgrid.NewLedGrid(gridSize, nil)

    animCtrl = NewAnimationController(refreshRate)
    animCtrl.Stop()

	canvas := NewCanvas(pixCtrl, ledGrid)
	grid := NewGrid(pixCtrl, ledGrid)

    animCtrl.DrawArea = canvas

	if runInteractive {
		sceneId = -1
		for {
			fmt.Printf("Animations:\n")
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("Canvas based\n")
			fmt.Printf("---------------------------------------\n")
			for i, scene := range canvasSceneList {
				if ch >= 'a' && ch <= 'z' && i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("Grid based\n")
			fmt.Printf("---------------------------------------\n")
			for i, scene := range gridSceneList {
				if ch >= 'A' && ch <= 'Z' && i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'A'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")
			fmt.Scanf("%s\n", &input)
			ch = input[0]
			if ch == '0' {
				break
			}
			if ch >= 'a' && ch <= 'z' {
				sceneId = int(ch - 'a')
				if sceneId < 0 || sceneId >= len(canvasSceneList) {
					continue
				}
				animCtrl.Stop()
				animCtrl.DelAllAnim()
				animCtrl.DrawArea = canvas
				animCtrl.Continue()
				canvas.DelAll()
				canvasSceneList[sceneId].fnc(canvas)
			}
			if ch >= 'A' && ch <= 'Z' {
				sceneId = int(ch - 'A')
				if sceneId < 0 || sceneId >= len(gridSceneList) {
					continue
				}
				animCtrl.Stop()
				animCtrl.DelAllAnim()
				animCtrl.DrawArea = grid
				animCtrl.Continue()
				grid.DelAll()
				gridSceneList[sceneId].fnc(grid)
			}
		}
	} else {
		if ch >= 'a' && ch <= 'z' {
			sceneId = int(ch - 'a')
			if sceneId >= 0 && sceneId < len(canvasSceneList) {
				animCtrl.DrawArea = canvas
				canvasSceneList[sceneId].fnc(canvas)
			}
		}
		if ch >= 'A' && ch <= 'Z' {
			sceneId = int(ch - 'A')
			if sceneId >= 0 && sceneId < len(gridSceneList) {
				animCtrl.DrawArea = grid
				gridSceneList[sceneId].fnc(grid)
			}
		}
		animCtrl.Continue()
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	animCtrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

	fmt.Printf("Canvas statistics:\n")
	fmt.Printf("  animation: %v\n", animCtrl.animWatch)
	fmt.Printf("  painting : %v\n", canvas.paintWatch)
	fmt.Printf("  sending  : %v\n", canvas.sendWatch)
	fmt.Printf("Grid statistics:\n")
	fmt.Printf("  animation: %v\n", animCtrl.animWatch)
	fmt.Printf("  painting : %v\n", grid.paintWatch)
	fmt.Printf("  sending  : %v\n", grid.sendWatch)
}
