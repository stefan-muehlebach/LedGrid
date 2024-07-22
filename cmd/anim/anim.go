package main

import (
	"flag"
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

var (
	width            = 30
	height           = 10
	gridSize         = image.Point{width, height}
	pixelHost        = "raspi-3"
	pixelPort   uint = 5333
	gammaValue       = 3.0
	refreshRate      = 30 * time.Millisecond
	backAlpha        = 1.0

	AnimCtrl Animator
)

//----------------------------------------------------------------------------

func RegularPolygonTest(c *Canvas) {
	posList := []geom.Point{
		ConvertPos(geom.Point{-5.5, 4.5}),
		ConvertPos(geom.Point{34.5, 4.5}),
	}
	posCenter := ConvertPos(geom.Point{14.5, 4.5})
	smallSize := ConvertSize(geom.Point{9.0, 9.0})
	largeSize := ConvertSize(geom.Point{80.0, 80.0})

	polyList := make([]*RegularPolygon, 9)

	aSeq := NewSequence()
	for n := 3; n <= 8; n++ {
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
		animFade := NewColorAnimation(&polyList[n].BorderColor, col.Alpha(0.0), time.Second)

		aGrp := NewGroup(animPos, animAngle)
		aGrp.duration = dur
		aObjSeq := NewSequence(aGrp, animSize, animFade)
		aSeq.Add(aObjSeq)
	}
	aSeq.RepeatCount = AnimationRepeatForever
	aSeq.Start()
}

func GroupTest(ctrl *Canvas) {
	// ctrl.Stop()

	rPos1 := ConvertPos(geom.Point{4.5, 4.5})
	rPos2 := ConvertPos(geom.Point{27.5, 4.5})
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
	// aGroup.Duration = 4*time.Second
	aGroup.RepeatCount = AnimationRepeatForever

	// ctrl.Save("gobs/GroupTest.gob")
	// ctrl.Continue()
	aGroup.Start()
}

func ReadGroupTest(ctrl *Canvas) {
	ctrl.Stop()

	ctrl.Load("gobs/GroupTest.gob")
	ctrl.Continue()

	// fh, err := os.Open("AnimationProgram.gob")
	// if err != nil {
	// 	log.Fatalf("Couldn't create file: %v", err)
	// }
	// gobDecoder := gob.NewDecoder(fh)
	// err = gobDecoder.Decode(&c)
	// if err != nil {
	// 	log.Fatalf("Couldn't decode data: %v", err)
	// }
	// fh.Close()

	// log.Printf("Controller : %+v\n", c)
	// log.Printf("ObjList[0] : (%T) %+v\n", c.ObjList[0], c.ObjList[0])
	// log.Printf("AnimList[0]: (%T) %+v\n", c.AnimList[0], c.AnimList[0])

	// ctrl.Continue()
}

func SequenceTest(ctrl *Canvas) {
	// ctrl.Stop()

	rPos := ConvertPos(geom.Point{14.5, 4.5})
	rSize1 := ConvertSize(geom.Point{29.0, 9.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})

	r := NewRectangle(rPos, rSize1, colornames.SkyBlue)
	ctrl.Add(r)

	aSize1 := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aColor1 := NewColorAnimation(&r.BorderColor, colornames.OrangeRed, time.Second/2)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 2
	aColor2 := NewColorAnimation(&r.BorderColor, colornames.OrangeRed, time.Second/2)
	aSize2 := NewSizeAnimation(&r.Size, rSize1, time.Second)
	aSize2.Cont = true
	aColor3 := NewColorAnimation(&r.BorderColor, colornames.SkyBlue, 3*time.Second/2)
	aColor3.Cont = true

	aSeq := NewSequence(aSize1, aColor1, aColor2, aSize2, aColor3)
	aSeq.duration = 8 * time.Second
	aSeq.RepeatCount = AnimationRepeatForever

	// ctrl.Save("gobs/SequenceTest.gob")
	// ctrl.Continue()
	aSeq.Start()
}

func TimelineTest(ctrl *Canvas) {
	// ctrl.Stop()

	r3Pos1 := ConvertPos(geom.Point{6.5, 4.5})
	r3Size1 := ConvertSize(geom.Point{9.0, 5.0})
	r4Pos1 := ConvertPos(geom.Point{21.5, 4.5})
	r4Size1 := ConvertSize(geom.Point{9.0, 5.0})

	r3 := NewRectangle(r3Pos1, r3Size1, colornames.GreenYellow)
	r4 := NewRectangle(r4Pos1, r4Size1, colornames.SkyBlue)
	ctrl.Add(r3, r4)

	aAngle1 := NewFloatAnimation(&r3.Angle, math.Pi, time.Second)
	aAngle2 := NewFloatAnimation(&r3.Angle, 0.0, time.Second)
	aAngle2.Cont = true

	aAngle3 := NewFloatAnimation(&r4.Angle, -math.Pi, time.Second)
	aAngle4 := NewFloatAnimation(&r4.Angle, 0.0, time.Second)
	aAngle4.Cont = true

	aColor1 := NewColorAnimation(&r3.BorderColor, colornames.OrangeRed, 200*time.Millisecond)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 3
	aColor2 := NewColorAnimation(&r3.BorderColor, colornames.Purple, 500*time.Millisecond)
	aColor3 := NewColorAnimation(&r3.BorderColor, colornames.GreenYellow, 500*time.Millisecond)
	aColor3.Cont = true

	aColor4 := NewColorAnimation(&r4.BorderColor, colornames.OrangeRed, 200*time.Millisecond)
	aColor4.AutoReverse = true
	aColor4.RepeatCount = 3
	aColor5 := NewColorAnimation(&r4.BorderColor, colornames.Purple, 500*time.Millisecond)
	aColor6 := NewColorAnimation(&r4.BorderColor, colornames.SkyBlue, 500*time.Millisecond)
	aColor6.Cont = true

	tl := NewTimeline(5 * time.Second)
	tl.RepeatCount = AnimationRepeatForever
	// Timeline positions for the first rectangle
	tl.Add(300*time.Millisecond, aColor1)
	tl.Add(1800*time.Millisecond, aAngle1)
	tl.Add(2300*time.Millisecond, aColor2)
	tl.Add(2900*time.Millisecond, aAngle2)
	tl.Add(3400*time.Millisecond, aColor3)
	// Timeline positions for the second rectangle
	tl.Add(500*time.Millisecond, aColor4)
	tl.Add(2000*time.Millisecond, aAngle3)
	tl.Add(2500*time.Millisecond, aColor5)
	tl.Add(3100*time.Millisecond, aAngle4)
	tl.Add(3600*time.Millisecond, aColor6)

	// for i, pos := range tl.posList {
	// 	log.Printf("[%d]: %+v", i, pos)
	// }

	// ctrl.Save("gobs/TimelineTest.gob")
	// ctrl.Continue()
	tl.Start()
}

func PathTest(ctrl *Canvas) {
	// ctrl.Stop()
	duration := 4 * time.Second
	pathA := FullCirclePathA
	pathB := FullCirclePathB

	pos1 := ConvertPos(geom.Point{1.5, 4.5})
	pos2 := ConvertPos(geom.Point{14.5, 1.5})
	pos3 := ConvertPos(geom.Point{27.5, 4.5})
	pos4 := ConvertPos(geom.Point{14.5, 7.5})
	cSize := ConvertSize(geom.Point{3.0, 3.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)
	ctrl.Add(c1, c2, c3, c4)

	c1Path := NewPathAnimation(&c1.Pos, pathB, ConvertSize(geom.Point{26.0, 6.0}), duration)
	c1Path.AutoReverse = true
	c3Path := NewPathAnimation(&c3.Pos, pathB, ConvertSize(geom.Point{-26.0, -6.0}), duration)
	c3Path.AutoReverse = true

	c2Path := NewPathAnimation(&c2.Pos, pathA, ConvertSize(geom.Point{10.0, 6.0}), duration)
	c2Path.AutoReverse = true
	c4Path := NewPathAnimation(&c4.Pos, pathA, ConvertSize(geom.Point{-10.0, -6.0}), duration)
	c4Path.AutoReverse = true

	aGrp := NewGroup(c1Path, c2Path, c3Path, c4Path)
	// aGrp.SetDuration(5 * time.Second)
	aGrp.RepeatCount = AnimationRepeatForever

	// ctrl.Save("gobs/PathTest.gob")
	// ctrl.Continue()
	aGrp.Start()
}

func PolygonPathTest(ctrl *Canvas) {
	// ctrl.Stop()

	cPos := ConvertPos(geom.Point{1, 1})
	cSize := ConvertSize(geom.Point{2, 2})

	polyPath1 := NewPolygonPath(
		ConvertPos(geom.Point{1, 1}),
		ConvertPos(geom.Point{28, 1}),
		ConvertPos(geom.Point{28, 8}),
		ConvertPos(geom.Point{1, 8}),

		ConvertPos(geom.Point{1, 2}),
		ConvertPos(geom.Point{27, 2}),
		ConvertPos(geom.Point{27, 7}),
		ConvertPos(geom.Point{2, 7}),

		ConvertPos(geom.Point{2, 3}),
		ConvertPos(geom.Point{26, 3}),
		ConvertPos(geom.Point{26, 6}),
		ConvertPos(geom.Point{3, 6}),

		ConvertPos(geom.Point{3, 4}),
		ConvertPos(geom.Point{25, 4}),
		ConvertPos(geom.Point{25, 5}),
		ConvertPos(geom.Point{4, 5}),
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
	// ctrl.Stop()

	rect := geom.Rectangle{Min: ConvertPos(geom.Point{1.0, 1.0}), Max: ConvertPos(geom.Point{28.0, 8.0})}
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

	// ctrl.Save("gobs/RandomWalk.gob")
	// ctrl.Continue()

	aPos1.Start()
	aPos2.Start()
}

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
	pos1 := ConvertPos(geom.Point{2.0, 2.0})
	pos2 := ConvertPos(geom.Point{8.0, 7.0})
	pos3 := ConvertPos(geom.Point{14.0, 2.0})
	pos4 := ConvertPos(geom.Point{20.0, 7.0})
	pos5 := ConvertPos(geom.Point{26.0, 2.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)
	c5 := NewEllipse(pos5, cSize, colornames.YellowGreen)

	stepRD := ConvertSize(geom.Point{6.0, 5.0})
	stepLU := stepRD.Neg()
	stepRU := ConvertSize(geom.Point{6.0, -5.0})
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
	c1Pos1 := ConvertPos(geom.Point{26.5, 4.5})
	c1Size1 := ConvertSize(geom.Point{10.0, 10.0})
	c1Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c1PosSize := ConvertSize(geom.Point{-24.0, -5.0})
	c2Pos := ConvertPos(geom.Point{2.5, 4.5})
	c2Size1 := ConvertSize(geom.Point{5.0, 5.0})
	c2Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c2PosSize := ConvertSize(geom.Point{24.0, 7.0})

	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.LedColorModel.Convert(colornames.DeepSkyBlue).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Lime).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Teal).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.SkyBlue).(ledgrid.LedColor),
		// ledgrid.ColorMap["DeepSkyBlue"].Color(0),
		// ledgrid.ColorMap["Lime"].Color(0),
		// ledgrid.ColorMap["Teal"].Color(0),
		// ledgrid.ColorMap["SkyBlue"].Color(0),
	)

	c1 := NewEllipse(c1Pos1, c1Size1, colornames.Gold)

	c1pos := NewPathAnimation(&c1.Pos, FullCirclePathB, c1PosSize, 4*time.Second)
	c1pos.RepeatCount = AnimationRepeatForever
	c1pos.Curve = AnimationLinear

	// c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	// c1pos.AutoReverse = true
	// c1pos.RepeatCount = AnimationRepeatForever

	c1size := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1size.AutoReverse = true
	c1size.RepeatCount = AnimationRepeatForever

	c1bcolor := NewColorAnimation(&c1.BorderColor, colornames.OrangeRed, time.Second)
	c1bcolor.AutoReverse = true
	c1bcolor.RepeatCount = AnimationRepeatForever

	// c1fcolor := NewColorAnimation(&c1.FillColor, colornames.OrangeRed.Alpha(0.4), time.Second)
	// c1fcolor.AutoReverse = true
	// c1fcolor.RepeatCount = AnimationRepeatForever

	c2 := NewEllipse(c2Pos, c2Size1, colornames.Lime)

	c2pos := NewPathAnimation(&c2.Pos, FullCirclePathB, c2PosSize, 4*time.Second)
	c2pos.RepeatCount = AnimationRepeatForever
	c2pos.Curve = AnimationLinear

	// c2angle := NewFloatAnimation(&c2.Angle, ConstFloat(2*math.Pi), 4*time.Second)
	// c2angle.RepeatCount = AnimationRepeatForever
	// c2angle.Curve = AnimationLinear

	c2size := NewSizeAnimation(&c2.Size, c2Size2, time.Second)
	c2size.AutoReverse = true
	c2size.RepeatCount = AnimationRepeatForever

	c2color := NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
	c2color.RepeatCount = AnimationRepeatForever
	c2color.Curve = AnimationLinear

	ctrl.Add(c2, c1)

	c1pos.Start()
	c1size.Start()
	c1bcolor.Start()
	// c1fcolor.Start()
	c2pos.Start()
	// c2angle.Start()
	c2size.Start()
	c2color.Start()
}

