//go:build ignore
// +build ignore

package ledgrid

import (
	"math"
	"slices"
	"time"
)


// Der Shader verwendet zur Berechnung der darzustellenden Farben
// math. Funktionen. Dazu wird gedanklich ueber das gesamte LedGrid ein
// Koordinatensystem gelegt, welches math. korrekt ist, seinen Ursprung in der
// Mitte des LedGrid hat und so dimensioniert ist, dass der Betrag der
// groessten Koordinaten immer 1.0 ist. Mit Hilfe einer Funktion des Typs
// ShaderFuncType werden dann die Farben berechnet. Die Parameter x und y sind
// Koordinaten im erwaehnten math. Koordinatensystem und mit t wird ein
// Zeitwert (in Sekunden und Bruchteilen davon) an die Funktion uebergeben.
// Der Rueckgabewert ist eine Fliesskommazahl in [0,1] und wird verwendet,
// um aus einer Palette einen Farbwert zu erhalten.

// Jeder Shader basiert auf einer Funktion mit diesem Profil. x und y sind
// Koordinaten des darzustellenden Punktes (siehe Text oben für die
// Dimensionierung des Koord.system), t ist ein fortlaufender Zeitparameter
// und p ist ein Slice von Parametern, die für diesen Shader verwendet werden
// (siehe dazu auch den Typ ShaderParam).
type ShaderFuncType func(x, y, t float64, p []ShaderParam) float64

// Jeder Shader kann über mehrere Parameter gesteuert werden. Ein Parameter
// ist ein Record vom Typ ShaderParam.
type ShaderParam struct {
    // Jeder Parameter hat einen Namen, der beispielsweise in einem GUI oder
    // TUI angezeigt werden kann.
	Name                         string
    // Val ist der aktuelle Wert des Parameters. Im Moment gibt es nur
    // Parameter mit Fliesskommawerten (TO DO: ev. generischen Typ erstellen,
    // der auch andere Werte aufnehmen kann).
	Val                          float64
    // LowerBound und UpperBound sind die Grenzen, in denen der Parameter
    // verändert werden kann und Step ist die Schrittweite, falls das GUI/TUI
    // eine schrittweise Veränderung des Wertes zulässt.
	LowerBound, UpperBound, Step float64
}

type Shader struct {
	VisualizableEmbed
	lg                 *LedGrid
	field              [][]float64
	dPixel, xMin, yMax float64
	fnc                ShaderFuncType
	params             []ShaderParam
    params2            []*Bounded[float64]
	Pal                Colorable
}

func NewShader(lg *LedGrid, shr ShaderRecord, pal Colorable) *Shader {
	s := &Shader{}
	s.VisualizableEmbed.Init("Shader")
	s.lg = lg
	size := lg.Rect.Size()
	s.field = make([][]float64, size.Y)
	for i := range size.Y {
		s.field[i] = make([]float64, size.X)
	}
	s.dPixel = 2.0 / float64(max(size.X, size.Y)-1)
	ratio := float64(size.X) / float64(size.Y)
	if ratio > 1.0 {
		s.xMin = -1.0
		s.yMax = ratio * 1.0
	} else if ratio < 1.0 {
		s.xMin = ratio * -1.0
		s.yMax = 1.0
	} else {
		s.xMin = -1.0
		s.yMax = 1.0
	}
	s.Pal = pal
	s.SetShaderData(shr)
	s.Update(0)
	return s
}

func (s *Shader) SetShaderData(shr ShaderRecord) {
	s.name = shr.n
	s.fnc = shr.f
	s.params = slices.Clone(shr.p)
}

func (s *Shader) ParamList() ([]ShaderParam) {
    return s.params
}

func (s *Shader) Param(name string) float64 {
	for _, param := range s.params {
		if param.Name == name {
			return param.Val
		}
	}
	return 0.0
}

func (s *Shader) SetParam(name string, val float64) {
	for _, param := range s.params {
		if param.Name == name {
			param.Val = val
			return
		}
	}
}

func (s *Shader) Update(dt time.Duration) bool {
	var col, row int
	var x, y float64

	dt = s.AnimatableEmbed.Update(dt)
	y = s.yMax
	for row = range s.field {
		x = s.xMin
		for col = range s.field[row] {
			s.field[row][col] = s.fnc(x, y, s.t0.Seconds(), s.params)
			x += s.dPixel
		}
		y -= s.dPixel
	}
	return true
}

func (s *Shader) Draw() {
	var col, row int

	for row = range s.lg.Bounds().Dy() {
		for col = range s.lg.Bounds().Dx() {
			shaderColor := s.Pal.Color(s.field[row][col])
			s.lg.MixLedColor(col, row, shaderColor, Max)
		}
	}
}


// Im folgenden Abschnitt sind ein paar vordefinierte Shader zusammengestellt.
type ShaderRecord struct {
	n string
	f ShaderFuncType
	p []ShaderParam
}

