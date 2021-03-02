module github.com/cima-lexis/lexisdn

go 1.16

require (
	github.com/meteocima/dewetra2wrf v0.0.0-20201207130544-9602d05b70bf
	github.com/meteocima/radar2wrf v1.7.0
	github.com/stretchr/testify v1.6.1 // indirect
)

replace github.com/meteocima/dewetra2wrf => ../dewetra2wrf

replace github.com/meteocima/radar2wrf => ../radar2wrf