func CircleAnimation(ctrl *Canvas) {
	c1Pos1 := ConvertPos(geom.Point{1.5, 4.5})
	// c1Pos2 := ConvertPos(geom.Point{14.5, 4.5})
	c1Pos3 := ConvertPos(geom.Point{27.5, 4.5})

	c1Size1 := ConvertSize(geom.Point{3.0, 3.0})
	c1Size2 := ConvertSize(geom.Point{9.0, 9.0})

	// c1Pos1 := ConvertPos(geom.Point{26.5, 4.5})
	// c1Pos2 := ConvertPos(geom.Point{2.5, 4.5})

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
	r1Pos1 := ConvertPos(geom.Point{0.5, 4.5})
	r1Pos2 := ConvertPos(geom.Point{13.5, 4.5})
	r2Pos1 := ConvertPos(geom.Point{28.5, 4.5})
	r2Pos2 := ConvertPos(geom.Point{15.5, 4.5})
	rSize1 := ConvertSize(geom.Point{27.0, 1.0})
	rSize2 := ConvertSize(geom.Point{1.0, 9.0})
	duration := 3 * time.Second / 2

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
	return y*30 + x
}

func pixCoord(idx int) (x, y int) {
	return idx % 30, idx / 30
}

func GlowingPixels(ctrl *Canvas) {
	var pixField []*Pixel

	pixField = make([]*Pixel, 10*30)

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
	// aTimel.Add(40*time.Second, aGrpGray)
	aTimel.RepeatCount = AnimationRepeatForever

	// ctrl.Save("gobs/GlowingPixels.gob")
	// ctrl.Continue()

	aTimel.Start()
}

