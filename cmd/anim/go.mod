module github.com/stefan-muehlebach/ledgrid/cmd/anim

go 1.23.0

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/korandiz/v4l v1.1.0
	github.com/stefan-muehlebach/gg v0.0.0-00010101000000-000000000000
	github.com/stefan-muehlebach/ledgrid v0.0.0-00010101000000-000000000000
	gocv.io/x/gocv v0.37.0
	golang.org/x/image v0.19.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.17.0 // indirect
	periph.io/x/conn/v3 v3.7.1 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
