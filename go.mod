module github.com/stefan-muehlebach/ledgrid

go 1.23.0

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/stefan-muehlebach/gg v0.0.0-00010101000000-000000000000
	github.com/vladimirvivien/go4vl v0.0.5
	gocv.io/x/gocv v0.37.0
	golang.org/x/image v0.19.0
	golang.org/x/term v0.23.0
	periph.io/x/conn/v3 v3.7.1
	periph.io/x/host/v3 v3.8.2
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)