type ShaderRecord2 struct {
	n string
	f ShaderFuncType
	p []*Bounded[float64]
}


var (
	PlasmaShader2 = ShaderRecord2{
		"Plasma",
		PlasmaShaderFunc,
		[]*Bounded[float64]{
			{Name: "p1", val: 1.2, init: 1.2, lb: 0.0, ub: 10.0, step: 0.1},
			{Name: "p2", val: 1.6, init: 1.6, lb: 0.0, ub: 10.0, step: 0.1},
			{Name: "p3", val: 3.0, init: 3.0, lb: 0.0, ub: 10.0, step: 0.1},
			{Name: "p4", val: 1.5, init: 1.5, lb: 0.0, ub: 10.0, step: 0.1},
			{Name: "p5", val: 5.0, init: 5.0, lb: 0.0, ub: 10.0, step: 0.1},
			{Name: "p6", val: 3.0, init: 5.0, lb: 0.0, ub: 10.0, step: 0.1},
		},
	}
	PlasmaShader = ShaderRecord{
		"Plasma",
		PlasmaShaderFunc,
		[]ShaderParam{
			{Name: "p1", Val: 1.2, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p2", Val: 1.6, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p3", Val: 3.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p4", Val: 1.5, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p5", Val: 5.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
			{Name: "p6", Val: 3.0, LowerBound: 0.0, UpperBound: 10.0, Step: 0.1},
		},
	}
	CircleShader = ShaderRecord{
		"Circle",
		CircleShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
	KaroShader = ShaderRecord{
		"Karo",
		KaroShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 1.0, LowerBound: 0.5, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
	LinearShader = ShaderRecord{
		"Linear",
		LinearShaderFunc,
		[]ShaderParam{
			{Name: "x", Val: 1.0, LowerBound: 0.0, UpperBound: 2.0, Step: 0.1},
			{Name: "y", Val: 0.0, LowerBound: 0.0, UpperBound: 2.0, Step: 0.1},
			{Name: "dir", Val: 1.0, LowerBound: -1.0, UpperBound: 1.0, Step: 2.0},
		},
	}
)

// Die beruehmt/beruechtigte Plasma-Animation. Die Parameter p1 - p6 sind eher
// als Konstanten zu verstehen und eignen sich nicht, um live veraendert
// zu werden.
func PlasmaShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	v1 := f1(x, y, t, p[0].Val)
	v2 := f2(x, y, t, p[1].Val, p[2].Val, p[3].Val)
	v3 := f3(x, y, t, p[4].Val, p[5].Val)
	v := (v1+v2+v3)/6.0 + 0.5
	return v
}
func f1(x, y, t, p1 float64) float64 {
	return math.Sin(x*p1 + t)
}
func f2(x, y, t, p1, p2, p3 float64) float64 {
	return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
}
func f3(x, y, t, p1, p2 float64) float64 {
	cx := 0.125*x + 0.5*math.Sin(t/p1)
	cy := 0.125*y + 0.5*math.Cos(t/p2)
	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
}

// Zeichnet verschachtelte Kreisflaechen. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func CircleShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x = p[0].Val * x / 4.0
	y = p[1].Val * y / 4.0
	t *= p[2].Val / 4.0
	return math.Abs(math.Mod(math.Sqrt(x*x+y*y)-t, 1.0))
}

// func CircleShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf(math.Sqrt(x*x+y*y)/(2.0*math.Sqrt2) + t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }

// Zeichnet verschachtelte Karomuster. Mit p1 kann die Geschw. und die
// Richtung der Anim. beeinflusst werden.
func KaroShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x = p[0].Val * x / 4.0
	y = p[1].Val * y / 4.0
	t *= p[2].Val / 4.0
	return math.Abs(math.Mod(math.Abs(x)+math.Abs(y)-t, 1.0))
}

// func KaroShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf((math.Abs(x)+math.Abs(y))/2.0 + t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }

// Allgemeine Funktion fuer einen animierten Farbverlauf. Mit p1 steuert man
// die Geschwindigkeit der Animation und mit p2/p3 kann festgelegt werden,
// in welche Richtung (x oder y) der Verlauf erfolgen soll.
func LinearShaderFunc(x, y, t float64, p []ShaderParam) float64 {
	x = p[0].Val * (x + 1.0) / 4.0
	y = p[1].Val * (y + 1.0) / 4.0
	t *= p[2].Val / 4.0
	return math.Abs(math.Mod(x+y-t, 1.0))
}

// func LinearShaderFunc(x, y, t float64, p []float64) float64 {
// 	_, f := math.Modf(p[1]*x + p[2]*y + p[0]*t)
// 	if f < 0.0 {
// 		f += 1.0
// 	}
// 	return f
// }
