package ledgrid

import (
	"image"
	"image/color"
	"time"
)

// Fuer das gesamte Package gueltige Variablen, resp. Konstanten.
var (
	OutsideColor = BlackColor
)

const (
	defFramesPerSec = 30
)

var (
	framesPerSecond int
	frameRefresh    time.Duration
	frameRefreshMs  int64
	frameRefreshSec float64
)

func init() {
	framesPerSecond = defFramesPerSec
	frameRefresh = time.Second / time.Duration(framesPerSecond)
	frameRefreshMs = frameRefresh.Microseconds()
    // frameRefreshMs = 1000 / framesPerSecond
    frameRefreshSec = frameRefresh.Seconds()
	// frameRefreshSec = float64(frameRefreshMs) / 1000.0
}

type Nameable interface {
    Name() string
}

// Alles, was sich auf dem LedGrid darstellen (d.h. zeichnen laesst),
// implementiert das Interface Drawable.
type Drawable interface {
    // Mit diesen Methoden kann ermittelt, resp. festgelegt werden, ob das
    // Objekt dargestellt (d.h. gezeichnet) werden soll.
	Visible() bool
	SetVisible(v bool)
	// Zeichnet das Objekt auf dem LedGrid.
	// TO DO: die Art, wie die bestehenden mit den neuen Farben gemischt
	// werden sollen, ist aktuell nicht bestimmbar.
	Draw()
}

// Dieses Embedable kann fuer eine Default-Implementation des Drawable-
// Interfaces genutzt werden.
type DrawableEmbed struct {
    visible bool
}

func (d *DrawableEmbed) Init() {
    d.visible = false
}

func (d *DrawableEmbed) Visible() bool {
    return d.visible
}

func (d *DrawableEmbed) SetVisible(visible bool) {
    d.visible = visible
}

// Alles, was sich irgendwie animieren laesst, muss das Interface Animatable
// implementieren.
type Animatable interface {
	// Animiert das Objekt, wobei t ein Point-in-Time in Sekunden und
	// Bruchteilen ist und dt die Zeit in Sekunden seit dem letzten Aufruf.
	// Falls Update false retourniert, ist das Objekt mit der Animation
	// fertig, darf nicht mehr gezeichnet werden und kann vom aufrufenden
	// Programm geloescht werden.
	Update(dt time.Duration) bool
	// Ueber diese beiden Methoden kann festgelegt werden, ob das Objekt
	// animiert werden kann, d.h. auf Update() reagieren soll.
	Alive() bool
	SetAlive(alive bool)
    // Der Speedup-Faktor bestimmt, wie stark sich die Animation auf das
    // Objekt auswirken soll. Es ist ein Faktor, der mit den Werten
    // t und dt der Methode Update() multipliziert wird.
    Speedup() *Bounded[float64]
}

// Dieses Embedable kann fuer eine Default-Implementation des Animatable-
// Interfaces genutzt werden.
type AnimatableEmbed struct {
    alive bool
    t0 time.Duration
    speedup *Bounded[float64]
}

func (a *AnimatableEmbed) Init() {
    a.alive = false
    a.t0 = time.Duration(0)
    a.speedup = NewBounded("speedup", 1.0, 0.1, 2.0, 0.1)
}

func (a *AnimatableEmbed) Update(dt time.Duration) time.Duration {
    if !a.alive {
        return 0
    }
    dt = time.Duration(float64(dt) * a.speedup.val)
    a.t0 += dt
    return dt
}

func (a *AnimatableEmbed) Alive() bool {
    return a.alive
}

func (a *AnimatableEmbed) SetAlive(alive bool) {
    a.alive = alive
}

func (a *AnimatableEmbed) Speedup() *Bounded[float64] {
    return a.speedup
}

// Dieses Interface wird von allen Objekten implementiert, die sowohl animiert,
// als auch dargestellt werden koennen.
type Visualizable interface {
    Drawable
    Animatable
    // Objekte dieser Stufe muessen einen Namen haben.
    Name() string
    SetName(name string)
    // Mit diesen Methoden werden beide Eigenschaften (sichtbar und animierbar
    // zusammen aktiviert, resp. abgefragt.
    Active() bool
    SetActive(active bool)
}

// Dieses Embedable kann fuer eine Default-Implementation des Visualizable-
// Interfaces genutzt werden.
type VisualizableEmbed struct {
    DrawableEmbed
    AnimatableEmbed
    name string
}

func (v *VisualizableEmbed) Init(name string) {
    v.DrawableEmbed.Init()
    v.AnimatableEmbed.Init()
    v.name = name
}

func (v *VisualizableEmbed) Name() (string) {
    return v.name
}

func (v *VisualizableEmbed) SetName(name string) {
    v.name = name
}

func (v *VisualizableEmbed) Active() (bool) {
    return v.Visible() && v.Alive()
}

