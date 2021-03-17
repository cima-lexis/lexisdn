module github.com/cima-lexis/lexisdn

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/meteocima/dewetra2wrf v1.0.1
	github.com/meteocima/radar2wrf v1.7.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

// replace github.com/meteocima/dewetra2wrf => ../dewetra2wrf
// replace github.com/meteocima/radar2wrf => ../radar2wrf
