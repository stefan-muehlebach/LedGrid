package main

import (
	"encoding/gob"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg/fonts"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

var (
	// Damit wird die Groesse der Queues dimensioniert, welche zu und von
	// den Hintergrundprozessen fuehren.
	queueSize = 400

	// Alle Positionsdaten werden bei der Konvertierung um diesem Wert
	// verschoben. Da gg mit Fliesskommazahlen arbeitet, treffen Koordinaten
	// wie (1,5) nie direkt auf ein Pixel, sondern immer dazwischen.
	displ = geom.Point{0.5, 0.5}

	// Mit oversize wird ein Vergroesserungsfaktor beschrieben, der fuer alle
	// Zeichenoperationen verwendet wird. Damit wird ein insgesamt weicheres
	// Bild erzielt.
	oversize = 10.0

	// Ueber die NewXXX-Funktionen koennen die Objekte einfacher erzeugt
	// werden. Die Fuellfarbe ist gleich der Randfarbe, hat aber einen
	// niedrigeren Alpha-Wert, der mit dieser Konstante definiert werden
	// kann.
	fillAlpha = 0.4

	doLog = false
)

// Ein Canvas ist eine animierbare Zeichenflaeche. Ihr koennen eine beliebige
// Anzahl von zeichenbaren Objekten (Interface CanvasObject) hinzugefuegt
// werden.
type Canvas struct {
	ObjList                          []CanvasObject
	AnimList                         []Animation
	pixCtrl                          ledgrid.PixelClient
	ledGrid                          *ledgrid.LedGrid
	canvas                           *image.RGBA
	gc                               *gg.Context
	scaler                           draw.Scaler
	ticker                           *time.Ticker
	quit                             bool
	animMutex                        *sync.Mutex
	animPit                          time.Time
	logFile                          io.Writer
	animWatch, paintWatch, sendWatch *Stopwatch
}

func NewCanvas(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Canvas {
	var err error

	c := &Canvas{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
	c.canvas = image.NewRGBA(image.Rectangle{Max: c.ledGrid.Rect.Max.Mul(int(oversize))})
	c.gc = gg.NewContextForRGBA(c.canvas)
	c.ObjList = make([]CanvasObject, 0)
	c.scaler = draw.BiLinear.NewScaler(c.ledGrid.Rect.Dx(), c.ledGrid.Rect.Dy(), c.canvas.Rect.Dx(), c.canvas.Rect.Dy())
	c.ticker = time.NewTicker(refreshRate)
	c.AnimList = make([]Animation, 0)
	c.animMutex = &sync.Mutex{}
	if doLog {
		c.logFile, err = os.Create("canvas.log")
		if err != nil {
			log.Fatalf("Couldn't create logfile: %v", err)
		}
	}
	c.animWatch = NewStopwatch()
	c.paintWatch = NewStopwatch()
	c.sendWatch = NewStopwatch()
	go c.backgroundThread()
	return c
}

func (c *Canvas) Close() {
	c.DelAllAnim()
	c.DelAll()
	c.quit = true
}

// Fuegt der Zeichenflaeche weitere Objekte hinzu. Der Zufgriff auf den
// entsprechenden Slice wird nicht synchronisiert.
func (c *Canvas) Add(objs ...CanvasObject) {
	c.ObjList = append(c.ObjList, objs...)
}

// Loescht alle Objekte von der Zeichenflaeche.
func (c *Canvas) DelAll() {
	c.ObjList = c.ObjList[:0]
}

// Fuegt weitere Animationen hinzu. Der Zugriff auf den entsprechenden Slice
// wird synchronisiert, da die Bearbeitung der Animationen durch den
// Background-Thread ebenfalls relativ haeufig auf den Slice zugreift.
func (c *Canvas) AddAnim(anims ...Animation) {
	c.animMutex.Lock()
	c.AnimList = append(c.AnimList, anims...)
	c.animMutex.Unlock()
}

// Loescht alle Animationen.
func (c *Canvas) DelAllAnim() {
	c.animMutex.Lock()
	c.AnimList = c.AnimList[:0]
	c.animMutex.Unlock()
}

// Loescht eine einzelne Animation.
func (c *Canvas) DelAnim(anim Animation) {
	c.animMutex.Lock()
	for i, a := range c.AnimList {
		if a == anim {
			c.AnimList[i] = nil
			return
		}
	}
	c.animMutex.Unlock()
}

// Mit Stop koennen die Animationen und die Darstellung auf der Hardware
// unterbunden werden.
func (c *Canvas) Stop() {
	c.ticker.Stop()
}

// Setzt die Animationen wieder fort.
// TO DO: Die Fortsetzung sollte fuer eine:n Beobachter:in nahtlos erfolgen.
// Im Moment tut es das nicht - man muesste sich bei den Methoden und Ideen
// von AnimationEmbed bedienen.
func (c *Canvas) Continue() {
	c.ticker.Reset(refreshRate)
}

// Mit den folgenden 4 Methoden verfolge ich das ambitionierte Ziel, die
// Animationen in irgendeiner Form serialisierbar zu machen, damit in ferner
// Zukunft die Animationen vollstaendig auf den Rechner des Pixelcontrollers
// verlegt werden koennen und das netzwerkbedingte "Ruckeln" der
// Vergangenheit angehoert.

func (c *Canvas) Save(fileName string) {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	c.Write(fh)
	fh.Close()
}

func (c *Canvas) Write(w io.Writer) {
	gobEncoder := gob.NewEncoder(w)
	err := gobEncoder.Encode(c)
	if err != nil {
		log.Fatalf("Couldn't encode data: %v", err)
	}
}

func (c *Canvas) Load(fileName string) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	c.Read(fh)
	fh.Close()
}

func (c *Canvas) Read(r io.Reader) {
	gobDecoder := gob.NewDecoder(r)
	err := gobDecoder.Decode(c)
	if err != nil {
		log.Fatalf("Couldn't decode data: %v", err)
	}
}

// Hier sind wichtige aber private Methoden, darum in Kleinbuchstaben und
// darum noch sehr wenig Kommentare.
func (c *Canvas) backgroundThread() {
	backColor := colornames.Black
	numCores := runtime.NumCPU()
	animChan := make(chan int, queueSize)
	doneChan := make(chan bool, queueSize)

	for range numCores {
		go c.animationThread(animChan, doneChan)
	}

	lastPit := time.Now()
	for c.animPit = range c.ticker.C {
		if doLog {
			delay := c.animPit.Sub(lastPit)
			lastPit = c.animPit
			fmt.Fprintf(c.logFile, "delay: %v\n", delay)
		}
		if c.quit {
			break
		}
		c.animWatch.Start()
		numAnims := 0
		for i, anim := range c.AnimList {
			if anim == nil || anim.IsStopped() {
				continue
			}
			numAnims++
			animChan <- i
		}
		for range numAnims {
			<-doneChan
		}
		c.animWatch.Stop()

		c.paintWatch.Start()
		c.gc.SetFillColor(backColor)
		c.gc.Clear()
		for _, obj := range c.ObjList {
			obj.Draw(c.gc)
		}
		c.paintWatch.Stop()

		c.sendWatch.Start()
		c.scaler.Scale(c.ledGrid, c.ledGrid.Rect, c.canvas, c.canvas.Rect, draw.Over, nil)
		c.pixCtrl.Draw(c.ledGrid)
		c.sendWatch.Stop()
	}
	close(doneChan)
	close(animChan)
}

func (c *Canvas) animationThread(animChan <-chan int, doneChan chan<- bool) {
	for animId := range animChan {
		c.AnimList[animId].Update(c.animPit)
		// if !c.AnimList[animId].Update(c.animPit) {
		//     c.AnimList[animId] = nil
		// }
		doneChan <- true
	}
}

// Damit werden die jeweiligen Graphik-Objekte beim Package gob registriert,
// um sie binaer zu exportieren.
func init() {
	gob.Register(ledgrid.LedColor{})

	gob.Register(&Ellipse{})
	gob.Register(&Rectangle{})
	gob.Register(&Line{})
	gob.Register(&Pixel{})
}

// Mit ConvertPos muessen alle Positionsdaten konvertiert werden.
func ConvertPos(p geom.Point) geom.Point {
	return p.Add(displ).Mul(oversize)
}

// ConvertSize dagegen wird fuer die Konvertierung aller Groessenangaben
// verwendet.
func ConvertSize(s geom.Point) geom.Point {
	return s.Mul(oversize)
}

// Einzelne Laengen werden mit ConvertLen konvertiert.
func ConvertLen(l float64) float64 {
	return l * oversize
}

type ColorConvertFunc func(ledgrid.LedColor) ledgrid.LedColor

func ApplyAlpha(c ledgrid.LedColor) ledgrid.LedColor {
    alpha := float64(c.A)/255.0
	return c.Alpha(alpha * fillAlpha)
}

// Alle Objekte, die durch den Controller auf dem LED-Grid dargestellt werden
// sollen, muessen im Minimum die Methode Draw implementieren, durch welche
// sie auf einem gg-Kontext gezeichnet werden.
type CanvasObject interface {
	Draw(gc *gg.Context)
}

var (
	defFont     = fonts.Seaford
	defFontSize = ConvertLen(10.0)
)

type Text struct {
	Pos   geom.Point
	Angle float64
	Color ledgrid.LedColor
	Face  font.Face
	text  string
}

func NewText(pos geom.Point, text string, color ledgrid.LedColor) *Text {
	t := &Text{Pos: pos, Color: color, Face: fonts.NewFace(defFont, defFontSize),
		text: text}
	return t
}

func (t *Text) Draw(gc *gg.Context) {
	if t.Angle != 0.0 {
		gc.Push()
		gc.RotateAbout(t.Angle, t.Pos.X, t.Pos.Y)
		defer gc.Pop()
	}
	gc.SetStrokeColor(t.Color)
	gc.SetFontFace(t.Face)
	gc.DrawStringAnchored(t.text, t.Pos.X, t.Pos.Y, 0.5, 0.5)
}

// Mit Ellipse sind alle kreisartigen Objekte abgedeckt. Pos bezeichnet die
// Position des Mittelpunktes und mit Size ist die Breite und Hoehe des
// gesamten Objektes gemeint. Falls ein Rand gezeichnet werden soll, muss
// BorderWith einen Wert >0 enthalten und FillColor, resp. BorderColor
// enthalten die Farben fuer Rand und Flaeche.
type Ellipse struct {
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
}

// Erzeugt eine 'klassische' Ellipse mit einer Randbreite von einem Pixel und
// setzt die Fuellfarbe gleich Randfarbe mit Alpha-Wert von 0.3.
// Will man die einzelnen Werte flexibler verwenden, empfiehlt sich die
// Erzeugung mittels &Ellipse{...}.
func NewEllipse(pos, size geom.Point, borderColor ledgrid.LedColor) *Ellipse {
	e := &Ellipse{Pos: pos, Size: size, BorderWidth: ConvertLen(1.0),
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	e.FillColor = ledgrid.Transparent
	return e
}

func (e *Ellipse) Draw(gc *gg.Context) {
	if e.Angle != 0.0 {
		gc.Push()
		gc.RotateAbout(e.Angle, e.Pos.X, e.Pos.Y)
		defer gc.Pop()
	}
	gc.DrawEllipse(e.Pos.X, e.Pos.Y, e.Size.X/2, e.Size.Y/2)
	gc.SetStrokeWidth(e.BorderWidth)
	gc.SetStrokeColor(e.BorderColor)
	if e.FillColor == ledgrid.Transparent {
		gc.SetFillColor(e.FillColorFnc(e.BorderColor))
	} else {
		gc.SetFillColor(e.FillColor)
	}
	gc.FillStroke()
}

// Rectangle ist fuer alle rechteckigen Objekte vorgesehen. Pos bezeichnet
// den Mittelpunkt des Objektes und Size die Breite, rsep. Hoehe.
type Rectangle struct {
	Pos, Size              geom.Point
	Angle                  float64
	BorderWidth            float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
}

func NewRectangle(pos, size geom.Point, borderColor ledgrid.LedColor) *Rectangle {
	r := &Rectangle{Pos: pos, Size: size, BorderWidth: ConvertLen(1.0),
		BorderColor: borderColor, FillColorFnc: ApplyAlpha}
	return r
}

func (r *Rectangle) Draw(gc *gg.Context) {
	if r.Angle != 0.0 {
		gc.Push()
		gc.RotateAbout(r.Angle, r.Pos.X, r.Pos.Y)
		defer gc.Pop()
	}
	gc.DrawRectangle(r.Pos.X-r.Size.X/2, r.Pos.Y-r.Size.Y/2, r.Size.X, r.Size.Y)
	gc.SetStrokeWidth(r.BorderWidth)
	gc.SetStrokeColor(r.BorderColor)
	if r.FillColorFnc != nil {
		gc.SetFillColor(r.FillColorFnc(r.BorderColor))
	} else {
		gc.SetFillColor(r.FillColor)
	}
	gc.FillStroke()
}

type RegularPolygon struct {
    Pos, Size geom.Point
    Angle float64
    BorderWidth float64
	BorderColor, FillColor ledgrid.LedColor
	FillColorFnc           ColorConvertFunc
    numPoints int
}

func NewRegularPolygon(numPoints int, pos, size geom.Point, borderColor ledgrid.LedColor) *RegularPolygon {
    p := &RegularPolygon{Pos: pos, Size: size, Angle: 0.0, BorderWidth: ConvertLen(1.0),
        BorderColor: borderColor, FillColorFnc: ApplyAlpha, numPoints: numPoints}
    return p
}

func (p *RegularPolygon) Draw(gc *gg.Context) {
    gc.DrawRegularPolygon(p.numPoints, p.Pos.X, p.Pos.Y, p.Size.X/2.0, p.Angle)
    gc.SetStrokeWidth(p.BorderWidth)
    gc.SetStrokeColor(p.BorderColor)
	if p.FillColorFnc != nil {
		gc.SetFillColor(p.FillColorFnc(p.BorderColor))
	} else {
		gc.SetFillColor(p.FillColor)
	}
	gc.FillStroke()
}

// Fuer Geraden ist dieser Datentyp vorgesehen, der von Pos1 nach Pos2
// verlaeuft.
type Line struct {
	Pos1, Pos2 geom.Point
	Width      float64
	Color      ledgrid.LedColor
}

func NewLine(pos1, pos2 geom.Point, col ledgrid.LedColor) *Line {
	l := &Line{Pos1: pos1, Pos2: pos2, Width: ConvertLen(1.0), Color: col}
	return l
}

func (l *Line) Draw(gc *gg.Context) {
	gc.SetStrokeWidth(l.Width)
	gc.SetStrokeColor(l.Color)
	gc.DrawLine(l.Pos1.X, l.Pos1.Y, l.Pos2.X, l.Pos2.Y)
	gc.Stroke()
}

// Will man ein einzelnes Pixel zeichnen, so eignet sich dieser Typ. Er wird
// ueber die Zeichenfunktion DrawPoint im gg-Kontext realisiert und hat einen
// Radius von 0.5*sqrt(2).
type Pixel struct {
	Pos   geom.Point
	Color ledgrid.LedColor
	Alpha float64
}

func NewPixel(pos geom.Point, col ledgrid.LedColor) *Pixel {
	p := &Pixel{Pos: pos, Color: col, Alpha: 1.0}
	return p
}

func (p *Pixel) Draw(gc *gg.Context) {
	gc.SetFillColor(p.Color.Alpha(p.Alpha))
	gc.DrawPoint(p.Pos.X, p.Pos.Y, ConvertLen(0.5*math.Sqrt2))
	gc.Fill()
}


