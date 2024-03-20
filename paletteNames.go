package ledgrid

var (
	// PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
	// Er wird dynamisch aus den Schluesselwoertern des Maps PaletteMap
	// aufgebaut (siehe Funktion init ganz am Ende des Files).
	PaletteNames []string

	// PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
	// vom Typ Palette.
	PaletteMap = map[string]*Palette{
		"FractalDefault": FractalDefault,
		"EarthAndSky":    EarthAndSky,
		"Neon":           Neon,
		"AAA":            AAA,
		"Adribble":       Adribble,
		"Brazil":         Brazil,
		"BW01":           BW01,
		"BW02":           BW02,
		"BW03":           BW03,
		"Cake":           Cake,
		"Castle":         Castle,
		"CGA":            CGA,
		"Civil":          Civil,
		"Cold00":         Cold00,
		"Cold01":         Cold01,
		"Cold02":         Cold02,
		"Cold03":         Cold03,
		"Cold04":         Cold04,
		"Cold05":         Cold05,
		"ColorBird":      ColorBird,
		"Corn":           Corn,
		"CycleRed":       CycleRed,
		"FadeAll":        FadeAll,
		"FadeBlue":       FadeBlue,
		"FadeGreen":      FadeGreen,
		"FadeRed":        FadeRed,
		"Fame575":        Fame575,
		"Fire":           Fire,
		"Fizz":           Fizz,
		"Fold":           Fold,
		"Hipster":        Hipster,
		"KittyHC":        KittyHC,
		"Kitty":          Kitty,
		"Lantern":        Lantern,
		"LeBron":         LeBron,
		"Lemming":        Lemming,
		"MiamiVice":      MiamiVice,
		"MIUSA":          MIUSA,
		"ML581AT":        ML581AT,
		"NewSeason":      NewSeason,
		"Nightspell":     Nightspell,
		"Rainbows":       Rainbows,
		"Rasta":          Rasta,
		"RGB":            RGB,
		"Simpson":        Simpson,
		"Smurf":          Smurf,
		"Spring":         Spring,
		"Sunset":         Sunset,
		"Violet":         Violet,
		"Warmth":         Warmth,
		"Wayyou":         Wayyou,
		"Weparted":       Weparted,
	}

	// In diesem Block werden die Paletten konkret erstellt. Im Moment
	// koennen so nur Paletten mit aequidistanten Farbstuetzstellen
	// erzeugt werden.
	FractalDefault = NewPaletteWithColors(FractalDefaultColors)
	EarthAndSky    = NewPaletteWithColors(EarthAndSkyColors)
	Neon           = NewPaletteWithColors(NeonColors)
	AAA            = NewPaletteWithColors(AAAColors)
	Adribble       = NewPaletteWithColors(AdribbleColors)
	Brazil         = NewPaletteWithColors(BrazilColors)
	BW01           = NewPaletteWithColors(BW01Colors)
	BW02           = NewPaletteWithColors(BW02Colors)
	BW03           = NewPaletteWithColors(BW03Colors)
	Cake           = NewPaletteWithColors(CakeColors)
	Castle         = NewPaletteWithColors(CastleColors)
	CGA            = NewPaletteWithColors(CGAColors)
	Civil          = NewPaletteWithColors(CivilColors)
	Cold00         = NewPaletteWithColors(Cold00Colors)
	Cold01         = NewPaletteWithColors(Cold01Colors)
	Cold02         = NewPaletteWithColors(Cold02Colors)
	Cold03         = NewPaletteWithColors(Cold03Colors)
	Cold04         = NewPaletteWithColors(Cold04Colors)
	Cold05         = NewPaletteWithColors(Cold05Colors)
	ColorBird      = NewPaletteWithColors(ColorBirdColors)
	Corn           = NewPaletteWithColors(CornColors)
	CycleRed       = NewPaletteWithColors(CycleRedColors)
	FadeAll        = NewPaletteWithColors(FadeAllColors)
	FadeBlue       = NewPaletteWithColors(FadeBlueColors)
	FadeGreen      = NewPaletteWithColors(FadeGreenColors)
	FadeRed        = NewPaletteWithColors(FadeRedColors)
	Fame575        = NewPaletteWithColors(Fame575Colors)
	Fire           = NewPaletteWithColors(FireColors)
	Fizz           = NewPaletteWithColors(FizzColors)
	Fold           = NewPaletteWithColors(FoldColors)
	Hipster        = NewPaletteWithColors(HipsterColors)
	KittyHC        = NewPaletteWithColors(KittyHCColors)
	Kitty          = NewPaletteWithColors(KittyColors)
	Lantern        = NewPaletteWithColors(LanternColors)
	LeBron         = NewPaletteWithColors(LeBronColors)
	Lemming        = NewPaletteWithColors(LemmingColors)
	MiamiVice      = NewPaletteWithColors(MiamiViceColors)
	MIUSA          = NewPaletteWithColors(MIUSAColors)
	ML581AT        = NewPaletteWithColors(ML581ATColors)
	NewSeason      = NewPaletteWithColors(NewSeasonColors)
	Nightspell     = NewPaletteWithColors(NightspellColors)
	Rainbows       = NewPaletteWithColors(RainbowsColors)
	Rasta          = NewPaletteWithColors(RastaColors)
	RGB            = NewPaletteWithColors(RGBColors)
	Simpson        = NewPaletteWithColors(SimpsonColors)
	Smurf          = NewPaletteWithColors(SmurfColors)
	Spring         = NewPaletteWithColors(SpringColors)
	Sunset         = NewPaletteWithColors(SunsetColors)
	Violet         = NewPaletteWithColors(VioletColors)
	Warmth         = NewPaletteWithColors(WarmthColors)
	Wayyou         = NewPaletteWithColors(WayyouColors)
	Weparted       = NewPaletteWithColors(WepartedColors)

	// Schliesslich folgen die Farblisten, welche fuer die einzelnen Paletten
	// verwendet werden.
	FractalDefaultColors = []LedColor{
		{0x01, 0x06, 0x62, 0xFF},
		{0xDA, 0xF9, 0xFE, 0xFF},
		{0xFE, 0xBC, 0x08, 0xFF},
		{0xDA, 0xF9, 0xFE, 0xFF},
		{0x01, 0x06, 0x62, 0xFF},
	}

	EarthAndSkyColors = []LedColor{
		{0xff, 0xff, 0xff, 0xFF},
		{0xff, 0xff, 0x00, 0xFF},
		{0xff, 0x33, 0x00, 0xFF},
		{0x00, 0x00, 0x99, 0xFF},
		{0x00, 0x66, 0xff, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}

	NeonColors = []LedColor{
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x99, 0xcc, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0x00, 0x99, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}

	AAAColors = []LedColor{
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0xff, 0xFF},
		{0x00, 0x80, 0xff, 0xFF},
		{0x80, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
	}
	AdribbleColors = []LedColor{
		{0x3D, 0x4C, 0x53, 0xFF},
		{0x70, 0xB7, 0xBA, 0xFF},
		{0xF1, 0x43, 0x3F, 0xFF},
		{0xE7, 0xE1, 0xD4, 0xFF},
		{0xFF, 0xFF, 0xFF, 0xFF},
	}
	BrazilColors = []LedColor{
		{0x00, 0x8c, 0x53, 0xFF},
		{0x2e, 0x00, 0xe4, 0xFF},
		{0xdf, 0xea, 0x00, 0xFF},
	}
	BW01Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	BW02Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	BW03Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	CakeColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xBD, 0x41, 0x41, 0xFF},
		{0xD9, 0x75, 0x70, 0xFF},
		{0xF2, 0xE8, 0xDF, 0xFF},
		{0xB2, 0xCF, 0xC0, 0xFF},
		{0x71, 0x9D, 0x98, 0xFF},
	}
	CastleColors = []LedColor{
		{0x4B, 0x34, 0x5C, 0xFF},
		{0x94, 0x62, 0x82, 0xFF},
		{0xE5, 0xA1, 0x9B, 0xFF},
	}
	CGAColors = []LedColor{
		{0xd3, 0x51, 0x7d, 0xFF},
		{0x15, 0xa0, 0xbf, 0xFF},
		{0xff, 0xc0, 0x62, 0xFF},
	}
	CivilColors = []LedColor{
		{0x36, 0x2F, 0x2D, 0xFF},
		{0x4C, 0x4C, 0x4C, 0xFF},
		{0x94, 0xB7, 0x3E, 0xFF},
		{0xB5, 0xC0, 0xAF, 0xFF},
		{0xFA, 0xFD, 0xF2, 0xFF},
	}
	Cold00Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
	}
	Cold01Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
	}
	Cold02Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
		{0x00, 0x80, 0x80, 0xFF},
	}
	Cold03Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
		{0x00, 0x80, 0x80, 0xFF},
		{0x80, 0xff, 0xff, 0xFF},
	}
	Cold04Colors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
		{0x00, 0x80, 0x80, 0xFF},
		{0x80, 0xff, 0xff, 0xFF},
		{0x80, 0x80, 0xff, 0xFF},
	}
	Cold05Colors = []LedColor{
		{0x80, 0xff, 0xff, 0xFF},
		{0x00, 0x80, 0xff, 0xFF},
		{0xc0, 0xc0, 0xc0, 0xFF},
		{0x80, 0xa3, 0xff, 0xFF},
		{0x58, 0x58, 0x58, 0xFF},
		{0xb9, 0xf4, 0xf9, 0xFF},
	}
	ColorBirdColors = []LedColor{
		{0x1F, 0xA6, 0x98, 0xFF},
		{0xC5, 0xD9, 0x32, 0xFF},
		{0xF2, 0x59, 0x22, 0xFF},
		{0x40, 0x1E, 0x11, 0xFF},
		{0xD7, 0x19, 0x5A, 0xFF},
	}
	CornColors = []LedColor{
		{0x29, 0x23, 0x1F, 0xFF},
		{0xEB, 0xE1, 0xCC, 0xFF},
		{0xDB, 0x9B, 0x1A, 0xFF},
	}
	FadeAllColors = []LedColor{
		{0xff, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0x00, 0xFF},
		{0x00, 0xff, 0x00, 0xFF},
		{0x00, 0xff, 0xff, 0xFF},
		{0x00, 0x00, 0xff, 0xFF},
		{0xff, 0x00, 0xff, 0xFF},
	}
	FadeRedColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0x00, 0x00, 0xFF},
	}
	FadeGreenColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0x00, 0xFF},
	}
	FadeBlueColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0x00, 0xff, 0xFF},
	}
	CycleRedColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0x00, 0x00, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
	}
	Fame575Colors = []LedColor{
		{0x54, 0x0c, 0x0d, 0xFF},
		{0xfb, 0x74, 0x23, 0xFF},
		{0xf9, 0xf4, 0x8e, 0xFF},
		{0x41, 0x76, 0xc4, 0xFF},
		{0x5a, 0xaf, 0x2e, 0xFF},
	}
	FireColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0x00, 0x40, 0xFF},
		{0xff, 0x00, 0x00, 0xFF},
		{0xff, 0xff, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	FizzColors = []LedColor{
		{0x58, 0x8F, 0x27, 0xFF},
		{0x04, 0xBF, 0xBF, 0xFF},
		{0xF7, 0xE9, 0x67, 0xFF},
	}
	FoldColors = []LedColor{
		{0x2A, 0x03, 0x08, 0xFF},
		{0x92, 0x4F, 0x1B, 0xFF},
		{0xE2, 0xAC, 0x3F, 0xFF},
		{0xF8, 0xED, 0xC6, 0xFF},
		{0x7B, 0xA5, 0x8D, 0xFF},
	}
	HipsterColors = []LedColor{
		{0x8D, 0x0E, 0x6B, 0xFF},
		{0x05, 0x34, 0x34, 0xFF},
		{0x00, 0x80, 0x80, 0xFF},
		{0x4D, 0x42, 0x08, 0xFF},
		{0xCD, 0xAD, 0x0A, 0xFF},
	}
	KittyHCColors = []LedColor{
		{0xc7, 0x56, 0xa7, 0xFF},
		{0xe0, 0xdd, 0x00, 0xFF},
		{0xc9, 0xcd, 0xd0, 0xFF},
	}
	KittyColors = []LedColor{
		{0x9f, 0x45, 0x6b, 0xFF},
		{0x4f, 0x7a, 0x9a, 0xFF},
		{0xe6, 0xc8, 0x4c, 0xFF},
	}
	LanternColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x0D, 0x9A, 0x0D, 0xFF},
		{0xFF, 0xFF, 0xFF, 0xFF},
	}
	// LanternColors = []LedColor{
	// 	LedColor{0x04, 0x04, 0x02, 0xFF},
	// 	LedColor{0xBA, 0xE7, 0xAE, 0xFF},
	// 	LedColor{0xFC, 0xBD, 0x00, 0xFF},
	// 	LedColor{0x3E, 0xA8, 0x8E, 0xFF},
	// 	LedColor{0x22, 0x22, 0x01, 0xFF},
	// }
	LeBronColors = []LedColor{
		{0x3e, 0x3e, 0x3e, 0xFF},
		{0xd4, 0xb6, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	LemmingColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x00, 0x00, 0xff, 0xFF},
		{0x00, 0xff, 0x00, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	MiamiViceColors = []LedColor{
		{0x1b, 0xe3, 0xff, 0xFF},
		{0xff, 0x82, 0xdc, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	MIUSAColors = []LedColor{
		{0x50, 0x4b, 0x46, 0xFF},
		{0x1a, 0x3c, 0x53, 0xFF},
		{0xa0, 0x00, 0x28, 0xFF},
	}
	ML581ATColors = []LedColor{
		{0x69, 0x96, 0x55, 0xFF},
		{0xf2, 0x6a, 0x36, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	NewSeasonColors = []LedColor{
		{0x4C, 0x39, 0x33, 0xFF},
		{0x00, 0x5B, 0x5B, 0xFF},
		{0xFE, 0x55, 0x02, 0xFF},
		{0xFE, 0xBF, 0x51, 0xFF},
		{0xF8, 0xEA, 0xB2, 0xFF},
	}
	NightspellColors = []LedColor{
		{0xED, 0xEE, 0xBA, 0xFF},
		{0xFE, 0xA8, 0x1C, 0xFF},
		{0xB2, 0x01, 0x52, 0xFF},
		{0x4B, 0x0B, 0x44, 0xFF},
		{0x24, 0x0F, 0x37, 0xFF},
	}
	VioletColors = []LedColor{
		{0x4B, 0x0B, 0x44, 0xFF},
		{0x4B, 0x0B, 0x44, 0xFF},
	}
	RainbowsColors = []LedColor{
		{0x49, 0x2D, 0x61, 0xFF},
		{0x04, 0x80, 0x91, 0xFF},
		{0x61, 0xC1, 0x55, 0xFF},
		{0xF2, 0xD4, 0x3F, 0xFF},
		{0xD1, 0x02, 0x6C, 0xFF},
	}
	RastaColors = []LedColor{
		{0xdc, 0x32, 0x3c, 0xFF},
		{0xf0, 0xcb, 0x58, 0xFF},
		{0x3c, 0x82, 0x5e, 0xFF},
	}
	RGBColors = []LedColor{
		{0xff, 0x00, 0x00, 0xFF},
		{0x00, 0xff, 0x00, 0xFF},
		{0x00, 0x00, 0xff, 0xFF},
	}
	SimpsonColors = []LedColor{
		{0xd9, 0xc2, 0x3e, 0xFF},
		{0xa9, 0x6a, 0x95, 0xFF},
		{0x7d, 0x95, 0x4b, 0xFF},
		{0x4b, 0x39, 0x6b, 0xFF},
	}
	SmurfColors = []LedColor{
		{0x1d, 0x16, 0x28, 0xFF},
		{0x44, 0xbd, 0xf4, 0xFF},
		{0xe3, 0x1e, 0x3a, 0xFF},
		{0xe8, 0xb1, 0x18, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}
	SpringColors = []LedColor{
		{0x1D, 0x19, 0x29, 0xFF},
		{0xFE, 0x53, 0x24, 0xFF},
		{0xA8, 0x14, 0x3A, 0xFF},
		{0x66, 0xA5, 0x95, 0xFF},
		{0xCF, 0xBD, 0x81, 0xFF},
	}
	SunsetColors = []LedColor{
		{0x00, 0x00, 0x00, 0xFF},
		{0x35, 0x24, 0x38, 0xFF},
		{0x61, 0x29, 0x49, 0xFF},
		{0x99, 0x28, 0x39, 0xFF},
		{0xED, 0x61, 0x37, 0xFF},
		{0xD6, 0xAB, 0x5C, 0xFF},
	}
	WarmthColors = []LedColor{
		{0xFC, 0xEB, 0xB6, 0xFF},
		{0x5E, 0x41, 0x2F, 0xFF},
		{0xF0, 0x78, 0x18, 0xFF},
		{0x78, 0xC0, 0xA8, 0xFF},
		{0xF0, 0xA8, 0x30, 0xFF},
	}
	WayyouColors = []LedColor{
		{0x1C, 0x21, 0x30, 0xFF},
		{0x02, 0x8F, 0x76, 0xFF},
		{0xB3, 0xE0, 0x99, 0xFF},
		{0xFF, 0xEA, 0xAD, 0xFF},
		{0xD1, 0x43, 0x34, 0xFF},
	}
	WepartedColors = []LedColor{
		{0x02, 0x7B, 0x7F, 0xFF},
		{0xFF, 0xA5, 0x88, 0xFF},
		{0xD6, 0x29, 0x57, 0xFF},
		{0xBF, 0x1E, 0x62, 0xFF},
		{0x57, 0x2E, 0x4F, 0xFF},
	}
)

// Die init-Funktion wird aktuell dafuer verwendet, beim Programmstart
// den Slice PaletteNames aufzubauen.
func init() {
	PaletteNames = make([]string, 0)
	for name, _ := range PaletteMap {
		PaletteNames = append(PaletteNames, name)
	}
}
