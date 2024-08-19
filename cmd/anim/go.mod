module ledgrid/cmd/anim

go 1.23.0

replace (
	github.com/stefan-muehlebach/gg => ../../../gg
	github.com/stefan-muehlebach/ledgrid => ../..
)

require (
	github.com/stefan-muehlebach/gg v1.3.4
	github.com/stefan-muehlebach/ledgrid v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.19.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/vladimirvivien/go4vl v0.0.5 // indirect
	gocv.io/x/gocv v0.37.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	periph.io/x/conn/v3 v3.7.1 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