func (v *VisualizableEmbed) SetActive(active bool) {
    v.SetVisible(active)
    v.SetAlive(active)
}

// Einige der Objekte (wie beispielsweise Shader) koennen zusaetzlich mit
// Parametern gesteuert werden. Damit diese Steuerung so generisch wie
// moeglich ist, haben alle parametrisierbaren Typen dieses Interface zu
// implementieren.
type Parametrizable interface {
    ParamList() ([]*Bounded[float64])
}

// Alles, was im Sinne einer Farbpalette Farben erzeugen kann, implementiert
// das Colorable Interface.
type Colorable interface {
    Nameable
	// Liefert in Abhaengigkeit des Parameters v eine Farbe aus der Palette
	// zurueck. v kann vielfaeltig verwendet werden, bsp. als Parameter im
	// Intervall [0,1] oder als Index (natuerliche Zahl) einer Farbenliste.
	Color(v float64) LedColor
}

// Entspricht dem Bild, welches auf einem LED-Panel angezeigt werden kann.
// Implementiert die Interfaces image.Image und draw.Image, also die Methoden
// ColorModel, Bounds, At und Set.
type LedGrid struct {
	// Groesse des LedGrids. Falls dieses LedGrid Teil eines groesseren
	// Panels sein sollte, dann muss Rect.Min nicht unbedingt {0, 0} sein.
	Rect image.Rectangle
	// Enthaelt die Farbwerte red, green, blue (RGB) fuer jede LED, welche
	// das LedGrid ausmachen. Die Reihenfolge entspricht dabei der
	// Verkabelung, d.h. sie beginnt links oben mit der LED Nr. 0,
	// geht dann nach rechts und auf der zweiten Zeile wieder nach links und
	// so schlangenfoermig weiter.
	Pix []uint8
}

func NewLedGrid(r image.Rectangle) *LedGrid {
	g := &LedGrid{}
	g.Rect = r
	g.Pix = make([]uint8, 3*r.Dx()*r.Dy())
	return g
}

func (g *LedGrid) ColorModel() color.Model {
	return LedColorModel
}

func (g *LedGrid) Bounds() image.Rectangle {
	return g.Rect
}

func (g *LedGrid) At(x, y int) color.Color {
	return g.LedColorAt(x, y)
}

func (g *LedGrid) Set(x, y int, c color.Color) {
	c1 := LedColorModel.Convert(c).(LedColor)
	g.SetLedColor(x, y, c1)
}

// Dient dem schnelleren Zugriff auf den Farbwert einer bestimmten Zelle, resp.
// einer bestimmten LED. Analog zu At(), retourniert den Farbwert jedoch als
// LedColor-Typ.
func (g *LedGrid) LedColorAt(x, y int) LedColor {
	if !(image.Point{x, y}.In(g.Rect)) {
		// log.Printf("LedColorAt(): point outside LedGrid: %d, %d\n", x, y)
		return LedColor{}
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	return LedColor{slc[0], slc[1], slc[2], 0xFF}
}

// Analoge Methode zu Set(), jedoch ohne zeitaufwaendige Konvertierung.
func (g *LedGrid) SetLedColor(x, y int, c LedColor) {
	if !(image.Point{x, y}.In(g.Rect)) {
		// log.Printf("SetLedColor(): point outside LedGrid: %d, %d\n", x, y)
		return
	}
	idx := g.PixOffset(x, y)
	slc := g.Pix[idx : idx+3 : idx+3]
	slc[0] = c.R
	slc[1] = c.G
	slc[2] = c.B
}

func (g *LedGrid) MixLedColor(x, y int, c LedColor, mixType ColorMixType) {
	if !(image.Point{x, y}.In(g.Rect)) {
		// log.Printf("SetLedColor(): point outside LedGrid: %d, %d\n", x, y)
		return
	}
	bgCol := g.LedColorAt(x, y)
	g.SetLedColor(x, y, c.Mix(bgCol, mixType))
}

// Damit wird der Offset eines bestimmten Farbwerts innerhalb des Slices
// Pix berechnet. Dabei wird beruecksichtigt, dass das die LED's im LedGrid
// schlangenfoermig angeordnet sind, wobei die Reihe mit der LED links oben
// beginnt.
func (g *LedGrid) PixOffset(x, y int) int {
	var idx int

	idx = y * g.Rect.Dx()
	if y%2 == 0 {
		idx += x
	} else {
		idx += (g.Rect.Dx() - x - 1)
	}
	return 3 * idx
}

// Hier kommen nun die fuer das LedGrid spezifischen Funktionen.
func (g *LedGrid) Clear(c LedColor) {
	for idx := 0; idx < len(g.Pix); idx += 3 {
		slc := g.Pix[idx : idx+3 : idx+3]
		slc[0] = c.R
		slc[1] = c.G
		slc[2] = c.B
	}
}