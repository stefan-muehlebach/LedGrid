package ledgrid

import (
	"embed"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/stefan-muehlebach/gg/color"
)

//go:embed data/*.json
var colorFiles embed.FS

type JsonPalette struct {
	ID       int
	Name     string `json:"Title"`
	IsCyclic bool
	Colors   []LedColor
	Stops    []ColorStop
}

// type StopsPalette struct {
// 	ID    int
// 	Name  string
// 	Stops []ColorStop
// }

// type ListPalette struct {
// 	ID       int
// 	Name     string
// 	IsCyclic bool
// 	Colors   []LedColor
// }

func ReadJsonData(fileName string) []JsonPalette {
    var jsonPaletteList []JsonPalette

	data, err := colorFiles.ReadFile(filepath.Join("data", fileName))
	if err != nil {
		log.Fatalf("ReadFile failed: %v", err)
	}
	err = json.Unmarshal(data, &jsonPaletteList)
	if err != nil {
		if err, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("Unmarshal failed in %s: %+v at offset %d", fileName, err, err.Offset)
		} else {
			log.Fatalf("Unmarshal failed in %s: %+v (%T)", fileName, err, err)
		}
	}
    return jsonPaletteList
}

func ReadJsonPalette(fileName string) {
	var colorListJson []JsonPalette

	data, err := colorFiles.ReadFile(filepath.Join("data", fileName))
	if err != nil {
		log.Fatalf("ReadFile failed: %v", err)
	}
	err = json.Unmarshal(data, &colorListJson)
	if err != nil {
		if err, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("Unmarshal failed in %s: %+v at offset %d", fileName, err, err.Offset)
		} else {
			log.Fatalf("Unmarshal failed in %s: %+v (%T)", fileName, err, err)
		}
	}

	// log.Printf("ReadPixelPalettes(): %d entries unmarshalled", len(colorListJson))
	for _, rec := range colorListJson {
		// log.Printf("%+v", rec)
		if len(rec.Colors) > 0 {
			pal := NewGradientPaletteByList(rec.Name, rec.IsCyclic, rec.Colors...)
			PaletteNames = append(PaletteNames, rec.Name)
			PaletteMap[rec.Name] = pal
			// PaletteList = append(PaletteList, pal)
		} else if len(rec.Stops) > 0 {
			pal := NewGradientPalette(rec.Name, rec.Stops...)
			PaletteNames = append(PaletteNames, rec.Name)
			PaletteMap[rec.Name] = pal
			// PaletteList = append(PaletteList, pal)
		} else {
			log.Printf("Palette '%s' has no colors", rec.Name)
		}
	}
}

func ReadNamedColors() {
	for _, colorName := range color.Names {
		ColorNames = append(ColorNames, colorName)
		pal := NewUniformPalette(colorName, LedColorModel.Convert(color.Map[colorName]).(LedColor))
		ColorMap[colorName] = pal
	}
    colorName := "Transparent"
    ColorNames = append(ColorNames, colorName)
    pal := NewUniformPalette(colorName, Transparent)
	ColorMap[colorName] = pal
}

func init() {
	ReadJsonPalette("palettes.json")
	ReadJsonPalette("colourlovers.json")
	ReadNamedColors()
}