func MovingText(c *Canvas) {
	tPos := ConvertPos(geom.Point{14.5, 4.5})

	t := NewText(tPos, "Beni", colornames.OrangeRed)

	aAngle := NewFloatAnimation(&t.Angle, 2*math.Pi, 4*time.Second)
	aAngle.RepeatCount = AnimationRepeatForever

	c.Add(t)
	aAngle.Start()
}

func GlowingGridPixels(g *Grid) {
	for y := range g.ledGrid.Rect.Dy() {
		for x := range g.ledGrid.Rect.Dx() {
			pos := image.Point{x, y}
			col := colornames.DimGray.Interpolate(colornames.DarkGrey, rand.Float64())
			pix := NewGridPixel(pos, col)
			g.Add(pix)

        		dur := time.Second + rand.N(400*time.Millisecond)
        		aAlpha := NewAlphaAnimation(&pix.Color.A, 127, dur)
        		aAlpha.AutoReverse = true
        		aAlpha.RepeatCount = AnimationRepeatForever
        		aAlpha.Start()
        }
    }

	aGrpPurple := NewGroup()
	aGrpYellow := NewGroup()
	aGrpGreen := NewGroup()

    for _, obj := range g.ObjList {
        if pix, ok := obj.(*GridPixel); ok {
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
	}

	aTimel := NewTimeline(40 * time.Second)
	aTimel.Add(10*time.Second, aGrpPurple)
	aTimel.Add(20*time.Second, aGrpYellow)
	aTimel.Add(30*time.Second, aGrpGreen)
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
	basePt := image.Point{10, 9}
	baseColor1 := colornames.SkyBlue

    pix := NewGridPixel(basePt, colornames.OrangeRed)
	txt := NewGridText(basePt, baseColor1, " ")
	g.Add(txt, pix)

	go func() {
		for ch := 0x20; ch < 0x7f; ch++ {
			txt.Text = fmt.Sprintf("%c", ch)
			time.Sleep(1 * time.Second)
		}
	}()
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
	var sceneId int
	var input string
	var runInteractive bool

	flag.StringVar(&input, "scene", input, "no menu: direct play")
	flag.BoolVar(&doLog, "log", doLog, "enable logging")
	flag.Parse()

	if len(input) > 0 {
		runInteractive = false
	} else {
		runInteractive = true
	}

	canvasSceneList := []canvasSceneRecord{
		{"(Regular) Polygon test", RegularPolygonTest},
		{"Group test", GroupTest},
		// {"Group test (from saved program)", ReadGroupTest},
		{"Sequence test", SequenceTest},
		{"Timeline test", TimelineTest},
		{"Path test", PathTest},
		{"PolygonPath Test", PolygonPathTest},
		{"Random walk", RandomWalk},
		{"Piiiiixels", Piiiiixels},
		{"Circling Circles", CirclingCircles},
		{"Chasing Circles", ChasingCircles},
		{"Circle Animation", CircleAnimation},
		{"Pushing Rectangles", PushingRectangles},
		{"Glowing Pixels", GlowingPixels},
		{"Moving Text", MovingText},
	}

	gridSceneList := []gridSceneRecord{
		{"Glowing Pixels (Grid)", GlowingGridPixels},
		{"Random Pixels", RandomGridPixels},
		{"Text on a grid", TextOnGrid},
		{"Walking pixel", WalkingPixelOnGrid},
	}

	pixCtrl := ledgrid.NewNetPixelClient(pixelHost, pixelPort)
	pixCtrl.SetGamma(gammaValue, gammaValue, gammaValue)
	pixCtrl.SetMaxBright(255, 255, 255)

    // modLayout := ledgrid.ModuleLayout{
    //     {ledgrid.Module{ledgrid.ModLR, ledgrid.Rot270}},
    //     {ledgrid.Module{ledgrid.ModLR, ledgrid.Rot270}},
    //     {ledgrid.Module{ledgrid.ModRL, ledgrid.Rot000}},
    // }
	ledGrid := ledgrid.NewLedGrid(gridSize, nil)
	canvas := NewCanvas(pixCtrl, ledGrid)
	canvas.Stop()
	grid := NewGrid(pixCtrl, ledGrid)
	grid.Stop()

	if runInteractive {
		sceneId = -1
		for {
			fmt.Printf("Animations:\n")
			fmt.Printf("---------------------------------------\n")
			for i, scene := range canvasSceneList {
				if i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")
			for i, scene := range gridSceneList {
				if i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'A'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")
			fmt.Scanf("%s\n", &input)
			if input[0] == '0' {
				break
			}
			if input[0] >= 'a' && input[0] <= 'z' {
				sceneId = int(input[0] - 'a')
				if sceneId < 0 || sceneId >= len(canvasSceneList) {
					continue
				}
				if AnimCtrl != nil {
					AnimCtrl.Stop()
					AnimCtrl.DelAllAnim()
				}
				AnimCtrl = canvas
				AnimCtrl.Continue()
				canvas.DelAll()
				canvasSceneList[sceneId].fnc(canvas)
			}
			if input[0] >= 'A' && input[0] <= 'Z' {
				sceneId = int(input[0] - 'A')
				if sceneId < 0 || sceneId >= len(gridSceneList) {
					continue
				}
				if AnimCtrl != nil {
					AnimCtrl.Stop()
					AnimCtrl.DelAllAnim()
				}
				AnimCtrl = grid
				AnimCtrl.Continue()
				grid.DelAll()
				gridSceneList[sceneId].fnc(grid)
			}
		}
	} else {
		if input[0] >= 'a' && input[0] <= 'z' {
			sceneId = int(input[0] - 'a')
			if sceneId >= 0 && sceneId < len(canvasSceneList) {
				AnimCtrl = canvas
				canvasSceneList[sceneId].fnc(canvas)
			}
		}
		if input[0] >= 'A' && input[0] <= 'Z' {
			sceneId = int(input[0] - 'A')
			if sceneId >= 0 && sceneId < len(gridSceneList) {
				AnimCtrl = grid
				gridSceneList[sceneId].fnc(grid)
			}
		}
		AnimCtrl.Continue()
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	AnimCtrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

	fmt.Printf("Canvas statistics:\n")
	fmt.Printf("  animation: %v\n", canvas.animWatch)
	fmt.Printf("  painting : %v\n", canvas.paintWatch)
	fmt.Printf("  sending  : %v\n", canvas.sendWatch)
	fmt.Printf("Grid statistics:\n")
	fmt.Printf("  animation: %v\n", grid.animWatch)
	fmt.Printf("  painting : %v\n", grid.paintWatch)
	fmt.Printf("  sending  : %v\n", grid.sendWatch)
}
