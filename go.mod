module github.com/cima-lexis/lexisdn

go 1.16

require (
	github.com/cweill/gotests v1.6.0 // indirect
	github.com/go-delve/delve v1.6.0 // indirect
	github.com/haya14busa/goplay v1.0.0 // indirect
	github.com/meteocima/dewetra2wrf v0.0.0-20201207130544-9602d05b70bf
	github.com/meteocima/radar2wrf v1.7.0
	github.com/pkg/profile v0.0.0-20170413231811-06b906832ed0 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/tools/gopls v0.6.6 // indirect
)

replace github.com/meteocima/dewetra2wrf => ../dewetra2wrf

replace github.com/meteocima/radar2wrf => ../radar2wrf
